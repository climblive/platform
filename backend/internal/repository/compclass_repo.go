package repository

import (
	"context"
	"database/sql"

	"github.com/climblive/platform/backend/internal/database"
	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

func (d *Database) GetCompClass(ctx context.Context, tx domain.Transaction, compClassID domain.CompClassID) (domain.CompClass, error) {
	record, err := d.WithTx(tx).GetCompClass(ctx, int32(compClassID))
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return domain.CompClass{}, errors.Wrap(domain.ErrNotFound, 0)
	case err != nil:
		return domain.CompClass{}, errors.Wrap(err, 0)
	}

	return compClassToDomain(record.CompClass), nil
}

func (d *Database) GetCompClassesByContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) ([]domain.CompClass, error) {
	records, err := d.WithTx(tx).GetCompClassesByContest(ctx, int32(contestID))
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	compClasses := make([]domain.CompClass, 0)

	for _, record := range records {
		compClasses = append(compClasses, compClassToDomain(record.CompClass))
	}

	return compClasses, nil
}

func (d *Database) StoreCompClass(ctx context.Context, tx domain.Transaction, compClass domain.CompClass) (domain.CompClass, error) {
	params := database.UpsertCompClassParams{
		ID:          int32(compClass.ID),
		OrganizerID: int32(compClass.Ownership.OrganizerID),
		ContestID:   int32(compClass.ContestID),
		Name:        compClass.Name,
		Description: makeNullString(compClass.Description),
		Color:       sql.NullString{String: "", Valid: false},
		TimeBegin:   compClass.TimeBegin,
		TimeEnd:     compClass.TimeEnd,
	}

	insertID, err := d.WithTx(tx).UpsertCompClass(ctx, params)
	if err != nil {
		return domain.CompClass{}, errors.Wrap(err, 0)
	}

	if insertID != 0 {
		compClass.ID = domain.CompClassID(insertID)
	}

	return compClass, err
}

func (d *Database) DeleteCompClass(ctx context.Context, tx domain.Transaction, compClassID domain.CompClassID) error {
	err := d.WithTx(tx).DeleteCompClass(ctx, int32(compClassID))
	if err != nil {
		return errors.Wrap(err, 0)
	}

	return nil
}
