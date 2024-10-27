package events

import (
	"sync"

	"github.com/climblive/platform/backend/internal/domain"
)

type broker struct {
	mu            sync.RWMutex
	subscriptions map[domain.SubscriptionID]*Subscription
}

func NewBroker() domain.EventBroker {
	return &broker{
		mu:            sync.RWMutex{},
		subscriptions: make(map[domain.SubscriptionID]*Subscription),
	}
}

func (b *broker) Subscribe(filter domain.EventFilter, bufferCapacity int) (domain.SubscriptionID, domain.EventReader) {
	b.mu.Lock()
	defer b.mu.Unlock()

	subscription := NewSubscription(filter, bufferCapacity)

	b.subscriptions[subscription.ID] = subscription

	return subscription.ID, subscription
}

func (b *broker) Unsubscribe(subscriptionID domain.SubscriptionID) {
	b.mu.Lock()
	defer b.mu.Unlock()

	delete(b.subscriptions, subscriptionID)
}

func (b *broker) Dispatch(contestID domain.ContestID, event any) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	eventName := eventName(event)

	for _, subscription := range b.subscriptions {
		switch subscription.Filter.ContestID {
		case 0:
		case contestID:
		default:
			continue
		}

		subscription.Post(domain.EventEnvelope{
			Name: eventName,
			Data: event,
		})
	}
}

func eventName(event any) string {
	switch event.(type) {
	case domain.ContenderEnteredEvent:
		return "CONTENDER_ENTERED"
	case domain.ContenderSwitchedClassEvent:
		return "CONTENDER_SWITCHED_CLASS"
	case domain.ContenderWithdrewFromFinalsEvent:
		return "CONTENDER_WITHDREW_FROM_FINALS"
	case domain.ContenderReenteredFinalsEvent:
		return "CONTENDER_REENTERED_FINALS"
	case domain.ContenderDisqualifiedEvent:
		return "CONTENDER_DISQUALIFIED"
	case domain.ContenderRequalifiedEvent:
		return "CONTENDER_REQUALIFIED"
	case domain.AscentRegisteredEvent:
		return "ASCENT_REGISTERED"
	case domain.AscentDeregisteredEvent:
		return "ASCENT_DEREGISTERED"
	case domain.ProblemAddedEvent:
		return "PROBLEM_ADDED"
	case domain.ProblemUpdatedEvent:
		return "PROBLEM_UPDATED"
	case domain.ProblemDeletedEvent:
		return "PROBLEM_DELETED"
	case domain.ContenderPublicInfoUpdatedEvent:
		return "CONTENDER_PUBLIC_INFO_UPDATED"
	case domain.ContenderScoreUpdatedEvent:
		return "CONTENDER_SCORE_UPDATED"
	default:
		return "UNKNOWN"
	}
}
