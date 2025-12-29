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

func (d *Database) GetOrganizerInvitesByOrganizer(ctx context.Context, tx domain.Transaction, organizerID domain.OrganizerID) ([]domain.OrganizerInvite, error) {
	records, err := d.WithTx(tx).GetOrganizerInvitesByOrganizer(ctx, int32(organizerID))
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, errors.Wrap(domain.ErrNotFound, 0)
	case err != nil:
		return nil, errors.Wrap(err, 0)
	}

	invites := make([]domain.OrganizerInvite, 0)

	for _, record := range records {
		invites = append(invites, organizerInviteToDomain(record.OrganizerInvite, record.Name))
	}

	return invites, nil
}

func (d *Database) GetOrganizerInvite(ctx context.Context, tx domain.Transaction, inviteID domain.OrganizerInviteID) (domain.OrganizerInvite, error) {
	record, err := d.WithTx(tx).GetOrganizerInvite(ctx, inviteID.String())
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return domain.OrganizerInvite{}, errors.Wrap(domain.ErrNotFound, 0)
	case err != nil:
		return domain.OrganizerInvite{}, errors.Wrap(err, 0)
	}

	return organizerInviteToDomain(record.OrganizerInvite, record.Name), nil
}

func (d *Database) StoreOrganizerInvite(ctx context.Context, tx domain.Transaction, invite domain.OrganizerInvite) error {
	params := database.InsertOrganizerInviteParams{
		ID:          invite.ID.String(),
		OrganizerID: int32(invite.OrganizerID),
		ExpiresAt:   invite.ExpiresAt,
	}

	err := d.WithTx(tx).InsertOrganizerInvite(ctx, params)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	return nil
}

func (d *Database) DeleteOrganizerInvite(ctx context.Context, tx domain.Transaction, inviteID domain.OrganizerInviteID) error {
	err := d.WithTx(tx).DeleteOrganizerInvite(ctx, inviteID.String())
	if err != nil {
		return errors.Wrap(err, 0)
	}

	return nil
}
