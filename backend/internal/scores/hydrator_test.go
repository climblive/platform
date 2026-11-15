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

	fakedTickID := testutils.RandomResourceID[domain.TickID]()
	fakedCompetitionID := testutils.RandomResourceID[domain.CompetitionID]()
	fakedProblemID := testutils.RandomResourceID[domain.ProblemID]()
	fakedContenderID := testutils.RandomResourceID[domain.ContenderID]()
	fakedCategoryID := testutils.RandomResourceID[domain.CategoryID]()

	now := time.Now()

	mockedRepo.
		On("GetCompetition", mock.Anything, nil, fakedCompetitionID).
		Return(domain.Competition{
			ID:                 fakedCompetitionID,
			QualifyingProblems: 10,
			Finalists:          7,
			UsePoints:          true,
			PooledPoints:       true,
			MaxAttempts:        10,
			PointDeduction:     3,
		}, nil)

	mockedRepo.
		On("GetProblemsByCompetition", mock.Anything, nil, fakedCompetitionID).
		Return([]domain.Problem{
			{
				ID:           fakedProblemID,
				Zone1Enabled: true,
				Zone2Enabled: true,
				ProblemValue: domain.ProblemValue{
					PointsTop:   100,
					PointsZone1: 50,
					PointsZone2: 75,
					FlashBonus:  10,
				},
			},
		}, nil)

	mockedRepo.
		On("GetContendersByCompetition", mock.Anything, nil, fakedCompetitionID).
		Return([]domain.Contender{
			{
				ID:                  fakedContenderID,
				CategoryID:          fakedCategoryID,
				Disqualified:        true,
				WithdrawnFromFinals: true,
				Entered:             now,
			},
			{
				ID:         fakedContenderID + 1,
				CategoryID: 0,
			},
		}, nil)

	mockedRepo.
		On("GetTicksByCompetition", mock.Anything, nil, fakedCompetitionID).
		Return([]domain.Tick{
			{
				ID: fakedTickID,
				Ownership: domain.OwnershipData{
					ContenderID: &fakedContenderID,
				},
				Revision:      1,
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
		UsePoints:          true,
		PooledPoints:       true,
		MaxAttempts:        10,
		PointDeduction:     3,
	}).Return()

	mockedStore.On("SaveProblem", scores.Problem{
		ID:           fakedProblemID,
		Zone1Enabled: true,
		Zone2Enabled: true,
		ProblemValue: scores.ProblemValue{
			PointsTop:   1000,
			PointsZone1: 500,
			PointsZone2: 750,
			FlashBonus:  100,
		},
	}).Return()

	mockedStore.On("SaveContender", scores.Contender{
		ID:                  fakedContenderID,
		CategoryID:          fakedCategoryID,
		Disqualified:        true,
		WithdrawnFromFinals: true,
	}).Return()

	mockedStore.On("SaveTick", fakedContenderID, scores.Tick{
		ID:            fakedTickID,
		Revision:      1,
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
	err := hydrator.Hydrate(context.Background(), fakedCompetitionID, mockedStore)

	require.NoError(t, err)

	mockedRepo.AssertExpectations(t)
	mockedStore.AssertExpectations(t)
}
