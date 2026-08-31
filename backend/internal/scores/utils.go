package scores

import (
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
