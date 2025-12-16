package scores_test

import (
	"math/rand"
	"slices"
	"testing"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/scores"
	"github.com/stretchr/testify/assert"
)

func TestMemoryStore(t *testing.T) {
	t.Run("GetContender", func(t *testing.T) {
		store := scores.NewMemoryStore()

		contender := scores.Contender{
			ID:          domain.ContenderID(rand.Int()),
			CompClassID: domain.CompClassID(rand.Int()),
		}

		store.SaveContender(contender)

		t.Run("Found", func(t *testing.T) {
			result, found := store.GetContender(contender.ID)
			assert.True(t, found)
			assert.Equal(t, contender, result)
		})

		t.Run("NotFound", func(t *testing.T) {
			result, found := store.GetContender(contender.ID + 1)
			assert.False(t, found)
			assert.Empty(t, result)
		})
	})

	t.Run("GetContendersByCompClass", func(t *testing.T) {
		store := scores.NewMemoryStore()

		store.SaveContender(scores.Contender{
			ID:          1,
			CompClassID: 1,
		})
		store.SaveContender(scores.Contender{
			ID:          2,
			CompClassID: 1,
		})
		store.SaveContender(scores.Contender{
			ID:          3,
			CompClassID: 2,
		})
		store.SaveContender(scores.Contender{
			ID:          4,
			CompClassID: 3,
		})
		store.SaveContender(scores.Contender{
			ID:          5,
			CompClassID: 1,
		})
		store.SaveContender(scores.Contender{
			ID:          6,
			CompClassID: 2,
		})
		store.SaveContender(scores.Contender{
			ID:          7,
			CompClassID: 1,
		})

		t.Run("CompClassOne", func(t *testing.T) {
			var filteredContenders []domain.ContenderID
			for contender := range slices.Values(slices.Collect(store.GetContendersByCompClass(1))) {
				filteredContenders = append(filteredContenders, contender.ID)
			}

			assert.ElementsMatch(t, filteredContenders, []domain.ContenderID{1, 2, 5, 7})
		})

		t.Run("CompClassTwo", func(t *testing.T) {
			var filteredContenders []domain.ContenderID
			for contender := range slices.Values(slices.Collect(store.GetContendersByCompClass(2))) {
				filteredContenders = append(filteredContenders, contender.ID)
			}

			assert.ElementsMatch(t, filteredContenders, []domain.ContenderID{3, 6})
		})

		t.Run("CompClassThree", func(t *testing.T) {
			var filteredContenders []domain.ContenderID
			for contender := range slices.Values(slices.Collect(store.GetContendersByCompClass(3))) {
				filteredContenders = append(filteredContenders, contender.ID)
			}

			assert.ElementsMatch(t, filteredContenders, []domain.ContenderID{4})
		})

		t.Run("ExtractFirstItemOnly", func(t *testing.T) {
			for _ = range store.GetContendersByCompClass(1) {
				break
			}
		})
	})

	t.Run("GetAllContenders", func(t *testing.T) {
		store := scores.NewMemoryStore()

		c1 := scores.Contender{
			ID:          1,
			CompClassID: 1,
		}
		c2 := scores.Contender{
			ID:          2,
			CompClassID: 2,
		}
		c3 := scores.Contender{
			ID:          3,
			CompClassID: 3,
		}

		store.SaveContender(c1)
		store.SaveContender(c2)
		store.SaveContender(c3)

		assert.ElementsMatch(t, []scores.Contender{c1, c2, c3}, slices.Collect(store.GetAllContenders()))
	})

	t.Run("SaveContender", func(t *testing.T) {
		store := scores.NewMemoryStore()

		contender := scores.Contender{
			ID:                  domain.ContenderID(rand.Int()),
			CompClassID:         domain.CompClassID(rand.Int()),
			WithdrawnFromFinals: true,
			Disqualified:        true,
			Score:               123,
		}

		store.SaveContender(contender)

		result, found := store.GetContender(contender.ID)
		assert.True(t, found)
		assert.Equal(t, contender, result)
	})

	t.Run("GetCompClassIDs", func(t *testing.T) {
		store := scores.NewMemoryStore()

		store.SaveContender(scores.Contender{
			ID:          domain.ContenderID(rand.Int()),
			CompClassID: 1,
		})
		store.SaveContender(scores.Contender{
			ID:          domain.ContenderID(rand.Int()),
			CompClassID: 1,
		})
		store.SaveContender(scores.Contender{
			ID:          domain.ContenderID(rand.Int()),
			CompClassID: 2,
		})
		store.SaveContender(scores.Contender{
			ID:          domain.ContenderID(rand.Int()),
			CompClassID: 3,
		})
		store.SaveContender(scores.Contender{
			ID:          domain.ContenderID(rand.Int()),
			CompClassID: 1,
		})
		store.SaveContender(scores.Contender{
			ID:          domain.ContenderID(rand.Int()),
			CompClassID: 2,
		})
		store.SaveContender(scores.Contender{
			ID:          domain.ContenderID(rand.Int()),
			CompClassID: 1,
		})

		assert.ElementsMatch(t, []domain.CompClassID{1, 2, 3}, store.GetCompClassIDs())
	})

	t.Run("GetTicks", func(t *testing.T) {
		store := scores.NewMemoryStore()

		contenderID := domain.ContenderID(1)
		otherContenderID := domain.ContenderID(2)

		t1 := scores.Tick{
			ProblemID: 1,
		}
		t2 := scores.Tick{
			ProblemID: 2,
		}
		t3 := scores.Tick{
			ProblemID: 3,
		}
		t4 := scores.Tick{
			ProblemID: 4,
		}
		t5 := scores.Tick{
			ProblemID: 5,
		}

		store.SaveTick(contenderID, t1)
		store.SaveTick(contenderID, t2)
		store.SaveTick(contenderID, t3)

		store.SaveTick(otherContenderID, t4)
		store.SaveTick(otherContenderID, t5)

		assert.ElementsMatch(t, []scores.Tick{t1, t2, t3}, slices.Collect(store.GetTicks(contenderID)))
		assert.ElementsMatch(t, []scores.Tick{t4, t5}, slices.Collect(store.GetTicks(otherContenderID)))
	})

	t.Run("SaveTick", func(t *testing.T) {
		store := scores.NewMemoryStore()

		t1 := scores.Tick{
			ProblemID:     1,
			Top:           true,
			AttemptsTop:   7,
			Zone1:         true,
			AttemptsZone1: 2,
			Zone2:         true,
			AttemptsZone2: 3,
			Points:        1337,
		}

		store.SaveTick(1, t1)

		assert.ElementsMatch(t, []scores.Tick{t1}, slices.Collect(store.GetTicks(1)))

		t2 := scores.Tick{
			ProblemID:     2,
			Top:           false,
			AttemptsTop:   5,
			Zone1:         false,
			AttemptsZone1: 4,
			Zone2:         false,
			AttemptsZone2: 3,
			Points:        100,
		}

		store.SaveTick(1, t2)

		assert.ElementsMatch(t, []scores.Tick{t1, t2}, slices.Collect(store.GetTicks(1)))

		t2.Points = 123
		store.SaveTick(1, t2)

		assert.ElementsMatch(t, []scores.Tick{t1, t2}, slices.Collect(store.GetTicks(1)))
	})

	t.Run("DeleteTick", func(t *testing.T) {
		store := scores.NewMemoryStore()

		contenderID := domain.ContenderID(1)

		t1 := scores.Tick{
			ProblemID: 1,
		}
		t2 := scores.Tick{
			ProblemID: 2,
		}
		t3 := scores.Tick{
			ProblemID: 3,
		}

		store.SaveTick(contenderID, t1)
		store.SaveTick(contenderID, t2)
		store.SaveTick(contenderID, t3)

		assert.ElementsMatch(t, []scores.Tick{t1, t2, t3}, slices.Collect(store.GetTicks(contenderID)))

		store.DeleteTick(contenderID, t2.ProblemID)

		assert.ElementsMatch(t, []scores.Tick{t1, t3}, slices.Collect(store.GetTicks(contenderID)))
	})

	t.Run("GetProblem", func(t *testing.T) {
		store := scores.NewMemoryStore()

		problem := scores.Problem{
			ID: domain.ProblemID(rand.Int()),
		}

		store.SaveProblem(problem)

		t.Run("Found", func(t *testing.T) {
			result, found := store.GetProblem(problem.ID)
			assert.True(t, found)
			assert.Equal(t, problem, result)
		})

		t.Run("NotFound", func(t *testing.T) {
			result, found := store.GetProblem(problem.ID + 1)
			assert.False(t, found)
			assert.Empty(t, result)
		})
	})

	t.Run("SaveProblem", func(t *testing.T) {
		store := scores.NewMemoryStore()

		problem := scores.Problem{
			ID:          domain.ProblemID(rand.Int()),
			PointsTop:   200,
			PointsZone1: 100,
			PointsZone2: 150,
			FlashBonus:  25,
		}

		store.SaveProblem(problem)

		result, found := store.GetProblem(problem.ID)
		assert.True(t, found)
		assert.Equal(t, problem, result)
	})
}
