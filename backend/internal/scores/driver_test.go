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
	"github.com/stretchr/testify/require"
)

func TestEngineDriver(t *testing.T) {
	fakedContestID := domain.ContestID(rand.Int())
	fakedInstanceID := uuid.New()

	type fixture struct {
		broker       *eventBrokerMock
		subscription *events.Subscription
		driver       *scores.ScoreEngineDriver
	}

	makeFixture := func(bufferCapacity int) (fixture, func(t *testing.T)) {
		mockedEventBroker := new(eventBrokerMock)

		subscription := events.NewSubscription(domain.EventFilter{}, bufferCapacity)
		subscriptionID := uuid.New()

		filter := domain.NewEventFilter(
			fakedContestID,
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
		)

		mockedEventBroker.On("Subscribe", filter, 0).Return(subscriptionID, subscription)

		mockedEventBroker.On("Unsubscribe", subscriptionID).Return()

		mockedEventBroker.On("Dispatch", fakedContestID, domain.ScoreEngineStartedEvent{
			InstanceID: fakedInstanceID,
		}).Return()

		mockedEventBroker.On("Dispatch", fakedContestID, domain.ScoreEngineStoppedEvent{
			InstanceID: fakedInstanceID,
		}).Return()

		driver := scores.NewScoreEngineDriver(fakedContestID, fakedInstanceID, mockedEventBroker)

		awaitExpectations := func(t *testing.T) {
			mockedEventBroker.AssertExpectations(t)
		}

		return fixture{
			broker:       mockedEventBroker,
			subscription: subscription,
			driver:       driver,
		}, awaitExpectations
	}

	t.Run("StartAndStop", func(t *testing.T) {
		f, awaitExpectations := makeFixture(0)

		ctx, cancel := context.WithCancel(context.Background())

		wg, _ := f.driver.Run(ctx)

		cancel()

		wg.Wait()

		awaitExpectations(t)
	})

	t.Run("SubscriptionUnexpectedlyClosed", func(t *testing.T) {
		f, awaitExpectations := makeFixture(1)

		err := f.subscription.Post(domain.EventEnvelope{
			Data: domain.ContenderScoreUpdatedEvent{},
		})
		require.NoError(t, err)

		err = f.subscription.Post(domain.EventEnvelope{
			Data: domain.ContenderEnteredEvent{},
		})
		require.ErrorIs(t, err, events.ErrBufferFull)

		wg, _ := f.driver.Run(context.Background())

		wg.Wait()

		awaitExpectations(t)
	})

	t.Run("InstallEngine", func(t *testing.T) {
		f, awaitExpectations := makeFixture(0)

		ctx, cancel := context.WithCancel(context.Background())
		wg, installEngine := f.driver.Run(ctx)

		mockedEngine := new(scoreEngineMock)

		mockedEngine.On("Start").Return()
		mockedEngine.On("Stop").Return()
		mockedEngine.On("GetDirtyScores").Return([]domain.Score{})

		installEngine(mockedEngine)

		cancel()

		wg.Wait()

		awaitExpectations(t)
		mockedEngine.AssertExpectations(t)
	})

	t.Run("ReplayPendingEvents", func(t *testing.T) {
		f, awaitExpectations := makeFixture(0)

		ctx, cancel := context.WithCancel(context.Background())
		wg, installEngine := f.driver.Run(ctx)

		err := f.subscription.Post(domain.EventEnvelope{
			Data: domain.AscentRegisteredEvent{
				ContenderID: 1,
				ProblemID:   1,
				Top:         true,
				AttemptsTop: 10,
			},
		})
		require.NoError(t, err)

		err = f.subscription.Post(domain.EventEnvelope{
			Data: domain.AscentDeregisteredEvent{
				ContenderID: 1,
				ProblemID:   2,
			},
		})
		require.NoError(t, err)

		err = f.subscription.Post(domain.EventEnvelope{
			Data: domain.ContenderEnteredEvent{
				ContenderID: 2,
				CompClassID: 1,
			},
		})
		require.NoError(t, err)

		mockedEngine := new(scoreEngineMock)

		mockedEngine.On("Start").Run(func(args mock.Arguments) {
			cancel()
		}).Return()
		mockedEngine.On("Stop").Return()
		mockedEngine.On("GetDirtyScores").Return([]domain.Score{})

		mockedEngine.On("HandleAscentRegistered", domain.AscentRegisteredEvent{
			ContenderID: 1,
			ProblemID:   1,
			Top:         true,
			AttemptsTop: 10,
		}).Return()

		mockedEngine.On("HandleAscentDeregistered", domain.AscentDeregisteredEvent{
			ContenderID: 1,
			ProblemID:   2,
		}).Return()

		mockedEngine.On("HandleContenderEntered", domain.ContenderEnteredEvent{
			ContenderID: 2,
			CompClassID: 1,
		}).Return()

		time.Sleep(time.Millisecond)

		mockedEngine.AssertNotCalled(t, "HandleAscentRegistered")
		mockedEngine.AssertNotCalled(t, "HandleAscentDeregistered")
		mockedEngine.AssertNotCalled(t, "HandleContenderEntered")

		installEngine(mockedEngine)

		wg.Wait()

		awaitExpectations(t)
		mockedEngine.AssertExpectations(t)
	})

	t.Run("HandleEvents", func(t *testing.T) {
		f, awaitExpectations := makeFixture(0)

		ctx, cancel := context.WithCancel(context.Background())
		wg, installEngine := f.driver.Run(ctx)

		events := []any{}

		events = append(events, domain.ContenderEnteredEvent{
			ContenderID: 1,
			CompClassID: 1,
		})
		events = append(events, domain.ContenderSwitchedClassEvent{
			ContenderID: 1,
			CompClassID: 1,
		})
		events = append(events, domain.ContenderWithdrewFromFinalsEvent{
			ContenderID: 1,
		})
		events = append(events, domain.ContenderReenteredFinalsEvent{
			ContenderID: 1,
		})
		events = append(events, domain.ContenderDisqualifiedEvent{
			ContenderID: 1,
		})
		events = append(events, domain.ContenderRequalifiedEvent{
			ContenderID: 1,
		})
		events = append(events, domain.AscentRegisteredEvent{
			ContenderID:  1,
			ProblemID:    1,
			Top:          true,
			AttemptsTop:  999,
			Zone:         true,
			AttemptsZone: 42,
		})
		events = append(events, domain.AscentDeregisteredEvent{
			ContenderID: 1,
			ProblemID:   1,
		})
		events = append(events, domain.ProblemAddedEvent{
			ProblemID:  1,
			PointsTop:  1000,
			PointsZone: 500,
			FlashBonus: 100,
		})

		for _, event := range events {
			err := f.subscription.Post(domain.EventEnvelope{
				Data: event,
			})

			require.NoError(t, err)
		}

		mockedEngine := new(scoreEngineMock)

		mockedEngine.On("Start").Return()
		mockedEngine.On("Stop").Return()
		mockedEngine.On("GetDirtyScores").Return([]domain.Score{})

		mockedEngine.On("HandleContenderEntered", domain.ContenderEnteredEvent{
			ContenderID: 1,
			CompClassID: 1,
		}).Return()
		mockedEngine.On("HandleContenderSwitchedClass", domain.ContenderSwitchedClassEvent{
			ContenderID: 1,
			CompClassID: 1,
		}).Return()
		mockedEngine.On("HandleContenderWithdrewFromFinals", domain.ContenderWithdrewFromFinalsEvent{
			ContenderID: 1,
		}).Return()
		mockedEngine.On("HandleContenderReenteredFinals", domain.ContenderReenteredFinalsEvent{
			ContenderID: 1,
		}).Return()
		mockedEngine.On("HandleContenderDisqualified", domain.ContenderDisqualifiedEvent{
			ContenderID: 1,
		}).Return()
		mockedEngine.On("HandleContenderRequalified", domain.ContenderRequalifiedEvent{
			ContenderID: 1,
		}).Return()
		mockedEngine.On("HandleAscentRegistered", domain.AscentRegisteredEvent{
			ContenderID:  1,
			ProblemID:    1,
			Top:          true,
			AttemptsTop:  999,
			Zone:         true,
			AttemptsZone: 42,
		}).Return()
		mockedEngine.On("HandleAscentDeregistered", domain.AscentDeregisteredEvent{
			ContenderID: 1,
			ProblemID:   1,
		}).Return()
		mockedEngine.On("HandleProblemAdded", domain.ProblemAddedEvent{
			ProblemID:  1,
			PointsTop:  1000,
			PointsZone: 500,
			FlashBonus: 100,
		}).Run(func(mock.Arguments) { cancel() }).Return()

		installEngine(mockedEngine)

		wg.Wait()

		awaitExpectations(t)
		mockedEngine.AssertExpectations(t)
	})

	t.Run("SetScoringRules", func(t *testing.T) {
		f, awaitExpectations := makeFixture(0)

		ctx, cancel := context.WithCancel(context.Background())
		wg, installEngine := f.driver.Run(ctx)

		mockedEngine := new(scoreEngineMock)

		engineIgnition := make(chan struct{})

		mockedEngine.On("Start").Run(func(mock.Arguments) { close(engineIgnition) }).Return()
		mockedEngine.On("Stop").Return()
		mockedEngine.On("GetDirtyScores").Return([]domain.Score{})

		newRules := &jackpotRules{}

		mockedEngine.On("ReplaceScoringRules", newRules).Run(func(mock.Arguments) { cancel() }).Return()

		installEngine(mockedEngine)

		<-engineIgnition

		f.driver.SetScoringRules(newRules)

		wg.Wait()

		awaitExpectations(t)
		mockedEngine.AssertExpectations(t)
	})

	t.Run("SetRanker", func(t *testing.T) {
		f, awaitExpectations := makeFixture(0)

		ctx, cancel := context.WithCancel(context.Background())
		wg, installEngine := f.driver.Run(ctx)

		mockedEngine := new(scoreEngineMock)

		engineIgnition := make(chan struct{})

		mockedEngine.On("Start").Run(func(mock.Arguments) { close(engineIgnition) }).Return()
		mockedEngine.On("Stop").Return()
		mockedEngine.On("GetDirtyScores").Return([]domain.Score{})

		newRanker := &fakeRanker{}

		mockedEngine.On("ReplaceRanker", newRanker).Run(func(mock.Arguments) { cancel() }).Return()

		installEngine(mockedEngine)

		<-engineIgnition

		f.driver.SetRanker(newRanker)

		wg.Wait()

		awaitExpectations(t)
		mockedEngine.AssertExpectations(t)
	})

	t.Run("PublishScores", func(t *testing.T) {
		f, awaitExpectations := makeFixture(0)

		ctx, cancel := context.WithCancel(context.Background())
		wg, installEngine := f.driver.Run(ctx)

		now := time.Now()

		mockedEngine := new(scoreEngineMock)

		mockedEngine.On("Start").Return()
		mockedEngine.On("Stop").Return()
		mockedEngine.On("GetDirtyScores").Return([]domain.Score{
			{
				ContenderID: 1,
				Timestamp:   now,
				Score:       100,
				Placement:   1,
				RankOrder:   0,
				Finalist:    true,
			},
			{
				ContenderID: 2,
				Timestamp:   now,
				Score:       200,
				Placement:   2,
				RankOrder:   1,
				Finalist:    true,
			},
			{
				ContenderID: 3,
				Timestamp:   now,
				Score:       300,
				Placement:   3,
				RankOrder:   2,
				Finalist:    false,
			},
		})

		f.broker.
			On("Dispatch", fakedContestID,
				domain.ContenderScoreUpdatedEvent{
					ContenderID: 1,
					Timestamp:   now,
					Score:       100,
					Placement:   1,
					RankOrder:   0,
					Finalist:    true,
				},
			).Return().
			On("Dispatch", fakedContestID,
				domain.ContenderScoreUpdatedEvent{
					ContenderID: 2,
					Timestamp:   now,
					Score:       200,
					Placement:   2,
					RankOrder:   1,
					Finalist:    true,
				},
			).Return().
			On("Dispatch", fakedContestID,
				domain.ContenderScoreUpdatedEvent{
					ContenderID: 3,
					Timestamp:   now,
					Score:       300,
					Placement:   3,
					RankOrder:   2,
					Finalist:    false,
				},
			).Return().
			On("Dispatch", fakedContestID,
				[]domain.ContenderScoreUpdatedEvent{
					{
						ContenderID: 1,
						Timestamp:   now,
						Score:       100,
						Placement:   1,
						RankOrder:   0,
						Finalist:    true,
					},
					{
						ContenderID: 2,
						Timestamp:   now,
						Score:       200,
						Placement:   2,
						RankOrder:   1,
						Finalist:    true,
					},
					{
						ContenderID: 3,
						Timestamp:   now,
						Score:       300,
						Placement:   3,
						RankOrder:   2,
						Finalist:    false,
					},
				},
			).Return()

		installEngine(mockedEngine)

		cancel()

		wg.Wait()

		awaitExpectations(t)
		mockedEngine.AssertExpectations(t)
	})
}

