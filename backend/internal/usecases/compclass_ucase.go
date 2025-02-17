package usecases

import (
	"context"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

type compClassUseCaseRepository interface {
	domain.Transactor

	GetCompClassesByContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) ([]domain.CompClass, error)
}

type CompClassUseCase struct {
	Repo compClassUseCaseRepository
}

func (uc *CompClassUseCase) GetCompClassesByContest(ctx context.Context, contestID domain.ContestID) ([]domain.CompClass, error) {
	compClasses, err := uc.Repo.GetCompClassesByContest(ctx, nil, contestID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	return compClasses, nil
}
