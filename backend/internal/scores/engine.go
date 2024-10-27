package scores

import (
	"context"
	"iter"
	"log/slog"
	"slices"

	"github.com/climblive/platform/backend/internal/domain"
)

type ScoringRules interface {
	CalculateScore(tickPointValues iter.Seq[int]) int
}

type Ranker interface {
	RankContenders(contenders iter.Seq[*Contender]) []domain.Score
}

type ScoreEngine struct {
	contestID   domain.ContestID
	ranker      Ranker
	eventBroker domain.EventBroker
	rules       ScoringRules

	problems   map[domain.ProblemID]*Problem
	contenders map[domain.ContenderID]*Contender
	scores     DiffMap[domain.ContenderID, domain.Score]
}

func NewScoreEngine(contestID domain.ContestID, eventBroker domain.EventBroker, rules ScoringRules, ranker Ranker) ScoreEngine {
	engine := ScoreEngine{
		contestID:   contestID,
		ranker:      ranker,
		eventBroker: eventBroker,
		rules:       rules,
		problems:    make(map[domain.ProblemID]*Problem),
		contenders:  make(map[domain.ContenderID]*Contender),
		scores:      NewDiffMap[domain.ContenderID, domain.Score](CompareScore),
	}

	return engine
}

func (e *ScoreEngine) Run(ctx context.Context) {
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
	subscriptionID, eventReader := e.eventBroker.Subscribe(filter, 0)

	defer e.eventBroker.Unsubscribe(subscriptionID)

	for {
		event, err := eventReader.AwaitEvent(ctx)
		if err != nil {
			panic(err)
		}

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
	}
}

func (e *ScoreEngine) HandleContenderEntered(event domain.ContenderEnteredEvent) {
	contender := Contender{
		ID:          event.ContenderID,
		CompClassID: event.CompClassID,
		Ticks:       make(map[domain.ProblemID]*Tick),
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

	if contender.CompClassID == event.CompClassID {
		return
	}

	reRankList := []domain.CompClassID{
		contender.CompClassID,
		event.CompClassID,
	}

	contender.CompClassID = event.CompClassID

	for compClassID := range slices.Values(reRankList) {
		scores := e.ranker.RankContenders(FilterByClass(e.contenders, compClassID))

		for score := range slices.Values(scores) {
			e.scores.Set(domain.ContenderID(score.ContenderID), score)
		}
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
