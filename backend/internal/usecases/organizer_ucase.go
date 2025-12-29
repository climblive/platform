package usecases

import (
	"context"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

type organizerUseCaseRepository interface {
	domain.Transactor

	GetOrganizer(ctx context.Context, tx domain.Transaction, organizerID domain.OrganizerID) (domain.Organizer, error)
	StoreOrganizer(ctx context.Context, tx domain.Transaction, organizer domain.Organizer) (domain.Organizer, error)
	GetOrganizerInvitesByOrganizer(ctx context.Context, tx domain.Transaction, organizerID domain.OrganizerID) ([]domain.OrganizerInvite, error)
	GetOrganizerInvite(ctx context.Context, tx domain.Transaction, inviteID domain.OrganizerInviteID) (domain.OrganizerInvite, error)
	StoreOrganizerInvite(ctx context.Context, tx domain.Transaction, invite domain.OrganizerInvite) error
	DeleteOrganizerInvite(ctx context.Context, tx domain.Transaction, inviteID domain.OrganizerInviteID) error
	AddUserToOrganizer(ctx context.Context, tx domain.Transaction, userID domain.UserID, organizerID domain.OrganizerID) error
	GetUserByUsername(ctx context.Context, tx domain.Transaction, username string) (domain.User, error)
}

type OrganizerUseCase struct {
	Authorizer    domain.Authorizer
	Repo          organizerUseCaseRepository
	UUIDGenerator domain.UUIDGenerator
}

func (uc *OrganizerUseCase) CreateOrganizer(ctx context.Context, template domain.OrganizerTemplate) (domain.Organizer, error) {
	if template.Name == "" {
		return domain.Organizer{}, errors.Wrap(domain.ErrInvalidData, 0)
	}

	authentication, err := uc.Authorizer.GetAuthentication(ctx)
	if err != nil {
		return domain.Organizer{}, errors.Wrap(err, 0)
	}

	user, err := uc.Repo.GetUserByUsername(ctx, nil, authentication.Username)
	if err != nil {
		return domain.Organizer{}, errors.Wrap(err, 0)
	}

	organizer := domain.Organizer{
		Name: template.Name,
	}

	tx, err := uc.Repo.Begin()
	if err != nil {
		return domain.Organizer{}, errors.Wrap(err, 0)
	}

	organizer, err = uc.Repo.StoreOrganizer(ctx, tx, organizer)
	if err != nil {
		tx.Rollback()
		return domain.Organizer{}, errors.Wrap(err, 0)
	}

	err = uc.Repo.AddUserToOrganizer(ctx, tx, user.ID, organizer.ID)
	if err != nil {
		tx.Rollback()
		return domain.Organizer{}, errors.Wrap(err, 0)
	}

	err = tx.Commit()
	if err != nil {
		return domain.Organizer{}, errors.Wrap(err, 0)
	}

	organizer.Ownership = domain.OwnershipData{OrganizerID: organizer.ID}

	return organizer, nil
}

func (uc *OrganizerUseCase) GetOrganizer(ctx context.Context, organizerID domain.OrganizerID) (domain.Organizer, error) {
	organizer, err := uc.Repo.GetOrganizer(ctx, nil, organizerID)
	if err != nil {
		return domain.Organizer{}, errors.Wrap(err, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, organizer.Ownership); err != nil {
		return domain.Organizer{}, errors.Wrap(err, 0)
	}

	return organizer, nil
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
		ID:          domain.OrganizerInviteID(uc.UUIDGenerator.Generate()),
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

	err = uc.Repo.DeleteOrganizerInvite(ctx, nil, inviteID)
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

	now := time.Now()

	if now.Equal(invite.ExpiresAt) || now.After(invite.ExpiresAt) {
		return errors.Wrap(domain.ErrExpired, 0)
	}

	authentication, err := uc.Authorizer.GetAuthentication(ctx)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	user, err := uc.Repo.GetUserByUsername(ctx, nil, authentication.Username)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	tx, err := uc.Repo.Begin()
	if err != nil {
		return errors.Wrap(err, 0)
	}

	err = uc.Repo.AddUserToOrganizer(ctx, tx, user.ID, invite.OrganizerID)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, 0)
	}

	err = uc.Repo.DeleteOrganizerInvite(ctx, tx, inviteID)
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
