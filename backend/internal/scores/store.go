package scores

import (
	"iter"
	"maps"
	"slices"

	"github.com/climblive/platform/backend/internal/domain"
)

type MemoryStore struct {
	problems   map[domain.ProblemID]Problem
	contenders map[domain.ContenderID]Contender
	ticks      map[domain.ContenderID][]Tick
	scores     *DiffMap[domain.ContenderID, domain.Score]
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		problems:   make(map[domain.ProblemID]Problem),
		contenders: make(map[domain.ContenderID]Contender),
		ticks:      make(map[domain.ContenderID][]Tick),
		scores:     NewDiffMap[domain.ContenderID](CompareScore),
	}
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

func (s *MemoryStore) GetTicks(contenderID domain.ContenderID) iter.Seq[Tick] {
	return slices.Values(s.ticks[contenderID])
}

func (s *MemoryStore) GetTicksByProblem(problemID domain.ProblemID) iter.Seq[Tick] {
	return func(yield func(Tick) bool) {
		for _, contenderTicks := range s.ticks {
			for _, tick := range contenderTicks {
				if tick.ProblemID != problemID {
					continue
				}

				if !yield(tick) {
					return
				}
			}
		}
	}
}

func (s *MemoryStore) SaveTick(contenderID domain.ContenderID, tick Tick) {
	cmp := func(t Tick) bool {
		return t.ProblemID == tick.ProblemID
	}

	contenderTicks := s.ticks[contenderID]

	if i := slices.IndexFunc(contenderTicks, cmp); i != -1 {
		contenderTicks[i] = tick
	} else {
		s.ticks[contenderID] = append(contenderTicks, tick)
	}
}

func (s *MemoryStore) DeleteTick(contenderID domain.ContenderID, problemID domain.ProblemID) {
	cmp := func(t Tick) bool {
		return t.ProblemID == problemID
	}

	s.ticks[contenderID] = slices.DeleteFunc(s.ticks[contenderID], cmp)
}

func (s *MemoryStore) GetProblem(problemID domain.ProblemID) (Problem, bool) {
	problem, ok := s.problems[problemID]
	return problem, ok
}

func (s *MemoryStore) SaveProblem(problem Problem) {
	s.problems[problem.ID] = problem
}

func (s *MemoryStore) SaveScore(score domain.Score) {
	s.scores.Set(score.ContenderID, score)
}

func (s *MemoryStore) GetDirtyScores() []domain.Score {
	return s.scores.Commit()
}
