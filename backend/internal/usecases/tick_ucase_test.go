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
	fakedContenderID := randomResourceID[domain.ContenderID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: randomResourceID[domain.OrganizerID](),
		ContenderID: &fakedContenderID,
	}

	makeMocks := func() (*repositoryMock, *authorizerMock) {
		mockedRepo := new(repositoryMock)
		mockedAuthorizer := new(authorizerMock)

		fakedContender := domain.Contender{
			ID:        fakedContenderID,
			Ownership: fakedOwnership,
		}

		mockedRepo.
			On("GetContender", mock.Anything, mock.Anything, fakedContenderID).
			Return(fakedContender, nil)

		return mockedRepo, mockedAuthorizer
	}

	t.Run("HappyPath", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

		fakedTicks := []domain.Tick{
			{
				ID: randomResourceID[domain.TickID](),
			},
		}

		mockedRepo.
			On("GetTicksByContender", mock.Anything, mock.Anything, fakedContenderID).
			Return(fakedTicks, nil)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.ContenderRole, nil)

		ucase := usecases.TickUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		ticks, err := ucase.GetTicksByContender(context.Background(), fakedContenderID)

		require.NoError(t, err)
		assert.Equal(t, fakedTicks, ticks)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		ucase := usecases.TickUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		ticks, err := ucase.GetTicksByContender(context.Background(), fakedContenderID)

		assert.ErrorIs(t, err, domain.ErrNoOwnership)
		assert.Nil(t, ticks)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})
}

