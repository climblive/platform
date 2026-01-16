package scores

import (
	"iter"
	"slices"

	"github.com/climblive/platform/backend/internal/domain"
)

type Rules struct {
	QualifyingProblems int
	Finalists          int
}

type ScoringRules interface {
	CalculateScore(points iter.Seq[int]) int
}

type Ranker interface {
	RankContenders(contenders iter.Seq[Contender]) []domain.Score
}

type EngineStore interface {
	GetRules() Rules
	SaveRules(Rules)

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
	GetDirtyScores() []domain.Score
}

type DefaultScoreEngine struct {
	ranker Ranker
	rules  ScoringRules
	store  EngineStore
}

func NewDefaultScoreEngine(store EngineStore) *DefaultScoreEngine {
	return &DefaultScoreEngine{
		store: store,
		rules: &HardestProblems{
			Number: store.GetRules().QualifyingProblems,
		},
		ranker: NewBasicRanker(store.GetRules().Finalists),
	}
}

func (e *DefaultScoreEngine) Start() {
	rules := e.store.GetRules()

	e.rules = &HardestProblems{
		Number: rules.QualifyingProblems,
	}
	e.ranker = NewBasicRanker(rules.Finalists)

	for contender := range e.store.GetAllContenders() {
		ticks := e.store.GetTicks(contender.ID)

		var scoredTicks iter.Seq[Tick] = func(yield func(Tick) bool) {
			for tick := range ticks {
				problem, found := e.store.GetProblem(tick.ProblemID)
				if !found {
					continue
				}

				tick.Score(problem)
				e.store.SaveTick(contender.ID, tick)

				yield(tick)
			}
		}

		contender.Score = e.rules.CalculateScore(Points(scoredTicks))

		if contender.Disqualified {
			contender.Score = 0
		}

		e.store.SaveContender(contender)
	}

	e.rankCompClasses(e.store.GetCompClassIDs()...)
}

func (e *DefaultScoreEngine) Stop() {
}

func (e *DefaultScoreEngine) HandleRulesUpdated(event domain.RulesUpdatedEvent) {
	rules := Rules{
		QualifyingProblems: event.QualifyingProblems,
		Finalists:          event.Finalists,
	}

	e.store.SaveRules(rules)

	e.rules = &HardestProblems{
		Number: rules.QualifyingProblems,
	}
	e.ranker = NewBasicRanker(rules.Finalists)

	for contender := range e.store.GetAllContenders() {
		contender.Score = e.rules.CalculateScore(Points(e.store.GetTicks(contender.ID)))
		e.store.SaveContender(contender)
	}

	e.rankCompClasses(e.store.GetCompClassIDs()...)
}

func (e *DefaultScoreEngine) HandleContenderEntered(event domain.ContenderEnteredEvent) {
	contender := Contender{
		ID:          event.ContenderID,
		CompClassID: event.CompClassID,
	}

	e.store.SaveContender(contender)

	e.rankCompClasses(contender.CompClassID)
}

func (e *DefaultScoreEngine) HandleContenderSwitchedClass(event domain.ContenderSwitchedClassEvent) {
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

func (e *DefaultScoreEngine) HandleContenderWithdrewFromFinals(event domain.ContenderWithdrewFromFinalsEvent) {
	contender, found := e.store.GetContender(event.ContenderID)
	if !found {
		return
	}

	contender.WithdrawnFromFinals = true

	e.store.SaveContender(contender)

	e.rankCompClasses(contender.CompClassID)
}

func (e *DefaultScoreEngine) HandleContenderReenteredFinals(event domain.ContenderReenteredFinalsEvent) {
	contender, found := e.store.GetContender(event.ContenderID)
	if !found {
		return
	}

	contender.WithdrawnFromFinals = false

	e.store.SaveContender(contender)

	e.rankCompClasses(contender.CompClassID)
}

func (e *DefaultScoreEngine) HandleContenderDisqualified(event domain.ContenderDisqualifiedEvent) {
	contender, found := e.store.GetContender(event.ContenderID)
	if !found {
		return
	}

	contender.Disqualified = true
	contender.Score = 0

	e.store.SaveContender(contender)

	e.rankCompClasses(contender.CompClassID)
}

func (e *DefaultScoreEngine) HandleContenderRequalified(event domain.ContenderRequalifiedEvent) {
	contender, found := e.store.GetContender(event.ContenderID)
	if !found {
		return
	}

	contender.Disqualified = false
	contender.Score = e.rules.CalculateScore(Points(e.store.GetTicks(contender.ID)))

	e.store.SaveContender(contender)

	e.rankCompClasses(contender.CompClassID)
}

func (e *DefaultScoreEngine) HandleAscentRegistered(event domain.AscentRegisteredEvent) {
	tick := Tick{
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

func (e *DefaultScoreEngine) HandleAscentDeregistered(event domain.AscentDeregisteredEvent) {
	contender, found := e.store.GetContender(event.ContenderID)
	if !found {
		return
	}

	e.store.DeleteTick(event.ContenderID, event.ProblemID)

	if contender.Disqualified {
		return
	}

	contender.Score = e.rules.CalculateScore(Points(e.store.GetTicks(contender.ID)))
	e.store.SaveContender(contender)

	e.rankCompClasses(contender.CompClassID)
}

func (e *DefaultScoreEngine) HandleProblemAdded(event domain.ProblemAddedEvent) {
	problem := Problem{
		ID:          event.ProblemID,
		PointsZone1: event.PointsZone1,
		PointsZone2: event.PointsZone2,
		PointsTop:   event.PointsTop,
		FlashBonus:  event.FlashBonus,
	}

	e.store.SaveProblem(problem)
}

func (e *DefaultScoreEngine) HandleProblemUpdated(event domain.ProblemUpdatedEvent) {
	problem := Problem{
		ID:          event.ProblemID,
		PointsZone1: event.PointsZone1,
		PointsZone2: event.PointsZone2,
		PointsTop:   event.PointsTop,
		FlashBonus:  event.FlashBonus,
	}

	e.store.SaveProblem(problem)

	e.Start()
}

func (e *DefaultScoreEngine) GetDirtyScores() []domain.Score {
	return e.store.GetDirtyScores()
}

func (e *DefaultScoreEngine) rankCompClasses(compClassIDs ...domain.CompClassID) {
	for _, compClassID := range compClassIDs {
		scores := e.ranker.RankContenders(e.store.GetContendersByCompClass(compClassID))

		for score := range slices.Values(scores) {
			e.store.SaveScore(score)
		}
	}
}
