package scores

import (
	"encoding/binary"
	"iter"
	"slices"

	"github.com/climblive/platform/backend/internal/domain"
)

type Rules struct {
	QualifyingProblems int
	Finalists          int
	UsePoints          bool
	PooledPoints       bool
	MaxAttempts        int
	PointDeduction     domain.Tenths
}

type ScoringRules interface {
	CalculateScore(points iter.Seq[domain.Tenths]) domain.Tenths
}

type Ranker interface {
	RankCompetitors(competitors iter.Seq[Competitor]) []domain.Result
}

type EffectType int8

const (
	EffectTypeCalculatePointValues EffectType = iota
	EffectTypeScoreCompetitor
	EffectTypeRankCategory
)

type EncodedEffect = [9]byte

type Effect interface {
	Type() EffectType
	Encode() EncodedEffect
}

type EffectScoreCompetitor struct {
	CompetitorID domain.CompetitorID
}

func (e EffectScoreCompetitor) Type() EffectType {
	return EffectTypeScoreCompetitor
}

func (e EffectScoreCompetitor) Encode() EncodedEffect {
	var data EncodedEffect
	data[0] = byte(EffectTypeScoreCompetitor)
	binary.LittleEndian.PutUint32(data[1:], uint32(e.CompetitorID))
	return data
}

type EffectRankCategory struct {
	CategoryID domain.CategoryID
}

func (e EffectRankCategory) Type() EffectType {
	return EffectTypeRankCategory
}

func (e EffectRankCategory) Encode() EncodedEffect {
	var data EncodedEffect
	data[0] = byte(EffectTypeRankCategory)
	binary.LittleEndian.PutUint32(data[1:], uint32(e.CategoryID))
	return data
}

type EffectCalculatePointValues struct {
	CategoryID domain.CategoryID
	ProblemID  domain.ProblemID
}

func (e EffectCalculatePointValues) Type() EffectType {
	return EffectTypeCalculatePointValues
}

func (e EffectCalculatePointValues) Encode() EncodedEffect {
	var data EncodedEffect
	data[0] = byte(EffectTypeCalculatePointValues)
	binary.LittleEndian.PutUint32(data[1:], uint32(e.ProblemID))
	binary.LittleEndian.PutUint32(data[5:], uint32(e.CategoryID))
	return data
}

type EngineStore interface {
	GetRules() Rules
	SaveRules(Rules)

	GetCompetitor(domain.CompetitorID) (Competitor, bool)
	GetCompetitorsByCategory(domain.CategoryID) iter.Seq[Competitor]
	SaveCompetitor(Competitor)

	GetCategoryIDs() []domain.CategoryID

	GetTicksByCompetitor(domain.CompetitorID) iter.Seq[Tick]
	GetTick(domain.CompetitorID, domain.ProblemID) (Tick, bool)
	SaveTick(domain.CompetitorID, Tick)
	DeleteTick(domain.TickID)
	GetTicksByProblem(domain.CategoryID, domain.ProblemID) iter.Seq[Tick]

	GetProblem(domain.ProblemID) (Problem, bool)
	SaveProblem(Problem)
	GetAllProblems() iter.Seq[Problem]

	GetPointValue(domain.CompetitorID, domain.ProblemID) (domain.PointValue, bool)
	SavePointValue(domain.CompetitorID, domain.ProblemID, domain.PointValue)
	GetDirtyPointValues() []domain.PointValue

	SaveResult(domain.Result)
	GetDirtyResults() []domain.Result
}

type DefaultScoreEngine struct {
	store EngineStore
}

func NewDefaultScoreEngine(store EngineStore) *DefaultScoreEngine {
	return &DefaultScoreEngine{
		store: store,
	}
}

