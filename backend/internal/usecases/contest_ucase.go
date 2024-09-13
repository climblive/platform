package usecases

import (
	"context"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

type contestUseCaseRepository interface {
	domain.Transactor

	GetContest(ctx context.Context, tx domain.Transaction, contestID domain.ResourceID) (domain.Contest, error)
}

type ContestUseCase struct {
	Repo contestUseCaseRepository
}

func (uc *ContestUseCase) GetContest(ctx context.Context, contestID domain.ResourceID) (domain.Contest, error) {
	contest, err := uc.Repo.GetContest(ctx, nil, contestID)
	if err != nil {
		return domain.Contest{}, errors.Wrap(err, 0)
	}

	return contest, nil
}

func (uc *ContestUseCase) GetContestsByOrganizer(ctx context.Context, organizerID domain.ResourceID) ([]domain.Contest, error) {
	panic("not implemented")
}

func (uc *ContestUseCase) UpdateContest(ctx context.Context, contestID domain.ResourceID, contest domain.Contest) (domain.Contest, error) {
	panic("not implemented")
}

func (uc *ContestUseCase) DeleteContest(ctx context.Context, contestID domain.ResourceID) error {
	panic("not implemented")
}

func (uc *ContestUseCase) DuplicateContest(ctx context.Context, contestID domain.ResourceID) (domain.Contest, error) {
	panic("not implemented")
}

func (uc *ContestUseCase) CreateContest(ctx context.Context, organizerID domain.ResourceID, contest domain.Contest) (domain.Contest, error) {
	panic("not implemented")
}

func (uc *ContestUseCase) GetScores(ctx context.Context, contestID domain.ResourceID) ([]domain.Score, error) {
	panic("not implemented")
}
