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
	ListScoreEnginesByContest(ctx context.Context, contestID domain.ContestID) ([]domain.ScoreEngineDescriptor, error)
	StopScoreEngine(ctx context.Context, instanceID domain.ScoreEngineInstanceID) error
	StartScoreEngine(ctx context.Context, contestID domain.ContestID, terminatedBy time.Time) (domain.ScoreEngineInstanceID, error)
}

type scoreEngineUseCaseRepository interface {
	GetContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) (domain.Contest, error)
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

func (uc *ScoreEngineUseCase) ListScoreEnginesByContest(ctx context.Context, contestID domain.ContestID) ([]domain.ScoreEngineDescriptor, error) {
	contest, err := uc.Repo.GetContest(ctx, nil, contestID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, contest.Ownership); err != nil {
		return nil, errors.Wrap(err, 0)
	}

	engines, err := uc.ScoreEngineManager.ListScoreEnginesByContest(ctx, contestID)
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

	contest, err := uc.Repo.GetContest(ctx, nil, engine.ContestID)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, contest.Ownership); err != nil {
		return errors.Wrap(err, 0)
	}

	err = uc.ScoreEngineManager.StopScoreEngine(ctx, instanceID)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	return nil
}

func (uc *ScoreEngineUseCase) StartScoreEngine(ctx context.Context, contestID domain.ContestID, terminatedBy time.Time) (domain.ScoreEngineInstanceID, error) {
	contest, err := uc.Repo.GetContest(ctx, nil, contestID)
	if err != nil {
		return uuid.Nil(), errors.Wrap(err, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, contest.Ownership); err != nil {
		return uuid.Nil(), errors.Wrap(err, 0)
	}

	if contest.TimeBegin.IsZero() || contest.TimeEnd.IsZero() {
		return uuid.Nil(), errors.Wrap(domain.ErrNotAllowed, 0)
	}

	now := time.Now()

	if now.Before(contest.TimeBegin.Add(-1 * time.Hour)) {
		return uuid.Nil(), errors.Wrap(domain.ErrNotAllowed, 0)
	}

	if terminatedBy.Before(now) {
		return uuid.Nil(), errors.Wrap(domain.ErrNotAllowed, 0)
	}

	if terminatedBy.Sub(now) > 12*time.Hour {
		return uuid.Nil(), errors.Wrap(domain.ErrNotAllowed, 0)
	}

	instanceID, err := uc.ScoreEngineManager.StartScoreEngine(ctx, contestID, terminatedBy)
	if err != nil {
		return uuid.Nil(), errors.Wrap(err, 0)
	}

	return instanceID, nil
}
