package testutils

import (
	"sync/atomic"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/scores"
	"github.com/climblive/platform/backend/internal/usecases/validators"
)

var counter atomic.Int64

func RandomResourceID[T domain.ResourceIDType]() T {
	return T(counter.Add(1))
}

func MustTick[T domain.Tick | scores.Tick](tick T) T {
	var domainTick domain.Tick

	switch tick := any(tick).(type) {
	case domain.Tick:
		domainTick = tick
	case scores.Tick:
		domainTick = domain.Tick{
			ID:            tick.ID,
			Ownership:     domain.OwnershipData{},
			Timestamp:     time.Time{},
			CompetitionID: 0,
			Revision:      tick.Revision,
			ProblemID:     tick.ProblemID,
			Zone1:         tick.Zone1,
			AttemptsZone1: tick.AttemptsZone1,
			Zone2:         tick.Zone2,
			AttemptsZone2: tick.AttemptsZone2,
			Top:           tick.Top,
			AttemptsTop:   tick.AttemptsTop,
		}
	}

	if err := (validators.TickValidator{}).Validate(domainTick); err != nil {
		panic(err)
	}

	return tick
}
