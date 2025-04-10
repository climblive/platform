package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/climblive/platform/backend/internal/database"
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
		contest.TimeBegin = timeBegin
	}

	if timeEnd, ok := record.TimeEnd.(time.Time); ok {
		contest.TimeEnd = timeEnd
	}

	return contest, nil
}

func (d *Database) GetContestsByOrganizer(ctx context.Context, tx domain.Transaction, organizerID domain.OrganizerID) ([]domain.Contest, error) {
	records, err := d.WithTx(tx).GetContestsByOrganizer(ctx, int32(organizerID))
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	contests := make([]domain.Contest, 0)

	for _, record := range records {
		contest := contestToDomain(record.Contest)

		if timeBegin, ok := record.TimeBegin.(time.Time); ok {
			contest.TimeBegin = timeBegin
		}

		if timeEnd, ok := record.TimeEnd.(time.Time); ok {
			contest.TimeEnd = timeEnd
		}

		contests = append(contests, contest)
	}

	return contests, nil
}

func (d *Database) GetContestsCurrentlyRunningOrByStartTime(ctx context.Context, tx domain.Transaction, earliestStartTime, latestStartTime time.Time) ([]domain.Contest, error) {
	records, err := d.WithTx(tx).GetContestsCurrentlyRunningOrByStartTime(ctx, database.GetContestsCurrentlyRunningOrByStartTimeParams{
		EarliestStartTime: earliestStartTime,
		LatestStartTime:   latestStartTime,
	})
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	contests := make([]domain.Contest, 0)

	for _, record := range records {
		contest := contestToDomain(database.Contest{
			ID:                 record.ID,
			OrganizerID:        record.OrganizerID,
			Protected:          record.Protected,
			SeriesID:           record.SeriesID,
			Name:               record.Name,
			Description:        record.Description,
			Location:           record.Location,
			FinalEnabled:       record.FinalEnabled,
			QualifyingProblems: record.QualifyingProblems,
			Finalists:          record.Finalists,
			Rules:              record.Rules,
			GracePeriod:        record.GracePeriod,
		})

		if timeBegin, ok := record.TimeBegin.(time.Time); ok {
			contest.TimeBegin = timeBegin
		}

		if timeEnd, ok := record.TimeEnd.(time.Time); ok {
			contest.TimeEnd = timeEnd
		}

		contests = append(contests, contest)
	}

	return contests, nil
}

func (d *Database) StoreContest(ctx context.Context, tx domain.Transaction, contest domain.Contest) (domain.Contest, error) {
	params := database.UpsertContestParams{
		ID:                 int32(contest.ID),
		OrganizerID:        int32(contest.Ownership.OrganizerID),
		SeriesID:           makeNullInt32(int32(contest.SeriesID)),
		Name:               contest.Name,
		Description:        makeNullString(contest.Description),
		Location:           makeNullString(contest.Location),
		FinalEnabled:       contest.Finalists > 0,
		QualifyingProblems: int32(contest.QualifyingProblems),
		Finalists:          int32(contest.Finalists),
		Rules:              makeNullString(contest.Rules),
		GracePeriod:        int32(contest.GracePeriod),
	}

	insertID, err := d.WithTx(tx).UpsertContest(ctx, params)
	if err != nil {
		return domain.Contest{}, errors.Wrap(err, 0)
	}

	if insertID != 0 {
		contest.ID = domain.ContestID(insertID)
	}

	return contest, err
}
