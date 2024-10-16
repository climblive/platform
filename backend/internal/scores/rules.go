package scores

import (
	"iter"
	"slices"
)

type HardestProblems struct {
	Number int
}

func (r *HardestProblems) CalculateScore(tickPointValues iter.Seq[int]) int {
	score := 0

	n := 0
	for _, points := range slices.Backward(slices.Sorted(tickPointValues)) {
		if n >= r.Number {
			break
		}

		score += points
		n++
	}

	return score
}
