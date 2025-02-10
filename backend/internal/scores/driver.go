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
	logger     *slog.Logger
	contestID  domain.ContestID
	instanceID domain.ScoreEngineInstanceID

	eventBroker   domain.EventBroker
	pendingEvents []domain.EventEnvelope

	engine ScoreEngine

	running    atomic.Bool
	sideQuests chan func()

	publishToken bool
}

func NewScoreEngineDriver(
	contestID domain.ContestID,
	instanceID domain.ScoreEngineInstanceID,
	eventBroker domain.EventBroker,
) *ScoreEngineDriver {
	logger := slog.New(slog.Default().Handler()).
		With("contest_id", contestID).
		With("instance_id", instanceID)

	return &ScoreEngineDriver{
		logger:        logger,
		contestID:     contestID,
		instanceID:    instanceID,
		eventBroker:   eventBroker,
		pendingEvents: make([]domain.EventEnvelope, 0),
		sideQuests:    make(chan func()),
	}
}

func (d *ScoreEngineDriver) Run(ctx context.Context, safe bool) (*sync.WaitGroup, func(ScoreEngine)) {
	wg := new(sync.WaitGroup)
	ready := make(chan struct{}, 1)

	wg.Add(1)

	engineReceiver := make(chan ScoreEngine, 1)

	installEngine := func(engine ScoreEngine) {
		engineReceiver <- engine
		close(engineReceiver)
	}

	go func() {
		defer func() {
			if !safe {
				return
			}

			if r := recover(); r != nil {
				d.logger.Error("score engine panicked", "error", r)
			}
		}()

		defer wg.Done()

		d.run(ctx, ready, engineReceiver)
	}()

	<-ready

	return wg, installEngine
}

func (d *ScoreEngineDriver) run(
	ctx context.Context,
	ready chan<- struct{},
	engineReceiver chan ScoreEngine,
) {
	defer func() {
		close(d.sideQuests)

		for range d.sideQuests {
		}
	}()

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

	subscriptionID, eventReader := d.eventBroker.Subscribe(filter, 0)
	d.logger.Info("score engine subscribed", "subscription_id", subscriptionID)

	defer d.eventBroker.Unsubscribe(subscriptionID)

	close(ready)

	d.eventBroker.Dispatch(d.contestID, domain.ScoreEngineStartedEvent{
		InstanceID: d.instanceID,
	})

	defer d.eventBroker.Dispatch(d.contestID, domain.ScoreEngineStoppedEvent{
		InstanceID: d.instanceID,
	})

	events := eventReader.EventsChan(ctx)

	d.processEvents(ctx, events, engineReceiver)

	if ctx.Err() == nil {
		d.logger.Warn("subscription closed unexpectedly")
	}

	d.logger.Info("score engine shutting down")
}

func (d *ScoreEngineDriver) processEvents(
	ctx context.Context,
	events <-chan domain.EventEnvelope,
	engineReceiver chan ScoreEngine,
) {
PreLoop:
	for {
		select {
		case event, open := <-events:
			if !open {
				return
			}

			d.pendingEvents = append(d.pendingEvents, event)
		case engine := <-engineReceiver:
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
	d.pendingEvents = nil

	ticker := time.Tick(100 * time.Millisecond)

	for {
		select {
		case event, open := <-events:
			if !open {
				return
			}

			d.handleEvent(event)
		case <-ticker:
			d.publishToken = false

			n := d.publishUpdatedScores()
			if n == 0 {
				d.publishToken = true
			}
		case f := <-d.sideQuests:
			f()
		case <-ctx.Done():
			return
		}

		if d.publishToken {
			n := d.publishUpdatedScores()
			if n > 0 {
				d.publishToken = false
			}
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

func (d *ScoreEngineDriver) publishUpdatedScores() int {
	scores := d.engine.GetDirtyScores()

	var batch []domain.ContenderScoreUpdatedEvent

	for score := range slices.Values(scores) {
		d.eventBroker.Dispatch(d.contestID, domain.ContenderScoreUpdatedEvent(score))

		batch = append(batch, domain.ContenderScoreUpdatedEvent(score))
	}

	if len(batch) > 0 {
		d.eventBroker.Dispatch(d.contestID, batch)
	}

	return len(scores)
}
