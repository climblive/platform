package scores_test

import (
	"slices"
	"testing"

	"github.com/climblive/platform/backend/internal/scores"
	"github.com/stretchr/testify/assert"
)

func TestHardestProblems(t *testing.T) {
	rules := scores.HardestProblems{
		Number: 5,
	}

	tickPointValues := []int{
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

	score := rules.CalculateScore(slices.Values(tickPointValues))

	assert.Equal(t, 1450, score)
}

func TestHardestProblems_NoLimit(t *testing.T) {
	rules := scores.HardestProblems{
		Number: 0,
	}

	tickPointValues := []int{
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

	score := rules.CalculateScore(slices.Values(tickPointValues))

	assert.Equal(t, 1775, score)
}
