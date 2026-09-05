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
	params := database.UpsertTickParams{
		OrganizerID:   int32(tick.Ownership.OrganizerID),
		ContestID:     int32(tick.ContestID),
		ContenderID:   int32(*tick.Ownership.ContenderID),
		ProblemID:     int32(tick.ProblemID),
		Timestamp:     tick.Timestamp,
		Top:           tick.Top,
		AttemptsTop:   int32(tick.AttemptsTop),
		Zone1:         tick.Zone1,
		AttemptsZone1: int32(tick.AttemptsZone1),
		Zone2:         tick.Zone2,
		AttemptsZone2: int32(tick.AttemptsZone2),
	}

	err := d.WithTx(tx).UpsertTick(ctx, params)
	if err != nil {
		return domain.Tick{}, errors.Wrap(err, 0)
	}

	return tick, nil
}

func (d *Database) DeleteTick(ctx context.Context, tx domain.Transaction, contenderID domain.ContenderID, problemID domain.ProblemID) error {
	err := d.WithTx(tx).DeleteTick(ctx, database.DeleteTickParams{
		ContenderID: int32(contenderID),
		ProblemID:   int32(problemID),
	})
	if err != nil {
		return errors.Wrap(err, 0)
	}

	return nil
}

func (d *Database) GetTickByContenderAndProblem(ctx context.Context, tx domain.Transaction, contenderID domain.ContenderID, problemID domain.ProblemID) (domain.Tick, error) {
	record, err := d.WithTx(tx).GetTickByContenderAndProblem(ctx, database.GetTickByContenderAndProblemParams{
		ContenderID: int32(contenderID),
		ProblemID:   int32(problemID),
	})
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return domain.Tick{}, errors.Wrap(domain.ErrNotFound, 0)
	case err != nil:
		return domain.Tick{}, errors.Wrap(err, 0)
	}

	return tickToDomain(record.Tick), nil
}

func (d *Database) GetTicksByProblem(ctx context.Context, tx domain.Transaction, problemID domain.ProblemID) ([]domain.Tick, error) {
	records, err := d.WithTx(tx).GetTicksByProblem(ctx, int32(problemID))
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	ticks := make([]domain.Tick, 0)

	for _, record := range records {
		ticks = append(ticks, tickToDomain(record.Tick))
	}

	return ticks, nil
}
