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

	t.Run("InvalidData", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("StoreCompClass", mock.Anything, nil, mock.AnythingOfType("domain.CompClass")).
			Return(domain.CompClass{}, nil)

		ucase := usecases.CompClassUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		validTemplate := func() domain.CompClassTemplate {
			return domain.CompClassTemplate{
				Name:        "Females",
				Description: "Female climbers",
				TimeBegin:   now,
				TimeEnd:     now.Add(time.Hour),
			}
		}

		_, err := ucase.CreateCompClass(context.Background(), fakedContestID, validTemplate())

		require.NoError(t, err)

		t.Run("EmptyName", func(t *testing.T) {
			tmpl := validTemplate()
			tmpl.Name = ""

			_, err := ucase.CreateCompClass(context.Background(), fakedContestID, tmpl)

			assert.ErrorIs(t, err, domain.ErrInvalidData)
		})

		t.Run("TimeEndBeforeTimeBegin", func(t *testing.T) {
			tmpl := validTemplate()
			tmpl.TimeEnd = tmpl.TimeBegin.Add(-time.Nanosecond)

			_, err := ucase.CreateCompClass(context.Background(), fakedContestID, tmpl)

			assert.ErrorIs(t, err, domain.ErrInvalidData)
		})

		t.Run("TotalDurationExceedingTwelveHours", func(t *testing.T) {
			tmpl := validTemplate()
			tmpl.TimeEnd = tmpl.TimeBegin.Add(12*time.Hour + time.Nanosecond)

			_, err := ucase.CreateCompClass(context.Background(), fakedContestID, tmpl)

			assert.ErrorIs(t, err, domain.ErrInvalidData)
		})

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
