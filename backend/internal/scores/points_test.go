package scores_test

import (
	"testing"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/scores"
	"github.com/climblive/platform/backend/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestCalculatePoints(t *testing.T) {
	problem := scores.ProblemValue{
		PointsTop:   1000,
		PointsZone1: 500,
		PointsZone2: 750,
		FlashBonus:  100,
	}

	tests := []struct {
		name     string
		tick     scores.Tick
		expected domain.Tenths
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
			expected: 1100,
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
			expected: 1000,
		},
		{
			name: "Zone1WithSeveralAttempts",
			tick: scores.Tick{
				AttemptsTop:   999,
				Zone1:         true,
				AttemptsZone1: 999,
				AttemptsZone2: 999,
			},
			expected: 500,
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
			expected: 750,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, scores.CalculatePoints(problem, tt.tick, scores.Rules{}))
		})
	}
}

func TestCalculatePoints_HighestValueFeatureWins(t *testing.T) {
	problem := scores.ProblemValue{
		PointsTop:   1000,
		PointsZone1: 500,
		PointsZone2: 750,
		FlashBonus:  100,
	}

	t.Run("Zone1WinsOverZone2", func(t *testing.T) {
		tick := scores.Tick{
			AttemptsTop:   4,
			Zone2:         true,
			AttemptsZone2: 4,
			Zone1:         true,
			AttemptsZone1: 1,
		}

		rules := scores.Rules{PointDeduction: 100}

		assert.Equal(t, domain.Tenths(500), scores.CalculatePoints(problem, tick, rules))
	})

	t.Run("Zone2WinsOverTop", func(t *testing.T) {
		tick := scores.Tick{
			Top:           true,
			AttemptsTop:   4,
			Zone2:         true,
			AttemptsZone2: 1,
			Zone1:         true,
			AttemptsZone1: 1,
		}

		rules := scores.Rules{PointDeduction: 100}

		assert.Equal(t, domain.Tenths(750), scores.CalculatePoints(problem, tick, rules))
	})
}

func TestCalculatePoints_PointDeduction(t *testing.T) {
	value := scores.ProblemValue{
		PointsTop:   1000,
		PointsZone1: 500,
		PointsZone2: 750,
		FlashBonus:  100,
	}

	t.Run("NoAscent", func(t *testing.T) {
		tick := scores.Tick{
			AttemptsTop:   3,
			AttemptsZone2: 3,
			AttemptsZone1: 3,
		}

		rules := scores.Rules{PointDeduction: 100}

		assert.Equal(t, domain.Tenths(0), scores.CalculatePoints(value, tick, rules))
	})

	t.Run("Flash", func(t *testing.T) {
		tick := scores.Tick{
			Top:           true,
			AttemptsTop:   1,
			Zone2:         true,
			AttemptsZone2: 1,
			Zone1:         true,
			AttemptsZone1: 1,
		}

		rules := scores.Rules{PointDeduction: 10}

		assert.Equal(t, domain.Tenths(1100), scores.CalculatePoints(value, tick, rules))
	})

	t.Run("Zone1", func(t *testing.T) {
		tick := scores.Tick{
			AttemptsTop:   5,
			AttemptsZone2: 3,
			Zone1:         true,
			AttemptsZone1: 2,
		}

		rules := scores.Rules{PointDeduction: 1}

		assert.Equal(t, domain.Tenths(499), scores.CalculatePoints(value, tick, rules))
	})

	t.Run("Zone2", func(t *testing.T) {
		tick := scores.Tick{
			AttemptsTop:   5,
			Zone2:         true,
			AttemptsZone2: 3,
			Zone1:         true,
			AttemptsZone1: 2,
		}

		rules := scores.Rules{PointDeduction: 1}

		assert.Equal(t, domain.Tenths(748), scores.CalculatePoints(value, tick, rules))
	})

	t.Run("Top", func(t *testing.T) {
		tick := scores.Tick{
			Top:           true,
			AttemptsTop:   5,
			Zone2:         true,
			AttemptsZone2: 3,
			Zone1:         true,
			AttemptsZone1: 2,
		}

		rules := scores.Rules{PointDeduction: 1}

		assert.Equal(t, domain.Tenths(996), scores.CalculatePoints(value, tick, rules))
	})

	t.Run("FloorAtZero", func(t *testing.T) {
		tick := scores.Tick{
			Top:           true,
			AttemptsTop:   999,
			Zone2:         true,
			AttemptsZone2: 999,
			Zone1:         true,
			AttemptsZone1: 999,
		}

		rules := scores.Rules{PointDeduction: 21_474_836_470}

		assert.Equal(t, domain.Tenths(0), scores.CalculatePoints(value, tick, rules))
	})
}

