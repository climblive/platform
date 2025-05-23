package repository

import (
	"context"

	"github.com/climblive/platform/backend/internal/database"
	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

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
