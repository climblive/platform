package repository

import (
	"context"
	"database/sql"

	"github.com/climblive/platform/backend/internal/database"
	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

func (d *Database) GetRafflesByContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) ([]domain.Raffle, error) {
	records, err := d.WithTx(tx).GetRafflesByContest(ctx, int32(contestID))
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	raffles := make([]domain.Raffle, 0)

	for _, record := range records {
		raffles = append(raffles, raffleToDomain(record.Raffle))
	}

	return raffles, nil
}

func (d *Database) GetRaffle(ctx context.Context, tx domain.Transaction, raffleID domain.RaffleID) (domain.Raffle, error) {
	record, err := d.WithTx(tx).GetRaffle(ctx, int32(raffleID))
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return domain.Raffle{}, errors.Wrap(domain.ErrNotFound, 0)
	case err != nil:
		return domain.Raffle{}, errors.Wrap(err, 0)
	}

	return raffleToDomain(record.Raffle), nil
}

func (d *Database) StoreRaffle(ctx context.Context, tx domain.Transaction, raffle domain.Raffle) (domain.Raffle, error) {
	params := database.UpsertRaffleParams{
		ID:          int32(raffle.ID),
		ContestID:   int32(raffle.ContestID),
		OrganizerID: int32(raffle.Ownership.OrganizerID),
	}

	insertID, err := d.WithTx(tx).UpsertRaffle(ctx, params)
	if err != nil {
		return domain.Raffle{}, errors.Wrap(err, 0)
	}

	if insertID != 0 {
		raffle.ID = domain.RaffleID(insertID)
	}

	return raffle, err
}

func (d *Database) GetRaffleWinners(ctx context.Context, tx domain.Transaction, raffleID domain.RaffleID) ([]domain.RaffleWinner, error) {
	records, err := d.WithTx(tx).GetRaffleWinners(ctx, int32(raffleID))
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	winners := make([]domain.RaffleWinner, 0)

	for _, record := range records {
		winners = append(winners, raffleWinnerToDomain(record.RaffleWinner, record.Name.String))
	}

	return winners, nil
}

func (d *Database) StoreRaffleWinner(ctx context.Context, tx domain.Transaction, winner domain.RaffleWinner) (domain.RaffleWinner, error) {
	params := database.UpsertRaffleWinnerParams{
		ID:          int32(winner.ID),
		RaffleID:    int32(winner.RaffleID),
		OrganizerID: int32(winner.Ownership.OrganizerID),
		ContenderID: int32(winner.ContenderID),
		Timestamp:   winner.Timestamp,
	}

	insertID, err := d.WithTx(tx).UpsertRaffleWinner(ctx, params)
	if err != nil {
		return domain.RaffleWinner{}, errors.Wrap(err, 0)
	}

	if insertID != 0 {
		winner.ID = domain.RaffleWinnerID(insertID)
	}

	return winner, nil
}

func (d *Database) DeleteRaffle(ctx context.Context, tx domain.Transaction, raffleID domain.RaffleID) error {
	err := d.WithTx(tx).DeleteRaffle(ctx, int32(raffleID))
	if err != nil {
		return errors.Wrap(err, 0)
	}

	return nil
}

func (d *Database) DeleteRaffleWinner(ctx context.Context, tx domain.Transaction, raffleWinnerID domain.RaffleWinnerID) error {
	err := d.WithTx(tx).DeleteRaffleWinner(ctx, int32(raffleWinnerID))
	if err != nil {
		return errors.Wrap(err, 0)
	}

	return nil
}