type scoreEngineMock struct {
	mock.Mock
}

func (m *scoreEngineMock) Start() {
	m.Called()
}

func (m *scoreEngineMock) Stop() {
	m.Called()
}

func (m *scoreEngineMock) ReplaceScoringRules(rules scores.ScoringRules) {
	m.Called(rules)
}

func (m *scoreEngineMock) ReplaceRanker(ranker scores.Ranker) {
	m.Called(ranker)
}

func (m *scoreEngineMock) HandleContenderEntered(event domain.ContenderEnteredEvent) {
	m.Called(event)
}

func (m *scoreEngineMock) HandleContenderSwitchedClass(event domain.ContenderSwitchedClassEvent) {
	m.Called(event)
}

func (m *scoreEngineMock) HandleContenderWithdrewFromFinals(event domain.ContenderWithdrewFromFinalsEvent) {
	m.Called(event)
}

func (m *scoreEngineMock) HandleContenderReenteredFinals(event domain.ContenderReenteredFinalsEvent) {
	m.Called(event)
}

func (m *scoreEngineMock) HandleContenderDisqualified(event domain.ContenderDisqualifiedEvent) {
	m.Called(event)
}

func (m *scoreEngineMock) HandleContenderRequalified(event domain.ContenderRequalifiedEvent) {
	m.Called(event)
}

func (m *scoreEngineMock) HandleAscentRegistered(event domain.AscentRegisteredEvent) {
	m.Called(event)
}

func (m *scoreEngineMock) HandleAscentDeregistered(event domain.AscentDeregisteredEvent) {
	m.Called(event)
}

func (m *scoreEngineMock) HandleProblemAdded(event domain.ProblemAddedEvent) {
	m.Called(event)
}

func (m *scoreEngineMock) GetDirtyScores() []domain.Score {
	args := m.Called()
	return args.Get(0).([]domain.Score)
}
