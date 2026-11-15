package testutils

import (
	"sync/atomic"
	"uuid"

	"github.com/climblive/platform/backend/internal/domain"
)

var counter atomic.Int64

func RandomCompetitionID() domain.CompetitionID {
	return domain.CompetitionID(uuid.New())
}

func RandomResourceID[T domain.ResourceIDType]() T {
	return T(counter.Add(1))
}
