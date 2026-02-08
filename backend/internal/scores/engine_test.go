package scores_test

import (
	"iter"
	"slices"
	"testing"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/scores"
	"github.com/climblive/platform/backend/internal/utils/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

		fakedCompClass1ID := testutils.RandomResourceID[domain.CompClassID]()
		fakedCompClass2ID := testutils.RandomResourceID[domain.CompClassID]()

		fakedProblem1ID := testutils.RandomResourceID[domain.ProblemID]()
		fakedProblem2ID := testutils.RandomResourceID[domain.ProblemID]()

		fakedContender1ID := testutils.RandomResourceID[domain.ContenderID]()
		fakedContender2ID := testutils.RandomResourceID[domain.ContenderID]()
		fakedContender3ID := testutils.RandomResourceID[domain.ContenderID]()
		fakedContender4ID := testutils.RandomResourceID[domain.ContenderID]()

		f.store.On("SaveRules", scores.Rules{
			QualifyingProblems: 10,
			Finalists:          7,
			UsePoints:          true,
			PooledPoints:       true,
		}).Return()

		f.store.On("GetCompClassIDs").Return([]domain.CompClassID{
			fakedCompClass1ID,
			fakedCompClass2ID,
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
			On("GetContendersByCompClass", fakedCompClass1ID).Return(slices.Values([]scores.Contender{
			{
				ID: fakedContender1ID,
			},
			{
				ID: fakedContender2ID,
			},
		})).
			On("GetContendersByCompClass", fakedCompClass2ID).Return(slices.Values([]scores.Contender{
			{
				ID: fakedContender3ID,
			},
			{
				ID: fakedContender4ID,
			},
		}))

		effects := slices.Collect(f.engine.HandleRulesUpdated(domain.RulesUpdatedEvent{
			QualifyingProblems: 10,
			Finalists:          7,
			UsePoints:          true,
			PooledPoints:       true,
		}))

		require.ElementsMatch(t, effects, []scores.Effect{
			scores.EffectCalculateProblemValue{CompClassID: fakedCompClass1ID, ProblemID: fakedProblem1ID},
			scores.EffectCalculateProblemValue{CompClassID: fakedCompClass1ID, ProblemID: fakedProblem2ID},
			scores.EffectCalculateProblemValue{CompClassID: fakedCompClass2ID, ProblemID: fakedProblem1ID},
			scores.EffectCalculateProblemValue{CompClassID: fakedCompClass2ID, ProblemID: fakedProblem2ID},

			scores.EffectScoreContender{ContenderID: fakedContender1ID},
			scores.EffectScoreContender{ContenderID: fakedContender2ID},
			scores.EffectScoreContender{ContenderID: fakedContender3ID},
			scores.EffectScoreContender{ContenderID: fakedContender4ID},

			scores.EffectRankClass{CompClassID: fakedCompClass1ID},
			scores.EffectRankClass{CompClassID: fakedCompClass2ID},
		})

		awaitExpectations(t)
	})

	t.Run("Start", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedCompClass1ID := testutils.RandomResourceID[domain.CompClassID]()
		fakedCompClass2ID := testutils.RandomResourceID[domain.CompClassID]()

		fakedProblem1ID := testutils.RandomResourceID[domain.ProblemID]()
		fakedProblem2ID := testutils.RandomResourceID[domain.ProblemID]()

		fakedContender1ID := testutils.RandomResourceID[domain.ContenderID]()
		fakedContender2ID := testutils.RandomResourceID[domain.ContenderID]()
		fakedContender3ID := testutils.RandomResourceID[domain.ContenderID]()
		fakedContender4ID := testutils.RandomResourceID[domain.ContenderID]()

		f.store.On("GetCompClassIDs").Return([]domain.CompClassID{
			fakedCompClass1ID,
			fakedCompClass2ID,
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
			On("GetContendersByCompClass", fakedCompClass1ID).Return(slices.Values([]scores.Contender{
			{
				ID: fakedContender1ID,
			},
			{
				ID: fakedContender2ID,
			},
		})).
			On("GetContendersByCompClass", fakedCompClass2ID).Return(slices.Values([]scores.Contender{
			{
				ID: fakedContender3ID,
			},
			{
				ID: fakedContender4ID,
			},
		}))

		effects := slices.Collect(f.engine.Start())

		require.ElementsMatch(t, effects, []scores.Effect{
			scores.EffectCalculateProblemValue{CompClassID: fakedCompClass1ID, ProblemID: fakedProblem1ID},
			scores.EffectCalculateProblemValue{CompClassID: fakedCompClass1ID, ProblemID: fakedProblem2ID},
			scores.EffectCalculateProblemValue{CompClassID: fakedCompClass2ID, ProblemID: fakedProblem1ID},
			scores.EffectCalculateProblemValue{CompClassID: fakedCompClass2ID, ProblemID: fakedProblem2ID},

			scores.EffectScoreContender{ContenderID: fakedContender1ID},
			scores.EffectScoreContender{ContenderID: fakedContender2ID},
			scores.EffectScoreContender{ContenderID: fakedContender3ID},
			scores.EffectScoreContender{ContenderID: fakedContender4ID},

			scores.EffectRankClass{CompClassID: fakedCompClass1ID},
			scores.EffectRankClass{CompClassID: fakedCompClass2ID},
		})

		awaitExpectations(t)
	})

	t.Run("ContenderEntered", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedContenderID := testutils.RandomResourceID[domain.ContenderID]()
		fakedCompClassID := testutils.RandomResourceID[domain.CompClassID]()

		f.store.On("SaveContender", scores.Contender{
			ID:          fakedContenderID,
			CompClassID: fakedCompClassID,
		}).Return()

		effects := slices.Collect(f.engine.HandleContenderEntered(domain.ContenderEnteredEvent{
			ContenderID: fakedContenderID,
			CompClassID: fakedCompClassID,
		}))

		require.ElementsMatch(t, effects, []scores.Effect{
			scores.EffectRankClass{CompClassID: fakedCompClassID},
		})

		awaitExpectations(t)
	})

	t.Run("ContenderSwitchedClass_ContenderNotFound", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedContenderID := testutils.RandomResourceID[domain.ContenderID]()
		fakedCompClassID := testutils.RandomResourceID[domain.CompClassID]()

		f.store.
			On("GetContender", fakedContenderID).
			Return(scores.Contender{}, false)

		effects := f.engine.HandleContenderSwitchedClass(
			domain.ContenderSwitchedClassEvent{
				ContenderID: fakedContenderID,
				CompClassID: fakedCompClassID,
			})

		assert.Empty(t, effects)

		awaitExpectations(t)
	})

	t.Run("ContenderSwitchedClass_SameClass", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedContenderID := testutils.RandomResourceID[domain.ContenderID]()
		fakedCompClassID := testutils.RandomResourceID[domain.CompClassID]()

		f.store.
			On("GetContender", fakedContenderID).
			Return(scores.Contender{
				ID:          fakedContenderID,
				CompClassID: fakedCompClassID,
			}, true)

		effects := f.engine.HandleContenderSwitchedClass(domain.ContenderSwitchedClassEvent{
			ContenderID: fakedContenderID,
			CompClassID: fakedCompClassID,
		})

		assert.Empty(t, effects)

		awaitExpectations(t)
	})

	t.Run("ContenderSwitchedClass", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedContenderID := testutils.RandomResourceID[domain.ContenderID]()
		fakedOldCompClassID := testutils.RandomResourceID[domain.CompClassID]()
		fakedNewCompClassID := testutils.RandomResourceID[domain.CompClassID]()

		fakedProblem1ID := testutils.RandomResourceID[domain.ProblemID]()
		fakedProblem2ID := testutils.RandomResourceID[domain.ProblemID]()
		fakedProblem3ID := testutils.RandomResourceID[domain.ProblemID]()

		f.store.
			On("GetContender", fakedContenderID).
			Return(scores.Contender{
				ID:                  fakedContenderID,
				CompClassID:         fakedOldCompClassID,
				Disqualified:        false,
				WithdrawnFromFinals: false,
				Score:               123,
			}, true)

		f.store.On("SaveContender", scores.Contender{
			ID:                  fakedContenderID,
			CompClassID:         fakedNewCompClassID,
			Disqualified:        false,
			WithdrawnFromFinals: false,
			Score:               123,
		}).Return()

		f.store.On("GetTicksByContender", fakedContenderID).
			Return(slices.Values([]scores.Tick{
				{
					ProblemID: fakedProblem1ID,
				},
				{
					ProblemID: fakedProblem2ID,
				},
				{
					ProblemID: fakedProblem3ID,
				},
			}))

		effects := slices.Collect(f.engine.HandleContenderSwitchedClass(domain.ContenderSwitchedClassEvent{
			ContenderID: fakedContenderID,
			CompClassID: fakedNewCompClassID,
		}))

		require.ElementsMatch(t, effects, []scores.Effect{
			scores.EffectCalculateProblemValue{CompClassID: fakedOldCompClassID, ProblemID: fakedProblem1ID},
			scores.EffectCalculateProblemValue{CompClassID: fakedOldCompClassID, ProblemID: fakedProblem2ID},
			scores.EffectCalculateProblemValue{CompClassID: fakedOldCompClassID, ProblemID: fakedProblem3ID},

			scores.EffectCalculateProblemValue{CompClassID: fakedNewCompClassID, ProblemID: fakedProblem1ID},
			scores.EffectCalculateProblemValue{CompClassID: fakedNewCompClassID, ProblemID: fakedProblem2ID},
			scores.EffectCalculateProblemValue{CompClassID: fakedNewCompClassID, ProblemID: fakedProblem3ID},

			scores.EffectRankClass{CompClassID: fakedOldCompClassID},
			scores.EffectRankClass{CompClassID: fakedNewCompClassID},
		})

		awaitExpectations(t)
	})

	t.Run("ContenderWithdrewFromFinals_ContenderNotFound", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedContenderID := testutils.RandomResourceID[domain.ContenderID]()

		f.store.
			On("GetContender", fakedContenderID).
			Return(scores.Contender{}, false)

		effects := f.engine.HandleContenderWithdrewFromFinals(domain.ContenderWithdrewFromFinalsEvent{
			ContenderID: fakedContenderID,
		})

		assert.Empty(t, effects)

		awaitExpectations(t)
	})

	t.Run("ContenderWithdrewFromFinals", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedContenderID := testutils.RandomResourceID[domain.ContenderID]()
		fakedCompClassID := testutils.RandomResourceID[domain.CompClassID]()

		f.store.
			On("GetContender", fakedContenderID).
			Return(scores.Contender{
				ID:           fakedContenderID,
				CompClassID:  fakedCompClassID,
				Disqualified: true,
				Score:        123,
			}, true)

		f.store.On("SaveContender", scores.Contender{
			ID:                  fakedContenderID,
			CompClassID:         fakedCompClassID,
			Disqualified:        true,
			WithdrawnFromFinals: true,
			Score:               123,
		}).Return()

		effects := slices.Collect(f.engine.HandleContenderWithdrewFromFinals(domain.ContenderWithdrewFromFinalsEvent{
			ContenderID: fakedContenderID,
		}))

		require.ElementsMatch(t, effects, []scores.Effect{
			scores.EffectRankClass{CompClassID: fakedCompClassID},
		})

		awaitExpectations(t)
	})

	t.Run("ContenderReenteredFinals_ContenderNotFound", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedContenderID := testutils.RandomResourceID[domain.ContenderID]()

		f.store.
			On("GetContender", fakedContenderID).
			Return(scores.Contender{}, false)

		effects := f.engine.HandleContenderReenteredFinals(domain.ContenderReenteredFinalsEvent{
			ContenderID: fakedContenderID,
		})

		assert.Empty(t, effects)

		awaitExpectations(t)
	})

	t.Run("ContenderReenteredFinals", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedContenderID := testutils.RandomResourceID[domain.ContenderID]()
		fakedCompClassID := testutils.RandomResourceID[domain.CompClassID]()

		f.store.
			On("GetContender", fakedContenderID).
			Return(scores.Contender{
				ID:                  fakedContenderID,
				CompClassID:         fakedCompClassID,
				Disqualified:        true,
				WithdrawnFromFinals: true,
				Score:               123,
			}, true)

		f.store.On("SaveContender", scores.Contender{
			ID:                  fakedContenderID,
			CompClassID:         fakedCompClassID,
			Disqualified:        true,
			WithdrawnFromFinals: false,
			Score:               123,
		}).Return()

		effects := slices.Collect(f.engine.HandleContenderReenteredFinals(domain.ContenderReenteredFinalsEvent{
			ContenderID: fakedContenderID,
		}))

		require.ElementsMatch(t, effects, []scores.Effect{
			scores.EffectRankClass{CompClassID: fakedCompClassID},
		})

		awaitExpectations(t)
	})

	t.Run("ContenderDisqualified_ContenderNotFound", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedContenderID := testutils.RandomResourceID[domain.ContenderID]()

		f.store.
			On("GetContender", fakedContenderID).
			Return(scores.Contender{}, false)

		effects := f.engine.HandleContenderDisqualified(domain.ContenderDisqualifiedEvent{
			ContenderID: fakedContenderID,
		})

		assert.Empty(t, effects)

		awaitExpectations(t)
	})

	t.Run("ContenderDisqualified", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedContenderID := testutils.RandomResourceID[domain.ContenderID]()
		fakedCompClassID := testutils.RandomResourceID[domain.CompClassID]()

		fakedProblem1ID := testutils.RandomResourceID[domain.ProblemID]()
		fakedProblem2ID := testutils.RandomResourceID[domain.ProblemID]()
		fakedProblem3ID := testutils.RandomResourceID[domain.ProblemID]()

		f.store.
			On("GetContender", fakedContenderID).
			Return(scores.Contender{
				ID:                  fakedContenderID,
				CompClassID:         fakedCompClassID,
				WithdrawnFromFinals: true,
				Score:               123,
			}, true)

		f.store.
			On("GetTicksByContender", fakedContenderID).
			Return(slices.Values([]scores.Tick{
				{
					ProblemID: fakedProblem1ID,
				},
				{
					ProblemID: fakedProblem2ID,
				},
				{
					ProblemID: fakedProblem3ID,
				},
			}))

		f.store.On("SaveContender", scores.Contender{
			ID:                  fakedContenderID,
			CompClassID:         fakedCompClassID,
			Disqualified:        true,
			WithdrawnFromFinals: true,
			Score:               123,
		}).Return()

		effects := slices.Collect(f.engine.HandleContenderDisqualified(domain.ContenderDisqualifiedEvent{
			ContenderID: fakedContenderID,
		}))

		require.ElementsMatch(t, effects, []scores.Effect{
			scores.EffectCalculateProblemValue{CompClassID: fakedCompClassID, ProblemID: fakedProblem1ID},
			scores.EffectCalculateProblemValue{CompClassID: fakedCompClassID, ProblemID: fakedProblem2ID},
			scores.EffectCalculateProblemValue{CompClassID: fakedCompClassID, ProblemID: fakedProblem3ID},
			scores.EffectScoreContender{ContenderID: fakedContenderID},
			scores.EffectRankClass{CompClassID: fakedCompClassID},
		})

		awaitExpectations(t)
	})

	t.Run("ContenderRequalified_ContenderNotFound", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedContenderID := testutils.RandomResourceID[domain.ContenderID]()

		f.store.
			On("GetContender", fakedContenderID).
			Return(scores.Contender{}, false)

		effects := f.engine.HandleContenderRequalified(domain.ContenderRequalifiedEvent{
			ContenderID: fakedContenderID,
		})

		assert.Empty(t, effects)

		awaitExpectations(t)
	})

	t.Run("ContenderRequalified", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedContenderID := testutils.RandomResourceID[domain.ContenderID]()
		fakedCompClassID := testutils.RandomResourceID[domain.CompClassID]()

		fakedProblem1ID := testutils.RandomResourceID[domain.ProblemID]()
		fakedProblem2ID := testutils.RandomResourceID[domain.ProblemID]()
		fakedProblem3ID := testutils.RandomResourceID[domain.ProblemID]()

		f.store.
			On("GetContender", fakedContenderID).
			Return(scores.Contender{
				ID:                  fakedContenderID,
				CompClassID:         fakedCompClassID,
				Disqualified:        true,
				WithdrawnFromFinals: true,
				Score:               0,
			}, true)

		f.store.
			On("GetTicksByContender", fakedContenderID).
			Return(slices.Values([]scores.Tick{
				{
					ProblemID: fakedProblem1ID,
				},
				{
					ProblemID: fakedProblem2ID,
				},
				{
					ProblemID: fakedProblem3ID,
				},
			}))

		f.store.On("SaveContender", scores.Contender{
			ID:                  fakedContenderID,
			CompClassID:         fakedCompClassID,
			Disqualified:        false,
			WithdrawnFromFinals: true,
			Score:               0,
		}).Return()

		effects := slices.Collect(f.engine.HandleContenderRequalified(domain.ContenderRequalifiedEvent{
			ContenderID: fakedContenderID,
		}))

		require.ElementsMatch(t, effects, []scores.Effect{
			scores.EffectCalculateProblemValue{CompClassID: fakedCompClassID, ProblemID: fakedProblem1ID},
			scores.EffectCalculateProblemValue{CompClassID: fakedCompClassID, ProblemID: fakedProblem2ID},
			scores.EffectCalculateProblemValue{CompClassID: fakedCompClassID, ProblemID: fakedProblem3ID},
			scores.EffectScoreContender{ContenderID: fakedContenderID},
			scores.EffectRankClass{CompClassID: fakedCompClassID},
		})

		awaitExpectations(t)
	})

	t.Run("ProblemAdded", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedProblemID := testutils.RandomResourceID[domain.ProblemID]()

		f.store.
			On("SaveProblem", scores.Problem{
				ID: fakedProblemID,
				ProblemValue: domain.ProblemValue{
					PointsTop:   100,
					PointsZone1: 50,
					PointsZone2: 75,
					FlashBonus:  10,
				},
			}).
			Return()

		f.engine.HandleProblemAdded(domain.ProblemAddedEvent{
			ProblemID: fakedProblemID,
			ProblemValue: domain.ProblemValue{
				PointsTop:   100,
				PointsZone1: 50,
				PointsZone2: 75,
				FlashBonus:  10,
			},
		})

		awaitExpectations(t)
	})

	t.Run("AscentRegistered_ContenderNotFound", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedContenderID := testutils.RandomResourceID[domain.ContenderID]()
		fakedProblemID := testutils.RandomResourceID[domain.ProblemID]()

		f.store.
			On("GetContender", fakedContenderID).
			Return(scores.Contender{}, false)

		f.engine.HandleAscentRegistered(domain.AscentRegisteredEvent{
			ContenderID:   fakedContenderID,
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

		fakedContenderID := testutils.RandomResourceID[domain.ContenderID]()
		fakedCompClassID := testutils.RandomResourceID[domain.CompClassID]()
		fakedProblemID := testutils.RandomResourceID[domain.ProblemID]()

		f.store.
			On("GetContender", fakedContenderID).
			Return(scores.Contender{
				ID:           fakedContenderID,
				CompClassID:  fakedCompClassID,
				Disqualified: true,
			}, true)

		f.store.
			On("SaveTick", fakedContenderID, scores.Tick{
				ContenderID:   fakedContenderID,
				ProblemID:     fakedProblemID,
				Top:           true,
				AttemptsTop:   5,
				Zone1:         true,
				AttemptsZone1: 2,
				Zone2:         true,
				AttemptsZone2: 3,
			}).
			Return()

		f.engine.HandleAscentRegistered(domain.AscentRegisteredEvent{
			ContenderID:   fakedContenderID,
			ProblemID:     fakedProblemID,
			Top:           true,
			AttemptsTop:   5,
			Zone1:         true,
			AttemptsZone1: 2,
			Zone2:         true,
			AttemptsZone2: 3,
		})

		awaitExpectations(t)
	})

	t.Run("AscentRegistered", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedContenderID := testutils.RandomResourceID[domain.ContenderID]()
		fakedCompClassID := testutils.RandomResourceID[domain.CompClassID]()
		fakedProblemID := testutils.RandomResourceID[domain.ProblemID]()

		f.store.
			On("GetContender", fakedContenderID).
			Return(scores.Contender{
				ID:          fakedContenderID,
				CompClassID: fakedCompClassID,
			}, true)

		f.store.
			On("SaveTick", fakedContenderID, scores.Tick{
				ContenderID:   fakedContenderID,
				ProblemID:     fakedProblemID,
				Top:           true,
				AttemptsTop:   5,
				Zone1:         true,
				AttemptsZone1: 2,
				Zone2:         true,
				AttemptsZone2: 3,
			}).
			Return()

		f.engine.HandleAscentRegistered(domain.AscentRegisteredEvent{
			ContenderID:   fakedContenderID,
			ProblemID:     fakedProblemID,
			Top:           true,
			AttemptsTop:   5,
			Zone1:         true,
			AttemptsZone1: 2,
			Zone2:         true,
			AttemptsZone2: 3,
		})

		awaitExpectations(t)
	})

	t.Run("AscentDeregistered_ContenderNotFound", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedContenderID := testutils.RandomResourceID[domain.ContenderID]()
		fakedProblemID := testutils.RandomResourceID[domain.ProblemID]()

		f.store.
			On("GetContender", fakedContenderID).
			Return(scores.Contender{}, false)

		f.engine.HandleAscentDeregistered(domain.AscentDeregisteredEvent{
			ContenderID: fakedContenderID,
			ProblemID:   fakedProblemID,
		})

		awaitExpectations(t)
	})

	t.Run("AscentDeregistered_Disqualified", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedContenderID := testutils.RandomResourceID[domain.ContenderID]()
		fakedCompClassID := testutils.RandomResourceID[domain.CompClassID]()
		fakedProblemID := testutils.RandomResourceID[domain.ProblemID]()

		f.store.
			On("GetContender", fakedContenderID).
			Return(scores.Contender{
				ID:           fakedContenderID,
				CompClassID:  fakedCompClassID,
				Disqualified: true,
			}, true)

		f.store.
			On("DeleteTick", fakedContenderID, fakedProblemID).
			Return()

		f.engine.HandleAscentDeregistered(domain.AscentDeregisteredEvent{
			ContenderID: fakedContenderID,
			ProblemID:   fakedProblemID,
		})

		awaitExpectations(t)
	})

	t.Run("AscentDeregistered", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedContenderID := testutils.RandomResourceID[domain.ContenderID]()
		fakedCompClassID := testutils.RandomResourceID[domain.CompClassID]()
		fakedProblemID := testutils.RandomResourceID[domain.ProblemID]()

		f.store.
			On("GetContender", fakedContenderID).
			Return(scores.Contender{
				ID:          fakedContenderID,
				CompClassID: fakedCompClassID,
			}, true)

		f.store.
			On("DeleteTick", fakedContenderID, fakedProblemID).
			Return()

		f.engine.HandleAscentDeregistered(domain.AscentDeregisteredEvent{
			ContenderID: fakedContenderID,
			ProblemID:   fakedProblemID,
		})

		awaitExpectations(t)
	})

	t.Run("GetDirtyScores", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		now := time.Now()

		fakedContenderID1 := testutils.RandomResourceID[domain.ContenderID]()
		fakedContenderID2 := testutils.RandomResourceID[domain.ContenderID]()
		fakedContenderID3 := testutils.RandomResourceID[domain.ContenderID]()

		fakedScores := []domain.Score{
			{
				ContenderID: fakedContenderID1,
				Timestamp:   now,
				Score:       100,
				Placement:   1,
				RankOrder:   0,
				Finalist:    true,
			},
			{
				ContenderID: fakedContenderID2,
				Timestamp:   now,
				Score:       200,
				Placement:   2,
				RankOrder:   1,
				Finalist:    true,
			},
			{
				ContenderID: fakedContenderID3,
				Timestamp:   now,
				Score:       300,
				Placement:   3,
				RankOrder:   2,
				Finalist:    false,
			},
		}

		f.store.
			On("GetDirtyScores").
			Return(fakedScores)

		scores := f.engine.GetDirtyScores()

		assert.Equal(t, fakedScores, scores)

		awaitExpectations(t)
	})

	t.Run("CalculateProblemValue_UsePointsFalse", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedCompClassID := testutils.RandomResourceID[domain.CompClassID]()
		fakedProblemID := testutils.RandomResourceID[domain.ProblemID]()

		f.store.
			On("GetRules").
			Return(scores.Rules{
				UsePoints: false,
			})

		effects := f.engine.CalculateProblemValue(fakedCompClassID, fakedProblemID)

		assert.Nil(t, effects)

		awaitExpectations(t)
	})

	t.Run("CalculateProblemValue_ProblemNotFound", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedCompClassID := testutils.RandomResourceID[domain.CompClassID]()
		fakedProblemID := testutils.RandomResourceID[domain.ProblemID]()

		f.store.
			On("GetRules").
			Return(scores.Rules{
				UsePoints: true,
			})

		f.store.
			On("GetProblem", fakedProblemID).
			Return(scores.Problem{}, false)

		effects := f.engine.CalculateProblemValue(fakedCompClassID, fakedProblemID)

		assert.Nil(t, effects)

		awaitExpectations(t)
	})

	t.Run("CalculateProblemValue_NonPooledPoints", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedCompClassID := testutils.RandomResourceID[domain.CompClassID]()
		fakedProblemID := testutils.RandomResourceID[domain.ProblemID]()

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
				ProblemValue: domain.ProblemValue{
					PointsZone1: 100,
					PointsZone2: 200,
					PointsTop:   500,
					FlashBonus:  50,
				},
			}, true)

		f.store.
			On("SaveProblemValue", fakedCompClassID, fakedProblemID, scores.ProblemValue{
				ProblemID:   fakedProblemID,
				CompClassID: fakedCompClassID,
				ProblemValue: domain.ProblemValue{
					PointsZone1: 100,
					PointsZone2: 200,
					PointsTop:   500,
					FlashBonus:  50,
				},
			}).Return()

		effects := f.engine.CalculateProblemValue(fakedCompClassID, fakedProblemID)

		assert.Nil(t, effects)

		awaitExpectations(t)
	})

	t.Run("CalculateProblemValue_PooledPoints_ValueUnchanged", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedCompClassID := testutils.RandomResourceID[domain.CompClassID]()
		fakedProblemID := testutils.RandomResourceID[domain.ProblemID]()

		fakedContender1ID := testutils.RandomResourceID[domain.ContenderID]()

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
				ProblemValue: domain.ProblemValue{
					PointsZone1: 100,
					PointsZone2: 200,
					PointsTop:   500,
					FlashBonus:  100,
				},
			}, true)

		f.store.
			On("GetTicksByProblem", fakedCompClassID, fakedProblemID).
			Return(slices.Values([]scores.Tick{
				{
					ContenderID: fakedContender1ID,
					ProblemID:   fakedProblemID,
					Zone1:       true,
					Zone2:       true,
					Top:         true,
					AttemptsTop: 1,
				},
			}))

		f.store.
			On("GetContender", fakedContender1ID).
			Return(scores.Contender{
				ID:           fakedContender1ID,
				Disqualified: false,
			}, true)

		f.store.
			On("SaveProblemValue", fakedCompClassID, fakedProblemID, scores.ProblemValue{
				ProblemID:   fakedProblemID,
				CompClassID: fakedCompClassID,
				ProblemValue: domain.ProblemValue{
					PointsZone1: 100,
					PointsZone2: 200,
					PointsTop:   500,
					FlashBonus:  100,
				},
			}).Return()

		effects := f.engine.CalculateProblemValue(fakedCompClassID, fakedProblemID)

		require.Nil(t, effects)

		awaitExpectations(t)
	})

	t.Run("CalculateProblemValue_PooledPoints_UpdatedValue", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedCompClassID := testutils.RandomResourceID[domain.CompClassID]()
		fakedProblemID := testutils.RandomResourceID[domain.ProblemID]()

		fakedContender1ID := testutils.RandomResourceID[domain.ContenderID]()
		fakedContender2ID := testutils.RandomResourceID[domain.ContenderID]()
		fakedContender3ID := testutils.RandomResourceID[domain.ContenderID]()
		fakedContender4ID := testutils.RandomResourceID[domain.ContenderID]()
		fakedContender5ID := testutils.RandomResourceID[domain.ContenderID]()

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
				ProblemValue: domain.ProblemValue{
					PointsZone1: 50,
					PointsZone2: 100,
					PointsTop:   500,
					FlashBonus:  100,
				},
			}, true)

		f.store.
			On("GetTicksByProblem", fakedCompClassID, fakedProblemID).
			Return(slices.Values([]scores.Tick{
				{
					ContenderID: fakedContender1ID,
					ProblemID:   fakedProblemID,
					Zone1:       true,
					Zone2:       true,
					Top:         true,
					AttemptsTop: 1,
				},
				{
					ContenderID: fakedContender2ID,
					ProblemID:   fakedProblemID,
					Zone1:       true,
					Zone2:       true,
					Top:         true,
					AttemptsTop: 1,
				},
				{
					ContenderID: fakedContender3ID,
					ProblemID:   fakedProblemID,
					Zone1:       true,
					Zone2:       true,
					Top:         true,
					AttemptsTop: 999,
				},
				{
					ContenderID: fakedContender4ID,
					ProblemID:   fakedProblemID,
					Zone1:       true,
					Zone2:       true,
					Top:         false,
					AttemptsTop: 999,
				},
				{
					ContenderID: fakedContender5ID,
					ProblemID:   fakedProblemID,
					Zone1:       true,
					Zone2:       false,
					Top:         false,
					AttemptsTop: 999,
				},
			}))

		f.store.
			On("GetContender", fakedContender1ID).
			Return(scores.Contender{
				ID:           fakedContender1ID,
				Disqualified: false,
			}, true).
			On("GetContender", fakedContender2ID).
			Return(scores.Contender{
				ID:           fakedContender2ID,
				Disqualified: false,
			}, true).
			On("GetContender", fakedContender3ID).
			Return(scores.Contender{
				ID:           fakedContender3ID,
				Disqualified: false,
			}, true).
			On("GetContender", fakedContender4ID).
			Return(scores.Contender{
				ID:           fakedContender4ID,
				Disqualified: false,
			}, true).
			On("GetContender", fakedContender5ID).
			Return(scores.Contender{
				ID:           fakedContender5ID,
				Disqualified: false,
			}, true)

		f.store.
			On("SaveProblemValue", fakedCompClassID, fakedProblemID, scores.ProblemValue{
				ProblemID:   fakedProblemID,
				CompClassID: fakedCompClassID,
				ProblemValue: domain.ProblemValue{
					PointsZone1: 10,
					PointsZone2: 25,
					PointsTop:   166,
					FlashBonus:  50,
				},
			}).Return()

		effects := slices.Collect(f.engine.CalculateProblemValue(fakedCompClassID, fakedProblemID))

		require.ElementsMatch(t, effects, []scores.Effect{
			scores.EffectScoreContender{ContenderID: fakedContender1ID},
			scores.EffectScoreContender{ContenderID: fakedContender2ID},
			scores.EffectScoreContender{ContenderID: fakedContender3ID},
			scores.EffectScoreContender{ContenderID: fakedContender4ID},
			scores.EffectScoreContender{ContenderID: fakedContender5ID},
		})

		awaitExpectations(t)
	})

	t.Run("CalculateProblemValue_PooledPoints_ExcludeDisqualified", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedCompClassID := testutils.RandomResourceID[domain.CompClassID]()
		fakedProblemID := testutils.RandomResourceID[domain.ProblemID]()

		fakedContender1ID := testutils.RandomResourceID[domain.ContenderID]()
		fakedContender2ID := testutils.RandomResourceID[domain.ContenderID]()
		fakedContender3ID := testutils.RandomResourceID[domain.ContenderID]()

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
				ProblemValue: domain.ProblemValue{
					PointsZone1: 50,
					PointsZone2: 100,
					PointsTop:   500,
					FlashBonus:  25,
				},
			}, true)

		f.store.
			On("GetTicksByProblem", fakedCompClassID, fakedProblemID).
			Return(slices.Values([]scores.Tick{
				{
					ContenderID: fakedContender1ID,
					ProblemID:   fakedProblemID,
					Zone1:       true,
					Zone2:       true,
					Top:         true,
					AttemptsTop: 1,
				},
				{
					ContenderID: fakedContender2ID,
					ProblemID:   fakedProblemID,
					Zone1:       true,
					Zone2:       true,
					Top:         true,
					AttemptsTop: 1,
				},
				{
					ContenderID: fakedContender3ID,
					ProblemID:   fakedProblemID,
					Zone1:       true,
					Zone2:       true,
					Top:         true,
					AttemptsTop: 1,
				},
			}))

		f.store.
			On("GetContender", fakedContender1ID).
			Return(scores.Contender{
				ID:           fakedContender1ID,
				Disqualified: false,
			}, true).
			On("GetContender", fakedContender2ID).
			Return(scores.Contender{
				ID:           fakedContender2ID,
				Disqualified: false,
			}, true).
			On("GetContender", fakedContender3ID).
			Return(scores.Contender{
				ID:           fakedContender3ID,
				Disqualified: true,
			}, true)

		f.store.
			On("SaveProblemValue", fakedCompClassID, fakedProblemID, scores.ProblemValue{
				ProblemID:   fakedProblemID,
				CompClassID: fakedCompClassID,
				ProblemValue: domain.ProblemValue{
					PointsZone1: 25,
					PointsZone2: 50,
					PointsTop:   250,
					FlashBonus:  12,
				},
			}).Return()

		effects := slices.Collect(f.engine.CalculateProblemValue(fakedCompClassID, fakedProblemID))

		require.ElementsMatch(t, effects, []scores.Effect{
			scores.EffectScoreContender{ContenderID: fakedContender1ID},
			scores.EffectScoreContender{ContenderID: fakedContender2ID},
		})

		awaitExpectations(t)
	})

	t.Run("ScoreContender_ContenderNotFound", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedContenderID := testutils.RandomResourceID[domain.ContenderID]()

		f.store.
			On("GetContender", fakedContenderID).
			Return(scores.Contender{}, false)

		effects := f.engine.ScoreContender(fakedContenderID)

		assert.Nil(t, effects)

		awaitExpectations(t)
	})

	t.Run("ScoreContender_Disqualified", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedContenderID := testutils.RandomResourceID[domain.ContenderID]()
		fakedCompClassID := testutils.RandomResourceID[domain.CompClassID]()

		f.store.
			On("GetContender", fakedContenderID).
			Return(scores.Contender{
				ID:           fakedContenderID,
				CompClassID:  fakedCompClassID,
				Disqualified: true,
				Score:        100,
			}, true)

		f.store.
			On("SaveContender", scores.Contender{
				ID:           fakedContenderID,
				CompClassID:  fakedCompClassID,
				Disqualified: true,
				Score:        0,
			}).Return()

		effects := slices.Collect(f.engine.ScoreContender(fakedContenderID))

		require.ElementsMatch(t, effects, []scores.Effect{
			scores.EffectRankClass{CompClassID: fakedCompClassID},
		})

		awaitExpectations(t)
	})

	t.Run("ScoreContender_DisqualifiedNoScoreChange", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedContenderID := testutils.RandomResourceID[domain.ContenderID]()
		fakedCompClassID := testutils.RandomResourceID[domain.CompClassID]()

		f.store.
			On("GetContender", fakedContenderID).
			Return(scores.Contender{
				ID:           fakedContenderID,
				CompClassID:  fakedCompClassID,
				Disqualified: true,
				Score:        0,
			}, true)

		effects := f.engine.ScoreContender(fakedContenderID)

		assert.Nil(t, effects)

		awaitExpectations(t)
	})

	t.Run("ScoreContender_WithTicks", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedContenderID := testutils.RandomResourceID[domain.ContenderID]()
		fakedCompClassID := testutils.RandomResourceID[domain.CompClassID]()
		fakedProblem1ID := testutils.RandomResourceID[domain.ProblemID]()
		fakedProblem2ID := testutils.RandomResourceID[domain.ProblemID]()

		f.store.
			On("GetContender", fakedContenderID).
			Return(scores.Contender{
				ID:          fakedContenderID,
				CompClassID: fakedCompClassID,
				Score:       0,
			}, true)

		f.store.
			On("GetTicksByContender", fakedContenderID).
			Return(slices.Values([]scores.Tick{
				{
					ContenderID: fakedContenderID,
					ProblemID:   fakedProblem1ID,
					Top:         true,
					AttemptsTop: 1,
				},
				{
					ContenderID: fakedContenderID,
					ProblemID:   fakedProblem2ID,
					Top:         true,
					AttemptsTop: 2,
				},
			}))

		f.store.
			On("GetProblemValue", fakedCompClassID, fakedProblem1ID).
			Return(scores.ProblemValue{
				ProblemID:   fakedProblem1ID,
				CompClassID: fakedCompClassID,
				ProblemValue: domain.ProblemValue{
					PointsTop:  100,
					FlashBonus: 10,
				},
			}, true)

		f.store.
			On("GetProblemValue", fakedCompClassID, fakedProblem2ID).
			Return(scores.ProblemValue{
				ProblemID:   fakedProblem2ID,
				CompClassID: fakedCompClassID,
				ProblemValue: domain.ProblemValue{
					PointsTop:  100,
					FlashBonus: 10,
				},
			}, true)

		f.store.
			On("GetRules").
			Return(scores.Rules{
				QualifyingProblems: 0,
				UsePoints:          false,
			})

		f.store.
			On("SaveContender", scores.Contender{
				ID:          fakedContenderID,
				CompClassID: fakedCompClassID,
				Score:       210,
			}).Return()

		effects := slices.Collect(f.engine.ScoreContender(fakedContenderID))

		require.ElementsMatch(t, effects, []scores.Effect{
			scores.EffectRankClass{CompClassID: fakedCompClassID},
		})

		awaitExpectations(t)
	})

	t.Run("ScoreContender_WithProblemLimit", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedContenderID := testutils.RandomResourceID[domain.ContenderID]()
		fakedCompClassID := testutils.RandomResourceID[domain.CompClassID]()
		fakedProblem1ID := testutils.RandomResourceID[domain.ProblemID]()
		fakedProblem2ID := testutils.RandomResourceID[domain.ProblemID]()
		fakedProblem3ID := testutils.RandomResourceID[domain.ProblemID]()

		f.store.
			On("GetContender", fakedContenderID).
			Return(scores.Contender{
				ID:          fakedContenderID,
				CompClassID: fakedCompClassID,
				Score:       0,
			}, true)

		f.store.
			On("GetTicksByContender", fakedContenderID).
			Return(slices.Values([]scores.Tick{
				{
					ContenderID: fakedContenderID,
					ProblemID:   fakedProblem1ID,
					Top:         true,
					AttemptsTop: 1,
				},
				{
					ContenderID: fakedContenderID,
					ProblemID:   fakedProblem2ID,
					Top:         true,
					AttemptsTop: 1,
				},
				{
					ContenderID: fakedContenderID,
					ProblemID:   fakedProblem3ID,
					Top:         true,
					AttemptsTop: 1,
				},
			}))

		f.store.
			On("GetProblemValue", fakedCompClassID, fakedProblem1ID).
			Return(scores.ProblemValue{
				ProblemID:   fakedProblem1ID,
				CompClassID: fakedCompClassID,
				ProblemValue: domain.ProblemValue{
					PointsTop:  50,
					FlashBonus: 5,
				},
			}, true)

		f.store.
			On("GetProblemValue", fakedCompClassID, fakedProblem2ID).
			Return(scores.ProblemValue{
				ProblemID:   fakedProblem2ID,
				CompClassID: fakedCompClassID,
				ProblemValue: domain.ProblemValue{
					PointsTop:  100,
					FlashBonus: 10,
				},
			}, true)

		f.store.
			On("GetProblemValue", fakedCompClassID, fakedProblem3ID).
			Return(scores.ProblemValue{
				ProblemID:   fakedProblem3ID,
				CompClassID: fakedCompClassID,
				ProblemValue: domain.ProblemValue{
					PointsTop:  75,
					FlashBonus: 7,
				},
			}, true)

		f.store.
			On("GetRules").
			Return(scores.Rules{
				QualifyingProblems: 2,
				UsePoints:          true,
			})

		f.store.
			On("SaveContender", scores.Contender{
				ID:          fakedContenderID,
				CompClassID: fakedCompClassID,
				Score:       192,
			}).Return()

		effects := slices.Collect(f.engine.ScoreContender(fakedContenderID))

		require.ElementsMatch(t, effects, []scores.Effect{
			scores.EffectRankClass{CompClassID: fakedCompClassID},
		})

		awaitExpectations(t)
	})

	t.Run("ScoreContender_NoScoreChange", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedContenderID := testutils.RandomResourceID[domain.ContenderID]()
		fakedCompClassID := testutils.RandomResourceID[domain.CompClassID]()
		fakedProblem1ID := testutils.RandomResourceID[domain.ProblemID]()

		f.store.
			On("GetContender", fakedContenderID).
			Return(scores.Contender{
				ID:          fakedContenderID,
				CompClassID: fakedCompClassID,
				Score:       110,
			}, true)

		f.store.
			On("GetTicksByContender", fakedContenderID).
			Return(slices.Values([]scores.Tick{
				{
					ContenderID: fakedContenderID,
					ProblemID:   fakedProblem1ID,
					Top:         true,
					AttemptsTop: 1,
				},
			}))

		f.store.
			On("GetProblemValue", fakedCompClassID, fakedProblem1ID).
			Return(scores.ProblemValue{
				ProblemID:   fakedProblem1ID,
				CompClassID: fakedCompClassID,
				ProblemValue: domain.ProblemValue{
					PointsTop:  100,
					FlashBonus: 10,
				},
			}, true)

		f.store.
			On("GetRules").
			Return(scores.Rules{
				QualifyingProblems: 0,
				UsePoints:          false,
			})

		effects := f.engine.ScoreContender(fakedContenderID)

		assert.Nil(t, effects)

		awaitExpectations(t)
	})

	t.Run("ScoreContender_ProblemValueNotFound", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedContenderID := testutils.RandomResourceID[domain.ContenderID]()
		fakedCompClassID := testutils.RandomResourceID[domain.CompClassID]()
		fakedProblem1ID := testutils.RandomResourceID[domain.ProblemID]()
		fakedProblem2ID := testutils.RandomResourceID[domain.ProblemID]()

		f.store.
			On("GetContender", fakedContenderID).
			Return(scores.Contender{
				ID:          fakedContenderID,
				CompClassID: fakedCompClassID,
				Score:       0,
			}, true)

		f.store.
			On("GetTicksByContender", fakedContenderID).
			Return(slices.Values([]scores.Tick{
				{
					ContenderID: fakedContenderID,
					ProblemID:   fakedProblem1ID,
					Top:         true,
					AttemptsTop: 1,
				},
				{
					ContenderID: fakedContenderID,
					ProblemID:   fakedProblem2ID,
					Top:         true,
					AttemptsTop: 1,
				},
			}))

		f.store.
			On("GetProblemValue", fakedCompClassID, fakedProblem1ID).
			Return(scores.ProblemValue{
				ProblemID:   fakedProblem1ID,
				CompClassID: fakedCompClassID,
				ProblemValue: domain.ProblemValue{
					PointsTop:  100,
					FlashBonus: 10,
				},
			}, true)

		f.store.
			On("GetProblemValue", fakedCompClassID, fakedProblem2ID).
			Return(scores.ProblemValue{}, false)

		f.store.
			On("GetRules").
			Return(scores.Rules{
				QualifyingProblems: 0,
				UsePoints:          false,
			})

		f.store.
			On("SaveContender", scores.Contender{
				ID:          fakedContenderID,
				CompClassID: fakedCompClassID,
				Score:       110,
			}).Return()

		effects := slices.Collect(f.engine.ScoreContender(fakedContenderID))

		require.ElementsMatch(t, effects, []scores.Effect{
			scores.EffectRankClass{CompClassID: fakedCompClassID},
		})

		awaitExpectations(t)
	})
}

