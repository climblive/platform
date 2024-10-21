package scores

import "github.com/climblive/platform/backend/internal/domain"

type Contender struct {
	ID                  domain.ContenderID
	CompClassID         domain.CompClassID
	Disqualified        bool
	WithdrawnFromFinals bool
	Ticks               map[domain.ProblemID]*Tick
	Score               int
}

func (c Contender) Compare(other Contender) int {
	if c.Score == other.Score {
		return int(c.ID) - int(other.ID)
	}

	return other.Score - c.Score
}

type Tick struct {
	ProblemID    domain.ProblemID
	Top          bool
	AttemptsTop  int
	Zone         bool
	AttemptsZone int
	Points       int
}

func (t *Tick) Score(problem Problem) {
	if t.Zone {
		t.Points = problem.PointsZone
	}

	if t.Top {
		t.Points = problem.PointsTop
	}

	if t.AttemptsTop == 1 {
		t.Points += problem.FlashBonus
	}
}

type Problem struct {
	ID         domain.ProblemID
	PointsTop  int
	PointsZone int
	FlashBonus int
}
