package scores_test

import (
	"testing"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/scores"
	"github.com/stretchr/testify/assert"
)

func TestScoreTick(t *testing.T) {
	problem := scores.ProblemValue{
		ProblemValue: domain.ProblemValue{
			PointsTop:   100,
			PointsZone1: 50,
			PointsZone2: 75,
			FlashBonus:  10,
		},
	}

	t.Run("NoAttempts", func(t *testing.T) {
		previousPoints := 1_000

		tick := scores.Tick{
			Top:           false,
			AttemptsTop:   0,
			Zone1:         false,
			AttemptsZone1: 0,
			Zone2:         false,
			AttemptsZone2: 0,
			Points:        previousPoints,
		}

		tick.Score(problem)

		assert.Equal(t, 0, tick.Points)
	})

	t.Run("SingleAttemptNoLuck", func(t *testing.T) {
		tick := scores.Tick{
			Top:           false,
			AttemptsTop:   1,
			Zone1:         false,
			AttemptsZone1: 1,
			Zone2:         false,
			AttemptsZone2: 1,
		}

		tick.Score(problem)

		assert.Equal(t, 0, tick.Points)
	})

	t.Run("Flash", func(t *testing.T) {
		tick := scores.Tick{
			Top:           true,
			AttemptsTop:   1,
			Zone1:         true,
			AttemptsZone1: 1,
			Zone2:         true,
			AttemptsZone2: 1,
		}

		tick.Score(problem)

		assert.Equal(t, 110, tick.Points)
	})

	t.Run("TopWithSeveralAttempts", func(t *testing.T) {
		tick := scores.Tick{
			Top:           true,
			AttemptsTop:   999,
			Zone1:         true,
			AttemptsZone1: 999,
			Zone2:         true,
			AttemptsZone2: 999,
		}

		tick.Score(problem)

		assert.Equal(t, 100, tick.Points)
	})

	t.Run("Zone1WithSeveralAttempts", func(t *testing.T) {
		tick := scores.Tick{
			Top:           false,
			AttemptsTop:   999,
			Zone1:         true,
			AttemptsZone1: 999,
			Zone2:         false,
			AttemptsZone2: 999,
		}

		tick.Score(problem)

		assert.Equal(t, 50, tick.Points)
	})
	t.Run("Zone2WithSeveralAttempts", func(t *testing.T) {
		tick := scores.Tick{
			Top:           false,
			AttemptsTop:   999,
			Zone1:         true,
			AttemptsZone1: 999,
			Zone2:         true,
			AttemptsZone2: 999,
		}

		tick.Score(problem)

		assert.Equal(t, 75, tick.Points)
	})
}

func TestCompareContender(t *testing.T) {
	t.Run("ByScore", func(t *testing.T) {
		c1 := scores.Contender{
			Score: 200,
		}

		c2 := scores.Contender{
			Score: 100,
		}

		assert.Less(t, c1.Compare(c2), 0)
		assert.Greater(t, c2.Compare(c1), 0)
	})

	t.Run("TieBreak", func(t *testing.T) {
		c1 := scores.Contender{
			ID:    1,
			Score: 100,
		}

		c2 := scores.Contender{
			ID:    2,
			Score: 100,
		}

		assert.Less(t, c1.Compare(c2), 0)
		assert.Greater(t, c2.Compare(c1), 0)
	})
}
