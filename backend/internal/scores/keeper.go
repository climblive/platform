package scores

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
)

const persistInterval = time.Minute
const lastDitchPersistTimeout = 5 * time.Second

type keeperRepository interface {
	domain.Transactor

	StoreScore(ctx context.Context, tx domain.Transaction, score domain.Score) (domain.Score, error)
}

type Keeper struct {
	mu          sync.RWMutex
	eventBroker domain.EventBroker
	scores      map[domain.ContenderID]domain.Score
	repo        keeperRepository
}

func NewScoreKeeper(eventBroker domain.EventBroker, repo keeperRepository) *Keeper {
	return &Keeper{
		eventBroker: eventBroker,
		scores:      make(map[domain.ContenderID]domain.Score),
		repo:        repo,
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

	events := eventReader.EventsChan(ctx)
	ticker := time.Tick(persistInterval)

ConsumeEvents:
	for {
		select {
		case event, open := <-events:
			if !open {
				break ConsumeEvents
			}

			switch ev := event.Data.(type) {
			case domain.ContenderScoreUpdatedEvent:
				k.HandleContenderScoreUpdated(ev)
			}
		case <-ticker:
			k.persistScores(ctx)
		case <-ctx.Done():
			slog.Info("subscription closed", "reason", ctx.Err())
			break ConsumeEvents
		}
	}

	if ctx.Err() == nil {
		slog.Warn("subscription closed unexpectedly")
	}

	slog.Info("score keeper shutting down")

	if len(k.scores) > 0 {
		ctxWithDeadline, cancel := context.WithTimeout(context.Background(), lastDitchPersistTimeout)
		defer cancel()

		slog.Warn("making a last-ditch attempt to persist scores", "timeout", lastDitchPersistTimeout)
		k.persistScores(ctxWithDeadline)
	}
}

func (k *Keeper) persistScores(ctx context.Context) {
	takeFirst := func() (domain.ContenderID, domain.Score) {
		k.mu.Lock()
		defer k.mu.Unlock()

		var contenderID domain.ContenderID
		var score domain.Score
		for contenderID, score = range k.scores {
			break
		}

		if contenderID == 0 {
			return 0, domain.Score{}
		}

		delete(k.scores, contenderID)

		return contenderID, score
	}

	putBack := func(contenderID domain.ContenderID, score domain.Score) {
		k.mu.Lock()
		defer k.mu.Unlock()

		if _, found := k.scores[contenderID]; !found {
			k.scores[contenderID] = score
		}
	}

	for ctx.Err() == nil {
		contenderID, score := takeFirst()

		if contenderID == 0 {
			break
		}

		_, err := k.repo.StoreScore(ctx, nil, score)
		if err != nil {
			slog.Error("failed to persist score",
				"contender_id", contenderID,
				"error", err)

			putBack(contenderID, score)
		}
	}

	if ctx.Err() != nil && len(k.scores) > 0 {
		slog.Warn("not all scores where persisted",
			"reason", ctx.Err(),
			"left_in_memory", len(k.scores),
		)
	}
}

func (k *Keeper) HandleContenderScoreUpdated(event domain.ContenderScoreUpdatedEvent) {
	k.mu.Lock()
	defer k.mu.Unlock()

	k.scores[event.ContenderID] = domain.Score(event)
}

func (k *Keeper) GetScore(contenderID domain.ContenderID) (domain.Score, error) {
	k.mu.RLock()
	defer k.mu.RUnlock()

	if score, found := k.scores[contenderID]; found {
		return score, nil
	}

	return domain.Score{}, domain.ErrNotFound
}
