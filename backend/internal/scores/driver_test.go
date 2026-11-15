package scores_test

import (
	"context"
	"iter"
	"testing"
	"testing/synctest"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/events"
	"github.com/climblive/platform/backend/internal/scores"
	"github.com/climblive/platform/backend/internal/testutils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"uuid"
)

var noEffects iter.Seq[scores.Effect]

func TestEngineDriver(t *testing.T) {
	fakedCompetitionID := testutils.RandomResourceID[domain.CompetitionID]()
	fakedInstanceID := uuid.New()

	type fixture struct {
		broker       *eventBrokerMock
		subscription *events.Subscription
		driver       *scores.ScoreEngineDriver
	}

	makeFixture := func(bufferCapacity int) (fixture, func(t *testing.T)) {
		mockedEventBroker := new(eventBrokerMock)

		filter := domain.NewEventFilter(
			fakedCompetitionID,
			0,
			"COMPETITOR_ENTERED",
			"COMPETITOR_SWITCHED_CATEGORY",
			"COMPETITOR_WITHDREW_FROM_FINALS",
			"COMPETITOR_REENTERED_FINALS",
			"COMPETITOR_DISQUALIFIED",
			"COMPETITOR_REQUALIFIED",
			"ASCENT_REGISTERED",
			"ASCENT_DEREGISTERED",
			"PROBLEM_ADDED",
			"PROBLEM_UPDATED",
			"RULES_UPDATED",
		)

		filters := []domain.EventFilter{filter}

		subscription := events.NewSubscription(filters, bufferCapacity)
		subscriptionID := uuid.New()

		mockedEventBroker.On("Subscribe", filters, 0).Return(subscriptionID, subscription)

		mockedEventBroker.On("Unsubscribe", subscriptionID).Return()

		mockedEventBroker.On("Dispatch", fakedCompetitionID, domain.ScoreEngineStartedEvent{
			InstanceID: fakedInstanceID,
		}).Return()

		mockedEventBroker.On("Dispatch", fakedCompetitionID, domain.ScoreEngineStoppedEvent{
			InstanceID: fakedInstanceID,
		}).Return()

		driver := scores.NewScoreEngineDriver(fakedCompetitionID, fakedInstanceID, mockedEventBroker)

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
			Data: domain.CompetitorResultUpdatedEvent{},
		})
		require.NoError(t, err)

		err = f.subscription.Post(domain.EventEnvelope{
			Data: domain.CompetitorEnteredEvent{},
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

		mockedEngine.On("Start").Run(func(args mock.Arguments) { cancel() }).Return(noEffects)
		mockedEngine.On("Stop").Return()
		mockedEngine.On("GetDirtyResults").Return([]domain.Result{})
		mockedEngine.On("GetDirtyPointValues").Return([]domain.PointValue{})

		installEngine(mockedEngine)

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
				CompetitorID: 1,
				ProblemID:    1,
				Top:          true,
				AttemptsTop:  10,
			},
		})
		require.NoError(t, err)

		err = f.subscription.Post(domain.EventEnvelope{
			Data: domain.AscentDeregisteredEvent{
				CompetitorID: 1,
				ProblemID:    2,
			},
		})
		require.NoError(t, err)

		err = f.subscription.Post(domain.EventEnvelope{
			Data: domain.CompetitorEnteredEvent{
				CompetitorID: 2,
				CategoryID:   1,
			},
		})
		require.NoError(t, err)

		mockedEngine := new(scoreEngineMock)

		mockedEngine.On("Start").Run(func(args mock.Arguments) {
			cancel()
		}).Return(noEffects)
		mockedEngine.On("Stop").Return()
		mockedEngine.On("GetDirtyResults").Return([]domain.Result{})
		mockedEngine.On("GetDirtyPointValues").Return([]domain.PointValue{})

		mockedEngine.On("HandleAscentRegistered", domain.AscentRegisteredEvent{
			CompetitorID: 1,
			ProblemID:    1,
			Top:          true,
			AttemptsTop:  10,
		}).Return(noEffects)

		mockedEngine.On("HandleAscentDeregistered", domain.AscentDeregisteredEvent{
			CompetitorID: 1,
			ProblemID:    2,
		}).Return(noEffects)

		mockedEngine.On("HandleCompetitorEntered", domain.CompetitorEnteredEvent{
			CompetitorID: 2,
			CategoryID:   1,
		}).Return(noEffects)

		time.Sleep(time.Millisecond)

		mockedEngine.AssertNotCalled(t, "HandleAscentRegistered")
		mockedEngine.AssertNotCalled(t, "HandleAscentDeregistered")
		mockedEngine.AssertNotCalled(t, "HandleCompetitorEntered")

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

		events = append(events, domain.RulesUpdatedEvent{
			CompetitionID:      1,
			QualifyingProblems: 10,
			Finalists:          7,
		})
		events = append(events, domain.CompetitorEnteredEvent{
			CompetitorID: 1,
			CategoryID:   1,
		})
		events = append(events, domain.CompetitorSwitchedCategoryEvent{
			CompetitorID: 1,
			CategoryID:   1,
		})
		events = append(events, domain.CompetitorWithdrewFromFinalsEvent{
			CompetitorID: 1,
		})
		events = append(events, domain.CompetitorReenteredFinalsEvent{
			CompetitorID: 1,
		})
		events = append(events, domain.CompetitorDisqualifiedEvent{
			CompetitorID: 1,
		})
		events = append(events, domain.CompetitorRequalifiedEvent{
			CompetitorID: 1,
		})
		events = append(events, domain.AscentRegisteredEvent{
			CompetitorID:  1,
			ProblemID:     1,
			Top:           true,
			AttemptsTop:   999,
			Zone1:         true,
			AttemptsZone1: 7,
			Zone2:         true,
			AttemptsZone2: 42,
		})
		events = append(events, domain.AscentDeregisteredEvent{
			CompetitorID: 1,
			ProblemID:    1,
		})
		events = append(events, domain.ProblemAddedEvent{
			ProblemID: 1,
			ProblemValue: domain.ProblemValue{
				PointsTop:   1000,
				PointsZone1: 500,
				PointsZone2: 750,
				FlashBonus:  100,
			},
		})

		for _, event := range events {
			err := f.subscription.Post(domain.EventEnvelope{
				Data: event,
			})

			require.NoError(t, err)
		}

		mockedEngine := new(scoreEngineMock)

		mockedEngine.On("Start").Return(noEffects)
		mockedEngine.On("Stop").Return()
		mockedEngine.On("GetDirtyResults").Return([]domain.Result{})
		mockedEngine.On("GetDirtyPointValues").Return([]domain.PointValue{})

		mockedEngine.On("HandleRulesUpdated", domain.RulesUpdatedEvent{
			CompetitionID:      1,
			QualifyingProblems: 10,
			Finalists:          7,
		}).Return(noEffects)
		mockedEngine.On("HandleCompetitorEntered", domain.CompetitorEnteredEvent{
			CompetitorID: 1,
			CategoryID:   1,
		}).Return(noEffects)
		mockedEngine.On("HandleCompetitorSwitchedCategory", domain.CompetitorSwitchedCategoryEvent{
			CompetitorID: 1,
			CategoryID:   1,
		}).Return(noEffects)
		mockedEngine.On("HandleCompetitorWithdrewFromFinals", domain.CompetitorWithdrewFromFinalsEvent{
			CompetitorID: 1,
		}).Return(noEffects)
		mockedEngine.On("HandleCompetitorReenteredFinals", domain.CompetitorReenteredFinalsEvent{
			CompetitorID: 1,
		}).Return(noEffects)
		mockedEngine.On("HandleCompetitorDisqualified", domain.CompetitorDisqualifiedEvent{
			CompetitorID: 1,
		}).Return(noEffects)
		mockedEngine.On("HandleCompetitorRequalified", domain.CompetitorRequalifiedEvent{
			CompetitorID: 1,
		}).Return(noEffects)
		mockedEngine.On("HandleAscentRegistered", domain.AscentRegisteredEvent{
			CompetitorID:  1,
			ProblemID:     1,
			Top:           true,
			AttemptsTop:   999,
			Zone1:         true,
			AttemptsZone1: 7,
			Zone2:         true,
			AttemptsZone2: 42,
		}).Return(noEffects)
		mockedEngine.On("HandleAscentDeregistered", domain.AscentDeregisteredEvent{
			CompetitorID: 1,
			ProblemID:    1,
		}).Return(noEffects)
		mockedEngine.On("HandleProblemAdded", domain.ProblemAddedEvent{
			ProblemID: 1,
			ProblemValue: domain.ProblemValue{
				PointsTop:   1000,
				PointsZone1: 500,
				PointsZone2: 750,
				FlashBonus:  100,
			},
		}).Run(func(mock.Arguments) { cancel() }).Return(noEffects)

		installEngine(mockedEngine)

		wg.Wait()

		awaitExpectations(t)
		mockedEngine.AssertExpectations(t)
	})

	t.Run("PublishResults", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			f, awaitExpectations := makeFixture(0)

			ctx, cancel := context.WithCancel(context.Background())
			wg, installEngine := f.driver.Run(ctx)

			now := time.Now()

			mockedEngine := new(scoreEngineMock)

			mockedEngine.On("Start").Return(noEffects)
			mockedEngine.On("Stop").Return()
			mockedEngine.On("GetDirtyResults").Return([]domain.Result{
				{
					CompetitorID: 1,
					Timestamp:    now,
					Score:        "100p",
					Placement:    1,
					RankOrder:    0,
					Finalist:     true,
				},
				{
					CompetitorID: 2,
					Timestamp:    now,
					Score:        "200p",
					Placement:    2,
					RankOrder:    1,
					Finalist:     true,
				},
				{
					CompetitorID: 3,
					Timestamp:    now,
					Score:        "300p",
					Placement:    3,
					RankOrder:    2,
					Finalist:     false,
				},
			})
			mockedEngine.On("GetDirtyPointValues").Return([]domain.PointValue{})

			f.broker.
				On("Dispatch", fakedCompetitionID,
					domain.CompetitorResultUpdatedEvent{
						CompetitorID: 1,
						Timestamp:    now,
						Score:        "100p",
						Placement:    1,
						RankOrder:    0,
						Finalist:     true,
					},
				).Return().
				On("Dispatch", fakedCompetitionID,
					domain.CompetitorResultUpdatedEvent{
						CompetitorID: 2,
						Timestamp:    now,
						Score:        "200p",
						Placement:    2,
						RankOrder:    1,
						Finalist:     true,
					},
				).Return().
				On("Dispatch", fakedCompetitionID,
					domain.CompetitorResultUpdatedEvent{
						CompetitorID: 3,
						Timestamp:    now,
						Score:        "300p",
						Placement:    3,
						RankOrder:    2,
						Finalist:     false,
					},
				).Return().
				On("Dispatch", fakedCompetitionID,
					[]domain.CompetitorResultUpdatedEvent{
						{
							CompetitorID: 1,
							Timestamp:    now,
							Score:        "100p",
							Placement:    1,
							RankOrder:    0,
							Finalist:     true,
						},
						{
							CompetitorID: 2,
							Timestamp:    now,
							Score:        "200p",
							Placement:    2,
							RankOrder:    1,
							Finalist:     true,
						},
						{
							CompetitorID: 3,
							Timestamp:    now,
							Score:        "300p",
							Placement:    3,
							RankOrder:    2,
							Finalist:     false,
						},
					},
				).Return()

			installEngine(mockedEngine)
			synctest.Wait()

			cancel()

			wg.Wait()

			awaitExpectations(t)
			mockedEngine.AssertExpectations(t)
		})
	})
}

