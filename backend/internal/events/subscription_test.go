package events_test

import (
	"context"
	"math/rand"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostAndReceive(t *testing.T) {
	randomNumber := rand.Int()

	subscription := events.NewSubscription(domain.EventFilter{}, 1)

	err := subscription.Post(domain.EventEnvelope{
		Data: randomNumber,
	})
	require.NoError(t, err)

	event, err := subscription.AwaitEvent(context.Background())

	require.NoError(t, err)
	assert.Equal(t, randomNumber, event.Data)
}

func TestFIFO(t *testing.T) {
	subscription := events.NewSubscription(domain.EventFilter{}, 3)

	for k := 1; k <= 3; k++ {
		err := subscription.Post(domain.EventEnvelope{
			Data: k,
		})

		require.NoError(t, err)
	}

	for k := 1; k <= 3; k++ {
		event, err := subscription.AwaitEvent(context.Background())

		require.NoError(t, err)
		assert.Equal(t, k, event.Data)
	}
}

func TestAwait(t *testing.T) {
	randomNumber := rand.Int()
	var wg sync.WaitGroup

	subscription := events.NewSubscription(domain.EventFilter{}, 1)

	wg.Add(1)

	var event domain.EventEnvelope
	var err error

	go func() {
		event, err = subscription.AwaitEvent(context.Background())

		wg.Done()
	}()

	time.Sleep(100 * time.Millisecond)

	postErr := subscription.Post(domain.EventEnvelope{
		Data: randomNumber,
	})
	require.NoError(t, postErr)

	wg.Wait()

	require.NoError(t, err)
	assert.Equal(t, randomNumber, event.Data)
}

func TestBufferFull(t *testing.T) {
	bufferCapacity := 10
	var wg sync.WaitGroup

	subscription := events.NewSubscription(domain.EventFilter{}, bufferCapacity)

	wg.Add(1)

	go func() {
		for range bufferCapacity {
			err := subscription.Post(domain.EventEnvelope{
				Data: "Something",
			})
			assert.NoError(t, err)
		}

		err := subscription.Post(domain.EventEnvelope{})
		assert.ErrorIs(t, err, events.ErrBufferFull)

		wg.Done()
	}()

	wg.Wait()

	for range 10 {
		event, err := subscription.AwaitEvent(context.Background())

		assert.Equal(t, domain.EventEnvelope{Data: "Something"}, event)
		require.NoError(t, err)
	}

	event, err := subscription.AwaitEvent(context.Background())

	assert.Empty(t, event)
	require.ErrorIs(t, err, events.ErrBufferFull)
}

func TestTerminate(t *testing.T) {
	subscription := events.NewSubscription(domain.EventFilter{}, 0)

	err := subscription.Post(domain.EventEnvelope{
		Data: "Something",
	})
	assert.NoError(t, err)

	subscription.Terminate()

	event, err := subscription.AwaitEvent(context.Background())

	assert.Equal(t, domain.EventEnvelope{Data: "Something"}, event)
	require.NoError(t, err)

	event, err = subscription.AwaitEvent(context.Background())

	assert.Empty(t, event)
	require.ErrorIs(t, err, events.ErrTerminated)
}

func TestAwaitCancelled(t *testing.T) {
	subscription := events.NewSubscription(domain.EventFilter{}, 1)
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)

	var event domain.EventEnvelope
	var err error

	go func() {
		event, err = subscription.AwaitEvent(ctx)

		wg.Done()
	}()

	time.Sleep(100 * time.Millisecond)

	cancel()
	wg.Wait()

	assert.Empty(t, event)
	require.ErrorIs(t, err, context.Canceled)
}

func TestContextCancelledPreAwait(t *testing.T) {
	subscription := events.NewSubscription(domain.EventFilter{}, 1)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	event, err := subscription.AwaitEvent(ctx)

	assert.Empty(t, event)
	require.ErrorIs(t, err, context.Canceled)
}

