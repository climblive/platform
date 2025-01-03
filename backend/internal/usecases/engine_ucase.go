package usecases

import (
	"context"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
	"github.com/google/uuid"
)

type scoreEngineManager interface {
	ListScoreEnginesByContest(ctx context.Context, contestID domain.ContestID) ([]domain.ScoreEngineInstanceID, error)
	StopScoreEngine(ctx context.Context, instanceID domain.ScoreEngineInstanceID) error
	StartScoreEngine(ctx context.Context, contestID domain.ContestID) (domain.ScoreEngineInstanceID, error)
	ReverseLoopupScoreEngine(ctx context.Context, instanceID domain.ScoreEngineInstanceID) (domain.ContestID, error)
}

type scoreEngineUseCaseRepository interface {
	GetContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) (domain.Contest, error)
}

type ScoreEngineUseCase struct {
	Authorizer         domain.Authorizer
	Repo               scoreEngineUseCaseRepository
	ScoreEngineManager scoreEngineManager
}

func (uc *ScoreEngineUseCase) ListScoreEnginesByContest(ctx context.Context, contestID domain.ContestID) ([]domain.ScoreEngineInstanceID, error) {
	contest, err := uc.Repo.GetContest(ctx, nil, contestID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, contest.Ownership); err != nil {
		return nil, errors.Wrap(err, 0)
	}

	return uc.ScoreEngineManager.ListScoreEnginesByContest(ctx, contestID)
}

func (uc *ScoreEngineUseCase) StopScoreEngine(ctx context.Context, instanceID domain.ScoreEngineInstanceID) error {
	contestID, err := uc.ScoreEngineManager.ReverseLoopupScoreEngine(ctx, instanceID)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	contest, err := uc.Repo.GetContest(ctx, nil, contestID)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, contest.Ownership); err != nil {
		return errors.Wrap(err, 0)
	}

	return uc.ScoreEngineManager.StopScoreEngine(ctx, instanceID)
}

func (uc *ScoreEngineUseCase) StartScoreEngine(ctx context.Context, contestID domain.ContestID) (domain.ScoreEngineInstanceID, error) {
	contest, err := uc.Repo.GetContest(ctx, nil, contestID)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, contest.Ownership); err != nil {
		return uuid.Nil, errors.Wrap(err, 0)
	}

	return uc.ScoreEngineManager.StartScoreEngine(ctx, contestID)
}
