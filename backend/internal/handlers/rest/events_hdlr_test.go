package rest_test

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/events"
	"github.com/climblive/platform/backend/internal/handlers/rest"
	"github.com/climblive/platform/backend/internal/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestEventsHandler(t *testing.T) {
	fakedContestID := testutils.RandomResourceID[domain.ContestID]()
	fakedContenderID := testutils.RandomResourceID[domain.ContenderID]()

	makeMocks := func(bufferCapacity int, filters []domain.EventFilter) (*repositoryMock, *eventBrokerMock, *events.Subscription) {
		mockedRepository := new(repositoryMock)

		mockedEventBroker := new(eventBrokerMock)

		subscription := events.NewSubscription(filters, bufferCapacity)
		subscriptionID := uuid.New()

		mockedEventBroker.On("Subscribe", filters, 1000).Return(subscriptionID, subscription)

		mockedEventBroker.On("Unsubscribe", subscriptionID).Return()

		return mockedRepository, mockedEventBroker, subscription
	}

	makeContenderFilters := func(contestID domain.ContestID, contenderID domain.ContenderID) []domain.EventFilter {
		return []domain.EventFilter{
			domain.NewEventFilter(
				contestID,
				contenderID,
				"CONTENDER_PUBLIC_INFO_UPDATED",
				"CONTENDER_SCORE_UPDATED",
				"ASCENT_REGISTERED",
				"ASCENT_DEREGISTERED",
				"RAFFLE_WINNER_DRAWN",
				"POINT_VALUE_UPDATED",
			),
			domain.NewEventFilter(
				contestID,
				0,
				"RULES_UPDATED",
				"PROBLEM_ADDED",
				"PROBLEM_UPDATED",
				"PROBLEM_DELETED",
			),
		}
	}

	t.Run("ConnectAndDisconnect", func(t *testing.T) {
		mockedRepository, mockedEventBroker, _ := makeMocks(0, makeContenderFilters(fakedContestID, fakedContenderID))

		mockedRepository.On("GetContender", mock.Anything, mock.Anything, fakedContenderID).Return(domain.Contender{
			ID:        fakedContenderID,
			ContestID: fakedContestID,
		}, nil)

		mux := rest.NewMux()
		rest.InstallEventHandler(mux, mockedEventBroker, mockedRepository, 0)

		server := httptest.NewServer(mux)

		resp, err := http.Get(server.URL + fmt.Sprintf("/contenders/%v/events", fakedContenderID))
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "text/event-stream", resp.Header.Get("Content-Type"))
		assert.Equal(t, "no", resp.Header.Get("X-Accel-Buffering"))
		assert.Equal(t, "no-store", resp.Header.Get("Cache-Control"))
		assert.Equal(t, "keep-alive", resp.Header.Get("Connection"))

		buf := bufio.NewReader(resp.Body)
		line, _, err := buf.ReadLine()
		require.NoError(t, err)
		assert.Equal(t, "retry: 5000", string(line))

		_ = resp.Body.Close()

		server.Close()

		mockedEventBroker.AssertExpectations(t)
	})

	t.Run("ReceivePing", func(t *testing.T) {
		mockedRepository, mockedEventBroker, _ := makeMocks(0, makeContenderFilters(fakedContestID, fakedContenderID))

		mockedRepository.On("GetContender", mock.Anything, mock.Anything, fakedContenderID).Return(domain.Contender{
			ID:        fakedContenderID,
			ContestID: fakedContestID,
		}, nil)

		mux := rest.NewMux()
		rest.InstallEventHandler(mux, mockedEventBroker, mockedRepository, time.Millisecond)

		server := httptest.NewServer(mux)

		resp, err := http.Get(server.URL + fmt.Sprintf("/contenders/%v/events", fakedContenderID))
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		buf := bufio.NewReader(resp.Body)

		var lines []string

		for i := 0; i < 3; i++ {
			line, _, err := buf.ReadLine()
			require.NoError(t, err)

			lines = append(lines, string(line))
		}

		assert.Equal(t, []string{"retry: 5000", "", ":"}, lines)

		_ = resp.Body.Close()

		server.Close()

		mockedEventBroker.AssertExpectations(t)
	})

	t.Run("ReceiveEvent", func(t *testing.T) {
		mockedRepository, mockedEventBroker, subscription := makeMocks(0, makeContenderFilters(fakedContestID, fakedContenderID))

		mockedRepository.On("GetContender", mock.Anything, mock.Anything, fakedContenderID).Return(domain.Contender{
			ID:        fakedContenderID,
			ContestID: fakedContestID,
		}, nil)

		err := subscription.Post(domain.EventEnvelope{
			Data: domain.ContenderScoreUpdatedEvent{
				Timestamp:   time.Date(2024, 12, 01, 00, 00, 00, 0, time.UTC),
				ContenderID: fakedContenderID,
				Score:       "100p",
				Placement:   10,
				Finalist:    true,
				RankOrder:   9,
			},
		})
		require.NoError(t, err)

		mux := rest.NewMux()
		rest.InstallEventHandler(mux, mockedEventBroker, mockedRepository, time.Hour)

		server := httptest.NewServer(mux)

		resp, err := http.Get(server.URL + fmt.Sprintf("/contenders/%v/events", fakedContenderID))
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		buf := bufio.NewReader(resp.Body)

		var lines []string

		for i := 0; i < 4; i++ {
			line, _, err := buf.ReadLine()
			require.NoError(t, err)

			lines = append(lines, string(line))
		}

		assert.Equal(t, []string{
			"retry: 5000",
			"",
			"event: CONTENDER_SCORE_UPDATED",
			fmt.Sprintf(`data: {"timestamp":"2024-12-01T00:00:00Z","contenderId":%v,"score":"100p","placement":10,"finalist":true,"rankOrder":9}`, fakedContenderID),
		}, lines)

		_ = resp.Body.Close()

		server.Close()

		mockedEventBroker.AssertExpectations(t)
	})

	t.Run("SubscriptionUnexpectedlyClosed", func(t *testing.T) {
		mockedRepository, mockedEventBroker, subscription := makeMocks(1, makeContenderFilters(fakedContestID, fakedContenderID))

		mockedRepository.On("GetContender", mock.Anything, mock.Anything, fakedContenderID).Return(domain.Contender{
			ID:        fakedContenderID,
			ContestID: fakedContestID,
		}, nil)

		err := subscription.Post(domain.EventEnvelope{
			Data: domain.ContenderScoreUpdatedEvent{},
		})
		require.NoError(t, err)

		err = subscription.Post(domain.EventEnvelope{
			Data: domain.ContenderScoreUpdatedEvent{},
		})
		require.ErrorIs(t, err, events.ErrBufferFull)

		mux := rest.NewMux()
		rest.InstallEventHandler(mux, mockedEventBroker, mockedRepository, 0)

		server := httptest.NewServer(mux)

		resp, err := http.Get(server.URL + fmt.Sprintf("/contenders/%v/events", fakedContenderID))
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		server.Close()

		mockedEventBroker.AssertExpectations(t)
	})

	t.Run("ContestEvents", func(t *testing.T) {
		mockedRepository, mockedEventBroker, _ := makeMocks(0, []domain.EventFilter{domain.NewEventFilter(
			fakedContestID,
			0,
			"CONTENDER_PUBLIC_INFO_UPDATED",
			"[]CONTENDER_SCORE_UPDATED",
			"SCORE_ENGINE_STARTED",
			"SCORE_ENGINE_STOPPED",
		)})

		mux := rest.NewMux()
		rest.InstallEventHandler(mux, mockedEventBroker, mockedRepository, time.Hour)

		server := httptest.NewServer(mux)

		resp, err := http.Get(server.URL + fmt.Sprintf("/contests/%v/events", fakedContestID))
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		_ = resp.Body.Close()

		server.Close()

		mockedEventBroker.AssertExpectations(t)
	})
}

type eventBrokerMock struct {
	mock.Mock
}

func (m *eventBrokerMock) Dispatch(contestID domain.ContestID, event any) {
	m.Called(contestID, event)
}

func (m *eventBrokerMock) Subscribe(filters []domain.EventFilter, bufferCapacity int) (domain.SubscriptionID, domain.EventReader) {
	args := m.Called(filters, bufferCapacity)
	return args.Get(0).(domain.SubscriptionID), args.Get(1).(domain.EventReader)
}

func (m *eventBrokerMock) Unsubscribe(subscriptionID domain.SubscriptionID) {
	m.Called(subscriptionID)
}

type repositoryMock struct {
	mock.Mock
}

func (m *repositoryMock) GetContender(ctx context.Context, tx domain.Transaction, contenderID domain.ContenderID) (domain.Contender, error) {
	args := m.Called(ctx, tx, contenderID)
	return args.Get(0).(domain.Contender), args.Error(1)
}
