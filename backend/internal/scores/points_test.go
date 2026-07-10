package scores_test

import (
	"testing"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/scores"
	"github.com/climblive/platform/backend/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestCalculatePoints(t *testing.T) {
	problem := domain.ProblemValue{
		PointsTop:   100,
		PointsZone1: 50,
		PointsZone2: 75,
		FlashBonus:  10,
	}

	tests := []struct {
		name     string
		tick     scores.Tick
		expected int
	}{
		{
			name:     "NoAttempts",
			tick:     scores.Tick{},
			expected: 0,
		},
		{
			name: "SingleAttemptNoLuck",
			tick: scores.Tick{
				AttemptsTop:   1,
				AttemptsZone1: 1,
				AttemptsZone2: 1,
			},
			expected: 0,
		},
		{
			name: "Flash",
			tick: scores.Tick{
				Top:           true,
				AttemptsTop:   1,
				Zone1:         true,
				AttemptsZone1: 1,
				Zone2:         true,
				AttemptsZone2: 1,
			},
			expected: 110,
		},
		{
			name: "TopWithSeveralAttempts",
			tick: scores.Tick{
				Top:           true,
				AttemptsTop:   999,
				Zone1:         true,
				AttemptsZone1: 999,
				Zone2:         true,
				AttemptsZone2: 999,
			},
			expected: 100,
		},
		{
			name: "Zone1WithSeveralAttempts",
			tick: scores.Tick{
				AttemptsTop:   999,
				Zone1:         true,
				AttemptsZone1: 999,
				AttemptsZone2: 999,
			},
			expected: 50,
		},
		{
			name: "Zone2WithSeveralAttempts",
			tick: scores.Tick{
				AttemptsTop:   999,
				Zone1:         true,
				AttemptsZone1: 999,
				Zone2:         true,
				AttemptsZone2: 999,
			},
			expected: 75,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, scores.CalculatePoints(problem, tt.tick))
		})
	}
}