type scoreEngineMock struct {
	mock.Mock
}

func (m *scoreEngineMock) Start() iter.Seq[scores.Effect] {
	args := m.Called()
	return args.Get(0).(iter.Seq[scores.Effect])
}

func (m *scoreEngineMock) Stop() {
	m.Called()
}

func (m *scoreEngineMock) HandleRulesUpdated(event domain.RulesUpdatedEvent) iter.Seq[scores.Effect] {
	args := m.Called(event)
	return args.Get(0).(iter.Seq[scores.Effect])
}

func (m *scoreEngineMock) HandleCompetitorEntered(event domain.CompetitorEnteredEvent) iter.Seq[scores.Effect] {
	args := m.Called(event)
	return args.Get(0).(iter.Seq[scores.Effect])
}

func (m *scoreEngineMock) HandleCompetitorSwitchedCategory(event domain.CompetitorSwitchedCategoryEvent) iter.Seq[scores.Effect] {
	args := m.Called(event)
	return args.Get(0).(iter.Seq[scores.Effect])
}

func (m *scoreEngineMock) HandleCompetitorWithdrewFromFinals(event domain.CompetitorWithdrewFromFinalsEvent) iter.Seq[scores.Effect] {
	args := m.Called(event)
	return args.Get(0).(iter.Seq[scores.Effect])
}

