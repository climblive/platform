package scores_test

import (
	"context"
	"testing"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/events"
	"github.com/climblive/platform/backend/internal/scores"
	"github.com/climblive/platform/backend/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"uuid"
)

func TestScoreEngineManager(t *testing.T) {
	t.Run("StartAndStop", func(t *testing.T) {
		mockedRepo := new(repositoryMock)
		mockedStoreHydrator := new(engineStoreHydratorMock)
		mockedEventBroker := new(eventBrokerMock)

		mockedRepo.
			On("GetCompetitionsCurrentlyRunningOrByStartTime", mock.Anything, mock.Anything, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).
			Return([]domain.Competition{}, nil)

		mngr := scores.NewScoreEngineManager(mockedRepo, mockedStoreHydrator, mockedEventBroker, time.Hour)

		ctx, cancel := context.WithCancel(context.Background())

		wg := mngr.Run(ctx)

		cancel()
		wg.Wait()

		mockedRepo.AssertExpectations(t)
		mockedStoreHydrator.AssertExpectations(t)
		mockedEventBroker.AssertExpectations(t)
	})

	t.Run("LoadSingleEngine", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		mockedRepo := new(repositoryMock)
		mockedStoreHydrator := new(engineStoreHydratorMock)
		mockedEventBroker := new(eventBrokerMock)

		fakedSubscriptionID := domain.SubscriptionID(uuid.New())
		fakedCompetitionID := testutils.RandomResourceID[domain.CompetitionID]()

		now := time.Now()

		mockedRepo.
			On("GetCompetitionsCurrentlyRunningOrByStartTime", mock.Anything, mock.Anything, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).
			Return([]domain.Competition{
				{
					ID:                 fakedCompetitionID,
					QualifyingProblems: 10,
					Finalists:          7,
					TimeBegin:          now,
					TimeEnd:            now,
				},
			}, nil)

		mockedRepo.
			On("GetCompetition", mock.Anything, mock.Anything, fakedCompetitionID).
			Return(domain.Competition{
				ID:                 fakedCompetitionID,
				QualifyingProblems: 10,
				Finalists:          7,
				TimeBegin:          now,
				TimeEnd:            now,
			}, nil)

		mockedEventBroker.
			On("Subscribe", mock.Anything, mock.Anything).
			Return(fakedSubscriptionID, events.NewSubscription([]domain.EventFilter{{}}, 1000))

		mockedEventBroker.
			On("Unsubscribe", fakedSubscriptionID).
			Return()

		mockedEventBroker.
			On("Dispatch", fakedCompetitionID, mock.AnythingOfType("domain.ScoreEngineStartedEvent")).
			On("Dispatch", fakedCompetitionID, mock.AnythingOfType("domain.ScoreEngineStoppedEvent"))

		mockedStoreHydrator.
			On("Hydrate", mock.Anything, fakedCompetitionID, mock.AnythingOfType("*scores.MemoryStore")).
			Run(func(args mock.Arguments) {
				cancel()
			}).
			Return(nil)

		mngr := scores.NewScoreEngineManager(mockedRepo, mockedStoreHydrator, mockedEventBroker, time.Hour)

		wg := mngr.Run(ctx)

		<-ctx.Done()
		wg.Wait()

		mockedRepo.AssertExpectations(t)
		mockedStoreHydrator.AssertExpectations(t)
		mockedEventBroker.AssertExpectations(t)
	})

	t.Run("ManageScoreEngines", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		mockedRepo := new(repositoryMock)
		mockedStoreHydrator := new(engineStoreHydratorMock)
		mockedEventBroker := new(eventBrokerMock)

		fakedSubscriptionID := domain.SubscriptionID(uuid.New())
		fakedCompetitionID := testutils.RandomResourceID[domain.CompetitionID]()

		now := time.Now()

		mockedRepo.
			On("GetCompetitionsCurrentlyRunningOrByStartTime", mock.Anything, mock.Anything, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).
			Return([]domain.Competition{}, nil)

		mockedRepo.
			On("GetCompetition", mock.Anything, mock.Anything, fakedCompetitionID).
			Return(domain.Competition{
				ID:                 fakedCompetitionID,
				QualifyingProblems: 10,
				Finalists:          7,
				TimeBegin:          now,
				TimeEnd:            now,
			}, nil)

		mockedEventBroker.
			On("Subscribe", mock.Anything, mock.Anything).
			Return(fakedSubscriptionID, events.NewSubscription([]domain.EventFilter{{}}, 1000))

		mockedEventBroker.
			On("Unsubscribe", fakedSubscriptionID).
			Return()

		mockedStoreHydrator.
			On("Hydrate", mock.Anything, fakedCompetitionID, mock.AnythingOfType("*scores.MemoryStore")).
			Return(nil)

		mockedEventBroker.
			On("Dispatch", fakedCompetitionID, mock.AnythingOfType("domain.ScoreEngineStartedEvent")).Return().
			On("Dispatch", fakedCompetitionID, mock.AnythingOfType("domain.ScoreEngineStoppedEvent")).Return()

		mngr := scores.NewScoreEngineManager(mockedRepo, mockedStoreHydrator, mockedEventBroker, time.Hour)

		wg := mngr.Run(ctx)

		instanceID, err := mngr.StartScoreEngine(context.Background(), fakedCompetitionID, time.Now().Add(time.Hour))

		require.NoError(t, err)
		assert.NotEmpty(t, instanceID)

		_, err = mngr.StartScoreEngine(context.Background(), fakedCompetitionID, time.Now().Add(time.Hour))

		require.ErrorIs(t, err, scores.ErrAlreadyStarted)

		scoreEngineDescriptor, err := mngr.GetScoreEngine(context.Background(), instanceID)

		require.NoError(t, err)
		assert.Equal(t, instanceID, scoreEngineDescriptor.InstanceID)
		assert.Equal(t, fakedCompetitionID, scoreEngineDescriptor.CompetitionID)

		instances, err := mngr.ListScoreEnginesByCompetition(context.Background(), fakedCompetitionID)

		require.NoError(t, err)
		assert.ElementsMatch(t, []domain.ScoreEngineDescriptor{{
			InstanceID:    instanceID,
			CompetitionID: fakedCompetitionID,
		}}, instances)

		allInstances, err := mngr.ListScoreEngines(context.Background())

		require.NoError(t, err)
		assert.ElementsMatch(t, []domain.ScoreEngineDescriptor{{
			InstanceID:    instanceID,
			CompetitionID: fakedCompetitionID,
		}}, allInstances)

		err = mngr.StopScoreEngine(context.Background(), instanceID)

		require.NoError(t, err)

		scoreEngineDescriptor, err = mngr.GetScoreEngine(context.Background(), instanceID)

		require.ErrorIs(t, err, domain.ErrNotFound)
		assert.Empty(t, scoreEngineDescriptor)

		instances, err = mngr.ListScoreEnginesByCompetition(context.Background(), fakedCompetitionID)

		require.NoError(t, err)
		require.Len(t, instances, 0)

		allInstances, err = mngr.ListScoreEngines(context.Background())

		require.NoError(t, err)
		require.Len(t, allInstances, 0)

		cancel()

		wg.Wait()

		mockedRepo.AssertExpectations(t)
		mockedStoreHydrator.AssertExpectations(t)
		mockedEventBroker.AssertExpectations(t)
	})
}

