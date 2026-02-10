package usecases_test

import (
	"context"
	"testing"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/usecases"
	"github.com/climblive/platform/backend/internal/usecases/validators"
	"github.com/climblive/platform/backend/internal/utils/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetCompClass(t *testing.T) {
	t.Parallel()

	fakedCompClassID := testutils.RandomResourceID[domain.CompClassID]()
	fakedOrganizerID := testutils.RandomResourceID[domain.OrganizerID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: fakedOrganizerID,
	}

	fakedCompClass := domain.CompClass{
		ID:        fakedCompClassID,
		Ownership: fakedOwnership,
	}

	mockedRepo := new(repositoryMock)

	mockedRepo.
		On("GetCompClass", mock.Anything, nil, fakedCompClassID).
		Return(fakedCompClass, nil)

	ucase := usecases.CompClassUseCase{
		Repo: mockedRepo,
	}

	compClass, err := ucase.GetCompClass(context.Background(), fakedCompClassID)

	require.NoError(t, err)
	assert.Equal(t, fakedCompClass, compClass)

	mockedRepo.AssertExpectations(t)
}

func TestGetCompClassesByContest(t *testing.T) {
	t.Parallel()

	fakedContestID := testutils.RandomResourceID[domain.ContestID]()

	fakedCompClasses := []domain.CompClass{{
		ID:        testutils.RandomResourceID[domain.CompClassID](),
		ContestID: fakedContestID,
	},
	}

	mockedRepo := new(repositoryMock)

	mockedRepo.
		On("GetCompClassesByContest", mock.Anything, mock.Anything, fakedContestID).
		Return(fakedCompClasses, nil)

	ucase := usecases.CompClassUseCase{
		Repo: mockedRepo,
	}

	compClasses, err := ucase.GetCompClassesByContest(context.Background(), fakedContestID)

	require.NoError(t, err)
	assert.Equal(t, fakedCompClasses, compClasses)

	mockedRepo.AssertExpectations(t)
}

func TestCreateCompClass(t *testing.T) {
	fakedOrganizerID := testutils.RandomResourceID[domain.OrganizerID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: fakedOrganizerID,
	}
	fakedContestID := testutils.RandomResourceID[domain.ContestID]()
	fakedCompClassID := testutils.RandomResourceID[domain.CompClassID]()

	now := time.Now()

	makeMocks := func() (*repositoryMock, *authorizerMock) {
		mockedRepo := new(repositoryMock)

		mockedRepo.
			On("GetContest", mock.Anything, nil, fakedContestID).
			Return(domain.Contest{
				ID:        fakedContestID,
				Ownership: fakedOwnership,
			}, nil)

		mockedAuthorizer := new(authorizerMock)

		return mockedRepo, mockedAuthorizer
	}

	t.Run("HappyCase", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("StoreCompClass", mock.Anything, nil,
				domain.CompClass{
					Ownership:   fakedOwnership,
					ContestID:   fakedContestID,
					Name:        "Females",
					Description: "Female climbers",
					TimeBegin:   now,
					TimeEnd:     now.Add(time.Hour),
				},
			).
			Return(
				domain.CompClass{
					ID:          fakedCompClassID,
					Ownership:   fakedOwnership,
					ContestID:   fakedContestID,
					Name:        "Females",
					Description: "Female climbers",
					TimeBegin:   now,
					TimeEnd:     now.Add(time.Hour),
				}, nil)

		ucase := usecases.CompClassUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		compClass, err := ucase.CreateCompClass(context.Background(), fakedContestID, domain.CompClassTemplate{
			Name:        "Females",
			Description: "Female climbers",
			TimeBegin:   now,
			TimeEnd:     now.Add(time.Hour),
		})

		require.NoError(t, err)
		assert.Equal(t, fakedCompClassID, compClass.ID)
		assert.Equal(t, fakedOwnership, compClass.Ownership)
		assert.Equal(t, "Females", compClass.Name)
		assert.Equal(t, "Female climbers", compClass.Description)
		assert.Equal(t, now, compClass.TimeBegin)
		assert.Equal(t, now.Add(time.Hour), compClass.TimeEnd)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("ValidatorIsInvoked", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		ucase := usecases.CompClassUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		_, err := ucase.CreateCompClass(context.Background(), fakedContestID, domain.CompClassTemplate{})

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.True(t, validators.CompClassValidator{}.IsValidationError(err))

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		ucase := usecases.CompClassUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		_, err := ucase.CreateCompClass(context.Background(), fakedContestID, domain.CompClassTemplate{})

		require.ErrorIs(t, err, domain.ErrNoOwnership)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})
}

