package usecases_test

import (
	"context"
	"fmt"
	"testing"
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
		mockedRepo, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("StoreContest", mock.Anything, nil,
				domain.Contest{
					Ownership:          fakedOwnership,
					Location:           "The garage",
					SeriesID:           0,
					Name:               "Swedish Championships",
					Description:        "Who is the best climber in Sweden?",
					QualifyingProblems: 10,
					Finalists:          7,
					Rules:              "No rules!",
					GracePeriod:        time.Hour,
				},
			).
			Return(domain.Contest{
				ID:                 fakedContestID,
				Ownership:          fakedOwnership,
				Location:           "The garage",
				SeriesID:           0,
				Name:               "Swedish Championships",
				Description:        "Who is the best climber in Sweden?",
				QualifyingProblems: 10,
				Finalists:          7,
				Rules:              "No rules!",
				GracePeriod:        time.Hour,
			}, nil)

		ucase := usecases.ContestUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		contest, err := ucase.CreateContest(context.Background(), fakedOrganizerID, domain.ContestTemplate{
			Location:           "The garage",
			Name:               "Swedish Championships",
			Description:        "Who is the best climber in Sweden?",
			QualifyingProblems: 10,
			Finalists:          7,
			Rules:              "No rules!",
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
		assert.Equal(t, "No rules!", contest.Rules)
		assert.Equal(t, time.Hour, contest.GracePeriod)
		assert.Empty(t, contest.TimeBegin)
		assert.Empty(t, contest.TimeEnd)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
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

	t.Run("RulesAreSanitized", func(t *testing.T) {
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
			Name:               "Swedish Championships",
			Description:        "Who is the best climber in Sweden?",
			QualifyingProblems: 10,
			Finalists:          7,
			Rules:              `<a href="javascript:alert('XSS1')" onmouseover="alert('XSS2')">XSS<a>`,
			GracePeriod:        time.Hour,
		})

		require.NoError(t, err)
		assert.Equal(t, "XSS", contest.Rules)

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
		Rules:              "No rules!",
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
		assert.Equal(t, "No rules!", duplicatedContest.Rules)
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

func TestPatchContest(t *testing.T) {
	t.Parallel()

	fakedContestID := randomResourceID[domain.ContestID]()
	fakedOrganizerID := randomResourceID[domain.OrganizerID]()
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
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("GetContest", mock.Anything, nil, fakedContestID).
			Return(domain.Contest{
				ID:        fakedContestID,
				Ownership: fakedOwnership,
			}, nil)

		mockedRepo.
			On("StoreContest", mock.Anything, nil,
				domain.Contest{
					ID:                 fakedContestID,
					Ownership:          fakedOwnership,
					Location:           "The garage",
					SeriesID:           domain.SeriesID(1),
					Name:               "Swedish Championships",
					Description:        "Who is the best climber in Sweden?",
					QualifyingProblems: 20,
					Finalists:          5,
					Rules:              "No rules!",
					GracePeriod:        time.Hour,
				},
			).
			Return(domain.Contest{
				ID:                 fakedContestID,
				Ownership:          fakedOwnership,
				Location:           "The garage",
				SeriesID:           domain.SeriesID(1),
				Name:               "Swedish Championships",
				Description:        "Who is the best climber in Sweden?",
				QualifyingProblems: 20,
				Finalists:          5,
				Rules:              "No rules!",
				GracePeriod:        time.Hour,
			}, nil)

		ucase := usecases.ContestUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		patch := domain.ContestPatch{
			Location:           domain.NewPatch("The garage"),
			SeriesID:           domain.NewPatch(domain.SeriesID(1)),
			Name:               domain.NewPatch("Swedish Championships"),
			Description:        domain.NewPatch("Who is the best climber in Sweden?"),
			QualifyingProblems: domain.NewPatch(20),
			Finalists:          domain.NewPatch(5),
			Rules:              domain.NewPatch("No rules!"),
			GracePeriod:        domain.NewPatch(time.Hour),
		}

		contest, err := ucase.PatchContest(context.Background(), fakedContestID, patch)

		require.NoError(t, err)
		assert.Equal(t, "The garage", contest.Location)
		assert.Equal(t, domain.SeriesID(1), contest.SeriesID)
		assert.Equal(t, "Swedish Championships", contest.Name)
		assert.Equal(t, "Who is the best climber in Sweden?", contest.Description)
		assert.Equal(t, 20, contest.QualifyingProblems)
		assert.Equal(t, 5, contest.Finalists)
		assert.Equal(t, "No rules!", contest.Rules)
		assert.Equal(t, time.Hour, contest.GracePeriod)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("ArchiveContest", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

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
			}, nil)

		mockedRepo.
			On("StoreContest", mock.Anything, nil,
				domain.Contest{
					ID:        fakedContestID,
					Ownership: fakedOwnership,
					Archived:  true,
					Name:      "Swedish Championships",
				},
			).
			Return(domain.Contest{
				ID:        fakedContestID,
				Ownership: fakedOwnership,
				Archived:  true,
				Name:      "Swedish Championships",
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
		mockedRepo, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("GetContest", mock.Anything, nil, fakedContestID).
			Return(domain.Contest{
				ID:        fakedContestID,
				Ownership: fakedOwnership,
				Name:      "Swedish Championships",
				Archived:  true,
			}, nil)

		mockedRepo.
			On("StoreContest", mock.Anything, nil,
				domain.Contest{
					ID:        fakedContestID,
					Ownership: fakedOwnership,
					Archived:  false,
					Name:      "Swedish Championships",
				},
			).
			Return(domain.Contest{
				ID:        fakedContestID,
				Ownership: fakedOwnership,
				Archived:  false,
				Name:      "Swedish Championships",
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
		mockedRepo, mockedAuthorizer := makeMocks()

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
		mockedRepo, mockedAuthorizer := makeMocks()

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

		_, err := ucase.PatchContest(context.Background(), fakedContestID, domain.ContestPatch{
			Name: domain.NewPatch("Norweigan Championships"),
		})

		assert.ErrorIs(t, err, domain.ErrArchived)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})
}
