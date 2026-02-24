package stresstest

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"runtime"
	"sync"
	"time"
)

type (
	UseCase struct{}

	// AllocationPattern defines how memory is allocated during stress test
	AllocationPattern string

	Result struct {
		Duration         string  `json:"duration"`
		Goroutines       int     `json:"goroutines"`
		Allocations      int     `json:"allocations"`
		Pattern          string  `json:"pattern"`
		GCCollections    uint32  `json:"gc_collections"`
		GCPauseTotalMs   float64 `json:"gc_pause_total_ms"`
		GCCPUFraction    float64 `json:"gc_cpu_fraction"`
		HeapObjectsDelta int64   `json:"heap_objects_delta"`
		MemoryDeltaMB    float64 `json:"memory_delta_mb"`
		FinalAllocMB     float64 `json:"final_alloc_mb"`
		NumGoroutines    int     `json:"num_goroutines"`
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

func (uc UseCase) Execute(_ context.Context, allocations, goroutines int, pattern AllocationPattern) Result {
	start := time.Now()

	var memBefore, memAfter runtime.MemStats
	runtime.ReadMemStats(&memBefore)

	// For long-lived pattern, keep references alive
	var longLived [][]byte
	var mu sync.Mutex

	var wg sync.WaitGroup

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			var localLongLived [][]byte

			for j := 0; j < allocations; j++ {
				var data []byte

				switch pattern {
				case PatternShortLived:
					// Short-lived: allocate and let GC collect immediately
					data = make([]byte, 1024*1024)
					data[0] = byte(id)
					hash := sha256.Sum256(data)
					_ = hex.EncodeToString(hash[:])

				case PatternLongLived:
					// Long-lived: keep references to prevent GC
					data = make([]byte, 1024*1024)
					data[0] = byte(id)
					localLongLived = append(localLongLived, data)

				case PatternMixed:
					// Mixed: 50% short-lived, 50% long-lived
					data = make([]byte, 1024*1024)
					data[0] = byte(id)
					if j%2 == 0 {
						localLongLived = append(localLongLived, data)
					} else {
						hash := sha256.Sum256(data)
						_ = hex.EncodeToString(hash[:])
					}
				}

				if j%10 == 0 {
					runtime.Gosched()
				}
			}

			// Transfer long-lived data to shared slice
			if len(localLongLived) > 0 {
				mu.Lock()
				longLived = append(longLived, localLongLived...)
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()
	runtime.ReadMemStats(&memAfter)

	duration := time.Since(start)
	totalAllocMB := float64(allocations * goroutines) // Each allocation is 1MB

	return Result{
		Duration:         duration.String(),
		Goroutines:       goroutines,
		Allocations:      allocations,
		Pattern:          string(pattern),
		GCCollections:    memAfter.NumGC - memBefore.NumGC,
		GCPauseTotalMs:   float64(memAfter.PauseTotalNs-memBefore.PauseTotalNs) / 1e6,
		GCCPUFraction:    memAfter.GCCPUFraction,
		HeapObjectsDelta: int64(memAfter.HeapObjects - memBefore.HeapObjects),
		MemoryDeltaMB:    float64(int64(memAfter.Alloc)-int64(memBefore.Alloc)) / 1024 / 1024,
		FinalAllocMB:     float64(memAfter.Alloc) / 1024 / 1024,
		NumGoroutines:    runtime.NumGoroutine(),
		ThroughputMBps:   totalAllocMB / duration.Seconds(),
	}
}
