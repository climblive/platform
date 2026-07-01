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

func HypotheticalSecondBestTop(t *testing.T) {
	tests := []struct {
		name     string
		input    scores.Tick
		expected scores.Tick
	}{
		{
			name:  "EmptyTick",
			input: scores.Tick{},
			expected: scores.Tick{
				Zone1:         true,
				AttemptsZone1: 999,
				Zone2:         true,
				AttemptsZone2: 999,
				Top:           true,
				AttemptsTop:   999,
			},
		},
		{
			name: "NormalizesToTop",
			input: scores.Tick{
				ContenderID:   1,
				ProblemID:     2,
				Zone1:         true,
				AttemptsZone1: 2,
				Zone2:         true,
				AttemptsZone2: 3,
				AttemptsTop:   4,
			},
			expected: scores.Tick{
				ContenderID:   1,
				ProblemID:     2,
				Zone1:         true,
				AttemptsZone1: 999,
				Zone2:         true,
				AttemptsZone2: 999,
				Top:           true,
				AttemptsTop:   999,
			},
		},
		{
			name: "DropsExistingFlashAttemptCounts",
			input: scores.Tick{
				Zone1:         true,
				AttemptsZone1: 2,
				Zone2:         true,
				AttemptsZone2: 1,
				Top:           true,
				AttemptsTop:   1,
			},
			expected: scores.Tick{
				Zone1:         true,
				AttemptsZone1: 999,
				Zone2:         true,
				AttemptsZone2: 999,
				Top:           true,
				AttemptsTop:   999,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := tt.input

			actual := scores.HypotheticalSecondBestTop(tt.input)

			assert.Equal(t, tt.expected, actual)
			assert.Equal(t, original, tt.input)
		})
	}
}

func HypotheticalBestTop(t *testing.T) {
	input := scores.Tick{
		ContenderID:   1,
		ProblemID:     2,
		Zone1:         true,
		AttemptsZone1: 999,
		Zone2:         true,
		AttemptsZone2: 999,
		Top:           true,
		AttemptsTop:   999,
	}

	assert.Equal(t, scores.Tick{
		ContenderID:   1,
		ProblemID:     2,
		Zone1:         true,
		AttemptsZone1: 1,
		Zone2:         true,
		AttemptsZone2: 1,
		Top:           true,
		AttemptsTop:   1,
	}, scores.HypotheticalBestTop(input))

	assert.Equal(t, scores.Tick{
		ContenderID:   1,
		ProblemID:     2,
		Zone1:         true,
		AttemptsZone1: 999,
		Zone2:         true,
		AttemptsZone2: 999,
		Top:           true,
		AttemptsTop:   999,
	}, input)
}
