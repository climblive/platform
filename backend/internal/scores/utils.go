package scores

import (
	"iter"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
)

func CompareScore(s1, s2 domain.Score) bool {
	s1.Timestamp = time.Time{}
	s2.Timestamp = time.Time{}

	return s1 == s2
}

func Points(ticks iter.Seq[Tick]) iter.Seq[int] {
	return func(yield func(int) bool) {
		for tick := range ticks {
			if !yield(tick.Points) {
				return
			}
		}
	}
}
