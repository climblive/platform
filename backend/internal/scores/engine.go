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
	EffectTypeCalculatePointValues EffectType = iota
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

type EffectCalculatePointValues struct {
	CompClassID domain.CompClassID
	ProblemID   domain.ProblemID
}

func (e EffectCalculatePointValues) Encode() EncodedEffect {
	var data EncodedEffect
	data[0] = byte(EffectTypeCalculatePointValues)
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
	GetTick(domain.ContenderID, domain.ProblemID) (Tick, bool)
	SaveTick(domain.ContenderID, Tick)
	DeleteTick(domain.ContenderID, domain.ProblemID)
	GetTicksByProblem(domain.CompClassID, domain.ProblemID) iter.Seq[Tick]

	GetProblem(domain.ProblemID) (Problem, bool)
	SaveProblem(Problem)
	GetAllProblems() iter.Seq[Problem]

	GetPointValue(domain.ContenderID, domain.ProblemID) (PointValue, bool)
	SavePointValue(domain.ContenderID, domain.ProblemID, PointValue)
	GetDirtyPointValues() []PointValue

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
				if !yield(EffectCalculatePointValues{CompClassID: compClassID, ProblemID: problem.ID}) {
					return
				}
			}

			for contender := range e.store.GetContendersByCompClass(compClassID) {
				if !yield(EffectScoreContender{ContenderID: contender.ID}) {
					return
				}
			}
		}

		for _, compClassID := range e.store.GetCompClassIDs() {
			if !yield(EffectRankClass{CompClassID: compClassID}) {
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
	}

	e.store.SaveRules(rules)

	return e.Start()
}

func (e *DefaultScoreEngine) HandleContenderEntered(event domain.ContenderEnteredEvent) iter.Seq[Effect] {
	contender := Contender{
		ID:                  event.ContenderID,
		CompClassID:         event.CompClassID,
		Disqualified:        false,
		WithdrawnFromFinals: false,
		Score:               0,
	}

	e.store.SaveContender(contender)

	return func(yield func(Effect) bool) {
		for problem := range e.store.GetAllProblems() {
			if !yield(EffectCalculatePointValues{CompClassID: contender.CompClassID, ProblemID: problem.ID}) {
				return
			}
		}

		if !yield(EffectRankClass{CompClassID: contender.CompClassID}) {
			return
		}
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
		for problem := range e.store.GetAllProblems() {
			if !yield(EffectCalculatePointValues{CompClassID: oldCompClassID, ProblemID: problem.ID}) {
				return
			}

			if !yield(EffectCalculatePointValues{CompClassID: event.CompClassID, ProblemID: problem.ID}) {
				return
			}
		}

		if !yield(EffectRankClass{CompClassID: oldCompClassID}) {
			return
		}

		if !yield(EffectRankClass{CompClassID: event.CompClassID}) {
			return
		}
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
		if !yield(EffectRankClass{CompClassID: contender.CompClassID}) {
			return
		}
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
		if !yield(EffectRankClass{CompClassID: contender.CompClassID}) {
			return
		}
	}
}

func (e *DefaultScoreEngine) HandleContenderDisqualified(event domain.ContenderDisqualifiedEvent) iter.Seq[Effect] {
	contender, found := e.store.GetContender(event.ContenderID)
	if !found {
		return nil
	}

	contender.Disqualified = true

	e.store.SaveContender(contender)

	ticks := e.store.GetTicksByContender(contender.ID)

	return func(yield func(Effect) bool) {
		for tick := range ticks {
			if !yield(EffectCalculatePointValues{CompClassID: contender.CompClassID, ProblemID: tick.ProblemID}) {
				return
			}
		}

		if !yield(EffectScoreContender{ContenderID: contender.ID}) {
			return
		}

		if !yield(EffectRankClass{CompClassID: contender.CompClassID}) {
			return
		}
	}
}

func (e *DefaultScoreEngine) HandleContenderRequalified(event domain.ContenderRequalifiedEvent) iter.Seq[Effect] {
	contender, found := e.store.GetContender(event.ContenderID)
	if !found {
		return nil
	}

	contender.Disqualified = false

	e.store.SaveContender(contender)

	ticks := e.store.GetTicksByContender(contender.ID)

	return func(yield func(Effect) bool) {
		for tick := range ticks {
			if !yield(EffectCalculatePointValues{CompClassID: contender.CompClassID, ProblemID: tick.ProblemID}) {
				return
			}
		}

		if !yield(EffectScoreContender{ContenderID: contender.ID}) {
			return
		}

		if !yield(EffectRankClass{CompClassID: contender.CompClassID}) {
			return
		}
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
		if !yield(EffectCalculatePointValues{CompClassID: contender.CompClassID, ProblemID: event.ProblemID}) {
			return
		}

		if !yield(EffectScoreContender{ContenderID: contender.ID}) {
			return
		}
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
		if !yield(EffectCalculatePointValues{CompClassID: contender.CompClassID, ProblemID: event.ProblemID}) {
			return
		}

		if !yield(EffectScoreContender{ContenderID: contender.ID}) {
			return
		}
	}
}

