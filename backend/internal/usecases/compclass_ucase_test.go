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

func TestGetCompClassesByContest(t *testing.T) {
	t.Parallel()

	fakedContestID := randomResourceID[domain.ContestID]()

	fakedCompClasses := []domain.CompClass{{
		ID:        randomResourceID[domain.CompClassID](),
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
	fakedOrganizerID := randomResourceID[domain.OrganizerID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: fakedOrganizerID,
	}
	fakedContestID := randomResourceID[domain.ContestID]()
	fakedCompClassID := randomResourceID[domain.CompClassID]()

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
		assert.True(t, usecases.CompClassValidator{}.IsValidationError(err))

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

	fakedCompClassID := randomResourceID[domain.CompClassID]()
	fakedOrganizerID := randomResourceID[domain.OrganizerID]()
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

func TestCompClassValidator(t *testing.T) {
	now := time.Now()

	validator := usecases.CompClassValidator{}

	validCompClass := func() domain.CompClass {
		return domain.CompClass{
			Name:        "Females",
			Description: "Female climbers",
			TimeBegin:   now,
			TimeEnd:     now.Add(time.Hour),
		}
	}

	t.Run("ValidData", func(t *testing.T) {
		err := validator.Validate(validCompClass())
		assert.NoError(t, err)
	})

	t.Run("EmptyName", func(t *testing.T) {
		compClass := validCompClass()
		compClass.Name = ""

		err := validator.Validate(compClass)

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.True(t, validator.IsValidationError(err))
	})

	t.Run("TimeEndBeforeTimeBegin", func(t *testing.T) {
		compClass := validCompClass()
		compClass.TimeEnd = compClass.TimeBegin.Add(-time.Nanosecond)

		err := validator.Validate(compClass)

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.True(t, validator.IsValidationError(err))
	})

	t.Run("TotalDurationExceedingTwelveHours", func(t *testing.T) {
		compClass := validCompClass()
		compClass.TimeEnd = compClass.TimeBegin.Add(12*time.Hour + time.Nanosecond)

		err := validator.Validate(compClass)

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.True(t, validator.IsValidationError(err))
	})
}