func (m *scoreEngineMock) HandleCompetitorReenteredFinals(event domain.CompetitorReenteredFinalsEvent) iter.Seq[scores.Effect] {
	args := m.Called(event)
	return args.Get(0).(iter.Seq[scores.Effect])
}

func (m *scoreEngineMock) HandleCompetitorDisqualified(event domain.CompetitorDisqualifiedEvent) iter.Seq[scores.Effect] {
	args := m.Called(event)
	return args.Get(0).(iter.Seq[scores.Effect])
}

func (m *scoreEngineMock) HandleCompetitorRequalified(event domain.CompetitorRequalifiedEvent) iter.Seq[scores.Effect] {
	args := m.Called(event)
	return args.Get(0).(iter.Seq[scores.Effect])
}

func (m *scoreEngineMock) HandleAscentRegistered(event domain.AscentRegisteredEvent) iter.Seq[scores.Effect] {
	args := m.Called(event)
	return args.Get(0).(iter.Seq[scores.Effect])
}

func (m *scoreEngineMock) HandleAscentDeregistered(event domain.AscentDeregisteredEvent) iter.Seq[scores.Effect] {
	args := m.Called(event)
	return args.Get(0).(iter.Seq[scores.Effect])
}

func (m *scoreEngineMock) HandleProblemAdded(event domain.ProblemAddedEvent) iter.Seq[scores.Effect] {
	args := m.Called(event)
	return args.Get(0).(iter.Seq[scores.Effect])
}

func (m *scoreEngineMock) HandleProblemUpdated(event domain.ProblemUpdatedEvent) iter.Seq[scores.Effect] {
	args := m.Called(event)
	return args.Get(0).(iter.Seq[scores.Effect])
}

func (m *scoreEngineMock) RankCategory(categoryID domain.CategoryID) {
	m.Called(categoryID)
}

func (m *scoreEngineMock) ScoreCompetitor(competitorID domain.CompetitorID) iter.Seq[scores.Effect] {
	args := m.Called(competitorID)
	return args.Get(0).(iter.Seq[scores.Effect])
}

func (m *scoreEngineMock) CalculatePointValues(categoryID domain.CategoryID, problemID domain.ProblemID) iter.Seq[scores.Effect] {
	args := m.Called(categoryID, problemID)
	return args.Get(0).(iter.Seq[scores.Effect])
}

func (m *scoreEngineMock) GetDirtyResults() []domain.Result {
	args := m.Called()
	return args.Get(0).([]domain.Result)
}

func (m *scoreEngineMock) GetDirtyPointValues() []domain.PointValue {
	args := m.Called()
	return args.Get(0).([]domain.PointValue)
}
