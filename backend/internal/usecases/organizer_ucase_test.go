package usecases_test

import (
	"context"
	"testing"
	"testing/synctest"
	"time"

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

func TestDeleteOrganizerInvite(t *testing.T) {
	fakedInviteID := domain.OrganizerInviteID(uuid.New())
	fakedInvite := domain.OrganizerInvite{ID: fakedInviteID}

	t.Run("HappyPath", func(t *testing.T) {
		mockedRepo := new(repositoryMock)

		mockedRepo.
			On("GetOrganizerInvite", mock.Anything, mock.Anything, fakedInviteID).
			Return(fakedInvite, nil)

		mockedRepo.
			On("DeleteOrganizerInvite", mock.Anything, mock.Anything, fakedInviteID).
			Return(nil)

		ucase := usecases.OrganizerUseCase{
			Repo: mockedRepo,
		}

		err := ucase.DeleteOrganizerInvite(context.Background(), fakedInviteID)

		require.NoError(t, err)

		mockedRepo.AssertExpectations(t)
	})
}

func TestCreateOrganizerInvite(t *testing.T) {
	fakedOrganizerID := randomResourceID[domain.OrganizerID]()
	fakedOwnership := domain.OwnershipData{OrganizerID: fakedOrganizerID}
	fakedOrganizer := domain.Organizer{ID: fakedOrganizerID, Ownership: fakedOwnership}
	fakedInviteID := domain.OrganizerInviteID(uuid.New())

	t.Run("HappyPath", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			mockedRepo := new(repositoryMock)
			mockedAuthorizer := new(authorizerMock)
			mockedUUIDGenerator := new(uuidGeneratorMock)

			mockedRepo.
				On("GetOrganizer", mock.Anything, mock.Anything, fakedOrganizerID).
				Return(fakedOrganizer, nil)

			mockedAuthorizer.
				On("HasOwnership", mock.Anything, fakedOwnership).
				Return(domain.OrganizerRole, nil)

			mockedUUIDGenerator.
				On("Generate").
				Return(uuid.UUID(fakedInviteID))

			mockedRepo.
				On("StoreOrganizerInvite", mock.Anything, mock.Anything, domain.OrganizerInvite{
					ID:          fakedInviteID,
					OrganizerID: fakedOrganizerID,
					ExpiresAt:   time.Now().Add(7 * 24 * time.Hour),
				}).
				Return(nil)

			ucase := usecases.OrganizerUseCase{
				Repo:          mockedRepo,
				Authorizer:    mockedAuthorizer,
				UUIDGenerator: mockedUUIDGenerator,
			}

			invite, err := ucase.CreateOrganizerInvite(context.Background(), fakedOrganizerID)

			require.NoError(t, err)
			assert.Equal(t, fakedOrganizerID, invite.OrganizerID)
			assert.Equal(t, fakedInviteID, invite.ID)
			assert.Equal(t, time.Now().Add(7*24*time.Hour), invite.ExpiresAt)

			mockedUUIDGenerator.AssertExpectations(t)
			mockedAuthorizer.AssertExpectations(t)
			mockedRepo.AssertExpectations(t)
		})
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

		invite, err := ucase.CreateOrganizerInvite(context.Background(), fakedOrganizerID)

		assert.ErrorIs(t, err, domain.ErrNoOwnership)
		assert.Equal(t, domain.OrganizerInvite{}, invite)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})
}

func TestAcceptOrganizerInvite(t *testing.T) {
	fakedInviteID := domain.OrganizerInviteID(uuid.New())
	fakedOrganizerID := randomResourceID[domain.OrganizerID]()
	fakedUserID := randomResourceID[domain.UserID]()
	fakedUsername := "alice"

	t.Run("HappyPath", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			mockedRepo := new(repositoryMock)
			mockedAuthorizer := new(authorizerMock)
			mockedTx := new(transactionMock)

			fakedInvite := domain.OrganizerInvite{
				ID:          fakedInviteID,
				OrganizerID: fakedOrganizerID,
				ExpiresAt:   time.Now().Add(time.Nanosecond),
			}

			fakedUser := domain.User{ID: fakedUserID, Username: fakedUsername}

			mockedRepo.
				On("GetOrganizerInvite", mock.Anything, mock.Anything, fakedInviteID).
				Return(fakedInvite, nil)

			mockedAuthorizer.
				On("GetAuthentication", mock.Anything).
				Return(domain.Authentication{Username: fakedUsername}, nil)

			mockedAuthorizer.
				On("HasOwnership", mock.Anything, domain.OwnershipData{}).
				Return(domain.NilRole, nil)

			mockedRepo.
				On("GetUserByUsername", mock.Anything, mock.Anything, fakedUsername).
				Return(fakedUser, nil)

			mockedRepo.
				On("Begin").
				Return(mockedTx, nil)

			mockedRepo.
				On("AddUserToOrganizer", mock.Anything, mockedTx, fakedUserID, fakedOrganizerID).
				Return(nil)

			mockedRepo.
				On("DeleteOrganizerInvite", mock.Anything, mockedTx, fakedInviteID).
				Return(nil)

			mockedTx.
				On("Commit").
				Return(nil)

			ucase := usecases.OrganizerUseCase{
				Repo:       mockedRepo,
				Authorizer: mockedAuthorizer,
			}

			err := ucase.AcceptOrganizerInvite(context.Background(), fakedInviteID)

			require.NoError(t, err)

			mockedAuthorizer.AssertExpectations(t)
			mockedRepo.AssertExpectations(t)
			mockedTx.AssertExpectations(t)
		})
	})

	t.Run("ExpiredInvite", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			mockedRepo := new(repositoryMock)

			fakedInvite := domain.OrganizerInvite{
				ID:          fakedInviteID,
				OrganizerID: fakedOrganizerID,
				ExpiresAt:   time.Now(),
			}

			mockedRepo.
				On("GetOrganizerInvite", mock.Anything, mock.Anything, fakedInviteID).
				Return(fakedInvite, nil)

			ucase := usecases.OrganizerUseCase{
				Repo: mockedRepo,
			}

			err := ucase.AcceptOrganizerInvite(context.Background(), fakedInviteID)

			assert.ErrorIs(t, err, domain.ErrExpired)

			mockedRepo.AssertExpectations(t)
		})
	})
}
