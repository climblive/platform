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

	t.Run("GetCompetitor", func(t *testing.T) {
		store := scores.NewMemoryStore()

		competitor := scores.Competitor{
			ID:         testutils.RandomResourceID[domain.CompetitorID](),
			CategoryID: testutils.RandomResourceID[domain.CategoryID](),
		}

		store.SaveCompetitor(competitor)

		t.Run("Found", func(t *testing.T) {
			result, found := store.GetCompetitor(competitor.ID)
			assert.True(t, found)
			assert.Equal(t, competitor, result)
		})

		t.Run("NotFound", func(t *testing.T) {
			result, found := store.GetCompetitor(competitor.ID + 1)
			assert.False(t, found)
			assert.Empty(t, result)
		})
	})

	t.Run("GetCompetitorsByCategory", func(t *testing.T) {
		store := scores.NewMemoryStore()

		store.SaveCompetitor(scores.Competitor{
			ID:         1,
			CategoryID: 1,
		})
		store.SaveCompetitor(scores.Competitor{
			ID:         2,
			CategoryID: 1,
		})
		store.SaveCompetitor(scores.Competitor{
			ID:         3,
			CategoryID: 2,
		})
		store.SaveCompetitor(scores.Competitor{
			ID:         4,
			CategoryID: 3,
		})
		store.SaveCompetitor(scores.Competitor{
			ID:         5,
			CategoryID: 1,
		})
		store.SaveCompetitor(scores.Competitor{
			ID:         6,
			CategoryID: 2,
		})
		store.SaveCompetitor(scores.Competitor{
			ID:         7,
			CategoryID: 1,
		})

		t.Run("CategoryOne", func(t *testing.T) {
			var filteredCompetitors []domain.CompetitorID
			for competitor := range slices.Values(slices.Collect(store.GetCompetitorsByCategory(1))) {
				filteredCompetitors = append(filteredCompetitors, competitor.ID)
			}

			assert.ElementsMatch(t, filteredCompetitors, []domain.CompetitorID{1, 2, 5, 7})
		})

		t.Run("CategoryTwo", func(t *testing.T) {
			var filteredCompetitors []domain.CompetitorID
			for competitor := range slices.Values(slices.Collect(store.GetCompetitorsByCategory(2))) {
				filteredCompetitors = append(filteredCompetitors, competitor.ID)
			}

			assert.ElementsMatch(t, filteredCompetitors, []domain.CompetitorID{3, 6})
		})

		t.Run("CategoryThree", func(t *testing.T) {
			var filteredCompetitors []domain.CompetitorID
			for competitor := range slices.Values(slices.Collect(store.GetCompetitorsByCategory(3))) {
				filteredCompetitors = append(filteredCompetitors, competitor.ID)
			}

			assert.ElementsMatch(t, filteredCompetitors, []domain.CompetitorID{4})
		})

		t.Run("ExtractFirstItemOnly", func(t *testing.T) {
			for _ = range store.GetCompetitorsByCategory(1) {
				break
			}
		})
	})

	t.Run("GetAllCompetitors", func(t *testing.T) {
		store := scores.NewMemoryStore()

		c1 := scores.Competitor{
			ID:         1,
			CategoryID: 1,
		}
		c2 := scores.Competitor{
			ID:         2,
			CategoryID: 2,
		}
		c3 := scores.Competitor{
			ID:         3,
			CategoryID: 3,
		}

		store.SaveCompetitor(c1)
		store.SaveCompetitor(c2)
		store.SaveCompetitor(c3)

		assert.ElementsMatch(t, []scores.Competitor{c1, c2, c3}, slices.Collect(store.GetAllCompetitors()))
	})

	t.Run("SaveCompetitor", func(t *testing.T) {
		store := scores.NewMemoryStore()

		competitor := scores.Competitor{
			ID:                  testutils.RandomResourceID[domain.CompetitorID](),
			CategoryID:          testutils.RandomResourceID[domain.CategoryID](),
			WithdrawnFromFinals: true,
			Disqualified:        true,
			Score: scores.Score{
				Points: 123,
			},
		}

		store.SaveCompetitor(competitor)

		result, found := store.GetCompetitor(competitor.ID)
		assert.True(t, found)
		assert.Equal(t, competitor, result)
	})

	t.Run("GetCategoryIDs", func(t *testing.T) {
		store := scores.NewMemoryStore()

		store.SaveCompetitor(scores.Competitor{
			ID:         testutils.RandomResourceID[domain.CompetitorID](),
			CategoryID: 1,
		})
		store.SaveCompetitor(scores.Competitor{
			ID:         testutils.RandomResourceID[domain.CompetitorID](),
			CategoryID: 1,
		})
		store.SaveCompetitor(scores.Competitor{
			ID:         testutils.RandomResourceID[domain.CompetitorID](),
			CategoryID: 2,
		})
		store.SaveCompetitor(scores.Competitor{
			ID:         testutils.RandomResourceID[domain.CompetitorID](),
			CategoryID: 3,
		})
		store.SaveCompetitor(scores.Competitor{
			ID:         testutils.RandomResourceID[domain.CompetitorID](),
			CategoryID: 1,
		})
		store.SaveCompetitor(scores.Competitor{
			ID:         testutils.RandomResourceID[domain.CompetitorID](),
			CategoryID: 2,
		})
		store.SaveCompetitor(scores.Competitor{
			ID:         testutils.RandomResourceID[domain.CompetitorID](),
			CategoryID: 1,
		})

		assert.ElementsMatch(t, []domain.CategoryID{1, 2, 3}, store.GetCategoryIDs())
	})

	t.Run("GetTicksByCompetitor", func(t *testing.T) {
		store := scores.NewMemoryStore()

		competitorID := domain.CompetitorID(1)
		otherCompetitorID := domain.CompetitorID(2)

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

		store.SaveTick(competitorID, t1)
		store.SaveTick(competitorID, t2)
		store.SaveTick(competitorID, t3)

		store.SaveTick(otherCompetitorID, t4)
		store.SaveTick(otherCompetitorID, t5)

		assert.ElementsMatch(t, []scores.Tick{t1, t2, t3}, slices.Collect(store.GetTicksByCompetitor(competitorID)))
		assert.ElementsMatch(t, []scores.Tick{t4, t5}, slices.Collect(store.GetTicksByCompetitor(otherCompetitorID)))
	})

	t.Run("GetTick", func(t *testing.T) {
		store := scores.NewMemoryStore()

		fakedCompetitorID := testutils.RandomResourceID[domain.CompetitorID]()
		fakedProblemID := testutils.RandomResourceID[domain.ProblemID]()

		fakedTick := scores.Tick{
			CompetitorID:  fakedCompetitorID,
			ProblemID:     fakedProblemID,
			Zone1:         true,
			AttemptsZone1: 1,
			Zone2:         true,
			AttemptsZone2: 2,
			Top:           true,
			AttemptsTop:   3,
		}

		store.SaveTick(fakedCompetitorID, fakedTick)

		tick, found := store.GetTick(fakedCompetitorID, fakedProblemID)

		assert.True(t, found)
		assert.Equal(t, fakedTick, tick)
	})

	t.Run("GetTick_NotFound", func(t *testing.T) {
		store := scores.NewMemoryStore()

		tick, found := store.GetTick(testutils.RandomResourceID[domain.CompetitorID](), testutils.RandomResourceID[domain.ProblemID]())

		assert.False(t, found)
		assert.Empty(t, tick)
	})

	t.Run("GetTicksByProblem", func(t *testing.T) {
		store := scores.NewMemoryStore()

		fakedCategoryID := testutils.RandomResourceID[domain.CategoryID]()
		fakedProblemID := testutils.RandomResourceID[domain.ProblemID]()

		fakedCompetitor1ID := testutils.RandomResourceID[domain.CompetitorID]()
		fakedCompetitor2ID := testutils.RandomResourceID[domain.CompetitorID]()
		fakedCompetitor3ID := testutils.RandomResourceID[domain.CompetitorID]()
		fakedCompetitor4ID := testutils.RandomResourceID[domain.CompetitorID]()

		t1 := scores.Tick{
			CompetitorID:  fakedCompetitor1ID,
			ProblemID:     fakedProblemID,
			Zone1:         true,
			AttemptsZone1: 1,
			Zone2:         true,
			AttemptsZone2: 2,
			Top:           true,
			AttemptsTop:   3,
		}
		t2 := scores.Tick{
			CompetitorID:  fakedCompetitor2ID,
			ProblemID:     fakedProblemID,
			Zone1:         false,
			AttemptsZone1: 4,
			Zone2:         false,
			AttemptsZone2: 5,
			Top:           false,
			AttemptsTop:   6,
		}

		store.SaveCompetitor(scores.Competitor{
			ID:         fakedCompetitor1ID,
			CategoryID: fakedCategoryID,
		})
		store.SaveCompetitor(scores.Competitor{
			ID:         fakedCompetitor2ID,
			CategoryID: fakedCategoryID,
		})
		store.SaveCompetitor(scores.Competitor{
			ID:         fakedCompetitor3ID,
			CategoryID: fakedCategoryID,
		})
		store.SaveCompetitor(scores.Competitor{
			ID:         fakedCompetitor4ID,
			CategoryID: testutils.RandomResourceID[domain.CategoryID](),
		})

		store.SaveTick(fakedCompetitor1ID, t1)
		store.SaveTick(fakedCompetitor2ID, t2)
		store.SaveTick(fakedCompetitor2ID, scores.Tick{
			CompetitorID:  fakedCompetitor2ID,
			ProblemID:     testutils.RandomResourceID[domain.ProblemID](),
			Zone1:         false,
			AttemptsZone1: 7,
			Zone2:         false,
			AttemptsZone2: 8,
			Top:           false,
			AttemptsTop:   9,
		})

		assert.ElementsMatch(t, []scores.Tick{t1, t2}, slices.Collect(store.GetTicksByProblem(fakedCategoryID, fakedProblemID)))
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

		assert.ElementsMatch(t, []scores.Tick{t1}, slices.Collect(store.GetTicksByCompetitor(1)))

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

		assert.ElementsMatch(t, []scores.Tick{t1, t2}, slices.Collect(store.GetTicksByCompetitor(1)))

		t2.AttemptsZone2 = 123
		store.SaveTick(1, t2)

		assert.ElementsMatch(t, []scores.Tick{t1, t2}, slices.Collect(store.GetTicksByCompetitor(1)))
	})

	t.Run("DeleteTick", func(t *testing.T) {
		store := scores.NewMemoryStore()

		competitorID := domain.CompetitorID(1)

		t1 := scores.Tick{
			ID:        10,
			ProblemID: 1,
		}
		t2 := scores.Tick{
			ID:        20,
			ProblemID: 2,
		}
		t3 := scores.Tick{
			ID:        30,
			ProblemID: 3,
		}

		store.SaveTick(competitorID, t1)
		store.SaveTick(competitorID, t2)
		store.SaveTick(competitorID, t3)

		assert.ElementsMatch(t, []scores.Tick{t1, t2, t3}, slices.Collect(store.GetTicksByCompetitor(competitorID)))

		store.DeleteTick(t2.ID)

		assert.ElementsMatch(t, []scores.Tick{t1, t3}, slices.Collect(store.GetTicksByCompetitor(competitorID)))
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
			ProblemValue: scores.ProblemValue{
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

		fakedCompetitorID := testutils.RandomResourceID[domain.CompetitorID]()
		fakedProblemID := testutils.RandomResourceID[domain.ProblemID]()

		pointValue := domain.PointValue{
			CompetitorID: fakedCompetitorID,
			ProblemID:    fakedProblemID,
			Current:      225,
			Zone1:        100,
			Zone2:        150,
			Top:          200,
			FlashBonus:   25,
		}

		store.SavePointValue(fakedCompetitorID, fakedProblemID, pointValue)

		t.Run("Found", func(t *testing.T) {
			result, found := store.GetPointValue(fakedCompetitorID, fakedProblemID)

			assert.True(t, found)
			assert.Equal(t, pointValue, result)
		})

		t.Run("NotFound", func(t *testing.T) {
			result, found := store.GetPointValue(testutils.RandomResourceID[domain.CompetitorID](), testutils.RandomResourceID[domain.ProblemID]())

			assert.False(t, found)
			assert.Empty(t, result)
		})
	})

	t.Run("GetDirtyPointValues", func(t *testing.T) {
		store := scores.NewMemoryStore()

		fakedCompetitorID := testutils.RandomResourceID[domain.CompetitorID]()
		fakedProblemID := testutils.RandomResourceID[domain.ProblemID]()

		pointValue := domain.PointValue{
			CompetitorID: fakedCompetitorID,
			ProblemID:    fakedProblemID,
			Current:      225,
			Zone1:        100,
			Zone2:        150,
			Top:          200,
			FlashBonus:   25,
		}

		store.SavePointValue(fakedCompetitorID, fakedProblemID, pointValue)

		dirtyPointValues := store.GetDirtyPointValues()

		assert.Len(t, dirtyPointValues, 1)
		assert.Equal(t, pointValue, dirtyPointValues[0])

		dirtyPointValues = store.GetDirtyPointValues()

		assert.Empty(t, dirtyPointValues)
	})

	t.Run("GetDirtyResults", func(t *testing.T) {
		store := scores.NewMemoryStore()

		fakedCompetitorID := testutils.RandomResourceID[domain.CompetitorID]()

		result := domain.Result{
			Timestamp:    time.Now(),
			CompetitorID: fakedCompetitorID,
			Score:        "1000p",
			Placement:    10,
			Finalist:     true,
			RankOrder:    9,
		}

		store.SaveResult(result)

		dirtyResults := store.GetDirtyResults()

		assert.Len(t, dirtyResults, 1)
		assert.Equal(t, result, dirtyResults[0])

		dirtyResults = store.GetDirtyResults()

		assert.Empty(t, dirtyResults)
	})
}
