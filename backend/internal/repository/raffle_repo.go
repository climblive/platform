package repository

import (
	"context"
	"database/sql"

	"github.com/climblive/platform/backend/internal/database"
	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
	"github.com/google/uuid"
)

func (d *Database) GetRafflesByContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) ([]domain.Raffle, error) {
	records, err := d.WithTx(tx).GetRafflesByContest(ctx, uuid.UUID(contestID))
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
	record, err := d.WithTx(tx).GetRaffle(ctx, uuid.UUID(raffleID))
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return domain.Raffle{}, errors.Wrap(domain.ErrNotFound, 0)
	case err != nil:
		return domain.Raffle{}, errors.Wrap(err, 0)
	}

	return raffleToDomain(record.Raffle), nil
}

func (d *Database) StoreRaffle(ctx context.Context, tx domain.Transaction, raffle domain.Raffle) (domain.Raffle, error) {
	if uuid.UUID(raffle.ID) == uuid.Nil {
		raffle.ID = domain.RaffleID(uuid.New())
	}

	params := database.UpsertRaffleParams{
		ID:          uuid.UUID(raffle.ID),
		ContestID:   uuid.UUID(raffle.ContestID),
		OrganizerID: uuid.UUID(raffle.Ownership.OrganizerID),
	}

	_, err := d.WithTx(tx).UpsertRaffle(ctx, params)
	if err != nil {
		return domain.Raffle{}, errors.Wrap(err, 0)
	}

	return raffle, err
}

func (d *Database) GetRaffleWinners(ctx context.Context, tx domain.Transaction, raffleID domain.RaffleID) ([]domain.RaffleWinner, error) {
	records, err := d.WithTx(tx).GetRaffleWinners(ctx, uuid.UUID(raffleID))
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
	if uuid.UUID(winner.ID) == uuid.Nil {
		winner.ID = domain.RaffleWinnerID(uuid.New())
	}

	params := database.UpsertRaffleWinnerParams{
		ID:          uuid.UUID(winner.ID),
		RaffleID:    uuid.UUID(winner.RaffleID),
		OrganizerID: uuid.UUID(winner.Ownership.OrganizerID),
		ContenderID: uuid.UUID(winner.ContenderID),
		Timestamp:   winner.Timestamp,
	}

	_, err := d.WithTx(tx).UpsertRaffleWinner(ctx, params)
	if err != nil {
		return domain.RaffleWinner{}, errors.Wrap(err, 0)
	}

	return winner, nil
}

func (d *Database) DeleteRaffle(ctx context.Context, tx domain.Transaction, raffleID domain.RaffleID) error {
	err := d.WithTx(tx).DeleteRaffle(ctx, uuid.UUID(raffleID))
	if err != nil {
		return errors.Wrap(err, 0)
	}

	return nil
}

func (d *Database) DeleteRaffleWinner(ctx context.Context, tx domain.Transaction, raffleWinnerID domain.RaffleWinnerID) error {
	err := d.WithTx(tx).DeleteRaffleWinner(ctx, uuid.UUID(raffleWinnerID))
	if err != nil {
		return errors.Wrap(err, 0)
	}

	return nil
}
