package rest_test

import (
	"bufio"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/events"
	"github.com/climblive/platform/backend/internal/handlers/rest"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestEventsHandler(t *testing.T) {
	makeMocks := func(bufferCapacity int, filter domain.EventFilter) (*eventBrokerMock, *events.Subscription) {
		mockedEventBroker := new(eventBrokerMock)

		subscription := events.NewSubscription(domain.EventFilter{}, bufferCapacity)
		subscriptionID := uuid.New()

		mockedEventBroker.On("Subscribe", filter, 1000).Return(subscriptionID, subscription)

		mockedEventBroker.On("Unsubscribe", subscriptionID).Return()

		return mockedEventBroker, subscription
	}

	t.Run("ConnectAndDisconnect", func(t *testing.T) {
		contenderID := domain.ContenderID(uuid.New())
		mockedEventBroker, _ := makeMocks(0, domain.NewEventFilter(
			domain.ContestID(uuid.Nil),
			contenderID,
			"CONTENDER_PUBLIC_INFO_UPDATED",
			"CONTENDER_SCORE_UPDATED",
			"ASCENT_REGISTERED",
			"ASCENT_DEREGISTERED",
		))

		mux := rest.NewMux()
		rest.InstallEventHandler(mux, mockedEventBroker, 0)

		server := httptest.NewServer(mux)

		resp, err := http.Get(server.URL + "/contenders/" + uuid.UUID(contenderID).String() + "/events")
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
		contenderID := domain.ContenderID(uuid.New())
		mockedEventBroker, _ := makeMocks(0, domain.NewEventFilter(
			domain.ContestID(uuid.Nil),
			contenderID,
			"CONTENDER_PUBLIC_INFO_UPDATED",
			"CONTENDER_SCORE_UPDATED",
			"ASCENT_REGISTERED",
			"ASCENT_DEREGISTERED",
		))

		mux := rest.NewMux()
		rest.InstallEventHandler(mux, mockedEventBroker, time.Millisecond)

		server := httptest.NewServer(mux)

		resp, err := http.Get(server.URL + "/contenders/" + uuid.UUID(contenderID).String() + "/events")
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
		contenderID := domain.ContenderID(uuid.New())
		mockedEventBroker, subscription := makeMocks(0, domain.NewEventFilter(
			domain.ContestID(uuid.Nil),
			contenderID,
			"CONTENDER_PUBLIC_INFO_UPDATED",
			"CONTENDER_SCORE_UPDATED",
			"ASCENT_REGISTERED",
			"ASCENT_DEREGISTERED",
		))

		err := subscription.Post(domain.EventEnvelope{
			Data: domain.ContenderScoreUpdatedEvent{
				Timestamp:   time.Date(2024, 12, 01, 00, 00, 00, 0, time.UTC),
				ContenderID: contenderID,
				Score:       100,
				Placement:   10,
				Finalist:    true,
				RankOrder:   9,
			},
		})
		require.NoError(t, err)

		mux := rest.NewMux()
		rest.InstallEventHandler(mux, mockedEventBroker, time.Hour)

		server := httptest.NewServer(mux)

		resp, err := http.Get(server.URL + "/contenders/" + uuid.UUID(contenderID).String() + "/events")
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		buf := bufio.NewReader(resp.Body)

		var lines []string

		for i := 0; i < 4; i++ {
			line, _, err := buf.ReadLine()
			require.NoError(t, err)

			lines = append(lines, string(line))
		}

		assert.Contains(t, lines[0], "retry: 5000")
		assert.Equal(t, "", lines[1])
		assert.Equal(t, "event: CONTENDER_SCORE_UPDATED", lines[2])
		assert.Contains(t, lines[3], `"timestamp":"2024-12-01T00:00:00Z"`)
		assert.Contains(t, lines[3], `"score":100`)
		assert.Contains(t, lines[3], `"placement":10`)

		_ = resp.Body.Close()

		server.Close()

		mockedEventBroker.AssertExpectations(t)
	})

	t.Run("SubscriptionUnexpectedlyClosed", func(t *testing.T) {
		contenderID := domain.ContenderID(uuid.New())
		mockedEventBroker, subscription := makeMocks(1, domain.NewEventFilter(
			domain.ContestID(uuid.Nil),
			contenderID,
			"CONTENDER_PUBLIC_INFO_UPDATED",
			"CONTENDER_SCORE_UPDATED",
			"ASCENT_REGISTERED",
			"ASCENT_DEREGISTERED",
		))

		err := subscription.Post(domain.EventEnvelope{
			Data: domain.ContenderScoreUpdatedEvent{},
		})
		require.NoError(t, err)

		err = subscription.Post(domain.EventEnvelope{
			Data: domain.ContenderScoreUpdatedEvent{},
		})
		require.ErrorIs(t, err, events.ErrBufferFull)

		mux := rest.NewMux()
		rest.InstallEventHandler(mux, mockedEventBroker, 0)

		server := httptest.NewServer(mux)

		resp, err := http.Get(server.URL + "/contenders/" + uuid.UUID(contenderID).String() + "/events")
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		server.Close()

		mockedEventBroker.AssertExpectations(t)
	})

	t.Run("ContestEvents", func(t *testing.T) {
		contestID := domain.ContestID(uuid.New())
		mockedEventBroker, _ := makeMocks(0, domain.NewEventFilter(
			contestID,
			domain.ContenderID(uuid.Nil),
			"CONTENDER_PUBLIC_INFO_UPDATED",
			"[]CONTENDER_SCORE_UPDATED",
			"SCORE_ENGINE_STARTED",
			"SCORE_ENGINE_STOPPED",
		))

		mux := rest.NewMux()
		rest.InstallEventHandler(mux, mockedEventBroker, time.Hour)

		server := httptest.NewServer(mux)

		resp, err := http.Get(server.URL + "/contests/" + uuid.UUID(contestID).String() + "/events")
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

func (m *eventBrokerMock) Subscribe(filter domain.EventFilter, bufferCapacity int) (domain.SubscriptionID, domain.EventReader) {
	args := m.Called(filter, bufferCapacity)
	return args.Get(0).(domain.SubscriptionID), args.Get(1).(domain.EventReader)
}

func (m *eventBrokerMock) Unsubscribe(subscriptionID domain.SubscriptionID) {
	m.Called(subscriptionID)
}
