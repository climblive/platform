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
