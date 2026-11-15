package events

import (
	"log/slog"
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

func (b *broker) Subscribe(filters []domain.EventFilter, bufferCapacity int) (domain.SubscriptionID, domain.EventReader) {
	b.mu.Lock()
	defer b.mu.Unlock()

	subscription := NewSubscription(filters, bufferCapacity)

	b.subscriptions[subscription.ID] = subscription

	return subscription.ID, subscription
}

func (b *broker) Unsubscribe(subscriptionID domain.SubscriptionID) {
	b.mu.Lock()
	defer b.mu.Unlock()

	delete(b.subscriptions, subscriptionID)
}

func (b *broker) Dispatch(competitionID domain.CompetitionID, event any) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	eventName := EventName(event)
	competitorID := extractCompetitorID(event)

	for _, subscription := range b.subscriptions {
		if !subscription.FiltersMatch(competitionID, competitorID, eventName) {
			continue
		}

		err := subscription.Post(domain.EventEnvelope{
			Data: event,
		})

		if err != nil {
			slog.Error("failed to post event", "subscription_id", subscription.ID, "error", err)
		}
	}
}

func EventName(event any) string {
	switch event.(type) {
	case domain.CompetitorEnteredEvent:
		return "COMPETITOR_ENTERED"
	case domain.CompetitorSwitchedCategoryEvent:
		return "COMPETITOR_SWITCHED_CATEGORY"
	case domain.CompetitorWithdrewFromFinalsEvent:
		return "COMPETITOR_WITHDREW_FROM_FINALS"
	case domain.CompetitorReenteredFinalsEvent:
		return "COMPETITOR_REENTERED_FINALS"
	case domain.CompetitorDisqualifiedEvent:
		return "COMPETITOR_DISQUALIFIED"
	case domain.CompetitorRequalifiedEvent:
		return "COMPETITOR_REQUALIFIED"
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
	case domain.RulesUpdatedEvent:
		return "RULES_UPDATED"
	case domain.CompetitorPublicInfoUpdatedEvent:
		return "COMPETITOR_PUBLIC_INFO_UPDATED"
	case domain.CompetitorResultUpdatedEvent:
		return "COMPETITOR_RESULT_UPDATED"
	case []domain.CompetitorResultUpdatedEvent:
		return "[]COMPETITOR_RESULT_UPDATED"
	case domain.ScoreEngineStartedEvent:
		return "SCORE_ENGINE_STARTED"
	case domain.ScoreEngineStoppedEvent:
		return "SCORE_ENGINE_STOPPED"
	case domain.PointValueUpdatedEvent:
		return "POINT_VALUE_UPDATED"
	case domain.RaffleWinnerDrawnEvent:
		return "RAFFLE_WINNER_DRAWN"
	default:
		return "UNKNOWN"
	}
}

func extractCompetitorID(event any) domain.CompetitorID {
	switch ev := event.(type) {
	case domain.CompetitorEnteredEvent:
		return ev.CompetitorID
	case domain.CompetitorSwitchedCategoryEvent:
		return ev.CompetitorID
	case domain.CompetitorWithdrewFromFinalsEvent:
		return ev.CompetitorID
	case domain.CompetitorReenteredFinalsEvent:
		return ev.CompetitorID
	case domain.CompetitorDisqualifiedEvent:
		return ev.CompetitorID
	case domain.CompetitorRequalifiedEvent:
		return ev.CompetitorID
	case domain.AscentRegisteredEvent:
		return ev.CompetitorID
	case domain.AscentDeregisteredEvent:
		return ev.CompetitorID
	case domain.CompetitorPublicInfoUpdatedEvent:
		return ev.CompetitorID
	case domain.CompetitorResultUpdatedEvent:
		return ev.CompetitorID
	case domain.PointValueUpdatedEvent:
		return ev.CompetitorID
	case domain.RaffleWinnerDrawnEvent:
		return ev.CompetitorID
	default:
		return 0
	}
}
