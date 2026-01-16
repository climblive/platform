package scores_test

import (
	"iter"
	"slices"
	"testing"
	"testing/synctest"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/scores"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDefaultScoreEngine(t *testing.T) {
	type fixture struct {
		store  *engineStoreMock
		engine *scores.DefaultScoreEngine
	}

	makeFixture := func() (fixture, func(t *testing.T)) {
		mockedStore := new(engineStoreMock)

		mockedStore.On("GetRules").Return(scores.Rules{
			QualifyingProblems: 10,
			Finalists:          7,
		})

		engine := scores.NewDefaultScoreEngine(mockedStore)

		awaitExpectations := func(t *testing.T) {
			mockedStore.AssertExpectations(t)
		}

		return fixture{
			store:  mockedStore,
			engine: engine,
		}, awaitExpectations
	}

	t.Run("ReplaceScoringRules", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			contender1ID := domain.ContenderID(uuid.New())
			contender2ID := domain.ContenderID(uuid.New())
			contender3ID := domain.ContenderID(uuid.New())
			class1ID := domain.CompClassID(uuid.New())
			class2ID := domain.CompClassID(uuid.New())
			class3ID := domain.CompClassID(uuid.New())

			f, awaitExpectations := makeFixture()

			f.store.On("SaveRules", scores.Rules{
				QualifyingProblems: 1,
				Finalists:          0,
			}).Return()

			f.store.On("GetAllContenders").
				Return(slices.Values([]scores.Contender{
					{ID: contender1ID, CompClassID: class1ID},
					{ID: contender2ID, CompClassID: class2ID},
					{ID: contender3ID, CompClassID: class3ID},
				}))

			f.store.
				On("GetTicks", contender1ID).
				Return(slices.Values([]scores.Tick{{Points: 100}, {Points: 200}, {Points: 300}})).
				On("GetTicks", contender2ID).
				Return(slices.Values([]scores.Tick{{Points: 400}, {Points: 500}})).
				On("GetTicks", contender3ID).
				Return(slices.Values([]scores.Tick{{Points: 600}}))

			f.store.
				On("SaveContender", scores.Contender{ID: contender1ID, CompClassID: class1ID, Score: 300}).Return().
				On("SaveContender", scores.Contender{ID: contender2ID, CompClassID: class2ID, Score: 500}).Return().
				On("SaveContender", scores.Contender{ID: contender3ID, CompClassID: class3ID, Score: 600}).Return()

			f.store.On("GetCompClassIDs").Return([]domain.CompClassID{class1ID, class2ID, class3ID})

			f.store.
				On("GetContendersByCompClass", class1ID).
				Return(slices.Values([]scores.Contender{{ID: contender1ID}})).
				On("GetContendersByCompClass", class2ID).
				Return(slices.Values([]scores.Contender{{ID: contender2ID}})).
				On("GetContendersByCompClass", class3ID).
				Return(slices.Values([]scores.Contender{{ID: contender3ID}}))

			f.store.On("SaveScore", domain.Score{Timestamp: time.Now(), ContenderID: contender1ID, Placement: 1}).Return()
			f.store.On("SaveScore", domain.Score{Timestamp: time.Now(), ContenderID: contender2ID, Placement: 1}).Return()
			f.store.On("SaveScore", domain.Score{Timestamp: time.Now(), ContenderID: contender3ID, Placement: 1}).Return()

			f.engine.HandleRulesUpdated(domain.RulesUpdatedEvent{
				QualifyingProblems: 1,
				Finalists:          0,
			})

			awaitExpectations(t)
		})
	})

	t.Run("Start", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			contender1ID := domain.ContenderID(uuid.New())
			contender2ID := domain.ContenderID(uuid.New())
			contender3ID := domain.ContenderID(uuid.New())
			class1ID := domain.CompClassID(uuid.New())
			class2ID := domain.CompClassID(uuid.New())
			class3ID := domain.CompClassID(uuid.New())
			problem1ID := domain.ProblemID(uuid.New())
			problem2ID := domain.ProblemID(uuid.New())
			problem3ID := domain.ProblemID(uuid.New())

			f, awaitExpectations := makeFixture()

			f.store.On("GetAllContenders").
				Return(slices.Values([]scores.Contender{
					{ID: contender1ID, CompClassID: class1ID},
					{ID: contender2ID, CompClassID: class2ID},
				}))

			f.store.
				On("GetTicks", contender1ID).
				Return(slices.Values([]scores.Tick{{ProblemID: problem1ID, Top: true}, {ProblemID: problem2ID, Top: true}, {ProblemID: problem3ID, Top: true}})).
				On("GetTicks", contender2ID).
				Return(slices.Values([]scores.Tick{{ProblemID: problem1ID, Zone1: true}, {ProblemID: problem2ID, Zone1: true}}))

			f.store.
				On("GetProblem", problem1ID).
				Return(scores.Problem{ID: problem1ID, PointsTop: 100, PointsZone1: 10}, true).
				On("GetProblem", problem2ID).
				Return(scores.Problem{ID: problem2ID, PointsTop: 200, PointsZone1: 20}, true).
				On("GetProblem", problem3ID).
				Return(scores.Problem{ID: problem3ID, PointsTop: 300}, true)

			f.store.
				On("SaveTick", contender1ID, scores.Tick{ProblemID: problem1ID, Top: true, Points: 100}).Return().
				On("SaveTick", contender1ID, scores.Tick{ProblemID: problem2ID, Top: true, Points: 200}).Return().
				On("SaveTick", contender1ID, scores.Tick{ProblemID: problem3ID, Top: true, Points: 300}).Return().
				On("SaveTick", contender2ID, scores.Tick{ProblemID: problem1ID, Zone1: true, Points: 10}).Return().
				On("SaveTick", contender2ID, scores.Tick{ProblemID: problem2ID, Zone1: true, Points: 20}).Return()

			f.store.
				On("SaveContender", scores.Contender{ID: contender1ID, CompClassID: class1ID, Score: 600}).Return().
				On("SaveContender", scores.Contender{ID: contender2ID, CompClassID: class2ID, Score: 30}).Return()

			f.store.On("GetCompClassIDs").Return([]domain.CompClassID{class1ID, class2ID, class3ID})

			f.store.
				On("GetContendersByCompClass", class1ID).
				Return(slices.Values([]scores.Contender{{ID: contender1ID, Score: 600}})).
				On("GetContendersByCompClass", class2ID).
				Return(slices.Values([]scores.Contender{{ID: contender2ID, Score: 30}})).
				On("GetContendersByCompClass", class3ID).
				Return(slices.Values([]scores.Contender{{ID: contender3ID, Score: 0}}))

			f.store.On("SaveScore", domain.Score{Timestamp: time.Now(), ContenderID: contender1ID, Score: 600, Placement: 1, Finalist: true, RankOrder: 0}).Return()
			f.store.On("SaveScore", domain.Score{Timestamp: time.Now(), ContenderID: contender2ID, Score: 30, Placement: 1, Finalist: true, RankOrder: 0}).Return()
			f.store.On("SaveScore", domain.Score{Timestamp: time.Now(), ContenderID: contender3ID, Score: 0, Placement: 1, Finalist: false, RankOrder: 0}).Return()

			f.engine.Start()

			awaitExpectations(t)
		})
	})

	t.Run("ContenderEntered", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			contender1ID := domain.ContenderID(uuid.New())
			contender2ID := domain.ContenderID(uuid.New())
			contender3ID := domain.ContenderID(uuid.New())
			class1ID := domain.CompClassID(uuid.New())

			f, awaitExpectations := makeFixture()

			f.store.On("SaveContender", scores.Contender{
				ID:          contender1ID,
				CompClassID: class1ID,
			}).Return()

			f.store.
				On("GetContendersByCompClass", class1ID).
				Return(slices.Values([]scores.Contender{{ID: contender1ID}, {ID: contender2ID}, {ID: contender3ID}}))

			f.store.On("SaveScore", domain.Score{Timestamp: time.Now(), ContenderID: contender1ID, Placement: 1, RankOrder: 0}).Return()
			f.store.On("SaveScore", domain.Score{Timestamp: time.Now(), ContenderID: contender2ID, Placement: 1, RankOrder: 1}).Return()
			f.store.On("SaveScore", domain.Score{Timestamp: time.Now(), ContenderID: contender3ID, Placement: 1, RankOrder: 2}).Return()

			f.engine.HandleContenderEntered(domain.ContenderEnteredEvent{
				ContenderID: contender1ID,
				CompClassID: class1ID,
			})

			awaitExpectations(t)
		})
	})

	t.Run("ContenderSwitchedClass_ContenderNotFound", func(t *testing.T) {
		contender1ID := domain.ContenderID(uuid.New())
		class1ID := domain.CompClassID(uuid.New())

		f, awaitExpectations := makeFixture()

		f.store.
			On("GetContender", contender1ID).
			Return(scores.Contender{}, false)

		f.engine.HandleContenderSwitchedClass(
			domain.ContenderSwitchedClassEvent{
				ContenderID: contender1ID,
				CompClassID: class1ID,
			})

		awaitExpectations(t)
	})

	t.Run("ContenderSwitchedClass_SameClass", func(t *testing.T) {
		contender1ID := domain.ContenderID(uuid.New())
		class1ID := domain.CompClassID(uuid.New())

		f, awaitExpectations := makeFixture()

		f.store.
			On("GetContender", contender1ID).
			Return(scores.Contender{
				ID:          contender1ID,
				CompClassID: class1ID,
			}, true)

		f.engine.HandleContenderSwitchedClass(domain.ContenderSwitchedClassEvent{
			ContenderID: contender1ID,
			CompClassID: class1ID,
		})

		awaitExpectations(t)
	})

	t.Run("ContenderSwitchedClass", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			contender1ID := domain.ContenderID(uuid.New())
			contender2ID := domain.ContenderID(uuid.New())
			contender3ID := domain.ContenderID(uuid.New())
			contender4ID := domain.ContenderID(uuid.New())
			class1ID := domain.CompClassID(uuid.New())
			class2ID := domain.CompClassID(uuid.New())

			f, awaitExpectations := makeFixture()

			f.store.
				On("GetContender", contender4ID).
				Return(scores.Contender{
					ID:                  contender4ID,
					CompClassID:         class1ID,
					Disqualified:        false,
					WithdrawnFromFinals: false,
					Score:               123,
				}, true)

			f.store.On("SaveContender", scores.Contender{
				ID:                  contender4ID,
				CompClassID:         class2ID,
				Disqualified:        false,
				WithdrawnFromFinals: false,
				Score:               123,
			}).Return()

			f.store.
				On("GetContendersByCompClass", class1ID).
				Return(slices.Values([]scores.Contender{{ID: contender1ID}, {ID: contender2ID}, {ID: contender3ID}})).
				On("GetContendersByCompClass", class2ID).
				Return(slices.Values([]scores.Contender{{ID: contender4ID}}))

			f.store.On("SaveScore", domain.Score{Timestamp: time.Now(), ContenderID: contender1ID, Placement: 1, RankOrder: 0}).Return()
			f.store.On("SaveScore", domain.Score{Timestamp: time.Now(), ContenderID: contender2ID, Placement: 1, RankOrder: 1}).Return()
			f.store.On("SaveScore", domain.Score{Timestamp: time.Now(), ContenderID: contender3ID, Placement: 1, RankOrder: 2}).Return()
			f.store.On("SaveScore", domain.Score{Timestamp: time.Now(), ContenderID: contender4ID, Placement: 1, RankOrder: 0}).Return()

			f.engine.HandleContenderSwitchedClass(domain.ContenderSwitchedClassEvent{
				ContenderID: contender4ID,
				CompClassID: class2ID,
			})

			awaitExpectations(t)
		})
	})

	t.Run("ContenderWithdrewFromFinals_ContenderNotFound", func(t *testing.T) {
		contender1ID := domain.ContenderID(uuid.New())

		f, awaitExpectations := makeFixture()

		f.store.
			On("GetContender", contender1ID).
			Return(scores.Contender{}, false)

		f.engine.HandleContenderWithdrewFromFinals(domain.ContenderWithdrewFromFinalsEvent{
			ContenderID: contender1ID,
		})

		awaitExpectations(t)
	})

	t.Run("ContenderWithdrewFromFinals", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			contender1ID := domain.ContenderID(uuid.New())
			contender2ID := domain.ContenderID(uuid.New())
			class1ID := domain.CompClassID(uuid.New())

			f, awaitExpectations := makeFixture()

			f.store.
				On("GetContender", contender1ID).
				Return(scores.Contender{
					ID:           contender1ID,
					CompClassID:  class1ID,
					Disqualified: true,
					Score:        123,
				}, true)

			f.store.On("SaveContender", scores.Contender{
				ID:                  contender1ID,
				CompClassID:         class1ID,
				Disqualified:        true,
				WithdrawnFromFinals: true,
				Score:               123,
			}).Return()

			f.store.
				On("GetContendersByCompClass", class1ID).
				Return(slices.Values([]scores.Contender{{ID: contender1ID}, {ID: contender2ID}}))

			f.store.On("SaveScore", domain.Score{Timestamp: time.Now(), ContenderID: contender1ID, Placement: 1, RankOrder: 0}).Return()
			f.store.On("SaveScore", domain.Score{Timestamp: time.Now(), ContenderID: contender2ID, Placement: 1, RankOrder: 1}).Return()

			f.engine.HandleContenderWithdrewFromFinals(domain.ContenderWithdrewFromFinalsEvent{
				ContenderID: contender1ID,
			})

			awaitExpectations(t)
		})
	})

	t.Run("ContenderReenteredFinals_ContenderNotFound", func(t *testing.T) {
		contender1ID := domain.ContenderID(uuid.New())

		f, awaitExpectations := makeFixture()

		f.store.
			On("GetContender", contender1ID).
			Return(scores.Contender{}, false)

		f.engine.HandleContenderReenteredFinals(domain.ContenderReenteredFinalsEvent{
			ContenderID: contender1ID,
		})

		awaitExpectations(t)
	})

	t.Run("ContenderReenteredFinals", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			contender1ID := domain.ContenderID(uuid.New())
			contender2ID := domain.ContenderID(uuid.New())
			class1ID := domain.CompClassID(uuid.New())

			f, awaitExpectations := makeFixture()

			f.store.
				On("GetContender", contender1ID).
				Return(scores.Contender{
					ID:                  contender1ID,
					CompClassID:         class1ID,
					Disqualified:        true,
					WithdrawnFromFinals: true,
					Score:               123,
				}, true)

			f.store.On("SaveContender", scores.Contender{
				ID:                  contender1ID,
				CompClassID:         class1ID,
				Disqualified:        true,
				WithdrawnFromFinals: false,
				Score:               123,
			}).Return()

			f.store.
				On("GetContendersByCompClass", class1ID).
				Return(slices.Values([]scores.Contender{{ID: contender1ID}, {ID: contender2ID}}))

			f.store.On("SaveScore", domain.Score{Timestamp: time.Now(), ContenderID: contender1ID, Placement: 1, RankOrder: 0}).Return()
			f.store.On("SaveScore", domain.Score{Timestamp: time.Now(), ContenderID: contender2ID, Placement: 1, RankOrder: 1}).Return()

			f.engine.HandleContenderReenteredFinals(domain.ContenderReenteredFinalsEvent{
				ContenderID: contender1ID,
			})

			awaitExpectations(t)
		})
	})

	t.Run("ContenderDisqualified_ContenderNotFound", func(t *testing.T) {
		contender1ID := domain.ContenderID(uuid.New())

		f, awaitExpectations := makeFixture()

		f.store.
			On("GetContender", contender1ID).
			Return(scores.Contender{}, false)

		f.engine.HandleContenderDisqualified(domain.ContenderDisqualifiedEvent{
			ContenderID: contender1ID,
		})

		awaitExpectations(t)
	})

	t.Run("ContenderDisqualified", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			contender1ID := domain.ContenderID(uuid.New())
			contender2ID := domain.ContenderID(uuid.New())
			class1ID := domain.CompClassID(uuid.New())

			f, awaitExpectations := makeFixture()

			f.store.
				On("GetContender", contender1ID).
				Return(scores.Contender{
					ID:                  contender1ID,
					CompClassID:         class1ID,
					WithdrawnFromFinals: true,
					Score:               123,
				}, true)

			f.store.On("SaveContender", scores.Contender{
				ID:                  contender1ID,
				CompClassID:         class1ID,
				Disqualified:        true,
				WithdrawnFromFinals: true,
				Score:               0,
			}).Return()

			f.store.
				On("GetContendersByCompClass", class1ID).
				Return(slices.Values([]scores.Contender{{ID: contender1ID}, {ID: contender2ID}}))

			f.store.On("SaveScore", domain.Score{Timestamp: time.Now(), ContenderID: contender1ID, Placement: 1, RankOrder: 0}).Return()
			f.store.On("SaveScore", domain.Score{Timestamp: time.Now(), ContenderID: contender2ID, Placement: 1, RankOrder: 1}).Return()

			f.engine.HandleContenderDisqualified(domain.ContenderDisqualifiedEvent{
				ContenderID: contender1ID,
			})

			awaitExpectations(t)
		})
	})

	t.Run("ContenderRequalified_ContenderNotFound", func(t *testing.T) {
		contender1ID := domain.ContenderID(uuid.New())

		f, awaitExpectations := makeFixture()

		f.store.
			On("GetContender", contender1ID).
			Return(scores.Contender{}, false)

		f.engine.HandleContenderRequalified(domain.ContenderRequalifiedEvent{
			ContenderID: contender1ID,
		})

		awaitExpectations(t)
	})

	t.Run("ContenderRequalified", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			contender1ID := domain.ContenderID(uuid.New())
			contender2ID := domain.ContenderID(uuid.New())
			class1ID := domain.CompClassID(uuid.New())

			f, awaitExpectations := makeFixture()

			f.store.
				On("GetContender", contender1ID).
				Return(scores.Contender{
					ID:                  contender1ID,
					CompClassID:         class1ID,
					Disqualified:        true,
					WithdrawnFromFinals: true,
					Score:               0,
				}, true)

			f.store.
				On("GetTicks", contender1ID).
				Return(slices.Values([]scores.Tick{{Points: 100}, {Points: 200}, {Points: 300}}))

			f.store.On("SaveContender", scores.Contender{
				ID:                  contender1ID,
				CompClassID:         class1ID,
				Disqualified:        false,
				WithdrawnFromFinals: true,
				Score:               600,
			}).Return()

			f.store.
				On("GetContendersByCompClass", class1ID).
				Return(slices.Values([]scores.Contender{{ID: contender1ID}, {ID: contender2ID}}))

			f.store.On("SaveScore", domain.Score{Timestamp: time.Now(), ContenderID: contender1ID, Placement: 1, RankOrder: 0}).Return()
			f.store.On("SaveScore", domain.Score{Timestamp: time.Now(), ContenderID: contender2ID, Placement: 1, RankOrder: 1}).Return()

			f.engine.HandleContenderRequalified(domain.ContenderRequalifiedEvent{
				ContenderID: contender1ID,
			})

			awaitExpectations(t)
		})
	})

	t.Run("ProblemAdded", func(t *testing.T) {
		problem1ID := domain.ProblemID(uuid.New())

		f, awaitExpectations := makeFixture()

		f.store.
			On("SaveProblem", scores.Problem{
				ID:          problem1ID,
				PointsTop:   100,
				PointsZone1: 50,
				PointsZone2: 75,
				FlashBonus:  10,
			}).
			Return()

		f.engine.HandleProblemAdded(domain.ProblemAddedEvent{
			ProblemID:   problem1ID,
			PointsTop:   100,
			PointsZone1: 50,
			PointsZone2: 75,
			FlashBonus:  10,
		})

		awaitExpectations(t)
	})

	t.Run("AscentRegistered_ContenderNotFound", func(t *testing.T) {
		contender1ID := domain.ContenderID(uuid.New())
		problem1ID := domain.ProblemID(uuid.New())

		f, awaitExpectations := makeFixture()

		f.store.
			On("GetContender", contender1ID).
			Return(scores.Contender{}, false)

		f.engine.HandleAscentRegistered(domain.AscentRegisteredEvent{
			ContenderID:   contender1ID,
			ProblemID:     problem1ID,
			Top:           true,
			AttemptsTop:   3,
			Zone1:         true,
			AttemptsZone1: 1,
			Zone2:         true,
			AttemptsZone2: 2,
		})

		awaitExpectations(t)
	})

	t.Run("AscentRegistered_ProblemNotFound", func(t *testing.T) {
		contender1ID := domain.ContenderID(uuid.New())
		class1ID := domain.CompClassID(uuid.New())
		problem1ID := domain.ProblemID(uuid.New())

		f, awaitExpectations := makeFixture()

		f.store.
			On("GetContender", contender1ID).
			Return(scores.Contender{
				ID:          contender1ID,
				CompClassID: class1ID,
			}, true)

		f.store.
			On("GetProblem", problem1ID).
			Return(scores.Problem{}, false)

		f.engine.HandleAscentRegistered(domain.AscentRegisteredEvent{
			ContenderID:   contender1ID,
			ProblemID:     problem1ID,
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
		contender1ID := domain.ContenderID(uuid.New())
		class1ID := domain.CompClassID(uuid.New())
		problem1ID := domain.ProblemID(uuid.New())

		f, awaitExpectations := makeFixture()

		f.store.
			On("GetContender", contender1ID).
			Return(scores.Contender{
				ID:           contender1ID,
				CompClassID:  class1ID,
				Disqualified: true,
			}, true)

		f.store.
			On("GetProblem", problem1ID).
			Return(scores.Problem{
				ID:          problem1ID,
				PointsTop:   100,
				PointsZone1: 50,
				PointsZone2: 75,
				FlashBonus:  10,
			}, true)

		f.store.
			On("SaveTick", contender1ID, scores.Tick{
				ProblemID:     problem1ID,
				Top:           true,
				AttemptsTop:   5,
				Zone1:         true,
				AttemptsZone1: 2,
				Zone2:         true,
				AttemptsZone2: 3,
				Points:        100,
			}).
			Return()

		f.engine.HandleAscentRegistered(domain.AscentRegisteredEvent{
			ContenderID:   contender1ID,
			ProblemID:     problem1ID,
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
		synctest.Test(t, func(t *testing.T) {
			contender1ID := domain.ContenderID(uuid.New())
			contender2ID := domain.ContenderID(uuid.New())
			class1ID := domain.CompClassID(uuid.New())
			problem1ID := domain.ProblemID(uuid.New())

			f, awaitExpectations := makeFixture()

			f.store.
				On("GetContender", contender1ID).
				Return(scores.Contender{
					ID:          contender1ID,
					CompClassID: class1ID,
				}, true)

			f.store.
				On("GetProblem", problem1ID).
				Return(scores.Problem{
					ID:          problem1ID,
					PointsTop:   100,
					PointsZone1: 50,
					PointsZone2: 75,
					FlashBonus:  10,
				}, true)

			f.store.
				On("SaveTick", contender1ID, scores.Tick{
					ProblemID:     problem1ID,
					Top:           true,
					AttemptsTop:   5,
					Zone1:         true,
					AttemptsZone1: 2,
					Zone2:         true,
					AttemptsZone2: 3,
					Points:        100,
				}).
				Return()

			f.store.
				On("GetTicks", contender1ID).
				Return(slices.Values([]scores.Tick{{Points: 100}, {Points: 200}, {Points: 300}}))

			f.store.
				On("SaveContender", scores.Contender{
					ID:          contender1ID,
					CompClassID: class1ID,
					Score:       600,
				}).
				Return()

			f.store.
				On("GetContendersByCompClass", class1ID).
				Return(slices.Values([]scores.Contender{{ID: contender1ID}, {ID: contender2ID}}))

			f.store.On("SaveScore", domain.Score{Timestamp: time.Now(), ContenderID: contender1ID, Placement: 1, RankOrder: 0}).Return()
			f.store.On("SaveScore", domain.Score{Timestamp: time.Now(), ContenderID: contender2ID, Placement: 1, RankOrder: 1}).Return()

			f.engine.HandleAscentRegistered(domain.AscentRegisteredEvent{
				ContenderID:   contender1ID,
				ProblemID:     problem1ID,
				Top:           true,
				AttemptsTop:   5,
				Zone1:         true,
				AttemptsZone1: 2,
				Zone2:         true,
				AttemptsZone2: 3,
			})

			awaitExpectations(t)
		})
	})

	t.Run("AscentDeregistered_ContenderNotFound", func(t *testing.T) {
		contender1ID := domain.ContenderID(uuid.New())
		problem1ID := domain.ProblemID(uuid.New())

		f, awaitExpectations := makeFixture()

		f.store.
			On("GetContender", contender1ID).
			Return(scores.Contender{}, false)

		f.engine.HandleAscentDeregistered(domain.AscentDeregisteredEvent{
			ContenderID: contender1ID,
			ProblemID:   problem1ID,
		})

		awaitExpectations(t)
	})

	t.Run("AscentDeregistered_Disqualified", func(t *testing.T) {
		contender1ID := domain.ContenderID(uuid.New())
		class1ID := domain.CompClassID(uuid.New())
		problem1ID := domain.ProblemID(uuid.New())

		f, awaitExpectations := makeFixture()

		f.store.
			On("GetContender", contender1ID).
			Return(scores.Contender{
				ID:           contender1ID,
				CompClassID:  class1ID,
				Disqualified: true,
			}, true)

		f.store.
			On("DeleteTick", contender1ID, problem1ID).
			Return()

		f.engine.HandleAscentDeregistered(domain.AscentDeregisteredEvent{
			ContenderID: contender1ID,
			ProblemID:   problem1ID,
		})

		awaitExpectations(t)
	})

	t.Run("AscentDeregistered", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			contender1ID := domain.ContenderID(uuid.New())
			contender2ID := domain.ContenderID(uuid.New())
			class1ID := domain.CompClassID(uuid.New())
			problem1ID := domain.ProblemID(uuid.New())

			f, awaitExpectations := makeFixture()

			f.store.
				On("GetContender", contender1ID).
				Return(scores.Contender{
					ID:          contender1ID,
					CompClassID: class1ID,
				}, true)

			f.store.
				On("DeleteTick", contender1ID, problem1ID).
				Return()

			f.store.
				On("GetTicks", contender1ID).
				Return(slices.Values([]scores.Tick{{Points: 100}, {Points: 200}, {Points: 300}}))

			f.store.
				On("SaveContender", scores.Contender{
					ID:          contender1ID,
					CompClassID: class1ID,
					Score:       600,
				}).
				Return()

			f.store.
				On("GetContendersByCompClass", class1ID).
				Return(slices.Values([]scores.Contender{{ID: contender1ID}, {ID: contender2ID}}))

			f.store.On("SaveScore", domain.Score{Timestamp: time.Now(), ContenderID: contender1ID, Placement: 1, RankOrder: 0}).Return()
			f.store.On("SaveScore", domain.Score{Timestamp: time.Now(), ContenderID: contender2ID, Placement: 1, RankOrder: 1}).Return()

			f.engine.HandleAscentDeregistered(domain.AscentDeregisteredEvent{
				ContenderID: contender1ID,
				ProblemID:   problem1ID,
			})

			awaitExpectations(t)
		})
	})

	t.Run("GetDirtyScores", func(t *testing.T) {
		contender1ID := domain.ContenderID(uuid.New())
		contender2ID := domain.ContenderID(uuid.New())
		contender3ID := domain.ContenderID(uuid.New())

		f, awaitExpectations := makeFixture()

		now := time.Now()

		fakedScores := []domain.Score{
			{
				ContenderID: contender1ID,
				Timestamp:   now,
				Score:       100,
				Placement:   1,
				RankOrder:   0,
				Finalist:    true,
			},
			{
				ContenderID: contender2ID,
				Timestamp:   now,
				Score:       200,
				Placement:   2,
				RankOrder:   1,
				Finalist:    true,
			},
			{
				ContenderID: contender3ID,
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

func (m *engineStoreMock) GetTicks(contenderID domain.ContenderID) iter.Seq[scores.Tick] {
	args := m.Called(contenderID)
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

func (m *engineStoreMock) SaveScore(score domain.Score) {
	m.Called(score)
}

func (m *engineStoreMock) GetDirtyScores() []domain.Score {
	args := m.Called()
	return args.Get(0).([]domain.Score)
}
