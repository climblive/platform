package usecases

import (
	"context"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

type statsUseCaseRepository interface {
	GetGlobalStats(ctx context.Context) (domain.GlobalStats, error)
}

type StatsUseCase struct {
	Repo statsUseCaseRepository
}

func (uc *StatsUseCase) GetGlobalStats(ctx context.Context) (domain.GlobalStats, error) {
	stats, err := uc.Repo.GetGlobalStats(ctx)
	if err != nil {
		return domain.GlobalStats{}, errors.Wrap(err, 0)
	}

	return stats, nil
}
