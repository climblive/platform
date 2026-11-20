package main

import (
	"testing"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestBuildTick(t *testing.T) {
	tests := []struct {
		name         string
		zone1Enabled bool
		zone2Enabled bool
	}{
		{name: "NoZones"},
		{name: "Zone1", zone1Enabled: true},
		{name: "Zone2", zone2Enabled: true},
		{name: "BothZones", zone1Enabled: true, zone2Enabled: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			problem := domain.Problem{
				ID:           42,
				Zone1Enabled: tt.zone1Enabled,
				Zone2Enabled: tt.zone2Enabled,
			}
			attemptCounts := make(map[int]bool)
			sawAttemptOnly := false
			sawFlash := false

			for range 1000 {
				tick := buildTick(problem)
				tick.Revision = 1
				tick = testutils.MustTick(tick)

				assert.Equal(t, problem.ID, tick.ProblemID)
				assert.GreaterOrEqual(t, tick.AttemptsTop, 1)
				assert.LessOrEqual(t, tick.AttemptsTop, 5)
				assert.GreaterOrEqual(t, tick.AttemptsZone1, 1)
				assert.LessOrEqual(t, tick.AttemptsZone1, tick.AttemptsZone2)
				assert.LessOrEqual(t, tick.AttemptsZone2, tick.AttemptsTop)

				if tick.AttemptsTop == 1 {
					sawFlash = true
					assert.True(t, tick.Top)
					assert.Equal(t, 1, tick.AttemptsZone1)
					assert.Equal(t, 1, tick.AttemptsZone2)
				} else {
					attemptCounts[tick.AttemptsTop] = true
				}

				if !tick.Zone1 {
					sawAttemptOnly = true
					assert.False(t, tick.Zone2)
					assert.False(t, tick.Top)
				}

				if !problem.Zone1Enabled && !tick.Top {
					assert.False(t, tick.Zone1)
					assert.False(t, tick.Zone2)
				} else if !problem.Zone2Enabled && !tick.Top {
					assert.False(t, tick.Zone2)
				}
			}

			assert.Greater(t, len(attemptCounts), 1)
			assert.True(t, sawAttemptOnly)
			assert.True(t, sawFlash)
		})
	}
}
