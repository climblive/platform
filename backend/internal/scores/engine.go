package scores

import (
	"context"
	"iter"
	"slices"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
)

type ScoringRules interface {
	CalculateScore(ticks iter.Seq[int]) int
	CalculatePoints(tick Tick, problem Problem) int
}

type XHardest struct {
	Number int
}

func (r *XHardest) CalculateScore(ticks iter.Seq[int]) int {
	score := 0

	n := 0
	for tick := range slices.Sorted(ticks) {
		n++
		score += tick

		if n >= r.Number {
			return score
		}
	}

	return score
}

func (r *XHardest) CalculatePoints(tick Tick, problem Problem) int {
	points := 0

	if tick.Zone {
		points = problem.PointsZone
	}

	if tick.Top {
		points = problem.PointsTop
	}

	if tick.AttemptsTop == 1 {
		points += problem.FlashBonus
	}

	return points
}

func NewScoreEngine(contestID domain.ResourceID, eventBroker domain.EventBroker, rules ScoringRules, numberOfFinalists int) ScoreEngine {
	engine := ScoreEngine{
		contestID:         contestID,
		numberOfFinalists: numberOfFinalists,
		eventBroker:       eventBroker,
		rules:             rules,
		problems:          make(map[int]*Problem),
		contenders:        make(map[int]*Contender),
	}

	return engine
}

type Contender struct {
	ID                  domain.ResourceID
	CompClassID         domain.ResourceID
	Disqualified        bool
	WithdrawnFromFinals bool
	Ticks               map[domain.ResourceID]*TickWithPoints
	Score               int
	Placement           int
	Finalist            bool
	RankOrder           int
}

type Tick struct {
	ProblemID    domain.ResourceID
	Top          bool
	AttemptsTop  int
	Zone         bool
	AttemptsZone int
}

type TickWithPoints struct {
	Tick
	Points int
}

type Problem struct {
	ID         domain.ResourceID
	PointsTop  int
	PointsZone int
	FlashBonus int
}

type ScoreEngine struct {
	contestID         domain.ResourceID
	numberOfFinalists int
	eventBroker       domain.EventBroker
	rules             ScoringRules

	problems   map[domain.ResourceID]*Problem
	contenders map[domain.ResourceID]*Contender
}

func InCompClass(contenders map[domain.ResourceID]*Contender, compClassID domain.ResourceID) iter.Seq[*Contender] {
	return func(yield func(*Contender) bool) {
		for _, contender := range contenders {
			if contender.CompClassID != compClassID {
				continue
			}

			if !yield(contender) {
				return
			}
		}
	}
}

