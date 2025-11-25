package stats

import (
	"context"

	"go-runtime-demo/internal/app/monitoring/domain"
)

type UseCase struct {
	monitor *domain.Monitor
}

func New(monitor *domain.Monitor) UseCase {
	return UseCase{
		monitor: monitor,
	}
}

func (uc UseCase) Execute(_ context.Context) domain.RuntimeStats {
	return uc.monitor.RuntimeStats()
}
