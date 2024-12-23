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

type ScoringRules interface {
	CalculateScore(points iter.Seq[int]) int
}

type Ranker interface {
	RankContenders(contenders iter.Seq[Contender]) []domain.Score
}

type EngineStore interface {
	GetContender(domain.ContenderID) (Contender, bool)
	GetContendersByCompClass(domain.CompClassID) iter.Seq[Contender]
	GetAllContenders() iter.Seq[Contender]
	SaveContender(Contender)

	GetCompClassIDs() []domain.CompClassID

	GetTicks(domain.ContenderID) iter.Seq[Tick]
	SaveTick(domain.ContenderID, Tick)
	DeleteTick(domain.ContenderID, domain.ProblemID)

	GetProblem(domain.ProblemID) (Problem, bool)
	SaveProblem(Problem)

	SaveScore(domain.Score)
	GetUnpublishedScores() []domain.Score
}

type ScoreEngine struct {
	logger      *slog.Logger
	contestID   domain.ContestID
	eventBroker domain.EventBroker

	ranker Ranker
	rules  ScoringRules
	store  EngineStore

	running    atomic.Bool
	sideQuests chan func()
}

func NewScoreEngine(
	contestID domain.ContestID,
	eventBroker domain.EventBroker,
	rules ScoringRules,
	ranker Ranker,
	store EngineStore,
) *ScoreEngine {
	logger := slog.New(slog.Default().Handler()).With("contest_id", contestID)

	return &ScoreEngine{
		logger:      logger,
		contestID:   contestID,
		eventBroker: eventBroker,
		ranker:      ranker,
		rules:       rules,
		store:       store,
		sideQuests:  make(chan func(), 1),
	}
}

func (e *ScoreEngine) Run(ctx context.Context) *sync.WaitGroup {
	wg := new(sync.WaitGroup)
	ready := make(chan struct{}, 1)

	wg.Add(1)

	filter := domain.NewEventFilter(
		e.contestID,
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
	)

	go e.run(ctx, filter, wg, ready)

	<-ready

	return wg
}

func (e *ScoreEngine) run(ctx context.Context, filter domain.EventFilter, wg *sync.WaitGroup, ready chan<- struct{}) {
	defer func() {
		if r := recover(); r != nil {
			e.logger.Error("score engine panicked", "error", r)
		}
	}()

	defer wg.Done()

	defer func() {
		close(e.sideQuests)

		for range e.sideQuests {
		}
	}()

	e.running.Store(true)
	defer e.running.Store(false)

	subscriptionID, eventReader := e.eventBroker.Subscribe(filter, 0)
	e.logger.Info("score engine subscribed", "subscription_id", subscriptionID)

	defer e.eventBroker.Unsubscribe(subscriptionID)

	close(ready)

	events := eventReader.EventsChan(ctx)
	ticker := time.Tick(100 * time.Millisecond)

MainLoop:
	for {
		select {
		case event, open := <-events:
			if !open {
				break MainLoop
			}

			e.handleEvent(event)
		case <-ticker:
			e.publishUpdatedScores()
		case f := <-e.sideQuests:
			f()
		case <-ctx.Done():
			break MainLoop
		}
	}

	if ctx.Err() == nil {
		e.logger.Warn("subscription closed unexpectedly")
	}

	e.logger.Info("score engine shutting down")
}

func (e *ScoreEngine) SetScoringRules(rules ScoringRules) {
	quest := func() {
		e.rules = rules

		for contender := range e.store.GetAllContenders() {
			contender.Score = e.rules.CalculateScore(Points(e.store.GetTicks(contender.ID)))
			e.store.SaveContender(contender)
		}

		e.rankCompClasses(e.store.GetCompClassIDs()...)
	}

	if e.running.Load() {
		e.sideQuests <- quest
	}
}

func (e *ScoreEngine) SetRanker(ranker Ranker) {
	quest := func() {
		e.ranker = ranker

		e.rankCompClasses(e.store.GetCompClassIDs()...)
	}

	if e.running.Load() {
		e.sideQuests <- quest
	}
}

func (e *ScoreEngine) handleEvent(event domain.EventEnvelope) {
	switch ev := event.Data.(type) {
	case domain.ContenderEnteredEvent:
		e.handleContenderEntered(ev)
	case domain.ContenderSwitchedClassEvent:
		e.handleContenderSwitchedClass(ev)
	case domain.ContenderWithdrewFromFinalsEvent:
		e.handleContenderWithdrewFromFinals(ev)
	case domain.ContenderReenteredFinalsEvent:
		e.handleContenderReenteredFinals(ev)
	case domain.ContenderDisqualifiedEvent:
		e.handleContenderDisqualified(ev)
	case domain.ContenderRequalifiedEvent:
		e.handleContenderRequalified(ev)
	case domain.AscentRegisteredEvent:
		e.handleAscentRegistered(ev)
	case domain.AscentDeregisteredEvent:
		e.handleAscentDeregistered(ev)
	case domain.ProblemAddedEvent:
		e.handleProblemAdded(ev)
	case domain.ProblemUpdatedEvent:
		e.logger.Warn("discarding unsupported event", "event", event)
	}
}

