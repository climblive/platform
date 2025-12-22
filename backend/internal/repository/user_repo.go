package repository

import (
	"context"

	"github.com/climblive/platform/backend/internal/database"
	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

func (d *Database) StoreUser(ctx context.Context, tx domain.Transaction, user domain.User) (domain.User, error) {
	params := database.UpsertUserParams{
		ID:       int32(user.ID),
		Username: user.Username,
		Admin:    user.Admin,
	}

	insertID, err := d.WithTx(tx).UpsertUser(ctx, params)
	if err != nil {
		return domain.User{}, errors.Wrap(err, 0)
	}

	user.ID = domain.UserID(insertID)

	return user, nil
}

func (d *Database) AddUserToOrganizer(ctx context.Context, tx domain.Transaction, userID domain.UserID, organizerID domain.OrganizerID) error {
	params := database.AddUserToOrganizerParams{
		UserID:      int32(userID),
		OrganizerID: int32(organizerID),
	}

	err := d.WithTx(tx).AddUserToOrganizer(ctx, params)

	if err != nil {
		return errors.Wrap(err, 0)
	}

	return nil
}

func (d *Database) GetUserByUsername(ctx context.Context, tx domain.Transaction, username string) (domain.User, error) {
	records, err := d.WithTx(tx).GetUserByUsername(ctx, username)
	if err != nil {
		return domain.User{}, errors.Wrap(err, 0)
	}

	if len(records) == 0 {
		return domain.User{}, errors.Wrap(domain.ErrNotFound, 0)
	}

	var user domain.User

	for index, record := range records {
		if index == 0 {
			user = userToDomain(record.User)
		}

		user.Organizers = append(user.Organizers, organizerToDomain(record.Organizer))
	}

	return user, nil
}

func (d *Database) GetUsersByOrganizer(ctx context.Context, tx domain.Transaction, organizerID domain.OrganizerID) ([]domain.User, error) {
	records, err := d.WithTx(tx).GetUsersByOrganizer(ctx, int32(organizerID))
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	users := make([]domain.User, 0)

	for _, record := range records {
		users = append(users, userToDomain(record.User))
	}

	return users, nil
}
