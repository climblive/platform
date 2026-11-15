package scores_test

import (
	"fmt"
	"math/rand"
	"slices"
	"testing"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/scores"
	"github.com/stretchr/testify/assert"
)

func TestBasicRankerPointsFormatting(t *testing.T) {
	tests := []struct {
		name     string
		points   domain.Tenths
		expected string
	}{
		{name: "Zero", points: 0, expected: "0p"},
		{name: "Integer", points: 250, expected: "25p"},
		{name: "Fractional", points: 247, expected: "24.7p"},
		{name: "LargeInteger", points: 10_000_000, expected: "1000000p"},
		{name: "LargeFractional", points: 10_000_001, expected: "1000000.1p"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ranker := scores.NewBasicRanker(0, true, false, false)
			competitors := []scores.Competitor{{ID: 1, Score: scores.Score{Points: tt.points}}}

			ranked := ranker.RankCompetitors(slices.Values(competitors))

			assert.Equal(t, tt.expected, ranked[0].Score)
		})
	}
}

func TestBasicRanker(t *testing.T) {
	makeCompetitors := func(count int) []scores.Competitor {
		competitors := make([]scores.Competitor, count)

		for n := range count {
			competitorID := domain.CompetitorID(n + 1)

			competitors[n] = scores.Competitor{
				ID: competitorID,
			}
		}

		return competitors
	}

	t.Run("CompetitionNotStarted", func(t *testing.T) {
		competitors := makeCompetitors(3)
		shuffleSlice(competitors)

		ranker := scores.NewBasicRanker(5, true, false, false)

		results := ranker.RankCompetitors(slices.Values(competitors))

		expected := []string{
			"i:1 p:1 r:0 f:-",
			"i:2 p:1 r:1 f:-",
			"i:3 p:1 r:2 f:-",
		}

		assert.Equal(t, expected, prettifyAll(results))
	})

	t.Run("ByPoints", func(t *testing.T) {
		t.Run("Simple", func(t *testing.T) {
			competitors := makeCompetitors(3)
			competitors[0].Points = 300
			competitors[1].Points = 200
			competitors[2].Points = 100

			shuffleSlice(competitors)

			ranker := scores.NewBasicRanker(5, true, false, false)

			results := ranker.RankCompetitors(slices.Values(competitors))

			expected := []string{
				"i:1 p:1 r:0 f:🏆",
				"i:2 p:2 r:1 f:🏆",
				"i:3 p:3 r:2 f:🏆",
			}

			assert.Equal(t, expected, prettifyAll(results))
		})

		t.Run("SharedPlacement", func(t *testing.T) {
			competitors := makeCompetitors(5)
			competitors[0].Points = 300
			competitors[1].Points = 200
			competitors[2].Points = 200
			competitors[3].Points = 200
			competitors[4].Points = 100

			shuffleSlice(competitors)

			ranker := scores.NewBasicRanker(5, true, false, false)

			results := ranker.RankCompetitors(slices.Values(competitors))

			expected := []string{
				"i:1 p:1 r:0 f:🏆",
				"i:2 p:2 r:1 f:🏆",
				"i:3 p:2 r:2 f:🏆",
				"i:4 p:2 r:3 f:🏆",
				"i:5 p:5 r:4 f:🏆",
			}

			assert.Equal(t, expected, prettifyAll(results))
		})

		t.Run("IgnoreAttempts", func(t *testing.T) {
			competitors := makeCompetitors(2)
			competitors[0].Points = 100
			competitors[0].Tops = 1
			competitors[0].AttemptsTops = 2
			competitors[1].Points = 100

			ranker := scores.NewBasicRanker(5, true, false, false)

			results := ranker.RankCompetitors(slices.Values(competitors))

			expected := []string{
				"i:1 p:1 r:0 f:🏆",
				"i:2 p:1 r:1 f:🏆",
			}

			assert.Equal(t, expected, prettifyAll(results))
		})

		t.Run("ExtraFinalists", func(t *testing.T) {
			competitors := makeCompetitors(10)
			competitors[0].Points = 500
			competitors[1].Points = 400
			competitors[2].Points = 300
			competitors[3].Points = 200
			competitors[4].Points = 100
			competitors[5].Points = 100
			competitors[6].Points = 100
			competitors[7].Points = 50
			competitors[8].Points = 50
			competitors[9].Points = 50

			shuffleSlice(competitors)

			ranker := scores.NewBasicRanker(5, true, false, false)

			results := ranker.RankCompetitors(slices.Values(competitors))

			expected := []string{
				"i:1 p:1 r:0 f:🏆",
				"i:2 p:2 r:1 f:🏆",
				"i:3 p:3 r:2 f:🏆",
				"i:4 p:4 r:3 f:🏆",
				"i:5 p:5 r:4 f:🏆",
				"i:6 p:5 r:5 f:🏆",
				"i:7 p:5 r:6 f:🏆",
				"i:8 p:8 r:7 f:-",
				"i:9 p:8 r:8 f:-",
				"i:10 p:8 r:9 f:-",
			}

			assert.Equal(t, expected, prettifyAll(results))
		})

		t.Run("WithdrawalsFromFinals", func(t *testing.T) {
			competitors := makeCompetitors(10)
			competitors[0].Points = 500
			competitors[1].Points = 400
			competitors[2].Points = 300
			competitors[3].Points = 200
			competitors[4].Points = 100
			competitors[5].Points = 100
			competitors[6].Points = 100
			competitors[7].Points = 50
			competitors[8].Points = 50
			competitors[9].Points = 0

			competitors[1].WithdrawnFromFinals = true
			competitors[2].WithdrawnFromFinals = true
			competitors[5].WithdrawnFromFinals = true

			shuffleSlice(competitors)

			ranker := scores.NewBasicRanker(5, true, false, false)

			results := ranker.RankCompetitors(slices.Values(competitors))

			expected := []string{
				"i:1 p:1 r:0 f:🏆",
				"i:2 p:2 r:1 f:-",
				"i:3 p:3 r:2 f:-",
				"i:4 p:4 r:3 f:🏆",
				"i:5 p:5 r:4 f:🏆",
				"i:6 p:5 r:5 f:-",
				"i:7 p:5 r:6 f:🏆",
				"i:8 p:8 r:7 f:🏆",
				"i:9 p:8 r:8 f:🏆",
				"i:10 p:10 r:9 f:-",
			}

			assert.Equal(t, expected, prettifyAll(results))
		})

		t.Run("DisqualifiedCompetitorsLast", func(t *testing.T) {
			competitors := makeCompetitors(5)
			competitors[0].Points = 0
			competitors[1].Points = 0
			competitors[2].Points = 0
			competitors[3].Points = 0
			competitors[4].Points = 0

			competitors[1].Disqualified = true

			shuffleSlice(competitors)

			ranker := scores.NewBasicRanker(5, true, false, false)

			results := ranker.RankCompetitors(slices.Values(competitors))

			expected := []string{
				"i:1 p:1 r:0 f:-",
				"i:3 p:1 r:1 f:-",
				"i:4 p:1 r:2 f:-",
				"i:5 p:1 r:3 f:-",
				"i:2 p:1 r:4 f:-",
			}

			assert.Equal(t, expected, prettifyAll(results))
		})
	})

	t.Run("ByAttempts", func(t *testing.T) {
		competitors := makeCompetitors(4)
		competitors[0].Tops = 1
		competitors[0].AttemptsTops = 2
		competitors[1].Tops = 1
		competitors[1].AttemptsTops = 2
		competitors[2].Tops = 1
		competitors[2].AttemptsTops = 3

		ranker := scores.NewBasicRanker(1, false, true, true)

		results := ranker.RankCompetitors(slices.Values(competitors))

		expected := []string{
			"i:1 p:1 r:0 f:🏆",
			"i:2 p:1 r:1 f:🏆",
			"i:3 p:3 r:2 f:-",
			"i:4 p:4 r:3 f:-",
		}

		assert.Equal(t, expected, prettifyAll(results))
	})
}

