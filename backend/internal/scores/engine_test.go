package scores_test

import (
	"context"
	"iter"
	"slices"
	"testing"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/events"
	"github.com/climblive/platform/backend/internal/scores"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestScoreEngine(t *testing.T) {
	mockedContestID := domain.ContestID(1)

	type fixture struct {
		broker       *eventBrokerMock
		rules        *scoringRulesMock
		ranker       *rankerMock
		store        *engineStoreMock
		subscription *events.Subscription
		engine       *scores.ScoreEngine
	}

	makeFixture := func(bufferCapacity int) fixture {
		mockedEventBroker := new(eventBrokerMock)
		mockedRanker := new(rankerMock)
		mockedRules := new(scoringRulesMock)
		mockedStore := new(engineStoreMock)

		subscription := events.NewSubscription(domain.EventFilter{}, bufferCapacity)
		subscriptionID := uuid.New()

		filter := domain.NewEventFilter(
			mockedContestID,
			0,
			"CONTENDER_ENTERED",
			"CONTENDER_SWITCHED_CLASS",
			"CONTENDER_WITHDREW_FROM_FINALS",
			"CONTENDER_REENTERED_FINALS",
			"CONTENDER_DISQUALIFIED",
			"CONTENDER_REQUALIFIED",
			"ASCENT_REGISTERED",
			"ASCENT_DEREGISTERED",
			"PROBLEM_ADDED",
			"PROBLEM_UPDATED",
		)

		mockedEventBroker.On("Subscribe", filter, 0).Return(subscriptionID, subscription)

		mockedEventBroker.On("Unsubscribe", subscriptionID).Return()

		engine := scores.NewScoreEngine(mockedContestID, mockedEventBroker, mockedRules, mockedRanker, mockedStore)

		return fixture{
			broker:       mockedEventBroker,
			rules:        mockedRules,
			ranker:       mockedRanker,
			store:        mockedStore,
			subscription: subscription,
			engine:       engine,
		}
	}

	t.Run("StartAndStop", func(t *testing.T) {
		f := makeFixture(0)

		ctx, cancel := context.WithCancel(context.Background())

		wg := f.engine.Run(ctx)

		cancel()

		wg.Wait()

		f.broker.AssertExpectations(t)
		f.rules.AssertExpectations(t)
		f.ranker.AssertExpectations(t)
		f.store.AssertExpectations(t)
	})

	t.Run("SubscriptionUnexpectedlyClosed", func(t *testing.T) {
		f := makeFixture(1)

		err := f.subscription.Post(domain.EventEnvelope{
			Name: "CONTENDER_SCORE_UPDATED",
			Data: domain.ContenderScoreUpdatedEvent{},
		})
		require.NoError(t, err)

		err = f.subscription.Post(domain.EventEnvelope{
			Name: "CONTENDER_ENTERED",
			Data: domain.ContenderEnteredEvent{},
		})
		require.ErrorIs(t, err, events.ErrBufferFull)

		wg := f.engine.Run(context.Background())

		wg.Wait()

		f.broker.AssertExpectations(t)
		f.rules.AssertExpectations(t)
		f.ranker.AssertExpectations(t)
		f.store.AssertExpectations(t)
	})

	t.Run("ContenderEntered", func(t *testing.T) {
		f := makeFixture(0)

		f.store.On("SaveContender", scores.Contender{
			ID:          1,
			CompClassID: 1,
		}).Return()

		f.store.
			On("GetContendersByCompClass", domain.CompClassID(1)).
			Return(slices.Values([]scores.Contender{{ID: 1}, {ID: 2}, {ID: 3}}))

		f.ranker.
			On("RankContenders", mock.MatchedBy(func(contenders iter.Seq[scores.Contender]) bool {
				return assert.ObjectsAreEqual(slices.Collect(contenders), []scores.Contender{{ID: 1}, {ID: 2}, {ID: 3}})
			})).
			Return([]domain.Score{{ContenderID: 1, Score: 100}, {ContenderID: 2, Score: 200}, {ContenderID: 3, Score: 300}})

		f.store.On("SaveScore", domain.Score{ContenderID: 1, Score: 100}).Return()
		f.store.On("SaveScore", domain.Score{ContenderID: 2, Score: 200}).Return()
		f.store.On("SaveScore", domain.Score{ContenderID: 3, Score: 300}).Return()

		err := f.subscription.Post(domain.EventEnvelope{
			Name: "CONTENDER_ENTERED",
			Data: domain.ContenderEnteredEvent{
				ContenderID: 1,
				CompClassID: 1,
			},
		})
		require.NoError(t, err)

		wg := f.engine.Run(context.Background())

		f.subscription.Terminate()

		wg.Wait()

		f.broker.AssertExpectations(t)
		f.rules.AssertExpectations(t)
		f.ranker.AssertExpectations(t)
		f.store.AssertExpectations(t)
	})
}

type rankerMock struct {
	mock.Mock
}

func (m *rankerMock) RankContenders(contenders iter.Seq[scores.Contender]) []domain.Score {
	args := m.Called(contenders)
	return args.Get(0).([]domain.Score)
}

type scoringRulesMock struct {
	mock.Mock
}

func (m *scoringRulesMock) CalculateScore(tickPointValues iter.Seq[int]) int {
	args := m.Called(tickPointValues)
	return args.Get(0).(int)
}

type engineStoreMock struct {
	mock.Mock
}

func (m *engineStoreMock) GetContender(contenderID domain.ContenderID) (scores.Contender, bool) {
	args := m.Called(contenderID)
	return args.Get(0).(scores.Contender), args.Bool(1)
}

func (m *engineStoreMock) GetContendersByCompClass(compClassID domain.CompClassID) iter.Seq[scores.Contender] {
	args := m.Called(compClassID)
	return args.Get(0).(iter.Seq[scores.Contender])
}

func (m *engineStoreMock) GetAllContenders() iter.Seq[scores.Contender] {
	args := m.Called()
	return args.Get(0).(iter.Seq[scores.Contender])
}

func (m *engineStoreMock) SaveContender(score scores.Contender) {
	m.Called(score)
}

func (m *engineStoreMock) GetCompClassIDs() []domain.CompClassID {
	args := m.Called()
	return args.Get(0).([]domain.CompClassID)
}

func (m *engineStoreMock) GetTicks(contenderID domain.ContenderID) iter.Seq[scores.Tick] {
	args := m.Called(contenderID)
	return args.Get(0).(iter.Seq[scores.Tick])
}

func (m *engineStoreMock) SaveTick(contenderID domain.ContenderID, tick scores.Tick) {
	m.Called(contenderID, tick)
}

func (m *engineStoreMock) DeleteTick(contenderID domain.ContenderID, problemID domain.ProblemID) {
	m.Called(contenderID, problemID)
}

func (m *engineStoreMock) GetProblem(problemID domain.ProblemID) (scores.Problem, bool) {
	args := m.Called(problemID)
	return args.Get(0).(scores.Problem), args.Bool(1)
}

func (m *engineStoreMock) SaveProblem(problem scores.Problem) {
	m.Called(problem)
}

func (m *engineStoreMock) SaveScore(score domain.Score) {
	m.Called(score)
}

func (m *engineStoreMock) GetUnpublishedScores() []domain.Score {
	args := m.Called()
	return args.Get(0).([]domain.Score)
}
