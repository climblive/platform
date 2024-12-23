package scores

import (
	"iter"
	"slices"
)

type HardestProblems struct {
	Number int
}

func (r *HardestProblems) CalculateScore(points iter.Seq[int]) int {
	score := 0

	n := 0
	for _, p := range slices.Backward(slices.Sorted(points)) {
		if n >= r.Number {
			break
		}

		score += p
		n++
	}

	return score
}
