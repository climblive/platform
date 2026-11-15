package events_test

import (
	"context"
	"testing"
	"testing/synctest"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContenderSwitchedCategory(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()

		broker := events.NewBroker()
		filter := domain.NewEventFilter(1, 2, "CONTENDER_SWITCHED_CATEGORY")
		subscriptionID, reader := broker.Subscribe([]domain.EventFilter{filter}, 1)
		defer broker.Unsubscribe(subscriptionID)
		channel := reader.EventsChan(ctx)

		event := domain.ContenderSwitchedCategoryEvent{ContenderID: 2, CategoryID: 3}
		assert.Equal(t, "CONTENDER_SWITCHED_CATEGORY", events.EventName(event))
		broker.Dispatch(1, event)
		synctest.Wait()

		select {
		case envelope, ok := <-channel:
			require.True(t, ok)
			assert.Equal(t, event, envelope.Data)
		default:
			t.Fatal("category switch event was not delivered")
		}
	})
}

func TestBlockingSubscriber(t *testing.T) {
	broker := events.NewBroker()
	filter := domain.EventFilter{
		CompetitionID: 1,
	}

	_, _ = broker.Subscribe([]domain.EventFilter{filter}, 1)

	for range 100 {
		broker.Dispatch(1, domain.ContenderEnteredEvent{
			ContenderID: 1,
		})
	}
}
