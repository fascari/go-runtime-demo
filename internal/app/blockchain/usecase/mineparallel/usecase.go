package mineparallel

import (
	"context"

	"go-runtime-demo/internal/app/blockchain/domain"
)

type UseCase struct {
	blockchain *domain.Blockchain
}

type Result struct {
	Blocks      []domain.Block `json:"blocks"`
	Duration    string         `json:"duration"`
	Goroutines  int            `json:"goroutines"`
	TotalBlocks int            `json:"total_blocks"`
}

func New(blockchain *domain.Blockchain) UseCase {
	return UseCase{
		blockchain: blockchain,
	}
}

func (uc UseCase) Execute(ctx context.Context, data string, numGoroutines int) Result {
	blocks, duration := uc.blockchain.MineParallel(data, numGoroutines)

	return Result{
		Blocks:      blocks,
		Duration:    duration.String(),
		Goroutines:  numGoroutines,
		TotalBlocks: uc.blockchain.Length(),
	}
}
