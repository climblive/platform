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
					Ownership:        fakedOwnership,
					CompetitionID:    fakedCompetitionID,
					Name:             "Females",
					Description:      "Female climbers",
					EligibilityBasis: domain.EligibilityBasisAge,
					MinAge:           new(18),
					MaxAge:           new(99),
					TimeBegin:        now,
					TimeEnd:          now.Add(time.Hour),
				},
			).
			Return(
				domain.Category{
					ID:               fakedCategoryID,
					Ownership:        fakedOwnership,
					CompetitionID:    fakedCompetitionID,
					Name:             "Females",
					Description:      "Female climbers",
					EligibilityBasis: domain.EligibilityBasisAge,
					MinAge:           new(18),
					MaxAge:           new(99),
					TimeBegin:        now,
					TimeEnd:          now.Add(time.Hour),
				}, nil)

		ucase := usecases.CategoryUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		category, err := ucase.CreateCategory(context.Background(), fakedCompetitionID, domain.CategoryTemplate{
			Name:             "Females",
			EligibilityBasis: domain.EligibilityBasisAge,
			MinAge:           new(18),
			MaxAge:           new(99),
			Description:      "Female climbers",
			TimeBegin:        now,
			TimeEnd:          now.Add(time.Hour),
		})

		require.NoError(t, err)
		assert.Equal(t, fakedCategoryID, category.ID)
		assert.Equal(t, fakedOwnership, category.Ownership)
		assert.Equal(t, "Females", category.Name)
		assert.Equal(t, "Female climbers", category.Description)
		assert.Equal(t, domain.EligibilityBasisAge, category.EligibilityBasis)
		assert.Equal(t, new(18), category.MinAge)
		assert.Equal(t, new(99), category.MaxAge)
		assert.Equal(t, now, category.TimeBegin)
		assert.Equal(t, now.Add(time.Hour), category.TimeEnd)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	for _, tc := range []struct {
		name         string
		basis        domain.EligibilityBasis
		minAge       *int
		maxAge       *int
		minBirthYear *int
		maxBirthYear *int
		invalid      bool
	}{
		{name: "BothAgesMissing", basis: domain.EligibilityBasisAge},
		{name: "MinimumAgeMissing", basis: domain.EligibilityBasisAge, maxAge: new(99)},
		{name: "MaximumAgeMissing", basis: domain.EligibilityBasisAge, minAge: new(18)},
		{name: "BirthYearSpan", basis: domain.EligibilityBasisBirthYear, minBirthYear: new(2010), maxBirthYear: new(2014)},
		{name: "MinimumBirthYearMissing", basis: domain.EligibilityBasisBirthYear, maxBirthYear: new(2009)},
		{name: "MaximumBirthYearMissing", basis: domain.EligibilityBasisBirthYear, minBirthYear: new(2012)},
		{name: "BothBirthYearsMissing", basis: domain.EligibilityBasisBirthYear},
		{name: "MissingBasis", invalid: true},
		{name: "UnknownBasis", basis: "unknown", invalid: true},
		{name: "InactiveAgeBounds", basis: domain.EligibilityBasisBirthYear, minAge: new(18), invalid: true},
		{name: "InactiveBirthYearBounds", basis: domain.EligibilityBasisAge, maxBirthYear: new(2014), invalid: true},
		{name: "ReversedBirthYears", basis: domain.EligibilityBasisBirthYear, minBirthYear: new(2014), maxBirthYear: new(2010), invalid: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mockedRepo, mockedAuthorizer := makeMocks()
			mockedAuthorizer.
				On("HasOwnership", mock.Anything, fakedOwnership).
				Return(domain.OrganizerRole, nil)
			mockedRepo.
				On("GetCategoriesByCompetition", mock.Anything, nil, fakedCompetitionID).
				Return([]domain.Category{}, nil)

			expected := domain.Category{
				Ownership:        fakedOwnership,
				CompetitionID:    fakedCompetitionID,
				Name:             "Females",
				Description:      "Retained for rollback",
				EligibilityBasis: tc.basis,
				MinAge:           tc.minAge,
				MaxAge:           tc.maxAge,
				MinBirthYear:     tc.minBirthYear,
				MaxBirthYear:     tc.maxBirthYear,
				TimeBegin:        now,
				TimeEnd:          now.Add(time.Hour),
			}
			if !tc.invalid {
				mockedRepo.On("StoreCategory", mock.Anything, nil, expected).Return(expected, nil).Once()
			}

			ucase := usecases.CategoryUseCase{Repo: mockedRepo, Authorizer: mockedAuthorizer}
			result, err := ucase.CreateCategory(context.Background(), fakedCompetitionID, domain.CategoryTemplate{
				Name:             "Females",
				Description:      expected.Description,
				EligibilityBasis: tc.basis,
				MinAge:           tc.minAge,
				MaxAge:           tc.maxAge,
				MinBirthYear:     tc.minBirthYear,
				MaxBirthYear:     tc.maxBirthYear,
				TimeBegin:        now,
				TimeEnd:          now.Add(time.Hour),
			})

			if tc.invalid {
				require.ErrorIs(t, err, domain.ErrInvalidData)
				mockedRepo.AssertNotCalled(t, "StoreCategory", mock.Anything, mock.Anything, mock.Anything)
			} else {
				require.NoError(t, err)
				assert.Equal(t, expected, result)
			}
			mockedRepo.AssertExpectations(t)
			mockedAuthorizer.AssertExpectations(t)
		})
	}

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

		_, err := ucase.CreateCategory(context.Background(), fakedCompetitionID, domain.CategoryTemplate{EligibilityBasis: domain.EligibilityBasisAge, MinAge: new(18), MaxAge: new(99)})

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
			On("GetCompetitorsByCategory", mock.Anything, nil, fakedCategoryID).
			Return([]domain.Competitor{}, nil)

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

	t.Run("CategoryHasCompetitors", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("GetCompetitorsByCategory", mock.Anything, nil, fakedCategoryID).
			Return([]domain.Competitor{{}}, nil)

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
					ID:               fakedCategoryID,
					Ownership:        fakedOwnership,
					Name:             "Females",
					Description:      "Female climbers",
					EligibilityBasis: domain.EligibilityBasisAge,
					MinAge:           new(18),
					MaxAge:           new(99),
					TimeBegin:        now,
					TimeEnd:          now.Add(time.Hour),
				},
			).
			Return(domain.Category{
				ID:               fakedCategoryID,
				Ownership:        fakedOwnership,
				Name:             "Females",
				Description:      "Female climbers",
				EligibilityBasis: domain.EligibilityBasisAge,
				MinAge:           new(18),
				MaxAge:           new(99),
				TimeBegin:        now,
				TimeEnd:          now.Add(time.Hour),
			}, nil)

		ucase := usecases.CategoryUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		patch := domain.CategoryPatch{
			Name:             domain.NewPatch("Females"),
			EligibilityBasis: domain.NewPatch(domain.EligibilityBasisAge),
			Description:      domain.NewPatch("Female climbers"),
			MinAge:           domain.NewPatch(new(18)),
			MaxAge:           domain.NewPatch(new(99)),
			TimeBegin:        domain.NewPatch(now),
			TimeEnd:          domain.NewPatch(now.Add(time.Hour)),
		}

		category, err := ucase.PatchCategory(context.Background(), fakedCategoryID, patch)

		require.NoError(t, err)
		assert.Equal(t, "Females", category.Name)
		assert.Equal(t, "Female climbers", category.Description)
		assert.Equal(t, new(18), category.MinAge)
		assert.Equal(t, new(99), category.MaxAge)
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

	for _, tc := range []struct {
		name              string
		patch             domain.CategoryPatch
		basis             domain.EligibilityBasis
		minAge            *int
		maxAge            *int
		minBirthYear      *int
		maxBirthYear      *int
		birthYearCategory bool
		invalid           bool
	}{
		{name: "UnchangedEligibility", basis: domain.EligibilityBasisAge, minAge: new(18), maxAge: new(99)},
		{name: "ChangeMinimumAge", patch: domain.CategoryPatch{MinAge: domain.NewPatch(new(16))}, basis: domain.EligibilityBasisAge, minAge: new(16), maxAge: new(99)},
		{name: "ChangeMaximumAge", patch: domain.CategoryPatch{MaxAge: domain.NewPatch(new(65))}, basis: domain.EligibilityBasisAge, minAge: new(18), maxAge: new(65)},
		{name: "SingleAge", patch: domain.CategoryPatch{MaxAge: domain.NewPatch(new(18))}, basis: domain.EligibilityBasisAge, minAge: new(18), maxAge: new(18)},
		{name: "ZeroMinimumAge", patch: domain.CategoryPatch{MinAge: domain.NewPatch(new(0))}, basis: domain.EligibilityBasisAge, minAge: new(0), maxAge: new(99)},
		{name: "ClearMinimumAge", patch: domain.CategoryPatch{MinAge: domain.NewPatch[*int](nil)}, basis: domain.EligibilityBasisAge, maxAge: new(99)},
		{name: "ClearMaximumAge", patch: domain.CategoryPatch{MaxAge: domain.NewPatch[*int](nil)}, basis: domain.EligibilityBasisAge, minAge: new(18)},
		{name: "SwitchToBirthYear", patch: domain.CategoryPatch{EligibilityBasis: domain.NewPatch(domain.EligibilityBasisBirthYear), MinAge: domain.NewPatch[*int](nil), MaxAge: domain.NewPatch[*int](nil), MinBirthYear: domain.NewPatch(new(2010)), MaxBirthYear: domain.NewPatch(new(2014))}, basis: domain.EligibilityBasisBirthYear, minBirthYear: new(2010), maxBirthYear: new(2014)},
		{name: "SwitchToAge", birthYearCategory: true, patch: domain.CategoryPatch{EligibilityBasis: domain.NewPatch(domain.EligibilityBasisAge), MinAge: domain.NewPatch(new(18)), MaxAge: domain.NewPatch(new(99)), MinBirthYear: domain.NewPatch[*int](nil), MaxBirthYear: domain.NewPatch[*int](nil)}, basis: domain.EligibilityBasisAge, minAge: new(18), maxAge: new(99)},
		{name: "ChangeBirthYear", birthYearCategory: true, patch: domain.CategoryPatch{MinBirthYear: domain.NewPatch(new(2012))}, basis: domain.EligibilityBasisBirthYear, minBirthYear: new(2012), maxBirthYear: new(2014)},
		{name: "ClearMinimumBirthYear", birthYearCategory: true, patch: domain.CategoryPatch{MinBirthYear: domain.NewPatch[*int](nil)}, basis: domain.EligibilityBasisBirthYear, maxBirthYear: new(2014)},
		{name: "ClearMaximumBirthYear", birthYearCategory: true, patch: domain.CategoryPatch{MaxBirthYear: domain.NewPatch[*int](nil)}, basis: domain.EligibilityBasisBirthYear, minBirthYear: new(2010)},
		{name: "ReversedBirthYears", birthYearCategory: true, patch: domain.CategoryPatch{MaxBirthYear: domain.NewPatch(new(2009))}, invalid: true},
		{name: "SwitchWithoutClearingBounds", patch: domain.CategoryPatch{EligibilityBasis: domain.NewPatch(domain.EligibilityBasisBirthYear)}, invalid: true},
		{name: "InactiveBounds", patch: domain.CategoryPatch{MaxBirthYear: domain.NewPatch(new(2014))}, invalid: true},
		{name: "InvalidBasis", patch: domain.CategoryPatch{EligibilityBasis: domain.NewPatch[domain.EligibilityBasis]("invalid")}, invalid: true},
		{name: "ReversedAgeSpan", patch: domain.CategoryPatch{MaxAge: domain.NewPatch(new(15))}, invalid: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mockedRepo := new(repositoryMock)
			mockedAuthorizer := new(authorizerMock)
			category := domain.Category{
				ID:               fakedCategoryID,
				Ownership:        fakedOwnership,
				Name:             "Open",
				Description:      "Retained for rollback",
				EligibilityBasis: domain.EligibilityBasisAge,
				MinAge:           new(18),
				MaxAge:           new(99),
				TimeBegin:        now,
				TimeEnd:          now.Add(time.Hour),
			}
			if tc.birthYearCategory {
				category.EligibilityBasis = domain.EligibilityBasisBirthYear
				category.MinAge, category.MaxAge = nil, nil
				category.MinBirthYear, category.MaxBirthYear = new(2010), new(2014)
			}
			mockedRepo.
				On("GetCategory", mock.Anything, nil, fakedCategoryID).
				Return(category, nil)
			mockedAuthorizer.
				On("HasOwnership", mock.Anything, fakedOwnership).
				Return(domain.OrganizerRole, nil)

			updated := category
			updated.MinAge, updated.MaxAge = tc.minAge, tc.maxAge
			updated.EligibilityBasis = tc.basis
			updated.MinBirthYear, updated.MaxBirthYear = tc.minBirthYear, tc.maxBirthYear
			if !tc.invalid {
				mockedRepo.
					On("StoreCategory", mock.Anything, nil, updated).
					Return(updated, nil).Once()
			}

			ucase := usecases.CategoryUseCase{Repo: mockedRepo, Authorizer: mockedAuthorizer}
			result, err := ucase.PatchCategory(context.Background(), fakedCategoryID, tc.patch)
			if tc.invalid {
				require.ErrorIs(t, err, domain.ErrInvalidData)
				mockedRepo.AssertNotCalled(t, "StoreCategory", mock.Anything, mock.Anything, mock.Anything)
			} else {
				require.NoError(t, err)
				assert.Equal(t, updated, result)
				assert.Equal(t, category.Description, result.Description)
			}

			mockedRepo.AssertExpectations(t)
			mockedAuthorizer.AssertExpectations(t)
		})
	}

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