func (e *ScoreEngine) handleContenderEntered(event domain.ContenderEnteredEvent) {
	contender := Contender{
		ID:          event.ContenderID,
		CompClassID: event.CompClassID,
	}

	e.store.SaveContender(contender)

	e.rankCompClasses(contender.CompClassID)
}

func (e *ScoreEngine) handleContenderSwitchedClass(event domain.ContenderSwitchedClassEvent) {
	contender, found := e.store.GetContender(event.ContenderID)
	if !found {
		return
	}

	if contender.CompClassID == event.CompClassID {
		return
	}

	compClassesToReRank := []domain.CompClassID{
		contender.CompClassID,
		event.CompClassID,
	}

	contender.CompClassID = event.CompClassID

	e.store.SaveContender(contender)

	e.rankCompClasses(compClassesToReRank...)
}

func (e *ScoreEngine) handleContenderWithdrewFromFinals(event domain.ContenderWithdrewFromFinalsEvent) {
	contender, found := e.store.GetContender(event.ContenderID)
	if !found {
		return
	}

	contender.WithdrawnFromFinals = true

	e.store.SaveContender(contender)

	e.rankCompClasses(contender.CompClassID)
}

func (e *ScoreEngine) handleContenderReenteredFinals(event domain.ContenderReenteredFinalsEvent) {
	contender, found := e.store.GetContender(event.ContenderID)
	if !found {
		return
	}

	contender.WithdrawnFromFinals = false

	e.store.SaveContender(contender)

	e.rankCompClasses(contender.CompClassID)
}

func (e *ScoreEngine) handleContenderDisqualified(event domain.ContenderDisqualifiedEvent) {
	contender, found := e.store.GetContender(event.ContenderID)
	if !found {
		return
	}

	contender.Disqualified = true
	contender.Score = 0

	e.store.SaveContender(contender)

	e.rankCompClasses(contender.CompClassID)
}

func (e *ScoreEngine) handleContenderRequalified(event domain.ContenderRequalifiedEvent) {
	contender, found := e.store.GetContender(event.ContenderID)
	if !found {
		return
	}

	contender.Disqualified = false
	contender.Score = e.rules.CalculateScore(Points(e.store.GetTicks(contender.ID)))

	e.store.SaveContender(contender)

	e.rankCompClasses(contender.CompClassID)
}

func (e *ScoreEngine) handleAscentRegistered(event domain.AscentRegisteredEvent) {
	tick := Tick{
		ProblemID:    event.ProblemID,
		Top:          event.Top,
		AttemptsTop:  event.AttemptsTop,
		Zone:         event.Zone,
		AttemptsZone: event.AttemptsZone,
	}

	contender, found := e.store.GetContender(event.ContenderID)
	if !found {
		return
	}

	problem, found := e.store.GetProblem(event.ProblemID)
	if !found {
		return
	}

	tick.Score(problem)
	e.store.SaveTick(event.ContenderID, tick)

	if contender.Disqualified {
		return
	}

	contender.Score = e.rules.CalculateScore(Points(e.store.GetTicks(contender.ID)))
	e.store.SaveContender(contender)

	e.rankCompClasses(contender.CompClassID)
}

func (e *ScoreEngine) handleAscentDeregistered(event domain.AscentDeregisteredEvent) {
	contender, found := e.store.GetContender(event.ContenderID)
	if !found {
		return
	}

	e.store.DeleteTick(event.ContenderID, event.ProblemID)

	if !contender.Disqualified {
		contender.Score = e.rules.CalculateScore(Points(e.store.GetTicks(contender.ID)))
		e.store.SaveContender(contender)
	}

	e.rankCompClasses(contender.CompClassID)
}

func (e *ScoreEngine) handleProblemAdded(event domain.ProblemAddedEvent) {
	problem := Problem{
		ID:         event.ProblemID,
		PointsTop:  event.PointsTop,
		PointsZone: event.PointsZone,
		FlashBonus: event.FlashBonus,
	}

	e.store.SaveProblem(problem)
}

func (e *ScoreEngine) GetUpdatedScores() []domain.Score {
	return e.store.GetUnpublishedScores()
}

func (e *ScoreEngine) rankCompClasses(compClassIDs ...domain.CompClassID) {
	for _, compClassID := range compClassIDs {
		scores := e.ranker.RankContenders(e.store.GetContendersByCompClass(compClassID))

		for score := range slices.Values(scores) {
			e.store.SaveScore(score)
		}
	}
}

func (e *ScoreEngine) publishUpdatedScores() {
	scores := e.store.GetUnpublishedScores()

	var batch []domain.ContenderScoreUpdatedEvent

	for score := range slices.Values(scores) {
		e.eventBroker.Dispatch(e.contestID, domain.ContenderScoreUpdatedEvent(score))

		batch = append(batch, domain.ContenderScoreUpdatedEvent(score))
	}

	if len(batch) > 0 {
		e.eventBroker.Dispatch(e.contestID, batch)
	}
}
