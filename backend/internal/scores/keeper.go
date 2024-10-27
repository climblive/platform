package scores

import (
	"context"

	"github.com/climblive/platform/backend/internal/domain"
)

type Keeper struct {
	eventBroker domain.EventBroker
	scores      map[domain.ContenderID]domain.Score
}

func NewScoreKeeper(eventBroker domain.EventBroker) *Keeper {
	return &Keeper{
		eventBroker: eventBroker,
		scores:      make(map[domain.ContenderID]domain.Score),
	}
}

func (k *Keeper) Run(ctx context.Context) {
	subscriptionID, eventReader := k.eventBroker.Subscribe(domain.EventFilter{}, 0)

	defer k.eventBroker.Unsubscribe(subscriptionID)

	for {
		event, err := eventReader.AwaitEvent(ctx)
		if err != nil {
			panic(err)
		}

		switch ev := event.Data.(type) {
		case domain.ContenderScoreUpdatedEvent:
			k.HandleContenderScoreUpdated(ev)
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

func (k *Keeper) GetScore(contenderID domain.ContenderID) (domain.Score, error) {
	if score, found := k.scores[contenderID]; found {
		return score, nil
	}

	return domain.Score{}, domain.ErrNotFound
}
