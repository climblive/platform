package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/climblive/platform/backend/internal/database"
	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

func (d *Database) GetCompetition(ctx context.Context, tx domain.Transaction, competitionID domain.CompetitionID) (domain.Competition, error) {
	record, err := d.WithTx(tx).GetCompetition(ctx, int32(competitionID))
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return domain.Competition{}, errors.Wrap(domain.ErrNotFound, 0)
	case err != nil:
		return domain.Competition{}, errors.Wrap(err, 0)
	}

	competition := competitionToDomain(record.Competition)

	if timeBegin, ok := record.TimeBegin.(time.Time); ok {
		competition.TimeBegin = timeBegin
	}

	if timeEnd, ok := record.TimeEnd.(time.Time); ok {
		competition.TimeEnd = timeEnd
	}

	competition.RegisteredContenders = int(record.RegisteredContenders)

	return competition, nil
}

func (d *Database) GetAllCompetitions(ctx context.Context, tx domain.Transaction) ([]domain.Competition, error) {
	records, err := d.WithTx(tx).GetAllCompetitions(ctx)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	competitions := make([]domain.Competition, 0)

	for _, record := range records {
		competition := competitionToDomain(record.Competition)

		if timeBegin, ok := record.TimeBegin.(time.Time); ok {
			competition.TimeBegin = timeBegin
		}

		if timeEnd, ok := record.TimeEnd.(time.Time); ok {
			competition.TimeEnd = timeEnd
		}

		competition.RegisteredContenders = int(record.RegisteredContenders)

		competitions = append(competitions, competition)
	}

	return competitions, nil
}

func (d *Database) GetCompetitionsByOrganizer(ctx context.Context, tx domain.Transaction, organizerID domain.OrganizerID) ([]domain.Competition, error) {
	records, err := d.WithTx(tx).GetCompetitionsByOrganizer(ctx, int32(organizerID))
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	competitions := make([]domain.Competition, 0)

	for _, record := range records {
		competition := competitionToDomain(record.Competition)

		if timeBegin, ok := record.TimeBegin.(time.Time); ok {
			competition.TimeBegin = timeBegin
		}

		if timeEnd, ok := record.TimeEnd.(time.Time); ok {
			competition.TimeEnd = timeEnd
		}

		competition.RegisteredContenders = int(record.RegisteredContenders)

		competitions = append(competitions, competition)
	}

	return competitions, nil
}

func (d *Database) GetCompetitionsCurrentlyRunningOrByStartTime(ctx context.Context, tx domain.Transaction, earliestStartTime, latestStartTime time.Time) ([]domain.Competition, error) {
	records, err := d.WithTx(tx).GetCompetitionsCurrentlyRunningOrByStartTime(ctx, database.GetCompetitionsCurrentlyRunningOrByStartTimeParams{
		EarliestStartTime: earliestStartTime,
		LatestStartTime:   latestStartTime,
	})
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	competitions := make([]domain.Competition, 0)

	for _, record := range records {
		competition := competitionToDomain(database.Competition{
			ID:                 record.ID,
			OrganizerID:        record.OrganizerID,
			SeriesID:           record.SeriesID,
			Name:               record.Name,
			Description:        record.Description,
			ArchivedAt:         record.ArchivedAt,
			Created:            record.Created,
			Location:           record.Location,
			Country:            record.Country,
			QualifyingProblems: record.QualifyingProblems,
			Finalists:          record.Finalists,
			UsePoints:          record.UsePoints,
			PooledPoints:       record.PooledPoints,
			MaxAttempts:        record.MaxAttempts,
			PointDeduction:     record.PointDeduction,
			Info:               record.Info,
			GracePeriod:        record.GracePeriod,
			NameRetentionTime:  record.NameRetentionTime,
		})

		if timeBegin, ok := record.TimeBegin.(time.Time); ok {
			competition.TimeBegin = timeBegin
		}

		if timeEnd, ok := record.TimeEnd.(time.Time); ok {
			competition.TimeEnd = timeEnd
		}

		competitions = append(competitions, competition)
	}

	return competitions, nil
}

func (d *Database) StoreCompetition(ctx context.Context, tx domain.Transaction, competition domain.Competition) (domain.Competition, error) {
	params := database.UpsertCompetitionParams{
		ID:                 int32(competition.ID),
		OrganizerID:        int32(competition.Ownership.OrganizerID),
		ArchivedAt:         makeNullTime(competition.ArchivedAt),
		SeriesID:           makeNullInt32(int32(competition.SeriesID)),
		Name:               competition.Name,
		Description:        makeNullString(competition.Description),
		Location:           makeNullString(competition.Location),
		Country:            competition.Country,
		QualifyingProblems: int32(competition.QualifyingProblems),
		Finalists:          int32(competition.Finalists),
		UsePoints:          competition.UsePoints,
		PooledPoints:       competition.PooledPoints,
		MaxAttempts:        int32(competition.MaxAttempts),
		PointDeduction:     competition.PointDeduction.Float64(),
		Info:               makeNullString(competition.Info),
		GracePeriod:        int32(competition.GracePeriod / time.Minute),
		NameRetentionTime:  int32(competition.NameRetentionTime / time.Minute),
		Created:            competition.Created,
	}

	insertID, err := d.WithTx(tx).UpsertCompetition(ctx, params)
	if err != nil {
		return domain.Competition{}, errors.Wrap(err, 0)
	}

	if insertID != 0 {
		competition.ID = domain.CompetitionID(insertID)
	}

	return competition, err
}

func (d *Database) DeleteCompetition(ctx context.Context, tx domain.Transaction, competitionID domain.CompetitionID) error {
	err := d.WithTx(tx).DeleteCompetition(ctx, int32(competitionID))
	if err != nil {
		return errors.Wrap(err, 0)
	}

	return nil
}
