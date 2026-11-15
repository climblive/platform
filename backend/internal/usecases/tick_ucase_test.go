package usecases_test

import (
	"context"
	"testing"
	"testing/synctest"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/testutils"
	"github.com/climblive/platform/backend/internal/usecases"
	"github.com/climblive/platform/backend/internal/usecases/validators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetTicksByContender(t *testing.T) {
	fakedContenderID := testutils.RandomResourceID[domain.ContenderID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: testutils.RandomResourceID[domain.OrganizerID](),
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
				ID: testutils.RandomResourceID[domain.TickID](),
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

func TestGetTicksByCompetition(t *testing.T) {
	fakedCompetitionID := testutils.RandomResourceID[domain.CompetitionID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: testutils.RandomResourceID[domain.OrganizerID](),
	}

	makeMocks := func() (*repositoryMock, *authorizerMock) {
		mockedRepo := new(repositoryMock)
		mockedAuthorizer := new(authorizerMock)

		fakedCompetition := domain.Competition{
			ID:        fakedCompetitionID,
			Ownership: fakedOwnership,
		}

		mockedRepo.
			On("GetCompetition", mock.Anything, mock.Anything, fakedCompetitionID).
			Return(fakedCompetition, nil)

		return mockedRepo, mockedAuthorizer
	}

	t.Run("HappyPath", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

		fakedTicks := []domain.Tick{
			{
				ID: testutils.RandomResourceID[domain.TickID](),
			},
		}

		mockedRepo.
			On("GetTicksByCompetition", mock.Anything, mock.Anything, fakedCompetitionID).
			Return(fakedTicks, nil)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.ContenderRole, nil)

		ucase := usecases.TickUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		ticks, err := ucase.GetTicksByCompetition(context.Background(), fakedCompetitionID)

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

		ticks, err := ucase.GetTicksByCompetition(context.Background(), fakedCompetitionID)

		assert.ErrorIs(t, err, domain.ErrNoOwnership)
		assert.Nil(t, ticks)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})
}

