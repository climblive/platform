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

func TestBasicRanker(t *testing.T) {
	ranker := scores.NewBasicRanker(5)

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

		scores := ranker.RankContenders(slices.Values(contenders))

		expected := []string{
			"i:1 p:1 r:0 f:-",
			"i:2 p:1 r:1 f:-",
			"i:3 p:1 r:2 f:-",
		}

		assert.Equal(t, expected, prettifyAll(scores))
	})

	t.Run("Simple", func(t *testing.T) {
		contenders := makeContenders(3)
		contenders[0].Score = 300
		contenders[1].Score = 200
		contenders[2].Score = 100

		shuffleSlice(contenders)

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
		contenders[0].Score = 300
		contenders[1].Score = 200
		contenders[2].Score = 200
		contenders[3].Score = 200
		contenders[4].Score = 100

		shuffleSlice(contenders)

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

	t.Run("ExtraFinalists", func(t *testing.T) {
		contenders := makeContenders(10)
		contenders[0].Score = 500
		contenders[1].Score = 400
		contenders[2].Score = 300
		contenders[3].Score = 200
		contenders[4].Score = 100
		contenders[5].Score = 100
		contenders[6].Score = 100
		contenders[7].Score = 50
		contenders[8].Score = 50
		contenders[9].Score = 50

		shuffleSlice(contenders)

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
		contenders[0].Score = 500
		contenders[1].Score = 400
		contenders[2].Score = 300
		contenders[3].Score = 200
		contenders[4].Score = 100
		contenders[5].Score = 100
		contenders[6].Score = 100
		contenders[7].Score = 50
		contenders[8].Score = 50
		contenders[9].Score = 0

		contenders[1].WithdrawnFromFinals = true
		contenders[2].WithdrawnFromFinals = true
		contenders[5].WithdrawnFromFinals = true

		shuffleSlice(contenders)

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
		contenders[0].Score = 0
		contenders[1].Score = 0
		contenders[2].Score = 0
		contenders[3].Score = 0
		contenders[4].Score = 0

		contenders[1].Disqualified = true

		shuffleSlice(contenders)

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