func TestTick(t *testing.T) {
	fakedContenderID := testutils.RandomResourceID[domain.ContenderID]()
	fakedProblemID := testutils.RandomResourceID[domain.ProblemID]()

	makeFakes := func() (scores.Tick, scores.Tick, scores.Tick, scores.Tick, scores.Tick) {
		none := scores.Tick{
			ContenderID:   fakedContenderID,
			ProblemID:     fakedProblemID,
			Zone1:         false,
			Zone2:         false,
			Top:           false,
			AttemptsZone1: 10,
			AttemptsZone2: 10,
			AttemptsTop:   10,
		}

		zone1 := scores.Tick{
			ContenderID:   fakedContenderID,
			ProblemID:     fakedProblemID,
			Zone1:         true,
			Zone2:         false,
			Top:           false,
			AttemptsZone1: 10,
			AttemptsZone2: 10,
			AttemptsTop:   10,
		}

		zone2 := scores.Tick{
			ContenderID:   fakedContenderID,
			ProblemID:     fakedProblemID,
			Zone1:         true,
			Zone2:         true,
			Top:           false,
			AttemptsZone1: 10,
			AttemptsZone2: 10,
			AttemptsTop:   10,
		}

		top := scores.Tick{
			ContenderID:   fakedContenderID,
			ProblemID:     fakedProblemID,
			Zone1:         true,
			Zone2:         true,
			Top:           true,
			AttemptsZone1: 10,
			AttemptsZone2: 10,
			AttemptsTop:   10,
		}

		flash := scores.Tick{
			ContenderID:   fakedContenderID,
			ProblemID:     fakedProblemID,
			Zone1:         true,
			Zone2:         true,
			Top:           true,
			AttemptsZone1: 1,
			AttemptsZone2: 1,
			AttemptsTop:   1,
		}

		return none, zone1, zone2, top, flash
	}

	t.Run("TurnIntoZone1", func(t *testing.T) {
		none, zone1, zone2, top, flash := makeFakes()

		cases := []struct {
			name     string
			tick     scores.Tick
			expected scores.Tick
		}{
			{
				name: "None",
				tick: none,
				expected: scores.Tick{
					ContenderID:   fakedContenderID,
					ProblemID:     fakedProblemID,
					Zone1:         true,
					Zone2:         false,
					Top:           false,
					AttemptsZone1: 11,
					AttemptsZone2: 11,
					AttemptsTop:   11,
				},
			},
			{
				name: "Zone1",
				tick: zone1,
				expected: scores.Tick{
					ContenderID:   fakedContenderID,
					ProblemID:     fakedProblemID,
					Zone1:         true,
					Zone2:         false,
					Top:           false,
					AttemptsZone1: 10,
					AttemptsZone2: 10,
					AttemptsTop:   10,
				},
			},
			{
				name: "Zone2",
				tick: zone2,
				expected: scores.Tick{
					ContenderID:   fakedContenderID,
					ProblemID:     fakedProblemID,
					Zone1:         true,
					Zone2:         false,
					Top:           false,
					AttemptsZone1: 10,
					AttemptsZone2: 10,
					AttemptsTop:   10,
				},
			},
			{
				name: "Top",
				tick: top,
				expected: scores.Tick{
					ContenderID:   fakedContenderID,
					ProblemID:     fakedProblemID,
					Zone1:         true,
					Zone2:         false,
					Top:           false,
					AttemptsZone1: 10,
					AttemptsZone2: 10,
					AttemptsTop:   10,
				},
			},
			{
				name: "Flash",
				tick: flash,
				expected: scores.Tick{
					ContenderID:   fakedContenderID,
					ProblemID:     fakedProblemID,
					Zone1:         true,
					Zone2:         false,
					Top:           false,
					AttemptsZone1: 1,
					AttemptsZone2: 1,
					AttemptsTop:   1,
				},
			},
		}

		for _, tt := range cases {
			t.Run(tt.name, func(t *testing.T) {
				assert.Equal(t, tt.expected, tt.tick.TurnIntoZone1())
			})
		}
	})

	t.Run("TurnIntoZone2", func(t *testing.T) {
		none, zone1, zone2, top, flash := makeFakes()

		cases := []struct {
			name     string
			tick     scores.Tick
			expected scores.Tick
		}{
			{
				name: "None",
				tick: none,
				expected: scores.Tick{
					ContenderID:   fakedContenderID,
					ProblemID:     fakedProblemID,
					Zone1:         true,
					Zone2:         true,
					Top:           false,
					AttemptsZone1: 11,
					AttemptsZone2: 11,
					AttemptsTop:   11,
				},
			},
			{
				name: "Zone1",
				tick: zone1,
				expected: scores.Tick{
					ContenderID:   fakedContenderID,
					ProblemID:     fakedProblemID,
					Zone1:         true,
					Zone2:         true,
					Top:           false,
					AttemptsZone1: 10,
					AttemptsZone2: 11,
					AttemptsTop:   11,
				},
			},
			{
				name: "Zone2",
				tick: zone2,
				expected: scores.Tick{
					ContenderID:   fakedContenderID,
					ProblemID:     fakedProblemID,
					Zone1:         true,
					Zone2:         true,
					Top:           false,
					AttemptsZone1: 10,
					AttemptsZone2: 10,
					AttemptsTop:   10,
				},
			},
			{
				name: "Top",
				tick: top,
				expected: scores.Tick{
					ContenderID:   fakedContenderID,
					ProblemID:     fakedProblemID,
					Zone1:         true,
					Zone2:         true,
					Top:           false,
					AttemptsZone1: 10,
					AttemptsZone2: 10,
					AttemptsTop:   10,
				},
			},
			{
				name: "Flash",
				tick: flash,
				expected: scores.Tick{
					ContenderID:   fakedContenderID,
					ProblemID:     fakedProblemID,
					Zone1:         true,
					Zone2:         true,
					Top:           false,
					AttemptsZone1: 1,
					AttemptsZone2: 1,
					AttemptsTop:   1,
				},
			},
		}

		for _, tt := range cases {
			t.Run(tt.name, func(t *testing.T) {
				assert.Equal(t, tt.expected, tt.tick.TurnIntoZone2())
			})
		}
	})

	t.Run("TurnIntoRedpoint", func(t *testing.T) {
		none, zone1, zone2, top, flash := makeFakes()

		cases := []struct {
			name     string
			tick     scores.Tick
			expected scores.Tick
		}{
			{
				name: "None",
				tick: none,
				expected: scores.Tick{
					ContenderID:   fakedContenderID,
					ProblemID:     fakedProblemID,
					Zone1:         true,
					Zone2:         true,
					Top:           true,
					AttemptsZone1: 11,
					AttemptsZone2: 11,
					AttemptsTop:   11,
				},
			},
			{
				name: "Zone1",
				tick: zone1,
				expected: scores.Tick{
					ContenderID:   fakedContenderID,
					ProblemID:     fakedProblemID,
					Zone1:         true,
					Zone2:         true,
					Top:           true,
					AttemptsZone1: 10,
					AttemptsZone2: 11,
					AttemptsTop:   11,
				},
			},
			{
				name: "Zone2",
				tick: zone2,
				expected: scores.Tick{
					ContenderID:   fakedContenderID,
					ProblemID:     fakedProblemID,
					Zone1:         true,
					Zone2:         true,
					Top:           true,
					AttemptsZone1: 10,
					AttemptsZone2: 10,
					AttemptsTop:   11,
				},
			},
			{
				name: "Top",
				tick: top,
				expected: scores.Tick{
					ContenderID:   fakedContenderID,
					ProblemID:     fakedProblemID,
					Zone1:         true,
					Zone2:         true,
					Top:           true,
					AttemptsZone1: 10,
					AttemptsZone2: 10,
					AttemptsTop:   10,
				},
			},
			{
				name: "Flash",
				tick: flash,
				expected: scores.Tick{
					ContenderID:   fakedContenderID,
					ProblemID:     fakedProblemID,
					Zone1:         true,
					Zone2:         true,
					Top:           true,
					AttemptsZone1: 2,
					AttemptsZone2: 2,
					AttemptsTop:   2,
				},
			},
		}

		for _, tt := range cases {
			t.Run(tt.name, func(t *testing.T) {
				assert.Equal(t, tt.expected, tt.tick.TurnIntoRedpoint())
			})
		}
	})

	t.Run("TurnIntoFlash", func(t *testing.T) {
		none, zone1, zone2, top, flash := makeFakes()

		cases := []struct {
			name     string
			tick     scores.Tick
			expected scores.Tick
		}{
			{
				name: "None",
				tick: none,
				expected: scores.Tick{
					ContenderID:   fakedContenderID,
					ProblemID:     fakedProblemID,
					Zone1:         true,
					Zone2:         true,
					Top:           true,
					AttemptsZone1: 1,
					AttemptsZone2: 1,
					AttemptsTop:   1,
				},
			},
			{
				name: "Zone1",
				tick: zone1,
				expected: scores.Tick{
					ContenderID:   fakedContenderID,
					ProblemID:     fakedProblemID,
					Zone1:         true,
					Zone2:         true,
					Top:           true,
					AttemptsZone1: 1,
					AttemptsZone2: 1,
					AttemptsTop:   1,
				},
			},
			{
				name: "Zone2",
				tick: zone2,
				expected: scores.Tick{
					ContenderID:   fakedContenderID,
					ProblemID:     fakedProblemID,
					Zone1:         true,
					Zone2:         true,
					Top:           true,
					AttemptsZone1: 1,
					AttemptsZone2: 1,
					AttemptsTop:   1,
				},
			},
			{
				name: "Top",
				tick: top,
				expected: scores.Tick{
					ContenderID:   fakedContenderID,
					ProblemID:     fakedProblemID,
					Zone1:         true,
					Zone2:         true,
					Top:           true,
					AttemptsZone1: 1,
					AttemptsZone2: 1,
					AttemptsTop:   1,
				},
			},
			{
				name: "Flash",
				tick: flash,
				expected: scores.Tick{
					ContenderID:   fakedContenderID,
					ProblemID:     fakedProblemID,
					Zone1:         true,
					Zone2:         true,
					Top:           true,
					AttemptsZone1: 1,
					AttemptsZone2: 1,
					AttemptsTop:   1,
				},
			},
		}

		for _, tt := range cases {
			t.Run(tt.name, func(t *testing.T) {
				assert.Equal(t, tt.expected, tt.tick.TurnIntoFlash())
			})
		}
	})
}
