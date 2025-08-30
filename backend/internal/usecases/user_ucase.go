package usecases

import (
	"context"
	"log/slog"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

type userUseCaseRepository interface {
	domain.Transactor

	GetUserByUsername(ctx context.Context, tx domain.Transaction, username string) (domain.User, error)
	GetAllOrganizers(ctx context.Context, tx domain.Transaction) ([]domain.Organizer, error)
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

	slog.Error("GetSelf", "username", authentication.Username)
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
