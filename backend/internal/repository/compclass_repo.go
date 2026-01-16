package repository

import (
	"context"
	"database/sql"

	"github.com/climblive/platform/backend/internal/database"
	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
	"github.com/google/uuid"
)

func (d *Database) GetCompClass(ctx context.Context, tx domain.Transaction, compClassID domain.CompClassID) (domain.CompClass, error) {
	record, err := d.WithTx(tx).GetCompClass(ctx, uuid.UUID(compClassID))
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return domain.CompClass{}, errors.Wrap(domain.ErrNotFound, 0)
	case err != nil:
		return domain.CompClass{}, errors.Wrap(err, 0)
	}

	return compClassToDomain(record.CompClass), nil
}

func (d *Database) GetCompClassesByContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) ([]domain.CompClass, error) {
	records, err := d.WithTx(tx).GetCompClassesByContest(ctx, uuid.UUID(contestID))
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
	if uuid.UUID(compClass.ID) == uuid.Nil {
		compClass.ID = domain.CompClassID(uuid.New())
	}

	params := database.UpsertCompClassParams{
		ID:          uuid.UUID(compClass.ID),
		OrganizerID: uuid.UUID(compClass.Ownership.OrganizerID),
		ContestID:   uuid.UUID(compClass.ContestID),
		Name:        compClass.Name,
		Description: makeNullString(compClass.Description),
		TimeBegin:   compClass.TimeBegin,
		TimeEnd:     compClass.TimeEnd,
	}

	_, err := d.WithTx(tx).UpsertCompClass(ctx, params)
	if err != nil {
		return domain.CompClass{}, errors.Wrap(err, 0)
	}

	return compClass, err
}

func (d *Database) DeleteCompClass(ctx context.Context, tx domain.Transaction, compClassID domain.CompClassID) error {
	err := d.WithTx(tx).DeleteCompClass(ctx, uuid.UUID(compClassID))
	if err != nil {
		return errors.Wrap(err, 0)
	}

	return nil
}
