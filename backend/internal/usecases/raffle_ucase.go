package usecases

import (
	"context"
	"crypto/rand"
	"math/big"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

type raffleUseCaseRepository interface {
	domain.Transactor

	StoreRaffle(ctx context.Context, tx domain.Transaction, raffle domain.Raffle) (domain.Raffle, error)
	StoreRaffleWinner(ctx context.Context, tx domain.Transaction, winner domain.RaffleWinner) (domain.RaffleWinner, error)
	GetContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) (domain.Contest, error)
	GetRaffle(ctx context.Context, tx domain.Transaction, raffleID domain.RaffleID) (domain.Raffle, error)
	GetRafflesByContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) ([]domain.Raffle, error)
	GetContendersByContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) ([]domain.Contender, error)
	GetRaffleWinners(ctx context.Context, tx domain.Transaction, raffleID domain.RaffleID) ([]domain.RaffleWinner, error)
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

func (uc *RaffleUseCase) DrawRaffleWinner(ctx context.Context, raffleID domain.RaffleID) (domain.RaffleWinner, error) {
	raffle, err := uc.Repo.GetRaffle(ctx, nil, raffleID)
	if err != nil {
		return domain.RaffleWinner{}, errors.Wrap(err, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, raffle.Ownership); err != nil {
		return domain.RaffleWinner{}, errors.Wrap(err, 0)
	}

	contenders, err := uc.Repo.GetContendersByContest(ctx, nil, raffle.ContestID)
	if err != nil {
		return domain.RaffleWinner{}, errors.Wrap(err, 0)
	}

	if len(contenders) == 0 {
		return domain.RaffleWinner{}, domain.ErrAllWinnersDrawn
	}

	winners, err := uc.Repo.GetRaffleWinners(ctx, nil, raffleID)
	if err != nil {
		return domain.RaffleWinner{}, errors.Wrap(err, 0)
	}

	winnersSet := make(map[domain.ContenderID]struct{})
	for _, winner := range winners {
		winnersSet[winner.ContenderID] = struct{}{}
	}

	candidates := make([]domain.Contender, 0)

	for _, contender := range contenders {
		if contender.Entered.IsZero() {
			continue
		}

		if _, alreadyDrawn := winnersSet[contender.ID]; !alreadyDrawn {
			candidates = append(candidates, contender)
		}
	}

	if len(candidates) == 0 {
		return domain.RaffleWinner{}, domain.ErrAllWinnersDrawn
	}

	winnerIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(candidates))))
	if err != nil {
		return domain.RaffleWinner{}, errors.Wrap(err, 0)
	}

	winner := domain.RaffleWinner{
		Ownership:     raffle.Ownership,
		RaffleID:      raffle.ID,
		ContenderID:   candidates[winnerIndex.Int64()].ID,
		ContenderName: candidates[winnerIndex.Int64()].Name,
		Timestamp:     time.Now(),
	}

	createdWinner, err := uc.Repo.StoreRaffleWinner(ctx, nil, winner)
	if err != nil {
		return domain.RaffleWinner{}, errors.Wrap(err, 0)
	}

	return createdWinner, nil
}

func (uc *RaffleUseCase) GetRaffleWinners(ctx context.Context, raffleID domain.RaffleID) ([]domain.RaffleWinner, error) {
	raffle, err := uc.Repo.GetRaffle(ctx, nil, raffleID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, raffle.Ownership); err != nil {
		return nil, errors.Wrap(err, 0)
	}

	winners, err := uc.Repo.GetRaffleWinners(ctx, nil, raffleID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	return winners, nil
}

func (uc *RaffleUseCase) GetRaffleWinnersByContest(ctx context.Context, contestID domain.ContestID) ([]domain.RaffleWinner, error) {
	raffles, err := uc.Repo.GetRafflesByContest(ctx, nil, contestID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	allWinners := make([]domain.RaffleWinner, 0)
	for _, raffle := range raffles {
		winners, err := uc.Repo.GetRaffleWinners(ctx, nil, raffle.ID)
		if err != nil {
			return nil, errors.Wrap(err, 0)
		}
		allWinners = append(allWinners, winners...)
	}

	return allWinners, nil
}
