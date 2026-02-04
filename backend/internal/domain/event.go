package domain

import (
	"context"
	"slices"

	"github.com/google/uuid"
)

type SubscriptionID = uuid.UUID

type EventFilter struct {
	ContestID   ContestID
	ContenderID ContenderID
	EventTypes  map[string]struct{}
}

func NewEventFilter(contestID ContestID, contenderID ContenderID, eventTypes ...string) EventFilter {
	filter := EventFilter{
		ContestID:   contestID,
		ContenderID: contenderID,
		EventTypes:  nil,
	}

	if len(eventTypes) > 0 {
		filter.EventTypes = map[string]struct{}{}
	}

	for eventType := range slices.Values(eventTypes) {
		filter.EventTypes[eventType] = struct{}{}
	}

	return filter
}

type EventBroker interface {
	Dispatch(contestID ContestID, event any)
	Subscribe(filter EventFilter, bufferCapacity int) (SubscriptionID, EventReader)
	Unsubscribe(subscriptionID SubscriptionID)
}

type EventReader interface {
	EventsChan(ctx context.Context) <-chan EventEnvelope
}

type EventEnvelope struct {
	Data any
}
