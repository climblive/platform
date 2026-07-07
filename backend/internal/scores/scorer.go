package scores

import (
	"iter"
	"slices"
)

type Scorer struct {
	ProblemLimit int
}

func (r *Scorer) CalculateScore(points iter.Seq[int]) int {
	score := 0

	if r.ProblemLimit == 0 {
		for p := range points {
			score += p
		}

		return score
	}

	n := 0
	for _, p := range slices.Backward(slices.Sorted(points)) {
		if n >= r.ProblemLimit {
			break
		}

		score += p
		n++
	}

	return score
}
