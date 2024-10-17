package scores

import (
	"context"
	"iter"
	"log/slog"
	"slices"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/events"
)

type ScoringRules interface {
	CalculateScore(tickPointValues iter.Seq[int]) int
}

type Ranker interface {
	RankContenders(contenders iter.Seq[*Contender]) []domain.Score
}

type ScoreEngine struct {
	contestID   domain.ResourceID
	ranker      Ranker
	eventBroker domain.EventBroker
	rules       ScoringRules

	problems   map[domain.ResourceID]*Problem
	contenders map[domain.ResourceID]*Contender
	scores     DiffMap[domain.ContenderID, domain.Score]
}

func NewScoreEngine(contestID domain.ResourceID, eventBroker domain.EventBroker, rules ScoringRules, ranker Ranker) ScoreEngine {
	engine := ScoreEngine{
		contestID:   contestID,
		ranker:      ranker,
		eventBroker: eventBroker,
		rules:       rules,
		problems:    make(map[int]*Problem),
		contenders:  make(map[int]*Contender),
		scores:      NewDiffMap[domain.ContenderID, domain.Score](),
	}

	return engine
}

func (e *ScoreEngine) Run(ctx context.Context) {
	events := make(chan domain.EventContainer, events.EventChannelBufferSize)
	subscriptionID := e.eventBroker.Subscribe(domain.EventFilter{
		ContestID: e.contestID,
	}, events)

	defer e.eventBroker.Unsubscribe(subscriptionID)

	for {
		select {
		case event := <-events:
			switch ev := event.Data.(type) {
			case domain.ContenderEnteredEvent:
				e.HandleContenderEntered(ev)
			case domain.ContenderSwitchedClassEvent:
				e.HandleContenderSwitchedClass(ev)
			case domain.ContenderWithdrewFromFinalsEvent:
				e.HandleContenderWithdrewFromFinals(ev)
			case domain.ContenderReenteredFinalsEvent:
				e.HandleContenderReenteredFinals(ev)
			case domain.ContenderDisqualifiedEvent:
				e.HandleContenderDisqualified(ev)
			case domain.ContenderRequalifiedEvent:
				e.HandleContenderRequalified(ev)
			case domain.AscentRegisteredEvent:
				e.HandleAscentRegistered(ev)
			case domain.AscentDeregisteredEvent:
				e.HandleAscentDeregistered(ev)
			case domain.ProblemAddedEvent:
				e.HandleProblemAdded(ev)
			case domain.ProblemUpdatedEvent, domain.ProblemDeletedEvent:
				slog.Warn("discarding unsupported event", "event", event)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (e *ScoreEngine) HandleContenderEntered(event domain.ContenderEnteredEvent) {
	contender := Contender{
		ID:          event.ContenderID,
		CompClassID: event.CompClassID,
		Ticks:       make(map[int]*Tick),
	}

	e.contenders[event.ContenderID] = &contender

	scores := e.ranker.RankContenders(FilterByClass(e.contenders, contender.CompClassID))

	for score := range slices.Values(scores) {
		e.scores.Set(domain.ContenderID(score.ContenderID), score)
	}

	e.PublishUpdatedScores()
}

func (e *ScoreEngine) HandleContenderSwitchedClass(event domain.ContenderSwitchedClassEvent) {
	contender, found := e.contenders[event.ContenderID]
	if !found {
		return
	}

	contender.CompClassID = event.CompClassID

	scores := e.ranker.RankContenders(FilterByClass(e.contenders, contender.CompClassID))

	for score := range slices.Values(scores) {
		e.scores.Set(domain.ContenderID(score.ContenderID), score)
	}

	e.PublishUpdatedScores()
}
func (e *ScoreEngine) HandleContenderWithdrewFromFinals(event domain.ContenderWithdrewFromFinalsEvent) {
	contender, found := e.contenders[event.ContenderID]
	if !found {
		return
	}

	contender.WithdrawnFromFinals = true

	scores := e.ranker.RankContenders(FilterByClass(e.contenders, contender.CompClassID))

	for score := range slices.Values(scores) {
		e.scores.Set(domain.ContenderID(score.ContenderID), score)
	}

	e.PublishUpdatedScores()
}
func (e *ScoreEngine) HandleContenderReenteredFinals(event domain.ContenderReenteredFinalsEvent) {
	contender, found := e.contenders[event.ContenderID]
	if !found {
		return
	}

	contender.WithdrawnFromFinals = false

	scores := e.ranker.RankContenders(FilterByClass(e.contenders, contender.CompClassID))

	for score := range slices.Values(scores) {
		e.scores.Set(domain.ContenderID(score.ContenderID), score)
	}

	e.PublishUpdatedScores()
}
func (e *ScoreEngine) HandleContenderDisqualified(event domain.ContenderDisqualifiedEvent) {
	contender, found := e.contenders[event.ContenderID]
	if !found {
		return
	}

	contender.Disqualified = true
	contender.Score = 0

	scores := e.ranker.RankContenders(FilterByClass(e.contenders, contender.CompClassID))

	for score := range slices.Values(scores) {
		e.scores.Set(domain.ContenderID(score.ContenderID), score)
	}

	e.PublishUpdatedScores()
}
func (e *ScoreEngine) HandleContenderRequalified(event domain.ContenderRequalifiedEvent) {
	contender, found := e.contenders[event.ContenderID]
	if !found {
		return
	}

	contender.Disqualified = false
	e.ScoreContender(contender)

	scores := e.ranker.RankContenders(FilterByClass(e.contenders, contender.CompClassID))

	for score := range slices.Values(scores) {
		e.scores.Set(domain.ContenderID(score.ContenderID), score)
	}

	e.PublishUpdatedScores()
}

func (e *ScoreEngine) HandleAscentRegistered(event domain.AscentRegisteredEvent) {
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

	e.ScoreContender(contender)

	scores := e.ranker.RankContenders(FilterByClass(e.contenders, contender.CompClassID))

	for score := range slices.Values(scores) {
		e.scores.Set(domain.ContenderID(score.ContenderID), score)
	}

	e.PublishUpdatedScores()
}

func (e *ScoreEngine) HandleAscentDeregistered(event domain.AscentDeregisteredEvent) {
	contender, found := e.contenders[event.ContenderID]

	if !found {
		return
	}

	delete(contender.Ticks, event.ProblemID)

	e.ScoreContender(contender)

	scores := e.ranker.RankContenders(FilterByClass(e.contenders, contender.CompClassID))

	for score := range slices.Values(scores) {
		e.scores.Set(domain.ContenderID(score.ContenderID), score)
	}

	e.PublishUpdatedScores()
}

func (e *ScoreEngine) HandleProblemAdded(event domain.ProblemAddedEvent) {
	e.problems[event.ProblemID] = &Problem{
		ID:         event.ProblemID,
		PointsTop:  event.PointsTop,
		PointsZone: event.PointsZone,
		FlashBonus: event.FlashBonus,
	}
}

func (e *ScoreEngine) ScoreContender(contender *Contender) {
	tickPointValues := func(ticks map[domain.ResourceID]*Tick) iter.Seq[int] {
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

func (e *ScoreEngine) PublishUpdatedScores() {
	diff := e.scores.Commit()

	for score := range slices.Values(diff) {
		e.eventBroker.Dispatch(e.contestID, domain.ContenderScoreUpdatedEvent{
			Timestamp:   score.Timestamp,
			ContenderID: score.ContenderID,
			Score:       score.Score,
			Placement:   score.Placement,
			Finalist:    score.Finalist,
			RankOrder:   score.RankOrder,
		})
	}
}
