package scores_test

import (
	"iter"
	"slices"
	"testing"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/scores"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDefaultScoreEngine(t *testing.T) {
	type fixture struct {
		rules  *scoringRulesMock
		ranker *rankerMock
		store  *engineStoreMock
		engine *scores.DefaultScoreEngine
	}

	makeFixture := func() (fixture, func(t *testing.T)) {
		mockedRanker := new(rankerMock)
		mockedRules := new(scoringRulesMock)
		mockedStore := new(engineStoreMock)

		engine := scores.NewDefaultScoreEngine(mockedRanker, mockedRules, mockedStore)

		awaitExpectations := func(t *testing.T) {
			mockedRules.AssertExpectations(t)
			mockedRanker.AssertExpectations(t)
			mockedStore.AssertExpectations(t)
		}

		return fixture{
			rules:  mockedRules,
			ranker: mockedRanker,
			store:  mockedStore,
			engine: engine,
		}, awaitExpectations
	}

	t.Run("ReplaceScoringRules", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		f.store.On("GetAllContenders").
			Return(slices.Values([]scores.Contender{
				{ID: 1, CompClassID: 1},
				{ID: 2, CompClassID: 2},
				{ID: 3, CompClassID: 3},
			}))

		f.store.
			On("GetTicks", domain.ContenderID(1)).
			Return(slices.Values([]scores.Tick{{Points: 100}, {Points: 200}, {Points: 300}})).
			On("GetTicks", domain.ContenderID(2)).
			Return(slices.Values([]scores.Tick{{Points: 400}, {Points: 500}})).
			On("GetTicks", domain.ContenderID(3)).
			Return(slices.Values([]scores.Tick{{Points: 600}}))

		f.store.
			On("SaveContender", scores.Contender{ID: 1, CompClassID: 1, Score: 3_000_000}).Return().
			On("SaveContender", scores.Contender{ID: 2, CompClassID: 2, Score: 2_000_000}).Return().
			On("SaveContender", scores.Contender{ID: 3, CompClassID: 3, Score: 1_000_000}).Return()

		f.store.On("GetCompClassIDs").Return([]domain.CompClassID{1, 2, 3})

		f.store.
			On("GetContendersByCompClass", domain.CompClassID(1)).
			Return(slices.Values([]scores.Contender{{ID: 1}})).
			On("GetContendersByCompClass", domain.CompClassID(2)).
			Return(slices.Values([]scores.Contender{{ID: 2}})).
			On("GetContendersByCompClass", domain.CompClassID(3)).
			Return(slices.Values([]scores.Contender{{ID: 3}}))

		f.ranker.
			On("RankContenders", iterMatcher([]scores.Contender{{ID: 1}})).
			Return([]domain.Score{{ContenderID: 1, Placement: 1}}).
			On("RankContenders", iterMatcher([]scores.Contender{{ID: 2}})).
			Return([]domain.Score{{ContenderID: 2, Placement: 2}}).
			On("RankContenders", iterMatcher([]scores.Contender{{ID: 3}})).
			Return([]domain.Score{{ContenderID: 3, Placement: 3}})

		f.store.On("SaveScore", domain.Score{ContenderID: 1, Placement: 1}).Return()
		f.store.On("SaveScore", domain.Score{ContenderID: 2, Placement: 2}).Return()
		f.store.On("SaveScore", domain.Score{ContenderID: 3, Placement: 3}).Return()

		f.engine.ReplaceScoringRules(&jackpotRules{})

		awaitExpectations(t)
	})

	t.Run("ReplaceRanker", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		f.store.On("GetCompClassIDs").Return([]domain.CompClassID{1, 2, 3})

		f.store.
			On("GetContendersByCompClass", domain.CompClassID(1)).
			Return(slices.Values([]scores.Contender{{ID: 1}})).
			On("GetContendersByCompClass", domain.CompClassID(2)).
			Return(slices.Values([]scores.Contender{{ID: 2}})).
			On("GetContendersByCompClass", domain.CompClassID(3)).
			Return(slices.Values([]scores.Contender{{ID: 3}}))

		f.store.On("SaveScore", domain.Score{ContenderID: 1, Placement: 1_000}).Return()
		f.store.On("SaveScore", domain.Score{ContenderID: 2, Placement: 1_000}).Return()
		f.store.On("SaveScore", domain.Score{ContenderID: 3, Placement: 1_000}).Return()

		f.engine.ReplaceRanker(&fakeRanker{})

		awaitExpectations(t)
	})

	t.Run("Start", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		f.store.On("GetAllContenders").
			Return(slices.Values([]scores.Contender{
				{ID: 1, CompClassID: 1},
				{ID: 2, CompClassID: 2},
			}))

		f.store.
			On("GetTicks", domain.ContenderID(1)).
			Return(slices.Values([]scores.Tick{{ProblemID: 1, Top: true}, {ProblemID: 2, Top: true}, {ProblemID: 3, Top: true}})).
			On("GetTicks", domain.ContenderID(2)).
			Return(slices.Values([]scores.Tick{{ProblemID: 1, Zone1: true}, {ProblemID: 2, Zone1: true}}))

		f.store.
			On("GetProblem", domain.ProblemID(1)).
			Return(scores.Problem{ID: 1, PointsTop: 100, PointsZone1: 10}, true).
			On("GetProblem", domain.ProblemID(2)).
			Return(scores.Problem{ID: 2, PointsTop: 200, PointsZone1: 20}, true).
			On("GetProblem", domain.ProblemID(3)).
			Return(scores.Problem{ID: 3, PointsTop: 300}, true)

		f.store.
			On("SaveTick", domain.ContenderID(1), scores.Tick{ProblemID: 1, Top: true, Points: 100}).Return().
			On("SaveTick", domain.ContenderID(1), scores.Tick{ProblemID: 2, Top: true, Points: 200}).Return().
			On("SaveTick", domain.ContenderID(1), scores.Tick{ProblemID: 3, Top: true, Points: 300}).Return().
			On("SaveTick", domain.ContenderID(2), scores.Tick{ProblemID: 1, Zone1: true, Points: 10}).Return().
			On("SaveTick", domain.ContenderID(2), scores.Tick{ProblemID: 2, Zone1: true, Points: 20}).Return()

		f.rules.
			On("CalculateScore", iterMatcher([]int{100, 200, 300})).
			Return(600).
			On("CalculateScore", iterMatcher([]int{10, 20})).
			Return(30)

		f.store.
			On("SaveContender", scores.Contender{ID: 1, CompClassID: 1, Score: 600}).Return().
			On("SaveContender", scores.Contender{ID: 2, CompClassID: 2, Score: 30}).Return()

		f.store.On("GetCompClassIDs").Return([]domain.CompClassID{1, 2, 3})

		f.store.
			On("GetContendersByCompClass", domain.CompClassID(1)).
			Return(slices.Values([]scores.Contender{{ID: 1}})).
			On("GetContendersByCompClass", domain.CompClassID(2)).
			Return(slices.Values([]scores.Contender{{ID: 2}})).
			On("GetContendersByCompClass", domain.CompClassID(3)).
			Return(slices.Values([]scores.Contender{{ID: 3}}))

		f.ranker.
			On("RankContenders", iterMatcher([]scores.Contender{{ID: 1}})).
			Return([]domain.Score{{ContenderID: 1, Placement: 1}}).
			On("RankContenders", iterMatcher([]scores.Contender{{ID: 2}})).
			Return([]domain.Score{{ContenderID: 2, Placement: 2}}).
			On("RankContenders", iterMatcher([]scores.Contender{{ID: 3}})).
			Return([]domain.Score{{ContenderID: 3, Placement: 3}})

		f.store.On("SaveScore", domain.Score{ContenderID: 1, Placement: 1}).Return()
		f.store.On("SaveScore", domain.Score{ContenderID: 2, Placement: 2}).Return()
		f.store.On("SaveScore", domain.Score{ContenderID: 3, Placement: 3}).Return()

		f.engine.Start()

		awaitExpectations(t)
	})

	t.Run("ContenderEntered", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		f.store.On("SaveContender", scores.Contender{
			ID:          1,
			CompClassID: 1,
		}).Return()

		f.store.
			On("GetContendersByCompClass", domain.CompClassID(1)).
			Return(slices.Values([]scores.Contender{{ID: 1}, {ID: 2}, {ID: 3}}))

		f.ranker.
			On("RankContenders", iterMatcher([]scores.Contender{{ID: 1}, {ID: 2}, {ID: 3}})).
			Return([]domain.Score{{ContenderID: 1, Placement: 1}, {ContenderID: 2, Placement: 2}, {ContenderID: 3, Placement: 3}})

		f.store.On("SaveScore", domain.Score{ContenderID: 1, Placement: 1}).Return()
		f.store.On("SaveScore", domain.Score{ContenderID: 2, Placement: 2}).Return()
		f.store.On("SaveScore", domain.Score{ContenderID: 3, Placement: 3}).Return()

		f.engine.HandleContenderEntered(domain.ContenderEnteredEvent{
			ContenderID: 1,
			CompClassID: 1,
		})

		awaitExpectations(t)
	})

	t.Run("ContenderSwitchedClass_ContenderNotFound", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		f.store.
			On("GetContender", domain.ContenderID(1)).
			Return(scores.Contender{}, false)

		f.engine.HandleContenderSwitchedClass(
			domain.ContenderSwitchedClassEvent{
				ContenderID: 1,
				CompClassID: 1,
			})

		awaitExpectations(t)
	})

	t.Run("ContenderSwitchedClass_SameClass", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		f.store.
			On("GetContender", domain.ContenderID(1)).
			Return(scores.Contender{
				ID:          1,
				CompClassID: 1,
			}, true)

		f.engine.HandleContenderSwitchedClass(domain.ContenderSwitchedClassEvent{
			ContenderID: 1,
			CompClassID: 1,
		})

		awaitExpectations(t)
	})

	t.Run("ContenderSwitchedClass", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		f.store.
			On("GetContender", domain.ContenderID(4)).
			Return(scores.Contender{
				ID:                  4,
				CompClassID:         1,
				Disqualified:        false,
				WithdrawnFromFinals: false,
				Score:               123,
			}, true)

		f.store.On("SaveContender", scores.Contender{
			ID:                  4,
			CompClassID:         2,
			Disqualified:        false,
			WithdrawnFromFinals: false,
			Score:               123,
		}).Return()

		f.store.
			On("GetContendersByCompClass", domain.CompClassID(1)).
			Return(slices.Values([]scores.Contender{{ID: 1}, {ID: 2}, {ID: 3}})).
			On("GetContendersByCompClass", domain.CompClassID(2)).
			Return(slices.Values([]scores.Contender{{ID: 4}}))

		f.ranker.
			On("RankContenders", iterMatcher([]scores.Contender{{ID: 1}, {ID: 2}, {ID: 3}})).
			Return([]domain.Score{{ContenderID: 1, Placement: 1}, {ContenderID: 2, Placement: 2}, {ContenderID: 3, Placement: 3}}).
			On("RankContenders", iterMatcher([]scores.Contender{{ID: 4}})).
			Return([]domain.Score{{ContenderID: 4, Placement: 4}})

		f.store.On("SaveScore", domain.Score{ContenderID: 1, Placement: 1}).Return()
		f.store.On("SaveScore", domain.Score{ContenderID: 2, Placement: 2}).Return()
		f.store.On("SaveScore", domain.Score{ContenderID: 3, Placement: 3}).Return()
		f.store.On("SaveScore", domain.Score{ContenderID: 4, Placement: 4}).Return()

		f.engine.HandleContenderSwitchedClass(domain.ContenderSwitchedClassEvent{
			ContenderID: 4,
			CompClassID: 2,
		})

		awaitExpectations(t)
	})

	t.Run("ContenderWithdrewFromFinals_ContenderNotFound", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		f.store.
			On("GetContender", domain.ContenderID(1)).
			Return(scores.Contender{}, false)

		f.engine.HandleContenderWithdrewFromFinals(domain.ContenderWithdrewFromFinalsEvent{
			ContenderID: 1,
		})

		awaitExpectations(t)
	})

	t.Run("ContenderWithdrewFromFinals", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		f.store.
			On("GetContender", domain.ContenderID(1)).
			Return(scores.Contender{
				ID:           1,
				CompClassID:  1,
				Disqualified: true,
				Score:        123,
			}, true)

		f.store.On("SaveContender", scores.Contender{
			ID:                  1,
			CompClassID:         1,
			Disqualified:        true,
			WithdrawnFromFinals: true,
			Score:               123,
		}).Return()

		f.store.
			On("GetContendersByCompClass", domain.CompClassID(1)).
			Return(slices.Values([]scores.Contender{{ID: 1}, {ID: 2}}))

		f.ranker.
			On("RankContenders", iterMatcher([]scores.Contender{{ID: 1}, {ID: 2}})).
			Return([]domain.Score{{ContenderID: 1, Placement: 1}, {ContenderID: 2, Placement: 2}})

		f.store.On("SaveScore", domain.Score{ContenderID: 1, Placement: 1}).Return()
		f.store.On("SaveScore", domain.Score{ContenderID: 2, Placement: 2}).Return()

		f.engine.HandleContenderWithdrewFromFinals(domain.ContenderWithdrewFromFinalsEvent{
			ContenderID: 1,
		})

		awaitExpectations(t)
	})

	t.Run("ContenderReenteredFinals_ContenderNotFound", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		f.store.
			On("GetContender", domain.ContenderID(1)).
			Return(scores.Contender{}, false)

		f.engine.HandleContenderReenteredFinals(domain.ContenderReenteredFinalsEvent{
			ContenderID: 1,
		})

		awaitExpectations(t)
	})

	t.Run("ContenderReenteredFinals", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		f.store.
			On("GetContender", domain.ContenderID(1)).
			Return(scores.Contender{
				ID:                  1,
				CompClassID:         1,
				Disqualified:        true,
				WithdrawnFromFinals: true,
				Score:               123,
			}, true)

		f.store.On("SaveContender", scores.Contender{
			ID:                  1,
			CompClassID:         1,
			Disqualified:        true,
			WithdrawnFromFinals: false,
			Score:               123,
		}).Return()

		f.store.
			On("GetContendersByCompClass", domain.CompClassID(1)).
			Return(slices.Values([]scores.Contender{{ID: 1}, {ID: 2}}))

		f.ranker.
			On("RankContenders", iterMatcher([]scores.Contender{{ID: 1}, {ID: 2}})).
			Return([]domain.Score{{ContenderID: 1, Placement: 1}, {ContenderID: 2, Placement: 2}})

		f.store.On("SaveScore", domain.Score{ContenderID: 1, Placement: 1}).Return()
		f.store.On("SaveScore", domain.Score{ContenderID: 2, Placement: 2}).Return()

		f.engine.HandleContenderReenteredFinals(domain.ContenderReenteredFinalsEvent{
			ContenderID: 1,
		})

		awaitExpectations(t)
	})

	t.Run("ContenderDisqualified_ContenderNotFound", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		f.store.
			On("GetContender", domain.ContenderID(1)).
			Return(scores.Contender{}, false)

		f.engine.HandleContenderDisqualified(domain.ContenderDisqualifiedEvent{
			ContenderID: 1,
		})

		awaitExpectations(t)
	})

	t.Run("ContenderDisqualified", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		f.store.
			On("GetContender", domain.ContenderID(1)).
			Return(scores.Contender{
				ID:                  1,
				CompClassID:         1,
				WithdrawnFromFinals: true,
				Score:               123,
			}, true)

		f.store.On("SaveContender", scores.Contender{
			ID:                  1,
			CompClassID:         1,
			Disqualified:        true,
			WithdrawnFromFinals: true,
			Score:               0,
		}).Return()

		f.store.
			On("GetContendersByCompClass", domain.CompClassID(1)).
			Return(slices.Values([]scores.Contender{{ID: 1}, {ID: 2}}))

		f.ranker.
			On("RankContenders", iterMatcher([]scores.Contender{{ID: 1}, {ID: 2}})).
			Return([]domain.Score{{ContenderID: 1, Placement: 1}, {ContenderID: 2, Placement: 2}})

		f.store.On("SaveScore", domain.Score{ContenderID: 1, Placement: 1}).Return()
		f.store.On("SaveScore", domain.Score{ContenderID: 2, Placement: 2}).Return()

		f.engine.HandleContenderDisqualified(domain.ContenderDisqualifiedEvent{
			ContenderID: 1,
		})

		awaitExpectations(t)
	})

	t.Run("ContenderRequalified_ContenderNotFound", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		f.store.
			On("GetContender", domain.ContenderID(1)).
			Return(scores.Contender{}, false)

		f.engine.HandleContenderRequalified(domain.ContenderRequalifiedEvent{
			ContenderID: 1,
		})

		awaitExpectations(t)
	})

	t.Run("ContenderRequalified", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		f.store.
			On("GetContender", domain.ContenderID(1)).
			Return(scores.Contender{
				ID:                  1,
				CompClassID:         1,
				Disqualified:        true,
				WithdrawnFromFinals: true,
				Score:               0,
			}, true)

		f.store.
			On("GetTicks", domain.ContenderID(1)).
			Return(slices.Values([]scores.Tick{{Points: 100}, {Points: 200}, {Points: 300}}))

		f.rules.
			On("CalculateScore", iterMatcher([]int{100, 200, 300})).
			Return(123)

		f.store.On("SaveContender", scores.Contender{
			ID:                  1,
			CompClassID:         1,
			Disqualified:        false,
			WithdrawnFromFinals: true,
			Score:               123,
		}).Return()

		f.store.
			On("GetContendersByCompClass", domain.CompClassID(1)).
			Return(slices.Values([]scores.Contender{{ID: 1}, {ID: 2}}))

		f.ranker.
			On("RankContenders", iterMatcher([]scores.Contender{{ID: 1}, {ID: 2}})).
			Return([]domain.Score{{ContenderID: 1, Placement: 1}, {ContenderID: 2, Placement: 2}})

		f.store.On("SaveScore", domain.Score{ContenderID: 1, Placement: 1}).Return()
		f.store.On("SaveScore", domain.Score{ContenderID: 2, Placement: 2}).Return()

		f.engine.HandleContenderRequalified(domain.ContenderRequalifiedEvent{
			ContenderID: 1,
		})

		awaitExpectations(t)
	})

	t.Run("ProblemAdded", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		f.store.
			On("SaveProblem", scores.Problem{
				ID:          1,
				PointsTop:   100,
				PointsZone1: 50,
				PointsZone2: 75,
				FlashBonus:  10,
			}).
			Return()

		f.engine.HandleProblemAdded(domain.ProblemAddedEvent{
			ProblemID:   1,
			PointsTop:   100,
			PointsZone1: 50,
			PointsZone2: 75,
			FlashBonus:  10,
		})

		awaitExpectations(t)
	})

	t.Run("AscentRegistered_ContenderNotFound", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		f.store.
			On("GetContender", domain.ContenderID(1)).
			Return(scores.Contender{}, false)

		f.engine.HandleAscentRegistered(domain.AscentRegisteredEvent{
			ContenderID:   1,
			ProblemID:     1,
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
		f, awaitExpectations := makeFixture()

		f.store.
			On("GetContender", domain.ContenderID(1)).
			Return(scores.Contender{
				ID:          1,
				CompClassID: 1,
			}, true)

		f.store.
			On("GetProblem", domain.ProblemID(1)).
			Return(scores.Problem{}, false)

		f.engine.HandleAscentRegistered(domain.AscentRegisteredEvent{
			ContenderID:   1,
			ProblemID:     1,
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

		f.store.
			On("GetContender", domain.ContenderID(1)).
			Return(scores.Contender{
				ID:           1,
				CompClassID:  1,
				Disqualified: true,
			}, true)

		f.store.
			On("GetProblem", domain.ProblemID(1)).
			Return(scores.Problem{
				ID:          1,
				PointsTop:   100,
				PointsZone1: 50,
				PointsZone2: 75,
				FlashBonus:  10,
			}, true)

		f.store.
			On("SaveTick", domain.ContenderID(1), scores.Tick{
				ProblemID:     1,
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
			ContenderID:   1,
			ProblemID:     1,
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

		f.store.
			On("GetContender", domain.ContenderID(1)).
			Return(scores.Contender{
				ID:          1,
				CompClassID: 1,
			}, true)

		f.store.
			On("GetProblem", domain.ProblemID(1)).
			Return(scores.Problem{
				ID:          1,
				PointsTop:   100,
				PointsZone1: 50,
				PointsZone2: 75,
				FlashBonus:  10,
			}, true)

		f.store.
			On("SaveTick", domain.ContenderID(1), scores.Tick{
				ProblemID:     1,
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
			On("GetTicks", domain.ContenderID(1)).
			Return(slices.Values([]scores.Tick{{Points: 100}, {Points: 200}, {Points: 300}}))

		f.rules.
			On("CalculateScore", iterMatcher([]int{100, 200, 300})).
			Return(123)

		f.store.
			On("SaveContender", scores.Contender{
				ID:          1,
				CompClassID: 1,
				Score:       123,
			}).
			Return()

		f.store.
			On("GetContendersByCompClass", domain.CompClassID(1)).
			Return(slices.Values([]scores.Contender{{ID: 1}, {ID: 2}}))

		f.ranker.
			On("RankContenders", iterMatcher([]scores.Contender{{ID: 1}, {ID: 2}})).
			Return([]domain.Score{{ContenderID: 1, Placement: 1}, {ContenderID: 2, Placement: 2}})

		f.store.On("SaveScore", domain.Score{ContenderID: 1, Placement: 1}).Return()
		f.store.On("SaveScore", domain.Score{ContenderID: 2, Placement: 2}).Return()

		f.engine.HandleAscentRegistered(domain.AscentRegisteredEvent{
			ContenderID:   1,
			ProblemID:     1,
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

		f.store.
			On("GetContender", domain.ContenderID(1)).
			Return(scores.Contender{}, false)

		f.engine.HandleAscentDeregistered(domain.AscentDeregisteredEvent{
			ContenderID: 1,
			ProblemID:   1,
		})

		awaitExpectations(t)
	})

	t.Run("AscentDeregistered_Disqualified", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		f.store.
			On("GetContender", domain.ContenderID(1)).
			Return(scores.Contender{
				ID:           1,
				CompClassID:  1,
				Disqualified: true,
			}, true)

		f.store.
			On("DeleteTick", domain.ContenderID(1), domain.ProblemID(1)).
			Return()

		f.engine.HandleAscentDeregistered(domain.AscentDeregisteredEvent{
			ContenderID: 1,
			ProblemID:   1,
		})

		awaitExpectations(t)
	})

	t.Run("AscentDeregistered", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		f.store.
			On("GetContender", domain.ContenderID(1)).
			Return(scores.Contender{
				ID:          1,
				CompClassID: 1,
			}, true)

		f.store.
			On("DeleteTick", domain.ContenderID(1), domain.ProblemID(1)).
			Return()

		f.store.
			On("GetTicks", domain.ContenderID(1)).
			Return(slices.Values([]scores.Tick{{Points: 100}, {Points: 200}, {Points: 300}}))

		f.rules.
			On("CalculateScore", iterMatcher([]int{100, 200, 300})).
			Return(123)

		f.store.
			On("SaveContender", scores.Contender{
				ID:          1,
				CompClassID: 1,
				Score:       123,
			}).
			Return()

		f.store.
			On("GetContendersByCompClass", domain.CompClassID(1)).
			Return(slices.Values([]scores.Contender{{ID: 1}, {ID: 2}}))

		f.ranker.
			On("RankContenders", iterMatcher([]scores.Contender{{ID: 1}, {ID: 2}})).
			Return([]domain.Score{{ContenderID: 1, Placement: 1}, {ContenderID: 2, Placement: 2}})

		f.store.On("SaveScore", domain.Score{ContenderID: 1, Placement: 1}).Return()
		f.store.On("SaveScore", domain.Score{ContenderID: 2, Placement: 2}).Return()

		f.engine.HandleAscentDeregistered(domain.AscentDeregisteredEvent{
			ContenderID: 1,
			ProblemID:   1,
		})

		awaitExpectations(t)
	})

	t.Run("GetDirtyScores", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		now := time.Now()

		fakedScores := []domain.Score{
			{
				ContenderID: 1,
				Timestamp:   now,
				Score:       100,
				Placement:   1,
				RankOrder:   0,
				Finalist:    true,
			},
			{
				ContenderID: 2,
				Timestamp:   now,
				Score:       200,
				Placement:   2,
				RankOrder:   1,
				Finalist:    true,
			},
			{
				ContenderID: 3,
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

func iterMatcher[T comparable](expected []T) any {
	return mock.MatchedBy(func(values iter.Seq[T]) bool {
		return assert.ObjectsAreEqual(expected, slices.Collect(values))
	})
}

type rankerMock struct {
	mock.Mock
}

func (m *rankerMock) RankContenders(contenders iter.Seq[scores.Contender]) []domain.Score {
	args := m.Called(contenders)
	return args.Get(0).([]domain.Score)
}

type scoringRulesMock struct {
	mock.Mock
}

func (m *scoringRulesMock) CalculateScore(points iter.Seq[int]) int {
	args := m.Called(points)
	return args.Get(0).(int)
}

type engineStoreMock struct {
	mock.Mock
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

type jackpotRules struct{}

func (m *jackpotRules) CalculateScore(points iter.Seq[int]) int {
	return len(slices.Collect(points)) * 1_000_000
}

type fakeRanker struct{}

func (r *fakeRanker) RankContenders(contenders iter.Seq[scores.Contender]) []domain.Score {
	var scores []domain.Score

	for contender := range contenders {
		scores = append(scores, domain.Score{
			ContenderID: contender.ID,
			Placement:   1_000,
		})
	}

	return scores
}
