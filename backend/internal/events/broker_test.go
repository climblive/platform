package events_test

import (
	"testing"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/events"
	"github.com/google/uuid"
)

func TestBlockingSubscriber(t *testing.T) {
	broker := events.NewBroker()
	contestID := domain.ContestID(uuid.New())
	contenderID := domain.ContenderID(uuid.New())
	filter := domain.EventFilter{
		ContestID: contestID,
	}

	_, _ = broker.Subscribe(filter, 1)

	for range 100 {
		broker.Dispatch(contestID, domain.ContenderEnteredEvent{
			ContenderID: contenderID,
		})
	}
}
