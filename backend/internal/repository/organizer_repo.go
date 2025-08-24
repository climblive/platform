package repository

import (
	"context"
	"database/sql"

	"github.com/climblive/platform/backend/internal/database"
	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

func (d *Database) StoreOrganizer(ctx context.Context, tx domain.Transaction, organizer domain.Organizer) (domain.Organizer, error) {
	params := database.UpsertOrganizerParams{
		ID:   int32(organizer.ID),
		Name: organizer.Name,
	}

	insertID, err := d.WithTx(tx).UpsertOrganizer(ctx, params)
	if err != nil {
		return domain.Organizer{}, errors.Wrap(err, 0)
	}

	organizer.ID = domain.OrganizerID(insertID)

	return organizer, nil
}

func (d *Database) GetOrganizer(ctx context.Context, tx domain.Transaction, organizerID domain.OrganizerID) (domain.Organizer, error) {
	record, err := d.WithTx(tx).GetOrganizer(ctx, int32(organizerID))
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return domain.Organizer{}, errors.Wrap(domain.ErrNotFound, 0)
	case err != nil:
		return domain.Organizer{}, errors.Wrap(err, 0)
	}

	return organizerToDomain(record), nil
}

func (d *Database) GetAllOrganizers(ctx context.Context, tx domain.Transaction) ([]domain.Organizer, error) {
	records, err := d.WithTx(tx).GetAllOrganizers(ctx)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, errors.Wrap(domain.ErrNotFound, 0)
	case err != nil:
		return nil, errors.Wrap(err, 0)
	}

	organizers := make([]domain.Organizer, 0)

	for _, record := range records {
		organizers = append(organizers, organizerToDomain(record))
	}

	return organizers, nil
}
