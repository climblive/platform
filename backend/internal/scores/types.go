package scores

import "github.com/climblive/platform/backend/internal/domain"

type Contender struct {
	ID                  domain.ContenderID
	CompClassID         domain.CompClassID
	Disqualified        bool
	WithdrawnFromFinals bool
	Score               int
}

func (c Contender) Compare(other Contender) int {
	if c.Score == other.Score {
		return int(c.ID) - int(other.ID)
	}

	return other.Score - c.Score
}

type Tick struct {
	ContenderID   domain.ContenderID
	ProblemID     domain.ProblemID
	Zone1         bool
	AttemptsZone1 int
	Zone2         bool
	AttemptsZone2 int
	Top           bool
	AttemptsTop   int
	Points        int
}

func (t *Tick) Score(value domain.ProblemValue) {
	t.Points = 0

	if t.Zone1 {
		t.Points = value.PointsZone1
	}

	if t.Zone2 {
		t.Points = value.PointsZone2
	}

	if t.Top {
		t.Points = value.PointsTop
	}

	if t.Top && t.AttemptsTop == 1 {
		t.Points += value.FlashBonus
	}
}

type Problem struct {
	ID domain.ProblemID

	domain.ProblemValue
}
