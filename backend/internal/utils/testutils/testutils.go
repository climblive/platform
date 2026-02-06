package testutils

import (
	"sync/atomic"

	"github.com/climblive/platform/backend/internal/domain"
)

var counter atomic.Int64

func RandomResourceID[T domain.ResourceIDType]() T {
	return T(counter.Add(1))
}
