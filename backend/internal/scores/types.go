package scores

import "github.com/climblive/platform/backend/internal/domain"

type Contender struct {
	ID                  domain.ContenderID
	CompClassID         domain.CompClassID
	Disqualified        bool
	WithdrawnFromFinals bool
	Score               int
	Tops                int
	AttemptsTops        int
	Zones               int
	AttemptsZones       int
}

func (c Contender) Compare(other Contender) int {
	cmp := other.Score - c.Score
	if cmp != 0 {
		return cmp
	}

	cmp = other.Tops - c.Tops
	if cmp != 0 {
		return cmp
	}

	cmp = other.AttemptsTops - c.AttemptsTops
	if cmp != 0 {
		return cmp
	}

	cmp = other.Zones - c.Zones
	if cmp != 0 {
		return cmp
	}

	cmp = other.AttemptsZones - c.AttemptsZones
	if cmp != 0 {
		return cmp
	}

	return int(other.ID) - int(c.ID)
}

type Tick struct {
	ProblemID    domain.ProblemID
	Top          bool
	AttemptsTop  int
	Zone         bool
	AttemptsZone int
	Points       int
}

func (tick Tick) IsFlash() bool {
	return tick.Top && tick.AttemptsTop == 1
}

func (t *Tick) Score(problem Problem) {
	t.Points = 0

	if t.Zone {
		t.Points = problem.PointsZone
	}

	if t.Top {
		t.Points = problem.PointsTop
	}

	if t.Top && t.AttemptsTop == 1 {
		t.Points += problem.FlashBonus
	}
}

type Problem struct {
	ID         domain.ProblemID
	PointsTop  int
	PointsZone int
	FlashBonus int
}