func TestEventsChan(t *testing.T) {
	bufferCapacity := 10
	randomNumber := rand.Int()

	ctx, cancel := context.WithCancel(context.Background())

	subscription := events.NewSubscription(domain.EventFilter{}, bufferCapacity)

	go func() {
		for range bufferCapacity {
			err := subscription.Post(domain.EventEnvelope{
				Data: randomNumber,
			})
			assert.NoError(t, err)
		}
	}()

	events := subscription.EventsChan(ctx)

	for range bufferCapacity {
		event, open := <-events
		assert.Equal(t, domain.EventEnvelope{
			Data: randomNumber,
		}, event)
		assert.True(t, open)
	}

	cancel()

	event, open := <-events
	assert.Empty(t, event)
	assert.False(t, open)
}

func TestEventsChanBufferFull(t *testing.T) {
	bufferCapacity := 1

	ctx := context.Background()

	subscription := events.NewSubscription(domain.EventFilter{}, bufferCapacity)

	err := subscription.Post(domain.EventEnvelope{Data: "Something"})
	assert.NoError(t, err)

	err = subscription.Post(domain.EventEnvelope{Data: "Something"})
	assert.ErrorIs(t, err, events.ErrBufferFull)

	events := subscription.EventsChan(ctx)

	event, open := <-events
	assert.Equal(t, domain.EventEnvelope{Data: "Something"}, event)
	assert.True(t, open)

	event, open = <-events
	assert.Empty(t, event)
	assert.False(t, open)
}

func TestMatchFilter(t *testing.T) {
	t.Run("ContestMatchWildcard", func(t *testing.T) {
		subscription := events.NewSubscription(domain.NewEventFilter(0, 0), 0)

		match := subscription.FilterMatch(domain.ContestID(rand.Int()), 0, "A")

		assert.True(t, match)
	})

	t.Run("ContestMatch", func(t *testing.T) {
		subscription := events.NewSubscription(domain.NewEventFilter(1337, 0), 0)

		match := subscription.FilterMatch(domain.ContestID(1337), 0, "A")

		assert.True(t, match)
	})

	t.Run("ContestNoMatch", func(t *testing.T) {
		subscription := events.NewSubscription(domain.NewEventFilter(1337, 0), 0)

		match := subscription.FilterMatch(domain.ContestID(42), 0, "A")

		assert.False(t, match)
	})

	t.Run("ContenderMatchWildcard", func(t *testing.T) {
		subscription := events.NewSubscription(domain.NewEventFilter(0, 0), 0)

		match := subscription.FilterMatch(domain.ContestID(rand.Int()), domain.ContenderID(rand.Int()), "A")

		assert.True(t, match)
	})

	t.Run("ContenderMatch", func(t *testing.T) {
		subscription := events.NewSubscription(domain.NewEventFilter(0, 1337), 0)

		match := subscription.FilterMatch(domain.ContestID(rand.Int()), domain.ContenderID(1337), "A")

		assert.True(t, match)
	})

	t.Run("ContenderNoMatch", func(t *testing.T) {
		subscription := events.NewSubscription(domain.NewEventFilter(0, 1337), 0)

		match := subscription.FilterMatch(domain.ContestID(rand.Int()), domain.ContenderID(42), "A")

		assert.False(t, match)
	})

	t.Run("EventTypeMatch", func(t *testing.T) {
		subscription := events.NewSubscription(domain.NewEventFilter(0, 0, "A", "B", "C"), 0)

		for eventType := range slices.Values([]string{"A", "B", "C"}) {
			match := subscription.FilterMatch(domain.ContestID(rand.Int()), domain.ContenderID(rand.Int()), eventType)

			assert.True(t, match)
		}
	})

	t.Run("EventTypeNoMatch", func(t *testing.T) {
		subscription := events.NewSubscription(domain.NewEventFilter(0, 0, "A", "B", "C"), 0)

		match := subscription.FilterMatch(domain.ContestID(rand.Int()), domain.ContenderID(rand.Int()), "X")

		assert.False(t, match)
	})
}