func TestPutTick(t *testing.T) {
	fakedContenderID := testutils.RandomResourceID[domain.ContenderID]()
	fakedCompetitionID := testutils.RandomResourceID[domain.CompetitionID]()
	fakedCategoryID := testutils.RandomResourceID[domain.CategoryID]()
	fakedProblemID := testutils.RandomResourceID[domain.ProblemID]()

	now := time.Now()
	gracePeriod := 15 * time.Minute

	fakedOwnership := domain.OwnershipData{
		OrganizerID: testutils.RandomResourceID[domain.OrganizerID](),
		ContenderID: &fakedContenderID,
	}

	makeMocks := func(timeBegin, timeEnd time.Time) (*repositoryMock, *eventBrokerMock) {
		mockedRepo := new(repositoryMock)
		mockedEventBroker := new(eventBrokerMock)

		fakedContender := domain.Contender{
			ID:            fakedContenderID,
			Ownership:     fakedOwnership,
			CompetitionID: fakedCompetitionID,
			CategoryID:    fakedCategoryID,
		}

		mockedRepo.
			On("GetContender", mock.Anything, mock.Anything, fakedContenderID).
			Return(fakedContender, nil)

		mockedRepo.
			On("GetCompetition", mock.Anything, mock.Anything, fakedCompetitionID).
			Return(domain.Competition{
				ID:          fakedCompetitionID,
				GracePeriod: gracePeriod,
			}, nil)

		mockedRepo.
			On("GetCategory", mock.Anything, mock.Anything, fakedCategoryID).
			Return(domain.Category{
				ID:        fakedCategoryID,
				TimeBegin: timeBegin,
				TimeEnd:   timeEnd,
			}, nil)

		return mockedRepo, mockedEventBroker
	}

	t.Run("HappyPath", func(t *testing.T) {
		mockedRepo, mockedEventBroker := makeMocks(time.Now(), time.Now())
		mockedAuthorizer := new(authorizerMock)

		fakedTickID := testutils.RandomResourceID[domain.TickID]()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.ContenderRole, nil)

		mockedRepo.
			On("GetTickByContenderAndProblem", mock.Anything, nil, fakedContenderID, fakedProblemID).
			Return(domain.Tick{}, domain.ErrNotFound)

		mockedRepo.
			On("GetProblem", mock.Anything, mock.Anything, fakedProblemID).
			Return(domain.Problem{
				ID:            fakedProblemID,
				CompetitionID: fakedCompetitionID,
			}, nil)

		mockedRepo.
			On("StoreTick", mock.Anything, nil, mock.MatchedBy(func(tick domain.Tick) bool {
				tick.Timestamp = time.Time{}

				expected := domain.Tick{
					Ownership:     fakedOwnership,
					CompetitionID: fakedCompetitionID,
					Revision:      1,
					ProblemID:     fakedProblemID,
					Top:           true,
					AttemptsTop:   5,
					Zone1:         true,
					AttemptsZone1: 2,
					Zone2:         true,
					AttemptsZone2: 3,
				}

				return tick.Timestamp.Sub(now) < time.Second && tick == expected
			})).
			Return(domain.Tick{
				ID:            fakedTickID,
				Ownership:     fakedOwnership,
				Timestamp:     now,
				CompetitionID: fakedCompetitionID,
				Revision:      1,
				ProblemID:     fakedProblemID,
				Top:           true,
				AttemptsTop:   5,
				Zone1:         true,
				AttemptsZone1: 2,
				Zone2:         true,
				AttemptsZone2: 3,
			}, nil)

		mockedEventBroker.On("Dispatch", fakedCompetitionID, domain.AscentRegisteredEvent{
			TickID:        fakedTickID,
			Timestamp:     now,
			ContenderID:   fakedContenderID,
			ProblemID:     fakedProblemID,
			Revision:      1,
			Top:           true,
			AttemptsTop:   5,
			Zone1:         true,
			AttemptsZone1: 2,
			Zone2:         true,
			AttemptsZone2: 3,
		}).Return()

		ucase := usecases.TickUseCase{
			Repo:        mockedRepo,
			Authorizer:  mockedAuthorizer,
			EventBroker: mockedEventBroker,
		}

		tick, err := ucase.PutTick(context.Background(), fakedContenderID, domain.Tick{
			Revision:      1,
			ProblemID:     fakedProblemID,
			Top:           true,
			AttemptsTop:   5,
			Zone1:         true,
			AttemptsZone1: 2,
			Zone2:         true,
			AttemptsZone2: 3,
		})

		require.NoError(t, err)

		assert.Equal(t, fakedOwnership, tick.Ownership)
		assert.WithinDuration(t, time.Now(), tick.Timestamp, time.Minute)
		assert.Equal(t, fakedCompetitionID, tick.CompetitionID)
		assert.Equal(t, fakedProblemID, tick.ProblemID)
		assert.Equal(t, true, tick.Top)
		assert.Equal(t, 5, tick.AttemptsTop)
		assert.Equal(t, true, tick.Zone1)
		assert.Equal(t, 2, tick.AttemptsZone1)
		assert.Equal(t, true, tick.Zone2)
		assert.Equal(t, 3, tick.AttemptsZone2)

		mockedRepo.AssertExpectations(t)
		mockedEventBroker.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("CannotRegisterAscentBeforeCompetitionStart", func(t *testing.T) {
		mockedRepo, mockedEventBroker := makeMocks(time.Now().Add(time.Minute), time.Now().Add(time.Hour))
		mockedAuthorizer := new(authorizerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.ContenderRole, nil)

		ucase := usecases.TickUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		tick, err := ucase.PutTick(context.Background(), fakedContenderID, domain.Tick{
			ProblemID:     fakedProblemID,
			Top:           true,
			AttemptsTop:   5,
			Zone1:         true,
			AttemptsZone1: 2,
			Zone2:         true,
			AttemptsZone2: 3,
		})

		assert.ErrorIs(t, err, domain.ErrCompetitionNotStarted)
		assert.Empty(t, tick)

		mockedRepo.AssertExpectations(t)
		mockedEventBroker.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("ContenderCannotRegisterAscentAfterGracePeriod", func(t *testing.T) {
		mockedRepo, mockedEventBroker := makeMocks(time.Now().Add(-1*time.Hour), time.Now().Add(-1*gracePeriod))
		mockedAuthorizer := new(authorizerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.ContenderRole, nil)

		ucase := usecases.TickUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		tick, err := ucase.PutTick(context.Background(), fakedContenderID, domain.Tick{
			ProblemID:     fakedProblemID,
			Top:           true,
			AttemptsTop:   5,
			Zone1:         true,
			AttemptsZone1: 2,
			Zone2:         true,
			AttemptsZone2: 3,
		})

		assert.ErrorIs(t, err, domain.ErrCompetitionEnded)
		assert.Empty(t, tick)

		mockedRepo.AssertExpectations(t)
		mockedEventBroker.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("ProblemBelongsToDifferentCompetition", func(t *testing.T) {
		mockedRepo, mockedEventBroker := makeMocks(time.Now(), time.Now())
		mockedAuthorizer := new(authorizerMock)

		fakedOtherProblemID := fakedProblemID + 1

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.ContenderRole, nil)

		mockedRepo.
			On("GetProblem", mock.Anything, mock.Anything, fakedOtherProblemID).
			Return(domain.Problem{
				ID:            fakedOtherProblemID,
				CompetitionID: fakedCompetitionID + 1,
			}, nil)

		ucase := usecases.TickUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		tick, err := ucase.PutTick(context.Background(), fakedContenderID, domain.Tick{
			ProblemID:     fakedOtherProblemID,
			Top:           true,
			AttemptsTop:   5,
			Zone1:         true,
			AttemptsZone1: 2,
			Zone2:         true,
			AttemptsZone2: 3,
		})

		assert.ErrorIs(t, err, domain.ErrProblemNotInCompetition)
		assert.Empty(t, tick)

		mockedRepo.AssertExpectations(t)
		mockedEventBroker.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("OrganizerCanRegisterAscentAfterGracePeriod", func(t *testing.T) {
		mockedRepo, mockedEventBroker := makeMocks(time.Now().Add(-1*time.Hour), time.Now().Add(-1*gracePeriod))
		mockedAuthorizer := new(authorizerMock)

		fakedTickID := testutils.RandomResourceID[domain.TickID]()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("GetTickByContenderAndProblem", mock.Anything, nil, fakedContenderID, fakedProblemID).
			Return(domain.Tick{}, domain.ErrNotFound)

		mockedRepo.
			On("GetProblem", mock.Anything, mock.Anything, fakedProblemID).
			Return(domain.Problem{
				ID:            fakedProblemID,
				CompetitionID: fakedCompetitionID,
			}, nil)

		mockedRepo.
			On("StoreTick", mock.Anything, nil, mock.MatchedBy(func(tick domain.Tick) bool {
				tick.Timestamp = time.Time{}

				expected := domain.Tick{
					Ownership:     fakedOwnership,
					CompetitionID: fakedCompetitionID,
					Revision:      1,
					ProblemID:     fakedProblemID,
					Top:           true,
					AttemptsTop:   5,
					Zone1:         true,
					AttemptsZone1: 2,
					Zone2:         true,
					AttemptsZone2: 3,
				}

				return tick.Timestamp.Sub(now) < time.Second && tick == expected
			})).
			Return(domain.Tick{
				ID:            fakedTickID,
				Ownership:     fakedOwnership,
				Timestamp:     now,
				CompetitionID: fakedCompetitionID,
				Revision:      1,
				ProblemID:     fakedProblemID,
				Top:           true,
				AttemptsTop:   5,
				Zone1:         true,
				AttemptsZone1: 2,
				Zone2:         true,
				AttemptsZone2: 3,
			}, nil)

		mockedEventBroker.On("Dispatch", fakedCompetitionID, domain.AscentRegisteredEvent{
			TickID:        fakedTickID,
			Timestamp:     now,
			ContenderID:   fakedContenderID,
			ProblemID:     fakedProblemID,
			Revision:      1,
			Top:           true,
			AttemptsTop:   5,
			Zone1:         true,
			AttemptsZone1: 2,
			Zone2:         true,
			AttemptsZone2: 3,
		}).Return()

		ucase := usecases.TickUseCase{
			Repo:        mockedRepo,
			Authorizer:  mockedAuthorizer,
			EventBroker: mockedEventBroker,
		}

		tick, err := ucase.PutTick(context.Background(), fakedContenderID, domain.Tick{
			Revision:      1,
			ProblemID:     fakedProblemID,
			Top:           true,
			AttemptsTop:   5,
			Zone1:         true,
			AttemptsZone1: 2,
			Zone2:         true,
			AttemptsZone2: 3,
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

		tick, err := ucase.PutTick(context.Background(), fakedContenderID, domain.Tick{})

		assert.ErrorIs(t, err, domain.ErrNoOwnership)
		assert.Empty(t, tick)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("ValidatorIsInvoked", func(t *testing.T) {
		mockedRepo, mockedEventBroker := makeMocks(time.Now(), time.Now())
		mockedAuthorizer := new(authorizerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.ContenderRole, nil)

		mockedRepo.
			On("GetTickByContenderAndProblem", mock.Anything, nil, fakedContenderID, fakedProblemID).
			Return(domain.Tick{}, domain.ErrNotFound)

		mockedRepo.
			On("GetProblem", mock.Anything, mock.Anything, fakedProblemID).
			Return(domain.Problem{
				ID:            fakedProblemID,
				CompetitionID: fakedCompetitionID,
			}, nil)

		ucase := usecases.TickUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		_, err := ucase.PutTick(context.Background(), fakedContenderID, domain.Tick{
			Revision:  1,
			ProblemID: fakedProblemID,
			Top:       true,
			Zone2:     false,
			Zone1:     false,
		})

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.True(t, validators.TickValidator{}.IsValidationError(err))

		mockedRepo.AssertExpectations(t)
		mockedEventBroker.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("RejectsSupersededRevision", func(t *testing.T) {
		mockedRepo, mockedEventBroker := makeMocks(time.Now().Add(-time.Hour), time.Now().Add(time.Hour))
		mockedAuthorizer := new(authorizerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.ContenderRole, nil)

		mockedRepo.
			On("GetProblem", mock.Anything, nil, fakedProblemID).
			Return(domain.Problem{
				ID:            fakedProblemID,
				CompetitionID: fakedCompetitionID,
			}, nil)

		mockedRepo.
			On("GetTickByContenderAndProblem", mock.Anything, nil, fakedContenderID, fakedProblemID).
			Return(domain.Tick{Revision: 2}, nil)

		ucase := usecases.TickUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		tick, err := ucase.PutTick(context.Background(), fakedContenderID, domain.Tick{
			Revision:  2,
			ProblemID: fakedProblemID,
		})

		assert.ErrorIs(t, err, domain.ErrSuperseded)
		assert.Empty(t, tick)

		mockedRepo.AssertExpectations(t)
		mockedEventBroker.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("RejectsRevisionJumpExceedingMaximum", func(t *testing.T) {
		mockedRepo, mockedEventBroker := makeMocks(now.Add(-time.Hour), now.Add(time.Hour))
		mockedAuthorizer := new(authorizerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.ContenderRole, nil)

		mockedRepo.
			On("GetProblem", mock.Anything, nil, fakedProblemID).
			Return(domain.Problem{
				ID:            fakedProblemID,
				CompetitionID: fakedCompetitionID,
			}, nil)

		mockedRepo.
			On("GetTickByContenderAndProblem", mock.Anything, nil, fakedContenderID, fakedProblemID).
			Return(domain.Tick{Revision: 1}, nil)

		ucase := usecases.TickUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		tick, err := ucase.PutTick(context.Background(), fakedContenderID, domain.Tick{
			Revision:  1002,
			ProblemID: fakedProblemID,
		})

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.Empty(t, tick)

		mockedRepo.AssertExpectations(t)
		mockedEventBroker.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("UpdatesExistingTick", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			now := time.Now()
			mockedRepo, mockedEventBroker := makeMocks(now.Add(-time.Hour), now.Add(time.Hour))
			mockedAuthorizer := new(authorizerMock)

			fakedTickID := testutils.RandomResourceID[domain.TickID]()

			existingTick := domain.Tick{
				ID:            fakedTickID,
				Ownership:     fakedOwnership,
				Timestamp:     now.Add(-time.Minute),
				CompetitionID: fakedCompetitionID,
				ProblemID:     fakedProblemID,
				Revision:      1,
			}

			mockedRepo.
				On("GetTickByContenderAndProblem", mock.Anything, nil, fakedContenderID, fakedProblemID).
				Return(existingTick, nil)

			mockedAuthorizer.
				On("HasOwnership", mock.Anything, fakedOwnership).
				Return(domain.ContenderRole, nil)

			mockedRepo.
				On("GetProblem", mock.Anything, nil, fakedProblemID).
				Return(domain.Problem{
					ID:            fakedProblemID,
					CompetitionID: fakedCompetitionID,
				}, nil)

			mockedRepo.
				On("StoreTick", mock.Anything, nil, domain.Tick{
					ID:            fakedTickID,
					Ownership:     fakedOwnership,
					Timestamp:     now,
					CompetitionID: fakedCompetitionID,
					ProblemID:     fakedProblemID,
					Revision:      2,
					Top:           true,
					AttemptsTop:   3,
					Zone1:         true,
					AttemptsZone1: 1,
					Zone2:         true,
					AttemptsZone2: 2,
				}).
				Return(domain.Tick{
					ID:            fakedTickID,
					Ownership:     fakedOwnership,
					Timestamp:     now,
					CompetitionID: fakedCompetitionID,
					ProblemID:     fakedProblemID,
					Revision:      2,
					Top:           true,
					AttemptsTop:   3,
					Zone1:         true,
					AttemptsZone1: 1,
					Zone2:         true,
					AttemptsZone2: 2,
				}, nil)

			mockedEventBroker.
				On("Dispatch", fakedCompetitionID, domain.AscentRegisteredEvent{
					TickID:        fakedTickID,
					Timestamp:     now,
					ContenderID:   fakedContenderID,
					ProblemID:     fakedProblemID,
					Revision:      2,
					Top:           true,
					AttemptsTop:   3,
					Zone1:         true,
					AttemptsZone1: 1,
					Zone2:         true,
					AttemptsZone2: 2,
				}).
				Return()

			ucase := usecases.TickUseCase{
				Repo:        mockedRepo,
				Authorizer:  mockedAuthorizer,
				EventBroker: mockedEventBroker,
			}

			updatedTick, err := ucase.PutTick(context.Background(), fakedContenderID, domain.Tick{
				Revision:      2,
				ProblemID:     fakedProblemID,
				Top:           true,
				AttemptsTop:   3,
				Zone1:         true,
				AttemptsZone1: 1,
				Zone2:         true,
				AttemptsZone2: 2,
			})

			require.NoError(t, err)
			assert.Equal(t, fakedTickID, updatedTick.ID)
			assert.Equal(t, fakedProblemID, updatedTick.ProblemID)
			assert.Equal(t, fakedCompetitionID, updatedTick.CompetitionID)
			assert.True(t, updatedTick.Top)
			assert.Equal(t, 3, updatedTick.AttemptsTop)
			assert.True(t, updatedTick.Zone2)
			assert.Equal(t, 2, updatedTick.AttemptsZone2)
			assert.True(t, updatedTick.Zone1)
			assert.Equal(t, 1, updatedTick.AttemptsZone1)
			assert.Equal(t, now, updatedTick.Timestamp)

			mockedRepo.AssertExpectations(t)
			mockedAuthorizer.AssertExpectations(t)
			mockedEventBroker.AssertExpectations(t)
		})
	})
}

func TestDeleteTick(t *testing.T) {
	fakedTickID := testutils.RandomResourceID[domain.TickID]()
	fakedContenderID := testutils.RandomResourceID[domain.ContenderID]()
	fakedCompetitionID := testutils.RandomResourceID[domain.CompetitionID]()
	fakedCategoryID := testutils.RandomResourceID[domain.CategoryID]()
	fakedProblemID := testutils.RandomResourceID[domain.ProblemID]()

	gracePeriod := 15 * time.Minute

	fakedOwnership := domain.OwnershipData{
		OrganizerID: testutils.RandomResourceID[domain.OrganizerID](),
		ContenderID: &fakedContenderID,
	}

	makeMocks := func(timeEnd time.Time) (*repositoryMock, *eventBrokerMock) {
		mockedRepo := new(repositoryMock)
		mockedEventBroker := new(eventBrokerMock)

		fakedContender := domain.Contender{
			ID:            fakedContenderID,
			CompetitionID: fakedCompetitionID,
			CategoryID:    fakedCategoryID,
		}

		fakedTick := domain.Tick{
			ID:            fakedTickID,
			Ownership:     fakedOwnership,
			ProblemID:     fakedProblemID,
			CompetitionID: fakedCompetitionID,
		}

		mockedRepo.
			On("GetTick", mock.Anything, mock.Anything, fakedTickID).
			Return(fakedTick, nil)

		mockedRepo.
			On("GetContender", mock.Anything, mock.Anything, fakedContenderID).
			Return(fakedContender, nil)

		mockedRepo.
			On("GetCompetition", mock.Anything, mock.Anything, fakedCompetitionID).
			Return(domain.Competition{
				ID:          fakedCompetitionID,
				GracePeriod: gracePeriod,
			}, nil)

		mockedRepo.
			On("GetCategory", mock.Anything, mock.Anything, fakedCategoryID).
			Return(domain.Category{
				ID:      fakedCategoryID,
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

		mockedEventBroker.On("Dispatch", fakedCompetitionID, domain.AscentDeregisteredEvent{
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

		assert.ErrorIs(t, err, domain.ErrCompetitionEnded)

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

		mockedEventBroker.On("Dispatch", fakedCompetitionID, domain.AscentDeregisteredEvent{
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
