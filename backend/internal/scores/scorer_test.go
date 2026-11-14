package scores_test

import (
	"slices"
	"testing"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/scores"
	"github.com/stretchr/testify/assert"
)

func TestScorerWithProblemLimit(t *testing.T) {
	rules := scores.Scorer{
		ProblemLimit: 5,
	}

	tickPointValues := []domain.Tenths{
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

	assert.Equal(t, domain.Tenths(14500), score)
}

func TestScorer_NoLimit(t *testing.T) {
	rules := scores.Scorer{
		ProblemLimit: 0,
	}

	tickPointValues := []domain.Tenths{
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

	assert.Equal(t, domain.Tenths(17750), score)
}
