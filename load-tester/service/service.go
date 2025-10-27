package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"load-tester/adapter/go_gateway_adapter"
	"load-tester/adapter/py_gateway_adapter"

	"github.com/google/uuid"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Custom error types for timeout and dropped requests
var (
	ErrTimeout = fmt.Errorf("request timed out")
	ErrDropped = fmt.Errorf("request dropped by queue")
)

// service adapter
type adapter struct {
	goGatewayAdapter *go_gateway_adapter.Adapter
	pyGatewayAdapter *py_gateway_adapter.Adapter
}

// service
type Service struct {
	logger *logrus.Logger

	adapter *adapter

	metricsLock sync.Mutex
	proc        *process.Process // Current process info for resource monitoring
}

func NewService(
	logger *logrus.Logger,
	goGatewayAdapter *go_gateway_adapter.Adapter,
	pyGatewayAdapter *py_gateway_adapter.Adapter,
) *Service {
	proc, _ := process.NewProcess(int32(os.Getpid()))

	return &Service{
		logger: logger,

		adapter: &adapter{
			goGatewayAdapter: goGatewayAdapter,
			pyGatewayAdapter: pyGatewayAdapter,
		},

		proc: proc,
	}
}

// monitorResources periodically collects resource usage metrics
func (service *Service) monitorResources(ctx context.Context, metrics *Metrics) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// CPU Usage
			if cpuPercent, err := service.proc.CPUPercent(); err == nil {
				service.metricsLock.Lock()
				metrics.CPUUsage = append(metrics.CPUUsage, cpuPercent)
				service.metricsLock.Unlock()
			}

			// Memory Usage
			if memInfo, err := service.proc.MemoryInfo(); err == nil {
				service.metricsLock.Lock()
				metrics.MemoryUsage = append(metrics.MemoryUsage, memInfo.RSS)
				service.metricsLock.Unlock()
			}

			// Force GC to get accurate memory stats
			runtime.GC()
		}
	}
}

// preparePayload creates a copy of the payload
func (service *Service) preparePayload(param *BaseParam) (map[string]any, error) {
	now := time.Now()
	timestamp := now.Format("20060102150405")
	nanoID := fmt.Sprintf("%d", now.UnixNano()%100000000000)

	// Start with the required fields
	payload := map[string]any{}

	if param.Protocol == "bl2" {
		// TCP protocol format for py-core
		payload["id_message"] = timestamp + nanoID + uuid.New().String()
		payload["operation"] = "get_account_by_account_number"
		payload["params"] = map[string]any{
			"account_number": param.Payload.AccountNumber,
		}
	} else {
		// gRPC protocol - just the account number
		payload["account_number"] = param.Payload.AccountNumber
	}

	return payload, nil
}

// executeRequest executes a request based on the protocol
func (service *Service) executeRequest(ctx context.Context, param *BaseParam) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		// Create payload with unique system_trace_audit
		payload, err := service.preparePayload(param)
		if err != nil {
			return fmt.Errorf("failed to prepare payload: %v", err)
		}
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to marshal payload: %v", err)
		}

		// Execute request based on protocol
		switch param.Protocol {
		case "bl2":
			respMap, err := service.adapter.pyGatewayAdapter.SendRequest(payloadBytes)
			if err != nil {
				return err
			}

			// Check err_info for timeout and dropped cases
			if errInfo, ok := respMap["err_info"].(string); ok {
				if errInfo == "timeout" {
					return ErrTimeout
				}
				if errInfo == "dropped" {
					return ErrDropped
				}
			}

			// Check response status
			if status, ok := respMap["status"].(string); ok && status != "000" {
				responseCode := "unknown"
				responseMsg := "unknown error"

				if code, ok := respMap["response_code"].(string); ok {
					responseCode = code
				}
				if msg, ok := respMap["response_msg"].(string); ok {
					responseMsg = msg
				}

				return fmt.Errorf("request failed with code %s: %s", responseCode, responseMsg)
			}

			return nil

		case "grpc":
			_, err := service.adapter.goGatewayAdapter.GetAccountByAccountNumber(context.TODO(), &go_gateway_adapter.GetAccountByAccountNumberParams{
				AccountNumber: payload["account_number"].(string),
			})
			if err != nil {
				// Check for specific gRPC error codes
				if st, ok := status.FromError(err); ok {
					switch st.Code() {
					case codes.DeadlineExceeded:
						return ErrTimeout
					case codes.Aborted:
						return ErrDropped
					}
				}
				return err
			}

			return nil

		default:
			return fmt.Errorf("unsupported protocol: %s", param.Protocol)
		}
	}
}
