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
	HandleContenderEntered(event domain.ContenderEnteredEvent) iter.Seq[Effect]
	HandleContenderSwitchedClass(event domain.ContenderSwitchedClassEvent) iter.Seq[Effect]
	HandleContenderWithdrewFromFinals(event domain.ContenderWithdrewFromFinalsEvent) iter.Seq[Effect]
	HandleContenderReenteredFinals(event domain.ContenderReenteredFinalsEvent) iter.Seq[Effect]
	HandleContenderDisqualified(event domain.ContenderDisqualifiedEvent) iter.Seq[Effect]
	HandleContenderRequalified(event domain.ContenderRequalifiedEvent) iter.Seq[Effect]
	HandleAscentRegistered(event domain.AscentRegisteredEvent) iter.Seq[Effect]
	HandleAscentDeregistered(event domain.AscentDeregisteredEvent) iter.Seq[Effect]
	HandleProblemAdded(event domain.ProblemAddedEvent) iter.Seq[Effect]
	HandleProblemUpdated(event domain.ProblemUpdatedEvent) iter.Seq[Effect]

	RankCompClass(compClassID domain.CompClassID) iter.Seq[Effect]
	ScoreContender(contenderID domain.ContenderID) iter.Seq[Effect]
	CalculateProblemValue(compClassID domain.CompClassID, problemID domain.ProblemID) iter.Seq[Effect]

	GetDirtyScores() []domain.Score
}

type ScoreEngineDriver struct {
	logger     *slog.Logger
	contestID  domain.ContestID
	instanceID domain.ScoreEngineInstanceID

	eventBroker   domain.EventBroker
	pendingEvents []domain.EventEnvelope

	engine ScoreEngine

	running atomic.Bool

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
		"PROBLEM_UPDATED",
		"RULES_UPDATED",
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

func (d *ScoreEngineDriver) handleEvent(event domain.EventEnvelope) {
	var effects iter.Seq[Effect]

	switch ev := event.Data.(type) {
	case domain.RulesUpdatedEvent:
		effects = d.engine.HandleRulesUpdated(ev)
	case domain.ContenderEnteredEvent:
		effects = d.engine.HandleContenderEntered(ev)
	case domain.ContenderSwitchedClassEvent:
		effects = d.engine.HandleContenderSwitchedClass(ev)
	case domain.ContenderWithdrewFromFinalsEvent:
		effects = d.engine.HandleContenderWithdrewFromFinals(ev)
	case domain.ContenderReenteredFinalsEvent:
		effects = d.engine.HandleContenderReenteredFinals(ev)
	case domain.ContenderDisqualifiedEvent:
		effects = d.engine.HandleContenderDisqualified(ev)
	case domain.ContenderRequalifiedEvent:
		effects = d.engine.HandleContenderRequalified(ev)
	case domain.AscentRegisteredEvent:
		effects = d.engine.HandleAscentRegistered(ev)
	case domain.AscentDeregisteredEvent:
		effects = d.engine.HandleAscentDeregistered(ev)
	case domain.ProblemAddedEvent:
		effects = d.engine.HandleProblemAdded(ev)
	case domain.ProblemUpdatedEvent:
		effects = d.engine.HandleProblemUpdated(ev)
	}

	runner := &EffectRunner{
		queue:  make(map[EncodedEffect]Effect),
		driver: d,
	}
	runner.RunChainEffects(effects)
}

type EffectRunner struct {
	queue  map[EncodedEffect]Effect
	driver *ScoreEngineDriver
}

func (r *EffectRunner) RunChainEffects(effects iter.Seq[Effect]) {
	r.Run(func(yield func(Effect) bool) {
		for effect := range effects {
			switch effect.(type) {
			case EffectCalculateProblemValue:
			default:
				r.queue[effect.Encode()] = effect
				continue
			}

			if !yield(effect) {
				return
			}
		}
	})

	r.Run(func(yield func(Effect) bool) {
		for _, effect := range r.queue {
			switch effect.(type) {
			case EffectScoreContender:
			default:
				continue
			}

			if !yield(effect) {
				return
			}
		}
	})

	r.Run(func(yield func(Effect) bool) {
		for _, effect := range r.queue {
			switch effect.(type) {
			case EffectRankClass:
			default:
				continue
			}

			if !yield(effect) {
				return
			}
		}
	})
}

func (r *EffectRunner) Run(effects iter.Seq[Effect]) {
	for e := range effects {
		var chainEffects iter.Seq[Effect]

		switch effect := e.(type) {
		case EffectRankClass:
			r.driver.logger.Info("re-ranking comp class", "comp_class_id", effect.CompClassID)
			chainEffects = r.driver.engine.RankCompClass(effect.CompClassID)
		case EffectScoreContender:
			r.driver.logger.Info("re-scoring contender", "contender_id", effect.ContenderID)
			chainEffects = r.driver.engine.ScoreContender(effect.ContenderID)
		case EffectCalculateProblemValue:
			r.driver.logger.Info("re-calculating problem value", "comp_class_id", effect.CompClassID, "problem_id", effect.ProblemID)
			chainEffects = r.driver.engine.CalculateProblemValue(effect.CompClassID, effect.ProblemID)
		}

		if chainEffects == nil {
			continue
		}

		for chainEffect := range chainEffects {
			r.queue[chainEffect.Encode()] = chainEffect
		}
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
