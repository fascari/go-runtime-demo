package addblock

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

func (uc UseCase) Execute(ctx context.Context, data string) domain.Block {
	return uc.blockchain.AddBlock(data)
}
