package scores

import (
	"iter"
	"maps"
	"slices"

	"github.com/climblive/platform/backend/internal/domain"
)

type MemoryStore struct {
	rules         Rules
	problems      map[domain.ProblemID]Problem
	problemValues *DiffMap[struct {
		CompClassID domain.CompClassID
		ProblemID   domain.ProblemID
	}, ProblemValue]
	contenders map[domain.ContenderID]Contender
	ticks      map[domain.ContenderID]map[domain.ProblemID]Tick
	scores     *DiffMap[domain.ContenderID, domain.Score]
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		problems: make(map[domain.ProblemID]Problem),
		problemValues: NewDiffMap[struct {
			CompClassID domain.CompClassID
			ProblemID   domain.ProblemID
		}](CompareProblemValue),
		contenders: make(map[domain.ContenderID]Contender),
		ticks:      make(map[domain.ContenderID]map[domain.ProblemID]Tick),
		scores:     NewDiffMap[domain.ContenderID](CompareScore),
	}
}

func (s *MemoryStore) GetRules() Rules {
	return s.rules
}

func (s *MemoryStore) SaveRules(rules Rules) {
	s.rules = rules
}

func (s *MemoryStore) GetContender(contenderID domain.ContenderID) (Contender, bool) {
	contender, ok := s.contenders[contenderID]
	return contender, ok
}

func (s *MemoryStore) GetContendersByCompClass(compClassID domain.CompClassID) iter.Seq[Contender] {
	return func(yield func(Contender) bool) {
		for _, contender := range s.contenders {
			if contender.CompClassID != compClassID {
				continue
			}

			if !yield(contender) {
				return
			}
		}
	}
}

func (s *MemoryStore) GetAllContenders() iter.Seq[Contender] {
	return maps.Values(s.contenders)
}

func (s *MemoryStore) SaveContender(contender Contender) {
	s.contenders[contender.ID] = contender
}

func (s *MemoryStore) GetCompClassIDs() []domain.CompClassID {
	compClassIDs := make(map[domain.CompClassID]struct{})

	for contender := range maps.Values(s.contenders) {
		compClassIDs[contender.CompClassID] = struct{}{}
	}

	return slices.Collect(maps.Keys(compClassIDs))
}

func (s *MemoryStore) GetTicksByContender(contenderID domain.ContenderID) iter.Seq[Tick] {
	return maps.Values(s.ticks[contenderID])
}

func (s *MemoryStore) GetTicksByProblem(compClassID domain.CompClassID, problemID domain.ProblemID) iter.Seq[Tick] {
	return func(yield func(Tick) bool) {
		for _, contender := range s.contenders {
			if contender.CompClassID != compClassID {
				continue
			}

			contenderTicks := s.ticks[contender.ID]

			if contenderTicks == nil {
				continue
			}

			tick, found := contenderTicks[problemID]

			if !found {
				continue
			}

			if !yield(tick) {
				return
			}
		}
	}
}

func (s *MemoryStore) SaveTick(contenderID domain.ContenderID, tick Tick) {
	contenderTicks := s.ticks[contenderID]

	if contenderTicks == nil {
		contenderTicks = make(map[domain.ProblemID]Tick)
		s.ticks[contenderID] = contenderTicks
	}

	contenderTicks[tick.ProblemID] = tick
}

func (s *MemoryStore) DeleteTick(contenderID domain.ContenderID, problemID domain.ProblemID) {
	contenderTicks := s.ticks[contenderID]

	delete(contenderTicks, problemID)
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

func (s *MemoryStore) GetProblemValue(compClassID domain.CompClassID, problemID domain.ProblemID) (ProblemValue, bool) {
	key := struct {
		CompClassID domain.CompClassID
		ProblemID   domain.ProblemID
	}{compClassID, problemID}

	value, ok := s.problemValues.Get(key)

	return value, ok
}

func (s *MemoryStore) SaveProblemValue(compClassID domain.CompClassID, problemID domain.ProblemID, value ProblemValue) {
	key := struct {
		CompClassID domain.CompClassID
		ProblemID   domain.ProblemID
	}{compClassID, problemID}

	s.problemValues.Set(key, value)
}

func (s *MemoryStore) SaveScore(score domain.Score) {
	s.scores.Set(score.ContenderID, score)
}

func (s *MemoryStore) GetDirtyScores() []domain.Score {
	return s.scores.Commit()
}

func (s *MemoryStore) GetDirtyProblemValues() []ProblemValue {
	return s.problemValues.Commit()
}
