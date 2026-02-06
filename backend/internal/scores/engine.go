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
}

type ScoringRules interface {
	CalculateScore(points iter.Seq[int]) int
}

type Ranker interface {
	RankContenders(contenders iter.Seq[Contender]) []domain.Score
}

type EffectType int8

const (
	EffectTypeCalculateProblemValue EffectType = iota
	EffectTypeScoreContender
	EffectTypeRankClass
)

type EncodedEffect = [9]byte

type Effect interface {
	Encode() EncodedEffect
}

type EffectScoreContender struct {
	ContenderID domain.ContenderID
}

func (e EffectScoreContender) Encode() EncodedEffect {
	var data EncodedEffect
	data[0] = byte(EffectTypeScoreContender)
	binary.LittleEndian.PutUint32(data[1:], uint32(e.ContenderID))
	return data
}

type EffectRankClass struct {
	CompClassID domain.CompClassID
}

func (e EffectRankClass) Encode() EncodedEffect {
	var data EncodedEffect
	data[0] = byte(EffectTypeRankClass)
	binary.LittleEndian.PutUint32(data[1:], uint32(e.CompClassID))
	return data
}

type EffectCalculateProblemValue struct {
	CompClassID domain.CompClassID
	ProblemID   domain.ProblemID
}

func (e EffectCalculateProblemValue) Encode() EncodedEffect {
	var data EncodedEffect
	data[0] = byte(EffectTypeCalculateProblemValue)
	binary.LittleEndian.PutUint32(data[1:], uint32(e.ProblemID))
	binary.LittleEndian.PutUint32(data[5:], uint32(e.CompClassID))
	return data
}

type EngineStore interface {
	GetRules() Rules
	SaveRules(Rules)

	GetContender(domain.ContenderID) (Contender, bool)
	GetContendersByCompClass(domain.CompClassID) iter.Seq[Contender]
	SaveContender(Contender)

	GetCompClassIDs() []domain.CompClassID

	GetTicksByContender(domain.ContenderID) iter.Seq[Tick]
	SaveTick(domain.ContenderID, Tick)
	DeleteTick(domain.ContenderID, domain.ProblemID)
	GetTicksByProblem(domain.CompClassID, domain.ProblemID) iter.Seq[Tick]

	GetProblem(domain.ProblemID) (Problem, bool)
	SaveProblem(Problem)
	GetAllProblems() iter.Seq[Problem]

	GetProblemValue(domain.CompClassID, domain.ProblemID) (ProblemValue, bool)
	SaveProblemValue(domain.CompClassID, domain.ProblemID, ProblemValue)
	GetDirtyProblemValues() []ProblemValue

	SaveScore(domain.Score)
	GetDirtyScores() []domain.Score
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
		for _, compClassID := range e.store.GetCompClassIDs() {
			for problem := range e.store.GetAllProblems() {
				yield(EffectCalculateProblemValue{CompClassID: compClassID, ProblemID: problem.ID})
			}

			for contender := range e.store.GetContendersByCompClass(compClassID) {
				yield(EffectScoreContender{ContenderID: contender.ID})
			}
		}

		for _, compClassID := range e.store.GetCompClassIDs() {
			yield(EffectRankClass{CompClassID: compClassID})
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
	}

	e.store.SaveRules(rules)

	return e.Start()
}

func (e *DefaultScoreEngine) HandleContenderEntered(event domain.ContenderEnteredEvent) iter.Seq[Effect] {
	contender := Contender{
		ID:          event.ContenderID,
		CompClassID: event.CompClassID,
	}

	e.store.SaveContender(contender)

	return func(yield func(Effect) bool) {
		yield(EffectRankClass{CompClassID: contender.CompClassID})
	}
}

func (e *DefaultScoreEngine) HandleContenderSwitchedClass(event domain.ContenderSwitchedClassEvent) iter.Seq[Effect] {
	contender, found := e.store.GetContender(event.ContenderID)
	if !found {
		return nil
	}

	if contender.CompClassID == event.CompClassID {
		return nil
	}

	oldCompClassID := contender.CompClassID

	contender.CompClassID = event.CompClassID

	e.store.SaveContender(contender)

	return func(yield func(Effect) bool) {
		for tick := range e.store.GetTicksByContender(contender.ID) {
			yield(EffectCalculateProblemValue{CompClassID: oldCompClassID, ProblemID: tick.ProblemID})
			yield(EffectCalculateProblemValue{CompClassID: event.CompClassID, ProblemID: tick.ProblemID})
		}

		yield(EffectRankClass{CompClassID: oldCompClassID})
		yield(EffectRankClass{CompClassID: event.CompClassID})
	}
}

