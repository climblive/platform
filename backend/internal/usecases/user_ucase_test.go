package usecases_test

import (
	"context"
	"testing"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/usecases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetSelf(t *testing.T) {
	fakedUserID := randomResourceID[domain.UserID]()

	mockedRepo := new(repositoryMock)
	mockedAuthorizer := new(authorizerMock)

	fakedUser := domain.User{
		ID: fakedUserID,
	}

	mockedRepo.
		On("GetUserByUsername", mock.Anything, mock.Anything, "john").
		Return(fakedUser, nil)

	mockedAuthorizer.
		On("GetAuthentication", mock.Anything).
		Return(domain.Authentication{
			Username: "john",
		}, nil)

	mockedAuthorizer.
		On("HasOwnership", mock.Anything, mock.AnythingOfType("domain.OwnershipData")).
		Return(domain.NilRole, domain.ErrNoOwnership)

	ucase := usecases.UserUseCase{
		Repo:       mockedRepo,
		Authorizer: mockedAuthorizer,
	}

	user, err := ucase.GetSelf(context.Background())

	require.NoError(t, err)
	assert.Equal(t, fakedUser, user)

	mockedRepo.AssertExpectations(t)
	mockedAuthorizer.AssertExpectations(t)
}
