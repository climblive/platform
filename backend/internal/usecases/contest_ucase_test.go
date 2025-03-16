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
