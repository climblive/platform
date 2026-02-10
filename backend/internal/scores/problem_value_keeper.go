package scores

import (
	"context"
	"log/slog"
	"sync"

	"github.com/climblive/platform/backend/internal/domain"
)

type problemValueKey struct {
	ProblemID   domain.ProblemID
	CompClassID domain.CompClassID
}

type ProblemValueKeeper struct {
	mu            sync.RWMutex
	eventBroker   domain.EventBroker
	problemValues map[problemValueKey]domain.ProblemValue
}

func NewProblemValueKeeper(eventBroker domain.EventBroker) *ProblemValueKeeper {
	return &ProblemValueKeeper{
		eventBroker:   eventBroker,
		problemValues: make(map[problemValueKey]domain.ProblemValue),
		mu:            sync.RWMutex{},
	}
}

func (k *ProblemValueKeeper) Run(ctx context.Context, options ...func(*runOptions)) *sync.WaitGroup {
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
				slog.Error("problem value keeper panicked", "error", r)
			}
		}()

		defer wg.Done()

		k.run(ctx, ready)
	}()

	<-ready

	return wg
}

func (k *ProblemValueKeeper) run(ctx context.Context, ready chan<- struct{}) {
	filter := domain.NewEventFilter(
		0,
		0,
		"PROBLEM_VALUE_UPDATED",
	)

	subscriptionID, eventReader := k.eventBroker.Subscribe(filter, 0)
	defer k.eventBroker.Unsubscribe(subscriptionID)

	close(ready)

	events := eventReader.EventsChan(ctx)

EventLoop:
	for {
		select {
		case event, open := <-events:
			if !open {
				break EventLoop
			}

			switch ev := event.Data.(type) {
			case domain.ProblemValueUpdatedEvent:
				k.HandleProblemValueUpdated(ev)
			}
		case <-ctx.Done():
			slog.Info("problem value keeper subscription closed", "reason", ctx.Err())
			break EventLoop
		}
	}

	if ctx.Err() == nil {
		slog.Error("problem value keeper subscription closed unexpectedly")
	}

	slog.Info("problem value keeper shutting down")
}

func (k *ProblemValueKeeper) HandleProblemValueUpdated(event domain.ProblemValueUpdatedEvent) {
	k.mu.Lock()
	defer k.mu.Unlock()

	key := problemValueKey{
		ProblemID:   event.ProblemID,
		CompClassID: event.CompClassID,
	}

	k.problemValues[key] = event.ProblemValue
}

func (k *ProblemValueKeeper) GetProblemValue(problemID domain.ProblemID, compClassID domain.CompClassID) (domain.ProblemValue, bool) {
	k.mu.RLock()
	defer k.mu.RUnlock()

	key := problemValueKey{
		ProblemID:   problemID,
		CompClassID: compClassID,
	}

	value, found := k.problemValues[key]
	return value, found
}
