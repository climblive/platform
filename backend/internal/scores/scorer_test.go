package scores_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/climblive/platform/backend/internal/scores"
	"github.com/stretchr/testify/assert"
)

func TestScorerWithProblemLimit(t *testing.T) {
	rules := scores.Scorer{
		ProblemLimit: 5,
	}

	tickPointValues := []scores.Tenths{
		1000,
		500,
		250,
		2000,
		2500,
		5000,
		2000,
		3000,
		500,
		1000,
	}

	score := rules.CalculatePoints(slices.Values(tickPointValues))

	assert.Equal(t, scores.Tenths(14500), score)
}

func TestScorer_NoLimit(t *testing.T) {
	rules := scores.Scorer{
		ProblemLimit: 0,
	}

	tickPointValues := []scores.Tenths{
		1000,
		500,
		250,
		2000,
		2500,
		5000,
		2000,
		3000,
		500,
		1000,
	}

	score := rules.CalculatePoints(slices.Values(tickPointValues))

	assert.Equal(t, scores.Tenths(17750), score)
}

func TestScorerFractionalTotals(t *testing.T) {
	tests := []struct {
		name     string
		points   []scores.Tenths
		expected scores.Tenths
	}{
		{name: "Empty", expected: 0},
		{name: "TwoTenths", points: []scores.Tenths{1, 2}, expected: 3},
		{name: "RepeatedTenths", points: slices.Repeat([]scores.Tenths{1}, 8), expected: 8},
		{name: "PooledPoints", points: []scores.Tenths{333, 333, 333}, expected: 999},
		{name: "SmallFirst", points: []scores.Tenths{1, 1, 6}, expected: 8},
		{name: "LargeMiddle", points: []scores.Tenths{1, 6, 1}, expected: 8},
		{name: "LargeFirst", points: []scores.Tenths{6, 1, 1}, expected: 8},
		{
			name:     "LargeTotal",
			points:   slices.Repeat([]scores.Tenths{42_949_672_940}, 65536),
			expected: scores.Tenths(42_949_672_940) * 65536,
		},
	}

	for _, tt := range tests {
		for _, limit := range []int{0, len(tt.points)} {
			t.Run(fmt.Sprintf("%s/Limit%d", tt.name, limit), func(t *testing.T) {
				scorer := scores.Scorer{ProblemLimit: limit}
				assert.Equal(t, tt.expected, scorer.CalculatePoints(slices.Values(tt.points)))
			})
		}
	}
}

func TestScorerFractionalProblemLimit(t *testing.T) {
	scorer := scores.Scorer{ProblemLimit: 2}
	points := []scores.Tenths{332, 333, 1, 334}

	assert.Equal(t, scores.Tenths(667), scorer.CalculatePoints(slices.Values(points)))
}
