package usecases_test

import (
	"context"
	"testing"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/usecases"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetOrganizer(t *testing.T) {
	fakedOrganizerID := randomResourceID[domain.OrganizerID]()
	fakedOwnership := domain.OwnershipData{OrganizerID: fakedOrganizerID}
	fakedOrganizer := domain.Organizer{ID: fakedOrganizerID, Ownership: fakedOwnership}

	t.Run("HappyPath", func(t *testing.T) {
		mockedRepo := new(repositoryMock)
		mockedAuthorizer := new(authorizerMock)

		mockedRepo.
			On("GetOrganizer", mock.Anything, mock.Anything, fakedOrganizerID).
			Return(fakedOrganizer, nil)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		ucase := usecases.OrganizerUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		organizer, err := ucase.GetOrganizer(context.Background(), fakedOrganizerID)

		require.NoError(t, err)
		assert.Equal(t, fakedOrganizer, organizer)

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

		ucase := usecases.OrganizerUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		organizer, err := ucase.GetOrganizer(context.Background(), fakedOrganizerID)

		assert.ErrorIs(t, err, domain.ErrNoOwnership)
		assert.Equal(t, domain.Organizer{}, organizer)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})
}

func TestGetOrganizerInvitesByOrganizer(t *testing.T) {
	fakedOrganizerID := randomResourceID[domain.OrganizerID]()
	fakedOwnership := domain.OwnershipData{OrganizerID: fakedOrganizerID}
	fakedOrganizer := domain.Organizer{ID: fakedOrganizerID, Ownership: fakedOwnership}
	fakedInvites := []domain.OrganizerInvite{
		{ID: domain.OrganizerInviteID(uuid.New())},
		{ID: domain.OrganizerInviteID(uuid.New())},
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
			On("GetOrganizerInvitesByOrganizer", mock.Anything, mock.Anything, fakedOrganizerID).
			Return(fakedInvites, nil)

		ucase := usecases.OrganizerUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		invites, err := ucase.GetOrganizerInvitesByOrganizer(context.Background(), fakedOrganizerID)

		require.NoError(t, err)
		assert.Equal(t, fakedInvites, invites)

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

		ucase := usecases.OrganizerUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		invites, err := ucase.GetOrganizerInvitesByOrganizer(context.Background(), fakedOrganizerID)

		assert.ErrorIs(t, err, domain.ErrNoOwnership)
		assert.Nil(t, invites)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})
}

func TestGetOrganizerInvite(t *testing.T) {
	fakedInviteID := domain.OrganizerInviteID(uuid.New())
	fakedInvite := domain.OrganizerInvite{ID: fakedInviteID}

	t.Run("HappyPath", func(t *testing.T) {
		mockedRepo := new(repositoryMock)

		mockedRepo.
			On("GetOrganizerInvite", mock.Anything, mock.Anything, fakedInviteID).
			Return(fakedInvite, nil)

		ucase := usecases.OrganizerUseCase{
			Repo: mockedRepo,
		}

		invite, err := ucase.GetOrganizerInvite(context.Background(), fakedInviteID)

		require.NoError(t, err)
		assert.Equal(t, fakedInvite, invite)

		mockedRepo.AssertExpectations(t)
	})
}
