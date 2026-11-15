package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/climblive/platform/backend/internal/database"
	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

func (d *Database) GetCompetitor(ctx context.Context, tx domain.Transaction, competitorID domain.CompetitorID) (domain.Competitor, error) {
	record, err := d.WithTx(tx).GetCompetitor(ctx, int32(competitorID))
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return domain.Competitor{}, errors.Wrap(domain.ErrNotFound, 0)
	case err != nil:
		return domain.Competitor{}, errors.Wrap(err, 0)
	}

	competitor := competitorToDomain(record)

	return competitor, nil
}

func (d *Database) GetCompetitorByCode(ctx context.Context, tx domain.Transaction, registrationCode string) (domain.Competitor, error) {
	record, err := d.WithTx(tx).GetCompetitorByCode(ctx, registrationCode)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return domain.Competitor{}, errors.Wrap(domain.ErrNotFound, 0)
	case err != nil:
		return domain.Competitor{}, errors.Wrap(err, 0)
	}

	competitor := competitorToDomain(database.GetCompetitorRow(record))

	return competitor, nil
}

func (d *Database) GetCompetitorsByCategory(ctx context.Context, tx domain.Transaction, categoryID domain.CategoryID) ([]domain.Competitor, error) {
	records, err := d.WithTx(tx).GetCompetitorsByCategory(ctx, sql.NullInt32{Valid: true, Int32: int32(categoryID)})
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	competitors := make([]domain.Competitor, 0)

	for _, record := range records {
		competitor := competitorToDomain(database.GetCompetitorRow(record))

		competitors = append(competitors, competitor)
	}

	return competitors, nil
}

func (d *Database) GetCompetitorsByCompetition(ctx context.Context, tx domain.Transaction, competitionID domain.CompetitionID) ([]domain.Competitor, error) {
	records, err := d.WithTx(tx).GetCompetitorsByCompetition(ctx, int32(competitionID))
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	competitors := make([]domain.Competitor, 0)

	for _, record := range records {
		competitor := competitorToDomain(database.GetCompetitorRow(record))

		competitors = append(competitors, competitor)
	}

	return competitors, nil
}

func (d *Database) StoreCompetitor(ctx context.Context, tx domain.Transaction, competitor domain.Competitor) (domain.Competitor, error) {
	params := database.UpsertCompetitorParams{
		ID:                  int32(competitor.ID),
		OrganizerID:         int32(competitor.Ownership.OrganizerID),
		CompetitionID:       int32(competitor.CompetitionID),
		RegistrationCode:    competitor.RegistrationCode,
		Name:                makeNullString(competitor.Name),
		CategoryID:          makeNullInt32(int32(competitor.CategoryID)),
		Entered:             makeNullTime(competitor.Entered),
		Disqualified:        competitor.Disqualified,
		WithdrawnFromFinals: competitor.WithdrawnFromFinals,
		ScrubbedAt:          makeNullTime(competitor.ScrubbedAt),
		ScrubBefore:         makeNullTime(competitor.ScrubBefore),
	}

	insertID, err := d.WithTx(tx).UpsertCompetitor(ctx, params)
	if err != nil {
		return domain.Competitor{}, errors.Wrap(err, 0)
	}

	if insertID != 0 {
		competitor.ID = domain.CompetitorID(insertID)
	}

	return competitor, err
}

func (d *Database) DeleteCompetitor(ctx context.Context, tx domain.Transaction, competitorID domain.CompetitorID) error {
	err := d.WithTx(tx).DeleteCompetitor(ctx, int32(competitorID))
	if err != nil {
		return errors.Wrap(err, 0)
	}

	return nil
}

func (d *Database) GetNumberOfCompetitors(ctx context.Context, tx domain.Transaction, competitionID domain.CompetitionID) (int, error) {
	count, err := d.WithTx(tx).CountCompetitors(ctx, int32(competitionID))
	if err != nil {
		return 0, errors.Wrap(err, 0)
	}

	return int(count), nil
}

func (d *Database) GetScrubEligibleCompetitors(ctx context.Context, deadline time.Time) ([]domain.Competitor, error) {
	records, err := d.queries.GetScrubEligibleCompetitors(ctx, makeNullTime(deadline))
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	competitors := make([]domain.Competitor, 0)

	for _, record := range records {
		competitor := competitorToDomain(database.GetCompetitorRow(record))

		competitors = append(competitors, competitor)
	}

	return competitors, nil
}
