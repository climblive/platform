package scores_test

import (
	"testing"

	"github.com/climblive/platform/backend/internal/scores"
	"github.com/stretchr/testify/assert"
)

func TestScoreTick(t *testing.T) {
	problem := scores.Problem{
		PointsTop:  100,
		PointsZone: 50,
		FlashBonus: 10,
	}

	t.Run("NoAttempts", func(t *testing.T) {
		previousPoints := 1_000

		tick := scores.Tick{
			Top:          false,
			AttemptsTop:  0,
			Zone:         false,
			AttemptsZone: 0,
			Points:       previousPoints,
		}

		tick.Score(problem)

		assert.Equal(t, 0, tick.Points)
	})

	t.Run("SingleAttemptNoLuck", func(t *testing.T) {
		tick := scores.Tick{
			Top:          false,
			AttemptsTop:  1,
			Zone:         false,
			AttemptsZone: 1,
		}

		tick.Score(problem)

		assert.Equal(t, 0, tick.Points)
	})

	t.Run("Flash", func(t *testing.T) {
		tick := scores.Tick{
			Top:          true,
			AttemptsTop:  1,
			Zone:         true,
			AttemptsZone: 1,
		}

		tick.Score(problem)

		assert.Equal(t, 110, tick.Points)
	})

	t.Run("TopWithSeveralAttempts", func(t *testing.T) {
		tick := scores.Tick{
			Top:          true,
			AttemptsTop:  999,
			Zone:         true,
			AttemptsZone: 999,
		}

		tick.Score(problem)

		assert.Equal(t, 100, tick.Points)
	})

	t.Run("ZoneWithSeveralAttempts", func(t *testing.T) {
		tick := scores.Tick{
			Top:          false,
			AttemptsTop:  999,
			Zone:         true,
			AttemptsZone: 999,
		}

		tick.Score(problem)

		assert.Equal(t, 50, tick.Points)
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
