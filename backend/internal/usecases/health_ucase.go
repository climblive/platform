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

func (uc *HealthUseCase) GetHealth(_ context.Context) (domain.HealthStatus, error) {
	return domain.HealthStatus{
		ScoreEngineManager: uc.ScoreEngineManager.GetStatus(),
		ScoreKeeper:        uc.ScoreKeeper.GetStatus(),
		Scrubber:           uc.Scrubber.GetStatus(),
	}, nil
}
