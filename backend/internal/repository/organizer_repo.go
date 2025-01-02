package repository

import (
	"context"

	"github.com/climblive/platform/backend/internal/database"
	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

func (d *Database) StoreOrganizer(ctx context.Context, tx domain.Transaction, organizer domain.Organizer) (domain.Organizer, error) {
	params := database.UpsertOrganizerParams{
		ID:       int32(organizer.ID),
		Name:     organizer.Name,
		Homepage: makeNullString(organizer.Homepage),
	}

	insertID, err := d.WithTx(tx).UpsertOrganizer(ctx, params)
	if err != nil {
		return domain.Organizer{}, errors.Wrap(err, 0)
	}

	organizer.ID = domain.OrganizerID(insertID)

	return organizer, nil
}
