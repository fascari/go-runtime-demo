package addblock

import (
	"context"
	"runtime"
	"time"

	"go-runtime-demo/internal/app/blockchain/domain"
)

type (
	UseCase struct {
		blockchain *domain.Blockchain
	}

	Result struct {
		Block         domain.Block `json:"block"`
		Duration      string       `json:"duration"`
		GCRuns        uint32       `json:"gc_runs"`
		GCPauseMs     float64      `json:"gc_pause_ms"`
		HeapDeltaMB   float64      `json:"heap_delta_mb"`
		HeapObjects   uint64       `json:"heap_objects"`
		GCCPUFraction float64      `json:"gc_cpu_fraction"`
	}
)

func New(blockchain *domain.Blockchain) UseCase {
	return UseCase{
		blockchain: blockchain,
	}
}

func (uc UseCase) Execute(_ context.Context, data string) Result {
	var memBefore, memAfter runtime.MemStats
	runtime.ReadMemStats(&memBefore)

	start := time.Now()
	block := uc.blockchain.AddBlock(data)
	duration := time.Since(start)

	runtime.ReadMemStats(&memAfter)

	return Result{
		Block:         block,
		Duration:      duration.String(),
		GCRuns:        memAfter.NumGC - memBefore.NumGC,
		GCPauseMs:     float64(memAfter.PauseTotalNs-memBefore.PauseTotalNs) / 1e6,
		HeapDeltaMB:   float64(int64(memAfter.HeapAlloc)-int64(memBefore.HeapAlloc)) / 1024 / 1024,
		HeapObjects:   memAfter.HeapObjects,
		GCCPUFraction: memAfter.GCCPUFraction,
	}
}