type engineStoreMock struct {
	mock.Mock
}

func (m *engineStoreMock) GetRules() scores.Rules {
	args := m.Called()
	return args.Get(0).(scores.Rules)
}

func (m *engineStoreMock) SaveRules(rules scores.Rules) {
	m.Called(rules)
}

func (m *engineStoreMock) GetContender(contenderID domain.ContenderID) (scores.Contender, bool) {
	args := m.Called(contenderID)
	return args.Get(0).(scores.Contender), args.Bool(1)
}

func (m *engineStoreMock) GetContendersByCompClass(compClassID domain.CompClassID) iter.Seq[scores.Contender] {
	args := m.Called(compClassID)
	return args.Get(0).(iter.Seq[scores.Contender])
}

func (m *engineStoreMock) GetAllContenders() iter.Seq[scores.Contender] {
	args := m.Called()
	return args.Get(0).(iter.Seq[scores.Contender])
}

func (m *engineStoreMock) SaveContender(contender scores.Contender) {
	m.Called(contender)
}

func (m *engineStoreMock) GetCompClassIDs() []domain.CompClassID {
	args := m.Called()
	return args.Get(0).([]domain.CompClassID)
}

func (m *engineStoreMock) GetTicksByContender(contenderID domain.ContenderID) iter.Seq[scores.Tick] {
	args := m.Called(contenderID)
	return args.Get(0).(iter.Seq[scores.Tick])
}

