package domain_test

import (
	"encoding/json/v2"
	"math"
	"testing"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTenthsFromFloat64(t *testing.T) {
	oneTenth := 0.1
	twoTenths := 0.2
	tests := []struct {
		name     string
		points   float64
		expected domain.Tenths
	}{
		{name: "Zero", points: 0, expected: 0},
		{name: "Integer", points: 25, expected: 250},
		{name: "Fractional", points: 24.7, expected: 247},
		{name: "Negative", points: -0.3, expected: -3},
		{name: "FloatingPointNoise", points: oneTenth + twoTenths, expected: 3},
		{name: "MaximumInteger", points: 2_147_483_647, expected: 21_474_836_470},
		{name: "LargeFractional", points: 2_147_483_646.9, expected: 21_474_836_469},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, domain.TenthsFromFloat64(tt.points))
			assert.Equal(t, tt.expected, domain.TenthsFromFloat64(tt.expected.Float64()))
		})
	}
}

func TestTenthsMarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		value    domain.Tenths
		expected string
	}{
		{name: "Zero", value: 0, expected: "0"},
		{name: "Integer", value: 250, expected: "25"},
		{name: "Fractional", value: 247, expected: "24.7"},
		{name: "OneTenth", value: 1, expected: "0.1"},
		{name: "NegativeTenth", value: -1, expected: "-0.1"},
		{name: "NegativeInteger", value: -250, expected: "-25"},
		{name: "NegativeFractional", value: -247, expected: "-24.7"},
		{name: "LargeFractional", value: 21_474_836_469, expected: "2147483646.9"},
		{name: "BeyondFloatPrecision", value: 9_007_199_254_740_993, expected: "900719925474099.3"},
		{name: "Maximum", value: math.MaxInt64, expected: "922337203685477580.7"},
		{name: "Minimum", value: math.MinInt64, expected: "-922337203685477580.8"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := json.Marshal(tt.value)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, string(encoded))
			assert.Equal(t, tt.expected, tt.value.String())
		})
	}
}

func TestTenthsUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		encoded  string
		expected domain.Tenths
	}{
		{name: "Zero", encoded: "0", expected: 0},
		{name: "Integer", encoded: "25", expected: 250},
		{name: "Fractional", encoded: "24.7", expected: 247},
		{name: "OneTenth", encoded: "0.1", expected: 1},
		{name: "NegativeTenth", encoded: "-0.1", expected: -1},
		{name: "NegativeFractional", encoded: "-24.7", expected: -247},
		{name: "TrailingZeros", encoded: "24.700", expected: 247},
		{name: "Exponent", encoded: "247e-1", expected: 247},
		{name: "FloatingPointNoise", encoded: "0.30000000000000004", expected: 3},
		{name: "LargeInteger", encoded: "2147483647", expected: 21_474_836_470},
		{name: "LargeFractional", encoded: "2147483646.9", expected: 21_474_836_469},
		{name: "Null", encoded: "null", expected: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value := domain.Tenths(123)
			require.NoError(t, json.Unmarshal([]byte(tt.encoded), &value))
			assert.Equal(t, tt.expected, value)

			encoded, err := json.Marshal(value)
			require.NoError(t, err)
			var decoded domain.Tenths
			require.NoError(t, json.Unmarshal(encoded, &decoded))
			assert.Equal(t, value, decoded)
		})
	}
}

func TestTenthsUnmarshalJSONInvalid(t *testing.T) {
	for _, encoded := range []string{
		"0.15", "-0.15", "1e-2", "1e18", "-1e18", "1e309", "1e308",
		"922337203685477580.8", "-922337203685477785.6",
		`"0.1"`, "true", "{}", "[]", "NaN", "Infinity", "01", "1 2", "",
	} {
		t.Run(encoded, func(t *testing.T) {
			value := domain.Tenths(123)
			assert.Error(t, value.UnmarshalJSON([]byte(encoded)))
			assert.Equal(t, domain.Tenths(123), value)
		})
	}
}
