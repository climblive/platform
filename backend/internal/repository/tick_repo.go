package repository

import (
	"context"
	"database/sql"

	"github.com/climblive/platform/backend/internal/database"
	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

func (d *Database) GetTicksByCompetitor(ctx context.Context, tx domain.Transaction, competitorID domain.CompetitorID) ([]domain.Tick, error) {
	records, err := d.WithTx(tx).GetTicksByCompetitor(ctx, int32(competitorID))
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	ticks := make([]domain.Tick, 0)

	for _, record := range records {
		ticks = append(ticks, tickToDomain(record.Tick))
	}

	return ticks, nil
}

func (d *Database) GetTicksByCompetition(ctx context.Context, tx domain.Transaction, competitionID domain.CompetitionID) ([]domain.Tick, error) {
	records, err := d.WithTx(tx).GetTicksByCompetition(ctx, int32(competitionID))
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
		ID:            int32(tick.ID),
		OrganizerID:   int32(tick.Ownership.OrganizerID),
		CompetitionID: int32(tick.CompetitionID),
		CompetitorID:  int32(*tick.Ownership.CompetitorID),
		ProblemID:     int32(tick.ProblemID),
		Timestamp:     tick.Timestamp,
		Revision:      uint32(tick.Revision),
		Top:           tick.Top,
		AttemptsTop:   int32(tick.AttemptsTop),
		Zone1:         tick.Zone1,
		AttemptsZone1: int32(tick.AttemptsZone1),
		Zone2:         tick.Zone2,
		AttemptsZone2: int32(tick.AttemptsZone2),
	}

	result, err := d.WithTx(tx).UpsertTick(ctx, params)
	if err != nil {
		return domain.Tick{}, errors.Wrap(err, 0)
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return domain.Tick{}, errors.Wrap(err, 0)
	}

	if affectedRows == 0 {
		return domain.Tick{}, errors.Wrap(domain.ErrSuperseded, 0)
	}

	insertID, err := result.LastInsertId()
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

func (d *Database) GetTickByCompetitorAndProblem(ctx context.Context, tx domain.Transaction, competitorID domain.CompetitorID, problemID domain.ProblemID) (domain.Tick, error) {
	record, err := d.WithTx(tx).GetTickByCompetitorAndProblem(ctx, database.GetTickByCompetitorAndProblemParams{
		CompetitorID: int32(competitorID),
		ProblemID:    int32(problemID),
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
