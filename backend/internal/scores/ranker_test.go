package scores_test

import (
	"fmt"
	"math/rand"
	"slices"
	"testing"
	"testing/synctest"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/scores"
	"github.com/stretchr/testify/assert"
)

func TestBasicRanker(t *testing.T) {
	makeContenders := func(count int) []scores.Contender {
		contenders := make([]scores.Contender, count)

		for n := range count {
			contenderID := domain.ContenderID(n + 1)

			contenders[n] = scores.Contender{
				ID: contenderID,
			}
		}

		return contenders
	}

	t.Run("ContestNotStarted", func(t *testing.T) {
		contenders := makeContenders(3)
		shuffleSlice(contenders)

		ranker := scores.NewBasicRanker(5, true, false, false)

		scores := ranker.RankContenders(slices.Values(contenders))

		expected := []string{
			"i:1 p:1 r:0 f:-",
			"i:2 p:1 r:1 f:-",
			"i:3 p:1 r:2 f:-",
		}

		assert.Equal(t, expected, prettifyAll(scores))
	})

	t.Run("ByPoints", func(t *testing.T) {
		t.Run("Simple", func(t *testing.T) {
			contenders := makeContenders(3)
			contenders[0].Points = 300
			contenders[1].Points = 200
			contenders[2].Points = 100

			shuffleSlice(contenders)

			ranker := scores.NewBasicRanker(5, true, false, false)

			scores := ranker.RankContenders(slices.Values(contenders))

			expected := []string{
				"i:1 p:1 r:0 f:🏆",
				"i:2 p:2 r:1 f:🏆",
				"i:3 p:3 r:2 f:🏆",
			}

			assert.Equal(t, expected, prettifyAll(scores))
		})

		t.Run("SharedPlacement", func(t *testing.T) {
			contenders := makeContenders(5)
			contenders[0].Points = 300
			contenders[1].Points = 200
			contenders[2].Points = 200
			contenders[3].Points = 200
			contenders[4].Points = 100

			shuffleSlice(contenders)

			ranker := scores.NewBasicRanker(5, true, false, false)

			scores := ranker.RankContenders(slices.Values(contenders))

			expected := []string{
				"i:1 p:1 r:0 f:🏆",
				"i:2 p:2 r:1 f:🏆",
				"i:3 p:2 r:2 f:🏆",
				"i:4 p:2 r:3 f:🏆",
				"i:5 p:5 r:4 f:🏆",
			}

			assert.Equal(t, expected, prettifyAll(scores))
		})

		t.Run("IgnoreAttempts", func(t *testing.T) {
			contenders := makeContenders(2)
			contenders[0].Points = 100
			contenders[0].Tops = 1
			contenders[0].AttemptsTops = 2
			contenders[1].Points = 100

			ranker := scores.NewBasicRanker(5, true, false, false)

			scores := ranker.RankContenders(slices.Values(contenders))

			expected := []string{
				"i:1 p:1 r:0 f:🏆",
				"i:2 p:1 r:1 f:🏆",
			}

			assert.Equal(t, expected, prettifyAll(scores))
		})

		t.Run("ExtraFinalists", func(t *testing.T) {
			contenders := makeContenders(10)
			contenders[0].Points = 500
			contenders[1].Points = 400
			contenders[2].Points = 300
			contenders[3].Points = 200
			contenders[4].Points = 100
			contenders[5].Points = 100
			contenders[6].Points = 100
			contenders[7].Points = 50
			contenders[8].Points = 50
			contenders[9].Points = 50

			shuffleSlice(contenders)

			ranker := scores.NewBasicRanker(5, true, false, false)

			scores := ranker.RankContenders(slices.Values(contenders))

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

			assert.Equal(t, expected, prettifyAll(scores))
		})

		t.Run("WithdrawalsFromFinals", func(t *testing.T) {
			contenders := makeContenders(10)
			contenders[0].Points = 500
			contenders[1].Points = 400
			contenders[2].Points = 300
			contenders[3].Points = 200
			contenders[4].Points = 100
			contenders[5].Points = 100
			contenders[6].Points = 100
			contenders[7].Points = 50
			contenders[8].Points = 50
			contenders[9].Points = 0

			contenders[1].WithdrawnFromFinals = true
			contenders[2].WithdrawnFromFinals = true
			contenders[5].WithdrawnFromFinals = true

			shuffleSlice(contenders)

			ranker := scores.NewBasicRanker(5, true, false, false)

			scores := ranker.RankContenders(slices.Values(contenders))

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

			assert.Equal(t, expected, prettifyAll(scores))
		})

		t.Run("DisqualifiedContendersLast", func(t *testing.T) {
			contenders := makeContenders(5)
			contenders[0].Points = 0
			contenders[1].Points = 0
			contenders[2].Points = 0
			contenders[3].Points = 0
			contenders[4].Points = 0

			contenders[1].Disqualified = true

			shuffleSlice(contenders)

			ranker := scores.NewBasicRanker(5, true, false, false)

			scores := ranker.RankContenders(slices.Values(contenders))

			expected := []string{
				"i:1 p:1 r:0 f:-",
				"i:3 p:1 r:1 f:-",
				"i:4 p:1 r:2 f:-",
				"i:5 p:1 r:3 f:-",
				"i:2 p:1 r:4 f:-",
			}

			assert.Equal(t, expected, prettifyAll(scores))
		})
	})

	t.Run("ByAttempts", func(t *testing.T) {
		contenders := makeContenders(4)
		contenders[0].Tops = 1
		contenders[0].AttemptsTops = 2
		contenders[1].Tops = 1
		contenders[1].AttemptsTops = 2
		contenders[2].Tops = 1
		contenders[2].AttemptsTops = 3

		ranker := scores.NewBasicRanker(1, false, true, true)

		rankedScores := ranker.RankContenders(slices.Values(contenders))

		expected := []string{
			"i:1 p:1 r:0 f:🏆",
			"i:2 p:1 r:1 f:🏆",
			"i:3 p:3 r:2 f:-",
			"i:4 p:4 r:3 f:-",
		}

		assert.Equal(t, expected, prettifyAll(rankedScores))
	})
}

func shuffleSlice[T any](slice []T) {
	for i := range slice {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}

func prettify(score domain.Score) string {
	finalist := "-"
	if score.Finalist {
		finalist = "🏆"
	}
	return fmt.Sprintf("i:%v p:%d r:%d f:%v", score.ContenderID, score.Placement, score.RankOrder, finalist)
}

func prettifyAll(scores []domain.Score) []string {
	arr := make([]string, 0)

	for score := range slices.Values(scores) {
		arr = append(arr, prettify(score))
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
		{name: "Zone2", zone2Enabled: true, expected: "2t 3z₂"},
		{name: "BothZones", zone1Enabled: true, zone2Enabled: true, expected: "2t 3z₂ 4z₁"},
		{name: "Points", usePoints: true, zone1Enabled: true, zone2Enabled: true, expected: "100p"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				ranker := scores.NewBasicRanker(0, tt.usePoints, tt.zone1Enabled, tt.zone2Enabled)
				contenders := []scores.Contender{{
					ID:    1,
					Score: scores.Score{Tops: 2, Zone2s: 3, Zone1s: 4, Points: 100},
				}}
				ranked := ranker.RankContenders(slices.Values(contenders))
				assert.Equal(t, tt.expected, ranked[0].Score)
			})
		})
	}
}
