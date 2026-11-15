package scores

import (
	"time"

	"github.com/climblive/platform/backend/internal/domain"
)

func CompareResult(r1, r2 domain.Result) bool {
	r1.Timestamp = time.Time{}
	r2.Timestamp = time.Time{}

	return r1 == r2
}

func ComparePointValue(pv1, pv2 domain.PointValue) bool {
	return pv1 == pv2
}
