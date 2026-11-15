package usecases

import (
	"context"

	"github.com/climblive/platform/backend/internal/domain"
)

type HealthUseCase struct {
	ScoreEngineManager domain.StatusReporter
	ResultKeeper       domain.StatusReporter
	Scrubber           domain.StatusReporter
	PointValueKeeper   domain.StatusReporter
}

func (uc *HealthUseCase) GetHealth(_ context.Context) ([]domain.ServiceStatus, error) {
	return []domain.ServiceStatus{
		uc.ScoreEngineManager.GetStatus(),
		uc.ResultKeeper.GetStatus(),
		uc.Scrubber.GetStatus(),
		uc.PointValueKeeper.GetStatus(),
	}, nil
}