func (e *DefaultScoreEngine) Start() iter.Seq[Effect] {
	return func(yield func(Effect) bool) {
		for _, categoryID := range e.store.GetCategoryIDs() {
			for problem := range e.store.GetAllProblems() {
				if !yield(EffectCalculatePointValues{CategoryID: categoryID, ProblemID: problem.ID}) {
					return
				}
			}

			for competitor := range e.store.GetCompetitorsByCategory(categoryID) {
				if !yield(EffectScoreCompetitor{CompetitorID: competitor.ID}) {
					return
				}
			}

			if !yield(EffectRankCategory{CategoryID: categoryID}) {
				return
			}
		}
	}
}

func (e *DefaultScoreEngine) Stop() {
}

func (e *DefaultScoreEngine) HandleRulesUpdated(event domain.RulesUpdatedEvent) iter.Seq[Effect] {
	rules := Rules{
		QualifyingProblems: event.QualifyingProblems,
		Finalists:          event.Finalists,
		UsePoints:          event.UsePoints,
		PooledPoints:       event.PooledPoints,
		MaxAttempts:        event.MaxAttempts,
		PointDeduction:     event.PointDeduction,
	}

	e.store.SaveRules(rules)

	return e.Start()
}

func (e *DefaultScoreEngine) HandleCompetitorEntered(event domain.CompetitorEnteredEvent) iter.Seq[Effect] {
	competitor := Competitor{
		ID:                  event.CompetitorID,
		CategoryID:          event.CategoryID,
		Disqualified:        false,
		WithdrawnFromFinals: false,
		Score:               Score{},
	}

	e.store.SaveCompetitor(competitor)

	return func(yield func(Effect) bool) {
		for problem := range e.store.GetAllProblems() {
			if !yield(EffectCalculatePointValues{CategoryID: competitor.CategoryID, ProblemID: problem.ID}) {
				return
			}
		}

		if !yield(EffectRankCategory{CategoryID: competitor.CategoryID}) {
			return
		}
	}
}

func (e *DefaultScoreEngine) HandleCompetitorSwitchedCategory(event domain.CompetitorSwitchedCategoryEvent) iter.Seq[Effect] {
	competitor, found := e.store.GetCompetitor(event.CompetitorID)
	if !found {
		return nil
	}

	if competitor.CategoryID == event.CategoryID {
		return nil
	}

	oldCategoryID := competitor.CategoryID

	competitor.CategoryID = event.CategoryID

	e.store.SaveCompetitor(competitor)

	return func(yield func(Effect) bool) {
		for problem := range e.store.GetAllProblems() {
			if !yield(EffectCalculatePointValues{CategoryID: oldCategoryID, ProblemID: problem.ID}) {
				return
			}

			if !yield(EffectCalculatePointValues{CategoryID: event.CategoryID, ProblemID: problem.ID}) {
				return
			}
		}

		if !yield(EffectRankCategory{CategoryID: oldCategoryID}) {
			return
		}

		if !yield(EffectRankCategory{CategoryID: event.CategoryID}) {
			return
		}
	}
}

func (e *DefaultScoreEngine) HandleCompetitorWithdrewFromFinals(event domain.CompetitorWithdrewFromFinalsEvent) iter.Seq[Effect] {
	competitor, found := e.store.GetCompetitor(event.CompetitorID)
	if !found {
		return nil
	}

	competitor.WithdrawnFromFinals = true

	e.store.SaveCompetitor(competitor)

	return func(yield func(Effect) bool) {
		if !yield(EffectRankCategory{CategoryID: competitor.CategoryID}) {
			return
		}
	}
}

func (e *DefaultScoreEngine) HandleCompetitorReenteredFinals(event domain.CompetitorReenteredFinalsEvent) iter.Seq[Effect] {
	competitor, found := e.store.GetCompetitor(event.CompetitorID)
	if !found {
		return nil
	}

	competitor.WithdrawnFromFinals = false

	e.store.SaveCompetitor(competitor)

	return func(yield func(Effect) bool) {
		if !yield(EffectRankCategory{CategoryID: competitor.CategoryID}) {
			return
		}
	}
}

