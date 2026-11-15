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
			"CONTENDER_ENTERED",
			"CONTENDER_SWITCHED_CATEGORY",
			"CONTENDER_WITHDREW_FROM_FINALS",
			"CONTENDER_REENTERED_FINALS",
			"CONTENDER_DISQUALIFIED",
			"CONTENDER_REQUALIFIED",
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
			Data: domain.ContenderResultUpdatedEvent{},
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
				CategoryID:  1,
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
			ContenderID: 1,
			ProblemID:   1,
			Top:         true,
			AttemptsTop: 10,
		}).Return(noEffects)

		mockedEngine.On("HandleAscentDeregistered", domain.AscentDeregisteredEvent{
			ContenderID: 1,
			ProblemID:   2,
		}).Return(noEffects)

		mockedEngine.On("HandleContenderEntered", domain.ContenderEnteredEvent{
			ContenderID: 2,
			CategoryID:  1,
		}).Return(noEffects)

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

		events = append(events, domain.RulesUpdatedEvent{
			CompetitionID:      1,
			QualifyingProblems: 10,
			Finalists:          7,
		})
		events = append(events, domain.ContenderEnteredEvent{
			ContenderID: 1,
			CategoryID:  1,
		})
		events = append(events, domain.ContenderSwitchedCategoryEvent{
			ContenderID: 1,
			CategoryID:  1,
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
			ContenderID:   1,
			ProblemID:     1,
			Top:           true,
			AttemptsTop:   999,
			Zone1:         true,
			AttemptsZone1: 7,
			Zone2:         true,
			AttemptsZone2: 42,
		})
		events = append(events, domain.AscentDeregisteredEvent{
			ContenderID: 1,
			ProblemID:   1,
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
		mockedEngine.On("HandleContenderEntered", domain.ContenderEnteredEvent{
			ContenderID: 1,
			CategoryID:  1,
		}).Return(noEffects)
		mockedEngine.On("HandleContenderSwitchedCategory", domain.ContenderSwitchedCategoryEvent{
			ContenderID: 1,
			CategoryID:  1,
		}).Return(noEffects)
		mockedEngine.On("HandleContenderWithdrewFromFinals", domain.ContenderWithdrewFromFinalsEvent{
			ContenderID: 1,
		}).Return(noEffects)
		mockedEngine.On("HandleContenderReenteredFinals", domain.ContenderReenteredFinalsEvent{
			ContenderID: 1,
		}).Return(noEffects)
		mockedEngine.On("HandleContenderDisqualified", domain.ContenderDisqualifiedEvent{
			ContenderID: 1,
		}).Return(noEffects)
		mockedEngine.On("HandleContenderRequalified", domain.ContenderRequalifiedEvent{
			ContenderID: 1,
		}).Return(noEffects)
		mockedEngine.On("HandleAscentRegistered", domain.AscentRegisteredEvent{
			ContenderID:   1,
			ProblemID:     1,
			Top:           true,
			AttemptsTop:   999,
			Zone1:         true,
			AttemptsZone1: 7,
			Zone2:         true,
			AttemptsZone2: 42,
		}).Return(noEffects)
		mockedEngine.On("HandleAscentDeregistered", domain.AscentDeregisteredEvent{
			ContenderID: 1,
			ProblemID:   1,
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
					ContenderID: 1,
					Timestamp:   now,
					Score:       "100p",
					Placement:   1,
					RankOrder:   0,
					Finalist:    true,
				},
				{
					ContenderID: 2,
					Timestamp:   now,
					Score:       "200p",
					Placement:   2,
					RankOrder:   1,
					Finalist:    true,
				},
				{
					ContenderID: 3,
					Timestamp:   now,
					Score:       "300p",
					Placement:   3,
					RankOrder:   2,
					Finalist:    false,
				},
			})
			mockedEngine.On("GetDirtyPointValues").Return([]domain.PointValue{})

			f.broker.
				On("Dispatch", fakedCompetitionID,
					domain.ContenderResultUpdatedEvent{
						ContenderID: 1,
						Timestamp:   now,
						Score:       "100p",
						Placement:   1,
						RankOrder:   0,
						Finalist:    true,
					},
				).Return().
				On("Dispatch", fakedCompetitionID,
					domain.ContenderResultUpdatedEvent{
						ContenderID: 2,
						Timestamp:   now,
						Score:       "200p",
						Placement:   2,
						RankOrder:   1,
						Finalist:    true,
					},
				).Return().
				On("Dispatch", fakedCompetitionID,
					domain.ContenderResultUpdatedEvent{
						ContenderID: 3,
						Timestamp:   now,
						Score:       "300p",
						Placement:   3,
						RankOrder:   2,
						Finalist:    false,
					},
				).Return().
				On("Dispatch", fakedCompetitionID,
					[]domain.ContenderResultUpdatedEvent{
						{
							ContenderID: 1,
							Timestamp:   now,
							Score:       "100p",
							Placement:   1,
							RankOrder:   0,
							Finalist:    true,
						},
						{
							ContenderID: 2,
							Timestamp:   now,
							Score:       "200p",
							Placement:   2,
							RankOrder:   1,
							Finalist:    true,
						},
						{
							ContenderID: 3,
							Timestamp:   now,
							Score:       "300p",
							Placement:   3,
							RankOrder:   2,
							Finalist:    false,
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

func (m *scoreEngineMock) HandleContenderEntered(event domain.ContenderEnteredEvent) iter.Seq[scores.Effect] {
	args := m.Called(event)
	return args.Get(0).(iter.Seq[scores.Effect])
}

func (m *scoreEngineMock) HandleContenderSwitchedCategory(event domain.ContenderSwitchedCategoryEvent) iter.Seq[scores.Effect] {
	args := m.Called(event)
	return args.Get(0).(iter.Seq[scores.Effect])
}

func (m *scoreEngineMock) HandleContenderWithdrewFromFinals(event domain.ContenderWithdrewFromFinalsEvent) iter.Seq[scores.Effect] {
	args := m.Called(event)
	return args.Get(0).(iter.Seq[scores.Effect])
}

func (m *scoreEngineMock) HandleContenderReenteredFinals(event domain.ContenderReenteredFinalsEvent) iter.Seq[scores.Effect] {
	args := m.Called(event)
	return args.Get(0).(iter.Seq[scores.Effect])
}

func (m *scoreEngineMock) HandleContenderDisqualified(event domain.ContenderDisqualifiedEvent) iter.Seq[scores.Effect] {
	args := m.Called(event)
	return args.Get(0).(iter.Seq[scores.Effect])
}

func (m *scoreEngineMock) HandleContenderRequalified(event domain.ContenderRequalifiedEvent) iter.Seq[scores.Effect] {
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

func (m *scoreEngineMock) ScoreContender(contenderID domain.ContenderID) iter.Seq[scores.Effect] {
	args := m.Called(contenderID)
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
