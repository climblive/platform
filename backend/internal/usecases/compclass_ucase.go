package usecases

import (
	"context"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

type compClassUseCaseRepository interface {
	domain.Transactor

	GetCompClassesByContest(ctx context.Context, tx domain.Transaction, contestID domain.ResourceID) ([]domain.CompClass, error)
}

type CompClassUseCase struct {
	Repo compClassUseCaseRepository
}

func (uc *CompClassUseCase) GetCompClassesByContest(ctx context.Context, contestID domain.ResourceID) ([]domain.CompClass, error) {
	compClasses, err := uc.Repo.GetCompClassesByContest(ctx, nil, contestID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	return compClasses, nil
}

func (uc *CompClassUseCase) UpdateCompClass(ctx context.Context, compClassID domain.ResourceID, compClass domain.CompClass) (domain.CompClass, error) {
	panic("not implemented")
}

func (uc *CompClassUseCase) DeleteCompClass(ctx context.Context, compClassID domain.ResourceID) error {
	panic("not implemented")
}

func (uc *CompClassUseCase) CreateCompClass(ctx context.Context, contestID domain.ResourceID, compClass domain.CompClass) (domain.CompClass, error) {
	panic("not implemented")
}
