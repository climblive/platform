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

type ScoreEngine interface {
	Start()
	Stop()

	ReplaceScoringRules(rules ScoringRules)
	ReplaceRanker(ranker Ranker)

	HandleContenderEntered(event domain.ContenderEnteredEvent)
	HandleContenderSwitchedClass(event domain.ContenderSwitchedClassEvent)
	HandleContenderWithdrewFromFinals(event domain.ContenderWithdrewFromFinalsEvent)
	HandleContenderReenteredFinals(event domain.ContenderReenteredFinalsEvent)
	HandleContenderDisqualified(event domain.ContenderDisqualifiedEvent)
	HandleContenderRequalified(event domain.ContenderRequalifiedEvent)
	HandleAscentRegistered(event domain.AscentRegisteredEvent)
	HandleAscentDeregistered(event domain.AscentDeregisteredEvent)
	HandleProblemAdded(event domain.ProblemAddedEvent)

	GetDirtyScores() []domain.Score
}

type ScoreEngineDriver struct {
	logger    *slog.Logger
	contestID domain.ContestID

	eventBroker   domain.EventBroker
	pendingEvents []domain.EventEnvelope

	engine ScoreEngine

	running    atomic.Bool
	sideQuests chan func()
}

func NewScoreEngineDriver(
	contestID domain.ContestID,
	eventBroker domain.EventBroker,
) *ScoreEngineDriver {
	logger := slog.New(slog.Default().Handler()).With("contest_id", contestID)

	return &ScoreEngineDriver{
		logger:        logger,
		contestID:     contestID,
		eventBroker:   eventBroker,
		pendingEvents: make([]domain.EventEnvelope, 0),
		sideQuests:    make(chan func()),
	}
}

func (d *ScoreEngineDriver) Run(ctx context.Context) (*sync.WaitGroup, func(ScoreEngine)) {
	wg := new(sync.WaitGroup)
	ready := make(chan struct{}, 1)

	wg.Add(1)

	filter := domain.NewEventFilter(
		d.contestID,
		0,
		"CONTENDER_ENTERED",
		"CONTENDER_SWITCHED_CLASS",
		"CONTENDER_WITHDREW_FROM_FINALS",
		"CONTENDER_REENTERED_FINALS",
		"CONTENDER_DISQUALIFIED",
		"CONTENDER_REQUALIFIED",
		"ASCENT_REGISTERED",
		"ASCENT_DEREGISTERED",
		"PROBLEM_ADDED",
	)

	hose := make(chan ScoreEngine, 1)

	installEngine := func(engine ScoreEngine) {
		hose <- engine
		close(hose)
	}

	go d.run(ctx, filter, wg, ready, hose)

	<-ready

	return wg, installEngine
}

func (d *ScoreEngineDriver) run(ctx context.Context, filter domain.EventFilter, wg *sync.WaitGroup, ready chan<- struct{}, hose chan ScoreEngine) {
	defer func() {
		if r := recover(); r != nil {
			d.logger.Error("score engine panicked", "error", r)
		}
	}()

	defer wg.Done()

	defer func() {
		close(d.sideQuests)

		for range d.sideQuests {
		}
	}()

	subscriptionID, eventReader := d.eventBroker.Subscribe(filter, 0)
	d.logger.Info("score engine subscribed", "subscription_id", subscriptionID)

	defer d.eventBroker.Unsubscribe(subscriptionID)

	close(ready)

	events := eventReader.EventsChan(ctx)

	d.processEvents(ctx, events, hose)

	if ctx.Err() == nil {
		d.logger.Warn("subscription closed unexpectedly")
	}

	d.logger.Info("score engine shutting down")
}

func (d *ScoreEngineDriver) processEvents(ctx context.Context, events <-chan domain.EventEnvelope, hose chan ScoreEngine) {
PreLoop:
	for {
		select {
		case event, open := <-events:
			if !open {
				return
			}

			d.pendingEvents = append(d.pendingEvents, event)
		case engine := <-hose:
			d.engine = engine

			d.engine.Start()
			defer d.engine.Stop()

			break PreLoop
		case <-ctx.Done():
			return
		}
	}

	defer d.publishUpdatedScores()

	d.running.Store(true)
	defer d.running.Store(false)

	if len(d.pendingEvents) != 0 {
		d.logger.Info("replaying pending events", "count", len(d.pendingEvents))
	}

	for event := range slices.Values(d.pendingEvents) {
		d.handleEvent(event)
	}

	ticker := time.Tick(100 * time.Millisecond)

MainLoop:
	for {
		select {
		case event, open := <-events:
			if !open {
				return
			}

			d.handleEvent(event)
		case <-ticker:
			d.publishUpdatedScores()
		case f := <-d.sideQuests:
			f()
		case <-ctx.Done():
			break MainLoop
		}
	}
}

func (d *ScoreEngineDriver) SetScoringRules(rules ScoringRules) {
	quest := func() {
		d.engine.ReplaceScoringRules(rules)
	}

	if d.running.Load() {
		d.sideQuests <- quest
	}
}

func (d *ScoreEngineDriver) SetRanker(ranker Ranker) {
	quest := func() {
		d.engine.ReplaceRanker(ranker)
	}

	if d.running.Load() {
		d.sideQuests <- quest
	}
}

func (d *ScoreEngineDriver) handleEvent(event domain.EventEnvelope) {
	switch ev := event.Data.(type) {
	case domain.ContenderEnteredEvent:
		d.engine.HandleContenderEntered(ev)
	case domain.ContenderSwitchedClassEvent:
		d.engine.HandleContenderSwitchedClass(ev)
	case domain.ContenderWithdrewFromFinalsEvent:
		d.engine.HandleContenderWithdrewFromFinals(ev)
	case domain.ContenderReenteredFinalsEvent:
		d.engine.HandleContenderReenteredFinals(ev)
	case domain.ContenderDisqualifiedEvent:
		d.engine.HandleContenderDisqualified(ev)
	case domain.ContenderRequalifiedEvent:
		d.engine.HandleContenderRequalified(ev)
	case domain.AscentRegisteredEvent:
		d.engine.HandleAscentRegistered(ev)
	case domain.AscentDeregisteredEvent:
		d.engine.HandleAscentDeregistered(ev)
	case domain.ProblemAddedEvent:
		d.engine.HandleProblemAdded(ev)
	}
}

func (d *ScoreEngineDriver) publishUpdatedScores() {
	scores := d.engine.GetDirtyScores()

	var batch []domain.ContenderScoreUpdatedEvent

	for score := range slices.Values(scores) {
		d.eventBroker.Dispatch(d.contestID, domain.ContenderScoreUpdatedEvent(score))

		batch = append(batch, domain.ContenderScoreUpdatedEvent(score))
	}

	if len(batch) > 0 {
		d.eventBroker.Dispatch(d.contestID, batch)
	}
}