type repositoryMock struct {
	mock.Mock
}

func (m *repositoryMock) GetCompetitionsCurrentlyRunningOrByStartTime(ctx context.Context, tx domain.Transaction, earliestStartTime, latestStartTime time.Time) ([]domain.Competition, error) {
	args := m.Called(ctx, tx, earliestStartTime, latestStartTime)
	return args.Get(0).([]domain.Competition), args.Error(1)
}

func (m *repositoryMock) GetCompetition(ctx context.Context, tx domain.Transaction, competitionID domain.CompetitionID) (domain.Competition, error) {
	args := m.Called(ctx, tx, competitionID)
	return args.Get(0).(domain.Competition), args.Error(1)
}

func (m *repositoryMock) GetContendersByCompetition(ctx context.Context, tx domain.Transaction, competitionID domain.CompetitionID) ([]domain.Contender, error) {
	args := m.Called(ctx, tx, competitionID)
	return args.Get(0).([]domain.Contender), args.Error(1)
}

func (m *repositoryMock) GetProblemsByCompetition(ctx context.Context, tx domain.Transaction, competitionID domain.CompetitionID) ([]domain.Problem, error) {
	args := m.Called(ctx, tx, competitionID)
	return args.Get(0).([]domain.Problem), args.Error(1)
}

func (m *repositoryMock) GetTicksByCompetition(ctx context.Context, tx domain.Transaction, competitionID domain.CompetitionID) ([]domain.Tick, error) {
	args := m.Called(ctx, tx, competitionID)
	return args.Get(0).([]domain.Tick), args.Error(1)
}

func (m *repositoryMock) StoreResult(ctx context.Context, tx domain.Transaction, result domain.Result) error {
	args := m.Called(ctx, tx, result)
	return args.Error(0)
}

type eventBrokerMock struct {
	mock.Mock
}

func (m *eventBrokerMock) Dispatch(competitionID domain.CompetitionID, event any) {
	m.Called(competitionID, event)
}

func (m *eventBrokerMock) Subscribe(filters []domain.EventFilter, bufferCapacity int) (domain.SubscriptionID, domain.EventReader) {
	args := m.Called(filters, bufferCapacity)
	return args.Get(0).(domain.SubscriptionID), args.Get(1).(domain.EventReader)
}

func (m *eventBrokerMock) Unsubscribe(subscriptionID domain.SubscriptionID) {
	m.Called(subscriptionID)
}

type engineStoreHydratorMock struct {
	mock.Mock
}

func (m *engineStoreHydratorMock) Hydrate(ctx context.Context, competitionID domain.CompetitionID, store scores.EngineStore) error {
	args := m.Called(ctx, competitionID, store)
	return args.Error(0)
}
