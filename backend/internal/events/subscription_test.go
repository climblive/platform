package events_test

import (
	"context"
	"math/rand"
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

	subscription.Post(domain.EventEnvelope{
		Name: "TEST",
		Data: randomNumber,
	})

	event, err := subscription.Await(context.Background())

	require.NoError(t, err)
	assert.Equal(t, randomNumber, event.Data)
}

func TestAwait(t *testing.T) {
	randomNumber := rand.Int()
	var wg sync.WaitGroup

	subscription := events.NewSubscription(domain.EventFilter{}, 1)

	wg.Add(1)

	var event domain.EventEnvelope
	var err error

	go func() {
		event, err = subscription.Await(context.Background())

		wg.Done()
	}()

	time.Sleep(100 * time.Millisecond)

	subscription.Post(domain.EventEnvelope{
		Name: "TEST",
		Data: randomNumber,
	})

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
			err := subscription.Post(domain.EventEnvelope{})
			assert.NoError(t, err)
		}

		err := subscription.Post(domain.EventEnvelope{})
		assert.ErrorIs(t, err, events.ErrBufferFull)

		wg.Done()

	}()

	wg.Wait()

	for range 2 {
		event, err := subscription.Await(context.Background())

		assert.Empty(t, event)
		require.ErrorIs(t, err, events.ErrBufferFull)
	}
}

func TestAwaitCancelled(t *testing.T) {
	subscription := events.NewSubscription(domain.EventFilter{}, 1)
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)

	var event domain.EventEnvelope
	var err error

	go func() {
		event, err = subscription.Await(ctx)

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

	event, err := subscription.Await(ctx)

	assert.Empty(t, event)
	require.ErrorIs(t, err, context.Canceled)
}
