package scores

import (
	"context"
	"iter"
	"log/slog"
	"maps"
	"slices"
	"sync"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
)

type ScoringRules interface {
	CalculateScore(tickPointValues iter.Seq[int]) int
}

type Ranker interface {
	RankContenders(contenders iter.Seq[*Contender]) []domain.Score
}

type ScoreEngine struct {
	logger      *slog.Logger
	contestID   domain.ContestID
	eventBroker domain.EventBroker

	sync.Mutex
	ranker Ranker
	rules  ScoringRules

	problems   map[domain.ProblemID]*Problem
	contenders map[domain.ContenderID]*Contender
	scores     *DiffMap[domain.ContenderID, domain.Score]
}

func NewScoreEngine(contestID domain.ContestID, eventBroker domain.EventBroker, rules ScoringRules, ranker Ranker) *ScoreEngine {
	logger := slog.New(slog.Default().Handler()).With("contest_id", contestID)

	return &ScoreEngine{
		logger:      logger,
		contestID:   contestID,
		eventBroker: eventBroker,
		ranker:      ranker,
		rules:       rules,
		problems:    make(map[domain.ProblemID]*Problem),
		contenders:  make(map[domain.ContenderID]*Contender),
		scores:      NewDiffMap[domain.ContenderID, domain.Score](CompareScore),
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
		"PROBLEM_DELETED",
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

	subscriptionID, eventReader := e.eventBroker.Subscribe(filter, 0)
	e.logger.Info("score engine subscribed", "subscription_id", subscriptionID)

	defer e.eventBroker.Unsubscribe(subscriptionID)

	close(ready)

	events := eventReader.EventsChan(ctx)
	ticker := time.Tick(100 * time.Millisecond)

ConsumeEvents:
	for {
		select {
		case event, open := <-events:
			if !open && ctx.Err() == nil {
				e.logger.Warn("subscription closed unexpectedly")
				break ConsumeEvents
			}

			e.HandleEvent(event)
		case <-ticker:
			e.publishUpdatedScores()
		case <-ctx.Done():
			break ConsumeEvents
		}
	}

	e.logger.Info("score engine shutting down")
}

func (e *ScoreEngine) SetScoringRules(rules ScoringRules) {
	e.Lock()
	defer e.Unlock()

	e.rules = rules

	for contender := range maps.Values(e.contenders) {
		e.scoreContender(contender)
	}

	e.rankCompClasses(CompClasses(e.contenders))
}

func (e *ScoreEngine) SetRanker(ranker Ranker) {
	e.Lock()
	defer e.Unlock()

	e.ranker = ranker

	e.rankCompClasses(CompClasses(e.contenders))
}

func (e *ScoreEngine) HandleEvent(event domain.EventEnvelope) {
	e.Lock()
	defer e.Unlock()

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
	case domain.ProblemUpdatedEvent, domain.ProblemDeletedEvent:
		e.logger.Warn("discarding unsupported event", "event", event)
	}
}

func (e *ScoreEngine) handleContenderEntered(event domain.ContenderEnteredEvent) {
	contender := Contender{
		ID:          event.ContenderID,
		CompClassID: event.CompClassID,
		Ticks:       make(map[domain.ProblemID]*Tick),
	}

	e.contenders[event.ContenderID] = &contender

	e.rankCompClass(contender.CompClassID)
}

func (e *ScoreEngine) handleContenderSwitchedClass(event domain.ContenderSwitchedClassEvent) {
	contender, found := e.contenders[event.ContenderID]
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

	e.rankCompClasses(slices.Values(compClassesToReRank))
}

func (e *ScoreEngine) handleContenderWithdrewFromFinals(event domain.ContenderWithdrewFromFinalsEvent) {
	contender, found := e.contenders[event.ContenderID]
	if !found {
		return
	}

	contender.WithdrawnFromFinals = true

	e.rankCompClass(contender.CompClassID)
}

func (e *ScoreEngine) handleContenderReenteredFinals(event domain.ContenderReenteredFinalsEvent) {
	contender, found := e.contenders[event.ContenderID]
	if !found {
		return
	}

	contender.WithdrawnFromFinals = false

	e.rankCompClass(contender.CompClassID)
}

func (e *ScoreEngine) handleContenderDisqualified(event domain.ContenderDisqualifiedEvent) {
	contender, found := e.contenders[event.ContenderID]
	if !found {
		return
	}

	contender.Disqualified = true
	contender.Score = 0

	e.rankCompClass(contender.CompClassID)
}

func (e *ScoreEngine) handleContenderRequalified(event domain.ContenderRequalifiedEvent) {
	contender, found := e.contenders[event.ContenderID]
	if !found {
		return
	}

	contender.Disqualified = false
	e.scoreContender(contender)

	e.rankCompClass(contender.CompClassID)
}

func (e *ScoreEngine) handleAscentRegistered(event domain.AscentRegisteredEvent) {
	tick := &Tick{
		ProblemID:    event.ProblemID,
		Top:          event.Top,
		AttemptsTop:  event.AttemptsTop,
		Zone:         event.Zone,
		AttemptsZone: event.AttemptsZone,
	}

	contender, found := e.contenders[event.ContenderID]
	if !found {
		return
	}

	problem, found := e.problems[event.ProblemID]
	if !found {
		return
	}

	contender.Ticks[event.ProblemID] = tick

	tick.Score(*problem)

	e.scoreContender(contender)

	e.rankCompClass(contender.CompClassID)
}

func (e *ScoreEngine) handleAscentDeregistered(event domain.AscentDeregisteredEvent) {
	contender, found := e.contenders[event.ContenderID]
	if !found {
		return
	}

	delete(contender.Ticks, event.ProblemID)

	e.scoreContender(contender)

	e.rankCompClass(contender.CompClassID)
}

func (e *ScoreEngine) handleProblemAdded(event domain.ProblemAddedEvent) {
	e.problems[event.ProblemID] = &Problem{
		ID:         event.ProblemID,
		PointsTop:  event.PointsTop,
		PointsZone: event.PointsZone,
		FlashBonus: event.FlashBonus,
	}
}

func (e *ScoreEngine) scoreContender(contender *Contender) {
	tickPointValues := func(ticks map[domain.ProblemID]*Tick) iter.Seq[int] {
		return func(yield func(int) bool) {
			for _, tick := range ticks {
				if !yield(tick.Points) {
					return
				}
			}
		}
	}

	contender.Score = e.rules.CalculateScore(tickPointValues(contender.Ticks))
}

func (e *ScoreEngine) rankCompClass(compClassID domain.CompClassID) {
	e.rankCompClasses(slices.Values([]domain.CompClassID{compClassID}))
}

func (e *ScoreEngine) rankCompClasses(compClassIDs iter.Seq[domain.CompClassID]) {
	for compClassID := range compClassIDs {
		scores := e.ranker.RankContenders(FilterByCompClass(e.contenders, compClassID))

		for score := range slices.Values(scores) {
			e.scores.Set(score.ContenderID, score)
		}
	}
}

func (e *ScoreEngine) publishUpdatedScores() {
	diff := e.scores.Commit()

	var batch []domain.ContenderScoreUpdatedEvent

	for score := range slices.Values(diff) {
		e.eventBroker.Dispatch(e.contestID, domain.ContenderScoreUpdatedEvent(score))

		batch = append(batch, domain.ContenderScoreUpdatedEvent(score))
	}

	if len(batch) > 0 {
		e.eventBroker.Dispatch(e.contestID, batch)
	}
}
