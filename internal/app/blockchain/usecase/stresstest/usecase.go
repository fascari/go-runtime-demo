package stresstest

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"runtime"
	"sync"
	"time"
)

type UseCase struct{}

type Result struct {
	Duration      string  `json:"duration"`
	Goroutines    int     `json:"goroutines"`
	Allocations   int     `json:"allocations"`
	GCCollections uint32  `json:"gc_collections"`
	MemoryDeltaMB float64 `json:"memory_delta_mb"`
	FinalAllocMB  float64 `json:"final_alloc_mb"`
	NumGoroutines int     `json:"num_goroutines"`
}

func New() UseCase {
	return UseCase{}
}

func (uc UseCase) Execute(ctx context.Context, allocations, goroutines int) Result {
	start := time.Now()

	var memBefore, memAfter runtime.MemStats
	runtime.ReadMemStats(&memBefore)

	var wg sync.WaitGroup

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			for j := 0; j < allocations; j++ {
				data := make([]byte, 1024*1024)
				data[0] = byte(id)

				hash := sha256.Sum256(data)
				_ = hex.EncodeToString(hash[:])

				if j%10 == 0 {
					runtime.Gosched()
				}
			}
		}(i)
	}

	wg.Wait()
	runtime.ReadMemStats(&memAfter)

	duration := time.Since(start)

	return Result{
		Duration:      duration.String(),
		Goroutines:    goroutines,
		Allocations:   allocations,
		GCCollections: memAfter.NumGC - memBefore.NumGC,
		MemoryDeltaMB: float64(int64(memAfter.Alloc)-int64(memBefore.Alloc)) / 1024 / 1024,
		FinalAllocMB:  float64(memAfter.Alloc) / 1024 / 1024,
		NumGoroutines: runtime.NumGoroutine(),
	}
}

