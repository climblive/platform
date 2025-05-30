package usecases

import (
	"context"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

type raffleUseCaseRepository interface {
	domain.Transactor

	StoreRaffle(ctx context.Context, tx domain.Transaction, raffle domain.Raffle) (domain.Raffle, error)
	GetContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) (domain.Contest, error)
	GetRaffle(ctx context.Context, tx domain.Transaction, raffleID domain.RaffleID) (domain.Raffle, error)
	GetRafflesByContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) ([]domain.Raffle, error)
}

type RaffleUseCase struct {
	Authorizer domain.Authorizer
	Repo       raffleUseCaseRepository
}

func (uc *RaffleUseCase) GetRaffle(ctx context.Context, raffleID domain.RaffleID) (domain.Raffle, error) {
	raffle, err := uc.Repo.GetRaffle(ctx, nil, raffleID)
	if err != nil {
		return domain.Raffle{}, errors.Wrap(err, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, raffle.Ownership); err != nil {
		return domain.Raffle{}, errors.Wrap(err, 0)
	}

	return raffle, nil
}

func (uc *RaffleUseCase) GetRafflesByContest(ctx context.Context, contestID domain.ContestID) ([]domain.Raffle, error) {
	contest, err := uc.Repo.GetContest(ctx, nil, contestID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, contest.Ownership); err != nil {
		return nil, errors.Wrap(err, 0)
	}

	raffles, err := uc.Repo.GetRafflesByContest(ctx, nil, contestID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	return raffles, nil
}

func (uc *RaffleUseCase) CreateRaffle(ctx context.Context, contestID domain.ContestID) (domain.Raffle, error) {
	contest, err := uc.Repo.GetContest(ctx, nil, contestID)
	if err != nil {
		return domain.Raffle{}, errors.Wrap(err, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, contest.Ownership); err != nil {
		return domain.Raffle{}, errors.Wrap(err, 0)
	}

	raffle := domain.Raffle{
		Ownership: contest.Ownership,
		ContestID: contestID,
	}

	createdRaffle, err := uc.Repo.StoreRaffle(ctx, nil, raffle)
	if err != nil {
		return domain.Raffle{}, errors.Wrap(err, 0)
	}

	return createdRaffle, nil
}
