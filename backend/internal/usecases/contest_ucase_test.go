package usecases_test

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"testing/synctest"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/scores"
	"github.com/climblive/platform/backend/internal/usecases"
	"github.com/climblive/platform/backend/internal/usecases/validators"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetContest(t *testing.T) {
	fakedContestID := randomResourceID[domain.ContestID]()

	fakedContest := domain.Contest{
		ID: fakedContestID,
	}

	mockedRepo := new(repositoryMock)

	mockedRepo.
		On("GetContest", mock.Anything, mock.Anything, fakedContestID).
		Return(fakedContest, nil)

	ucase := usecases.ContestUseCase{
		Repo: mockedRepo,
	}

	contest, err := ucase.GetContest(context.Background(), fakedContestID)

	require.NoError(t, err)
	assert.Equal(t, fakedContestID, contest.ID)

	mockedRepo.AssertExpectations(t)
}

func TestGetScoreboard(t *testing.T) {
	fakedContestID := randomResourceID[domain.ContestID]()
	fakedCompClassID := randomResourceID[domain.CompClassID]()

	mockedRepo := new(repositoryMock)
	mockedScoreKeeper := new(scoreKeeperMock)

	currentTime := time.Now()

	var contenders []domain.Contender

	for i := 1; i <= 10; i++ {
		contenderID := domain.ContenderID(i)

		fakedContender := domain.Contender{
			ID:                  contenderID,
			CompClassID:         fakedCompClassID,
			Name:                fmt.Sprintf("Climber %d", i),
			WithdrawnFromFinals: true,
			Disqualified:        true,
			Score: &domain.Score{
				ContenderID: contenderID,
				Score:       i * 10,
				Placement:   i,
				Finalist:    true,
				RankOrder:   i - 1,
				Timestamp:   currentTime,
			},
		}

		contenders = append(contenders, fakedContender)
	}

	mockedRepo.
		On("GetContendersByContest", mock.Anything, mock.Anything, fakedContestID).
		Return(contenders, nil)

	future := currentTime.Add(time.Minute)

	mockedScoreKeeper.On("GetScore", domain.ContenderID(1)).Return(domain.Score{
		Timestamp:   future,
		ContenderID: domain.ContenderID(1),
		Score:       1234,
		Placement:   42,
		Finalist:    false,
		RankOrder:   1337,
	}, nil)

	mockedScoreKeeper.On("GetScore", mock.Anything).Return(domain.Score{}, errMock)

	ucase := usecases.ContestUseCase{
		Repo:        mockedRepo,
		ScoreKeeper: mockedScoreKeeper,
	}

	scoreboard, err := ucase.GetScoreboard(context.Background(), fakedContestID)

	require.NoError(t, err)

	assert.Len(t, scoreboard, 10)

	assert.Equal(t, domain.ContenderID(1), scoreboard[0].ContenderID)
	assert.Equal(t, fakedCompClassID, scoreboard[0].CompClassID)
	assert.Equal(t, "Climber 1", scoreboard[0].Name)
	assert.Equal(t, true, scoreboard[0].WithdrawnFromFinals)
	assert.Equal(t, true, scoreboard[0].Disqualified)
	assert.NotNil(t, scoreboard[0].Score)
	assert.Equal(t, future, scoreboard[0].Score.Timestamp)
	assert.Equal(t, 1234, scoreboard[0].Score.Score)
	assert.Equal(t, 42, scoreboard[0].Score.Placement)
	assert.False(t, scoreboard[0].Score.Finalist)
	assert.Equal(t, 1337, scoreboard[0].Score.RankOrder)

	for i := 2; i <= 10; i++ {
		entry := scoreboard[i-1]

		assert.Equal(t, domain.ContenderID(i), entry.ContenderID)
		assert.Equal(t, fakedCompClassID, entry.CompClassID)
		assert.Equal(t, fmt.Sprintf("Climber %d", i), entry.Name)
		assert.Equal(t, true, entry.WithdrawnFromFinals)
		assert.Equal(t, true, entry.Disqualified)
		assert.NotNil(t, entry.Score)
		assert.Equal(t, currentTime, entry.Score.Timestamp)
		assert.Equal(t, i*10, entry.Score.Score)
		assert.Equal(t, i, entry.Score.Placement)
		assert.Equal(t, i-1, entry.Score.RankOrder)
		assert.True(t, entry.Score.Finalist)
	}

	mockedRepo.AssertExpectations(t)
	mockedScoreKeeper.AssertExpectations(t)
}

func TestGetScoreboard_Empty(t *testing.T) {
	fakedContestID := randomResourceID[domain.ContestID]()

	mockedRepo := new(repositoryMock)

	mockedRepo.
		On("GetContendersByContest", mock.Anything, mock.Anything, fakedContestID).
		Return([]domain.Contender{}, nil)

	ucase := usecases.ContestUseCase{
		Repo: mockedRepo,
	}

	scoreboard, err := ucase.GetScoreboard(context.Background(), fakedContestID)

	require.NoError(t, err)
	assert.NotNil(t, scoreboard)
	assert.Len(t, scoreboard, 0)

	mockedRepo.AssertExpectations(t)
}

