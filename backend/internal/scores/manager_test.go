package scores_test

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/events"
	"github.com/climblive/platform/backend/internal/scores"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

func TestScoreEngineManager(t *testing.T) {
	t.Run("StartAndStop", func(t *testing.T) {
		mockedRepo := new(repositoryMock)
		mockedStoreHydrator := new(hydratorMock)
		mockedEventBroker := new(eventBrokerMock)

		mockedRepo.
			On("GetContestsCurrentlyRunningOrByStartTime", mock.Anything, mock.Anything, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).
			Return([]domain.Contest{}, nil)

		mngr := scores.NewScoreEngineManager(mockedRepo, mockedStoreHydrator, mockedEventBroker)

		ctx, cancel := context.WithCancel(context.Background())

		wg := mngr.Run(ctx)

		cancel()
		wg.Wait()
	})

	t.Run("LoadSingleEngine", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		mockedRepo := new(repositoryMock)
		mockedStoreHydrator := new(hydratorMock)
		mockedEventBroker := new(eventBrokerMock)

		mockedSubscriptionID := domain.SubscriptionID(uuid.New())
		mockedContestID := domain.ContestID(rand.Int())

		now := time.Now()

		mockedRepo.
			On("GetContestsCurrentlyRunningOrByStartTime", mock.Anything, mock.Anything, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).
			Return([]domain.Contest{
				{
					ID:                 mockedContestID,
					QualifyingProblems: 10,
					Finalists:          7,
					TimeBegin:          &now,
					TimeEnd:            &now,
				},
			}, nil)

		mockedEventBroker.
			On("Subscribe", mock.Anything, mock.Anything).
			Return(mockedSubscriptionID, events.NewSubscription(domain.EventFilter{}, 1000))

		mockedEventBroker.
			On("Unsubscribe", mockedSubscriptionID).
			Return()

		mockedStoreHydrator.
			On("Hydrate", mock.Anything, mockedContestID, mock.AnythingOfType("*scores.MemoryStore")).
			Run(func(args mock.Arguments) {
				cancel()
			}).
			Return(nil)

		mngr := scores.NewScoreEngineManager(mockedRepo, mockedStoreHydrator, mockedEventBroker)

		wg := mngr.Run(ctx)

		<-ctx.Done()
		wg.Wait()

		mockedRepo.AssertExpectations(t)
		mockedStoreHydrator.AssertExpectations(t)
		mockedEventBroker.AssertExpectations(t)
	})
}

type repositoryMock struct {
	mock.Mock
}

func (m *repositoryMock) GetContestsCurrentlyRunningOrByStartTime(ctx context.Context, tx domain.Transaction, earliestStartTime, latestStartTime time.Time) ([]domain.Contest, error) {
	args := m.Called(ctx, tx, earliestStartTime, latestStartTime)
	return args.Get(0).([]domain.Contest), args.Error(1)
}

func (m *repositoryMock) GetContendersByContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) ([]domain.Contender, error) {
	args := m.Called(ctx, tx, contestID)
	return args.Get(0).([]domain.Contender), args.Error(1)
}

func (m *repositoryMock) GetProblemsByContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) ([]domain.Problem, error) {
	args := m.Called(ctx, tx, contestID)
	return args.Get(0).([]domain.Problem), args.Error(1)
}

func (m *repositoryMock) GetTicksByContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) ([]domain.Tick, error) {
	args := m.Called(ctx, tx, contestID)
	return args.Get(0).([]domain.Tick), args.Error(1)
}

func (m *repositoryMock) StoreScore(ctx context.Context, tx domain.Transaction, score domain.Score) (domain.Score, error) {
	args := m.Called(ctx, tx, score)
	return args.Get(0).(domain.Score), args.Error(1)
}

type eventBrokerMock struct {
	mock.Mock
}

func (m *eventBrokerMock) Dispatch(contestID domain.ContestID, event any) {
	m.Called(contestID, event)
}

func (m *eventBrokerMock) Subscribe(filter domain.EventFilter, bufferCapacity int) (domain.SubscriptionID, domain.EventReader) {
	args := m.Called(filter, bufferCapacity)
	return args.Get(0).(domain.SubscriptionID), args.Get(1).(domain.EventReader)
}

func (m *eventBrokerMock) Unsubscribe(subscriptionID domain.SubscriptionID) {
	m.Called(subscriptionID)
}

type hydratorMock struct {
	mock.Mock
}

func (m *hydratorMock) Hydrate(ctx context.Context, contestID domain.ContestID, store scores.EngineStore) error {
	args := m.Called(ctx, contestID, store)
	return args.Error(0)
}
