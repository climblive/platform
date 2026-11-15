package scores

import (
	"context"
	"iter"
	"log/slog"
	"slices"
	"sync"
	"sync/atomic"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
)

type ScoreEngine interface {
	Start() iter.Seq[Effect]
	Stop()

	HandleRulesUpdated(event domain.RulesUpdatedEvent) iter.Seq[Effect]
	HandleCompetitorEntered(event domain.CompetitorEnteredEvent) iter.Seq[Effect]
	HandleCompetitorSwitchedCategory(event domain.CompetitorSwitchedCategoryEvent) iter.Seq[Effect]
	HandleCompetitorWithdrewFromFinals(event domain.CompetitorWithdrewFromFinalsEvent) iter.Seq[Effect]
	HandleCompetitorReenteredFinals(event domain.CompetitorReenteredFinalsEvent) iter.Seq[Effect]
	HandleCompetitorDisqualified(event domain.CompetitorDisqualifiedEvent) iter.Seq[Effect]
	HandleCompetitorRequalified(event domain.CompetitorRequalifiedEvent) iter.Seq[Effect]
	HandleAscentRegistered(event domain.AscentRegisteredEvent) iter.Seq[Effect]
	HandleAscentDeregistered(event domain.AscentDeregisteredEvent) iter.Seq[Effect]
	HandleProblemAdded(event domain.ProblemAddedEvent) iter.Seq[Effect]
	HandleProblemUpdated(event domain.ProblemUpdatedEvent) iter.Seq[Effect]

	RankCategory(categoryID domain.CategoryID)
	ScoreCompetitor(competitorID domain.CompetitorID) iter.Seq[Effect]
	CalculatePointValues(categoryID domain.CategoryID, problemID domain.ProblemID) iter.Seq[Effect]

	GetDirtyResults() []domain.Result
	GetDirtyPointValues() []domain.PointValue
}

type ScoreEngineDriver struct {
	logger        *slog.Logger
	competitionID domain.CompetitionID
	instanceID    domain.ScoreEngineInstanceID

	eventBroker   domain.EventBroker
	pendingEvents []domain.EventEnvelope

	engine ScoreEngine

	running atomic.Bool

	publishToken bool
}

func NewScoreEngineDriver(
	competitionID domain.CompetitionID,
	instanceID domain.ScoreEngineInstanceID,
	eventBroker domain.EventBroker,
) *ScoreEngineDriver {
	logger := slog.New(slog.Default().Handler()).
		With("competition_id", competitionID).
		With("instance_id", instanceID)

	return &ScoreEngineDriver{
		logger:        logger,
		competitionID: competitionID,
		instanceID:    instanceID,
		eventBroker:   eventBroker,
		pendingEvents: make([]domain.EventEnvelope, 0),
		engine:        nil,
		running:       atomic.Bool{},
		publishToken:  false,
	}
}

type runOptions struct {
	recoverPanics bool
}

func WithPanicRecovery() func(*runOptions) {
	return func(s *runOptions) {
		s.recoverPanics = true
	}
}

func (d *ScoreEngineDriver) Run(ctx context.Context, options ...func(*runOptions)) (*sync.WaitGroup, func(ScoreEngine)) {
	config := &runOptions{}
	for _, opt := range options {
		opt(config)
	}

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
			if !config.recoverPanics {
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
	filter := domain.NewEventFilter(
		d.competitionID,
		0,
		"COMPETITOR_ENTERED",
		"COMPETITOR_SWITCHED_CATEGORY",
		"COMPETITOR_WITHDREW_FROM_FINALS",
		"COMPETITOR_REENTERED_FINALS",
		"COMPETITOR_DISQUALIFIED",
		"COMPETITOR_REQUALIFIED",
		"ASCENT_REGISTERED",
		"ASCENT_DEREGISTERED",
		"PROBLEM_ADDED",
		"PROBLEM_UPDATED",
		"RULES_UPDATED",
	)

	subscriptionID, eventReader := d.eventBroker.Subscribe([]domain.EventFilter{filter}, 0)
	d.logger.Info("score engine subscribed", "subscription_id", subscriptionID)

	defer d.eventBroker.Unsubscribe(subscriptionID)

	close(ready)

	d.eventBroker.Dispatch(d.competitionID, domain.ScoreEngineStartedEvent{
		InstanceID: d.instanceID,
	})

	defer d.eventBroker.Dispatch(d.competitionID, domain.ScoreEngineStoppedEvent{
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

			effects := d.engine.Start()
			defer d.engine.Stop()

			runner := NewEffectRunner(d.engine)
			runner.RunEffects(effects)

			break PreLoop
		case <-ctx.Done():
			return
		}
	}

	defer d.publishDirtyUpdates()

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

			n := d.publishDirtyUpdates()
			if n == 0 {
				d.publishToken = true
			}
		case <-ctx.Done():
			return
		}

		if d.publishToken {
			n := d.publishDirtyUpdates()
			if n > 0 {
				d.publishToken = false
			}
		}
	}
}

func (d *ScoreEngineDriver) handleEvent(event domain.EventEnvelope) {
	var effects iter.Seq[Effect]

	switch ev := event.Data.(type) {
	case domain.RulesUpdatedEvent:
		effects = d.engine.HandleRulesUpdated(ev)
	case domain.CompetitorEnteredEvent:
		effects = d.engine.HandleCompetitorEntered(ev)
	case domain.CompetitorSwitchedCategoryEvent:
		effects = d.engine.HandleCompetitorSwitchedCategory(ev)
	case domain.CompetitorWithdrewFromFinalsEvent:
		effects = d.engine.HandleCompetitorWithdrewFromFinals(ev)
	case domain.CompetitorReenteredFinalsEvent:
		effects = d.engine.HandleCompetitorReenteredFinals(ev)
	case domain.CompetitorDisqualifiedEvent:
		effects = d.engine.HandleCompetitorDisqualified(ev)
	case domain.CompetitorRequalifiedEvent:
		effects = d.engine.HandleCompetitorRequalified(ev)
	case domain.AscentRegisteredEvent:
		effects = d.engine.HandleAscentRegistered(ev)
	case domain.AscentDeregisteredEvent:
		effects = d.engine.HandleAscentDeregistered(ev)
	case domain.ProblemAddedEvent:
		effects = d.engine.HandleProblemAdded(ev)
	case domain.ProblemUpdatedEvent:
		effects = d.engine.HandleProblemUpdated(ev)
	}

	runner := NewEffectRunner(d.engine)
	runner.RunEffects(effects)
}

func (d *ScoreEngineDriver) publishDirtyUpdates() int {
	publishCount := 0

	var batch []domain.CompetitorResultUpdatedEvent

	for result := range slices.Values(d.engine.GetDirtyResults()) {
		d.eventBroker.Dispatch(d.competitionID, domain.CompetitorResultUpdatedEvent(result))
		publishCount++

		batch = append(batch, domain.CompetitorResultUpdatedEvent(result))
	}

	for pointValue := range slices.Values(d.engine.GetDirtyPointValues()) {
		d.eventBroker.Dispatch(d.competitionID, domain.PointValueUpdatedEvent(pointValue))

		publishCount++
	}

	if len(batch) > 0 {
		d.eventBroker.Dispatch(d.competitionID, batch)
	}

	return publishCount
}
