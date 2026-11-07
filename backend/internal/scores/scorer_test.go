package scores_test

import (
	"slices"
	"testing"

	"github.com/climblive/platform/backend/internal/scores"
	"github.com/stretchr/testify/assert"
)

func TestScorerWithProblemLimit(t *testing.T) {
	rules := scores.Scorer{
		ProblemLimit: 5,
	}

	tickPointValues := []float64{
		100,
		50,
		25,
		200,
		250,
		500,
		200,
		300,
		50,
		100,
	}

	score := rules.CalculatePoints(slices.Values(tickPointValues))

	assert.Equal(t, 1450.0, score)
}

func TestScorer_NoLimit(t *testing.T) {
	rules := scores.Scorer{
		ProblemLimit: 0,
	}

	tickPointValues := []float64{
		100,
		50,
		25,
		200,
		250,
		500,
		200,
		300,
		50,
		100,
	}

	score := rules.CalculatePoints(slices.Values(tickPointValues))

	assert.Equal(t, 1775.0, score)
}

func TestScorerRoundsFractionalTotal(t *testing.T) {
	scorer := scores.Scorer{}
	assert.Equal(t, 0.3, scorer.CalculatePoints(slices.Values([]float64{0.1, 0.2})))
}
