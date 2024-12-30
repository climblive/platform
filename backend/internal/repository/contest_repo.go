package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

func (d *Database) GetContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) (domain.Contest, error) {
	record, err := d.WithTx(tx).GetContest(ctx, int32(contestID))
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return domain.Contest{}, errors.Wrap(domain.ErrNotFound, 0)
	case err != nil:
		return domain.Contest{}, errors.Wrap(err, 0)
	}

	contest := contestToDomain(record.Contest)

	if timeBegin, ok := record.TimeBegin.(time.Time); ok {
		contest.TimeBegin = &timeBegin
	}

	if timeEnd, ok := record.TimeEnd.(time.Time); ok {
		contest.TimeEnd = &timeEnd
	}

	return contest, nil
}

func (d *Database) GetContestsCurrentlyRunningOrByStartTime(ctx context.Context, tx domain.Transaction, earliestStartTime, latestStartTime time.Time) ([]domain.Contest, error) {
	records, err := d.WithTx(tx).GetContestsCurrentlyRunningOrByStartTime(ctx)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	contests := make([]domain.Contest, 0)

	for _, record := range records {
		contest := contestToDomain(record.Contest)

		if timeBegin, ok := record.TimeBegin.(time.Time); ok {
			contest.TimeBegin = &timeBegin
		}

		if timeEnd, ok := record.TimeEnd.(time.Time); ok {
			contest.TimeEnd = &timeEnd
		}

		contests = append(contests, contest)
	}

	return contests, nil
}
