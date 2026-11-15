package testutils

import (
	"math"
	"testing"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/scores"
	"github.com/stretchr/testify/assert"
)

func TestRandomResourceID(t *testing.T) {
	id1 := RandomResourceID[domain.CompetitorID]()
	id2 := RandomResourceID[domain.CompetitorID]()
	id3 := RandomResourceID[domain.CompetitorID]()

	assert.Greater(t, id1, domain.CompetitorID(0))
	assert.Equal(t, id1+1, id2)
	assert.Equal(t, id2+1, id3)
}

func TestMustTick(t *testing.T) {
	competitorID := RandomResourceID[domain.CompetitorID]()
	validTick := domain.Tick{
		ID: RandomResourceID[domain.TickID](),
		Ownership: domain.OwnershipData{
			OrganizerID:  RandomResourceID[domain.OrganizerID](),
			CompetitorID: &competitorID,
		},
		Timestamp:     time.Now(),
		CompetitionID: RandomResourceID[domain.CompetitionID](),
		ProblemID:     RandomResourceID[domain.ProblemID](),
		Revision:      1,
		Zone1:         true,
		AttemptsZone1: 2,
		Zone2:         true,
		AttemptsZone2: 3,
		Top:           true,
		AttemptsTop:   4,
	}

	check := func(t *testing.T, tick domain.Tick, wantPanic bool) {
		t.Helper()
		scoreTick := scores.Tick{
			ID:            tick.ID,
			Revision:      tick.Revision,
			CompetitorID:  competitorID,
			ProblemID:     tick.ProblemID,
			Zone1:         tick.Zone1,
			AttemptsZone1: tick.AttemptsZone1,
			Zone2:         tick.Zone2,
			AttemptsZone2: tick.AttemptsZone2,
			Top:           tick.Top,
			AttemptsTop:   tick.AttemptsTop,
		}

		t.Run("Domain", func(t *testing.T) {
			if wantPanic {
				assert.Panics(t, func() { MustTick(tick) })
				return
			}
			assert.Equal(t, tick, MustTick(tick))
		})
		t.Run("Scores", func(t *testing.T) {
			if wantPanic {
				assert.Panics(t, func() { MustScoreTick(scoreTick) })
				return
			}
			assert.Equal(t, scoreTick, MustScoreTick(scoreTick))
		})
	}

	valid := map[string]domain.Tick{
		"Top":             validTick,
		"NoAttempts":      {Revision: 1},
		"Attempted":       {Revision: 1, AttemptsZone1: 3, AttemptsZone2: 3, AttemptsTop: 3},
		"Zone1":           {Revision: 1, Zone1: true, AttemptsZone1: 1, AttemptsZone2: 2, AttemptsTop: 3},
		"Zone2":           {Revision: 1, Zone1: true, AttemptsZone1: 1, Zone2: true, AttemptsZone2: 2, AttemptsTop: 3},
		"Flash":           {Revision: 1, Zone1: true, AttemptsZone1: 1, Zone2: true, AttemptsZone2: 1, Top: true, AttemptsTop: 1},
		"MaximumAttempts": {Revision: 1, Zone1: true, AttemptsZone1: 999, Zone2: true, AttemptsZone2: 999, Top: true, AttemptsTop: 999},
		"MaximumRevision": {Revision: math.MaxUint32},
	}
	for name, tick := range valid {
		t.Run(name, func(t *testing.T) {
			check(t, tick, false)
		})
	}

	invalid := map[string]func(*domain.Tick){
		"MissingRevision":       func(tick *domain.Tick) { tick.Revision = 0 },
		"RevisionTooLarge":      func(tick *domain.Tick) { tick.Revision = math.MaxUint32 + 1 },
		"NegativeZone1Attempts": func(tick *domain.Tick) { tick.AttemptsZone1 = -1 },
		"NegativeZone2Attempts": func(tick *domain.Tick) { tick.AttemptsZone2 = -1 },
		"NegativeTopAttempts":   func(tick *domain.Tick) { tick.AttemptsTop = -1 },
		"Zone1AttemptsTooLarge": func(tick *domain.Tick) { tick.AttemptsZone1 = 1000 },
		"Zone2AttemptsTooLarge": func(tick *domain.Tick) { tick.AttemptsZone2 = 1000 },
		"TopAttemptsTooLarge":   func(tick *domain.Tick) { tick.AttemptsTop = 1000 },
		"Zone1WithoutAttempts":  func(tick *domain.Tick) { tick.AttemptsZone1 = 0 },
		"Zone2WithoutAttempts": func(tick *domain.Tick) {
			tick.AttemptsZone1 = 0
			tick.AttemptsZone2 = 0
			tick.Zone1 = false
		},
		"TopWithoutAttempts": func(tick *domain.Tick) {
			tick.AttemptsZone1 = 0
			tick.AttemptsZone2 = 0
			tick.AttemptsTop = 0
			tick.Zone1 = false
			tick.Zone2 = false
		},
		"Zone1AttemptsExceedZone2": func(tick *domain.Tick) { tick.AttemptsZone1 = 4 },
		"Zone2AttemptsExceedTop":   func(tick *domain.Tick) { tick.AttemptsZone2 = 5 },
		"TopWithoutZone1":          func(tick *domain.Tick) { tick.Zone1 = false },
		"TopWithoutZone2":          func(tick *domain.Tick) { tick.Zone2 = false },
		"Zone2WithoutZone1": func(tick *domain.Tick) {
			tick.Top = false
			tick.Zone1 = false
		},
	}
	for name, mutate := range invalid {
		t.Run(name, func(t *testing.T) {
			tick := validTick
			mutate(&tick)
			check(t, tick, true)
		})
	}
}
