package usecases_test

import (
	"context"
	"testing"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/testutils"
	"github.com/climblive/platform/backend/internal/usecases"
	"github.com/climblive/platform/backend/internal/usecases/validators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetCategory(t *testing.T) {
	t.Parallel()

	fakedCategoryID := testutils.RandomResourceID[domain.CategoryID]()
	fakedOrganizerID := testutils.RandomResourceID[domain.OrganizerID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: fakedOrganizerID,
	}

	fakedCategory := domain.Category{
		ID:        fakedCategoryID,
		Ownership: fakedOwnership,
	}

	mockedRepo := new(repositoryMock)

	mockedRepo.
		On("GetCategory", mock.Anything, nil, fakedCategoryID).
		Return(fakedCategory, nil)

	ucase := usecases.CategoryUseCase{
		Repo: mockedRepo,
	}

	category, err := ucase.GetCategory(context.Background(), fakedCategoryID)

	require.NoError(t, err)
	assert.Equal(t, fakedCategory, category)

	mockedRepo.AssertExpectations(t)
}

func TestGetCategoriesByCompetition(t *testing.T) {
	t.Parallel()

	fakedCompetitionID := testutils.RandomResourceID[domain.CompetitionID]()

	fakedCategories := []domain.Category{{
		ID:            testutils.RandomResourceID[domain.CategoryID](),
		CompetitionID: fakedCompetitionID,
	},
	}

	mockedRepo := new(repositoryMock)

	mockedRepo.
		On("GetCategoriesByCompetition", mock.Anything, mock.Anything, fakedCompetitionID).
		Return(fakedCategories, nil)

	ucase := usecases.CategoryUseCase{
		Repo: mockedRepo,
	}

	categories, err := ucase.GetCategoriesByCompetition(context.Background(), fakedCompetitionID)

	require.NoError(t, err)
	assert.Equal(t, fakedCategories, categories)

	mockedRepo.AssertExpectations(t)
}

