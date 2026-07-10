package scores_test

import (
	"context"
	"testing"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/scores"
	"github.com/climblive/platform/backend/internal/testutils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHydrate(t *testing.T) {
	mockedRepo := new(repositoryMock)
	mockedStore := new(engineStoreMock)

	fakedContestID := testutils.RandomResourceID[domain.ContestID]()
	fakedProblemID := testutils.RandomResourceID[domain.ProblemID]()
	fakedContenderID := testutils.RandomResourceID[domain.ContenderID]()
	fakedCompClassID := testutils.RandomResourceID[domain.CompClassID]()

	now := time.Now()

	mockedRepo.
		On("GetContest", mock.Anything, nil, fakedContestID).
		Return(domain.Contest{
			ID:                 fakedContestID,
			QualifyingProblems: 10,
			Finalists:          7,
		}, nil)

	mockedRepo.
		On("GetProblemsByContest", mock.Anything, nil, fakedContestID).
		Return([]domain.Problem{
			{
				ID: fakedProblemID,
				ProblemValue: domain.ProblemValue{
					PointsTop:   100,
					PointsZone1: 50,
					PointsZone2: 75,
					FlashBonus:  10,
				},
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
				Entered:             now,
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
				ProblemID:     fakedProblemID,
				Top:           true,
				AttemptsTop:   999,
				Zone1:         true,
				AttemptsZone1: 1,
				Zone2:         true,
				AttemptsZone2: 2,
			},
		}, nil)

	mockedStore.On("SaveRules", scores.Rules{
		QualifyingProblems: 10,
		Finalists:          7,
	}).Return()

	mockedStore.On("SaveProblem", scores.Problem{
		ID: fakedProblemID,
		ProblemValue: domain.ProblemValue{
			PointsTop:   100,
			PointsZone1: 50,
			PointsZone2: 75,
			FlashBonus:  10,
		},
	}).Return()

	mockedStore.On("SaveContender", scores.Contender{
		ID:                  fakedContenderID,
		CompClassID:         fakedCompClassID,
		Disqualified:        true,
		WithdrawnFromFinals: true,
	}).Return()

	mockedStore.On("SaveTick", fakedContenderID, scores.Tick{
		ContenderID:   fakedContenderID,
		ProblemID:     fakedProblemID,
		Top:           true,
		AttemptsTop:   999,
		Zone1:         true,
		AttemptsZone1: 1,
		Zone2:         true,
		AttemptsZone2: 2,
	}).Return()

	hydrator := &scores.StandardEngineStoreHydrator{Repo: mockedRepo}
	err := hydrator.Hydrate(context.Background(), fakedContestID, mockedStore)

	require.NoError(t, err)

	mockedRepo.AssertExpectations(t)
	mockedStore.AssertExpectations(t)
}
