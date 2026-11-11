package scores

import (
	"iter"
	"slices"
)

type Scorer struct {
	ProblemLimit int
}

func (r *Scorer) CalculatePoints(points iter.Seq[Tenths]) Tenths {
	var score Tenths

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
