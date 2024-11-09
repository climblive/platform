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
		mockedEventBroker := new(eventBrokerMock)

		mockedRepo.
			On("GetContestsRunningOrAboutToStart", mock.Anything, mock.Anything, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).
			Return([]domain.Contest{}, nil)

		mngr := scores.NewScoreEngineManager(mockedRepo, mockedEventBroker)

		ctx, cancel := context.WithCancel(context.Background())

		wg := mngr.Run(ctx)

		mockedRepo.AssertExpectations(t)
		mockedEventBroker.AssertExpectations(t)

		cancel()
		wg.Wait()
	})

	t.Run("LoadSingleEngine", func(t *testing.T) {
		mockedRepo := new(repositoryMock)
		mockedEventBroker := new(eventBrokerMock)

		mockedContestID := domain.ContestID(rand.Int())
		mockedProblemID := domain.ProblemID(rand.Int())
		mockedContenderID := domain.ContenderID(rand.Int())
		mockedCompClassID := domain.CompClassID(rand.Int())

		mockedEventBroker.
			On("Subscribe", mock.Anything, mock.Anything).
			Return(domain.SubscriptionID(uuid.New()), events.NewSubscription(domain.EventFilter{}, 1000))

		mockedRepo.
			On("GetContestsRunningOrAboutToStart", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]domain.Contest{
				{
					ID:                 mockedContestID,
					QualifyingProblems: 10,
					Finalists:          7,
				},
			}, nil)

		mockedRepo.
			On("GetProblemsByContest", mock.Anything, mock.Anything, mockedContestID).
			Return([]domain.Problem{
				{
					ID:         mockedProblemID,
					PointsTop:  100,
					PointsZone: 50,
					FlashBonus: 10,
				},
			}, nil)

		mockedRepo.
			On("GetContendersByContest", mock.Anything, mock.Anything, mockedContestID).
			Return([]domain.Contender{
				{
					ID:                  mockedContenderID,
					CompClassID:         mockedCompClassID,
					Disqualified:        true,
					WithdrawnFromFinals: true,
				},
			}, nil)

		mockedRepo.
			On("GetTicksByContest", mock.Anything, mock.Anything, mockedContestID).
			Return([]domain.Tick{
				{
					Ownership: domain.OwnershipData{
						ContenderID: &mockedContenderID,
					},
					ProblemID:    mockedProblemID,
					Top:          true,
					AttemptsTop:  999,
					Zone:         true,
					AttemptsZone: 1,
				},
			}, nil)

		mockedEventBroker.On("Dispatch", mockedContestID, domain.ProblemAddedEvent{
			ProblemID:  mockedProblemID,
			PointsTop:  100,
			PointsZone: 50,
			FlashBonus: 10,
		}).Return()

		mockedEventBroker.On("Dispatch", mockedContestID, domain.ContenderEnteredEvent{
			ContenderID: mockedContenderID,
			CompClassID: mockedCompClassID,
		}).Return()

		mockedEventBroker.On("Dispatch", mockedContestID, domain.ContenderWithdrewFromFinalsEvent{
			ContenderID: mockedContenderID,
		}).Return()

		mockedEventBroker.On("Dispatch", mockedContestID, domain.ContenderDisqualifiedEvent{
			ContenderID: mockedContenderID,
		}).Return()

		mockedEventBroker.On("Dispatch", mockedContestID, domain.AscentRegisteredEvent{
			ContenderID:  mockedContenderID,
			ProblemID:    mockedProblemID,
			Top:          true,
			AttemptsTop:  999,
			Zone:         true,
			AttemptsZone: 1,
		}).Return()

		mngr := scores.NewScoreEngineManager(mockedRepo, mockedEventBroker)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		wg := mngr.Run(ctx)

		mockedRepo.AssertExpectations(t)
		mockedEventBroker.AssertExpectations(t)

		cancel()
		wg.Wait()
	})
}

type repositoryMock struct {
	mock.Mock
}

func (m *repositoryMock) GetContestsRunningOrAboutToStart(ctx context.Context, tx domain.Transaction, earliestStartTime, latestStartTime time.Time) ([]domain.Contest, error) {
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
