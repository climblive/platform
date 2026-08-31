package scores_test

import (
	"slices"
	"testing"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/scores"
	"github.com/climblive/platform/backend/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestMemoryStore(t *testing.T) {
	t.Run("SaveRules", func(t *testing.T) {
		store := scores.NewMemoryStore()

		rules := scores.Rules{
			QualifyingProblems: 10,
			Finalists:          7,
		}

		store.SaveRules(rules)

		assert.Equal(t, rules, store.GetRules())
	})

	t.Run("GetContender", func(t *testing.T) {
		store := scores.NewMemoryStore()

		contender := scores.Contender{
			ID:          testutils.RandomResourceID[domain.ContenderID](),
			CompClassID: testutils.RandomResourceID[domain.CompClassID](),
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
			ID:                  testutils.RandomResourceID[domain.ContenderID](),
			CompClassID:         testutils.RandomResourceID[domain.CompClassID](),
			WithdrawnFromFinals: true,
			Disqualified:        true,
			Score: scores.Score{
				Points: 123,
			},
		}

		store.SaveContender(contender)

		result, found := store.GetContender(contender.ID)
		assert.True(t, found)
		assert.Equal(t, contender, result)
	})

	t.Run("GetCompClassIDs", func(t *testing.T) {
		store := scores.NewMemoryStore()

		store.SaveContender(scores.Contender{
			ID:          testutils.RandomResourceID[domain.ContenderID](),
			CompClassID: 1,
		})
		store.SaveContender(scores.Contender{
			ID:          testutils.RandomResourceID[domain.ContenderID](),
			CompClassID: 1,
		})
		store.SaveContender(scores.Contender{
			ID:          testutils.RandomResourceID[domain.ContenderID](),
			CompClassID: 2,
		})
		store.SaveContender(scores.Contender{
			ID:          testutils.RandomResourceID[domain.ContenderID](),
			CompClassID: 3,
		})
		store.SaveContender(scores.Contender{
			ID:          testutils.RandomResourceID[domain.ContenderID](),
			CompClassID: 1,
		})
		store.SaveContender(scores.Contender{
			ID:          testutils.RandomResourceID[domain.ContenderID](),
			CompClassID: 2,
		})
		store.SaveContender(scores.Contender{
			ID:          testutils.RandomResourceID[domain.ContenderID](),
			CompClassID: 1,
		})

		assert.ElementsMatch(t, []domain.CompClassID{1, 2, 3}, store.GetCompClassIDs())
	})

	t.Run("GetTicksByContender", func(t *testing.T) {
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

		assert.ElementsMatch(t, []scores.Tick{t1, t2, t3}, slices.Collect(store.GetTicksByContender(contenderID)))
		assert.ElementsMatch(t, []scores.Tick{t4, t5}, slices.Collect(store.GetTicksByContender(otherContenderID)))
	})

	t.Run("GetTick", func(t *testing.T) {
		store := scores.NewMemoryStore()

		fakedContenderID := testutils.RandomResourceID[domain.ContenderID]()
		fakedProblemID := testutils.RandomResourceID[domain.ProblemID]()

		fakedTick := scores.Tick{
			ContenderID:   fakedContenderID,
			ProblemID:     fakedProblemID,
			Zone1:         true,
			AttemptsZone1: 1,
			Zone2:         true,
			AttemptsZone2: 2,
			Top:           true,
			AttemptsTop:   3,
		}

		store.SaveTick(fakedContenderID, fakedTick)

		tick, found := store.GetTick(fakedContenderID, fakedProblemID)

		assert.True(t, found)
		assert.Equal(t, fakedTick, tick)
	})

	t.Run("GetTick_NotFound", func(t *testing.T) {
		store := scores.NewMemoryStore()

		tick, found := store.GetTick(testutils.RandomResourceID[domain.ContenderID](), testutils.RandomResourceID[domain.ProblemID]())

		assert.False(t, found)
		assert.Empty(t, tick)
	})

	t.Run("GetTicksByProblem", func(t *testing.T) {
		store := scores.NewMemoryStore()

		fakedCompClassID := testutils.RandomResourceID[domain.CompClassID]()
		fakedProblemID := testutils.RandomResourceID[domain.ProblemID]()

		fakedContender1ID := testutils.RandomResourceID[domain.ContenderID]()
		fakedContender2ID := testutils.RandomResourceID[domain.ContenderID]()
		fakedContender3ID := testutils.RandomResourceID[domain.ContenderID]()
		fakedContender4ID := testutils.RandomResourceID[domain.ContenderID]()

		t1 := scores.Tick{
			ContenderID:   fakedContender1ID,
			ProblemID:     fakedProblemID,
			Zone1:         true,
			AttemptsZone1: 1,
			Zone2:         true,
			AttemptsZone2: 2,
			Top:           true,
			AttemptsTop:   3,
		}
		t2 := scores.Tick{
			ContenderID:   fakedContender2ID,
			ProblemID:     fakedProblemID,
			Zone1:         false,
			AttemptsZone1: 4,
			Zone2:         false,
			AttemptsZone2: 5,
			Top:           false,
			AttemptsTop:   6,
		}

		store.SaveContender(scores.Contender{
			ID:          fakedContender1ID,
			CompClassID: fakedCompClassID,
		})
		store.SaveContender(scores.Contender{
			ID:          fakedContender2ID,
			CompClassID: fakedCompClassID,
		})
		store.SaveContender(scores.Contender{
			ID:          fakedContender3ID,
			CompClassID: fakedCompClassID,
		})
		store.SaveContender(scores.Contender{
			ID:          fakedContender4ID,
			CompClassID: testutils.RandomResourceID[domain.CompClassID](),
		})

		store.SaveTick(fakedContender1ID, t1)
		store.SaveTick(fakedContender2ID, t2)
		store.SaveTick(fakedContender2ID, scores.Tick{
			ContenderID:   fakedContender2ID,
			ProblemID:     testutils.RandomResourceID[domain.ProblemID](),
			Zone1:         false,
			AttemptsZone1: 7,
			Zone2:         false,
			AttemptsZone2: 8,
			Top:           false,
			AttemptsTop:   9,
		})

		assert.ElementsMatch(t, []scores.Tick{t1, t2}, slices.Collect(store.GetTicksByProblem(fakedCompClassID, fakedProblemID)))
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
		}

		store.SaveTick(1, t1)

		assert.ElementsMatch(t, []scores.Tick{t1}, slices.Collect(store.GetTicksByContender(1)))

		t2 := scores.Tick{
			ProblemID:     2,
			Top:           false,
			AttemptsTop:   5,
			Zone1:         false,
			AttemptsZone1: 4,
			Zone2:         false,
			AttemptsZone2: 3,
		}

		store.SaveTick(1, t2)

		assert.ElementsMatch(t, []scores.Tick{t1, t2}, slices.Collect(store.GetTicksByContender(1)))

		t2.AttemptsZone2 = 123
		store.SaveTick(1, t2)

		assert.ElementsMatch(t, []scores.Tick{t1, t2}, slices.Collect(store.GetTicksByContender(1)))
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

		assert.ElementsMatch(t, []scores.Tick{t1, t2, t3}, slices.Collect(store.GetTicksByContender(contenderID)))

		store.DeleteTick(contenderID, t2.ProblemID)

		assert.ElementsMatch(t, []scores.Tick{t1, t3}, slices.Collect(store.GetTicksByContender(contenderID)))
	})

	t.Run("GetProblem", func(t *testing.T) {
		store := scores.NewMemoryStore()

		problem := scores.Problem{
			ID: testutils.RandomResourceID[domain.ProblemID](),
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
			ID: testutils.RandomResourceID[domain.ProblemID](),
			ProblemValue: domain.ProblemValue{
				PointsTop:   200,
				PointsZone1: 100,
				PointsZone2: 150,
				FlashBonus:  25,
			},
		}

		store.SaveProblem(problem)

		result, found := store.GetProblem(problem.ID)
		assert.True(t, found)
		assert.Equal(t, problem, result)
	})

	t.Run("GetAllProblems", func(t *testing.T) {
		store := scores.NewMemoryStore()

		p1 := scores.Problem{
			ID: testutils.RandomResourceID[domain.ProblemID](),
		}
		p2 := scores.Problem{
			ID: testutils.RandomResourceID[domain.ProblemID](),
		}
		p3 := scores.Problem{
			ID: testutils.RandomResourceID[domain.ProblemID](),
		}

		store.SaveProblem(p1)
		store.SaveProblem(p2)
		store.SaveProblem(p3)

		assert.ElementsMatch(t, []scores.Problem{p1, p2, p3}, slices.Collect(store.GetAllProblems()))
	})

	t.Run("GetPointValue", func(t *testing.T) {
		store := scores.NewMemoryStore()

		fakedContenderID := testutils.RandomResourceID[domain.ContenderID]()
		fakedProblemID := testutils.RandomResourceID[domain.ProblemID]()

		pointValue := domain.PointValue{
			ContenderID: fakedContenderID,
			ProblemID:   fakedProblemID,
			Current:     225,
			Zone1:       100,
			Zone2:       150,
			Top:         200,
			FlashBonus:  25,
		}

		store.SavePointValue(fakedContenderID, fakedProblemID, pointValue)

		t.Run("Found", func(t *testing.T) {
			result, found := store.GetPointValue(fakedContenderID, fakedProblemID)

			assert.True(t, found)
			assert.Equal(t, pointValue, result)
		})

		t.Run("NotFound", func(t *testing.T) {
			result, found := store.GetPointValue(testutils.RandomResourceID[domain.ContenderID](), testutils.RandomResourceID[domain.ProblemID]())

			assert.False(t, found)
			assert.Empty(t, result)
		})
	})

	t.Run("GetDirtyPointValues", func(t *testing.T) {
		store := scores.NewMemoryStore()

		fakedContenderID := testutils.RandomResourceID[domain.ContenderID]()
		fakedProblemID := testutils.RandomResourceID[domain.ProblemID]()

		pointValue := domain.PointValue{
			ContenderID: fakedContenderID,
			ProblemID:   fakedProblemID,
			Current:     225,
			Zone1:       100,
			Zone2:       150,
			Top:         200,
			FlashBonus:  25,
		}

		store.SavePointValue(fakedContenderID, fakedProblemID, pointValue)

		dirtyPointValues := store.GetDirtyPointValues()

		assert.Len(t, dirtyPointValues, 1)
		assert.Equal(t, pointValue, dirtyPointValues[0])

		dirtyPointValues = store.GetDirtyPointValues()

		assert.Empty(t, dirtyPointValues)
	})

	t.Run("GetDirtyScores", func(t *testing.T) {
		store := scores.NewMemoryStore()

		fakedContenderID := testutils.RandomResourceID[domain.ContenderID]()

		score := domain.Score{
			Timestamp:   time.Now(),
			ContenderID: fakedContenderID,
			Score:       "1000p",
			Placement:   10,
			Finalist:    true,
			RankOrder:   9,
		}

		store.SaveScore(score)

		dirtyScores := store.GetDirtyScores()

		assert.Len(t, dirtyScores, 1)
		assert.Equal(t, score, dirtyScores[0])

		dirtyScores = store.GetDirtyScores()

		assert.Empty(t, dirtyScores)
	})
}
