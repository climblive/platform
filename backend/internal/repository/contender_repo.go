package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/climblive/platform/backend/internal/database"
	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

func (d *Database) GetContender(ctx context.Context, tx domain.Transaction, contenderID domain.ContenderID) (domain.Contender, error) {
	record, err := d.WithTx(tx).GetContender(ctx, int32(contenderID))
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return domain.Contender{}, errors.Wrap(domain.ErrNotFound, 0)
	case err != nil:
		return domain.Contender{}, errors.Wrap(err, 0)
	}

	contender := contenderToDomain(record)

	return contender, nil
}

func (d *Database) GetContenderByCode(ctx context.Context, tx domain.Transaction, registrationCode string) (domain.Contender, error) {
	record, err := d.WithTx(tx).GetContenderByCode(ctx, registrationCode)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return domain.Contender{}, errors.Wrap(domain.ErrNotFound, 0)
	case err != nil:
		return domain.Contender{}, errors.Wrap(err, 0)
	}

	contender := contenderToDomain(database.GetContenderRow(record))

	return contender, nil
}

func (d *Database) GetContendersByCategory(ctx context.Context, tx domain.Transaction, categoryID domain.CategoryID) ([]domain.Contender, error) {
	records, err := d.WithTx(tx).GetContendersByCategory(ctx, sql.NullInt32{Valid: true, Int32: int32(categoryID)})
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	contenders := make([]domain.Contender, 0)

	for _, record := range records {
		contender := contenderToDomain(database.GetContenderRow(record))

		contenders = append(contenders, contender)
	}

	return contenders, nil
}

func (d *Database) GetContendersByCompetition(ctx context.Context, tx domain.Transaction, competitionID domain.CompetitionID) ([]domain.Contender, error) {
	records, err := d.WithTx(tx).GetContendersByCompetition(ctx, int32(competitionID))
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	contenders := make([]domain.Contender, 0)

	for _, record := range records {
		contender := contenderToDomain(database.GetContenderRow(record))

		contenders = append(contenders, contender)
	}

	return contenders, nil
}

func (d *Database) StoreContender(ctx context.Context, tx domain.Transaction, contender domain.Contender) (domain.Contender, error) {
	params := database.UpsertContenderParams{
		ID:                  int32(contender.ID),
		OrganizerID:         int32(contender.Ownership.OrganizerID),
		CompetitionID:       int32(contender.CompetitionID),
		RegistrationCode:    contender.RegistrationCode,
		Name:                makeNullString(contender.Name),
		CategoryID:          makeNullInt32(int32(contender.CategoryID)),
		Entered:             makeNullTime(contender.Entered),
		Disqualified:        contender.Disqualified,
		WithdrawnFromFinals: contender.WithdrawnFromFinals,
		ScrubbedAt:          makeNullTime(contender.ScrubbedAt),
		ScrubBefore:         makeNullTime(contender.ScrubBefore),
	}

	insertID, err := d.WithTx(tx).UpsertContender(ctx, params)
	if err != nil {
		return domain.Contender{}, errors.Wrap(err, 0)
	}

	if insertID != 0 {
		contender.ID = domain.ContenderID(insertID)
	}

	return contender, err
}

func (d *Database) DeleteContender(ctx context.Context, tx domain.Transaction, contenderID domain.ContenderID) error {
	err := d.WithTx(tx).DeleteContender(ctx, int32(contenderID))
	if err != nil {
		return errors.Wrap(err, 0)
	}

	return nil
}

func (d *Database) GetNumberOfContenders(ctx context.Context, tx domain.Transaction, competitionID domain.CompetitionID) (int, error) {
	count, err := d.WithTx(tx).CountContenders(ctx, int32(competitionID))
	if err != nil {
		return 0, errors.Wrap(err, 0)
	}

	return int(count), nil
}

func (d *Database) GetScrubEligibleContenders(ctx context.Context, deadline time.Time) ([]domain.Contender, error) {
	records, err := d.queries.GetScrubEligibleContenders(ctx, makeNullTime(deadline))
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	contenders := make([]domain.Contender, 0)

	for _, record := range records {
		contender := contenderToDomain(database.GetContenderRow(record))

		contenders = append(contenders, contender)
	}

	return contenders, nil
}
