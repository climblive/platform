package scores

import (
	"bytes"

	"github.com/climblive/platform/backend/internal/domain"
)

type Contender struct {
	ID                  domain.ContenderID
	CompClassID         domain.CompClassID
	Disqualified        bool
	WithdrawnFromFinals bool
	Score               int
}

func (c Contender) Compare(other Contender) int {
	if c.Score == other.Score {
		return bytes.Compare(c.ID[:], other.ID[:])
	}

	return other.Score - c.Score
}

type Tick struct {
	ProblemID     domain.ProblemID
	Zone1         bool
	AttemptsZone1 int
	Zone2         bool
	AttemptsZone2 int
	Top           bool
	AttemptsTop   int
	Points        int
}

func (t *Tick) Score(problem Problem) {
	t.Points = 0

	if t.Zone1 {
		t.Points = problem.PointsZone1
	}

	if t.Zone2 {
		t.Points = problem.PointsZone2
	}

	if t.Top {
		t.Points = problem.PointsTop
	}

	if t.Top && t.AttemptsTop == 1 {
		t.Points += problem.FlashBonus
	}
}

type Problem struct {
	ID          domain.ProblemID
	PointsZone1 int
	PointsZone2 int
	PointsTop   int
	FlashBonus  int
}
