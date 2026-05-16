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

type TickPool struct {
	Zone1 int
	Zone2 int
	Top   int
	Flash int
}

func (c TickPool) Add(tick Tick) TickPool {
	if tick.Zone1 {
		c.Zone1++
	}

	if tick.Zone2 {
		c.Zone2++
	}

	if tick.Top {
		c.Top++

		if tick.AttemptsTop == 1 {
			c.Flash++
		}
	}

	return c
}

func (c TickPool) Sub(tick Tick) TickPool {
	if tick.Zone1 {
		c.Zone1--
	}

	if tick.Zone2 {
		c.Zone2--
	}

	if tick.Top {
		c.Top--

		if tick.AttemptsTop == 1 {
			c.Flash--
		}
	}

	return c
}

func (c TickPool) CalculateProblemValue(value domain.ProblemValue) domain.ProblemValue {
	return domain.ProblemValue{
		PointsZone1: value.PointsZone1 / max(1, c.Zone1),
		PointsZone2: value.PointsZone2 / max(1, c.Zone2),
		PointsTop:   value.PointsTop / max(1, c.Top),
		FlashBonus:  value.FlashBonus / max(1, c.Flash),
	}
}