func shuffleSlice[T any](slice []T) {
	for i := range slice {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}

func prettify(result domain.Result) string {
	finalist := "-"
	if result.Finalist {
		finalist = "🏆"
	}
	return fmt.Sprintf("i:%v p:%d r:%d f:%v", result.CompetitorID, result.Placement, result.RankOrder, finalist)
}

func prettifyAll(results []domain.Result) []string {
	arr := make([]string, 0)

	for result := range slices.Values(results) {
		arr = append(arr, prettify(result))
	}

	return arr
}

func TestBasicRankerScoreFormat(t *testing.T) {
	for _, tt := range []struct {
		name         string
		usePoints    bool
		zone1Enabled bool
		zone2Enabled bool
		expected     string
	}{
		{name: "NoZones", expected: "2t"},
		{name: "Zone1", zone1Enabled: true, expected: "2t 4z₁"},
		{name: "Zone2", zone1Enabled: true, zone2Enabled: true, expected: "2t 3z₂ 4z₁"},
		{name: "Points", usePoints: true, zone1Enabled: true, zone2Enabled: true, expected: "100p"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ranker := scores.NewBasicRanker(0, tt.usePoints, tt.zone1Enabled, tt.zone2Enabled)
			competitors := []scores.Competitor{{
				ID:    1,
				Score: scores.Score{Tops: 2, Zone2s: 3, Zone1s: 4, Points: 1000},
			}}
			ranked := ranker.RankCompetitors(slices.Values(competitors))
			assert.Equal(t, tt.expected, ranked[0].Score)
		})
	}
}
