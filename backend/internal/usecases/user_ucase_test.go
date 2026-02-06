package usecases_test

import (
	"context"
	"testing"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/utils/testutils"
	"github.com/climblive/platform/backend/internal/usecases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetSelf(t *testing.T) {
	fakedUserID := testutils.RandomResourceID[domain.UserID]()

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

func TestGetUsersByOrganizer(t *testing.T) {
	fakedOrganizerID := testutils.RandomResourceID[domain.OrganizerID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: fakedOrganizerID,
	}
	fakedOrganizer := domain.Organizer{
		ID:        fakedOrganizerID,
		Ownership: fakedOwnership,
	}
	fakedUsers := []domain.User{
		{ID: testutils.RandomResourceID[domain.UserID](), Username: "alice"},
		{ID: testutils.RandomResourceID[domain.UserID](), Username: "bob"},
	}

	t.Run("HappyPath", func(t *testing.T) {
		mockedRepo := new(repositoryMock)
		mockedAuthorizer := new(authorizerMock)

		mockedRepo.
			On("GetOrganizer", mock.Anything, mock.Anything, fakedOrganizerID).
			Return(fakedOrganizer, nil)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("GetUsersByOrganizer", mock.Anything, mock.Anything, fakedOrganizerID).
			Return(fakedUsers, nil)

		ucase := usecases.UserUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		users, err := ucase.GetUsersByOrganizer(context.Background(), fakedOrganizerID)

		require.NoError(t, err)
		assert.Equal(t, fakedUsers, users)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("NoOwnership", func(t *testing.T) {
		mockedRepo := new(repositoryMock)
		mockedAuthorizer := new(authorizerMock)

		mockedRepo.
			On("GetOrganizer", mock.Anything, mock.Anything, fakedOrganizerID).
			Return(fakedOrganizer, nil)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		ucase := usecases.UserUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		_, err := ucase.GetUsersByOrganizer(context.Background(), fakedOrganizerID)

		assert.ErrorIs(t, err, domain.ErrNoOwnership)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})
}