func (e *DefaultScoreEngine) HandleCompetitorDisqualified(event domain.CompetitorDisqualifiedEvent) iter.Seq[Effect] {
	competitor, found := e.store.GetCompetitor(event.CompetitorID)
	if !found {
		return nil
	}

	competitor.Disqualified = true

	e.store.SaveCompetitor(competitor)

	return func(yield func(Effect) bool) {
		for problem := range e.store.GetAllProblems() {
			if !yield(EffectCalculatePointValues{CategoryID: competitor.CategoryID, ProblemID: problem.ID}) {
				return
			}
		}

		if !yield(EffectScoreCompetitor{CompetitorID: competitor.ID}) {
			return
		}

		if !yield(EffectRankCategory{CategoryID: competitor.CategoryID}) {
			return
		}
	}
}

func (e *DefaultScoreEngine) HandleCompetitorRequalified(event domain.CompetitorRequalifiedEvent) iter.Seq[Effect] {
	competitor, found := e.store.GetCompetitor(event.CompetitorID)
	if !found {
		return nil
	}

	competitor.Disqualified = false

	e.store.SaveCompetitor(competitor)

	return func(yield func(Effect) bool) {
		for problem := range e.store.GetAllProblems() {
			if !yield(EffectCalculatePointValues{CategoryID: competitor.CategoryID, ProblemID: problem.ID}) {
				return
			}
		}

		if !yield(EffectScoreCompetitor{CompetitorID: competitor.ID}) {
			return
		}

		if !yield(EffectRankCategory{CategoryID: competitor.CategoryID}) {
			return
		}
	}
}

func (e *DefaultScoreEngine) HandleAscentRegistered(event domain.AscentRegisteredEvent) iter.Seq[Effect] {
	tick := Tick{
		ID:            event.TickID,
		Revision:      event.Revision,
		CompetitorID:  event.CompetitorID,
		ProblemID:     event.ProblemID,
		Zone1:         event.Zone1,
		AttemptsZone1: event.AttemptsZone1,
		Zone2:         event.Zone2,
		AttemptsZone2: event.AttemptsZone2,
		Top:           event.Top,
		AttemptsTop:   event.AttemptsTop,
	}

	existingTick, found := e.store.GetTick(event.CompetitorID, event.ProblemID)
	if found && existingTick.Revision > event.Revision {
		return nil
	}

	competitor, found := e.store.GetCompetitor(event.CompetitorID)
	if !found {
		return nil
	}

	e.store.SaveTick(event.CompetitorID, tick)

	if competitor.Disqualified {
		return nil
	}

	return func(yield func(Effect) bool) {
		if !yield(EffectCalculatePointValues{CategoryID: competitor.CategoryID, ProblemID: event.ProblemID}) {
			return
		}

		if !yield(EffectScoreCompetitor{CompetitorID: competitor.ID}) {
			return
		}
	}
}

func (e *DefaultScoreEngine) HandleAscentDeregistered(event domain.AscentDeregisteredEvent) iter.Seq[Effect] {
	competitor, found := e.store.GetCompetitor(event.CompetitorID)
	if !found {
		return nil
	}

	e.store.DeleteTick(event.TickID)

	if competitor.Disqualified {
		return nil
	}

	return func(yield func(Effect) bool) {
		if !yield(EffectCalculatePointValues{CategoryID: competitor.CategoryID, ProblemID: event.ProblemID}) {
			return
		}

		if !yield(EffectScoreCompetitor{CompetitorID: competitor.ID}) {
			return
		}
	}
}