func (e *DefaultScoreEngine) HandleContenderWithdrewFromFinals(event domain.ContenderWithdrewFromFinalsEvent) iter.Seq[Effect] {
	contender, found := e.store.GetContender(event.ContenderID)
	if !found {
		return nil
	}

	contender.WithdrawnFromFinals = true

	e.store.SaveContender(contender)

	return func(yield func(Effect) bool) {
		yield(EffectRankClass{CompClassID: contender.CompClassID})
	}
}

func (e *DefaultScoreEngine) HandleContenderReenteredFinals(event domain.ContenderReenteredFinalsEvent) iter.Seq[Effect] {
	contender, found := e.store.GetContender(event.ContenderID)
	if !found {
		return nil
	}

	contender.WithdrawnFromFinals = false

	e.store.SaveContender(contender)

	return func(yield func(Effect) bool) {
		yield(EffectRankClass{CompClassID: contender.CompClassID})
	}
}

func (e *DefaultScoreEngine) HandleContenderDisqualified(event domain.ContenderDisqualifiedEvent) iter.Seq[Effect] {
	contender, found := e.store.GetContender(event.ContenderID)
	if !found {
		return nil
	}

	contender.Disqualified = true

	e.store.SaveContender(contender)

	return func(yield func(Effect) bool) {
		yield(EffectScoreContender{ContenderID: contender.ID})
		yield(EffectRankClass{CompClassID: contender.CompClassID})
	}
}

func (e *DefaultScoreEngine) HandleContenderRequalified(event domain.ContenderRequalifiedEvent) iter.Seq[Effect] {
	contender, found := e.store.GetContender(event.ContenderID)
	if !found {
		return nil
	}

	contender.Disqualified = false

	e.store.SaveContender(contender)

	return func(yield func(Effect) bool) {
		yield(EffectScoreContender{ContenderID: contender.ID})
		yield(EffectRankClass{CompClassID: contender.CompClassID})
	}
}

func (e *DefaultScoreEngine) HandleAscentRegistered(event domain.AscentRegisteredEvent) iter.Seq[Effect] {
	tick := Tick{
		ContenderID:   event.ContenderID,
		ProblemID:     event.ProblemID,
		Zone1:         event.Zone1,
		AttemptsZone1: event.AttemptsZone1,
		Zone2:         event.Zone2,
		AttemptsZone2: event.AttemptsZone2,
		Top:           event.Top,
		AttemptsTop:   event.AttemptsTop,
	}

	contender, found := e.store.GetContender(event.ContenderID)
	if !found {
		return nil
	}

	e.store.SaveTick(event.ContenderID, tick)

	if contender.Disqualified {
		return nil
	}

	return func(yield func(Effect) bool) {
		yield(EffectCalculateProblemValue{CompClassID: contender.CompClassID, ProblemID: event.ProblemID})
		yield(EffectScoreContender{ContenderID: contender.ID})
	}
}

func (e *DefaultScoreEngine) HandleAscentDeregistered(event domain.AscentDeregisteredEvent) iter.Seq[Effect] {
	contender, found := e.store.GetContender(event.ContenderID)
	if !found {
		return nil
	}

	e.store.DeleteTick(event.ContenderID, event.ProblemID)

	if contender.Disqualified {
		return nil
	}

	return func(yield func(Effect) bool) {
		yield(EffectCalculateProblemValue{CompClassID: contender.CompClassID, ProblemID: event.ProblemID})
		yield(EffectScoreContender{ContenderID: contender.ID})
	}
}

func (e *DefaultScoreEngine) HandleProblemAdded(event domain.ProblemAddedEvent) iter.Seq[Effect] {
	problem := Problem{
		ID:           event.ProblemID,
		ProblemValue: event.ProblemValue,
	}

	e.store.SaveProblem(problem)

	return nil
}

