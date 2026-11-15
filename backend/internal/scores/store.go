package scores

import (
	"iter"
	"maps"
	"slices"

	"github.com/climblive/platform/backend/internal/domain"
)

type MemoryStore struct {
	rules       Rules
	problems    map[domain.ProblemID]Problem
	pointValues *DiffMap[struct {
		CompetitorID domain.CompetitorID
		ProblemID    domain.ProblemID
	}, domain.PointValue]
	competitors map[domain.CompetitorID]Competitor
	ticks       map[domain.CompetitorID]map[domain.ProblemID]Tick
	results     *DiffMap[domain.CompetitorID, domain.Result]
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		problems: make(map[domain.ProblemID]Problem),
		pointValues: NewDiffMap[struct {
			CompetitorID domain.CompetitorID
			ProblemID    domain.ProblemID
		}](ComparePointValue),
		competitors: make(map[domain.CompetitorID]Competitor),
		ticks:       make(map[domain.CompetitorID]map[domain.ProblemID]Tick),
		results:     NewDiffMap[domain.CompetitorID](CompareResult),
		rules:       Rules{},
	}
}

func (s *MemoryStore) GetRules() Rules {
	return s.rules
}

func (s *MemoryStore) SaveRules(rules Rules) {
	s.rules = rules
}

func (s *MemoryStore) GetCompetitor(competitorID domain.CompetitorID) (Competitor, bool) {
	competitor, ok := s.competitors[competitorID]
	return competitor, ok
}

func (s *MemoryStore) GetCompetitorsByCategory(categoryID domain.CategoryID) iter.Seq[Competitor] {
	return func(yield func(Competitor) bool) {
		for _, competitor := range s.competitors {
			if competitor.CategoryID != categoryID {
				continue
			}

			if !yield(competitor) {
				return
			}
		}
	}
}

func (s *MemoryStore) GetAllCompetitors() iter.Seq[Competitor] {
	return maps.Values(s.competitors)
}

func (s *MemoryStore) SaveCompetitor(competitor Competitor) {
	s.competitors[competitor.ID] = competitor
}

func (s *MemoryStore) GetCategoryIDs() []domain.CategoryID {
	categoryIDs := make(map[domain.CategoryID]struct{})

	for competitor := range maps.Values(s.competitors) {
		categoryIDs[competitor.CategoryID] = struct{}{}
	}

	return slices.Collect(maps.Keys(categoryIDs))
}

func (s *MemoryStore) GetTicksByCompetitor(competitorID domain.CompetitorID) iter.Seq[Tick] {
	return maps.Values(s.ticks[competitorID])
}

func (s *MemoryStore) GetTick(competitorID domain.CompetitorID, problemID domain.ProblemID) (Tick, bool) {
	competitorTicks := s.ticks[competitorID]
	if competitorTicks == nil {
		return Tick{}, false
	}

	tick, ok := competitorTicks[problemID]
	return tick, ok
}

func (s *MemoryStore) GetTicksByProblem(categoryID domain.CategoryID, problemID domain.ProblemID) iter.Seq[Tick] {
	return func(yield func(Tick) bool) {
		for _, competitor := range s.competitors {
			if competitor.CategoryID != categoryID {
				continue
			}

			competitorTicks := s.ticks[competitor.ID]

			if competitorTicks == nil {
				continue
			}

			tick, found := competitorTicks[problemID]

			if !found {
				continue
			}

			if !yield(tick) {
				return
			}
		}
	}
}

func (s *MemoryStore) SaveTick(competitorID domain.CompetitorID, tick Tick) {
	competitorTicks := s.ticks[competitorID]

	if competitorTicks == nil {
		competitorTicks = make(map[domain.ProblemID]Tick)
		s.ticks[competitorID] = competitorTicks
	}

	competitorTicks[tick.ProblemID] = tick
}

func (s *MemoryStore) DeleteTick(tickID domain.TickID) {
	for _, competitorTicks := range s.ticks {
		for problemID, tick := range competitorTicks {
			if tick.ID == tickID {
				delete(competitorTicks, problemID)
				return
			}
		}
	}
}

func (s *MemoryStore) GetProblem(problemID domain.ProblemID) (Problem, bool) {
	problem, ok := s.problems[problemID]
	return problem, ok
}

func (s *MemoryStore) SaveProblem(problem Problem) {
	s.problems[problem.ID] = problem
}

func (s *MemoryStore) GetAllProblems() iter.Seq[Problem] {
	return maps.Values(s.problems)
}

func (s *MemoryStore) GetPointValue(competitorID domain.CompetitorID, problemID domain.ProblemID) (domain.PointValue, bool) {
	key := struct {
		CompetitorID domain.CompetitorID
		ProblemID    domain.ProblemID
	}{competitorID, problemID}

	value, ok := s.pointValues.Get(key)

	return value, ok
}

func (s *MemoryStore) SavePointValue(competitorID domain.CompetitorID, problemID domain.ProblemID, value domain.PointValue) {
	key := struct {
		CompetitorID domain.CompetitorID
		ProblemID    domain.ProblemID
	}{competitorID, problemID}

	s.pointValues.Set(key, value)
}

func (s *MemoryStore) SaveResult(result domain.Result) {
	s.results.Set(result.CompetitorID, result)
}

func (s *MemoryStore) GetDirtyResults() []domain.Result {
	return s.results.Commit()
}

func (s *MemoryStore) GetDirtyPointValues() []domain.PointValue {
	return s.pointValues.Commit()
}
