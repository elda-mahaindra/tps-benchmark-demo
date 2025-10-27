package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"load-tester/util/errs"

	"github.com/sirupsen/logrus"
)

type LoadBurstParam struct {
	BaseParam
	TotalReqs int `json:"total_reqs"` // Total number of requests to send
}

type LoadBurstResult = TestResult

// Burst testing, no limiter, burst all to goroutine
func (service *Service) LoadBurst(ctx context.Context, param *LoadBurstParam) (map[string]any, error) {
	const op errs.Op = "service/LoadBurst"

	service.logger.WithFields(logrus.Fields{
		"op":    op,
		"param": param,
	}).Info("Starting burst load test")

	metrics := &Metrics{
		Latencies:   make([]time.Duration, 0, param.TotalReqs),
		StartTime:   time.Now(),
		CPUUsage:    make([]float64, 0),
		MemoryUsage: make([]uint64, 0),
		Errors:      make(map[string]int),
	}

	// Create a wait group to track all goroutines
	var wg sync.WaitGroup
	wg.Add(param.TotalReqs)

	// Create channels for collecting metrics
	latencyChan := make(chan time.Duration, param.TotalReqs)
	errorChan := make(chan error, param.TotalReqs)

	// Start resource monitoring
	monitorCtx, cancelMonitor := context.WithCancel(ctx)
	defer cancelMonitor()
	go service.monitorResources(monitorCtx, metrics)

	// Create a done channel to signal completion
	done := make(chan struct{})

	// Start a goroutine to collect results
	go func() {
		for {
			select {
			case latency, ok := <-latencyChan:
				if !ok {
					return
				}
				metrics.Latencies = append(metrics.Latencies, latency)
			case err, ok := <-errorChan:
				if !ok {
					return
				}
				errStr := err.Error()
				metrics.Errors[errStr]++

				// Categorize errors
				switch err {
				case ErrTimeout:
					metrics.TimeoutRequests.Add(1)
				case ErrDropped:
					metrics.DroppedRequests.Add(1)
				default:
					metrics.FailedRequests.Add(1)
				}
			case <-done:
				return
			}
		}
	}()

	// Launch goroutines for each request
	for i := 0; i < param.TotalReqs; i++ {
		go func() {
			defer wg.Done()

			start := time.Now()
			err := service.executeRequest(ctx, &param.BaseParam)
			latency := time.Since(start)

			latencyChan <- latency
			if err != nil {
				errorChan <- err
			} else {
				metrics.SuccessfulRequests.Add(1)
			}
			metrics.TotalRequests.Add(1)
		}()
	}

	// Wait for all requests to complete
	wg.Wait()

	// Signal collection goroutine to stop and close channels
	close(latencyChan)
	close(errorChan)
	close(done)

	metrics.EndTime = time.Now()

	// Convert metrics to test result
	result := &LoadBurstResult{
		Summary: Summary{
			TestDuration:       metrics.EndTime.Sub(metrics.StartTime).String(),
			TotalRequests:      fmt.Sprintf("%d", metrics.TotalRequests.Load()),
			SuccessfulRequests: fmt.Sprintf("%d", metrics.SuccessfulRequests.Load()),
			FailedRequests:     fmt.Sprintf("%d", metrics.FailedRequests.Load()),
			TimeoutRequests:    fmt.Sprintf("%d", metrics.TimeoutRequests.Load()),
			DroppedRequests:    fmt.Sprintf("%d", metrics.DroppedRequests.Load()),
			AverageRPS:         fmt.Sprintf("%.2f", float64(metrics.TotalRequests.Load())/metrics.EndTime.Sub(metrics.StartTime).Seconds()),
			P95LatencyMs:       fmt.Sprintf("%.2f", calculatePercentileLatency(metrics.Latencies, 95)),
			P99LatencyMs:       fmt.Sprintf("%.2f", calculatePercentileLatency(metrics.Latencies, 99)),
		},
		ResourceMetrics: ResourceMetric{
			CPU: CPU{
				AverageUsage: calculateAverage(metrics.CPUUsage),
				PeakUsage:    calculatePeak(metrics.CPUUsage),
				UsagePattern: metrics.CPUUsage,
			},
			Memory: Memory{
				AverageMB: float64(calculateAverage(convertToFloat64(metrics.MemoryUsage))) / 1024 / 1024,
				PeakMB:    float64(calculatePeak(convertToFloat64(metrics.MemoryUsage))) / 1024 / 1024,
			},
		},
		Errors: metrics.Errors,
	}

	return map[string]any{
		"result": result,
	}, nil
}