func (e *ScoreEngine) Run(ctx context.Context, contestID domain.ResourceID) {
	events := make(chan domain.EventContainer, 1000)
	subscriptionID := e.eventBroker.Subscribe(domain.EventFilter{
		ContestID: contestID,
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
			case domain.ProblemUpdatedEvent:
				e.HandleProblemUpdated(ev)
			case domain.ProblemDeletedEvent:
				e.HandleProblemDeleted(ev)
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
		Ticks:       make(map[int]*TickWithPoints),
	}

	e.contenders[event.ContenderID] = &contender
}

func (e *ScoreEngine) HandleContenderSwitchedClass(event domain.ContenderSwitchedClassEvent) {
	if contender, found := e.contenders[event.ContenderID]; found {
		contender.CompClassID = event.CompClassID
	}
}
func (e *ScoreEngine) HandleContenderWithdrewFromFinals(event domain.ContenderWithdrewFromFinalsEvent) {
	if contender, found := e.contenders[event.ContenderID]; found {
		contender.WithdrawnFromFinals = true
	}
}
func (e *ScoreEngine) HandleContenderReenteredFinals(event domain.ContenderReenteredFinalsEvent) {
	if contender, found := e.contenders[event.ContenderID]; found {
		contender.WithdrawnFromFinals = false
	}
}
func (e *ScoreEngine) HandleContenderDisqualified(event domain.ContenderDisqualifiedEvent) {
	if contender, found := e.contenders[event.ContenderID]; found {
		contender.Disqualified = true
	}
}
func (e *ScoreEngine) HandleContenderRequalified(event domain.ContenderRequalifiedEvent) {
	if contender, found := e.contenders[event.ContenderID]; found {
		contender.Disqualified = false
	}
}

func (e *ScoreEngine) HandleAscentRegistered(event domain.AscentRegisteredEvent) {
	tick := &TickWithPoints{
		Tick: Tick{
			ProblemID:    event.ProblemID,
			Top:          event.Top,
			AttemptsTop:  event.AttemptsTop,
			Zone:         event.Zone,
			AttemptsZone: event.AttemptsZone,
		},
	}

	contender, found := e.contenders[event.ContenderID]
	if !found {
		panic(domain.ErrNotFound)
	}

	contender.Ticks[event.ProblemID] = tick

	problem, found := e.problems[event.ProblemID]
	if !found {
		panic(domain.ErrNotFound)
	}

	tick.Points = e.rules.CalculatePoints(tick.Tick, *problem)

	var contenders []Contender

	for contender := range InCompClass(e.contenders, contender.CompClassID) {
		contenders = append(contenders, *contender)
	}

	e.ScoreContenders(contenders)
	e.RankContenders(contenders)
	e.Commit(contenders)
}

func (e *ScoreEngine) HandleAscentDeregistered(event domain.AscentDeregisteredEvent) {
	contender, found := e.contenders[event.ContenderID]
	if !found {
		panic(domain.ErrNotFound)
	}

	delete(contender.Ticks, event.ProblemID)

	var contenders []Contender

	for contender := range InCompClass(e.contenders, contender.CompClassID) {
		contenders = append(contenders, *contender)
	}

	e.ScoreContenders(contenders)
	e.RankContenders(contenders)
	e.Commit(contenders)
}

func (e *ScoreEngine) HandleProblemAdded(event domain.ProblemAddedEvent) {
	e.problems[event.ProblemID] = &Problem{
		ID:         event.ProblemID,
		PointsTop:  event.PointsTop,
		PointsZone: event.PointsZone,
		FlashBonus: event.FlashBonus,
	}
}

func (e *ScoreEngine) HandleProblemUpdated(event domain.ProblemUpdatedEvent) {
	e.problems[event.ProblemID] = &Problem{
		ID:         event.ProblemID,
		PointsTop:  event.PointsTop,
		PointsZone: event.PointsZone,
		FlashBonus: event.FlashBonus,
	}
}

func (e *ScoreEngine) HandleProblemDeleted(event domain.ProblemDeletedEvent) {
	delete(e.problems, event.ProblemID)
}

func (e *ScoreEngine) ScoreContenders(contenders []Contender) {
	tickValues := func(ticks map[domain.ResourceID]*TickWithPoints) iter.Seq[int] {
		return func(yield func(int) bool) {
			for _, tick := range ticks {
				if !yield(tick.Points) {
					return
				}
			}
		}
	}

	for i, contender := range contenders {
		contender.Score = e.rules.CalculateScore(tickValues(contender.Ticks))

		contenders[i] = contender
	}
}

func (e *ScoreEngine) RankContenders(contenders []Contender) {
	slices.SortFunc(contenders, func(c1, c2 Contender) int {
		return c2.Score - c1.Score
	})

	var previousContender *Contender
	var placement int
	var gap int
	var numberOfAssignedFinalists int
	var lastFinalistPlacement int

	for i, contender := range contenders {
		switch {
		case previousContender == nil:
			placement = 1
			gap = 0
		case contender.Score == previousContender.Score:
			gap++
		case contender.Score != previousContender.Score:
			placement += 1 + gap
			gap = 0
		}

		contender.Placement = placement
		contender.RankOrder = i

		switch {
		case contender.Score == 0:
			fallthrough
		case contender.WithdrawnFromFinals:
			contender.Finalist = false
		case numberOfAssignedFinalists < e.numberOfFinalists:
			contender.Finalist = true
			numberOfAssignedFinalists++
			lastFinalistPlacement = contender.Placement
		case contender.Placement == lastFinalistPlacement:
			contender.Finalist = true
		default:
			contender.Finalist = false
		}

		previousContender = &contender

		contenders[i] = contender
	}
}

func (e *ScoreEngine) Commit(contenders []Contender) {
	for contender := range slices.Values(contenders) {
		stored := e.contenders[contender.ID]

		if stored.Score != contender.Score || stored.Placement != contender.Placement || stored.Finalist != contender.Finalist || stored.RankOrder != contender.RankOrder {
			e.contenders[contender.ID] = &contender

			e.eventBroker.Dispatch(e.contestID, domain.ContenderScoreUpdatedEvent{
				Timestamp:   time.Now(),
				ContenderID: contender.ID,
				Score:       contender.Score,
				Placement:   contender.Placement,
				Finalist:    contender.Finalist,
				RankOrder:   contender.RankOrder,
			})
		}
	}
}
