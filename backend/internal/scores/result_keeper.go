package scores

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

const persistInterval = time.Minute
const lastDitchPersistTimeout = 10 * time.Second

type keeperRepository interface {
	StoreResult(ctx context.Context, tx domain.Transaction, result domain.Result) error
}

type Keeper struct {
	mu                     sync.RWMutex
	eventBroker            domain.EventBroker
	results                map[domain.CompetitorID]domain.Result
	repo                   keeperRepository
	externalPersistTrigger chan struct{}
	running                atomic.Bool
}

func NewResultKeeper(eventBroker domain.EventBroker, repo keeperRepository) *Keeper {
	return &Keeper{
		eventBroker:            eventBroker,
		results:                make(map[domain.CompetitorID]domain.Result),
		repo:                   repo,
		externalPersistTrigger: make(chan struct{}, 1),
		mu:                     sync.RWMutex{},
		running:                atomic.Bool{},
	}
}

func (k *Keeper) Run(ctx context.Context, options ...func(*runOptions)) *sync.WaitGroup {
	config := &runOptions{}
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
				slog.Error("result keeper panicked", "error", r)
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
		domain.CompetitionID{},
		0,
		"COMPETITOR_RESULT_UPDATED",
	)

	subscriptionID, eventReader := k.eventBroker.Subscribe([]domain.EventFilter{filter}, 0)
	defer k.eventBroker.Unsubscribe(subscriptionID)

	k.running.Store(true)
	defer k.running.Store(false)

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
			case domain.CompetitorResultUpdatedEvent:
				k.HandleCompetitorResultUpdated(ev)
			}
		case <-ticker:
			k.persistResults(ctx)
		case <-k.externalPersistTrigger:
			k.persistResults(ctx)
		case <-ctx.Done():
			slog.Info("subscription closed", "reason", ctx.Err())
			break EventLoop
		}
	}

	if ctx.Err() == nil {
		slog.Error("subscription closed unexpectedly")
	}

	slog.Info("result keeper shutting down")

	numResults := k.getNumResultsWithLock()
	if numResults > 0 {
		ctxWithDeadline, cancel := context.WithTimeout(context.Background(), lastDitchPersistTimeout)
		defer cancel()

		slog.Warn("making a last-ditch attempt to persist results", "timeout", lastDitchPersistTimeout)
		k.persistResults(ctxWithDeadline)
	}
}

func (k *Keeper) RequestPersist() {
	k.externalPersistTrigger <- struct{}{}
}

func (k *Keeper) persistResults(ctx context.Context) {
	takeFirst := func() (domain.CompetitorID, domain.Result) {
		k.mu.Lock()
		defer k.mu.Unlock()

		var competitorID domain.CompetitorID
		var result domain.Result
		for competitorID, result = range k.results {
			break
		}

		if competitorID == 0 {
			return 0, domain.Result{}
		}

		delete(k.results, competitorID)

		return competitorID, result
	}

	putBack := func(competitorID domain.CompetitorID, result domain.Result) {
		k.mu.Lock()
		defer k.mu.Unlock()

		if _, found := k.results[competitorID]; !found {
			k.results[competitorID] = result
		}
	}

	persistedResults := 0

IterateResults:
	for ctx.Err() == nil {
		competitorID, result := takeFirst()

		if competitorID == 0 {
			break
		}

		err := k.repo.StoreResult(ctx, nil, result)
		switch {
		case err == nil:
			persistedResults += 1
		case errors.Is(err, domain.ErrNotFound):
			slog.Warn("failed to persist result for non-existent competitor",
				"competitor_id", competitorID,
				"action", "drop",
				"error", err)

			continue
		default:
			slog.Error("failed to persist result",
				"competitor_id", competitorID,
				"action", "try_again_later",
				"error", err)

			putBack(competitorID, result)

			break IterateResults
		}
	}

	leftInMemory := k.getNumResultsWithLock()
	if leftInMemory > 0 {
		slog.Warn("not all results were persisted",
			"left_in_memory", leftInMemory,
		)
	} else if persistedResults > 0 {
		slog.Info("successfully persisted results", "num_results", persistedResults)
	}
}

func (k *Keeper) getNumResultsWithLock() int {
	k.mu.RLock()
	defer k.mu.RUnlock()

	return len(k.results)
}

func (k *Keeper) HandleCompetitorResultUpdated(event domain.CompetitorResultUpdatedEvent) {
	k.mu.Lock()
	defer k.mu.Unlock()

	k.results[event.CompetitorID] = domain.Result(event)
}

func (k *Keeper) GetResult(competitorID domain.CompetitorID) (domain.Result, error) {
	k.mu.RLock()
	defer k.mu.RUnlock()

	if result, found := k.results[competitorID]; found {
		return result, nil
	}

	return domain.Result{}, domain.ErrNotFound
}

func (k *Keeper) GetStatus() domain.ServiceStatus {
	return domain.ServiceStatus{Name: "ResultKeeper", Healthy: k.running.Load(), CheckedAt: time.Now()}
}
