package listblocks

import (
	"context"

	"go-runtime-demo/internal/app/blockchain/domain"
)

type UseCase struct {
	blockchain *domain.Blockchain
}

func New(blockchain *domain.Blockchain) UseCase {
	return UseCase{
		blockchain: blockchain,
	}
}

func (uc UseCase) Execute(ctx context.Context) []domain.Block {
	return uc.blockchain.Chain()
}