func TestGetContestsByOrganizer(t *testing.T) {
	fakedOrganizerID := randomResourceID[domain.OrganizerID]()
	fakedContestID := randomResourceID[domain.ContestID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: fakedOrganizerID,
	}

	makeMocks := func() (*repositoryMock, *authorizerMock) {
		mockedRepo := new(repositoryMock)

		mockedRepo.
			On("GetOrganizer", mock.Anything, nil, fakedOrganizerID).
			Return(domain.Organizer{
				ID:        fakedOrganizerID,
				Ownership: fakedOwnership,
			}, nil)

		mockedAuthorizer := new(authorizerMock)

		return mockedRepo, mockedAuthorizer
	}

	t.Run("HappyCase", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("GetContestsByOrganizer", mock.Anything, nil, fakedOrganizerID).
			Return([]domain.Contest{
				{
					ID:        fakedContestID,
					Ownership: fakedOwnership,
				},
			}, nil)

		ucase := usecases.ContestUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		contests, err := ucase.GetContestsByOrganizer(context.Background(), fakedOrganizerID)

		require.NoError(t, err)
		require.Len(t, contests, 1)
		assert.Equal(t, fakedContestID, contests[0].ID)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		ucase := usecases.ContestUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		_, err := ucase.GetContestsByOrganizer(context.Background(), fakedOrganizerID)

		require.ErrorIs(t, err, domain.ErrNoOwnership)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})
}

func TestGetAllContests(t *testing.T) {
	fakedOrganizerID := randomResourceID[domain.OrganizerID]()
	fakedContestID := randomResourceID[domain.ContestID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: fakedOrganizerID,
	}

	makeMocks := func() (*repositoryMock, *authorizerMock) {
		mockedRepo := new(repositoryMock)

		mockedAuthorizer := new(authorizerMock)

		return mockedRepo, mockedAuthorizer
	}

	t.Run("HappyCase", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mock.AnythingOfType("domain.OwnershipData")).
			Return(domain.AdminRole, nil)

		mockedRepo.
			On("GetAllContests", mock.Anything, nil).
			Return([]domain.Contest{
				{
					ID:        fakedContestID,
					Ownership: fakedOwnership,
				},
			}, nil)

		ucase := usecases.ContestUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		contests, err := ucase.GetAllContests(context.Background())

		require.NoError(t, err)
		require.Len(t, contests, 1)
		assert.Equal(t, fakedContestID, contests[0].ID)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mock.AnythingOfType("domain.OwnershipData")).
			Return(domain.NilRole, domain.ErrNoOwnership)

		ucase := usecases.ContestUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		_, err := ucase.GetAllContests(context.Background())

		require.ErrorIs(t, err, domain.ErrNoOwnership)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("NotAdmin", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mock.AnythingOfType("domain.OwnershipData")).
			Return(domain.OrganizerRole, nil)

		ucase := usecases.ContestUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		_, err := ucase.GetAllContests(context.Background())

		require.ErrorIs(t, err, domain.ErrNotAuthorized)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})
}
func TestCreateContest(t *testing.T) {
	fakedOrganizerID := randomResourceID[domain.OrganizerID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: fakedOrganizerID,
	}
	fakedContestID := randomResourceID[domain.ContestID]()

	makeMocks := func() (*repositoryMock, *authorizerMock) {
		mockedRepo := new(repositoryMock)

		mockedRepo.
			On("GetOrganizer", mock.Anything, nil, fakedOrganizerID).
			Return(domain.Organizer{
				ID:        fakedOrganizerID,
				Ownership: fakedOwnership,
			}, nil)

		mockedAuthorizer := new(authorizerMock)

		return mockedRepo, mockedAuthorizer
	}

	t.Run("HappyCase", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			mockedRepo, mockedAuthorizer := makeMocks()

			mockedAuthorizer.
				On("HasOwnership", mock.Anything, fakedOwnership).
				Return(domain.OrganizerRole, nil)

			mockedRepo.
				On("StoreContest", mock.Anything, nil,
					domain.Contest{
						Ownership:          fakedOwnership,
						Location:           "The garage",
						Country:            "SE",
						SeriesID:           0,
						Name:               "Swedish Championships",
						Description:        "Who is the best climber in Sweden?",
						QualifyingProblems: 10,
						Finalists:          7,
						Info:               "No rules!",
						GracePeriod:        time.Hour,
						Created:            time.Now(),
					},
				).
				Return(domain.Contest{
					ID:                 fakedContestID,
					Ownership:          fakedOwnership,
					Location:           "The garage",
					Country:            "SE",
					SeriesID:           0,
					Name:               "Swedish Championships",
					Description:        "Who is the best climber in Sweden?",
					QualifyingProblems: 10,
					Finalists:          7,
					Info:               "No rules!",
					GracePeriod:        time.Hour,
					Created:            time.Now(),
				}, nil)

			ucase := usecases.ContestUseCase{
				Repo:       mockedRepo,
				Authorizer: mockedAuthorizer,
			}

			contest, err := ucase.CreateContest(context.Background(), fakedOrganizerID, domain.ContestTemplate{
				Location:           "The garage",
				Country:            "SE",
				Name:               "Swedish Championships",
				Description:        "Who is the best climber in Sweden?",
				QualifyingProblems: 10,
				Finalists:          7,
				Info:               "No rules!",
				GracePeriod:        time.Hour,
			})

			require.NoError(t, err)
			assert.Equal(t, fakedContestID, contest.ID)
			assert.Equal(t, fakedOwnership, contest.Ownership)
			assert.False(t, contest.Archived)
			assert.Equal(t, "The garage", contest.Location)
			assert.Equal(t, "Swedish Championships", contest.Name)
			assert.Equal(t, "Who is the best climber in Sweden?", contest.Description)
			assert.Equal(t, 10, contest.QualifyingProblems)
			assert.Equal(t, 7, contest.Finalists)
			assert.Equal(t, "No rules!", contest.Info)
			assert.Equal(t, time.Hour, contest.GracePeriod)
			assert.Empty(t, contest.TimeBegin)
			assert.Empty(t, contest.TimeEnd)
			assert.Equal(t, time.Now(), contest.Created)

			mockedRepo.AssertExpectations(t)
			mockedAuthorizer.AssertExpectations(t)
		})
	})

	t.Run("ValidatorIsInvoked", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		ucase := usecases.ContestUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		_, err := ucase.CreateContest(context.Background(), fakedOrganizerID, domain.ContestTemplate{})

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.True(t, validators.ContestValidator{}.IsValidationError(err))

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("InfoIsSanitized", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("StoreContest", mock.Anything, nil, mock.AnythingOfType("domain.Contest")).
			Return(mirrorInstruction{}, nil)

		ucase := usecases.ContestUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		contest, err := ucase.CreateContest(context.Background(), fakedOrganizerID, domain.ContestTemplate{
			Location:           "The garage",
			Country:            "SE",
			Name:               "Swedish Championships",
			Description:        "Who is the best climber in Sweden?",
			QualifyingProblems: 10,
			Finalists:          7,
			Info:               `<a href="javascript:alert('XSS1')" onmouseover="alert('XSS2')">XSS<a>`,
			GracePeriod:        time.Hour,
		})

		require.NoError(t, err)
		assert.Equal(t, "XSS", contest.Info)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		ucase := usecases.ContestUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		_, err := ucase.CreateContest(context.Background(), fakedOrganizerID, domain.ContestTemplate{})

		require.ErrorIs(t, err, domain.ErrNoOwnership)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})
}

