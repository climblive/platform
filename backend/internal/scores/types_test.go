package scores_test

import (
	"testing"

	"github.com/climblive/platform/backend/internal/scores"
	"github.com/stretchr/testify/assert"
)

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

func TestTickPool_Add(t *testing.T) {
	t.Run("AddZone1", func(t *testing.T) {
		pool := scores.TickPool{}

		for range 5 {
			pool = pool.Add(scores.Tick{
				Zone1:         true,
				AttemptsZone1: 1,
			})
		}

		assert.Equal(t, scores.TickPool{
			Zone1: 5,
		}, pool)
	})

	t.Run("AddZone2", func(t *testing.T) {
		pool := scores.TickPool{}

		for range 5 {
			pool = pool.Add(scores.Tick{
				Zone1:         true,
				AttemptsZone1: 1,
				Zone2:         true,
				AttemptsZone2: 1,
			})
		}

		assert.Equal(t, scores.TickPool{
			Zone1: 5,
			Zone2: 5,
		}, pool)
	})

	t.Run("AddTopWithoutFlash", func(t *testing.T) {
		pool := scores.TickPool{}

		for range 5 {
			pool = pool.Add(scores.Tick{
				Zone1:         true,
				AttemptsZone1: 999,
				Zone2:         true,
				AttemptsZone2: 999,
				Top:           true,
				AttemptsTop:   999,
			})
		}

		assert.Equal(t, scores.TickPool{
			Zone1: 5,
			Zone2: 5,
			Top:   5,
		}, pool)
	})

	t.Run("AddTopWithFlash", func(t *testing.T) {
		pool := scores.TickPool{}

		for range 5 {
			pool = pool.Add(scores.Tick{
				Zone1:         true,
				AttemptsZone1: 1,
				Zone2:         true,
				AttemptsZone2: 1,
				Top:           true,
				AttemptsTop:   1,
			})
		}

		assert.Equal(t, scores.TickPool{
			Zone1: 5,
			Zone2: 5,
			Top:   5,
			Flash: 5,
		}, pool)
	})
}

func TestTickPool_Sub(t *testing.T) {
	t.Run("SubZone1", func(t *testing.T) {
		pool := scores.TickPool{
			Zone1: 5,
			Zone2: 5,
			Top:   5,
			Flash: 5,
		}

		for range 3 {
			pool = pool.Sub(scores.Tick{
				Zone1:         true,
				AttemptsZone1: 1,
			})
		}

		assert.Equal(t, scores.TickPool{
			Zone1: 2,
			Zone2: 5,
			Top:   5,
			Flash: 5,
		}, pool)
	})

	t.Run("SubZone2", func(t *testing.T) {
		pool := scores.TickPool{
			Zone1: 5,
			Zone2: 5,
			Top:   5,
			Flash: 5,
		}

		for range 3 {
			pool = pool.Sub(scores.Tick{
				Zone1:         true,
				AttemptsZone1: 1,
				Zone2:         true,
				AttemptsZone2: 1,
			})
		}

		assert.Equal(t, scores.TickPool{
			Zone1: 2,
			Zone2: 2,
			Top:   5,
			Flash: 5,
		}, pool)
	})

	t.Run("SubTopWithoutFlash", func(t *testing.T) {
		pool := scores.TickPool{
			Zone1: 5,
			Zone2: 5,
			Top:   5,
			Flash: 5,
		}

		for range 3 {
			pool = pool.Sub(scores.Tick{
				Zone1:         true,
				AttemptsZone1: 999,
				Zone2:         true,
				AttemptsZone2: 999,
				Top:           true,
				AttemptsTop:   999,
			})
		}

		assert.Equal(t, scores.TickPool{
			Zone1: 2,
			Zone2: 2,
			Top:   2,
			Flash: 5,
		}, pool)
	})

	t.Run("SubTopWithFlash", func(t *testing.T) {
		pool := scores.TickPool{
			Zone1: 5,
			Zone2: 5,
			Top:   5,
			Flash: 5,
		}

		for range 3 {
			pool = pool.Sub(scores.Tick{
				Zone1:         true,
				AttemptsZone1: 1,
				Zone2:         true,
				AttemptsZone2: 1,
				Top:           true,
				AttemptsTop:   1,
			})
		}

		assert.Equal(t, scores.TickPool{
			Zone1: 2,
			Zone2: 2,
			Top:   2,
			Flash: 2,
		}, pool)
	})
}
