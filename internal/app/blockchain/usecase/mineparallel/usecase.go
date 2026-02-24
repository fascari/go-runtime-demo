package mineparallel

import (
	"context"
	"runtime"

	"go-runtime-demo/internal/app/blockchain/domain"
)

type (
	UseCase struct {
		blockchain *domain.Blockchain
	}

	Result struct {
		Blocks        []domain.Block `json:"blocks"`
		Duration      string         `json:"duration"`
		Goroutines    int            `json:"goroutines"`
		TotalBlocks   int            `json:"total_blocks"`
		GCRuns        uint32         `json:"gc_runs"`
		GCPauseMs     float64        `json:"gc_pause_ms"`
		HeapDeltaMB   float64        `json:"heap_delta_mb"`
		GCCPUFraction float64        `json:"gc_cpu_fraction"`
	}
)

func New(blockchain *domain.Blockchain) UseCase {
	return UseCase{
		blockchain: blockchain,
	}
}

func (uc UseCase) Execute(_ context.Context, data string, numGoroutines int) Result {
	var memBefore, memAfter runtime.MemStats
	runtime.ReadMemStats(&memBefore)

	blocks, duration := uc.blockchain.MineParallel(data, numGoroutines)

	runtime.ReadMemStats(&memAfter)

	return Result{
		Blocks:        blocks,
		Duration:      duration.String(),
		Goroutines:    numGoroutines,
		TotalBlocks:   uc.blockchain.Length(),
		GCRuns:        memAfter.NumGC - memBefore.NumGC,
		GCPauseMs:     float64(memAfter.PauseTotalNs-memBefore.PauseTotalNs) / 1e6,
		HeapDeltaMB:   float64(int64(memAfter.HeapAlloc)-int64(memBefore.HeapAlloc)) / 1024 / 1024,
		GCCPUFraction: memAfter.GCCPUFraction,
	}
}
