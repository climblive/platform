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
	mockedContenderID := randomResourceID[domain.ContenderID]()
	mockedOwnership := domain.OwnershipData{
		OrganizerID: randomResourceID[domain.OrganizerID](),
		ContenderID: &mockedContenderID,
	}

	makeMocks := func() (*repositoryMock, *authorizerMock) {
		mockedRepo := new(repositoryMock)
		mockedAuthorizer := new(authorizerMock)

		mockedContender := domain.Contender{
			ID:        mockedContenderID,
			Ownership: mockedOwnership,
		}

		mockedRepo.
			On("GetContender", mock.Anything, mock.Anything, mockedContenderID).
			Return(mockedContender, nil)

		return mockedRepo, mockedAuthorizer
	}

	t.Run("HappyPath", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

		mockedTicks := []domain.Tick{
			{
				ID: randomResourceID[domain.TickID](),
			},
		}

		mockedRepo.
			On("GetTicksByContender", mock.Anything, mock.Anything, mockedContenderID).
			Return(mockedTicks, nil)

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

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

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

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})
}

func TestCreateTick(t *testing.T) {
	mockedContenderID := randomResourceID[domain.ContenderID]()
	mockedContestID := randomResourceID[domain.ContestID]()
	mockedCompClassID := randomResourceID[domain.CompClassID]()
	mockedProblemID := randomResourceID[domain.ProblemID]()

	gracePeriod := 15 * time.Minute

	mockedOwnership := domain.OwnershipData{
		OrganizerID: randomResourceID[domain.OrganizerID](),
		ContenderID: &mockedContenderID,
	}

	makeMocks := func(timeEnd time.Time) (*repositoryMock, *eventBrokerMock) {
		mockedRepo := new(repositoryMock)
		mockedEventBroker := new(eventBrokerMock)

		mockedContender := domain.Contender{
			ID:          mockedContenderID,
			Ownership:   mockedOwnership,
			ContestID:   mockedContestID,
			CompClassID: mockedCompClassID,
		}

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

		return mockedRepo, mockedEventBroker
	}

	t.Run("HappyPath", func(t *testing.T) {
		mockedRepo, mockedEventBroker := makeMocks(time.Now())
		mockedAuthorizer := new(authorizerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mockedOwnership).
			Return(domain.ContenderRole, nil)

		mockedRepo.
			On("StoreTick", mock.Anything, mock.Anything, mock.AnythingOfType("domain.Tick")).
			Return(mirrorInstruction{}, nil)

		mockedEventBroker.On("Dispatch", mockedContestID, domain.AscentRegisteredEvent{
			ContenderID:  mockedContenderID,
			ProblemID:    mockedProblemID,
			Top:          true,
			AttemptsTop:  5,
			Zone:         true,
			AttemptsZone: 2,
		}).Return()

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

		assert.Equal(t, mockedOwnership, tick.Ownership)
		assert.WithinDuration(t, time.Now(), tick.Timestamp, time.Minute)
		assert.Equal(t, mockedContestID, tick.ContestID)
		assert.Equal(t, mockedProblemID, tick.ProblemID)
		assert.Equal(t, true, tick.Top)
		assert.Equal(t, 5, tick.AttemptsTop)
		assert.Equal(t, true, tick.Zone)
		assert.Equal(t, 2, tick.AttemptsZone)

		mockedRepo.AssertExpectations(t)
		mockedEventBroker.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("ContenderCannotRegisterAscentAfterGracePeriod", func(t *testing.T) {
		mockedRepo, mockedEventBroker := makeMocks(time.Now().Add(-1 * gracePeriod))
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

		mockedRepo.AssertExpectations(t)
		mockedEventBroker.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("OrganizerCanRegisterAscentAfterGracePeriod", func(t *testing.T) {
		mockedRepo, mockedEventBroker := makeMocks(time.Now().Add(-1 * gracePeriod))
		mockedAuthorizer := new(authorizerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mockedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("StoreTick", mock.Anything, mock.Anything, mock.AnythingOfType("domain.Tick")).
			Return(mirrorInstruction{}, nil)

		mockedEventBroker.On("Dispatch", mockedContestID, domain.AscentRegisteredEvent{
			ContenderID:  mockedContenderID,
			ProblemID:    mockedProblemID,
			Top:          true,
			AttemptsTop:  5,
			Zone:         true,
			AttemptsZone: 2,
		}).Return()

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

		mockedRepo.AssertExpectations(t)
		mockedEventBroker.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedRepo := new(repositoryMock)
		mockedAuthorizer := new(authorizerMock)

		mockedContender := domain.Contender{
			ID:        mockedContenderID,
			Ownership: mockedOwnership,
		}

		mockedRepo.
			On("GetContender", mock.Anything, mock.Anything, mockedContenderID).
			Return(mockedContender, nil)

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

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})
}

func TestDeleteTick(t *testing.T) {
	mockedTickID := randomResourceID[domain.TickID]()
	mockedContenderID := randomResourceID[domain.ContenderID]()
	mockedContestID := randomResourceID[domain.ContestID]()
	mockedCompClassID := randomResourceID[domain.CompClassID]()
	mockedProblemID := randomResourceID[domain.ProblemID]()

	gracePeriod := 15 * time.Minute

	mockedOwnership := domain.OwnershipData{
		OrganizerID: randomResourceID[domain.OrganizerID](),
		ContenderID: &mockedContenderID,
	}

	makeMocks := func(timeEnd time.Time) (*repositoryMock, *eventBrokerMock) {
		mockedRepo := new(repositoryMock)
		mockedEventBroker := new(eventBrokerMock)

		mockedContender := domain.Contender{
			ID:          mockedContenderID,
			ContestID:   mockedContestID,
			CompClassID: mockedCompClassID,
		}

		mockedTick := domain.Tick{
			ID:        mockedTickID,
			Ownership: mockedOwnership,
			ProblemID: mockedProblemID,
			ContestID: mockedContestID,
		}

		mockedRepo.
			On("GetTick", mock.Anything, mock.Anything, mockedTickID).
			Return(mockedTick, nil)

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

		return mockedRepo, mockedEventBroker
	}

	t.Run("HappyPath", func(t *testing.T) {
		mockedRepo, mockedEventBroker := makeMocks(time.Now())
		mockedAuthorizer := new(authorizerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mockedOwnership).
			Return(domain.ContenderRole, nil)

		mockedRepo.
			On("DeleteTick", mock.Anything, mock.Anything, mockedTickID).
			Return(nil)

		mockedEventBroker.On("Dispatch", mockedContestID, domain.AscentDeregisteredEvent{
			ContenderID: mockedContenderID,
			ProblemID:   mockedProblemID,
		}).Return()

		ucase := usecases.TickUseCase{
			Repo:        mockedRepo,
			Authorizer:  mockedAuthorizer,
			EventBroker: mockedEventBroker,
		}

		err := ucase.DeleteTick(context.Background(), mockedTickID)

		require.NoError(t, err)

		mockedRepo.AssertExpectations(t)
		mockedEventBroker.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("ContenderCannotDeregisterAscentAfterGracePeriod", func(t *testing.T) {
		mockedRepo, mockedEventBroker := makeMocks(time.Now().Add(-1 * gracePeriod))
		mockedAuthorizer := new(authorizerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mockedOwnership).
			Return(domain.ContenderRole, nil)

		ucase := usecases.TickUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		err := ucase.DeleteTick(context.Background(), mockedTickID)

		assert.ErrorIs(t, err, domain.ErrContestEnded)

		mockedRepo.AssertExpectations(t)
		mockedEventBroker.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("OrganizerCanDeregisterAscentAfterGracePeriod", func(t *testing.T) {
		mockedRepo, mockedEventBroker := makeMocks(time.Now().Add(-1 * gracePeriod))
		mockedAuthorizer := new(authorizerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mockedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("DeleteTick", mock.Anything, mock.Anything, mockedTickID).
			Return(nil)

		mockedEventBroker.On("Dispatch", mockedContestID, domain.AscentDeregisteredEvent{
			ContenderID: mockedContenderID,
			ProblemID:   mockedProblemID,
		}).Return()

		ucase := usecases.TickUseCase{
			Repo:        mockedRepo,
			Authorizer:  mockedAuthorizer,
			EventBroker: mockedEventBroker,
		}

		err := ucase.DeleteTick(context.Background(), mockedTickID)

		require.NoError(t, err)

		mockedRepo.AssertExpectations(t)
		mockedEventBroker.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedRepo := new(repositoryMock)
		mockedAuthorizer := new(authorizerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mockedOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		mockedRepo.
			On("GetTick", mock.Anything, mock.Anything, mockedTickID).
			Return(domain.Tick{
				ID:        mockedTickID,
				Ownership: mockedOwnership,
			}, nil)

		ucase := usecases.TickUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		err := ucase.DeleteTick(context.Background(), mockedTickID)

		assert.ErrorIs(t, err, domain.ErrNoOwnership)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})
}