func (e *DefaultScoreEngine) HandleProblemAdded(event domain.ProblemAddedEvent) iter.Seq[Effect] {
	problem := Problem{
		ID: event.ProblemID,
		ProblemValue: ProblemValue{
			PointsZone1: domain.Tenths(event.PointsZone1) * 10,
			PointsZone2: domain.Tenths(event.PointsZone2) * 10,
			PointsTop:   domain.Tenths(event.PointsTop) * 10,
			FlashBonus:  domain.Tenths(event.FlashBonus) * 10,
		},
		Zone1Enabled: event.Zone1Enabled,
		Zone2Enabled: event.Zone2Enabled,
	}

	e.store.SaveProblem(problem)

	return func(yield func(Effect) bool) {
		for _, categoryID := range e.store.GetCategoryIDs() {
			if !yield(EffectRankCategory{CategoryID: categoryID}) {
				return
			}

			if !yield(EffectCalculatePointValues{CategoryID: categoryID, ProblemID: event.ProblemID}) {
				return
			}
		}
	}
}

func (e *DefaultScoreEngine) HandleProblemUpdated(event domain.ProblemUpdatedEvent) iter.Seq[Effect] {
	problem := Problem{
		ID: event.ProblemID,
		ProblemValue: ProblemValue{
			PointsZone1: domain.Tenths(event.PointsZone1) * 10,
			PointsZone2: domain.Tenths(event.PointsZone2) * 10,
			PointsTop:   domain.Tenths(event.PointsTop) * 10,
			FlashBonus:  domain.Tenths(event.FlashBonus) * 10,
		},
		Zone1Enabled: event.Zone1Enabled,
		Zone2Enabled: event.Zone2Enabled,
	}

	e.store.SaveProblem(problem)

	return func(yield func(Effect) bool) {
		for _, categoryID := range e.store.GetCategoryIDs() {
			if !yield(EffectRankCategory{CategoryID: categoryID}) {
				return
			}

			if !yield(EffectCalculatePointValues{CategoryID: categoryID, ProblemID: event.ProblemID}) {
				return
			}

			competitors := e.store.GetCompetitorsByCategory(categoryID)

			for competitor := range competitors {
				if !yield(EffectScoreCompetitor{CompetitorID: competitor.ID}) {
					return
				}
			}
		}
	}
}

func (e *DefaultScoreEngine) GetDirtyResults() []domain.Result {
	return e.store.GetDirtyResults()
}

func (e *DefaultScoreEngine) GetDirtyPointValues() []domain.PointValue {
	return e.store.GetDirtyPointValues()
}

func (e *DefaultScoreEngine) CalculatePointValues(categoryID domain.CategoryID, problemID domain.ProblemID) iter.Seq[Effect] {
	rules := e.store.GetRules()

	if !rules.UsePoints {
		return nil
	}

	problem, found := e.store.GetProblem(problemID)
	if !found {
		return nil
	}

	problemValue := problem.ProblemValue
	tickPool := TickPool{}

	if rules.PooledPoints {
		for tick := range e.store.GetTicksByProblem(categoryID, problemID) {
			competitor, found := e.store.GetCompetitor(tick.CompetitorID)
			if !found {
				continue
			}

			if competitor.Disqualified {
				continue
			}

			tickPool = tickPool.Add(tick)
		}

		problemValue = tickPool.CalculatePooledProblemValue(problem.ProblemValue)
	}

	affectedCompetitors := make([]domain.CompetitorID, 0)

	for competitor := range e.store.GetCompetitorsByCategory(categoryID) {
		tick, _ := e.store.GetTick(competitor.ID, problemID)

		hypotheticalBestZone1 := tick.TurnIntoZone1()
		hypotheticalBestZone2 := tick.TurnIntoZone2()
		hypotheticalBestTop := tick.TurnIntoTop()
		hypotheticalFlash := tick.TurnIntoFlash()

		pointValue := domain.PointValue{
			CompetitorID: competitor.ID,
			ProblemID:    problemID,
			Current:      0,
			Zone1:        0,
			Zone2:        0,
			Top:          0,
			FlashBonus:   0,
		}

		if !competitor.Disqualified {
			pointValue.Current = CalculatePoints(problemValue, tick, rules)

			hypotheticalProblemValue := problemValue
			if rules.PooledPoints {
				hypotheticalProblemValue = tickPool.Sub(tick).Add(hypotheticalFlash).CalculatePooledProblemValue(problem.ProblemValue)
			}

			if hypotheticalBestTop.AttemptsTop == 1 {
				pointValue.FlashBonus = hypotheticalProblemValue.FlashBonus
			}
			hypotheticalProblemValue.FlashBonus = 0

			pointValue.Zone1 = CalculatePoints(hypotheticalProblemValue, hypotheticalBestZone1, rules)
			pointValue.Zone2 = CalculatePoints(hypotheticalProblemValue, hypotheticalBestZone2, rules)
			pointValue.Top = CalculatePoints(hypotheticalProblemValue, hypotheticalBestTop, rules)
		}

		oldValue, hadPointValue := e.store.GetPointValue(competitor.ID, problemID)
		e.store.SavePointValue(competitor.ID, problemID, pointValue)

		if !hadPointValue || !ComparePointValue(oldValue, pointValue) {
			affectedCompetitors = append(affectedCompetitors, competitor.ID)
		}
	}

	if len(affectedCompetitors) == 0 {
		return nil
	}

	return func(yield func(Effect) bool) {
		for _, competitorID := range affectedCompetitors {
			if !yield(EffectScoreCompetitor{CompetitorID: competitorID}) {
				return
			}
		}
	}
}

