package usecases

import (
	"context"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

type tickUseCaseRepository interface {
	domain.Transactor

	GetContender(ctx context.Context, tx domain.Transaction, contenderID domain.ResourceID) (domain.Contender, error)
	GetTicksByContender(ctx context.Context, tx domain.Transaction, contenderID domain.ResourceID) ([]domain.Tick, error)
}

type TickUseCase struct {
	Repo       tickUseCaseRepository
	Authorizer domain.Authorizer
}

func (uc *TickUseCase) GetTicksByContender(ctx context.Context, contenderID domain.ResourceID) ([]domain.Tick, error) {
	contender, err := uc.Repo.GetContender(ctx, nil, contenderID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, contender.Ownership); err != nil {
		return nil, errors.Wrap(err, 0)
	}

	ticks, err := uc.Repo.GetTicksByContender(ctx, nil, contenderID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	return ticks, nil
}

func (uc *TickUseCase) GetTicksByProblem(ctx context.Context, problemID domain.ResourceID) ([]domain.Tick, error) {
	panic("not implemented")
}

func (uc *TickUseCase) DeleteTick(ctx context.Context, tickID domain.ResourceID) error {
	panic("not implemented")
}

func (uc *TickUseCase) CreateTick(ctx context.Context, contenderID domain.ResourceID, tick domain.Tick) (domain.Tick, error) {
	panic("not implemented")
}
