package domain

import (
	"context"
	"slices"

	"uuid"
)

type SubscriptionID = uuid.UUID

func (f EventFilter) Match(competitionID CompetitionID, competitorID CompetitorID, eventType string) bool {
	switch f.CompetitionID {
	case 0, competitionID:
	default:
		return false
	}

	switch f.CompetitorID {
	case 0, competitorID:
	default:
		return false
	}

	hasEventTypeFilters := len(f.EventTypes) > 0

	if _, found := f.EventTypes[eventType]; hasEventTypeFilters && !found {
		return false
	}

	return true
}

func NewEventFilter(competitionID CompetitionID, competitorID CompetitorID, eventTypes ...string) EventFilter {
	filter := EventFilter{
		CompetitionID: competitionID,
		CompetitorID:  competitorID,
		EventTypes:    nil,
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
	Dispatch(competitionID CompetitionID, event any)
	Subscribe(filters []EventFilter, bufferCapacity int) (SubscriptionID, EventReader)
	Unsubscribe(subscriptionID SubscriptionID)
}

type EventReader interface {
	EventsChan(ctx context.Context) <-chan EventEnvelope
}

type EventEnvelope struct {
	Data any
}