func (e *DefaultScoreEngine) ScoreCompetitor(competitorID domain.CompetitorID) iter.Seq[Effect] {
	competitor, found := e.store.GetCompetitor(competitorID)
	if !found {
		return nil
	}

	rules := e.store.GetRules()

	oldScore := competitor.Score

	competitor.Score = Score{}

	if !competitor.Disqualified {
		if rules.UsePoints {
			ticks := e.store.GetTicksByCompetitor(competitor.ID)

			var pointValues iter.Seq[domain.Tenths] = func(yield func(domain.Tenths) bool) {
				for tick := range ticks {
					value, found := e.store.GetPointValue(competitor.ID, tick.ProblemID)
					if !found {
						continue
					}

					if !yield(value.Current) {
						return
					}
				}
			}

			scorer := Scorer{
				ProblemLimit: rules.QualifyingProblems,
			}

			competitor.Points = scorer.CalculatePoints(pointValues)
		}

		ticks := e.store.GetTicksByCompetitor(competitor.ID)

		for tick := range ticks {
			problem, found := e.store.GetProblem(tick.ProblemID)
			if !found {
				continue
			}

			if tick.Top {
				competitor.Tops += 1
				competitor.AttemptsTops += tick.AttemptsTop
			}

			if tick.Zone2 && problem.Zone2Enabled {
				competitor.Zone2s += 1
				competitor.AttemptsZone2s += tick.AttemptsZone2
			}

			if tick.Zone1 && problem.Zone1Enabled {
				competitor.Zone1s += 1
				competitor.AttemptsZone1s += tick.AttemptsZone1
			}
		}
	}

	if competitor.Score == oldScore {
		return nil
	}

	e.store.SaveCompetitor(competitor)

	return func(yield func(Effect) bool) {
		yield(EffectRankCategory{CategoryID: competitor.CategoryID})
	}
}

func (e *DefaultScoreEngine) RankCategory(categoryID domain.CategoryID) {
	rules := e.store.GetRules()
	var zone1Enabled, zone2Enabled bool
	if !rules.UsePoints {
		for problem := range e.store.GetAllProblems() {
			zone1Enabled = zone1Enabled || problem.Zone1Enabled
			zone2Enabled = zone2Enabled || problem.Zone2Enabled
		}
	}
	ranker := NewBasicRanker(rules.Finalists, rules.UsePoints, zone1Enabled, zone2Enabled)

	results := ranker.RankCompetitors(e.store.GetCompetitorsByCategory(categoryID))

	for result := range slices.Values(results) {
		e.store.SaveResult(result)
	}
}
