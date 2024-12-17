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
)

func TestScoreEngine(t *testing.T) {
	mockedContestID := domain.ContestID(1)

	makeMocks := func(bufferCapacity int) (*eventBrokerMock, *rankerMock, *scoringRulesMock) {
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

		return mockedEventBroker, mockedRanker, mockedRules
	}

	t.Run("StartAndStop", func(t *testing.T) {
		mockedEventBroker, mockedRanker, mockedRules := makeMocks(0)
		engine := scores.NewScoreEngine(mockedContestID, mockedEventBroker, mockedRules, mockedRanker)

		ctx, cancel := context.WithCancel(context.Background())

		wg := engine.Run(ctx)

		cancel()

		wg.Wait()
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
