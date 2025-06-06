package usecases

import (
	"context"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

type userUseCaseRepository interface {
	domain.Transactor

	GetUserByUsername(ctx context.Context, tx domain.Transaction, username string) (domain.User, error)
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
		return domain.User{}, errors.Wrap(domain.ErrNotAuthorized, 0)
	}

	user, err := uc.Repo.GetUserByUsername(ctx, nil, authentication.Username)
	if err != nil {
		return domain.User{}, errors.Wrap(err, 0)
	}

	return user, nil
}