func TestDeleteCompClass(t *testing.T) {
	t.Parallel()

	fakedCompClassID := testutils.RandomResourceID[domain.CompClassID]()
	fakedOrganizerID := testutils.RandomResourceID[domain.OrganizerID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: fakedOrganizerID,
	}

	makeMocks := func() (*repositoryMock, *authorizerMock) {
		mockedRepo := new(repositoryMock)

		mockedRepo.
			On("GetCompClass", mock.Anything, nil, fakedCompClassID).
			Return(domain.CompClass{
				ID:        fakedCompClassID,
				Ownership: fakedOwnership,
			}, nil)

		mockedAuthorizer := new(authorizerMock)

		return mockedRepo, mockedAuthorizer
	}

	t.Run("HappyCase", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("GetContendersByCompClass", mock.Anything, nil, fakedCompClassID).
			Return([]domain.Contender{}, nil)

		mockedRepo.
			On("DeleteCompClass", mock.Anything, nil, fakedCompClassID).
			Return(nil)

		ucase := usecases.CompClassUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		err := ucase.DeleteCompClass(context.Background(), fakedCompClassID)

		require.NoError(t, err)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("CompClassHasContenders", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("GetContendersByCompClass", mock.Anything, nil, fakedCompClassID).
			Return([]domain.Contender{{}}, nil)

		ucase := usecases.CompClassUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		err := ucase.DeleteCompClass(context.Background(), fakedCompClassID)

		assert.ErrorIs(t, err, domain.ErrNotAllowed)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		ucase := usecases.CompClassUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		err := ucase.DeleteCompClass(context.Background(), fakedCompClassID)

		assert.ErrorIs(t, err, domain.ErrNoOwnership)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})
}
func TestPatchCompClass(t *testing.T) {
	t.Parallel()

	fakedCompClassID := testutils.RandomResourceID[domain.CompClassID]()
	fakedOrganizerID := testutils.RandomResourceID[domain.OrganizerID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: fakedOrganizerID,
	}

	now := time.Now()

	makeMocks := func() (*repositoryMock, *authorizerMock) {
		mockedRepo := new(repositoryMock)

		mockedRepo.
			On("GetCompClass", mock.Anything, nil, fakedCompClassID).
			Return(domain.CompClass{
				ID:        fakedCompClassID,
				Ownership: fakedOwnership,
			}, nil)

		mockedAuthorizer := new(authorizerMock)

		return mockedRepo, mockedAuthorizer
	}

	t.Run("HappyCase", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("StoreCompClass", mock.Anything, nil,
				domain.CompClass{
					ID:          fakedCompClassID,
					Ownership:   fakedOwnership,
					Name:        "Females",
					Description: "Female climbers",
					TimeBegin:   now,
					TimeEnd:     now.Add(time.Hour),
				},
			).
			Return(domain.CompClass{
				ID:          fakedCompClassID,
				Ownership:   fakedOwnership,
				Name:        "Females",
				Description: "Female climbers",
				TimeBegin:   now,
				TimeEnd:     now.Add(time.Hour),
			}, nil)

		ucase := usecases.CompClassUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		patch := domain.CompClassPatch{
			Name:        domain.NewPatch("Females"),
			Description: domain.NewPatch("Female climbers"),
			TimeBegin:   domain.NewPatch(now),
			TimeEnd:     domain.NewPatch(now.Add(time.Hour)),
		}

		compClass, err := ucase.PatchCompClass(context.Background(), fakedCompClassID, patch)

		require.NoError(t, err)
		assert.Equal(t, "Females", compClass.Name)
		assert.Equal(t, "Female climbers", compClass.Description)
		assert.Equal(t, now, compClass.TimeBegin)
		assert.Equal(t, now.Add(time.Hour), compClass.TimeEnd)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		ucase := usecases.CompClassUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		_, err := ucase.PatchCompClass(context.Background(), fakedCompClassID, domain.CompClassPatch{})

		assert.ErrorIs(t, err, domain.ErrNoOwnership)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("ValidatorIsInvoked", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		ucase := usecases.CompClassUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		_, err := ucase.PatchCompClass(context.Background(), fakedCompClassID, domain.CompClassPatch{})

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.True(t, validators.CompClassValidator{}.IsValidationError(err))

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})
}