func TestDuplicateContest(t *testing.T) {
	fakedContestID := randomResourceID[domain.ContestID]()
	fakedDuplicatedContestID := randomResourceID[domain.ContestID]()
	fakedCompClassID := randomResourceID[domain.CompClassID]()
	fakedProblemID := randomResourceID[domain.ProblemID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: randomResourceID[domain.OrganizerID](),
	}

	timeBegin := time.Now()
	timeEnd := timeBegin.Add(time.Hour)

	fakedContest := domain.Contest{
		ID:                 fakedContestID,
		Ownership:          fakedOwnership,
		Location:           "The garage",
		SeriesID:           randomResourceID[domain.SeriesID](),
		Name:               "Original Contest",
		Description:        "Who is the best climber in Sweden?",
		QualifyingProblems: 10,
		Finalists:          7,
		Info:               "No rules!",
		GracePeriod:        time.Hour,
		TimeBegin:          timeBegin,
		TimeEnd:            timeEnd,
	}

	fakedCompClass := domain.CompClass{
		ID:          fakedCompClassID,
		Ownership:   fakedOwnership,
		ContestID:   fakedContestID,
		Name:        "Males",
		Description: "Male climbers",
		TimeBegin:   timeBegin,
		TimeEnd:     timeEnd,
	}

	fakedProblem := domain.Problem{
		ID:                 fakedProblemID,
		Ownership:          fakedOwnership,
		ContestID:          fakedContestID,
		Number:             42,
		HoldColorPrimary:   "#FF0000",
		HoldColorSecondary: "#00FF00",
		Description:        "Test Problem",
		Zone1Enabled:       true,
		Zone2Enabled:       true,
		PointsTop:          100,
		PointsZone1:        50,
		PointsZone2:        75,
		FlashBonus:         20,
	}

	makeMocks := func() (*repositoryMock, *authorizerMock) {
		mockedRepo := new(repositoryMock)
		mockedAuthorizer := new(authorizerMock)

		return mockedRepo, mockedAuthorizer
	}

	t.Run("HappyCase", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

		fakedDuplicatedContest := fakedContest
		fakedDuplicatedContest.ID = fakedDuplicatedContestID
		fakedDuplicatedContest.Name = "Original Contest (Copy)"

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("GetContest", mock.Anything, nil, fakedContestID).
			Return(fakedContest, nil)

		mockedRepo.
			On("GetCompClassesByContest", mock.Anything, nil, fakedContestID).
			Return([]domain.CompClass{fakedCompClass}, nil)

		mockedRepo.
			On("GetProblemsByContest", mock.Anything, nil, fakedContestID).
			Return([]domain.Problem{fakedProblem}, nil)

		mockedTx := new(transactionMock)

		mockedRepo.
			On("Begin").
			Return(mockedTx, nil)

		mockedRepo.
			On("StoreContest", mock.Anything, mockedTx, mock.MatchedBy(func(contest domain.Contest) bool {
				expected := fakedDuplicatedContest
				expected.ID = 0

				return contest == expected
			})).
			Return(fakedDuplicatedContest, nil)

		mockedRepo.
			On("StoreCompClass", mock.Anything, mockedTx, mock.MatchedBy(func(compClass domain.CompClass) bool {
				expected := fakedCompClass
				expected.ID = 0
				expected.ContestID = fakedDuplicatedContestID

				return compClass == expected
			})).
			Return(domain.CompClass{}, nil)

		mockedRepo.
			On("StoreProblem", mock.Anything, mockedTx, mock.MatchedBy(func(problem domain.Problem) bool {
				expected := fakedProblem
				expected.ID = 0
				expected.ContestID = fakedDuplicatedContestID

				return problem == expected
			})).
			Return(domain.Problem{}, nil)

		mockedTx.
			On("Commit").
			Return(nil)

		ucase := usecases.ContestUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		duplicatedContest, err := ucase.DuplicateContest(context.Background(), fakedContestID)

		require.NoError(t, err)
		assert.Equal(t, fakedDuplicatedContestID, duplicatedContest.ID)
		assert.Equal(t, "Original Contest (Copy)", duplicatedContest.Name)
		assert.Equal(t, fakedOwnership, duplicatedContest.Ownership)
		assert.Equal(t, "The garage", duplicatedContest.Location)
		assert.Equal(t, fakedContest.SeriesID, duplicatedContest.SeriesID)
		assert.Equal(t, "Who is the best climber in Sweden?", duplicatedContest.Description)
		assert.Equal(t, 10, duplicatedContest.QualifyingProblems)
		assert.Equal(t, 7, duplicatedContest.Finalists)
		assert.Equal(t, "No rules!", duplicatedContest.Info)
		assert.Equal(t, time.Hour, duplicatedContest.GracePeriod)
		assert.Equal(t, timeBegin, duplicatedContest.TimeBegin)
		assert.NotZero(t, timeEnd, duplicatedContest.TimeEnd)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		mockedRepo.
			On("GetContest", mock.Anything, nil, fakedContestID).
			Return(fakedContest, nil)

		ucase := usecases.ContestUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		_, err := ucase.DuplicateContest(context.Background(), fakedContestID)

		require.ErrorIs(t, err, domain.ErrNoOwnership)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("ContestArchived", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("GetContest", mock.Anything, nil, fakedContestID).
			Return(domain.Contest{
				ID:        fakedContestID,
				Ownership: fakedOwnership,
				Archived:  true,
			}, nil)

		ucase := usecases.ContestUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		_, err := ucase.DuplicateContest(context.Background(), fakedContestID)

		require.ErrorIs(t, err, domain.ErrArchived)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})
}

