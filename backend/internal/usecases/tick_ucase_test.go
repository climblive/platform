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

func TestGetTicksByContender(t *testing.T) {
	mockedContenderID := domain.ResourceID(1)
	mockedOwnership := domain.OwnershipData{
		OrganizerID: 1,
		ContenderID: &mockedContenderID,
	}

	mockedContender := domain.Contender{
		ID:        mockedContenderID,
		Ownership: mockedOwnership,
	}

	mockedTicks := []domain.Tick{
		{
			ID: 1,
		},
	}

	mockedRepo := new(repositoryMock)

	mockedRepo.
		On("GetContender", mock.Anything, mock.Anything, mockedContenderID).
		Return(mockedContender, nil)

	mockedRepo.
		On("GetTicksByContender", mock.Anything, mock.Anything, mockedContenderID).
		Return(mockedTicks, nil)

	t.Run("HappyPath", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mockedOwnership).
			Return(domain.ContenderRole, nil)

		ucase := usecases.TickUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		ticks, err := ucase.GetTicksByContender(context.Background(), mockedContenderID)

		require.NoError(t, err)
		assert.Equal(t, mockedTicks, ticks)
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mockedOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		ucase := usecases.TickUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		ticks, err := ucase.GetTicksByContender(context.Background(), mockedContenderID)

		assert.ErrorIs(t, err, domain.ErrNoOwnership)
		assert.Nil(t, ticks)
	})
}

func TestCreateTick(t *testing.T) {
	mockedContenderID := randomResourceID()
	mockedContestID := randomResourceID()
	mockedCompClassID := randomResourceID()
	mockedProblemID := randomResourceID()

	gracePeriod := 15 * time.Minute

	mockedOwnership := domain.OwnershipData{
		OrganizerID: 1,
		ContenderID: &mockedContenderID,
	}

	mockedContender := domain.Contender{
		ID:          mockedContenderID,
		Ownership:   mockedOwnership,
		ContestID:   mockedContestID,
		CompClassID: mockedCompClassID,
	}

	mockedEventBroker := new(eventBrokerMock)

	makeMockedRepo := func(timeEnd time.Time) *repositoryMock {
		mockedRepo := new(repositoryMock)

		mockedRepo.
			On("GetContender", mock.Anything, mock.Anything, mockedContenderID).
			Return(mockedContender, nil)

		mockedRepo.
			On("GetContest", mock.Anything, mock.Anything, mockedContestID).
			Return(domain.Contest{
				ID:          mockedContestID,
				GracePeriod: gracePeriod,
			}, nil)

		mockedRepo.
			On("GetCompClass", mock.Anything, mock.Anything, mockedCompClassID).
			Return(domain.CompClass{
				ID:      mockedCompClassID,
				TimeEnd: timeEnd,
			}, nil)

		mockedRepo.
			On("GetProblem", mock.Anything, mock.Anything, mockedProblemID).
			Return(domain.Problem{
				ID: mockedProblemID,
			}, nil)

		mockedRepo.
			On("StoreTick", mock.Anything, mock.Anything, mock.AnythingOfType("domain.Tick")).
			Return(mirrorInstruction{}, nil)

		mockedEventBroker.On("Dispatch", mockedContestID, mock.Anything).Return()

		return mockedRepo
	}

	t.Run("HappyPath", func(t *testing.T) {
		mockedRepo := makeMockedRepo(time.Now())
		mockedAuthorizer := new(authorizerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mockedOwnership).
			Return(domain.ContenderRole, nil)

		ucase := usecases.TickUseCase{
			Repo:        mockedRepo,
			Authorizer:  mockedAuthorizer,
			EventBroker: mockedEventBroker,
		}

		tick, err := ucase.CreateTick(context.Background(), mockedContenderID, domain.Tick{
			ProblemID:    mockedProblemID,
			Top:          true,
			AttemptsTop:  5,
			Zone:         true,
			AttemptsZone: 2,
		})

		require.NoError(t, err)

		mockedRepo.AssertExpectations(t)

		assert.Equal(t, mockedOwnership, tick.Ownership)
		assert.WithinDuration(t, time.Now(), tick.Timestamp, time.Minute)
		assert.Equal(t, mockedContestID, tick.ContestID)
		assert.Equal(t, mockedProblemID, tick.ProblemID)
		assert.Equal(t, true, tick.Top)
		assert.Equal(t, 5, tick.AttemptsTop)
		assert.Equal(t, true, tick.Zone)
		assert.Equal(t, 2, tick.AttemptsZone)

		mockedEventBroker.AssertCalled(t, "Dispatch", mockedContestID, domain.AscentRegisteredEvent{
			ContenderID:  mockedContenderID,
			ProblemID:    mockedProblemID,
			Top:          true,
			AttemptsTop:  5,
			Zone:         true,
			AttemptsZone: 2,
		})
	})

	t.Run("ContenderCannotRegisterAscentAfterGracePeriod", func(t *testing.T) {
		mockedRepo := makeMockedRepo(time.Now().Add(-1 * gracePeriod))
		mockedAuthorizer := new(authorizerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mockedOwnership).
			Return(domain.ContenderRole, nil)

		ucase := usecases.TickUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		tick, err := ucase.CreateTick(context.Background(), mockedContenderID, domain.Tick{
			ProblemID:    mockedProblemID,
			Top:          true,
			AttemptsTop:  5,
			Zone:         true,
			AttemptsZone: 2,
		})

		assert.ErrorIs(t, err, domain.ErrContestEnded)
		assert.Empty(t, tick)

	})

	t.Run("OrganizerCanRegisterAscentAfterGracePeriod", func(t *testing.T) {
		mockedRepo := makeMockedRepo(time.Now().Add(-1 * gracePeriod))
		mockedAuthorizer := new(authorizerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mockedOwnership).
			Return(domain.OrganizerRole, nil)

		ucase := usecases.TickUseCase{
			Repo:        mockedRepo,
			Authorizer:  mockedAuthorizer,
			EventBroker: mockedEventBroker,
		}

		tick, err := ucase.CreateTick(context.Background(), mockedContenderID, domain.Tick{
			ProblemID:    mockedProblemID,
			Top:          true,
			AttemptsTop:  5,
			Zone:         true,
			AttemptsZone: 2,
		})

		require.NoError(t, err)
		assert.NotEmpty(t, tick)
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedRepo := makeMockedRepo(time.Now())
		mockedAuthorizer := new(authorizerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mockedOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		ucase := usecases.TickUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		tick, err := ucase.CreateTick(context.Background(), mockedContenderID, domain.Tick{})

		assert.ErrorIs(t, err, domain.ErrNoOwnership)
		assert.Empty(t, tick)
	})
}
