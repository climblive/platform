package usecases

import (
	"context"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

type problemUseCaseRepository interface {
	domain.Transactor

	GetProblemsByContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) ([]domain.Problem, error)
}

type ProblemUseCase struct {
	Repo problemUseCaseRepository
}

func (uc *ProblemUseCase) GetProblemsByContest(ctx context.Context, contestID domain.ContestID) ([]domain.Problem, error) {
	problems, err := uc.Repo.GetProblemsByContest(ctx, nil, contestID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	return problems, nil
}
