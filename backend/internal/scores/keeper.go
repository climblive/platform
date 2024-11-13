package scores

import (
	"context"
	"log/slog"
	"sync"

	"github.com/climblive/platform/backend/internal/domain"
)

type Keeper struct {
	mu          sync.RWMutex
	eventBroker domain.EventBroker
	scores      map[domain.ContenderID]domain.Score
}

func NewScoreKeeper(eventBroker domain.EventBroker) *Keeper {
	return &Keeper{
		eventBroker: eventBroker,
		scores:      make(map[domain.ContenderID]domain.Score),
	}
}

func (k *Keeper) Run(ctx context.Context) *sync.WaitGroup {
	wg := new(sync.WaitGroup)
	ready := make(chan struct{}, 1)

	filter := domain.NewEventFilter(
		0,
		0,
		"CONTENDER_SCORE_UPDATED",
	)

	wg.Add(1)

	go k.run(ctx, filter, wg, ready)

	<-ready

	return wg
}

func (k *Keeper) run(ctx context.Context, filter domain.EventFilter, wg *sync.WaitGroup, ready chan<- struct{}) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("score keeper panicked", "error", r)
		}
	}()

	defer wg.Done()

	subscriptionID, eventReader := k.eventBroker.Subscribe(filter, 0)
	defer k.eventBroker.Unsubscribe(subscriptionID)

	close(ready)

	for {
		event, err := eventReader.AwaitEvent(ctx)
		switch err {
		case nil:
		case context.Canceled, context.DeadlineExceeded:
			slog.Info("score keeper shutting down", "reason", err.Error())

			return
		default:
			panic(err)
		}

		switch ev := event.Data.(type) {
		case domain.ContenderScoreUpdatedEvent:
			k.HandleContenderScoreUpdated(ev)
		}
	}
}

func (k *Keeper) HandleContenderScoreUpdated(event domain.ContenderScoreUpdatedEvent) {
	k.mu.Lock()
	defer k.mu.Unlock()

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
	k.mu.RLock()
	defer k.mu.RUnlock()

	if score, found := k.scores[contenderID]; found {
		return score, nil
	}

	return domain.Score{}, domain.ErrNotFound
}
