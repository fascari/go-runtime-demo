package gcbenchmark

import (
	"context"
	"runtime"
	"time"
)

type (
	UseCase struct{}

	// AllocationPattern defines how memory is allocated
	AllocationPattern string

	Input struct {
		Allocations int               `json:"allocations"`
		SizeKB      int               `json:"size_kb"`
		Pattern     AllocationPattern `json:"pattern"`
	}

	Result struct {
		DurationMs       float64 `json:"duration_ms"`
		AllocBeforeMB    float64 `json:"alloc_before_mb"`
		AllocAfterMB     float64 `json:"alloc_after_mb"`
		HeapObjectsDelta int64   `json:"heap_objects_delta"`
		GCRuns           uint32  `json:"gc_runs"`
		GCPauseTotalMs   float64 `json:"gc_pause_total_ms"`
		GCCPUFraction    float64 `json:"gc_cpu_fraction"`
		ThroughputMBps   float64 `json:"throughput_mb_per_sec"`
	}
)

const (
	PatternShortLived AllocationPattern = "short-lived"
	PatternLongLived  AllocationPattern = "long-lived"
	PatternMixed      AllocationPattern = "mixed"
)

func New() UseCase {
	return UseCase{}
}

func (uc UseCase) Execute(_ context.Context, input Input) Result {
	if input.Allocations <= 0 {
		input.Allocations = 10000
	}
	if input.SizeKB <= 0 {
		input.SizeKB = 1
	}
	if input.Pattern == "" {
		input.Pattern = PatternShortLived
	}

	var memBefore, memAfter runtime.MemStats
	runtime.ReadMemStats(&memBefore)

	start := time.Now()

	// Keep references for long-lived pattern
	var longLived [][]byte
	sizeBytes := input.SizeKB * 1024

	for i := 0; i < input.Allocations; i++ {
		data := make([]byte, sizeBytes)

		switch input.Pattern {
		case PatternLongLived:
			longLived = append(longLived, data)
		case PatternMixed:
			if i%2 == 0 {
				longLived = append(longLived, data)
			}
			// PatternShortLived: data goes out of scope and can be collected
		}

		// Simulate some work
		for j := range data {
			data[j] = byte(i + j)
		}
	}

	// Force a GC to see impact
	runtime.GC()

	duration := time.Since(start)
	runtime.ReadMemStats(&memAfter)

	totalAllocMB := float64(input.Allocations*input.SizeKB) / 1024

	return Result{
		DurationMs:       float64(duration.Microseconds()) / 1000,
		AllocBeforeMB:    float64(memBefore.Alloc) / 1024 / 1024,
		AllocAfterMB:     float64(memAfter.Alloc) / 1024 / 1024,
		HeapObjectsDelta: int64(memAfter.HeapObjects - memBefore.HeapObjects),
		GCRuns:           memAfter.NumGC - memBefore.NumGC,
		GCPauseTotalMs:   float64(memAfter.PauseTotalNs-memBefore.PauseTotalNs) / 1e6,
		GCCPUFraction:    memAfter.GCCPUFraction,
		ThroughputMBps:   totalAllocMB / duration.Seconds(),
	}
}
