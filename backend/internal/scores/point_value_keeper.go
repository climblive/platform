package scores

import (
	"context"
	"log/slog"
	"slices"
	"sync"
	"sync/atomic"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
)

const pointValueTTL = 24 * time.Hour
const pointValueSweepInterval = time.Minute

type keptPointValue struct {
	value     domain.PointValue
	expiresAt time.Time
}

type PointValueKeeper struct {
	mu          sync.RWMutex
	eventBroker domain.EventBroker
	pointValues map[domain.ContenderID]map[domain.ProblemID]keptPointValue
	running     atomic.Bool
}

func NewPointValueKeeper(eventBroker domain.EventBroker) *PointValueKeeper {
	return &PointValueKeeper{
		mu:          sync.RWMutex{},
		eventBroker: eventBroker,
		pointValues: make(map[domain.ContenderID]map[domain.ProblemID]keptPointValue),
		running:     atomic.Bool{},
	}
}

func (k *PointValueKeeper) Run(ctx context.Context, options ...func(*runOptions)) *sync.WaitGroup {
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
				slog.Error("point value keeper panicked", "error", r)
			}
		}()

		defer wg.Done()

		k.run(ctx, ready)
	}()

	<-ready

	return wg
}

func (k *PointValueKeeper) run(ctx context.Context, ready chan<- struct{}) {
	filter := domain.NewEventFilter(
		0,
		0,
		"POINT_VALUE_UPDATED",
	)

	subscriptionID, eventReader := k.eventBroker.Subscribe(filter, 0)
	defer k.eventBroker.Unsubscribe(subscriptionID)

	k.running.Store(true)
	defer k.running.Store(false)

	close(ready)

	events := eventReader.EventsChan(ctx)
	ticker := time.Tick(pointValueSweepInterval)

EventLoop:
	for {
		select {
		case event, open := <-events:
			if !open {
				break EventLoop
			}

			switch ev := event.Data.(type) {
			case domain.PointValueUpdatedEvent:
				k.HandlePointValueUpdated(ev)
			}
		case <-ticker:
			k.expungeExpired(time.Now())
		case <-ctx.Done():
			slog.Info("point value keeper subscription closed", "reason", ctx.Err())
			break EventLoop
		}
	}

	if ctx.Err() == nil {
		slog.Error("point value keeper subscription closed unexpectedly")
	}
}

func (k *PointValueKeeper) HandlePointValueUpdated(event domain.PointValueUpdatedEvent) {
	k.mu.Lock()
	defer k.mu.Unlock()

	contenderPointValues := k.pointValues[event.ContenderID]
	if contenderPointValues == nil {
		contenderPointValues = make(map[domain.ProblemID]keptPointValue)
		k.pointValues[event.ContenderID] = contenderPointValues
	}

	contenderPointValues[event.ProblemID] = keptPointValue{
		value:     domain.PointValue(event),
		expiresAt: time.Now().Add(pointValueTTL),
	}
}

func (k *PointValueKeeper) GetPointValues(contenderID domain.ContenderID) []domain.PointValue {
	k.mu.RLock()
	defer k.mu.RUnlock()

	now := time.Now()
	contenderPointValues := k.pointValues[contenderID]
	if contenderPointValues == nil {
		return []domain.PointValue{}
	}

	pointValues := make([]domain.PointValue, 0, len(contenderPointValues))

	for _, pointValue := range contenderPointValues {
		if now.After(pointValue.expiresAt) {
			continue
		}

		pointValues = append(pointValues, pointValue.value)
	}

	slices.SortFunc(pointValues, func(a, b domain.PointValue) int {
		if a.ProblemID == b.ProblemID {
			return int(a.ContenderID) - int(b.ContenderID)
		}

		return int(a.ProblemID) - int(b.ProblemID)
	})

	return pointValues
}

func (k *PointValueKeeper) expungeExpired(now time.Time) {
	k.mu.Lock()
	defer k.mu.Unlock()

	for contenderID, contenderPointValues := range k.pointValues {
		for problemID, pointValue := range contenderPointValues {
			if now.After(pointValue.expiresAt) {
				delete(contenderPointValues, problemID)
			}
		}

		if len(contenderPointValues) == 0 {
			delete(k.pointValues, contenderID)
		}
	}
}
