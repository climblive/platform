package scores

import (
	"cmp"

	"github.com/climblive/platform/backend/internal/domain"
)

type Score struct {
	Points         int
	Tops           int
	AttemptsTops   int
	Zone1s         int
	AttemptsZone1s int
	Zone2s         int
	AttemptsZone2s int
}

type Contender struct {
	ID                  domain.ContenderID
	CompClassID         domain.CompClassID
	Disqualified        bool
	WithdrawnFromFinals bool

	Score
}

func (c Contender) Compare(other Contender) int {
	if result := cmp.Compare(other.Points, c.Points); result != 0 {
		return result
	}

	if result := cmp.Compare(other.Tops, c.Tops); result != 0 {
		return result
	}

	if result := cmp.Compare(c.AttemptsTops, other.AttemptsTops); result != 0 {
		return result
	}

	if result := cmp.Compare(other.Zone2s, c.Zone2s); result != 0 {
		return result
	}

	if result := cmp.Compare(c.AttemptsZone2s, other.AttemptsZone2s); result != 0 {
		return result
	}

	if result := cmp.Compare(other.Zone1s, c.Zone1s); result != 0 {
		return result
	}

	if result := cmp.Compare(c.AttemptsZone1s, other.AttemptsZone1s); result != 0 {
		return result
	}

	return cmp.Compare(c.ID, other.ID)
}

type Tick struct {
	Revision      int
	ContenderID   domain.ContenderID
	ProblemID     domain.ProblemID
	Zone1         bool
	AttemptsZone1 int
	Zone2         bool
	AttemptsZone2 int
	Top           bool
	AttemptsTop   int
}

func (t Tick) TurnIntoZone1() Tick {
	if !t.Zone1 {
		t.AttemptsZone1 += 1
		t.AttemptsZone2 += 1
		t.AttemptsTop += 1
	}

	t.Zone1 = true
	t.Zone2 = false
	t.Top = false

	return t
}

func (t Tick) TurnIntoZone2() Tick {
	if !t.Zone1 {
		t.AttemptsZone1 += 1
	}

	if !t.Zone1 || !t.Zone2 {
		t.AttemptsZone2 += 1
		t.AttemptsTop += 1
	}

	t.Zone1 = true
	t.Zone2 = true
	t.Top = false

	return t
}

func (t Tick) TurnIntoRedpoint() Tick {
	if !t.Zone1 {
		t.AttemptsZone1 += 1
	}

	if !t.Zone1 || !t.Zone2 {
		t.AttemptsZone2 += 1
	}

	if !t.Zone1 || !t.Zone2 || !t.Top {
		t.AttemptsTop += 1
	}

	if t.AttemptsTop == 1 {
		t.AttemptsZone1 += 1
		t.AttemptsZone2 += 1
		t.AttemptsTop += 1
	}

	t.Zone1 = true
	t.Zone2 = true
	t.Top = true

	return t
}

func (t Tick) TurnIntoFlash() Tick {
	t.Zone1 = true
	t.AttemptsZone1 = 1

	t.Zone2 = true
	t.AttemptsZone2 = 1

	t.Top = true
	t.AttemptsTop = 1

	return t
}

type Problem struct {
	ID domain.ProblemID

	domain.ProblemValue

	Zone1Enabled bool
	Zone2Enabled bool
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

func (c TickPool) CalculatePooledProblemValue(value domain.ProblemValue) domain.ProblemValue {
	weightedValue := func(value int, divisor int) int {
		if divisor == 0 {
			return value
		}

		if value == 0 {
			return 0
		}

		return max(1, value/divisor)
	}

	return domain.ProblemValue{
		PointsZone1: weightedValue(value.PointsZone1, c.Zone1),
		PointsZone2: weightedValue(value.PointsZone2, c.Zone2),
		PointsTop:   weightedValue(value.PointsTop, c.Top),
		FlashBonus:  weightedValue(value.FlashBonus, c.Flash),
	}
}