func TestTransferContest(t *testing.T) {
	fakedContestID := randomResourceID[domain.ContestID]()
	fakedOldOrganizerID := randomResourceID[domain.OrganizerID]()
	fakedNewOrganizerID := randomResourceID[domain.OrganizerID]()
	fakedCompClassID := randomResourceID[domain.CompClassID]()
	fakedProblemID := randomResourceID[domain.ProblemID]()
	fakedContenderID := randomResourceID[domain.ContenderID]()
	fakedRaffleID := randomResourceID[domain.RaffleID]()
	fakedRaffleWinnerID := randomResourceID[domain.RaffleWinnerID]()
	fakedTickID := randomResourceID[domain.TickID]()
	fakedSeriesID := randomResourceID[domain.SeriesID]()

	now := time.Now()

	fakedOldOwnership := domain.OwnershipData{OrganizerID: fakedOldOrganizerID}
	fakedNewOwnership := domain.OwnershipData{OrganizerID: fakedNewOrganizerID}

	timeBegin := time.Now()
	timeEnd := timeBegin.Add(2 * time.Hour)

	fakedContest := domain.Contest{
		ID:                 fakedContestID,
		Ownership:          fakedOldOwnership,
		Location:           "Stockholm Climbing Center",
		SeriesID:           fakedSeriesID,
		Name:               "Swedish Open 2025",
		Description:        "National bouldering championship",
		QualifyingProblems: 8,
		Finalists:          6,
		Info:               "Standard IFSC rules apply",
		GracePeriod:        30 * time.Minute,
		TimeBegin:          timeBegin,
		TimeEnd:            timeEnd,
		Created:            now.Add(time.Duration(rand.Int())),
	}

	fakedCompClass := domain.CompClass{
		ID:          fakedCompClassID,
		Ownership:   fakedOldOwnership,
		ContestID:   fakedContestID,
		Name:        "Males",
		Description: "Male competitors",
		TimeBegin:   timeBegin,
		TimeEnd:     timeEnd,
	}

	fakedProblem := domain.Problem{
		ID:                 fakedProblemID,
		Ownership:          fakedOldOwnership,
		ContestID:          fakedContestID,
		Number:             15,
		HoldColorPrimary:   "#0000FF",
		HoldColorSecondary: "#FFFFFF",
		Description:        "Technical slab problem",
		Zone1Enabled:       true,
		Zone2Enabled:       true,
		PointsTop:          100,
		PointsZone1:        25,
		PointsZone2:        50,
		FlashBonus:         10,
	}

	fakedScore := domain.Score{
		Timestamp:   now.Add(time.Duration(rand.Int())),
		ContenderID: fakedContenderID,
		Score:       10,
		Placement:   3,
		Finalist:    true,
		RankOrder:   2,
	}

	fakedContender := domain.Contender{
		ID:                  fakedContenderID,
		Ownership:           domain.OwnershipData{OrganizerID: fakedOldOrganizerID, ContenderID: &fakedContenderID},
		ContestID:           fakedContestID,
		CompClassID:         fakedCompClassID,
		RegistrationCode:    "ABCD1234",
		Name:                "John Doe",
		Entered:             now.Add(time.Duration(rand.Int())),
		WithdrawnFromFinals: true,
		Disqualified:        true,
		Score:               &fakedScore,
	}

	fakedRaffle := domain.Raffle{
		ID:        fakedRaffleID,
		Ownership: fakedOldOwnership,
		ContestID: fakedContestID,
	}

	fakedWinner := domain.RaffleWinner{
		ID:            fakedRaffleWinnerID,
		Ownership:     fakedOldOwnership,
		RaffleID:      fakedRaffleID,
		ContenderID:   fakedContenderID,
		ContenderName: "John Doe",
		Timestamp:     now.Add(time.Duration(rand.Int())),
	}

	fakedTick := domain.Tick{
		ID:            fakedTickID,
		Ownership:     domain.OwnershipData{OrganizerID: fakedOldOrganizerID, ContenderID: &fakedContenderID},
		Timestamp:     now.Add(time.Duration(rand.Int())),
		ContestID:     fakedContestID,
		ProblemID:     fakedProblemID,
		Zone1:         true,
		AttemptsZone1: 2,
		Zone2:         true,
		AttemptsZone2: 3,
		Top:           true,
		AttemptsTop:   5,
	}

	makeMocks := func() (*repositoryMock, *authorizerMock, *transactionMock) {
		mockedRepo := new(repositoryMock)
		mockedAuthorizer := new(authorizerMock)
		mockedTx := new(transactionMock)

		return mockedRepo, mockedAuthorizer, mockedTx
	}

	t.Run("HappyCase", func(t *testing.T) {
		mockedRepo, mockedAuthorizer, mockedTx := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOldOwnership).
			Return(domain.OrganizerRole, nil).
			On("HasOwnership", mock.Anything, fakedNewOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("GetContest", mock.Anything, nil, fakedContestID).
			Return(fakedContest, nil)
		mockedRepo.
			On("GetOrganizer", mock.Anything, nil, fakedNewOrganizerID).
			Return(domain.Organizer{
				ID:        fakedNewOrganizerID,
				Ownership: fakedNewOwnership,
			}, nil)
		mockedRepo.On("GetCompClassesByContest", mock.Anything, nil, fakedContestID).
			Return([]domain.CompClass{fakedCompClass}, nil)
		mockedRepo.
			On("GetProblemsByContest", mock.Anything, nil, fakedContestID).
			Return([]domain.Problem{fakedProblem}, nil)
		mockedRepo.
			On("GetContendersByContest", mock.Anything, nil, fakedContestID).
			Return([]domain.Contender{fakedContender}, nil)
		mockedRepo.
			On("GetRafflesByContest", mock.Anything, nil, fakedContestID).
			Return([]domain.Raffle{fakedRaffle}, nil)
		mockedRepo.
			On("GetRaffleWinners", mock.Anything, nil, fakedRaffleID).
			Return([]domain.RaffleWinner{fakedWinner}, nil)
		mockedRepo.
			On("GetTicksByContest", mock.Anything, nil, fakedContestID).
			Return([]domain.Tick{fakedTick}, nil)

		mockedRepo.On("Begin").Return(mockedTx, nil)

		mockedRepo.On("DeleteTick", mock.Anything, mockedTx, fakedTickID).Return(nil)
		mockedRepo.On("DeleteRaffleWinner", mock.Anything, mockedTx, fakedRaffleWinnerID).Return(nil)
		mockedRepo.On("DeleteRaffle", mock.Anything, mockedTx, fakedRaffleID).Return(nil)
		mockedRepo.On("DeleteContender", mock.Anything, mockedTx, fakedContenderID).Return(nil)
		mockedRepo.On("DeleteProblem", mock.Anything, mockedTx, fakedProblemID).Return(nil)
		mockedRepo.On("DeleteCompClass", mock.Anything, mockedTx, fakedCompClassID).Return(nil)
		mockedRepo.On("DeleteContest", mock.Anything, mockedTx, fakedContestID).Return(nil)

		mockedRepo.
			On("StoreContest", mock.Anything, mockedTx, domain.Contest{
				ID:                 fakedContestID,
				Ownership:          fakedNewOwnership,
				Location:           "Stockholm Climbing Center",
				SeriesID:           fakedSeriesID,
				Name:               "Swedish Open 2025",
				Description:        "National bouldering championship",
				QualifyingProblems: 8,
				Finalists:          6,
				Info:               "Standard IFSC rules apply",
				GracePeriod:        30 * time.Minute,
				TimeBegin:          fakedContest.TimeBegin,
				TimeEnd:            fakedContest.TimeEnd,
				Created:            fakedContest.Created,
			}).
			Return(domain.Contest{}, nil)

		mockedRepo.
			On("StoreCompClass", mock.Anything, mockedTx, domain.CompClass{
				ID:          fakedCompClassID,
				Ownership:   fakedNewOwnership,
				ContestID:   fakedContestID,
				Name:        "Males",
				Description: "Male competitors",
				TimeBegin:   fakedCompClass.TimeBegin,
				TimeEnd:     fakedCompClass.TimeEnd,
			}).
			Return(domain.CompClass{}, nil)

		mockedRepo.
			On("StoreProblem", mock.Anything, mockedTx, domain.Problem{
				ID:                 fakedProblemID,
				Ownership:          fakedNewOwnership,
				ContestID:          fakedContestID,
				Number:             15,
				HoldColorPrimary:   "#0000FF",
				HoldColorSecondary: "#FFFFFF",
				Description:        "Technical slab problem",
				Zone1Enabled:       true,
				Zone2Enabled:       true,
				PointsTop:          100,
				PointsZone1:        25,
				PointsZone2:        50,
				FlashBonus:         10,
			}).Return(domain.Problem{}, nil)

		mockedRepo.
			On("StoreContender", mock.Anything, mockedTx, domain.Contender{
				ID:                  fakedContenderID,
				Ownership:           domain.OwnershipData{OrganizerID: fakedNewOrganizerID, ContenderID: &fakedContenderID},
				ContestID:           fakedContestID,
				CompClassID:         fakedCompClassID,
				RegistrationCode:    "ABCD1234",
				Name:                "John Doe",
				Entered:             fakedContender.Entered,
				WithdrawnFromFinals: true,
				Disqualified:        true,
				Score:               &fakedScore,
			}).
			Return(domain.Contender{}, nil)

		mockedRepo.
			On("StoreScore", mock.Anything, mockedTx, fakedScore).
			Return(nil)

		mockedRepo.
			On("StoreRaffle", mock.Anything, mockedTx, domain.Raffle{
				ID:        fakedRaffleID,
				Ownership: fakedNewOwnership,
				ContestID: fakedContestID,
			}).
			Return(domain.Raffle{}, nil)

		mockedRepo.
			On("StoreRaffleWinner", mock.Anything, mockedTx, domain.RaffleWinner{
				ID:            fakedRaffleWinnerID,
				Ownership:     fakedNewOwnership,
				RaffleID:      fakedRaffleID,
				ContenderID:   fakedContenderID,
				ContenderName: "John Doe",
				Timestamp:     fakedWinner.Timestamp,
			}).
			Return(domain.RaffleWinner{}, nil)

		mockedRepo.
			On("StoreTick", mock.Anything, mockedTx, domain.Tick{
				ID:            fakedTickID,
				Ownership:     domain.OwnershipData{OrganizerID: fakedNewOrganizerID, ContenderID: &fakedContenderID},
				Timestamp:     fakedTick.Timestamp,
				ContestID:     fakedContestID,
				ProblemID:     fakedProblemID,
				Zone1:         true,
				AttemptsZone1: 2,
				Zone2:         true,
				AttemptsZone2: 3,
				Top:           true,
				AttemptsTop:   5,
			}).
			Return(domain.Tick{}, nil)

		mockedTx.On("Commit").Return(nil)

		ucase := usecases.ContestUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		contest, err := ucase.TransferContest(context.Background(), fakedContestID, fakedNewOrganizerID)

		require.NoError(t, err)
		assert.Equal(t, fakedNewOrganizerID, contest.Ownership.OrganizerID)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
		mockedTx.AssertExpectations(t)
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedRepo, mockedAuthorizer, _ := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOldOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		mockedRepo.
			On("GetContest", mock.Anything, nil, fakedContestID).
			Return(fakedContest, nil)

		ucase := usecases.ContestUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		_, err := ucase.TransferContest(context.Background(), fakedContestID, fakedNewOrganizerID)

		assert.ErrorIs(t, err, domain.ErrNoOwnership)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("NewOrganizerUnauthorized", func(t *testing.T) {
		mockedRepo, mockedAuthorizer, _ := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOldOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("GetContest", mock.Anything, nil, fakedContestID).
			Return(fakedContest, nil)

		mockedRepo.
			On("GetOrganizer", mock.Anything, nil, fakedNewOrganizerID).
			Return(domain.Organizer{
				ID:        fakedNewOrganizerID,
				Ownership: fakedNewOwnership,
			}, nil)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedNewOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		ucase := usecases.ContestUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		_, err := ucase.TransferContest(context.Background(), fakedContestID, fakedNewOrganizerID)

		assert.ErrorIs(t, err, domain.ErrNoOwnership)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("ContestArchived", func(t *testing.T) {
		mockedRepo, _, _ := makeMocks()

		archivedContest := fakedContest
		archivedContest.Archived = true

		mockedRepo.
			On("GetContest", mock.Anything, nil, fakedContestID).
			Return(archivedContest, nil)

		ucase := usecases.ContestUseCase{
			Repo: mockedRepo,
		}

		_, err := ucase.TransferContest(context.Background(), fakedContestID, fakedNewOrganizerID)

		assert.ErrorIs(t, err, domain.ErrArchived)

		mockedRepo.AssertExpectations(t)
	})
}

func TestPatchContest(t *testing.T) {
	t.Parallel()

	fakedContestID := randomResourceID[domain.ContestID]()
	fakedOrganizerID := randomResourceID[domain.OrganizerID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: fakedOrganizerID,
	}

	makeMocks := func() (*repositoryMock, *authorizerMock, *eventBrokerMock) {
		mockedRepo := new(repositoryMock)

		mockedAuthorizer := new(authorizerMock)

		mockedEventBroker := new(eventBrokerMock)

		return mockedRepo, mockedAuthorizer, mockedEventBroker
	}

	t.Run("HappyCase", func(t *testing.T) {
		mockedRepo, mockedAuthorizer, mockedEventBroker := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("GetContest", mock.Anything, nil, fakedContestID).
			Return(domain.Contest{
				ID:        fakedContestID,
				Ownership: fakedOwnership,
				Country:   "SE",
			}, nil)

		mockedRepo.
			On("StoreContest", mock.Anything, nil,
				domain.Contest{
					ID:                 fakedContestID,
					Ownership:          fakedOwnership,
					Location:           "The garage",
					Country:            "SE",
					SeriesID:           domain.SeriesID(1),
					Name:               "Swedish Championships",
					Description:        "Who is the best climber in Sweden?",
					QualifyingProblems: 20,
					Finalists:          5,
					Info:               "No rules!",
					GracePeriod:        time.Hour,
				},
			).
			Return(domain.Contest{
				ID:                 fakedContestID,
				Ownership:          fakedOwnership,
				Location:           "The garage",
				Country:            "SE",
				SeriesID:           domain.SeriesID(1),
				Name:               "Swedish Championships",
				Description:        "Who is the best climber in Sweden?",
				QualifyingProblems: 20,
				Finalists:          5,
				Info:               "No rules!",
				GracePeriod:        time.Hour,
			}, nil)

		mockedEventBroker.
			On("Dispatch", fakedContestID, domain.RulesUpdatedEvent{
				QualifyingProblems: 20,
				Finalists:          5,
			}).
			Return(nil)

		ucase := usecases.ContestUseCase{
			Repo:        mockedRepo,
			Authorizer:  mockedAuthorizer,
			EventBroker: mockedEventBroker,
		}

		patch := domain.ContestPatch{
			Location:           domain.NewPatch("The garage"),
			Country:           domain.NewPatch("SE"),
			SeriesID:           domain.NewPatch(domain.SeriesID(1)),
			Name:               domain.NewPatch("Swedish Championships"),
			Description:        domain.NewPatch("Who is the best climber in Sweden?"),
			QualifyingProblems: domain.NewPatch(20),
			Finalists:          domain.NewPatch(5),
			Info:               domain.NewPatch("No rules!"),
			GracePeriod:        domain.NewPatch(time.Hour),
		}

		contest, err := ucase.PatchContest(context.Background(), fakedContestID, patch)

		require.NoError(t, err)
		assert.Equal(t, "The garage", contest.Location)
		assert.Equal(t, "SE", contest.Country)
		assert.Equal(t, domain.SeriesID(1), contest.SeriesID)
		assert.Equal(t, "Swedish Championships", contest.Name)
		assert.Equal(t, "Who is the best climber in Sweden?", contest.Description)
		assert.Equal(t, 20, contest.QualifyingProblems)
		assert.Equal(t, 5, contest.Finalists)
		assert.Equal(t, "No rules!", contest.Info)
		assert.Equal(t, time.Hour, contest.GracePeriod)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
		mockedEventBroker.AssertExpectations(t)
	})

	t.Run("ArchiveContest", func(t *testing.T) {
		mockedRepo, mockedAuthorizer, _ := makeMocks()

		mockedScoreEngineManager := new(scoreEngineManagerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("GetContest", mock.Anything, nil, fakedContestID).
			Return(domain.Contest{
				ID:        fakedContestID,
				Ownership: fakedOwnership,
				Name:      "Swedish Championships",
				Country:   "SE",
			}, nil)

		mockedRepo.
			On("StoreContest", mock.Anything, nil,
				domain.Contest{
					ID:        fakedContestID,
					Ownership: fakedOwnership,
					Archived:  true,
					Name:      "Swedish Championships",
					Country:   "SE",
				},
			).
			Return(domain.Contest{
				ID:        fakedContestID,
				Ownership: fakedOwnership,
				Archived:  true,
				Name:      "Swedish Championships",
				Country:   "SE",
			}, nil)

		fakedScoreEngineInstanceID := domain.ScoreEngineInstanceID(uuid.New())

		mockedScoreEngineManager.
			On("ListScoreEnginesByContest", mock.Anything, fakedContestID).
			Return([]scores.ScoreEngineDescriptor{
				{
					InstanceID: fakedScoreEngineInstanceID,
					ContestID:  fakedContestID,
				},
			}, nil)

		mockedScoreEngineManager.
			On("StopScoreEngine", mock.Anything, fakedScoreEngineInstanceID).
			Return(nil)

		ucase := usecases.ContestUseCase{
			Repo:               mockedRepo,
			Authorizer:         mockedAuthorizer,
			ScoreEngineManager: mockedScoreEngineManager,
		}

		patch := domain.ContestPatch{
			Archived: domain.NewPatch(true),
		}

		contest, err := ucase.PatchContest(context.Background(), fakedContestID, patch)

		require.NoError(t, err)
		assert.True(t, contest.Archived)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
		mockedScoreEngineManager.AssertExpectations(t)
	})

	t.Run("RestoreContest", func(t *testing.T) {
		mockedRepo, mockedAuthorizer, _ := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("GetContest", mock.Anything, nil, fakedContestID).
			Return(domain.Contest{
				ID:        fakedContestID,
				Ownership: fakedOwnership,
				Name:      "Swedish Championships",
				Country:   "SE",
				Archived:  true,
			}, nil)

		mockedRepo.
			On("StoreContest", mock.Anything, nil,
				domain.Contest{
					ID:        fakedContestID,
					Ownership: fakedOwnership,
					Archived:  false,
					Name:      "Swedish Championships",
					Country:   "SE",
				},
			).
			Return(domain.Contest{
				ID:        fakedContestID,
				Ownership: fakedOwnership,
				Archived:  false,
				Name:      "Swedish Championships",
				Country:   "SE",
			}, nil)

		ucase := usecases.ContestUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		patch := domain.ContestPatch{
			Archived: domain.NewPatch(false),
		}

		contest, err := ucase.PatchContest(context.Background(), fakedContestID, patch)

		require.NoError(t, err)
		assert.False(t, contest.Archived)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedRepo, mockedAuthorizer, _ := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		mockedRepo.
			On("GetContest", mock.Anything, nil, fakedContestID).
			Return(domain.Contest{
				ID:        fakedContestID,
				Ownership: fakedOwnership,
			}, nil)

		ucase := usecases.ContestUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		_, err := ucase.PatchContest(context.Background(), fakedContestID, domain.ContestPatch{})

		assert.ErrorIs(t, err, domain.ErrNoOwnership)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("ValidatorIsInvoked", func(t *testing.T) {
		mockedRepo, mockedAuthorizer, _ := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("GetContest", mock.Anything, nil, fakedContestID).
			Return(domain.Contest{
				ID:        fakedContestID,
				Ownership: fakedOwnership,
			}, nil)

		ucase := usecases.ContestUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		_, err := ucase.PatchContest(context.Background(), fakedContestID, domain.ContestPatch{})

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.True(t, validators.ContestValidator{}.IsValidationError(err))

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("ContestIsArchived", func(t *testing.T) {
		mockedRepo, mockedAuthorizer, _ := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("GetContest", mock.Anything, nil, fakedContestID).
			Return(domain.Contest{
				ID:        fakedContestID,
				Ownership: fakedOwnership,
				Archived:  true,
			}, nil)

		ucase := usecases.ContestUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		_, err := ucase.PatchContest(context.Background(), fakedContestID, domain.ContestPatch{
			Name: domain.NewPatch("Norweigan Championships"),
		})

		assert.ErrorIs(t, err, domain.ErrArchived)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})
}