func (e *DefaultScoreEngine) HandleProblemUpdated(event domain.ProblemUpdatedEvent) iter.Seq[Effect] {
	problem := Problem{
		ID:           event.ProblemID,
		ProblemValue: event.ProblemValue,
	}

	e.store.SaveProblem(problem)

	return func(yield func(Effect) bool) {
		for _, compClassID := range e.store.GetCompClassIDs() {
			if !yield(EffectCalculateProblemValue{CompClassID: compClassID, ProblemID: event.ProblemID}) {
				return
			}
		}
	}
}

func (e *DefaultScoreEngine) GetDirtyScores() []domain.Score {
	return e.store.GetDirtyScores()
}

func (e *DefaultScoreEngine) GetDirtyProblemValues() []ProblemValue {
	return e.store.GetDirtyProblemValues()
}

func (e *DefaultScoreEngine) CalculateProblemValue(compClassID domain.CompClassID, problemID domain.ProblemID) iter.Seq[Effect] {
	rules := e.store.GetRules()
	var affectedContenders []domain.ContenderID

	if !rules.UsePoints {
		return nil
	}

	problem, found := e.store.GetProblem(problemID)
	if !found {
		return nil
	}

	value := ProblemValue{
		ProblemID:    problem.ID,
		CompClassID:  compClassID,
		ProblemValue: problem.ProblemValue,
	}

	if rules.PooledPoints {
		numZone1 := 0
		numZone2 := 0
		numTop := 0
		numFlash := 0

		for tick := range e.store.GetTicksByProblem(compClassID, problemID) {
			if tick.Zone1 {
				numZone1++
			}

			if tick.Zone2 {
				numZone2++
			}

			if tick.Top {
				numTop++

				if tick.AttemptsTop == 1 {
					numFlash++
				}
			}

			affectedContenders = append(affectedContenders, tick.ContenderID)
		}

		value = ProblemValue{
			ProblemID:   problem.ID,
			CompClassID: compClassID,

			ProblemValue: domain.ProblemValue{
				PointsZone1: problem.PointsZone1 / max(1, numZone1),
				PointsZone2: problem.PointsZone2 / max(1, numZone2),
				PointsTop:   problem.PointsTop / max(1, numTop),
				FlashBonus:  problem.FlashBonus / max(1, numFlash),
			},
		}
	}

	e.store.SaveProblemValue(compClassID, problemID, value)

	return func(yield func(Effect) bool) {
		for _, contenderID := range affectedContenders {
			if !yield(EffectScoreContender{ContenderID: contenderID}) {
				return
			}
		}
	}
}

func (e *DefaultScoreEngine) ScoreContender(contenderID domain.ContenderID) iter.Seq[Effect] {
	contender, found := e.store.GetContender(contenderID)
	if !found {
		return nil
	}

	oldScore := contender.Score

	if contender.Disqualified {
		contender.Score = 0
	} else {
		ticks := e.store.GetTicksByContender(contender.ID)

		var scoredTicks iter.Seq[Tick] = func(yield func(Tick) bool) {
			for tick := range ticks {
				value, found := e.store.GetProblemValue(contender.CompClassID, tick.ProblemID)
				if !found {
					continue
				}

				tick.Score(value)

				if !yield(tick) {
					return
				}
			}
		}

		problemLimit := e.store.GetRules().QualifyingProblems
		if !e.store.GetRules().UsePoints {
			problemLimit = 0
		}

		scorer := Scorer{
			ProblemLimit: problemLimit,
		}

		contender.Score = scorer.CalculateScore(Points(scoredTicks))
	}

	if contender.Score == oldScore {
		return nil
	}

	e.store.SaveContender(contender)

	return func(yield func(Effect) bool) {
		yield(EffectRankClass{CompClassID: contender.CompClassID})
	}
}

func (e *DefaultScoreEngine) RankCompClass(compClassID domain.CompClassID) iter.Seq[Effect] {
	ranker := NewBasicRanker(e.store.GetRules().Finalists)

	scores := ranker.RankContenders(e.store.GetContendersByCompClass(compClassID))

	for score := range slices.Values(scores) {
		e.store.SaveScore(score)
	}

	return nil
}
