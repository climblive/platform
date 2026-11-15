package usecases

import (
	"context"
	"crypto/rand"
	"math/big"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

const maxRafflesPerCompetition = 10

type raffleUseCaseRepository interface {
	domain.Transactor

	StoreRaffle(ctx context.Context, tx domain.Transaction, raffle domain.Raffle) (domain.Raffle, error)
	StoreRaffleWinner(ctx context.Context, tx domain.Transaction, winner domain.RaffleWinner) (domain.RaffleWinner, error)
	GetCompetition(ctx context.Context, tx domain.Transaction, competitionID domain.CompetitionID) (domain.Competition, error)
	GetRaffle(ctx context.Context, tx domain.Transaction, raffleID domain.RaffleID) (domain.Raffle, error)
	GetRafflesByCompetition(ctx context.Context, tx domain.Transaction, competitionID domain.CompetitionID) ([]domain.Raffle, error)
	GetCompetitorsByCompetition(ctx context.Context, tx domain.Transaction, competitionID domain.CompetitionID) ([]domain.Competitor, error)
	GetRaffleWinners(ctx context.Context, tx domain.Transaction, raffleID domain.RaffleID) ([]domain.RaffleWinner, error)
	DeleteRaffleWinner(ctx context.Context, tx domain.Transaction, raffleWinnerID domain.RaffleWinnerID) error
	DeleteRaffle(ctx context.Context, tx domain.Transaction, raffleID domain.RaffleID) error
}

type RaffleUseCase struct {
	Authorizer  domain.Authorizer
	Repo        raffleUseCaseRepository
	EventBroker domain.EventBroker
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

func (uc *RaffleUseCase) GetRafflesByCompetition(ctx context.Context, competitionID domain.CompetitionID) ([]domain.Raffle, error) {
	competition, err := uc.Repo.GetCompetition(ctx, nil, competitionID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, competition.Ownership); err != nil {
		return nil, errors.Wrap(err, 0)
	}

	raffles, err := uc.Repo.GetRafflesByCompetition(ctx, nil, competitionID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	return raffles, nil
}

func (uc *RaffleUseCase) CreateRaffle(ctx context.Context, competitionID domain.CompetitionID) (domain.Raffle, error) {
	competition, err := uc.Repo.GetCompetition(ctx, nil, competitionID)
	if err != nil {
		return domain.Raffle{}, errors.Wrap(err, 0)
	}

	_, err = uc.Authorizer.HasOwnership(ctx, competition.Ownership)
	if err != nil {
		return domain.Raffle{}, errors.Wrap(err, 0)
	}

	raffles, err := uc.Repo.GetRafflesByCompetition(ctx, nil, competitionID)
	if err != nil {
		return domain.Raffle{}, errors.Wrap(err, 0)
	}

	if len(raffles) >= maxRafflesPerCompetition {
		return domain.Raffle{}, errors.New(domain.ErrLimitExceeded)
	}

	raffle := domain.Raffle{
		ID:            0,
		Ownership:     competition.Ownership,
		CompetitionID: competitionID,
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

	competitors, err := uc.Repo.GetCompetitorsByCompetition(ctx, nil, raffle.CompetitionID)
	if err != nil {
		return domain.RaffleWinner{}, errors.Wrap(err, 0)
	}

	if len(competitors) == 0 {
		return domain.RaffleWinner{}, domain.ErrAllWinnersDrawn
	}

	winners, err := uc.Repo.GetRaffleWinners(ctx, nil, raffleID)
	if err != nil {
		return domain.RaffleWinner{}, errors.Wrap(err, 0)
	}

	winnersSet := make(map[domain.CompetitorID]struct{})
	for _, winner := range winners {
		winnersSet[winner.CompetitorID] = struct{}{}
	}

	candidates := make([]domain.Competitor, 0)

	for _, competitor := range competitors {
		if competitor.Entered.IsZero() {
			continue
		}

		if !competitor.ScrubbedAt.IsZero() {
			continue
		}

		if competitor.Disqualified {
			continue
		}

		if _, alreadyDrawn := winnersSet[competitor.ID]; !alreadyDrawn {
			candidates = append(candidates, competitor)
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
		ID:                   0,
		Ownership:            raffle.Ownership,
		RaffleID:             raffle.ID,
		CompetitorID:         candidates[winnerIndex.Int64()].ID,
		CompetitorName:       candidates[winnerIndex.Int64()].Name,
		CompetitorScrubbedAt: candidates[winnerIndex.Int64()].ScrubbedAt,
		Timestamp:            time.Now(),
	}

	createdWinner, err := uc.Repo.StoreRaffleWinner(ctx, nil, winner)
	if err != nil {
		return domain.RaffleWinner{}, errors.Wrap(err, 0)
	}

	uc.EventBroker.Dispatch(raffle.CompetitionID, domain.RaffleWinnerDrawnEvent{
		RaffleID:     createdWinner.RaffleID,
		CompetitorID: createdWinner.CompetitorID,
		Timestamp:    createdWinner.Timestamp,
	})

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

func (uc *RaffleUseCase) DeleteRaffle(ctx context.Context, raffleID domain.RaffleID) error {
	raffle, err := uc.Repo.GetRaffle(ctx, nil, raffleID)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, raffle.Ownership); err != nil {
		return errors.Wrap(err, 0)
	}

	winners, err := uc.Repo.GetRaffleWinners(ctx, nil, raffleID)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	tx, err := uc.Repo.Begin()
	if err != nil {
		return errors.Wrap(err, 0)
	}
	defer tx.Rollback()

	for _, winner := range winners {
		err = uc.Repo.DeleteRaffleWinner(ctx, tx, winner.ID)
		if err != nil {
			return errors.Wrap(err, 0)
		}
	}

	err = uc.Repo.DeleteRaffle(ctx, tx, raffleID)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, 0)
	}

	return nil
}
