package usecases_test

import (
	"context"
	"testing"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/usecases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateCategoryRequiresAgeSpan(t *testing.T) {
	for _, tc := range []struct {
		name   string
		minAge *int
		maxAge *int
	}{
		{name: "both missing"},
		{name: "minimum missing", maxAge: new(99)},
		{name: "maximum missing", minAge: new(18)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repo := new(repositoryMock)
			auth := new(authorizerMock)
			repo.On("GetCompetition", mock.Anything, nil, domain.CompetitionID(1)).Return(domain.Competition{}, nil)
			auth.On("HasOwnership", mock.Anything, domain.OwnershipData{}).Return(domain.OrganizerRole, nil)
			repo.On("GetCategoriesByCompetition", mock.Anything, nil, domain.CompetitionID(1)).Return([]domain.Category{}, nil)
			uc := usecases.CategoryUseCase{Repo: repo, Authorizer: auth}
			_, err := uc.CreateCategory(context.Background(), 1, domain.CategoryTemplate{MinAge: tc.minAge, MaxAge: tc.maxAge})
			require.ErrorIs(t, err, domain.ErrInvalidData)
			repo.AssertNotCalled(t, "StoreCategory", mock.Anything, mock.Anything, mock.Anything)
			repo.AssertExpectations(t)
			auth.AssertExpectations(t)
		})
	}
}

func TestPatchCategoryAgeSpan(t *testing.T) {
	for _, tc := range []struct {
		name    string
		patch   domain.CategoryPatch
		minAge  int
		maxAge  int
		invalid bool
	}{
		{name: "change minimum", patch: domain.CategoryPatch{MinAge: domain.NewPatch(new(16))}, minAge: 16, maxAge: 99},
		{name: "change maximum", patch: domain.CategoryPatch{MaxAge: domain.NewPatch(new(65))}, minAge: 18, maxAge: 65},
		{name: "single age", patch: domain.CategoryPatch{MaxAge: domain.NewPatch(new(18))}, minAge: 18, maxAge: 18},
		{name: "zero minimum", patch: domain.CategoryPatch{MinAge: domain.NewPatch(new(0))}, minAge: 0, maxAge: 99},
		{name: "clear minimum", patch: domain.CategoryPatch{MinAge: domain.NewPatch[*int](nil)}, invalid: true},
		{name: "clear maximum", patch: domain.CategoryPatch{MaxAge: domain.NewPatch[*int](nil)}, invalid: true},
		{name: "reversed span", patch: domain.CategoryPatch{MaxAge: domain.NewPatch(new(15))}, invalid: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repo := new(repositoryMock)
			auth := new(authorizerMock)
			now := time.Now()
			category := domain.Category{ID: 3, Name: "Open", Description: "Retained for rollback", MinAge: 18, MaxAge: 99, TimeBegin: now, TimeEnd: now.Add(time.Hour)}
			repo.On("GetCategory", mock.Anything, nil, category.ID).Return(category, nil)
			auth.On("HasOwnership", mock.Anything, category.Ownership).Return(domain.OrganizerRole, nil)
			updated := category
			updated.MinAge, updated.MaxAge = tc.minAge, tc.maxAge
			if !tc.invalid {
				repo.On("StoreCategory", mock.Anything, nil, updated).Return(updated, nil).Once()
			}
			uc := usecases.CategoryUseCase{Repo: repo, Authorizer: auth}
			result, err := uc.PatchCategory(context.Background(), category.ID, tc.patch)
			if tc.invalid {
				require.ErrorIs(t, err, domain.ErrInvalidData)
				repo.AssertNotCalled(t, "StoreCategory", mock.Anything, mock.Anything, mock.Anything)
			} else {
				require.NoError(t, err)
				assert.Equal(t, updated, result)
				assert.Equal(t, category.Description, result.Description)
			}
			repo.AssertExpectations(t)
			auth.AssertExpectations(t)
		})
	}
}
