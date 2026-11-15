package usecases

import (
	"context"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
	"uuid"
)

type scoreEngineManager interface {
	GetScoreEngine(ctx context.Context, instanceID domain.ScoreEngineInstanceID) (domain.ScoreEngineDescriptor, error)
	ListScoreEngines(ctx context.Context) ([]domain.ScoreEngineDescriptor, error)
	ListScoreEnginesByCompetition(ctx context.Context, competitionID domain.CompetitionID) ([]domain.ScoreEngineDescriptor, error)
	StopScoreEngine(ctx context.Context, instanceID domain.ScoreEngineInstanceID) error
	StartScoreEngine(ctx context.Context, competitionID domain.CompetitionID, terminatedBy time.Time) (domain.ScoreEngineInstanceID, error)
}

type scoreEngineUseCaseRepository interface {
	GetCompetition(ctx context.Context, tx domain.Transaction, competitionID domain.CompetitionID) (domain.Competition, error)
}

type ScoreEngineUseCase struct {
	Authorizer         domain.Authorizer
	Repo               scoreEngineUseCaseRepository
	ScoreEngineManager scoreEngineManager
}

func (uc *ScoreEngineUseCase) ListScoreEngines(ctx context.Context) ([]domain.ScoreEngineDescriptor, error) {
	var role domain.AuthRole
	var err error

	if role, err = uc.Authorizer.HasOwnership(ctx, domain.OwnershipData{}); err != nil {
		return nil, errors.Wrap(err, 0)
	}

	if role != domain.AdminRole {
		return nil, domain.ErrNotAuthorized
	}

	engines, err := uc.ScoreEngineManager.ListScoreEngines(ctx)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	return engines, nil
}

func (uc *ScoreEngineUseCase) ListScoreEnginesByCompetition(ctx context.Context, competitionID domain.CompetitionID) ([]domain.ScoreEngineDescriptor, error) {
	competition, err := uc.Repo.GetCompetition(ctx, nil, competitionID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, competition.Ownership); err != nil {
		return nil, errors.Wrap(err, 0)
	}

	engines, err := uc.ScoreEngineManager.ListScoreEnginesByCompetition(ctx, competitionID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	return engines, nil
}

func (uc *ScoreEngineUseCase) StopScoreEngine(ctx context.Context, instanceID domain.ScoreEngineInstanceID) error {
	engine, err := uc.ScoreEngineManager.GetScoreEngine(ctx, instanceID)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	competition, err := uc.Repo.GetCompetition(ctx, nil, engine.CompetitionID)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, competition.Ownership); err != nil {
		return errors.Wrap(err, 0)
	}

	err = uc.ScoreEngineManager.StopScoreEngine(ctx, instanceID)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	return nil
}

func (uc *ScoreEngineUseCase) StartScoreEngine(ctx context.Context, competitionID domain.CompetitionID, terminatedBy time.Time) (domain.ScoreEngineInstanceID, error) {
	competition, err := uc.Repo.GetCompetition(ctx, nil, competitionID)
	if err != nil {
		return uuid.Nil(), errors.Wrap(err, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, competition.Ownership); err != nil {
		return uuid.Nil(), errors.Wrap(err, 0)
	}

	if competition.TimeBegin.IsZero() || competition.TimeEnd.IsZero() {
		return uuid.Nil(), errors.Wrap(domain.ErrNotAllowed, 0)
	}

	now := time.Now()

	if now.Before(competition.TimeBegin.Add(-1 * time.Hour)) {
		return uuid.Nil(), errors.Wrap(domain.ErrNotAllowed, 0)
	}

	if terminatedBy.Before(now) {
		return uuid.Nil(), errors.Wrap(domain.ErrNotAllowed, 0)
	}

	if terminatedBy.Sub(now) > 12*time.Hour {
		return uuid.Nil(), errors.Wrap(domain.ErrNotAllowed, 0)
	}

	instanceID, err := uc.ScoreEngineManager.StartScoreEngine(ctx, competitionID, terminatedBy)
	if err != nil {
		return uuid.Nil(), errors.Wrap(err, 0)
	}

	return instanceID, nil
}
