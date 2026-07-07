package scores_test

import (
	"iter"
	"testing"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/scores"
	"github.com/climblive/platform/backend/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type effectResolverMock struct {
	mock.Mock
}

func (m *effectResolverMock) RankCompClass(compClassID domain.CompClassID) {
	m.Called(compClassID)
}

func (m *effectResolverMock) ScoreContender(contenderID domain.ContenderID) iter.Seq[scores.Effect] {
	args := m.Called(contenderID)
	return args.Get(0).(iter.Seq[scores.Effect])
}

func (m *effectResolverMock) CalculatePointValues(compClassID domain.CompClassID, problemID domain.ProblemID) iter.Seq[scores.Effect] {
	args := m.Called(compClassID, problemID)
	return args.Get(0).(iter.Seq[scores.Effect])
}

func TestEffectRunner_RunEffects(t *testing.T) {
	t.Run("PreventUncontrolledChainReaction", func(t *testing.T) {
		mockedResolver := new(effectResolverMock)

		effectYielder := func(iterations int) iter.Seq[scores.Effect] {
			return func(yield func(scores.Effect) bool) {
				for range iterations {
					if !yield(scores.EffectCalculatePointValues{CompClassID: testutils.RandomResourceID[domain.CompClassID](), ProblemID: testutils.RandomResourceID[domain.ProblemID]()}) {
						return
					}

					if !yield(scores.EffectScoreContender{ContenderID: testutils.RandomResourceID[domain.ContenderID]()}) {
						return
					}

					if !yield(scores.EffectRankClass{CompClassID: testutils.RandomResourceID[domain.CompClassID]()}) {
						return
					}
				}
			}
		}

		effectCounters := map[scores.EffectType]int{
			scores.EffectTypeRankClass:            0,
			scores.EffectTypeScoreContender:       0,
			scores.EffectTypeCalculatePointValues: 0,
		}

		incrementCounter := func(effectType scores.EffectType) func(args mock.Arguments) {
			return func(args mock.Arguments) {
				effectCounters[effectType] += 1
			}
		}

		mockedResolver.
			On("RankCompClass", mock.Anything).
			Return().
			Run(incrementCounter(scores.EffectTypeRankClass))
		mockedResolver.
			On("ScoreContender", mock.Anything).
			Return(effectYielder(10)).
			Run(incrementCounter(scores.EffectTypeScoreContender))
		mockedResolver.
			On("CalculatePointValues", mock.Anything, mock.Anything).
			Return(effectYielder(10)).
			Run(incrementCounter(scores.EffectTypeCalculatePointValues))

		runner := scores.NewEffectRunner(mockedResolver)

		runner.RunEffects(effectYielder(10))

		assert.Equal(t, 10, effectCounters[scores.EffectTypeCalculatePointValues])
		assert.Equal(t, 110, effectCounters[scores.EffectTypeScoreContender])
		assert.Equal(t, 1210, effectCounters[scores.EffectTypeRankClass])

		mockedResolver.AssertExpectations(t)
	})

	t.Run("RunEffectsInOrder", func(t *testing.T) {
		mockedResolver := new(effectResolverMock)

		effectYielder := func(effects ...scores.Effect) iter.Seq[scores.Effect] {
			return func(yield func(scores.Effect) bool) {
				for _, effect := range effects {
					if !yield(effect) {
						return
					}
				}
			}
		}

		eff1 := scores.EffectCalculatePointValues{CompClassID: testutils.RandomResourceID[domain.CompClassID](), ProblemID: testutils.RandomResourceID[domain.ProblemID]()}
		eff2 := scores.EffectScoreContender{ContenderID: testutils.RandomResourceID[domain.ContenderID]()}
		eff3 := scores.EffectScoreContender{ContenderID: testutils.RandomResourceID[domain.ContenderID]()}
		eff4 := scores.EffectRankClass{CompClassID: testutils.RandomResourceID[domain.CompClassID]()}
		eff5 := scores.EffectRankClass{CompClassID: testutils.RandomResourceID[domain.CompClassID]()}

		recordedOrder := make([]scores.Effect, 0)

		recordOrder := func(eff scores.Effect) func(args mock.Arguments) {
			return func(mock.Arguments) {
				recordedOrder = append(recordedOrder, eff)
			}
		}

		mockedResolver.
			On("CalculatePointValues", eff1.CompClassID, eff1.ProblemID).
			Return(effectYielder(eff3)).
			Run(recordOrder(eff1))

		mockedResolver.
			On("ScoreContender", eff2.ContenderID).
			Return(effectYielder(eff5)).
			Run(recordOrder(eff2))

		mockedResolver.
			On("ScoreContender", eff3.ContenderID).
			Return(effectYielder()).
			Run(recordOrder(eff3))

		mockedResolver.
			On("RankCompClass", eff4.CompClassID).
			Return().
			Run(recordOrder(eff4))

		mockedResolver.
			On("RankCompClass", eff5.CompClassID).
			Return().
			Run(recordOrder(eff5))

		runner := scores.NewEffectRunner(mockedResolver)

		runner.RunEffects(effectYielder(eff4, eff2, eff1))

		assert.Equal(t, []scores.Effect{eff1, eff2, eff3, eff4, eff5}, recordedOrder)

		mockedResolver.AssertExpectations(t)
	})

	t.Run("DontRunDuplicateChainReactions", func(t *testing.T) {
		mockedResolver := new(effectResolverMock)

		effectYielder := func(effects ...scores.Effect) iter.Seq[scores.Effect] {
			return func(yield func(scores.Effect) bool) {
				for _, effect := range effects {
					if !yield(effect) {
						return
					}
				}
			}
		}

		eff1 := scores.EffectCalculatePointValues{CompClassID: testutils.RandomResourceID[domain.CompClassID](), ProblemID: testutils.RandomResourceID[domain.ProblemID]()}
		eff2 := scores.EffectScoreContender{ContenderID: testutils.RandomResourceID[domain.ContenderID]()}
		eff3 := scores.EffectRankClass{CompClassID: testutils.RandomResourceID[domain.CompClassID]()}

		mockedResolver.
			On("CalculatePointValues", eff1.CompClassID, eff1.ProblemID).
			Return(effectYielder(eff2, eff2, eff2, eff2, eff2))

		mockedResolver.
			On("ScoreContender", eff2.ContenderID).
			Return(effectYielder(eff3, eff3, eff3, eff3, eff3)).
			Once()

		mockedResolver.
			On("RankCompClass", eff3.CompClassID).
			Return().
			Once()

		runner := scores.NewEffectRunner(mockedResolver)

		runner.RunEffects(effectYielder(eff1))

		mockedResolver.AssertExpectations(t)
	})
}
