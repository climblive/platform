package usecases

import (
	"context"

	"github.com/climblive/platform/backend/internal/domain"
)

type HealthUseCase struct {
	ScoreEngineManager domain.StatusReporter
	ScoreKeeper        domain.StatusReporter
	Scrubber           domain.StatusReporter
}

func (uc *HealthUseCase) GetHealth(_ context.Context) ([]domain.ServiceStatus, error) {
	return []domain.ServiceStatus{
		uc.ScoreEngineManager.GetStatus(),
		uc.ScoreKeeper.GetStatus(),
		uc.Scrubber.GetStatus(),
	}, nil
}
