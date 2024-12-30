package repository

import (
	"context"
	"database/sql"

	"github.com/climblive/platform/backend/internal/database"
	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

func (d *Database) GetTicksByContender(ctx context.Context, tx domain.Transaction, contenderID domain.ContenderID) ([]domain.Tick, error) {
	records, err := d.WithTx(tx).GetTicksByContender(ctx, int32(contenderID))
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	ticks := make([]domain.Tick, 0)

	for _, record := range records {
		ticks = append(ticks, tickToDomain(record.Tick))
	}

	return ticks, nil
}

func (d *Database) GetTicksByContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) ([]domain.Tick, error) {
	records, err := d.WithTx(tx).GetTicksByContest(ctx, int32(contestID))
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	ticks := make([]domain.Tick, 0)

	for _, record := range records {
		ticks = append(ticks, tickToDomain(record.Tick))
	}

	return ticks, nil
}

func (d *Database) StoreTick(ctx context.Context, tx domain.Transaction, tick domain.Tick) (domain.Tick, error) {
	params := database.InsertTickParams{
		OrganizerID: int32(tick.Ownership.OrganizerID),
		ContestID:   int32(tick.ContestID),
		ContenderID: int32(*tick.Ownership.ContenderID),
		ProblemID:   int32(tick.ProblemID),
		Flash:       tick.AttemptsTop == 1,
		Timestamp:   tick.Timestamp,
	}

	insertID, err := d.WithTx(tx).InsertTick(ctx, params)
	if err != nil {
		return domain.Tick{}, errors.Wrap(err, 0)
	}

	tick.ID = domain.TickID(insertID)

	return tick, nil
}

func (d *Database) DeleteTick(ctx context.Context, tx domain.Transaction, tickID domain.TickID) error {
	err := d.WithTx(tx).DeleteTick(ctx, int32(tickID))
	if err != nil {
		return errors.Wrap(err, 0)
	}

	return nil
}

func (d *Database) GetTick(ctx context.Context, tx domain.Transaction, tickID domain.TickID) (domain.Tick, error) {
	record, err := d.WithTx(tx).GetTick(ctx, int32(tickID))
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return domain.Tick{}, errors.Wrap(domain.ErrNotFound, 0)
	case err != nil:
		return domain.Tick{}, errors.Wrap(err, 0)
	}

	return tickToDomain(record.Tick), nil
}