func TestCreateCategory(t *testing.T) {
	fakedOrganizerID := testutils.RandomResourceID[domain.OrganizerID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: fakedOrganizerID,
	}
	fakedCompetitionID := testutils.RandomResourceID[domain.CompetitionID]()
	fakedCategoryID := testutils.RandomResourceID[domain.CategoryID]()

	now := time.Now()

	makeMocks := func() (*repositoryMock, *authorizerMock) {
		mockedRepo := new(repositoryMock)

		mockedRepo.
			On("GetCompetition", mock.Anything, nil, fakedCompetitionID).
			Return(domain.Competition{
				ID:        fakedCompetitionID,
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
			On("GetCategoriesByCompetition", mock.Anything, nil, fakedCompetitionID).
			Return([]domain.Category{}, nil)

		mockedRepo.
			On("StoreCategory", mock.Anything, nil,
				domain.Category{
					Ownership:     fakedOwnership,
					CompetitionID: fakedCompetitionID,
					Name:          "Females",
					Description:   "Female climbers",
					TimeBegin:     now,
					TimeEnd:       now.Add(time.Hour),
				},
			).
			Return(
				domain.Category{
					ID:            fakedCategoryID,
					Ownership:     fakedOwnership,
					CompetitionID: fakedCompetitionID,
					Name:          "Females",
					Description:   "Female climbers",
					TimeBegin:     now,
					TimeEnd:       now.Add(time.Hour),
				}, nil)

		ucase := usecases.CategoryUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		category, err := ucase.CreateCategory(context.Background(), fakedCompetitionID, domain.CategoryTemplate{
			Name:        "Females",
			Description: "Female climbers",
			TimeBegin:   now,
			TimeEnd:     now.Add(time.Hour),
		})

		require.NoError(t, err)
		assert.Equal(t, fakedCategoryID, category.ID)
		assert.Equal(t, fakedOwnership, category.Ownership)
		assert.Equal(t, "Females", category.Name)
		assert.Equal(t, "Female climbers", category.Description)
		assert.Equal(t, now, category.TimeBegin)
		assert.Equal(t, now.Add(time.Hour), category.TimeEnd)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("ValidatorIsInvoked", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("GetCategoriesByCompetition", mock.Anything, nil, fakedCompetitionID).
			Return([]domain.Category{}, nil)

		ucase := usecases.CategoryUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		_, err := ucase.CreateCategory(context.Background(), fakedCompetitionID, domain.CategoryTemplate{})

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.True(t, validators.CategoryValidator{}.IsValidationError(err))

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("LimitExceeded", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		categories := make([]domain.Category, 20)
		mockedRepo.
			On("GetCategoriesByCompetition", mock.Anything, nil, fakedCompetitionID).
			Return(categories, nil)

		ucase := usecases.CategoryUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		_, err := ucase.CreateCategory(context.Background(), fakedCompetitionID, domain.CategoryTemplate{})

		require.ErrorIs(t, err, domain.ErrLimitExceeded)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		ucase := usecases.CategoryUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		_, err := ucase.CreateCategory(context.Background(), fakedCompetitionID, domain.CategoryTemplate{})

		require.ErrorIs(t, err, domain.ErrNoOwnership)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})
}

func TestDeleteCategory(t *testing.T) {
	t.Parallel()

	fakedCategoryID := testutils.RandomResourceID[domain.CategoryID]()
	fakedOrganizerID := testutils.RandomResourceID[domain.OrganizerID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: fakedOrganizerID,
	}

	makeMocks := func() (*repositoryMock, *authorizerMock) {
		mockedRepo := new(repositoryMock)

		mockedRepo.
			On("GetCategory", mock.Anything, nil, fakedCategoryID).
			Return(domain.Category{
				ID:        fakedCategoryID,
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
			On("GetContendersByCategory", mock.Anything, nil, fakedCategoryID).
			Return([]domain.Contender{}, nil)

		mockedRepo.
			On("DeleteCategory", mock.Anything, nil, fakedCategoryID).
			Return(nil)

		ucase := usecases.CategoryUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		err := ucase.DeleteCategory(context.Background(), fakedCategoryID)

		require.NoError(t, err)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("CategoryHasContenders", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("GetContendersByCategory", mock.Anything, nil, fakedCategoryID).
			Return([]domain.Contender{{}}, nil)

		ucase := usecases.CategoryUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		err := ucase.DeleteCategory(context.Background(), fakedCategoryID)

		assert.ErrorIs(t, err, domain.ErrNotAllowed)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		ucase := usecases.CategoryUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		err := ucase.DeleteCategory(context.Background(), fakedCategoryID)

		assert.ErrorIs(t, err, domain.ErrNoOwnership)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})
}
func TestPatchCategory(t *testing.T) {
	t.Parallel()

	fakedCategoryID := testutils.RandomResourceID[domain.CategoryID]()
	fakedOrganizerID := testutils.RandomResourceID[domain.OrganizerID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: fakedOrganizerID,
	}

	now := time.Now()

	makeMocks := func() (*repositoryMock, *authorizerMock) {
		mockedRepo := new(repositoryMock)

		mockedRepo.
			On("GetCategory", mock.Anything, nil, fakedCategoryID).
			Return(domain.Category{
				ID:        fakedCategoryID,
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
			On("StoreCategory", mock.Anything, nil,
				domain.Category{
					ID:          fakedCategoryID,
					Ownership:   fakedOwnership,
					Name:        "Females",
					Description: "Female climbers",
					TimeBegin:   now,
					TimeEnd:     now.Add(time.Hour),
				},
			).
			Return(domain.Category{
				ID:          fakedCategoryID,
				Ownership:   fakedOwnership,
				Name:        "Females",
				Description: "Female climbers",
				TimeBegin:   now,
				TimeEnd:     now.Add(time.Hour),
			}, nil)

		ucase := usecases.CategoryUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		patch := domain.CategoryPatch{
			Name:        domain.NewPatch("Females"),
			Description: domain.NewPatch("Female climbers"),
			TimeBegin:   domain.NewPatch(now),
			TimeEnd:     domain.NewPatch(now.Add(time.Hour)),
		}

		category, err := ucase.PatchCategory(context.Background(), fakedCategoryID, patch)

		require.NoError(t, err)
		assert.Equal(t, "Females", category.Name)
		assert.Equal(t, "Female climbers", category.Description)
		assert.Equal(t, now, category.TimeBegin)
		assert.Equal(t, now.Add(time.Hour), category.TimeEnd)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		ucase := usecases.CategoryUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		_, err := ucase.PatchCategory(context.Background(), fakedCategoryID, domain.CategoryPatch{})

		assert.ErrorIs(t, err, domain.ErrNoOwnership)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("ValidatorIsInvoked", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		ucase := usecases.CategoryUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		_, err := ucase.PatchCategory(context.Background(), fakedCategoryID, domain.CategoryPatch{})

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.True(t, validators.CategoryValidator{}.IsValidationError(err))

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})
}
