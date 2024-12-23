package scores_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/scores"
	"github.com/stretchr/testify/mock"
)

func TestHydrate(t *testing.T) {
	mockedRepo := new(repositoryMock)
	mockedStore := new(engineStoreMock)

	mockedContestID := domain.ContestID(rand.Int())
	mockedProblemID := domain.ProblemID(rand.Int())
	mockedContenderID := domain.ContenderID(rand.Int())
	mockedCompClassID := domain.CompClassID(rand.Int())

	now := time.Now()

	mockedRepo.
		On("GetProblemsByContest", mock.Anything, nil, mockedContestID).
		Return([]domain.Problem{
			{
				ID:         mockedProblemID,
				PointsTop:  100,
				PointsZone: 50,
				FlashBonus: 10,
			},
		}, nil)

	mockedRepo.
		On("GetContendersByContest", mock.Anything, nil, mockedContestID).
		Return([]domain.Contender{
			{
				ID:                  mockedContenderID,
				CompClassID:         mockedCompClassID,
				Disqualified:        true,
				WithdrawnFromFinals: true,
				Entered:             &now,
			},
		}, nil)

	mockedRepo.
		On("GetTicksByContest", mock.Anything, nil, mockedContestID).
		Return([]domain.Tick{
			{
				Ownership: domain.OwnershipData{
					ContenderID: &mockedContenderID,
				},
				ProblemID:    mockedProblemID,
				Top:          true,
				AttemptsTop:  999,
				Zone:         true,
				AttemptsZone: 1,
			},
		}, nil)

	mockedStore.On("SaveProblem", scores.Problem{
		ID:         mockedProblemID,
		PointsTop:  100,
		PointsZone: 50,
		FlashBonus: 10,
	}).Return()

	mockedStore.On("SaveContender", scores.Contender{
		ID:                  mockedContenderID,
		CompClassID:         mockedCompClassID,
		Disqualified:        true,
		WithdrawnFromFinals: true,
	}).Return()

	mockedStore.On("SaveTick", mockedContenderID, scores.Tick{
		ProblemID:    mockedProblemID,
		Top:          true,
		AttemptsTop:  999,
		Zone:         true,
		AttemptsZone: 1,
	}).Return()
}
