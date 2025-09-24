package usecases

import (
	"context"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
	"github.com/google/uuid"
)

type organizerUseCaseRepository interface {
	domain.Transactor

	GetOrganizer(ctx context.Context, tx domain.Transaction, organizerID domain.OrganizerID) (domain.Organizer, error)
	GetOrganizerInvitesByOrganizer(ctx context.Context, tx domain.Transaction, organizerID domain.OrganizerID) ([]domain.OrganizerInvite, error)
	GetOrganizerInvite(ctx context.Context, tx domain.Transaction, inviteID domain.OrganizerInviteID) (domain.OrganizerInvite, error)
	StoreOrganizerInvite(ctx context.Context, tx domain.Transaction, invite domain.OrganizerInvite) error
	DeleteOrganizerInvite(ctx context.Context, tx domain.Transaction, inviteID domain.OrganizerInviteID) error
	AddUserToOrganizer(ctx context.Context, tx domain.Transaction, userID domain.UserID, organizerID domain.OrganizerID) error
	GetUserByUsername(ctx context.Context, tx domain.Transaction, username string) (domain.User, error)
}

type OrganizerUseCase struct {
	Authorizer domain.Authorizer
	Repo       organizerUseCaseRepository
}

func (uc *OrganizerUseCase) GetOrganizerInvitesByOrganizer(ctx context.Context, organizerID domain.OrganizerID) ([]domain.OrganizerInvite, error) {
	organizer, err := uc.Repo.GetOrganizer(ctx, nil, organizerID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, organizer.Ownership); err != nil {
		return nil, errors.Wrap(err, 0)
	}

	invites, err := uc.Repo.GetOrganizerInvitesByOrganizer(ctx, nil, organizerID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	return invites, nil
}

func (uc *OrganizerUseCase) GetOrganizerInvite(ctx context.Context, inviteID domain.OrganizerInviteID) (domain.OrganizerInvite, error) {
	invite, err := uc.Repo.GetOrganizerInvite(ctx, nil, inviteID)
	if err != nil {
		return domain.OrganizerInvite{}, errors.Wrap(err, 0)
	}

	return invite, nil
}

func (uc *OrganizerUseCase) CreateOrganizerInvite(ctx context.Context, organizerID domain.OrganizerID) (domain.OrganizerInvite, error) {
	organizer, err := uc.Repo.GetOrganizer(ctx, nil, organizerID)
	if err != nil {
		return domain.OrganizerInvite{}, errors.Wrap(err, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, organizer.Ownership); err != nil {
		return domain.OrganizerInvite{}, errors.Wrap(err, 0)
	}

	invite := domain.OrganizerInvite{
		ID:          domain.OrganizerInviteID(uuid.New()),
		OrganizerID: organizerID,
		ExpiresAt:   time.Now().Add(7 * 24 * time.Hour),
	}

	if err := uc.Repo.StoreOrganizerInvite(ctx, nil, invite); err != nil {
		return domain.OrganizerInvite{}, errors.Wrap(err, 0)
	}

	return invite, nil
}

func (uc *OrganizerUseCase) DeleteOrganizerInvite(ctx context.Context, inviteID domain.OrganizerInviteID) error {
	_, err := uc.Repo.GetOrganizerInvite(ctx, nil, inviteID)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	uc.Repo.DeleteOrganizerInvite(ctx, nil, inviteID)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	return nil
}

func (uc *OrganizerUseCase) AcceptOrganizerInvite(ctx context.Context, inviteID domain.OrganizerInviteID) error {
	invite, err := uc.Repo.GetOrganizerInvite(ctx, nil, inviteID)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	if time.Now().After(invite.ExpiresAt) {
		return errors.Wrap(domain.ErrExpired, 0)
	}

	authentication, err := uc.Authorizer.GetAuthentication(ctx)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	_, _ = uc.Authorizer.HasOwnership(ctx, domain.OwnershipData{})

	user, err := uc.Repo.GetUserByUsername(ctx, nil, authentication.Username)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	tx, err := uc.Repo.Begin()
	if err != nil {
		return errors.Wrap(err, 0)
	}

	err = uc.Repo.AddUserToOrganizer(ctx, nil, user.ID, invite.OrganizerID)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, 0)
	}

	err = uc.Repo.DeleteOrganizerInvite(ctx, nil, inviteID)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, 0)
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, 0)
	}

	return nil
}
