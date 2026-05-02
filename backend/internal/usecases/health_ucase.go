package usecases

import (
	"context"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

type HealthUseCase struct {
	Authorizer         domain.Authorizer
	ScoreEngineManager domain.StatusReporter
	ScoreKeeper        domain.StatusReporter
	Scrubber           domain.StatusReporter
}

func (uc *HealthUseCase) GetHealth(ctx context.Context) (domain.HealthStatus, error) {
	role, err := uc.Authorizer.HasOwnership(ctx, domain.OwnershipData{})
	if err != nil {
		return domain.HealthStatus{}, errors.Wrap(err, 0)
	}

	if !role.OneOf(domain.AdminRole) {
		return domain.HealthStatus{}, errors.Wrap(domain.ErrInsufficientRole, 0)
	}

	return domain.HealthStatus{
		ScoreEngineManager: uc.ScoreEngineManager.GetStatus(),
		ScoreKeeper:        uc.ScoreKeeper.GetStatus(),
		Scrubber:           uc.Scrubber.GetStatus(),
	}, nil
}
