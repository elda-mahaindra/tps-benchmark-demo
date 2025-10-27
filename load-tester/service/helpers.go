package service

import (
	"math"
	"sort"
	"time"
)

// calculatePercentile calculates the nth percentile from a slice of durations
func calculatePercentile(values []time.Duration, percentile float64) time.Duration {
	if len(values) == 0 {
		return 0
	}

	// Create a copy to avoid modifying the original slice
	sorted := make([]time.Duration, len(values))
	copy(sorted, values)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})

	index := int(float64(len(sorted)-1) * (percentile / 100))
	return sorted[index]
}

// calculateAverage calculates the average of float64 values
func calculateAverage(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	var sum float64
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

// calculatePeak returns the highest value in the slice
func calculatePeak(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	peak := values[0]
	for _, v := range values[1:] {
		if v > peak {
			peak = v
		}
	}
	return peak
}

// calculatePeakUint64 returns the highest uint64 value
func calculatePeakUint64(values []uint64) uint64 {
	if len(values) == 0 {
		return 0
	}

	peak := values[0]
	for _, v := range values[1:] {
		if v > peak {
			peak = v
		}
	}
	return peak
}

// calculatePeakInt returns the highest int value
func calculatePeakInt(values []int) int {
	if len(values) == 0 {
		return 0
	}

	peak := values[0]
	for _, v := range values[1:] {
		if v > peak {
			peak = v
		}
	}
	return peak
}

// convertToFloat64 converts uint64 slice to float64 slice
func convertToFloat64(values []uint64) []float64 {
	result := make([]float64, len(values))
	for i, v := range values {
		result[i] = float64(v)
	}
	return result
}

// Helper function for min value
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func calculatePercentileLatency(latencies []time.Duration, percentile float64) float64 {
	if len(latencies) == 0 {
		return 0
	}

	// Convert to milliseconds and sort
	ms := make([]float64, len(latencies))
	for i, lat := range latencies {
		ms[i] = float64(lat.Milliseconds())
	}
	sort.Float64s(ms)

	// Calculate percentile index
	index := int(math.Ceil((percentile/100)*float64(len(ms)))) - 1
	if index < 0 {
		index = 0
	}

	return ms[index]
}
