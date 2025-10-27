package service

import (
	"sync/atomic"
	"time"
)

// Metrics holds the test execution metrics
type Metrics struct {
	TotalRequests      atomic.Int64
	SuccessfulRequests atomic.Int64
	FailedRequests     atomic.Int64
	TimeoutRequests    atomic.Int64 // Requests that timed out
	DroppedRequests    atomic.Int64 // Requests rejected by queue
	Latencies          []time.Duration
	StartTime          time.Time
	EndTime            time.Time
	CPUUsage           []float64
	MemoryUsage        []uint64
	Errors             map[string]int
}

type BaseParam struct {
	ServiceName string  `json:"service_name"`
	Protocol    string  `json:"protocol"` // "bl2" or "grpc"
	Payload     Payload `json:"payload"`
}

type Payload struct {
	AccountNumber string `json:"account_number"`
}

type Summary struct {
	TestDuration       string `json:"test_duration"`
	TotalRequests      string `json:"total_requests"`
	SuccessfulRequests string `json:"successful_requests"`
	FailedRequests     string `json:"failed_requests"`
	TimeoutRequests    string `json:"timeout_requests"` // Requests that timed out
	DroppedRequests    string `json:"dropped_requests"` // Requests rejected by queue
	AverageRPS         string `json:"average_rps"`
	P95LatencyMs       string `json:"p95_latency_ms"`
	P99LatencyMs       string `json:"p99_latency_ms"`
}

type CPU struct {
	AverageUsage float64   `json:"average_usage"`
	PeakUsage    float64   `json:"peak_usage"`
	UsagePattern []float64 `json:"usage_pattern"`
}

type Memory struct {
	AverageMB float64 `json:"average_mb"`
	PeakMB    float64 `json:"peak_mb"`
}

type ResourceMetric struct {
	CPU    CPU    `json:"cpu"`
	Memory Memory `json:"memory"`
}

type TestResult struct {
	Summary         Summary        `json:"summary"`
	ResourceMetrics ResourceMetric `json:"resource_metric"`
	Errors          map[string]int `json:"errors"`
}