func TestCreateTick(t *testing.T) {
	fakedContenderID := randomResourceID[domain.ContenderID]()
	fakedContestID := randomResourceID[domain.ContestID]()
	fakedCompClassID := randomResourceID[domain.CompClassID]()
	fakedProblemID := randomResourceID[domain.ProblemID]()

	now := time.Now()
	gracePeriod := 15 * time.Minute

	fakedOwnership := domain.OwnershipData{
		OrganizerID: randomResourceID[domain.OrganizerID](),
		ContenderID: &fakedContenderID,
	}

	makeMocks := func(timeEnd time.Time) (*repositoryMock, *eventBrokerMock) {
		mockedRepo := new(repositoryMock)
		mockedEventBroker := new(eventBrokerMock)

		fakedContender := domain.Contender{
			ID:          fakedContenderID,
			Ownership:   fakedOwnership,
			ContestID:   fakedContestID,
			CompClassID: fakedCompClassID,
		}

		mockedRepo.
			On("GetContender", mock.Anything, mock.Anything, fakedContenderID).
			Return(fakedContender, nil)

		mockedRepo.
			On("GetContest", mock.Anything, mock.Anything, fakedContestID).
			Return(domain.Contest{
				ID:          fakedContestID,
				GracePeriod: gracePeriod,
			}, nil)

		mockedRepo.
			On("GetCompClass", mock.Anything, mock.Anything, fakedCompClassID).
			Return(domain.CompClass{
				ID:      fakedCompClassID,
				TimeEnd: timeEnd,
			}, nil)

		return mockedRepo, mockedEventBroker
	}

	t.Run("HappyPath", func(t *testing.T) {
		mockedRepo, mockedEventBroker := makeMocks(time.Now())
		mockedAuthorizer := new(authorizerMock)

		fakedTickID := randomResourceID[domain.TickID]()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.ContenderRole, nil)

		mockedRepo.
			On("GetProblem", mock.Anything, mock.Anything, fakedProblemID).
			Return(domain.Problem{
				ID:        fakedProblemID,
				ContestID: fakedContestID,
			}, nil)

		mockedRepo.
			On("StoreTick", mock.Anything, nil, mock.MatchedBy(func(tick domain.Tick) bool {
				tick.Timestamp = time.Time{}

				expected := domain.Tick{
					Ownership:    fakedOwnership,
					ContestID:    fakedContestID,
					ProblemID:    fakedProblemID,
					Top:          true,
					AttemptsTop:  5,
					Zone:         true,
					AttemptsZone: 2,
				}

				return tick.Timestamp.Sub(now) < time.Second && tick == expected
			})).
			Return(domain.Tick{
				ID:           fakedTickID,
				Ownership:    fakedOwnership,
				Timestamp:    now,
				ContestID:    fakedContestID,
				ProblemID:    fakedProblemID,
				Top:          true,
				AttemptsTop:  5,
				Zone:         true,
				AttemptsZone: 2,
			}, nil)

		mockedEventBroker.On("Dispatch", fakedContestID, domain.AscentRegisteredEvent{
			TickID:       fakedTickID,
			Timestamp:    now,
			ContenderID:  fakedContenderID,
			ProblemID:    fakedProblemID,
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

		tick, err := ucase.CreateTick(context.Background(), fakedContenderID, domain.Tick{
			ProblemID:    fakedProblemID,
			Top:          true,
			AttemptsTop:  5,
			Zone:         true,
			AttemptsZone: 2,
		})

		require.NoError(t, err)

		assert.Equal(t, fakedOwnership, tick.Ownership)
		assert.WithinDuration(t, time.Now(), tick.Timestamp, time.Minute)
		assert.Equal(t, fakedContestID, tick.ContestID)
		assert.Equal(t, fakedProblemID, tick.ProblemID)
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
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.ContenderRole, nil)

		ucase := usecases.TickUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		tick, err := ucase.CreateTick(context.Background(), fakedContenderID, domain.Tick{
			ProblemID:    fakedProblemID,
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

	t.Run("ProblemBelongsToDifferentContest", func(t *testing.T) {
		mockedRepo, mockedEventBroker := makeMocks(time.Now())
		mockedAuthorizer := new(authorizerMock)

		fakedOtherProblemID := fakedProblemID + 1

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.ContenderRole, nil)

		mockedRepo.
			On("GetProblem", mock.Anything, mock.Anything, fakedOtherProblemID).
			Return(domain.Problem{
				ID:        fakedOtherProblemID,
				ContestID: fakedContestID + 1,
			}, nil)

		ucase := usecases.TickUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		tick, err := ucase.CreateTick(context.Background(), fakedContenderID, domain.Tick{
			ProblemID:    fakedOtherProblemID,
			Top:          true,
			AttemptsTop:  5,
			Zone:         true,
			AttemptsZone: 2,
		})

		assert.ErrorIs(t, err, domain.ErrProblemNotInContest)
		assert.Empty(t, tick)

		mockedRepo.AssertExpectations(t)
		mockedEventBroker.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("OrganizerCanRegisterAscentAfterGracePeriod", func(t *testing.T) {
		mockedRepo, mockedEventBroker := makeMocks(time.Now().Add(-1 * gracePeriod))
		mockedAuthorizer := new(authorizerMock)

		fakedTickID := randomResourceID[domain.TickID]()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("GetProblem", mock.Anything, mock.Anything, fakedProblemID).
			Return(domain.Problem{
				ID:        fakedProblemID,
				ContestID: fakedContestID,
			}, nil)

		mockedRepo.
			On("StoreTick", mock.Anything, nil, mock.MatchedBy(func(tick domain.Tick) bool {
				tick.Timestamp = time.Time{}

				expected := domain.Tick{
					Ownership:    fakedOwnership,
					ContestID:    fakedContestID,
					ProblemID:    fakedProblemID,
					Top:          true,
					AttemptsTop:  5,
					Zone:         true,
					AttemptsZone: 2,
				}

				return tick.Timestamp.Sub(now) < time.Second && tick == expected
			})).
			Return(domain.Tick{
				ID:           fakedTickID,
				Ownership:    fakedOwnership,
				Timestamp:    now,
				ContestID:    fakedContestID,
				ProblemID:    fakedProblemID,
				Top:          true,
				AttemptsTop:  5,
				Zone:         true,
				AttemptsZone: 2,
			}, nil)

		mockedEventBroker.On("Dispatch", fakedContestID, domain.AscentRegisteredEvent{
			TickID:       fakedTickID,
			Timestamp:    now,
			ContenderID:  fakedContenderID,
			ProblemID:    fakedProblemID,
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

		tick, err := ucase.CreateTick(context.Background(), fakedContenderID, domain.Tick{
			ProblemID:    fakedProblemID,
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

		fakedContender := domain.Contender{
			ID:        fakedContenderID,
			Ownership: fakedOwnership,
		}

		mockedRepo.
			On("GetContender", mock.Anything, mock.Anything, fakedContenderID).
			Return(fakedContender, nil)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		ucase := usecases.TickUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		tick, err := ucase.CreateTick(context.Background(), fakedContenderID, domain.Tick{})

		assert.ErrorIs(t, err, domain.ErrNoOwnership)
		assert.Empty(t, tick)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})
}

func TestDeleteTick(t *testing.T) {
	fakedTickID := randomResourceID[domain.TickID]()
	fakedContenderID := randomResourceID[domain.ContenderID]()
	fakedContestID := randomResourceID[domain.ContestID]()
	fakedCompClassID := randomResourceID[domain.CompClassID]()
	fakedProblemID := randomResourceID[domain.ProblemID]()

	gracePeriod := 15 * time.Minute

	fakedOwnership := domain.OwnershipData{
		OrganizerID: randomResourceID[domain.OrganizerID](),
		ContenderID: &fakedContenderID,
	}

	makeMocks := func(timeEnd time.Time) (*repositoryMock, *eventBrokerMock) {
		mockedRepo := new(repositoryMock)
		mockedEventBroker := new(eventBrokerMock)

		fakedContender := domain.Contender{
			ID:          fakedContenderID,
			ContestID:   fakedContestID,
			CompClassID: fakedCompClassID,
		}

		fakedTick := domain.Tick{
			ID:        fakedTickID,
			Ownership: fakedOwnership,
			ProblemID: fakedProblemID,
			ContestID: fakedContestID,
		}

		mockedRepo.
			On("GetTick", mock.Anything, mock.Anything, fakedTickID).
			Return(fakedTick, nil)

		mockedRepo.
			On("GetContender", mock.Anything, mock.Anything, fakedContenderID).
			Return(fakedContender, nil)

		mockedRepo.
			On("GetContest", mock.Anything, mock.Anything, fakedContestID).
			Return(domain.Contest{
				ID:          fakedContestID,
				GracePeriod: gracePeriod,
			}, nil)

		mockedRepo.
			On("GetCompClass", mock.Anything, mock.Anything, fakedCompClassID).
			Return(domain.CompClass{
				ID:      fakedCompClassID,
				TimeEnd: timeEnd,
			}, nil)

		return mockedRepo, mockedEventBroker
	}

	t.Run("HappyPath", func(t *testing.T) {
		mockedRepo, mockedEventBroker := makeMocks(time.Now())
		mockedAuthorizer := new(authorizerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.ContenderRole, nil)

		mockedRepo.
			On("DeleteTick", mock.Anything, mock.Anything, fakedTickID).
			Return(nil)

		mockedEventBroker.On("Dispatch", fakedContestID, domain.AscentDeregisteredEvent{
			TickID:      fakedTickID,
			ContenderID: fakedContenderID,
			ProblemID:   fakedProblemID,
		}).Return()

		ucase := usecases.TickUseCase{
			Repo:        mockedRepo,
			Authorizer:  mockedAuthorizer,
			EventBroker: mockedEventBroker,
		}

		err := ucase.DeleteTick(context.Background(), fakedTickID)

		require.NoError(t, err)

		mockedRepo.AssertExpectations(t)
		mockedEventBroker.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("ContenderCannotDeregisterAscentAfterGracePeriod", func(t *testing.T) {
		mockedRepo, mockedEventBroker := makeMocks(time.Now().Add(-1 * gracePeriod))
		mockedAuthorizer := new(authorizerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.ContenderRole, nil)

		ucase := usecases.TickUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		err := ucase.DeleteTick(context.Background(), fakedTickID)

		assert.ErrorIs(t, err, domain.ErrContestEnded)

		mockedRepo.AssertExpectations(t)
		mockedEventBroker.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("OrganizerCanDeregisterAscentAfterGracePeriod", func(t *testing.T) {
		mockedRepo, mockedEventBroker := makeMocks(time.Now().Add(-1 * gracePeriod))
		mockedAuthorizer := new(authorizerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("DeleteTick", mock.Anything, mock.Anything, fakedTickID).
			Return(nil)

		mockedEventBroker.On("Dispatch", fakedContestID, domain.AscentDeregisteredEvent{
			TickID:      fakedTickID,
			ContenderID: fakedContenderID,
			ProblemID:   fakedProblemID,
		}).Return()

		ucase := usecases.TickUseCase{
			Repo:        mockedRepo,
			Authorizer:  mockedAuthorizer,
			EventBroker: mockedEventBroker,
		}

		err := ucase.DeleteTick(context.Background(), fakedTickID)

		require.NoError(t, err)

		mockedRepo.AssertExpectations(t)
		mockedEventBroker.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedRepo := new(repositoryMock)
		mockedAuthorizer := new(authorizerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		mockedRepo.
			On("GetTick", mock.Anything, mock.Anything, fakedTickID).
			Return(domain.Tick{
				ID:        fakedTickID,
				Ownership: fakedOwnership,
			}, nil)

		ucase := usecases.TickUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		err := ucase.DeleteTick(context.Background(), fakedTickID)

		assert.ErrorIs(t, err, domain.ErrNoOwnership)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})
}
