package usecases

import (
	"context"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

type userUseCaseRepository interface {
	domain.Transactor

	GetUserByUsername(ctx context.Context, tx domain.Transaction, username string) (domain.User, error)
	GetUsersByOrganizer(ctx context.Context, tx domain.Transaction, organizerID domain.OrganizerID) ([]domain.User, error)
	GetAllOrganizers(ctx context.Context, tx domain.Transaction) ([]domain.Organizer, error)
	GetOrganizer(ctx context.Context, tx domain.Transaction, organizerID domain.OrganizerID) (domain.Organizer, error)
}

type UserUseCase struct {
	Authorizer domain.Authorizer
	Repo       userUseCaseRepository
}

func (uc *UserUseCase) GetSelf(ctx context.Context) (domain.User, error) {
	authentication, err := uc.Authorizer.GetAuthentication(ctx)
	if err != nil {
		return domain.User{}, errors.Wrap(err, 0)
	}

	if authentication.Username == "" {
		return domain.User{}, errors.Wrap(domain.ErrNotAuthenticated, 0)
	}

	_, _ = uc.Authorizer.HasOwnership(ctx, domain.OwnershipData{
		OrganizerID: 0,
		ContenderID: nil,
	})

	user, err := uc.Repo.GetUserByUsername(ctx, nil, authentication.Username)
	if err != nil {
		return domain.User{}, errors.Wrap(err, 0)
	}

	if user.Admin {
		organizers, err := uc.Repo.GetAllOrganizers(ctx, nil)
		if err != nil {
			return domain.User{}, errors.Wrap(err, 0)
		}

		user.Organizers = organizers
	}

	return user, nil
}

func (uc *UserUseCase) GetUsersByOrganizer(ctx context.Context, organizerID domain.OrganizerID) ([]domain.User, error) {
	organizer, err := uc.Repo.GetOrganizer(ctx, nil, organizerID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, organizer.Ownership); err != nil {
		return nil, errors.Wrap(err, 0)
	}

	users, err := uc.Repo.GetUsersByOrganizer(ctx, nil, organizerID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	return users, nil
}
