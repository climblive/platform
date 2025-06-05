package usecases_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/usecases"
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
			PublicName:          fmt.Sprintf("Climber %d", i),
			ClubName:            "Testers' Climbing Club",
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
	assert.Equal(t, "Climber 1", scoreboard[0].PublicName)
	assert.Equal(t, "Testers' Climbing Club", scoreboard[0].ClubName)
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
		assert.Equal(t, fmt.Sprintf("Climber %d", i), entry.PublicName)
		assert.Equal(t, "Testers' Climbing Club", entry.ClubName)
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

	t.Run("InvalidData", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("StoreContest", mock.Anything, nil, mock.AnythingOfType("domain.Contest")).
			Return(domain.Contest{}, nil)

		ucase := usecases.ContestUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		validTemplate := func() domain.ContestTemplate {
			return domain.ContestTemplate{
				Location:           "The garage",
				Name:               "Swedish Championships",
				Description:        "Who is the best climber in Sweden?",
				QualifyingProblems: 10,
				Finalists:          7,
				Rules:              "No rules!",
				GracePeriod:        time.Hour,
			}
		}

		_, err := ucase.CreateContest(context.Background(), fakedOrganizerID, validTemplate())

		require.NoError(t, err)

		t.Run("EmptyName", func(t *testing.T) {
			tmpl := validTemplate()
			tmpl.Name = ""

			_, err := ucase.CreateContest(context.Background(), fakedOrganizerID, tmpl)

			assert.ErrorIs(t, err, domain.ErrInvalidData)
		})

		t.Run("NegativeFinalists", func(t *testing.T) {
			tmpl := validTemplate()
			tmpl.Finalists = -1

			_, err := ucase.CreateContest(context.Background(), fakedOrganizerID, tmpl)

			assert.ErrorIs(t, err, domain.ErrInvalidData)
		})

		t.Run("NegativeQualifyingProblems", func(t *testing.T) {
			tmpl := validTemplate()
			tmpl.QualifyingProblems = -1

			_, err := ucase.CreateContest(context.Background(), fakedOrganizerID, tmpl)

			assert.ErrorIs(t, err, domain.ErrInvalidData)
		})

		t.Run("NegativeGracePeriod", func(t *testing.T) {
			tmpl := validTemplate()
			tmpl.GracePeriod = -1 * time.Nanosecond

			_, err := ucase.CreateContest(context.Background(), fakedOrganizerID, tmpl)

			assert.ErrorIs(t, err, domain.ErrInvalidData)
		})

		t.Run("GracePeriodLongerThanOneHour", func(t *testing.T) {
			tmpl := validTemplate()
			tmpl.GracePeriod = time.Hour + time.Nanosecond

			_, err := ucase.CreateContest(context.Background(), fakedOrganizerID, tmpl)

			assert.ErrorIs(t, err, domain.ErrInvalidData)
		})

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
		PointsTop:          100,
		PointsZone:         50,
		FlashBonus:         20,
	}

	makeMocks := func() (*repositoryMock, *authorizerMock) {
		mockedRepo := new(repositoryMock)
		mockedAuthorizer := new(authorizerMock)

		mockedRepo.
			On("GetContest", mock.Anything, nil, fakedContestID).
			Return(fakedContest, nil)

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

		ucase := usecases.ContestUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		_, err := ucase.DuplicateContest(context.Background(), fakedContestID)

		require.ErrorIs(t, err, domain.ErrNoOwnership)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})
}
