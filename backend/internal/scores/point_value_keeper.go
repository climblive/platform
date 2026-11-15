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

const gcInterval = time.Minute

type keptPointValue struct {
	value     domain.PointValue
	expiresAt time.Time
}

type PointValueKeeper struct {
	mu          sync.RWMutex
	eventBroker domain.EventBroker
	pointValues map[domain.CompetitorID][]keptPointValue
	running     atomic.Bool
	ttl         time.Duration
}

func NewPointValueKeeper(eventBroker domain.EventBroker, ttl time.Duration) *PointValueKeeper {
	return &PointValueKeeper{
		mu:          sync.RWMutex{},
		eventBroker: eventBroker,
		pointValues: make(map[domain.CompetitorID][]keptPointValue),
		running:     atomic.Bool{},
		ttl:         ttl,
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

	subscriptionID, eventReader := k.eventBroker.Subscribe([]domain.EventFilter{filter}, 0)
	defer k.eventBroker.Unsubscribe(subscriptionID)

	k.running.Store(true)
	defer k.running.Store(false)

	close(ready)

	events := eventReader.EventsChan(ctx)
	ticker := time.Tick(gcInterval)

	slog.Info("point value keeper started", "ttl", k.ttl, "gc_interval", gcInterval)

EventLoop:
	for {
		select {
		case event, open := <-events:
			if !open {
				break EventLoop
			}

			switch ev := event.Data.(type) {
			case domain.PointValueUpdatedEvent:
				k.handlePointValueUpdated(ev)
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

	slog.Info("point value keeper shutting down")
}

func (k *PointValueKeeper) handlePointValueUpdated(event domain.PointValueUpdatedEvent) {
	k.mu.Lock()
	defer k.mu.Unlock()

	competitorPointValues := k.pointValues[event.CompetitorID]

	idx := slices.IndexFunc(competitorPointValues, func(kpv keptPointValue) bool {
		return kpv.value.ProblemID == event.ProblemID
	})

	updatedValue := keptPointValue{
		value:     domain.PointValue(event),
		expiresAt: time.Now().Add(k.ttl),
	}

	if idx == -1 {
		k.pointValues[event.CompetitorID] = append(competitorPointValues, updatedValue)
	} else {
		competitorPointValues[idx] = updatedValue
	}
}

func (k *PointValueKeeper) GetPointValues(competitorID domain.CompetitorID) []domain.PointValue {
	k.mu.RLock()
	defer k.mu.RUnlock()

	competitorPointValues := k.pointValues[competitorID]
	if competitorPointValues == nil {
		return []domain.PointValue{}
	}

	pointValues := make([]domain.PointValue, 0, len(competitorPointValues))

	for _, pointValue := range competitorPointValues {
		pointValues = append(pointValues, pointValue.value)
	}

	return pointValues
}

func (k *PointValueKeeper) expungeExpired(now time.Time) {
	k.mu.Lock()
	defer k.mu.Unlock()

	for competitorID, competitorPointValues := range k.pointValues {
		expiredCount := 0

		for _, pointValue := range competitorPointValues {
			if now.After(pointValue.expiresAt) {
				expiredCount += 1
			}
		}

		if expiredCount == 0 {
			continue
		}

		freshCompetitorPointValues := make([]keptPointValue, 0, len(competitorPointValues)-expiredCount)

		for _, pointValue := range competitorPointValues {
			if now.After(pointValue.expiresAt) {
				continue
			}

			freshCompetitorPointValues = append(freshCompetitorPointValues, pointValue)
		}

		k.pointValues[competitorID] = freshCompetitorPointValues
	}
}

func (k *PointValueKeeper) GetStatus() domain.ServiceStatus {
	return domain.ServiceStatus{Name: "PointValueKeeper", Healthy: k.running.Load(), CheckedAt: time.Now()}
}
