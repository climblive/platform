package usecases_test

import (
	"context"
	"testing"
	"testing/synctest"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/utils/testutils"
	"github.com/climblive/platform/backend/internal/usecases"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const whitespaceCharacters = "\u0009\u000A\u000B\u000C\u000D\u0020\u0085\u00A0\u1680\u2000\u2001\u2002\u2003\u2004\u2005\u2006\u2007\u2008\u2009\u200A\u2028\u2029\u202F\u205F\u3000"

func TestCreateOrganizer(t *testing.T) {
	fakedUsername := "alice"
	fakedUserID := testutils.RandomResourceID[domain.UserID]()
	fakedUser := domain.User{ID: fakedUserID, Username: fakedUsername}
	fakedAuthentication := domain.Authentication{Username: fakedUsername}
	fakedOrganizerID := testutils.RandomResourceID[domain.OrganizerID]()

	t.Run("HappyPath", func(t *testing.T) {
		mockedRepo := new(repositoryMock)
		mockedAuthorizer := new(authorizerMock)
		mockedTx := new(transactionMock)

		mockedAuthorizer.
			On("GetAuthentication", mock.Anything).
			Return(fakedAuthentication, nil)

		mockedRepo.
			On("GetUserByUsername", mock.Anything, mock.Anything, fakedUsername).
			Return(fakedUser, nil)

		mockedRepo.
			On("Begin").
			Return(mockedTx, nil)

		mockedRepo.
			On("StoreOrganizer", mock.Anything, mockedTx, domain.Organizer{
				Name: "Test Organizer",
			}).
			Return(domain.Organizer{
				ID: fakedOrganizerID,
				Ownership: domain.OwnershipData{
					OrganizerID: fakedOrganizerID,
				},
				Name: "Test Organizer",
			}, nil)

		mockedRepo.
			On("AddUserToOrganizer", mock.Anything, mockedTx, fakedUserID, fakedOrganizerID).
			Return(nil)

		mockedTx.
			On("Commit").
			Return(nil)

		ucase := usecases.OrganizerUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		organizer, err := ucase.CreateOrganizer(context.Background(), domain.OrganizerTemplate{
			Name: "Test Organizer",
		})

		require.NoError(t, err)
		assert.Equal(t, fakedOrganizerID, organizer.ID)
		assert.Equal(t, "Test Organizer", organizer.Name)
		assert.Equal(t, fakedOrganizerID, organizer.Ownership.OrganizerID)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
		mockedTx.AssertExpectations(t)
	})

	t.Run("EmptyName", func(t *testing.T) {
		mockedRepo := new(repositoryMock)

		ucase := usecases.OrganizerUseCase{
			Repo: mockedRepo,
		}

		organizer, err := ucase.CreateOrganizer(context.Background(), domain.OrganizerTemplate{
			Name: whitespaceCharacters,
		})

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.Equal(t, domain.Organizer{}, organizer)

		mockedRepo.AssertExpectations(t)
	})

	t.Run("ContenderCannotCreateOrganizer", func(t *testing.T) {
		mockedRepo := new(repositoryMock)
		mockedAuthorizer := new(authorizerMock)

		mockedAuthorizer.
			On("GetAuthentication", mock.Anything).
			Return(domain.Authentication{
				Regcode: "ABCD0001",
			}, nil)

		ucase := usecases.OrganizerUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		organizer, err := ucase.CreateOrganizer(context.Background(), domain.OrganizerTemplate{
			Name: "Test Organizer",
		})

		assert.ErrorIs(t, err, domain.ErrNotAuthenticated)
		assert.Equal(t, domain.Organizer{}, organizer)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})
}

func TestGetOrganizer(t *testing.T) {
	fakedOrganizerID := testutils.RandomResourceID[domain.OrganizerID]()
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
	fakedOrganizerID := testutils.RandomResourceID[domain.OrganizerID]()
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
	fakedOrganizerID := testutils.RandomResourceID[domain.OrganizerID]()
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
	fakedOrganizerID := testutils.RandomResourceID[domain.OrganizerID]()
	fakedUserID := testutils.RandomResourceID[domain.UserID]()
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

func TestPatchOrganizer(t *testing.T) {
	fakedOrganizerID := testutils.RandomResourceID[domain.OrganizerID]()
	fakedOwnership := domain.OwnershipData{OrganizerID: fakedOrganizerID}
	fakedOrganizer := domain.Organizer{
		ID:        fakedOrganizerID,
		Ownership: fakedOwnership,
		Name:      "Old Name",
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
			On("StoreOrganizer", mock.Anything, mock.Anything, domain.Organizer{
				ID:        fakedOrganizerID,
				Ownership: fakedOwnership,
				Name:      "New Name",
			}).
			Return(domain.Organizer{
				ID:        fakedOrganizerID,
				Ownership: fakedOwnership,
				Name:      "New Name",
			}, nil)

		ucase := usecases.OrganizerUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		patch := domain.OrganizerPatch{
			Name: domain.NewPatch("New Name"),
		}

		organizer, err := ucase.PatchOrganizer(context.Background(), fakedOrganizerID, patch)

		require.NoError(t, err)
		assert.Equal(t, "New Name", organizer.Name)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("EmptyName", func(t *testing.T) {
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

		patch := domain.OrganizerPatch{
			Name: domain.NewPatch(whitespaceCharacters),
		}

		organizer, err := ucase.PatchOrganizer(context.Background(), fakedOrganizerID, patch)

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.Equal(t, domain.Organizer{}, organizer)

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

		patch := domain.OrganizerPatch{
			Name: domain.NewPatch("New Name"),
		}

		_, err := ucase.PatchOrganizer(context.Background(), fakedOrganizerID, patch)

		assert.ErrorIs(t, err, domain.ErrNoOwnership)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})
}
