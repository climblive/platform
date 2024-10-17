package scores

import (
	"context"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/events"
)

type Keeper struct {
	eventBroker domain.EventBroker
	scores      map[domain.ResourceID]domain.Score
}

func NewScoreKeeper(eventBroker domain.EventBroker) *Keeper {
	return &Keeper{
		eventBroker: eventBroker,
		scores:      make(map[int]domain.Score),
	}
}

func (k *Keeper) Run(ctx context.Context) {
	events := make(chan domain.EventContainer, events.EventChannelBufferSize)
	subscriptionID := k.eventBroker.Subscribe(domain.EventFilter{}, events)

	defer k.eventBroker.Unsubscribe(subscriptionID)

	for {
		select {
		case event := <-events:
			switch ev := event.Data.(type) {
			case domain.ContenderScoreUpdatedEvent:
				k.HandleContenderScoreUpdated(ev)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (k *Keeper) HandleContenderScoreUpdated(event domain.ContenderScoreUpdatedEvent) {
	k.scores[event.ContenderID] = domain.Score{
		Timestamp:   event.Timestamp,
		ContenderID: event.ContenderID,
		Score:       event.Score,
		Placement:   event.Placement,
		Finalist:    event.Finalist,
		RankOrder:   event.RankOrder,
	}
}

func (k *Keeper) GetScore(contenderID domain.ResourceID) (domain.Score, error) {
	if score, found := k.scores[contenderID]; found {
		return score, nil
	}

	return domain.Score{}, domain.ErrNotFound
}
