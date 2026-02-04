package scores

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

const persistInterval = time.Minute
const lastDitchPersistTimeout = 10 * time.Second

type keeperRepository interface {
	StoreScore(ctx context.Context, tx domain.Transaction, score domain.Score) error
}

type Keeper struct {
	mu                     sync.RWMutex
	eventBroker            domain.EventBroker
	scores                 map[domain.ContenderID]domain.Score
	repo                   keeperRepository
	externalPersistTrigger chan struct{}
}

func NewScoreKeeper(eventBroker domain.EventBroker, repo keeperRepository) *Keeper {
	return &Keeper{
		eventBroker:            eventBroker,
		scores:                 make(map[domain.ContenderID]domain.Score),
		repo:                   repo,
		externalPersistTrigger: make(chan struct{}, 1),
		mu:                     sync.RWMutex{},
	}
}

func (k *Keeper) Run(ctx context.Context, options ...func(*runOptions)) *sync.WaitGroup {
	config := &runOptions{
		recoverPanics: false,
	}
	for _, opt := range options {
		opt(config)
	}

	wg := new(sync.WaitGroup)
	ready := make(chan struct{}, 1)

	wg.Add(1)

	go func() {
		defer func() {
			if !config.recoverPanics {
				return
			}

			if r := recover(); r != nil {
				slog.Error("score keeper panicked", "error", r)
			}
		}()

		defer wg.Done()

		k.run(ctx, ready)
	}()

	<-ready

	return wg
}

func (k *Keeper) run(ctx context.Context, ready chan<- struct{}) {
	filter := domain.NewEventFilter(
		0,
		0,
		"CONTENDER_SCORE_UPDATED",
	)

	subscriptionID, eventReader := k.eventBroker.Subscribe(filter, 0)
	defer k.eventBroker.Unsubscribe(subscriptionID)

	close(ready)

	events := eventReader.EventsChan(ctx)
	ticker := time.Tick(persistInterval)

EventLoop:
	for {
		select {
		case event, open := <-events:
			if !open {
				break EventLoop
			}

			switch ev := event.Data.(type) {
			case domain.ContenderScoreUpdatedEvent:
				k.HandleContenderScoreUpdated(ev)
			}
		case <-ticker:
			k.persistScores(ctx)
		case <-k.externalPersistTrigger:
			k.persistScores(ctx)
		case <-ctx.Done():
			slog.Info("subscription closed", "reason", ctx.Err())
			break EventLoop
		}
	}

	if ctx.Err() == nil {
		slog.Error("subscription closed unexpectedly")
	}

	slog.Info("score keeper shutting down")

	numScores := k.getNumScoresWithLock()
	if numScores > 0 {
		ctxWithDeadline, cancel := context.WithTimeout(context.Background(), lastDitchPersistTimeout)
		defer cancel()

		slog.Warn("making a last-ditch attempt to persist scores", "timeout", lastDitchPersistTimeout)
		k.persistScores(ctxWithDeadline)
	}
}

func (k *Keeper) RequestPersist() {
	k.externalPersistTrigger <- struct{}{}
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
			return 0, domain.Score{
				Timestamp:   time.Time{},
				ContenderID: 0,
				Score:       0,
				Placement:   0,
				Finalist:    false,
				RankOrder:   0,
			}
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

	persistedScores := 0

IterateScores:
	for ctx.Err() == nil {
		contenderID, score := takeFirst()

		if contenderID == 0 {
			break
		}

		err := k.repo.StoreScore(ctx, nil, score)
		switch {
		case err == nil:
			persistedScores += 1
		case errors.Is(err, domain.ErrNotFound):
			slog.Warn("failed to persist score for non-existent contender",
				"contender_id", contenderID,
				"action", "drop",
				"error", err)

			continue
		default:
			slog.Error("failed to persist score",
				"contender_id", contenderID,
				"action", "try_again_later",
				"error", err)

			putBack(contenderID, score)

			break IterateScores
		}
	}

	leftInMemory := k.getNumScoresWithLock()
	if leftInMemory > 0 {
		slog.Warn("not all scores where persisted",
			"left_in_memory", leftInMemory,
		)
	} else if persistedScores > 0 {
		slog.Info("successfully persisted scores", "num_scores", persistedScores)
	}
}

func (k *Keeper) getNumScoresWithLock() int {
	k.mu.RLock()
	defer k.mu.RUnlock()

	return len(k.scores)
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
