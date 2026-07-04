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

		mockedResolver.On("RankCompClass", mock.Anything).Return()
		mockedResolver.On("ScoreContender", mock.Anything).Return(effectYielder(10))
		mockedResolver.On("CalculatePointValues", mock.Anything, mock.Anything).Return(effectYielder(10))

		runner := scores.NewEffectRunner(mockedResolver)

		effectsResolved := runner.RunEffects(effectYielder(10))

		assert.Equal(t, int64(10+110+1210), effectsResolved)
	})
}