func TestCalculatePoints_MaximumTopAndFlashBonus(t *testing.T) {
	value := scores.ProblemValue{
		PointsTop:  21_474_836_470,
		FlashBonus: 21_474_836_470,
	}
	tick := scores.Tick{Top: true, AttemptsTop: 1}

	assert.Equal(t, domain.Tenths(42_949_672_940), scores.CalculatePoints(value, tick, scores.Rules{}))
}

func TestCalculatePoints_MaxAttempts(t *testing.T) {
	value := scores.ProblemValue{
		PointsTop:   1000,
		PointsZone1: 500,
		PointsZone2: 750,
		FlashBonus:  100,
	}

	t.Run("Flash", func(t *testing.T) {
		tick := scores.Tick{
			Top:           true,
			AttemptsTop:   1,
			Zone2:         true,
			AttemptsZone2: 1,
			Zone1:         true,
			AttemptsZone1: 1,
		}

		rules := scores.Rules{MaxAttempts: 1}

		assert.Equal(t, domain.Tenths(1100), scores.CalculatePoints(value, tick, rules))
	})

	t.Run("AtAttemptLimit", func(t *testing.T) {
		tick := scores.Tick{
			Top:           true,
			AttemptsTop:   10,
			Zone2:         true,
			AttemptsZone2: 10,
			Zone1:         true,
			AttemptsZone1: 10,
		}
		rules := scores.Rules{MaxAttempts: 10}

		assert.Equal(t, domain.Tenths(1000), scores.CalculatePoints(value, tick, rules))
	})

	t.Run("AboveAttemptLimit", func(t *testing.T) {
		tick := scores.Tick{
			Top:           true,
			AttemptsTop:   11,
			Zone2:         true,
			AttemptsZone2: 11,
			Zone1:         true,
			AttemptsZone1: 11,
		}
		rules := scores.Rules{MaxAttempts: 10}

		assert.Equal(t, domain.Tenths(0), scores.CalculatePoints(value, tick, rules))
	})

	t.Run("MaxAttemptsAffectsOnlyReachedFeatures", func(t *testing.T) {
		tick := scores.Tick{
			AttemptsTop:   11,
			AttemptsZone2: 11,
			Zone1:         true,
			AttemptsZone1: 2,
		}
		rules := scores.Rules{MaxAttempts: 10}

		assert.Equal(t, domain.Tenths(500), scores.CalculatePoints(value, tick, rules))
	})
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

	t.Run("TurnIntoTop", func(t *testing.T) {
		none, zone1, zone2, top, flash := makeFakes()

		cases := []struct {
			name     string
			tick     scores.Tick
			expected scores.Tick
		}{
			{
				name: "Unattempted",
				tick: scores.Tick{},
				expected: scores.Tick{
					Zone1:         true,
					Zone2:         true,
					Top:           true,
					AttemptsZone1: 1,
					AttemptsZone2: 1,
					AttemptsTop:   1,
				},
			},
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
					AttemptsZone1: 1,
					AttemptsZone2: 1,
					AttemptsTop:   1,
				},
			},
		}

		for _, tt := range cases {
			t.Run(tt.name, func(t *testing.T) {
				assert.Equal(t, tt.expected, tt.tick.TurnIntoTop())
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
