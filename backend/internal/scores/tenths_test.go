package scores

import (
	"testing"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestTenthsFromFloat64(t *testing.T) {
	oneTenth := 0.1
	twoTenths := 0.2
	tests := []struct {
		name     string
		points   float64
		expected Tenths
	}{
		{name: "Zero", points: 0, expected: 0},
		{name: "Integer", points: 25, expected: 250},
		{name: "Fractional", points: 24.7, expected: 247},
		{name: "FloatingPointNoise", points: oneTenth + twoTenths, expected: 3},
		{name: "MaximumInteger", points: 2_147_483_647, expected: 21_474_836_470},
		{name: "LargeFractional", points: 2_147_483_646.9, expected: 21_474_836_469},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tenthsFromFloat64(tt.points))
			assert.Equal(t, tt.expected, tenthsFromFloat64(tt.expected.Float64()))
		})
	}
}

func TestTenthsFloat64(t *testing.T) {
	tests := []struct {
		value    Tenths
		expected float64
	}{
		{value: 0, expected: 0},
		{value: 1, expected: 0.1},
		{value: 10, expected: 1},
		{value: 247, expected: 24.7},
		{value: 21_474_836_470, expected: 2_147_483_647},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, tt.value.Float64())
	}
}

func TestProblemValueFromDomain(t *testing.T) {
	value := problemValueFromDomain(domain.ProblemValue{
		PointsZone1: 1,
		PointsZone2: 25,
		PointsTop:   2_147_483_647,
		FlashBonus:  2_147_483_647,
	})

	assert.Equal(t, ProblemValue{
		PointsZone1: 10,
		PointsZone2: 250,
		PointsTop:   21_474_836_470,
		FlashBonus:  21_474_836_470,
	}, value)
}
