package events_test

import (
	"testing"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/events"
)

func TestBlockingSubscriber(t *testing.T) {
	broker := events.NewBroker()
	filter := domain.EventFilter{
		ContestID: 1,
	}

	_, _ = broker.Subscribe(filter, 1)

	for range 100 {
		broker.Dispatch(1, domain.ContenderEnteredEvent{
			ContenderID: 1,
		})
	}
}
