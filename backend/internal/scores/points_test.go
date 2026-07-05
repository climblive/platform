package scores_test

import (
	"testing"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/scores"
	"github.com/stretchr/testify/assert"
)

func TestCalculatePoints(t *testing.T) {
	problem := domain.ProblemValue{
		PointsTop:   100,
		PointsZone1: 50,
		PointsZone2: 75,
		FlashBonus:  10,
	}

	tests := []struct {
		name     string
		tick     scores.Tick
		expected int
	}{
		{
			name:     "NoAttempts",
			tick:     scores.Tick{},
			expected: 0,
		},
		{
			name: "SingleAttemptNoLuck",
			tick: scores.Tick{
				AttemptsTop:   1,
				AttemptsZone1: 1,
				AttemptsZone2: 1,
			},
			expected: 0,
		},
		{
			name: "Flash",
			tick: scores.Tick{
				Top:           true,
				AttemptsTop:   1,
				Zone1:         true,
				AttemptsZone1: 1,
				Zone2:         true,
				AttemptsZone2: 1,
			},
			expected: 110,
		},
		{
			name: "TopWithSeveralAttempts",
			tick: scores.Tick{
				Top:           true,
				AttemptsTop:   999,
				Zone1:         true,
				AttemptsZone1: 999,
				Zone2:         true,
				AttemptsZone2: 999,
			},
			expected: 100,
		},
		{
			name: "Zone1WithSeveralAttempts",
			tick: scores.Tick{
				AttemptsTop:   999,
				Zone1:         true,
				AttemptsZone1: 999,
				AttemptsZone2: 999,
			},
			expected: 50,
		},
		{
			name: "Zone2WithSeveralAttempts",
			tick: scores.Tick{
				AttemptsTop:   999,
				Zone1:         true,
				AttemptsZone1: 999,
				Zone2:         true,
				AttemptsZone2: 999,
			},
			expected: 75,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, scores.CalculatePoints(problem, tt.tick))
		})
	}
}

func TestHypotheticalBestZone1(t *testing.T) {
	t.Run("EmptyTick", func(t *testing.T) {
		hypothetical := scores.HypotheticalBestZone1(scores.Tick{})

		assert.Equal(t, scores.Tick{
			Zone1:         true,
			AttemptsZone1: 1,
		}, hypothetical)
	})

	t.Run("Topped", func(t *testing.T) {
		hypothetical := scores.HypotheticalBestZone1(scores.Tick{
			Zone1:         true,
			AttemptsZone1: 10,
			Zone2:         true,
			AttemptsZone2: 10,
			Top:           true,
			AttemptsTop:   10,
		})

		assert.Equal(t, scores.Tick{
			Zone1:         true,
			AttemptsZone1: 10,
			AttemptsZone2: 10,
			AttemptsTop:   10,
		}, hypothetical)
	})
}

func TestHypotheticalBestZone2(t *testing.T) {
	t.Run("EmptyTick", func(t *testing.T) {
		hypothetical := scores.HypotheticalBestZone2(scores.Tick{})

		assert.Equal(t, scores.Tick{
			Zone1:         true,
			AttemptsZone1: 1,
			Zone2:         true,
			AttemptsZone2: 1,
		}, hypothetical)
	})

	t.Run("ReachedZone1", func(t *testing.T) {
		hypothetical := scores.HypotheticalBestZone2(scores.Tick{
			Zone1:         true,
			AttemptsZone1: 5,
			AttemptsZone2: 5,
			AttemptsTop:   5,
		})

		assert.Equal(t, scores.Tick{
			Zone1:         true,
			AttemptsZone1: 5,
			Zone2:         true,
			AttemptsZone2: 6,
			AttemptsTop:   6,
		}, hypothetical)
	})

	t.Run("Topped", func(t *testing.T) {
		hypothetical := scores.HypotheticalBestZone2(scores.Tick{
			Zone1:         true,
			AttemptsZone1: 5,
			Zone2:         true,
			AttemptsZone2: 10,
			AttemptsTop:   10,
		})

		assert.Equal(t, scores.Tick{
			Zone1:         true,
			AttemptsZone1: 5,
			Zone2:         true,
			AttemptsZone2: 10,
			AttemptsTop:   10,
		}, hypothetical)
	})
}

func TestHypotheticalBestTopNoFlash(t *testing.T) {
	t.Run("EmptyTick", func(t *testing.T) {
		hypothetical := scores.HypotheticalBestTopNoFlash(scores.Tick{})

		assert.Equal(t, scores.Tick{
			Zone1:         true,
			AttemptsZone1: 2,
			Zone2:         true,
			AttemptsZone2: 2,
			Top:           true,
			AttemptsTop:   2,
		}, hypothetical)
	})

	t.Run("ReachedZone1", func(t *testing.T) {
		hypothetical := scores.HypotheticalBestTopNoFlash(scores.Tick{
			Zone1:         true,
			AttemptsZone1: 5,
			AttemptsZone2: 5,
			AttemptsTop:   5,
		})

		assert.Equal(t, scores.Tick{
			Zone1:         true,
			AttemptsZone1: 5,
			Zone2:         true,
			AttemptsZone2: 6,
			Top:           true,
			AttemptsTop:   6,
		}, hypothetical)
	})

	t.Run("ReachedZone2", func(t *testing.T) {
		hypothetical := scores.HypotheticalBestTopNoFlash(scores.Tick{
			Zone1:         true,
			AttemptsZone1: 5,
			Zone2:         true,
			AttemptsZone2: 10,
			AttemptsTop:   10,
		})

		assert.Equal(t, scores.Tick{
			Zone1:         true,
			AttemptsZone1: 5,
			Zone2:         true,
			AttemptsZone2: 10,
			Top:           true,
			AttemptsTop:   11,
		}, hypothetical)
	})

	t.Run("Flashed", func(t *testing.T) {
		hypothetical := scores.HypotheticalBestTopNoFlash(scores.Tick{
			Zone1:         true,
			AttemptsZone1: 1,
			Zone2:         true,
			AttemptsZone2: 1,
			Top:           true,
			AttemptsTop:   1,
		})

		assert.Equal(t, scores.Tick{
			Zone1:         true,
			AttemptsZone1: 2,
			Zone2:         true,
			AttemptsZone2: 2,
			Top:           true,
			AttemptsTop:   2,
		}, hypothetical)
	})

	t.Run("AlreadyReachedTop", func(t *testing.T) {
		hypothetical := scores.HypotheticalBestTopNoFlash(scores.Tick{
			Zone1:         true,
			AttemptsZone1: 2,
			Zone2:         true,
			AttemptsZone2: 2,
			Top:           true,
			AttemptsTop:   2,
		})

		assert.Equal(t, scores.Tick{
			Zone1:         true,
			AttemptsZone1: 2,
			Zone2:         true,
			AttemptsZone2: 2,
			Top:           true,
			AttemptsTop:   2,
		}, hypothetical)
	})
}
