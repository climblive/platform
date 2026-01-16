package repository

import (
	"context"
	"database/sql"

	"github.com/climblive/platform/backend/internal/database"
	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
	"github.com/google/uuid"
)

func (d *Database) GetTicksByContender(ctx context.Context, tx domain.Transaction, contenderID domain.ContenderID) ([]domain.Tick, error) {
	records, err := d.WithTx(tx).GetTicksByContender(ctx, uuid.UUID(contenderID))
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
	records, err := d.WithTx(tx).GetTicksByContest(ctx, uuid.UUID(contestID))
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
	if uuid.UUID(tick.ID) == uuid.Nil {
		tick.ID = domain.TickID(uuid.New())
	}

	params := database.InsertTickParams{
		ID:            uuid.UUID(tick.ID),
		OrganizerID:   uuid.UUID(tick.Ownership.OrganizerID),
		ContestID:     uuid.UUID(tick.ContestID),
		ContenderID:   uuid.UUID(*tick.Ownership.ContenderID),
		ProblemID:     uuid.UUID(tick.ProblemID),
		Timestamp:     tick.Timestamp,
		Top:           tick.Top,
		AttemptsTop:   int32(tick.AttemptsTop),
		Zone1:         tick.Zone1,
		AttemptsZone1: int32(tick.AttemptsZone1),
		Zone2:         tick.Zone2,
		AttemptsZone2: int32(tick.AttemptsZone2),
	}

	_, err := d.WithTx(tx).InsertTick(ctx, params)
	switch {
	case mysqlDuplicateKeyConstraintViolation.Is(err):
		return domain.Tick{}, errors.New(domain.ErrDuplicate)
	case err != nil:
		return domain.Tick{}, errors.Wrap(err, 0)
	}

	return tick, nil
}

func (d *Database) DeleteTick(ctx context.Context, tx domain.Transaction, tickID domain.TickID) error {
	err := d.WithTx(tx).DeleteTick(ctx, uuid.UUID(tickID))
	if err != nil {
		return errors.Wrap(err, 0)
	}

	return nil
}

func (d *Database) GetTick(ctx context.Context, tx domain.Transaction, tickID domain.TickID) (domain.Tick, error) {
	record, err := d.WithTx(tx).GetTick(ctx, uuid.UUID(tickID))
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return domain.Tick{}, errors.Wrap(domain.ErrNotFound, 0)
	case err != nil:
		return domain.Tick{}, errors.Wrap(err, 0)
	}

	return tickToDomain(record.Tick), nil
}

func (d *Database) GetTicksByProblem(ctx context.Context, tx domain.Transaction, problemID domain.ProblemID) ([]domain.Tick, error) {
	records, err := d.WithTx(tx).GetTicksByProblem(ctx, uuid.UUID(problemID))
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	ticks := make([]domain.Tick, 0)

	for _, record := range records {
		ticks = append(ticks, tickToDomain(record.Tick))
	}

	return ticks, nil
}
