package usecases

import (
	"context"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

type problemUseCaseRepository interface {
	domain.Transactor

	GetProblemsByContest(ctx context.Context, tx domain.Transaction, contestID domain.ResourceID) ([]domain.Problem, error)
}

type ProblemUseCase struct {
	Repo problemUseCaseRepository
}

func (uc *ProblemUseCase) GetProblemsByContest(ctx context.Context, contestID domain.ResourceID) ([]domain.Problem, error) {
	problems, err := uc.Repo.GetProblemsByContest(ctx, nil, contestID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	return problems, nil
}

func (uc *ProblemUseCase) UpdateProblem(ctx context.Context, problemID domain.ResourceID, problem domain.Problem) (domain.Problem, error) {
	panic("not implemented")
}

func (uc *ProblemUseCase) DeleteProblem(ctx context.Context, problemID domain.ResourceID) error {
	panic("not implemented")
}

func (uc *ProblemUseCase) CreateProblem(ctx context.Context, contestID domain.ResourceID, problem domain.Problem) (domain.Problem, error) {
	panic("not implemented")
}