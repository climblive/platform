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
}

type Problem struct {
	ID domain.ProblemID

	domain.ProblemValue
}
