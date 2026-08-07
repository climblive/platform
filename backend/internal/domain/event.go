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

func (f EventFilter) Match(contestID ContestID, contenderID ContenderID, eventType string) bool {
	switch f.ContestID {
	case 0, contestID:
	default:
		return false
	}

	switch f.ContenderID {
	case 0, contenderID:
	default:
		return false
	}

	hasEventTypeFilters := len(f.EventTypes) > 0

	if _, found := f.EventTypes[eventType]; hasEventTypeFilters && !found {
		return false
	}

	return true
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
	Subscribe(filters []EventFilter, bufferCapacity int) (SubscriptionID, EventReader)
	Unsubscribe(subscriptionID SubscriptionID)
}

type EventReader interface {
	EventsChan(ctx context.Context) <-chan EventEnvelope
}

type EventEnvelope struct {
	Data any
}
