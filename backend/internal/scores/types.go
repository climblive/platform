package scores

import "github.com/climblive/platform/backend/internal/domain"

type Contender struct {
	ID                  domain.ResourceID
	CompClassID         domain.ResourceID
	Disqualified        bool
	WithdrawnFromFinals bool
	Ticks               map[domain.ResourceID]*Tick
	Score               int
}

func (c Contender) Compare(other Contender) int {
	if c.Score == other.Score {
		return c.ID - other.ID
	}

	return other.Score - c.Score
}

type Tick struct {
	ProblemID    domain.ResourceID
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
	ID         domain.ResourceID
	PointsTop  int
	PointsZone int
	FlashBonus int
}
