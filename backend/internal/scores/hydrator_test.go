package scores_test

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/scores"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHydrate(t *testing.T) {
	mockedRepo := new(repositoryMock)
	mockedStore := new(engineStoreMock)

	fakedContestID := domain.ContestID(rand.Int())
	fakedProblemID := domain.ProblemID(rand.Int())
	fakedContenderID := domain.ContenderID(rand.Int())
	fakedCompClassID := domain.CompClassID(rand.Int())

	now := time.Now()

	mockedRepo.
		On("GetProblemsByContest", mock.Anything, nil, fakedContestID).
		Return([]domain.Problem{
			{
				ID:         fakedProblemID,
				PointsTop:  100,
				PointsZone: 50,
				FlashBonus: 10,
			},
		}, nil)

	mockedRepo.
		On("GetContendersByContest", mock.Anything, nil, fakedContestID).
		Return([]domain.Contender{
			{
				ID:                  fakedContenderID,
				CompClassID:         fakedCompClassID,
				Disqualified:        true,
				WithdrawnFromFinals: true,
				Entered:             &now,
			},
			{
				ID:          fakedContenderID + 1,
				CompClassID: 0,
			},
		}, nil)

	mockedRepo.
		On("GetTicksByContest", mock.Anything, nil, fakedContestID).
		Return([]domain.Tick{
			{
				Ownership: domain.OwnershipData{
					ContenderID: &fakedContenderID,
				},
				ProblemID:    fakedProblemID,
				Top:          true,
				AttemptsTop:  999,
				Zone:         true,
				AttemptsZone: 1,
			},
		}, nil)

	mockedStore.On("SaveProblem", scores.Problem{
		ID:         fakedProblemID,
		PointsTop:  100,
		PointsZone: 50,
		FlashBonus: 10,
	}).Return()

	mockedStore.On("SaveContender", scores.Contender{
		ID:                  fakedContenderID,
		CompClassID:         fakedCompClassID,
		Disqualified:        true,
		WithdrawnFromFinals: true,
	}).Return()

	mockedStore.On("SaveTick", fakedContenderID, scores.Tick{
		ProblemID:    fakedProblemID,
		Top:          true,
		AttemptsTop:  999,
		Zone:         true,
		AttemptsZone: 1,
	}).Return()

	hydrator := &scores.StandardEngineStoreHydrator{Repo: mockedRepo}
	err := hydrator.Hydrate(context.Background(), fakedContestID, mockedStore)

	require.NoError(t, err)

	mockedRepo.AssertExpectations(t)
	mockedStore.AssertExpectations(t)
}
