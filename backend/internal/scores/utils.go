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

func ComparePointValue(pv1, pv2 domain.PointValue) bool {
	return pv1 == pv2
}

func Points(values iter.Seq[domain.PointValue]) iter.Seq[int] {
	return func(yield func(int) bool) {
		for value := range values {
			if !yield(value.Current) {
				return
			}
		}
	}
}
