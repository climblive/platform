package scores_test

import (
	"context"
	"iter"
	"testing"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/events"
	"github.com/climblive/platform/backend/internal/scores"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestScoreEngine(t *testing.T) {
	mockedContestID := domain.ContestID(1)

	makeMocks := func(bufferCapacity int) (*eventBrokerMock, *rankerMock, *scoringRulesMock, *events.Subscription) {
		mockedEventBroker := new(eventBrokerMock)
		mockedRanker := new(rankerMock)
		mockedRules := new(scoringRulesMock)

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

		return mockedEventBroker, mockedRanker, mockedRules, subscription
	}

	t.Run("StartAndStop", func(t *testing.T) {
		mockedEventBroker, mockedRanker, mockedRules, _ := makeMocks(0)
		engine := scores.NewScoreEngine(mockedContestID, mockedEventBroker, mockedRules, mockedRanker)

		ctx, cancel := context.WithCancel(context.Background())

		wg := engine.Run(ctx)

		cancel()

		wg.Wait()

		mockedEventBroker.AssertExpectations(t)
		mockedRanker.AssertExpectations(t)
		mockedRules.AssertExpectations(t)
	})

	t.Run("SubscriptionUnexpectedlyClosed", func(t *testing.T) {
		mockedEventBroker, mockedRanker, mockedRules, subscription := makeMocks(1)
		engine := scores.NewScoreEngine(mockedContestID, mockedEventBroker, mockedRules, mockedRanker)

		err := subscription.Post(domain.EventEnvelope{
			Name: "CONTENDER_SCORE_UPDATED",
			Data: domain.ContenderScoreUpdatedEvent{},
		})
		require.NoError(t, err)

		err = subscription.Post(domain.EventEnvelope{
			Name: "CONTENDER_ENTERED",
			Data: domain.ContenderEnteredEvent{},
		})
		require.ErrorIs(t, err, events.ErrBufferFull)

		wg := engine.Run(context.Background())

		wg.Wait()

		mockedEventBroker.AssertExpectations(t)
		mockedRanker.AssertExpectations(t)
		mockedRules.AssertExpectations(t)
	})
}

type rankerMock struct {
	mock.Mock
}

func (m *rankerMock) RankContenders(contenders iter.Seq[*scores.Contender]) []domain.Score {
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
