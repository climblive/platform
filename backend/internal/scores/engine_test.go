package scores_test

import (
	"maps"
	"slices"
	"testing"
	"testing/synctest"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/scores"
	"github.com/climblive/platform/backend/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultScoreEngine(t *testing.T) {
	type fixture struct {
		store  *engineStoreMock
		engine *scores.DefaultScoreEngine
	}

	makeFixture := func() (fixture, func(t *testing.T)) {
		mockedStore := new(engineStoreMock)

		engine := scores.NewDefaultScoreEngine(mockedStore)

		awaitExpectations := func(t *testing.T) {
			mockedStore.AssertExpectations(t)
		}

		return fixture{
			store:  mockedStore,
			engine: engine,
		}, awaitExpectations
	}

	t.Run("HandleRulesUpdated", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedCompetitionID := testutils.RandomResourceID[domain.CompetitionID]()

		fakedCategory1ID := testutils.RandomResourceID[domain.CategoryID]()
		fakedCategory2ID := testutils.RandomResourceID[domain.CategoryID]()

		fakedProblem1ID := testutils.RandomResourceID[domain.ProblemID]()
		fakedProblem2ID := testutils.RandomResourceID[domain.ProblemID]()

		fakedCompetitor1ID := testutils.RandomResourceID[domain.CompetitorID]()
		fakedCompetitor2ID := testutils.RandomResourceID[domain.CompetitorID]()
		fakedCompetitor3ID := testutils.RandomResourceID[domain.CompetitorID]()
		fakedCompetitor4ID := testutils.RandomResourceID[domain.CompetitorID]()

		f.store.On("SaveRules", scores.Rules{
			QualifyingProblems: 10,
			Finalists:          7,
			UsePoints:          true,
			FavorTopAttempts:   true,
			PooledPoints:       true,
			MaxAttempts:        5,
			PointDeduction:     3,
		}).Return()

		f.store.On("GetCategoryIDs").Return([]domain.CategoryID{
			fakedCategory1ID,
			fakedCategory2ID,
		})

		f.store.On("GetAllProblems").Return(slices.Values([]scores.Problem{
			{
				ID: fakedProblem1ID,
			},
			{
				ID: fakedProblem2ID,
			},
		}))

		f.store.
			On("GetCompetitorsByCategory", fakedCategory1ID).Return(slices.Values([]scores.Competitor{
			{
				ID: fakedCompetitor1ID,
			},
			{
				ID: fakedCompetitor2ID,
			},
		})).
			On("GetCompetitorsByCategory", fakedCategory2ID).Return(slices.Values([]scores.Competitor{
			{
				ID: fakedCompetitor3ID,
			},
			{
				ID: fakedCompetitor4ID,
			},
		}))

		effects := slices.Collect(f.engine.HandleRulesUpdated(domain.RulesUpdatedEvent{
			CompetitionID:      fakedCompetitionID,
			QualifyingProblems: 10,
			Finalists:          7,
			UsePoints:          true,
			FavorTopAttempts:   true,
			PooledPoints:       true,
			MaxAttempts:        5,
			PointDeduction:     3,
		}))

		require.ElementsMatch(t, effects, []scores.Effect{
			scores.EffectCalculatePointValues{CategoryID: fakedCategory1ID, ProblemID: fakedProblem1ID},
			scores.EffectCalculatePointValues{CategoryID: fakedCategory1ID, ProblemID: fakedProblem2ID},
			scores.EffectCalculatePointValues{CategoryID: fakedCategory2ID, ProblemID: fakedProblem1ID},
			scores.EffectCalculatePointValues{CategoryID: fakedCategory2ID, ProblemID: fakedProblem2ID},

			scores.EffectScoreCompetitor{CompetitorID: fakedCompetitor1ID},
			scores.EffectScoreCompetitor{CompetitorID: fakedCompetitor2ID},
			scores.EffectScoreCompetitor{CompetitorID: fakedCompetitor3ID},
			scores.EffectScoreCompetitor{CompetitorID: fakedCompetitor4ID},

			scores.EffectRankCategory{CategoryID: fakedCategory1ID},
			scores.EffectRankCategory{CategoryID: fakedCategory2ID},
		})

		awaitExpectations(t)
	})

	t.Run("Start", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedCategory1ID := testutils.RandomResourceID[domain.CategoryID]()
		fakedCategory2ID := testutils.RandomResourceID[domain.CategoryID]()

		fakedProblem1ID := testutils.RandomResourceID[domain.ProblemID]()
		fakedProblem2ID := testutils.RandomResourceID[domain.ProblemID]()

		fakedCompetitor1ID := testutils.RandomResourceID[domain.CompetitorID]()
		fakedCompetitor2ID := testutils.RandomResourceID[domain.CompetitorID]()
		fakedCompetitor3ID := testutils.RandomResourceID[domain.CompetitorID]()
		fakedCompetitor4ID := testutils.RandomResourceID[domain.CompetitorID]()

		f.store.On("GetCategoryIDs").Return([]domain.CategoryID{
			fakedCategory1ID,
			fakedCategory2ID,
		})

		f.store.On("GetAllProblems").Return(slices.Values([]scores.Problem{
			{
				ID: fakedProblem1ID,
			},
			{
				ID: fakedProblem2ID,
			},
		}))

		f.store.
			On("GetCompetitorsByCategory", fakedCategory1ID).Return(slices.Values([]scores.Competitor{
			{
				ID: fakedCompetitor1ID,
			},
			{
				ID: fakedCompetitor2ID,
			},
		})).
			On("GetCompetitorsByCategory", fakedCategory2ID).Return(slices.Values([]scores.Competitor{
			{
				ID: fakedCompetitor3ID,
			},
			{
				ID: fakedCompetitor4ID,
			},
		}))

		effects := slices.Collect(f.engine.Start())

		require.ElementsMatch(t, effects, []scores.Effect{
			scores.EffectCalculatePointValues{CategoryID: fakedCategory1ID, ProblemID: fakedProblem1ID},
			scores.EffectCalculatePointValues{CategoryID: fakedCategory1ID, ProblemID: fakedProblem2ID},
			scores.EffectCalculatePointValues{CategoryID: fakedCategory2ID, ProblemID: fakedProblem1ID},
			scores.EffectCalculatePointValues{CategoryID: fakedCategory2ID, ProblemID: fakedProblem2ID},

			scores.EffectScoreCompetitor{CompetitorID: fakedCompetitor1ID},
			scores.EffectScoreCompetitor{CompetitorID: fakedCompetitor2ID},
			scores.EffectScoreCompetitor{CompetitorID: fakedCompetitor3ID},
			scores.EffectScoreCompetitor{CompetitorID: fakedCompetitor4ID},

			scores.EffectRankCategory{CategoryID: fakedCategory1ID},
			scores.EffectRankCategory{CategoryID: fakedCategory2ID},
		})

		awaitExpectations(t)
	})

	t.Run("CompetitorEntered", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedCompetitorID := testutils.RandomResourceID[domain.CompetitorID]()
		fakedCategoryID := testutils.RandomResourceID[domain.CategoryID]()

		fakedProblem1ID := testutils.RandomResourceID[domain.ProblemID]()
		fakedProblem2ID := testutils.RandomResourceID[domain.ProblemID]()
		fakedProblem3ID := testutils.RandomResourceID[domain.ProblemID]()

		f.store.On("SaveCompetitor", scores.Competitor{
			ID:         fakedCompetitorID,
			CategoryID: fakedCategoryID,
		}).Return()

		f.store.
			On("GetAllProblems").
			Return(slices.Values([]scores.Problem{
				{
					ID: fakedProblem1ID,
				},
				{
					ID: fakedProblem2ID,
				},
				{
					ID: fakedProblem3ID,
				},
			}))

		effects := slices.Collect(f.engine.HandleCompetitorEntered(domain.CompetitorEnteredEvent{
			CompetitorID: fakedCompetitorID,
			CategoryID:   fakedCategoryID,
		}))

		require.ElementsMatch(t, effects, []scores.Effect{
			scores.EffectCalculatePointValues{CategoryID: fakedCategoryID, ProblemID: fakedProblem1ID},
			scores.EffectCalculatePointValues{CategoryID: fakedCategoryID, ProblemID: fakedProblem2ID},
			scores.EffectCalculatePointValues{CategoryID: fakedCategoryID, ProblemID: fakedProblem3ID},
			scores.EffectRankCategory{CategoryID: fakedCategoryID},
		})

		awaitExpectations(t)
	})

	t.Run("CompetitorSwitchedCategory_CompetitorNotFound", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedCompetitorID := testutils.RandomResourceID[domain.CompetitorID]()
		fakedCategoryID := testutils.RandomResourceID[domain.CategoryID]()

		f.store.
			On("GetCompetitor", fakedCompetitorID).
			Return(scores.Competitor{}, false)

		effects := f.engine.HandleCompetitorSwitchedCategory(
			domain.CompetitorSwitchedCategoryEvent{
				CompetitorID: fakedCompetitorID,
				CategoryID:   fakedCategoryID,
			})

		assert.Empty(t, effects)

		awaitExpectations(t)
	})

	t.Run("CompetitorSwitchedCategory_SameCategory", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedCompetitorID := testutils.RandomResourceID[domain.CompetitorID]()
		fakedCategoryID := testutils.RandomResourceID[domain.CategoryID]()

		f.store.
			On("GetCompetitor", fakedCompetitorID).
			Return(scores.Competitor{
				ID:         fakedCompetitorID,
				CategoryID: fakedCategoryID,
			}, true)

		effects := f.engine.HandleCompetitorSwitchedCategory(domain.CompetitorSwitchedCategoryEvent{
			CompetitorID: fakedCompetitorID,
			CategoryID:   fakedCategoryID,
		})

		assert.Empty(t, effects)

		awaitExpectations(t)
	})

	t.Run("CompetitorSwitchedCategory", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedCompetitorID := testutils.RandomResourceID[domain.CompetitorID]()
		fakedOldCategoryID := testutils.RandomResourceID[domain.CategoryID]()
		fakedNewCategoryID := testutils.RandomResourceID[domain.CategoryID]()

		fakedProblem1ID := testutils.RandomResourceID[domain.ProblemID]()
		fakedProblem2ID := testutils.RandomResourceID[domain.ProblemID]()
		fakedProblem3ID := testutils.RandomResourceID[domain.ProblemID]()

		f.store.
			On("GetCompetitor", fakedCompetitorID).
			Return(scores.Competitor{
				ID:                  fakedCompetitorID,
				CategoryID:          fakedOldCategoryID,
				Disqualified:        false,
				WithdrawnFromFinals: false,
				Score: scores.Score{
					Points: 123,
				},
			}, true)

		f.store.On("SaveCompetitor", scores.Competitor{
			ID:                  fakedCompetitorID,
			CategoryID:          fakedNewCategoryID,
			Disqualified:        false,
			WithdrawnFromFinals: false,
			Score: scores.Score{
				Points: 123,
			},
		}).Return()

		f.store.
			On("GetAllProblems").
			Return(slices.Values([]scores.Problem{
				{
					ID: fakedProblem1ID,
				},
				{
					ID: fakedProblem2ID,
				},
				{
					ID: fakedProblem3ID,
				},
			}))

		effects := slices.Collect(f.engine.HandleCompetitorSwitchedCategory(domain.CompetitorSwitchedCategoryEvent{
			CompetitorID: fakedCompetitorID,
			CategoryID:   fakedNewCategoryID,
		}))

		require.ElementsMatch(t, effects, []scores.Effect{
			scores.EffectCalculatePointValues{CategoryID: fakedOldCategoryID, ProblemID: fakedProblem1ID},
			scores.EffectCalculatePointValues{CategoryID: fakedOldCategoryID, ProblemID: fakedProblem2ID},
			scores.EffectCalculatePointValues{CategoryID: fakedOldCategoryID, ProblemID: fakedProblem3ID},

			scores.EffectCalculatePointValues{CategoryID: fakedNewCategoryID, ProblemID: fakedProblem1ID},
			scores.EffectCalculatePointValues{CategoryID: fakedNewCategoryID, ProblemID: fakedProblem2ID},
			scores.EffectCalculatePointValues{CategoryID: fakedNewCategoryID, ProblemID: fakedProblem3ID},

			scores.EffectRankCategory{CategoryID: fakedOldCategoryID},
			scores.EffectRankCategory{CategoryID: fakedNewCategoryID},
		})

		awaitExpectations(t)
	})

	t.Run("CompetitorWithdrewFromFinals_CompetitorNotFound", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedCompetitorID := testutils.RandomResourceID[domain.CompetitorID]()

		f.store.
			On("GetCompetitor", fakedCompetitorID).
			Return(scores.Competitor{}, false)

		effects := f.engine.HandleCompetitorWithdrewFromFinals(domain.CompetitorWithdrewFromFinalsEvent{
			CompetitorID: fakedCompetitorID,
		})

		assert.Empty(t, effects)

		awaitExpectations(t)
	})

	t.Run("CompetitorWithdrewFromFinals", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedCompetitorID := testutils.RandomResourceID[domain.CompetitorID]()
		fakedCategoryID := testutils.RandomResourceID[domain.CategoryID]()

		f.store.
			On("GetCompetitor", fakedCompetitorID).
			Return(scores.Competitor{
				ID:           fakedCompetitorID,
				CategoryID:   fakedCategoryID,
				Disqualified: true,
				Score: scores.Score{
					Points: 123,
				},
			}, true)

		f.store.On("SaveCompetitor", scores.Competitor{
			ID:                  fakedCompetitorID,
			CategoryID:          fakedCategoryID,
			Disqualified:        true,
			WithdrawnFromFinals: true,
			Score: scores.Score{
				Points: 123,
			},
		}).Return()

		effects := slices.Collect(f.engine.HandleCompetitorWithdrewFromFinals(domain.CompetitorWithdrewFromFinalsEvent{
			CompetitorID: fakedCompetitorID,
		}))

		require.ElementsMatch(t, effects, []scores.Effect{
			scores.EffectRankCategory{CategoryID: fakedCategoryID},
		})

		awaitExpectations(t)
	})

	t.Run("CompetitorReenteredFinals_CompetitorNotFound", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedCompetitorID := testutils.RandomResourceID[domain.CompetitorID]()

		f.store.
			On("GetCompetitor", fakedCompetitorID).
			Return(scores.Competitor{}, false)

		effects := f.engine.HandleCompetitorReenteredFinals(domain.CompetitorReenteredFinalsEvent{
			CompetitorID: fakedCompetitorID,
		})

		assert.Empty(t, effects)

		awaitExpectations(t)
	})

	t.Run("CompetitorReenteredFinals", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedCompetitorID := testutils.RandomResourceID[domain.CompetitorID]()
		fakedCategoryID := testutils.RandomResourceID[domain.CategoryID]()

		f.store.
			On("GetCompetitor", fakedCompetitorID).
			Return(scores.Competitor{
				ID:                  fakedCompetitorID,
				CategoryID:          fakedCategoryID,
				Disqualified:        true,
				WithdrawnFromFinals: true,
				Score: scores.Score{
					Points: 123,
				},
			}, true)

		f.store.On("SaveCompetitor", scores.Competitor{
			ID:                  fakedCompetitorID,
			CategoryID:          fakedCategoryID,
			Disqualified:        true,
			WithdrawnFromFinals: false,
			Score: scores.Score{
				Points: 123,
			},
		}).Return()

		effects := slices.Collect(f.engine.HandleCompetitorReenteredFinals(domain.CompetitorReenteredFinalsEvent{
			CompetitorID: fakedCompetitorID,
		}))

		require.ElementsMatch(t, effects, []scores.Effect{
			scores.EffectRankCategory{CategoryID: fakedCategoryID},
		})

		awaitExpectations(t)
	})

	t.Run("CompetitorDisqualified_CompetitorNotFound", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedCompetitorID := testutils.RandomResourceID[domain.CompetitorID]()

		f.store.
			On("GetCompetitor", fakedCompetitorID).
			Return(scores.Competitor{}, false)

		effects := f.engine.HandleCompetitorDisqualified(domain.CompetitorDisqualifiedEvent{
			CompetitorID: fakedCompetitorID,
		})

		assert.Empty(t, effects)

		awaitExpectations(t)
	})

	t.Run("CompetitorDisqualified", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedCompetitorID := testutils.RandomResourceID[domain.CompetitorID]()
		fakedCategoryID := testutils.RandomResourceID[domain.CategoryID]()

		fakedProblem1ID := testutils.RandomResourceID[domain.ProblemID]()
		fakedProblem2ID := testutils.RandomResourceID[domain.ProblemID]()
		fakedProblem3ID := testutils.RandomResourceID[domain.ProblemID]()

		f.store.
			On("GetCompetitor", fakedCompetitorID).
			Return(scores.Competitor{
				ID:                  fakedCompetitorID,
				CategoryID:          fakedCategoryID,
				WithdrawnFromFinals: true,
				Score: scores.Score{
					Points: 123,
				},
			}, true)

		f.store.
			On("GetAllProblems").
			Return(slices.Values([]scores.Problem{
				{
					ID: fakedProblem1ID,
				},
				{
					ID: fakedProblem2ID,
				},
				{
					ID: fakedProblem3ID,
				},
			}))

		f.store.On("SaveCompetitor", scores.Competitor{
			ID:                  fakedCompetitorID,
			CategoryID:          fakedCategoryID,
			Disqualified:        true,
			WithdrawnFromFinals: true,
			Score: scores.Score{
				Points: 123,
			},
		}).Return()

		effects := slices.Collect(f.engine.HandleCompetitorDisqualified(domain.CompetitorDisqualifiedEvent{
			CompetitorID: fakedCompetitorID,
		}))

		require.ElementsMatch(t, effects, []scores.Effect{
			scores.EffectCalculatePointValues{CategoryID: fakedCategoryID, ProblemID: fakedProblem1ID},
			scores.EffectCalculatePointValues{CategoryID: fakedCategoryID, ProblemID: fakedProblem2ID},
			scores.EffectCalculatePointValues{CategoryID: fakedCategoryID, ProblemID: fakedProblem3ID},
			scores.EffectScoreCompetitor{CompetitorID: fakedCompetitorID},
			scores.EffectRankCategory{CategoryID: fakedCategoryID},
		})

		awaitExpectations(t)
	})

	t.Run("CompetitorRequalified_CompetitorNotFound", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedCompetitorID := testutils.RandomResourceID[domain.CompetitorID]()

		f.store.
			On("GetCompetitor", fakedCompetitorID).
			Return(scores.Competitor{}, false)

		effects := f.engine.HandleCompetitorRequalified(domain.CompetitorRequalifiedEvent{
			CompetitorID: fakedCompetitorID,
		})

		assert.Empty(t, effects)

		awaitExpectations(t)
	})

	t.Run("CompetitorRequalified", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedCompetitorID := testutils.RandomResourceID[domain.CompetitorID]()
		fakedCategoryID := testutils.RandomResourceID[domain.CategoryID]()

		fakedProblem1ID := testutils.RandomResourceID[domain.ProblemID]()
		fakedProblem2ID := testutils.RandomResourceID[domain.ProblemID]()
		fakedProblem3ID := testutils.RandomResourceID[domain.ProblemID]()

		f.store.
			On("GetCompetitor", fakedCompetitorID).
			Return(scores.Competitor{
				ID:                  fakedCompetitorID,
				CategoryID:          fakedCategoryID,
				Disqualified:        true,
				WithdrawnFromFinals: true,
				Score:               scores.Score{},
			}, true)

		f.store.
			On("GetAllProblems").
			Return(slices.Values([]scores.Problem{
				{
					ID: fakedProblem1ID,
				},
				{
					ID: fakedProblem2ID,
				},
				{
					ID: fakedProblem3ID,
				},
			}))

		f.store.On("SaveCompetitor", scores.Competitor{
			ID:                  fakedCompetitorID,
			CategoryID:          fakedCategoryID,
			Disqualified:        false,
			WithdrawnFromFinals: true,
			Score:               scores.Score{},
		}).Return()

		effects := slices.Collect(f.engine.HandleCompetitorRequalified(domain.CompetitorRequalifiedEvent{
			CompetitorID: fakedCompetitorID,
		}))

		require.ElementsMatch(t, effects, []scores.Effect{
			scores.EffectCalculatePointValues{CategoryID: fakedCategoryID, ProblemID: fakedProblem1ID},
			scores.EffectCalculatePointValues{CategoryID: fakedCategoryID, ProblemID: fakedProblem2ID},
			scores.EffectCalculatePointValues{CategoryID: fakedCategoryID, ProblemID: fakedProblem3ID},
			scores.EffectScoreCompetitor{CompetitorID: fakedCompetitorID},
			scores.EffectRankCategory{CategoryID: fakedCategoryID},
		})

		awaitExpectations(t)
	})

	t.Run("ProblemAdded", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedProblemID := testutils.RandomResourceID[domain.ProblemID]()

		fakedCategory1ID := testutils.RandomResourceID[domain.CategoryID]()
		fakedCategory2ID := testutils.RandomResourceID[domain.CategoryID]()
		fakedCategory3ID := testutils.RandomResourceID[domain.CategoryID]()

		f.store.
			On("SaveProblem", scores.Problem{
				ID:           fakedProblemID,
				Zone1Enabled: true,
				Zone2Enabled: true,
				ProblemValue: scores.ProblemValue{
					PointsTop:   1000,
					PointsZone1: 500,
					PointsZone2: 750,
					FlashBonus:  100,
				},
			}).
			Return()

		f.store.
			On("GetCategoryIDs").
			Return([]domain.CategoryID{
				fakedCategory1ID,
				fakedCategory2ID,
				fakedCategory3ID,
			})

		effects := slices.Collect(f.engine.HandleProblemAdded(domain.ProblemAddedEvent{
			ProblemID:    fakedProblemID,
			Zone1Enabled: true,
			Zone2Enabled: true,
			ProblemValue: domain.ProblemValue{
				PointsTop:   100,
				PointsZone1: 50,
				PointsZone2: 75,
				FlashBonus:  10,
			},
		}))

		require.ElementsMatch(t, effects, []scores.Effect{
			scores.EffectRankCategory{CategoryID: fakedCategory1ID},
			scores.EffectRankCategory{CategoryID: fakedCategory2ID},
			scores.EffectRankCategory{CategoryID: fakedCategory3ID},
			scores.EffectCalculatePointValues{CategoryID: fakedCategory1ID, ProblemID: fakedProblemID},
			scores.EffectCalculatePointValues{CategoryID: fakedCategory2ID, ProblemID: fakedProblemID},
			scores.EffectCalculatePointValues{CategoryID: fakedCategory3ID, ProblemID: fakedProblemID},
		})

		awaitExpectations(t)
	})

	t.Run("ProblemUpdated", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedProblemID := testutils.RandomResourceID[domain.ProblemID]()

		fakedCategory1ID := testutils.RandomResourceID[domain.CategoryID]()
		fakedCategory2ID := testutils.RandomResourceID[domain.CategoryID]()
		fakedCategory3ID := testutils.RandomResourceID[domain.CategoryID]()

		fakedCompetitor1ID := testutils.RandomResourceID[domain.CompetitorID]()
		fakedCompetitor2ID := testutils.RandomResourceID[domain.CompetitorID]()
		fakedCompetitor3ID := testutils.RandomResourceID[domain.CompetitorID]()

		f.store.
			On("SaveProblem", scores.Problem{
				ID:           fakedProblemID,
				Zone1Enabled: true,
				Zone2Enabled: true,
				ProblemValue: scores.ProblemValue{
					PointsTop:   1000,
					PointsZone1: 500,
					PointsZone2: 750,
					FlashBonus:  100,
				},
			}).
			Return()

		f.store.
			On("GetCategoryIDs").
			Return([]domain.CategoryID{
				fakedCategory1ID,
				fakedCategory2ID,
				fakedCategory3ID,
			})

		f.store.
			On("GetCompetitorsByCategory", fakedCategory1ID).
			Return(slices.Values([]scores.Competitor{
				{
					ID: fakedCompetitor1ID,
				},
			})).
			On("GetCompetitorsByCategory", fakedCategory2ID).
			Return(slices.Values([]scores.Competitor{
				{
					ID: fakedCompetitor2ID,
				},
			})).
			On("GetCompetitorsByCategory", fakedCategory3ID).
			Return(slices.Values([]scores.Competitor{
				{
					ID: fakedCompetitor3ID,
				},
			}))

		effects := slices.Collect(f.engine.HandleProblemUpdated(domain.ProblemUpdatedEvent{
			ProblemID:    fakedProblemID,
			Zone1Enabled: true,
			Zone2Enabled: true,
			ProblemValue: domain.ProblemValue{
				PointsTop:   100,
				PointsZone1: 50,
				PointsZone2: 75,
				FlashBonus:  10,
			},
		}))

		require.ElementsMatch(t, effects, []scores.Effect{
			scores.EffectRankCategory{CategoryID: fakedCategory1ID},
			scores.EffectRankCategory{CategoryID: fakedCategory2ID},
			scores.EffectRankCategory{CategoryID: fakedCategory3ID},
			scores.EffectCalculatePointValues{CategoryID: fakedCategory1ID, ProblemID: fakedProblemID},
			scores.EffectCalculatePointValues{CategoryID: fakedCategory2ID, ProblemID: fakedProblemID},
			scores.EffectCalculatePointValues{CategoryID: fakedCategory3ID, ProblemID: fakedProblemID},

			scores.EffectScoreCompetitor{CompetitorID: fakedCompetitor1ID},
			scores.EffectScoreCompetitor{CompetitorID: fakedCompetitor2ID},
			scores.EffectScoreCompetitor{CompetitorID: fakedCompetitor3ID},
		})

		awaitExpectations(t)
	})

	t.Run("AscentRegistered_CompetitorNotFound", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedCompetitorID := testutils.RandomResourceID[domain.CompetitorID]()
		fakedProblemID := testutils.RandomResourceID[domain.ProblemID]()

		f.store.
			On("GetTick", fakedCompetitorID, fakedProblemID).
			Return(scores.Tick{}, false)

		f.store.
			On("GetCompetitor", fakedCompetitorID).
			Return(scores.Competitor{}, false)

		f.engine.HandleAscentRegistered(domain.AscentRegisteredEvent{
			CompetitorID:  fakedCompetitorID,
			ProblemID:     fakedProblemID,
			Top:           true,
			AttemptsTop:   3,
			Zone1:         true,
			AttemptsZone1: 1,
			Zone2:         true,
			AttemptsZone2: 2,
		})

		awaitExpectations(t)
	})

	t.Run("AscentRegistered_Disqualified", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedTickID := testutils.RandomResourceID[domain.TickID]()
		fakedCompetitorID := testutils.RandomResourceID[domain.CompetitorID]()
		fakedCategoryID := testutils.RandomResourceID[domain.CategoryID]()
		fakedProblemID := testutils.RandomResourceID[domain.ProblemID]()

		f.store.
			On("GetTick", fakedCompetitorID, fakedProblemID).
			Return(scores.Tick{}, false)

		f.store.
			On("GetCompetitor", fakedCompetitorID).
			Return(scores.Competitor{
				ID:           fakedCompetitorID,
				CategoryID:   fakedCategoryID,
				Disqualified: true,
			}, true)

		f.store.
			On("SaveTick", fakedCompetitorID, scores.Tick{
				ID:            fakedTickID,
				Revision:      1,
				CompetitorID:  fakedCompetitorID,
				ProblemID:     fakedProblemID,
				Top:           true,
				AttemptsTop:   5,
				Zone1:         true,
				AttemptsZone1: 2,
				Zone2:         true,
				AttemptsZone2: 3,
			}).
			Return()

		effects := f.engine.HandleAscentRegistered(domain.AscentRegisteredEvent{
			TickID:        fakedTickID,
			Revision:      1,
			CompetitorID:  fakedCompetitorID,
			ProblemID:     fakedProblemID,
			Top:           true,
			AttemptsTop:   5,
			Zone1:         true,
			AttemptsZone1: 2,
			Zone2:         true,
			AttemptsZone2: 3,
		})

		require.Nil(t, effects)

		awaitExpectations(t)
	})

	t.Run("AscentRegistered_StaleRevision", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedCompetitorID := testutils.RandomResourceID[domain.CompetitorID]()
		fakedProblemID := testutils.RandomResourceID[domain.ProblemID]()

		f.store.
			On("GetTick", fakedCompetitorID, fakedProblemID).
			Return(scores.Tick{
				CompetitorID: fakedCompetitorID,
				ProblemID:    fakedProblemID,
				Revision:     2,
			}, true)

		effects := f.engine.HandleAscentRegistered(domain.AscentRegisteredEvent{
			CompetitorID: fakedCompetitorID,
			ProblemID:    fakedProblemID,
			Revision:     1,
		})

		require.Nil(t, effects)

		awaitExpectations(t)
	})

	t.Run("AscentRegistered", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedTickID := testutils.RandomResourceID[domain.TickID]()
		fakedCompetitorID := testutils.RandomResourceID[domain.CompetitorID]()
		fakedCategoryID := testutils.RandomResourceID[domain.CategoryID]()
		fakedProblemID := testutils.RandomResourceID[domain.ProblemID]()

		f.store.
			On("GetTick", fakedCompetitorID, fakedProblemID).
			Return(scores.Tick{}, false)

		f.store.
			On("GetCompetitor", fakedCompetitorID).
			Return(scores.Competitor{
				ID:         fakedCompetitorID,
				CategoryID: fakedCategoryID,
			}, true)

		f.store.
			On("SaveTick", fakedCompetitorID, scores.Tick{
				ID:            fakedTickID,
				Revision:      1,
				CompetitorID:  fakedCompetitorID,
				ProblemID:     fakedProblemID,
				Top:           true,
				AttemptsTop:   5,
				Zone1:         true,
				AttemptsZone1: 2,
				Zone2:         true,
				AttemptsZone2: 3,
			}).
			Return()

		effects := slices.Collect(f.engine.HandleAscentRegistered(domain.AscentRegisteredEvent{
			TickID:        fakedTickID,
			Revision:      1,
			CompetitorID:  fakedCompetitorID,
			ProblemID:     fakedProblemID,
			Top:           true,
			AttemptsTop:   5,
			Zone1:         true,
			AttemptsZone1: 2,
			Zone2:         true,
			AttemptsZone2: 3,
		}))

		require.ElementsMatch(t, effects, []scores.Effect{
			scores.EffectCalculatePointValues{CategoryID: fakedCategoryID, ProblemID: fakedProblemID},
			scores.EffectScoreCompetitor{CompetitorID: fakedCompetitorID},
		})

		awaitExpectations(t)
	})

	t.Run("AscentDeregistered_CompetitorNotFound", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedCompetitorID := testutils.RandomResourceID[domain.CompetitorID]()
		fakedProblemID := testutils.RandomResourceID[domain.ProblemID]()

		f.store.
			On("GetCompetitor", fakedCompetitorID).
			Return(scores.Competitor{}, false)

		f.engine.HandleAscentDeregistered(domain.AscentDeregisteredEvent{
			CompetitorID: fakedCompetitorID,
			ProblemID:    fakedProblemID,
		})

		awaitExpectations(t)
	})

	t.Run("AscentDeregistered_Disqualified", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedTickID := testutils.RandomResourceID[domain.TickID]()
		fakedCompetitorID := testutils.RandomResourceID[domain.CompetitorID]()
		fakedCategoryID := testutils.RandomResourceID[domain.CategoryID]()
		fakedProblemID := testutils.RandomResourceID[domain.ProblemID]()

		f.store.
			On("GetCompetitor", fakedCompetitorID).
			Return(scores.Competitor{
				ID:           fakedCompetitorID,
				CategoryID:   fakedCategoryID,
				Disqualified: true,
			}, true)

		f.store.
			On("DeleteTick", fakedTickID).
			Return()

		effects := f.engine.HandleAscentDeregistered(domain.AscentDeregisteredEvent{
			TickID:       fakedTickID,
			CompetitorID: fakedCompetitorID,
			ProblemID:    fakedProblemID,
		})

		require.Nil(t, effects)

		awaitExpectations(t)
	})

	t.Run("AscentDeregistered", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedTickID := testutils.RandomResourceID[domain.TickID]()
		fakedCompetitorID := testutils.RandomResourceID[domain.CompetitorID]()
		fakedCategoryID := testutils.RandomResourceID[domain.CategoryID]()
		fakedProblemID := testutils.RandomResourceID[domain.ProblemID]()

		f.store.
			On("GetCompetitor", fakedCompetitorID).
			Return(scores.Competitor{
				ID:         fakedCompetitorID,
				CategoryID: fakedCategoryID,
			}, true)

		f.store.
			On("DeleteTick", fakedTickID).
			Return()

		effects := slices.Collect(f.engine.HandleAscentDeregistered(domain.AscentDeregisteredEvent{
			TickID:       fakedTickID,
			CompetitorID: fakedCompetitorID,
			ProblemID:    fakedProblemID,
		}))

		require.ElementsMatch(t, effects, []scores.Effect{
			scores.EffectCalculatePointValues{CategoryID: fakedCategoryID, ProblemID: fakedProblemID},
			scores.EffectScoreCompetitor{CompetitorID: fakedCompetitorID},
		})

		awaitExpectations(t)
	})

	t.Run("GetDirtyResults", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		now := time.Now()

		fakedCompetitorID1 := testutils.RandomResourceID[domain.CompetitorID]()
		fakedCompetitorID2 := testutils.RandomResourceID[domain.CompetitorID]()
		fakedCompetitorID3 := testutils.RandomResourceID[domain.CompetitorID]()

		fakedResults := []domain.Result{
			{
				CompetitorID: fakedCompetitorID1,
				Timestamp:    now,
				Score:        "100p",
				Placement:    1,
				RankOrder:    0,
				Finalist:     true,
			},
			{
				CompetitorID: fakedCompetitorID2,
				Timestamp:    now,
				Score:        "200p",
				Placement:    2,
				RankOrder:    1,
				Finalist:     true,
			},
			{
				CompetitorID: fakedCompetitorID3,
				Timestamp:    now,
				Score:        "300p",
				Placement:    3,
				RankOrder:    2,
				Finalist:     false,
			},
		}

		f.store.
			On("GetDirtyResults").
			Return(fakedResults)

		results := f.engine.GetDirtyResults()

		assert.Equal(t, fakedResults, results)

		awaitExpectations(t)
	})

	t.Run("CalculatePointValues_UsePointsFalse", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedCategoryID := testutils.RandomResourceID[domain.CategoryID]()
		fakedProblemID := testutils.RandomResourceID[domain.ProblemID]()

		f.store.
			On("GetRules").
			Return(scores.Rules{
				UsePoints: false,
			})

		effects := f.engine.CalculatePointValues(fakedCategoryID, fakedProblemID)

		assert.Nil(t, effects)

		awaitExpectations(t)
	})

	t.Run("CalculatePointValues_ProblemNotFound", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedCategoryID := testutils.RandomResourceID[domain.CategoryID]()
		fakedProblemID := testutils.RandomResourceID[domain.ProblemID]()

		f.store.
			On("GetRules").
			Return(scores.Rules{
				UsePoints: true,
			})

		f.store.
			On("GetProblem", fakedProblemID).
			Return(scores.Problem{}, false)

		effects := f.engine.CalculatePointValues(fakedCategoryID, fakedProblemID)

		assert.Nil(t, effects)

		awaitExpectations(t)
	})

	t.Run("CalculatePointValues_NonPooledPoints", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedCategoryID := testutils.RandomResourceID[domain.CategoryID]()
		fakedProblemID := testutils.RandomResourceID[domain.ProblemID]()

		fakedCompetitor1ID := testutils.RandomResourceID[domain.CompetitorID]()
		fakedCompetitor2ID := testutils.RandomResourceID[domain.CompetitorID]()

		f.store.
			On("GetRules").
			Return(scores.Rules{
				UsePoints:    true,
				PooledPoints: false,
			})

		f.store.
			On("GetProblem", fakedProblemID).
			Return(scores.Problem{
				ID: fakedProblemID,
				ProblemValue: scores.ProblemValue{
					PointsZone1: 1000,
					PointsZone2: 2000,
					PointsTop:   5000,
					FlashBonus:  500,
				},
			}, true)

		f.store.
			On("GetCompetitorsByCategory", fakedCategoryID).
			Return(slices.Values([]scores.Competitor{
				{
					ID: fakedCompetitor1ID,
				},
				{
					ID: fakedCompetitor2ID,
				},
			}))

		f.store.
			On("GetTick", fakedCompetitor1ID, fakedProblemID).
			Return(scores.Tick{
				CompetitorID:  fakedCompetitor1ID,
				ProblemID:     fakedProblemID,
				Top:           true,
				AttemptsTop:   1,
				Zone2:         true,
				AttemptsZone2: 1,
				Zone1:         true,
				AttemptsZone1: 1,
			}, true).
			On("GetTick", fakedCompetitor2ID, fakedProblemID).
			Return(scores.Tick{}, false)

		f.store.
			On("GetPointValue", fakedCompetitor1ID, fakedProblemID).
			Return(domain.PointValue{}, false).
			On("GetPointValue", fakedCompetitor2ID, fakedProblemID).
			Return(domain.PointValue{}, false)

		f.store.
			On("SavePointValue", fakedCompetitor1ID, fakedProblemID, domain.PointValue{
				CompetitorID: fakedCompetitor1ID,
				ProblemID:    fakedProblemID,
				Current:      5500,
				Zone1:        1000,
				Zone2:        2000,
				Top:          5000,
				FlashBonus:   500,
			}).Return().
			On("SavePointValue", fakedCompetitor2ID, fakedProblemID, domain.PointValue{
				CompetitorID: fakedCompetitor2ID,
				ProblemID:    fakedProblemID,
				Current:      0,
				Zone1:        1000,
				Zone2:        2000,
				Top:          5000,
				FlashBonus:   500,
			}).Return()

		effects := slices.Collect(f.engine.CalculatePointValues(fakedCategoryID, fakedProblemID))

		require.ElementsMatch(t, effects, []scores.Effect{
			scores.EffectScoreCompetitor{CompetitorID: fakedCompetitor1ID},
			scores.EffectScoreCompetitor{CompetitorID: fakedCompetitor2ID},
		})

		awaitExpectations(t)
	})

	t.Run("CalculatePointValues_ValueUnchanged", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedCategoryID := testutils.RandomResourceID[domain.CategoryID]()
		fakedProblemID := testutils.RandomResourceID[domain.ProblemID]()

		fakedCompetitorID := testutils.RandomResourceID[domain.CompetitorID]()

		fakedCompetitor := scores.Competitor{
			ID: fakedCompetitorID,
		}

		fakedTick := scores.Tick{
			CompetitorID:  fakedCompetitorID,
			ProblemID:     fakedProblemID,
			Zone1:         true,
			AttemptsZone1: 1,
			Zone2:         true,
			AttemptsZone2: 1,
			Top:           true,
			AttemptsTop:   1,
		}

		f.store.
			On("GetRules").
			Return(scores.Rules{
				UsePoints: true,
			})

		f.store.
			On("GetProblem", fakedProblemID).
			Return(scores.Problem{
				ID: fakedProblemID,
				ProblemValue: scores.ProblemValue{
					PointsZone1: 1000,
					PointsZone2: 2000,
					PointsTop:   5000,
					FlashBonus:  1000,
				},
			}, true)

		f.store.
			On("GetCompetitorsByCategory", fakedCategoryID).
			Return(slices.Values([]scores.Competitor{fakedCompetitor}))

		f.store.
			On("GetTick", fakedCompetitorID, fakedProblemID).
			Return(fakedTick, true)

		f.store.
			On("GetPointValue", fakedCompetitorID, fakedProblemID).
			Return(domain.PointValue{
				CompetitorID: fakedCompetitorID,
				ProblemID:    fakedProblemID,
				Current:      6000,
				Zone1:        1000,
				Zone2:        2000,
				Top:          5000,
				FlashBonus:   1000,
			}, true)

		f.store.
			On("SavePointValue", fakedCompetitorID, fakedProblemID, domain.PointValue{
				CompetitorID: fakedCompetitorID,
				ProblemID:    fakedProblemID,
				Current:      6000,
				Zone1:        1000,
				Zone2:        2000,
				Top:          5000,
				FlashBonus:   1000,
			}).Return()

		effects := f.engine.CalculatePointValues(fakedCategoryID, fakedProblemID)

		require.Nil(t, effects)

		awaitExpectations(t)
	})

	t.Run("CalculatePointValues_PooledPoints_UpdatedValue", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedCategoryID := testutils.RandomResourceID[domain.CategoryID]()
		fakedProblemID := testutils.RandomResourceID[domain.ProblemID]()

		fakedCompetitor1ID := testutils.RandomResourceID[domain.CompetitorID]()
		fakedCompetitor2ID := testutils.RandomResourceID[domain.CompetitorID]()
		fakedCompetitor3ID := testutils.RandomResourceID[domain.CompetitorID]()
		fakedCompetitor4ID := testutils.RandomResourceID[domain.CompetitorID]()
		fakedCompetitor5ID := testutils.RandomResourceID[domain.CompetitorID]()
		fakedCompetitor6ID := testutils.RandomResourceID[domain.CompetitorID]()

		fakedTicksByCompetitor := map[domain.CompetitorID]scores.Tick{
			fakedCompetitor1ID: {
				CompetitorID:  fakedCompetitor1ID,
				ProblemID:     fakedProblemID,
				Zone1:         true,
				AttemptsZone1: 1,
				Zone2:         true,
				AttemptsZone2: 1,
				Top:           true,
				AttemptsTop:   1,
			},
			fakedCompetitor2ID: {
				CompetitorID:  fakedCompetitor2ID,
				ProblemID:     fakedProblemID,
				Zone1:         true,
				AttemptsZone1: 1,
				Zone2:         true,
				AttemptsZone2: 1,
				Top:           true,
				AttemptsTop:   1,
			},
			fakedCompetitor3ID: {
				CompetitorID:  fakedCompetitor3ID,
				ProblemID:     fakedProblemID,
				Zone1:         true,
				AttemptsZone1: 999,
				Zone2:         true,
				AttemptsZone2: 999,
				Top:           true,
				AttemptsTop:   999,
			},
			fakedCompetitor4ID: {
				CompetitorID:  fakedCompetitor4ID,
				ProblemID:     fakedProblemID,
				Zone1:         true,
				AttemptsZone1: 999,
				Zone2:         true,
				AttemptsZone2: 999,
				Top:           false,
				AttemptsTop:   999,
			},
			fakedCompetitor5ID: {
				CompetitorID:  fakedCompetitor5ID,
				ProblemID:     fakedProblemID,
				Zone1:         true,
				AttemptsZone1: 999,
				Zone2:         false,
				AttemptsZone2: 999,
				Top:           false,
				AttemptsTop:   999,
			},
		}

		f.store.
			On("GetRules").
			Return(scores.Rules{
				UsePoints:    true,
				PooledPoints: true,
			})

		f.store.
			On("GetProblem", fakedProblemID).
			Return(scores.Problem{
				ID: fakedProblemID,
				ProblemValue: scores.ProblemValue{
					PointsZone1: 500,
					PointsZone2: 1000,
					PointsTop:   5000,
					FlashBonus:  1000,
				},
			}, true)

		f.store.
			On("GetTicksByProblem", fakedCategoryID, fakedProblemID).
			Return(maps.Values(fakedTicksByCompetitor))

		f.store.
			On("GetCompetitor", fakedCompetitor1ID).
			Return(scores.Competitor{ID: fakedCompetitor1ID}, true).
			On("GetCompetitor", fakedCompetitor2ID).
			Return(scores.Competitor{ID: fakedCompetitor2ID}, true).
			On("GetCompetitor", fakedCompetitor3ID).
			Return(scores.Competitor{ID: fakedCompetitor3ID}, true).
			On("GetCompetitor", fakedCompetitor4ID).
			Return(scores.Competitor{ID: fakedCompetitor4ID}, true).
			On("GetCompetitor", fakedCompetitor5ID).
			Return(scores.Competitor{ID: fakedCompetitor5ID}, true)

		f.store.
			On("GetCompetitorsByCategory", fakedCategoryID).
			Return(slices.Values([]scores.Competitor{
				{ID: fakedCompetitor1ID},
				{ID: fakedCompetitor2ID},
				{ID: fakedCompetitor3ID},
				{ID: fakedCompetitor4ID},
				{ID: fakedCompetitor5ID},
				{ID: fakedCompetitor6ID},
			}))

		f.store.
			On("GetTick", fakedCompetitor1ID, fakedProblemID).
			Return(fakedTicksByCompetitor[fakedCompetitor1ID], true).
			On("GetTick", fakedCompetitor2ID, fakedProblemID).
			Return(fakedTicksByCompetitor[fakedCompetitor2ID], true).
			On("GetTick", fakedCompetitor3ID, fakedProblemID).
			Return(fakedTicksByCompetitor[fakedCompetitor3ID], true).
			On("GetTick", fakedCompetitor4ID, fakedProblemID).
			Return(fakedTicksByCompetitor[fakedCompetitor4ID], true).
			On("GetTick", fakedCompetitor5ID, fakedProblemID).
			Return(fakedTicksByCompetitor[fakedCompetitor5ID], true).
			On("GetTick", fakedCompetitor6ID, fakedProblemID).
			Return(scores.Tick{}, false)

		f.store.
			On("GetPointValue", fakedCompetitor1ID, fakedProblemID).
			Return(domain.PointValue{}, false).
			On("GetPointValue", fakedCompetitor2ID, fakedProblemID).
			Return(domain.PointValue{}, false).
			On("GetPointValue", fakedCompetitor3ID, fakedProblemID).
			Return(domain.PointValue{}, false).
			On("GetPointValue", fakedCompetitor4ID, fakedProblemID).
			Return(domain.PointValue{}, false).
			On("GetPointValue", fakedCompetitor5ID, fakedProblemID).
			Return(domain.PointValue{}, false).
			On("GetPointValue", fakedCompetitor6ID, fakedProblemID).
			Return(domain.PointValue{}, false)

		f.store.
			On("SavePointValue", fakedCompetitor1ID, fakedProblemID, domain.PointValue{
				CompetitorID: fakedCompetitor1ID,
				ProblemID:    fakedProblemID,
				Current:      2166,
				Zone1:        100,
				Zone2:        250,
				Top:          1666,
				FlashBonus:   500,
			}).Return().
			On("SavePointValue", fakedCompetitor2ID, fakedProblemID, domain.PointValue{
				CompetitorID: fakedCompetitor2ID,
				ProblemID:    fakedProblemID,
				Current:      2166,
				Zone1:        100,
				Zone2:        250,
				Top:          1666,
				FlashBonus:   500,
			}).Return().
			On("SavePointValue", fakedCompetitor3ID, fakedProblemID, domain.PointValue{
				CompetitorID: fakedCompetitor3ID,
				ProblemID:    fakedProblemID,
				Current:      1666,
				Zone1:        100,
				Zone2:        250,
				Top:          1666,
				FlashBonus:   0,
			}).Return().
			On("SavePointValue", fakedCompetitor4ID, fakedProblemID, domain.PointValue{
				CompetitorID: fakedCompetitor4ID,
				ProblemID:    fakedProblemID,
				Current:      250,
				Zone1:        100,
				Zone2:        250,
				Top:          1250,
				FlashBonus:   0,
			}).Return().
			On("SavePointValue", fakedCompetitor5ID, fakedProblemID, domain.PointValue{
				CompetitorID: fakedCompetitor5ID,
				ProblemID:    fakedProblemID,
				Current:      100,
				Zone1:        100,
				Zone2:        200,
				Top:          1250,
				FlashBonus:   0,
			}).Return().
			On("SavePointValue", fakedCompetitor6ID, fakedProblemID, domain.PointValue{
				CompetitorID: fakedCompetitor6ID,
				ProblemID:    fakedProblemID,
				Current:      0,
				Zone1:        83,
				Zone2:        200,
				Top:          1250,
				FlashBonus:   333,
			}).Return()

		effects := slices.Collect(f.engine.CalculatePointValues(fakedCategoryID, fakedProblemID))

		require.ElementsMatch(t, effects, []scores.Effect{
			scores.EffectScoreCompetitor{CompetitorID: fakedCompetitor1ID},
			scores.EffectScoreCompetitor{CompetitorID: fakedCompetitor2ID},
			scores.EffectScoreCompetitor{CompetitorID: fakedCompetitor3ID},
			scores.EffectScoreCompetitor{CompetitorID: fakedCompetitor4ID},
			scores.EffectScoreCompetitor{CompetitorID: fakedCompetitor5ID},
			scores.EffectScoreCompetitor{CompetitorID: fakedCompetitor6ID},
		})

		awaitExpectations(t)
	})

	t.Run("CalculatePointValues_PooledPoints_ExcludeDisqualified", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedCategoryID := testutils.RandomResourceID[domain.CategoryID]()
		fakedProblemID := testutils.RandomResourceID[domain.ProblemID]()

		fakedCompetitor1 := scores.Competitor{
			ID: testutils.RandomResourceID[domain.CompetitorID](),
		}

		fakedCompetitor2 := scores.Competitor{
			ID: testutils.RandomResourceID[domain.CompetitorID](),
		}

		fakedCompetitor3 := scores.Competitor{
			ID:           testutils.RandomResourceID[domain.CompetitorID](),
			Disqualified: true,
		}

		fakedTicksByCompetitor := map[domain.CompetitorID]scores.Tick{
			fakedCompetitor1.ID: {
				CompetitorID:  fakedCompetitor1.ID,
				ProblemID:     fakedProblemID,
				Zone1:         true,
				AttemptsZone1: 1,
				Zone2:         true,
				AttemptsZone2: 1,
				Top:           true,
				AttemptsTop:   1,
			},
			fakedCompetitor2.ID: {
				CompetitorID:  fakedCompetitor2.ID,
				ProblemID:     fakedProblemID,
				Zone1:         true,
				AttemptsZone1: 1,
				Zone2:         true,
				AttemptsZone2: 1,
				Top:           true,
				AttemptsTop:   1,
			},
			fakedCompetitor3.ID: {
				CompetitorID:  fakedCompetitor3.ID,
				ProblemID:     fakedProblemID,
				Zone1:         true,
				AttemptsZone1: 1,
				Zone2:         true,
				AttemptsZone2: 1,
				Top:           true,
				AttemptsTop:   1,
			},
		}

		f.store.
			On("GetRules").
			Return(scores.Rules{
				UsePoints:    true,
				PooledPoints: true,
			})

		f.store.
			On("GetProblem", fakedProblemID).
			Return(scores.Problem{
				ID: fakedProblemID,
				ProblemValue: scores.ProblemValue{
					PointsZone1: 500,
					PointsZone2: 1000,
					PointsTop:   5000,
					FlashBonus:  250,
				},
			}, true)

		f.store.
			On("GetTicksByProblem", fakedCategoryID, fakedProblemID).
			Return(maps.Values(fakedTicksByCompetitor))

		f.store.
			On("GetCompetitor", fakedCompetitor1.ID).
			Return(fakedCompetitor1, true).
			On("GetCompetitor", fakedCompetitor2.ID).
			Return(fakedCompetitor2, true).
			On("GetCompetitor", fakedCompetitor3.ID).
			Return(fakedCompetitor3, true)

		f.store.
			On("GetCompetitorsByCategory", fakedCategoryID).
			Return(slices.Values([]scores.Competitor{
				fakedCompetitor1,
				fakedCompetitor2,
				fakedCompetitor3,
			}))

		f.store.
			On("GetTick", fakedCompetitor1.ID, fakedProblemID).
			Return(fakedTicksByCompetitor[fakedCompetitor1.ID], true).
			On("GetTick", fakedCompetitor2.ID, fakedProblemID).
			Return(fakedTicksByCompetitor[fakedCompetitor2.ID], true).
			On("GetTick", fakedCompetitor3.ID, fakedProblemID).
			Return(fakedTicksByCompetitor[fakedCompetitor3.ID], true)

		f.store.
			On("GetPointValue", fakedCompetitor1.ID, fakedProblemID).
			Return(domain.PointValue{}, false).
			On("GetPointValue", fakedCompetitor2.ID, fakedProblemID).
			Return(domain.PointValue{}, false).
			On("GetPointValue", fakedCompetitor3.ID, fakedProblemID).
			Return(domain.PointValue{}, false)

		f.store.
			On("SavePointValue", fakedCompetitor1.ID, fakedProblemID, domain.PointValue{
				CompetitorID: fakedCompetitor1.ID,
				ProblemID:    fakedProblemID,
				Current:      2625,
				Zone1:        250,
				Zone2:        500,
				Top:          2500,
				FlashBonus:   125,
			}).Return().
			On("SavePointValue", fakedCompetitor2.ID, fakedProblemID, domain.PointValue{
				CompetitorID: fakedCompetitor2.ID,
				ProblemID:    fakedProblemID,
				Current:      2625,
				Zone1:        250,
				Zone2:        500,
				Top:          2500,
				FlashBonus:   125,
			}).Return().
			On("SavePointValue", fakedCompetitor3.ID, fakedProblemID, domain.PointValue{
				CompetitorID: fakedCompetitor3.ID,
				ProblemID:    fakedProblemID,
				Current:      0,
				Zone1:        0,
				Zone2:        0,
				Top:          0,
				FlashBonus:   0,
			}).Return()

		effects := slices.Collect(f.engine.CalculatePointValues(fakedCategoryID, fakedProblemID))

		require.ElementsMatch(t, effects, []scores.Effect{
			scores.EffectScoreCompetitor{CompetitorID: fakedCompetitor1.ID},
			scores.EffectScoreCompetitor{CompetitorID: fakedCompetitor2.ID},
			scores.EffectScoreCompetitor{CompetitorID: fakedCompetitor3.ID},
		})

		awaitExpectations(t)
	})

	t.Run("ScoreCompetitor_CompetitorNotFound", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedCompetitorID := testutils.RandomResourceID[domain.CompetitorID]()

		f.store.
			On("GetCompetitor", fakedCompetitorID).
			Return(scores.Competitor{}, false)

		effects := f.engine.ScoreCompetitor(fakedCompetitorID)

		assert.Nil(t, effects)

		awaitExpectations(t)
	})

	t.Run("ScoreCompetitor_Disqualified", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedCompetitorID := testutils.RandomResourceID[domain.CompetitorID]()
		fakedCategoryID := testutils.RandomResourceID[domain.CategoryID]()

		f.store.
			On("GetRules").
			Return(scores.Rules{})

		f.store.
			On("GetCompetitor", fakedCompetitorID).
			Return(scores.Competitor{
				ID:           fakedCompetitorID,
				CategoryID:   fakedCategoryID,
				Disqualified: true,
				Score: scores.Score{
					Points:       100,
					Tops:         1,
					AttemptsTops: 10,
				},
			}, true)

		f.store.
			On("SaveCompetitor", scores.Competitor{
				ID:           fakedCompetitorID,
				CategoryID:   fakedCategoryID,
				Disqualified: true,
				Score:        scores.Score{},
			}).Return()

		effects := slices.Collect(f.engine.ScoreCompetitor(fakedCompetitorID))

		require.ElementsMatch(t, effects, []scores.Effect{
			scores.EffectRankCategory{CategoryID: fakedCategoryID},
		})

		awaitExpectations(t)
	})

	t.Run("ScoreCompetitor_PointsAndAttempts", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedCompetitorID := testutils.RandomResourceID[domain.CompetitorID]()
		fakedCategoryID := testutils.RandomResourceID[domain.CategoryID]()
		fakedProblem1ID := testutils.RandomResourceID[domain.ProblemID]()
		fakedProblem2ID := testutils.RandomResourceID[domain.ProblemID]()
		fakedProblem3ID := testutils.RandomResourceID[domain.ProblemID]()
		fakedProblem4ID := testutils.RandomResourceID[domain.ProblemID]()

		f.store.
			On("GetCompetitor", fakedCompetitorID).
			Return(scores.Competitor{
				ID:         fakedCompetitorID,
				CategoryID: fakedCategoryID,
				Score:      scores.Score{},
			}, true)

		f.store.
			On("GetTicksByCompetitor", fakedCompetitorID).
			Return(slices.Values([]scores.Tick{
				{
					CompetitorID:  fakedCompetitorID,
					ProblemID:     fakedProblem1ID,
					Top:           false,
					AttemptsTop:   1_000_000,
					Zone2:         false,
					AttemptsZone2: 1_000_000,
					Zone1:         false,
					AttemptsZone1: 1_000_000,
				},
				{
					CompetitorID:  fakedCompetitorID,
					ProblemID:     fakedProblem2ID,
					Top:           true,
					AttemptsTop:   300,
					Zone2:         true,
					AttemptsZone2: 200,
					Zone1:         true,
					AttemptsZone1: 100,
				},
				{
					CompetitorID:  fakedCompetitorID,
					ProblemID:     fakedProblem3ID,
					Top:           true,
					AttemptsTop:   30,
					Zone2:         true,
					AttemptsZone2: 20,
					Zone1:         true,
					AttemptsZone1: 10,
				},
				{
					CompetitorID:  fakedCompetitorID,
					ProblemID:     fakedProblem4ID,
					Top:           true,
					AttemptsTop:   3,
					Zone2:         true,
					AttemptsZone2: 2,
					Zone1:         true,
					AttemptsZone1: 1,
				},
			}))

		f.store.
			On("GetPointValue", fakedCompetitorID, fakedProblem1ID).
			Return(domain.PointValue{
				CompetitorID: fakedCompetitorID,
				ProblemID:    fakedProblem1ID,
				Current:      0,
			}, true).
			On("GetPointValue", fakedCompetitorID, fakedProblem2ID).
			Return(domain.PointValue{
				CompetitorID: fakedCompetitorID,
				ProblemID:    fakedProblem2ID,
				Current:      100,
			}, true).
			On("GetPointValue", fakedCompetitorID, fakedProblem3ID).
			Return(domain.PointValue{
				CompetitorID: fakedCompetitorID,
				ProblemID:    fakedProblem3ID,
				Current:      110,
			}, true).
			On("GetPointValue", fakedCompetitorID, fakedProblem4ID).
			Return(domain.PointValue{
				CompetitorID: fakedCompetitorID,
				ProblemID:    fakedProblem4ID,
				Current:      0,
			}, true)

		f.store.
			On("GetProblem", fakedProblem1ID).
			Return(scores.Problem{
				Zone2Enabled: true,
				Zone1Enabled: true,
			}, true).
			On("GetProblem", fakedProblem2ID).
			Return(scores.Problem{
				Zone2Enabled: true,
				Zone1Enabled: true,
			}, true).
			On("GetProblem", fakedProblem3ID).
			Return(scores.Problem{
				Zone1Enabled: true,
			}, true).
			On("GetProblem", fakedProblem4ID).
			Return(scores.Problem{}, true)

		f.store.
			On("GetRules").
			Return(scores.Rules{
				UsePoints: true,
			})

		f.store.
			On("SaveCompetitor", scores.Competitor{
				ID:         fakedCompetitorID,
				CategoryID: fakedCategoryID,
				Score: scores.Score{
					Points:         210,
					Tops:           3,
					AttemptsTops:   300 + 30 + 3,
					Zone1s:         2,
					AttemptsZone1s: 100 + 10,
					Zone2s:         1,
					AttemptsZone2s: 200,
				},
			}).Return()

		effects := slices.Collect(f.engine.ScoreCompetitor(fakedCompetitorID))

		require.ElementsMatch(t, effects, []scores.Effect{
			scores.EffectRankCategory{CategoryID: fakedCategoryID},
		})

		awaitExpectations(t)
	})

	t.Run("ScoreCompetitor_NoPoints", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedCompetitorID := testutils.RandomResourceID[domain.CompetitorID]()
		fakedCategoryID := testutils.RandomResourceID[domain.CategoryID]()
		fakedProblem1ID := testutils.RandomResourceID[domain.ProblemID]()
		fakedProblem2ID := testutils.RandomResourceID[domain.ProblemID]()

		f.store.
			On("GetCompetitor", fakedCompetitorID).
			Return(scores.Competitor{
				ID:         fakedCompetitorID,
				CategoryID: fakedCategoryID,
				Score:      scores.Score{},
			}, true)

		f.store.
			On("GetTicksByCompetitor", fakedCompetitorID).
			Return(slices.Values([]scores.Tick{
				{
					CompetitorID: fakedCompetitorID,
					ProblemID:    fakedProblem1ID,
					Top:          true,
					AttemptsTop:  10,
				},
				{
					CompetitorID: fakedCompetitorID,
					ProblemID:    fakedProblem2ID,
					Top:          true,
					AttemptsTop:  20,
				},
			}))

		f.store.
			On("GetProblem", fakedProblem1ID).
			Return(scores.Problem{}, true).
			On("GetProblem", fakedProblem2ID).
			Return(scores.Problem{}, true)

		f.store.
			On("GetRules").
			Return(scores.Rules{})

		f.store.
			On("SaveCompetitor", scores.Competitor{
				ID:         fakedCompetitorID,
				CategoryID: fakedCategoryID,
				Score: scores.Score{
					Points:       0,
					Tops:         2,
					AttemptsTops: 30,
				},
			}).Return()

		effects := slices.Collect(f.engine.ScoreCompetitor(fakedCompetitorID))

		require.ElementsMatch(t, effects, []scores.Effect{
			scores.EffectRankCategory{CategoryID: fakedCategoryID},
		})

		awaitExpectations(t)
	})

	t.Run("ScoreCompetitor_WithProblemLimit", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedCompetitorID := testutils.RandomResourceID[domain.CompetitorID]()
		fakedCategoryID := testutils.RandomResourceID[domain.CategoryID]()
		fakedProblem1ID := testutils.RandomResourceID[domain.ProblemID]()
		fakedProblem2ID := testutils.RandomResourceID[domain.ProblemID]()
		fakedProblem3ID := testutils.RandomResourceID[domain.ProblemID]()

		f.store.
			On("GetCompetitor", fakedCompetitorID).
			Return(scores.Competitor{
				ID:         fakedCompetitorID,
				CategoryID: fakedCategoryID,
				Score:      scores.Score{},
			}, true)

		f.store.
			On("GetTicksByCompetitor", fakedCompetitorID).
			Return(slices.Values([]scores.Tick{
				{
					CompetitorID: fakedCompetitorID,
					ProblemID:    fakedProblem1ID,
					Top:          true,
					AttemptsTop:  1,
				},
				{
					CompetitorID: fakedCompetitorID,
					ProblemID:    fakedProblem2ID,
					Top:          true,
					AttemptsTop:  1,
				},
				{
					CompetitorID: fakedCompetitorID,
					ProblemID:    fakedProblem3ID,
					Top:          true,
					AttemptsTop:  1,
				},
			}))

		f.store.
			On("GetPointValue", fakedCompetitorID, fakedProblem1ID).
			Return(domain.PointValue{
				CompetitorID: fakedCompetitorID,
				ProblemID:    fakedProblem1ID,
				Current:      55,
			}, true).
			On("GetPointValue", fakedCompetitorID, fakedProblem2ID).
			Return(domain.PointValue{
				CompetitorID: fakedCompetitorID,
				ProblemID:    fakedProblem2ID,
				Current:      110,
			}, true)

		f.store.
			On("GetProblem", fakedProblem1ID).
			Return(scores.Problem{}, true).
			On("GetProblem", fakedProblem2ID).
			Return(scores.Problem{}, true).
			On("GetProblem", fakedProblem3ID).
			Return(scores.Problem{}, true)

		f.store.
			On("GetPointValue", fakedCompetitorID, fakedProblem3ID).
			Return(domain.PointValue{
				CompetitorID: fakedCompetitorID,
				ProblemID:    fakedProblem3ID,
				Current:      82,
			}, true)

		f.store.
			On("GetRules").
			Return(scores.Rules{
				UsePoints:          true,
				QualifyingProblems: 2,
			})

		f.store.
			On("SaveCompetitor", scores.Competitor{
				ID:         fakedCompetitorID,
				CategoryID: fakedCategoryID,
				Score: scores.Score{
					Points:       192,
					Tops:         3,
					AttemptsTops: 3,
				},
			}).Return()

		effects := slices.Collect(f.engine.ScoreCompetitor(fakedCompetitorID))

		require.ElementsMatch(t, effects, []scores.Effect{
			scores.EffectRankCategory{CategoryID: fakedCategoryID},
		})

		awaitExpectations(t)
	})

	t.Run("ScoreCompetitor_NoScoreChange", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedCompetitorID := testutils.RandomResourceID[domain.CompetitorID]()
		fakedCategoryID := testutils.RandomResourceID[domain.CategoryID]()
		fakedProblem1ID := testutils.RandomResourceID[domain.ProblemID]()

		f.store.
			On("GetCompetitor", fakedCompetitorID).
			Return(scores.Competitor{
				ID:         fakedCompetitorID,
				CategoryID: fakedCategoryID,
				Score: scores.Score{
					Points:       110,
					Tops:         1,
					AttemptsTops: 1,
				},
			}, true)

		f.store.
			On("GetTicksByCompetitor", fakedCompetitorID).
			Return(slices.Values([]scores.Tick{
				{
					CompetitorID: fakedCompetitorID,
					ProblemID:    fakedProblem1ID,
					Top:          true,
					AttemptsTop:  1,
				},
			}))

		f.store.
			On("GetPointValue", fakedCompetitorID, fakedProblem1ID).
			Return(domain.PointValue{
				CompetitorID: fakedCompetitorID,
				ProblemID:    fakedProblem1ID,
				Current:      110,
			}, true)

		f.store.
			On("GetProblem", fakedProblem1ID).
			Return(scores.Problem{}, true)

		f.store.
			On("GetRules").
			Return(scores.Rules{
				UsePoints:          true,
				QualifyingProblems: 0,
			})

		effects := f.engine.ScoreCompetitor(fakedCompetitorID)

		assert.Nil(t, effects)

		awaitExpectations(t)
	})

	t.Run("ScoreCompetitor_ProblemValueNotFound", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedCompetitorID := testutils.RandomResourceID[domain.CompetitorID]()
		fakedCategoryID := testutils.RandomResourceID[domain.CategoryID]()
		fakedProblem1ID := testutils.RandomResourceID[domain.ProblemID]()
		fakedProblem2ID := testutils.RandomResourceID[domain.ProblemID]()

		f.store.
			On("GetCompetitor", fakedCompetitorID).
			Return(scores.Competitor{
				ID:         fakedCompetitorID,
				CategoryID: fakedCategoryID,
				Score:      scores.Score{},
			}, true)

		f.store.
			On("GetTicksByCompetitor", fakedCompetitorID).
			Return(slices.Values([]scores.Tick{
				{
					CompetitorID: fakedCompetitorID,
					ProblemID:    fakedProblem1ID,
					Top:          true,
					AttemptsTop:  1,
				},
				{
					CompetitorID: fakedCompetitorID,
					ProblemID:    fakedProblem2ID,
					Top:          true,
					AttemptsTop:  1,
				},
			}))

		f.store.
			On("GetPointValue", fakedCompetitorID, fakedProblem1ID).
			Return(domain.PointValue{
				CompetitorID: fakedCompetitorID,
				ProblemID:    fakedProblem1ID,
				Current:      110,
			}, true)

		f.store.
			On("GetPointValue", fakedCompetitorID, fakedProblem2ID).
			Return(domain.PointValue{}, false)

		f.store.
			On("GetProblem", fakedProblem1ID).
			Return(scores.Problem{}, true).
			On("GetProblem", fakedProblem2ID).
			Return(scores.Problem{}, true)

		f.store.
			On("GetRules").
			Return(scores.Rules{
				UsePoints:          true,
				QualifyingProblems: 0,
			})

		f.store.
			On("SaveCompetitor", scores.Competitor{
				ID:         fakedCompetitorID,
				CategoryID: fakedCategoryID,
				Score: scores.Score{
					Points:       110,
					Tops:         2,
					AttemptsTops: 2,
				},
			}).Return()

		effects := slices.Collect(f.engine.ScoreCompetitor(fakedCompetitorID))

		require.ElementsMatch(t, effects, []scores.Effect{
			scores.EffectRankCategory{CategoryID: fakedCategoryID},
		})

		awaitExpectations(t)
	})

	t.Run("RankCategory", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			f, awaitExpectations := makeFixture()

			fakedCategoryID := testutils.RandomResourceID[domain.CategoryID]()
			fakedCompetitor1ID := testutils.RandomResourceID[domain.CompetitorID]()
			fakedCompetitor2ID := testutils.RandomResourceID[domain.CompetitorID]()
			fakedCompetitor3ID := testutils.RandomResourceID[domain.CompetitorID]()

			f.store.
				On("GetRules").
				Return(scores.Rules{
					Finalists: 2,
					UsePoints: true,
				})

			f.store.
				On("GetCompetitorsByCategory", fakedCategoryID).
				Return(slices.Values([]scores.Competitor{
					{
						ID:         fakedCompetitor1ID,
						CategoryID: fakedCategoryID,
						Score: scores.Score{
							Points: 1000,
						},
					},
					{
						ID:         fakedCompetitor2ID,
						CategoryID: fakedCategoryID,
						Score: scores.Score{
							Points: 2000,
						},
					},
					{
						ID:         fakedCompetitor3ID,
						CategoryID: fakedCategoryID,
						Score: scores.Score{
							Points: 1500,
						},
					},
				}))

			f.store.
				On("SaveResult", domain.Result{
					CompetitorID: fakedCompetitor2ID,
					Score:        "200p",
					Placement:    1,
					RankOrder:    0,
					Finalist:     true,
					Timestamp:    time.Now(),
				}).Return()

			f.store.
				On("SaveResult", domain.Result{
					CompetitorID: fakedCompetitor3ID,
					Score:        "150p",
					Placement:    2,
					RankOrder:    1,
					Finalist:     true,
					Timestamp:    time.Now(),
				}).Return()

			f.store.
				On("SaveResult", domain.Result{
					CompetitorID: fakedCompetitor1ID,
					Score:        "100p",
					Placement:    3,
					RankOrder:    2,
					Finalist:     false,
					Timestamp:    time.Now(),
				}).Return()

			f.engine.RankCategory(fakedCategoryID)

			awaitExpectations(t)
		})
	})

	t.Run("RankCategory_EnabledZones", func(t *testing.T) {
		for _, tt := range []struct {
			name     string
			problems []scores.Problem
			expected string
		}{
			{name: "NoZones", problems: []scores.Problem{{}}, expected: "2t"},
			{
				name:     "Zone1",
				problems: []scores.Problem{{Zone1Enabled: true}, {}},
				expected: "2t 4z₁",
			},
			{
				name:     "Zone2",
				problems: []scores.Problem{{Zone1Enabled: true}, {Zone1Enabled: true, Zone2Enabled: true}, {}},
				expected: "2t 3z₂ 4z₁",
			},
		} {
			t.Run(tt.name, func(t *testing.T) {
				synctest.Test(t, func(t *testing.T) {
					f, awaitExpectations := makeFixture()

					fakedCategoryID := testutils.RandomResourceID[domain.CategoryID]()
					fakedCompetitorID := testutils.RandomResourceID[domain.CompetitorID]()

					f.store.On("GetRules").Return(scores.Rules{})
					f.store.On("GetAllProblems").Return(slices.Values(tt.problems))
					f.store.
						On("GetCompetitorsByCategory", fakedCategoryID).
						Return(slices.Values([]scores.Competitor{{
							ID:         fakedCompetitorID,
							CategoryID: fakedCategoryID,
							Score:      scores.Score{Tops: 2, Zone2s: 3, Zone1s: 4},
						}}))
					f.store.
						On("SaveResult", domain.Result{
							CompetitorID: fakedCompetitorID,
							Score:        tt.expected,
							Placement:    1,
							Timestamp:    time.Now(),
						}).Return()

					f.engine.RankCategory(fakedCategoryID)

					awaitExpectations(t)
				})
			})
		}
	})
}