func (m *engineStoreMock) GetTicksByProblem(compClassID domain.CompClassID, problemID domain.ProblemID) iter.Seq[scores.Tick] {
	args := m.Called(compClassID, problemID)
	return args.Get(0).(iter.Seq[scores.Tick])
}

func (m *engineStoreMock) SaveTick(contenderID domain.ContenderID, tick scores.Tick) {
	m.Called(contenderID, tick)
}

func (m *engineStoreMock) DeleteTick(contenderID domain.ContenderID, problemID domain.ProblemID) {
	m.Called(contenderID, problemID)
}

func (m *engineStoreMock) GetProblem(problemID domain.ProblemID) (scores.Problem, bool) {
	args := m.Called(problemID)
	return args.Get(0).(scores.Problem), args.Bool(1)
}

func (m *engineStoreMock) SaveProblem(problem scores.Problem) {
	m.Called(problem)
}

func (m *engineStoreMock) GetAllProblems() iter.Seq[scores.Problem] {
	args := m.Called()
	return args.Get(0).(iter.Seq[scores.Problem])
}

func (m *engineStoreMock) GetProblemValue(compClassID domain.CompClassID, problemID domain.ProblemID) (scores.ProblemValue, bool) {
	args := m.Called(compClassID, problemID)
	return args.Get(0).(scores.ProblemValue), args.Bool(1)
}

func (m *engineStoreMock) SaveProblemValue(compClassID domain.CompClassID, problemID domain.ProblemID, value scores.ProblemValue) {
	m.Called(compClassID, problemID, value)
}

func (m *engineStoreMock) GetDirtyProblemValues() []scores.ProblemValue {
	args := m.Called()
	return args.Get(0).([]scores.ProblemValue)
}

func (m *engineStoreMock) SaveScore(score domain.Score) {
	m.Called(score)
}

func (m *engineStoreMock) GetDirtyScores() []domain.Score {
	args := m.Called()
	return args.Get(0).([]domain.Score)
}
