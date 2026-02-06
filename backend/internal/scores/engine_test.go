package scores_test

import (
	"iter"
	"slices"
	"testing"
	"testing/synctest"
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
		synctest.Test(t, func(t *testing.T) {
			f, awaitExpectations := makeFixture()

			f.store.On("SaveRules", scores.Rules{
				QualifyingProblems: 1,
				Finalists:          0,
			}).Return()

			f.engine.HandleRulesUpdated(domain.RulesUpdatedEvent{
				QualifyingProblems: 1,
				Finalists:          0,
			})

			awaitExpectations(t)
		})
	})

	t.Run("Start", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			f, awaitExpectations := makeFixture()

			f.engine.Start()

			awaitExpectations(t)
		})
	})

	t.Run("ContenderEntered", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
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
	})

	t.Run("ContenderSwitchedClass_ContenderNotFound", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedContenderID := testutils.RandomResourceID[domain.ContenderID]()
		fakedCompClassID := testutils.RandomResourceID[domain.CompClassID]()

		f.store.
			On("GetContender", fakedContenderID).
			Return(scores.Contender{}, false)

		f.engine.HandleContenderSwitchedClass(
			domain.ContenderSwitchedClassEvent{
				ContenderID: fakedContenderID,
				CompClassID: fakedCompClassID,
			})

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

		f.engine.HandleContenderSwitchedClass(domain.ContenderSwitchedClassEvent{
			ContenderID: fakedContenderID,
			CompClassID: fakedCompClassID,
		})

		awaitExpectations(t)
	})

	t.Run("ContenderSwitchedClass", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			f, awaitExpectations := makeFixture()

			fakedContenderID := testutils.RandomResourceID[domain.ContenderID]()
			fakedOldCompClassID := testutils.RandomResourceID[domain.CompClassID]()
			fakedNewCompClassID := testutils.RandomResourceID[domain.CompClassID]()

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

			f.engine.HandleContenderSwitchedClass(domain.ContenderSwitchedClassEvent{
				ContenderID: fakedContenderID,
				CompClassID: fakedNewCompClassID,
			})

			awaitExpectations(t)
		})
	})

	t.Run("ContenderWithdrewFromFinals_ContenderNotFound", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedContenderID := testutils.RandomResourceID[domain.ContenderID]()

		f.store.
			On("GetContender", fakedContenderID).
			Return(scores.Contender{}, false)

		f.engine.HandleContenderWithdrewFromFinals(domain.ContenderWithdrewFromFinalsEvent{
			ContenderID: fakedContenderID,
		})

		awaitExpectations(t)
	})

	t.Run("ContenderWithdrewFromFinals", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
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

			f.engine.HandleContenderWithdrewFromFinals(domain.ContenderWithdrewFromFinalsEvent{
				ContenderID: fakedContenderID,
			})

			awaitExpectations(t)
		})
	})

	t.Run("ContenderReenteredFinals_ContenderNotFound", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedContenderID := testutils.RandomResourceID[domain.ContenderID]()

		f.store.
			On("GetContender", fakedContenderID).
			Return(scores.Contender{}, false)

		f.engine.HandleContenderReenteredFinals(domain.ContenderReenteredFinalsEvent{
			ContenderID: fakedContenderID,
		})

		awaitExpectations(t)
	})

	t.Run("ContenderReenteredFinals", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
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

			f.engine.HandleContenderReenteredFinals(domain.ContenderReenteredFinalsEvent{
				ContenderID: fakedContenderID,
			})

			awaitExpectations(t)
		})
	})

	t.Run("ContenderDisqualified_ContenderNotFound", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedContenderID := testutils.RandomResourceID[domain.ContenderID]()

		f.store.
			On("GetContender", fakedContenderID).
			Return(scores.Contender{}, false)

		f.engine.HandleContenderDisqualified(domain.ContenderDisqualifiedEvent{
			ContenderID: fakedContenderID,
		})

		awaitExpectations(t)
	})

	t.Run("ContenderDisqualified", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			f, awaitExpectations := makeFixture()

			fakedContenderID := testutils.RandomResourceID[domain.ContenderID]()
			fakedCompClassID := testutils.RandomResourceID[domain.CompClassID]()

			f.store.
				On("GetContender", fakedContenderID).
				Return(scores.Contender{
					ID:                  fakedContenderID,
					CompClassID:         fakedCompClassID,
					WithdrawnFromFinals: true,
					Score:               123,
				}, true)

			f.store.On("SaveContender", scores.Contender{
				ID:                  fakedContenderID,
				CompClassID:         fakedCompClassID,
				Disqualified:        true,
				WithdrawnFromFinals: true,
				Score:               123,
			}).Return()

			f.engine.HandleContenderDisqualified(domain.ContenderDisqualifiedEvent{
				ContenderID: fakedContenderID,
			})

			awaitExpectations(t)
		})
	})

	t.Run("ContenderRequalified_ContenderNotFound", func(t *testing.T) {
		f, awaitExpectations := makeFixture()

		fakedContenderID := testutils.RandomResourceID[domain.ContenderID]()

		f.store.
			On("GetContender", fakedContenderID).
			Return(scores.Contender{}, false)

		f.engine.HandleContenderRequalified(domain.ContenderRequalifiedEvent{
			ContenderID: fakedContenderID,
		})

		awaitExpectations(t)
	})

	t.Run("ContenderRequalified", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
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
					Score:               0,
				}, true)

			f.store.On("SaveContender", scores.Contender{
				ID:                  fakedContenderID,
				CompClassID:         fakedCompClassID,
				Disqualified:        false,
				WithdrawnFromFinals: true,
				Score:               0,
			}).Return()

			f.engine.HandleContenderRequalified(domain.ContenderRequalifiedEvent{
				ContenderID: fakedContenderID,
			})

			awaitExpectations(t)
		})
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
		synctest.Test(t, func(t *testing.T) {
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
		synctest.Test(t, func(t *testing.T) {
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