func (e *DefaultScoreEngine) HandleProblemAdded(event domain.ProblemAddedEvent) iter.Seq[Effect] {
	problem := Problem{
		ID:           event.ProblemID,
		ProblemValue: event.ProblemValue,
	}

	e.store.SaveProblem(problem)

	return func(yield func(Effect) bool) {
		for _, compClassID := range e.store.GetCompClassIDs() {
			if !yield(EffectCalculatePointValues{CompClassID: compClassID, ProblemID: event.ProblemID}) {
				return
			}
		}
	}
}

func (e *DefaultScoreEngine) HandleProblemUpdated(event domain.ProblemUpdatedEvent) iter.Seq[Effect] {
	problem := Problem{
		ID:           event.ProblemID,
		ProblemValue: event.ProblemValue,
	}

	e.store.SaveProblem(problem)

	return func(yield func(Effect) bool) {
		for _, compClassID := range e.store.GetCompClassIDs() {
			if !yield(EffectCalculatePointValues{CompClassID: compClassID, ProblemID: event.ProblemID}) {
				return
			}
		}
	}
}

func (e *DefaultScoreEngine) GetDirtyScores() []domain.Score {
	return e.store.GetDirtyScores()
}

func (e *DefaultScoreEngine) GetDirtyPointValues() []PointValue {
	return e.store.GetDirtyPointValues()
}

func pointMaximum(value domain.ProblemValue, tick *Tick) int {
	maximum := value.PointsTop

	if tick == nil {
		maximum += value.FlashBonus
	}

	return maximum
}

func pointCurrent(value domain.ProblemValue, tick *Tick) int {
	if tick == nil {
		return 0
	}

	current := 0

	if tick.Zone1 {
		current = value.PointsZone1
	}

	if tick.Zone2 {
		current = value.PointsZone2
	}

	if tick.Top {
		current = value.PointsTop

		if tick.AttemptsTop == 1 {
			current += value.FlashBonus
		}
	}

	return current
}

func (e *DefaultScoreEngine) CalculatePointValues(compClassID domain.CompClassID, problemID domain.ProblemID) iter.Seq[Effect] {
	rules := e.store.GetRules()

	if !rules.UsePoints {
		return nil
	}

	problem, found := e.store.GetProblem(problemID)
	if !found {
		return nil
	}

	value := problem.ProblemValue

	if rules.PooledPoints {
		numZone1 := 0
		numZone2 := 0
		numTop := 0
		numFlash := 0

		for tick := range e.store.GetTicksByProblem(compClassID, problemID) {
			contender, found := e.store.GetContender(tick.ContenderID)
			if !found {
				continue
			}

			if contender.Disqualified {
				continue
			}

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
		}

		value = domain.ProblemValue{
			PointsZone1: problem.PointsZone1 / max(1, numZone1),
			PointsZone2: problem.PointsZone2 / max(1, numZone2),
			PointsTop:   problem.PointsTop / max(1, numTop),
			FlashBonus:  problem.FlashBonus / max(1, numFlash),
		}
	}

	affectedContenders := make([]domain.ContenderID, 0)

	for contender := range e.store.GetContendersByCompClass(compClassID) {
		tick, found := e.store.GetTick(contender.ID, problemID)
		var tickPtr *Tick
		if found {
			tickPtr = &tick
		}

		pointValue := PointValue{
			ContenderID: contender.ID,
			ProblemID:   problemID,
			Current:     pointCurrent(value, tickPtr),
			Maximum:     pointMaximum(value, tickPtr),
		}

		oldValue, found := e.store.GetPointValue(contender.ID, problemID)
		e.store.SavePointValue(contender.ID, problemID, pointValue)

		if !found || !ComparePointValue(oldValue, pointValue) {
			affectedContenders = append(affectedContenders, contender.ID)
		}
	}

	if len(affectedContenders) == 0 {
		return nil
	}

	return func(yield func(Effect) bool) {
		for _, contenderID := range affectedContenders {
			if !yield(EffectScoreContender{ContenderID: contenderID}) {
				return
			}
		}
	}
}

func (e *DefaultScoreEngine) CalculateProblemValue(compClassID domain.CompClassID, problemID domain.ProblemID) iter.Seq[Effect] {
	return e.CalculatePointValues(compClassID, problemID)
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

		var pointValues iter.Seq[PointValue] = func(yield func(PointValue) bool) {
			for tick := range ticks {
				value, found := e.store.GetPointValue(contender.ID, tick.ProblemID)
				if !found {
					continue
				}

				if !yield(value) {
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

		contender.Score = scorer.CalculateScore(Points(pointValues))
	}

	if contender.Score == oldScore {
		return nil
	}

	e.store.SaveContender(contender)

	return func(yield func(Effect) bool) {
		yield(EffectRankClass{CompClassID: contender.CompClassID})
	}
}

func (e *DefaultScoreEngine) RankCompClass(compClassID domain.CompClassID) {
	ranker := NewBasicRanker(e.store.GetRules().Finalists)

	scores := ranker.RankContenders(e.store.GetContendersByCompClass(compClassID))

	for score := range slices.Values(scores) {
		e.store.SaveScore(score)
	}
}
