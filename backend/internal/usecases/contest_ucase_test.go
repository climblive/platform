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
	mockedContestID := domain.ResourceID(1)

	mockedContest := domain.Contest{
		ID: mockedContestID,
	}

	mockedRepo := new(repositoryMock)

	mockedRepo.
		On("GetContest", mock.Anything, mock.Anything, mockedContestID).
		Return(mockedContest, nil)

	ucase := usecases.ContestUseCase{
		Repo: mockedRepo,
	}

	contest, err := ucase.GetContest(context.Background(), mockedContestID)

	require.NoError(t, err)
	assert.Equal(t, mockedContestID, contest.ID)
}

func TestGetScoreboard(t *testing.T) {
	mockedContestID := domain.ResourceID(1)
	mockedRepo := new(repositoryMock)
	mockedScoreKeeper := new(scoreKeeperMock)
	currentTime := time.Now()

	var contenders []domain.Contender

	for i := 1; i <= 10; i++ {
		mockedContender := domain.Contender{
			ID:                  i,
			CompClassID:         1,
			PublicName:          fmt.Sprintf("Climber %d", i),
			ClubName:            "Testers' Climbing Club",
			WithdrawnFromFinals: true,
			Disqualified:        true,
			Score:               i * 10,
			Placement:           i,
			Finalist:            true,
			ScoreUpdated:        &currentTime,
		}

		contenders = append(contenders, mockedContender)
	}

	mockedRepo.
		On("GetContendersByContest", mock.Anything, mock.Anything, mockedContestID).
		Return(contenders, nil)

	future := currentTime.Add(time.Minute)

	mockedScoreKeeper.On("GetScore", 1).Return(domain.Score{
		Timestamp:   &future,
		ContenderID: 1,
		Score:       1234,
		Placement:   42,
		Finalist:    false,
	}, nil)

	mockedScoreKeeper.On("GetScore", mock.Anything).Return(domain.Score{}, errMock)

	ucase := usecases.ContestUseCase{
		Repo:        mockedRepo,
		ScoreKeeper: mockedScoreKeeper,
	}

	scoreboard, err := ucase.GetScoreboard(context.Background(), mockedContestID)

	require.NoError(t, err)

	assert.Len(t, scoreboard, 10)

	assert.Equal(t, 1, scoreboard[0].ContenderID)
	assert.Equal(t, 1, scoreboard[0].CompClassID)
	assert.Equal(t, "Climber 1", scoreboard[0].PublicName)
	assert.Equal(t, "Testers' Climbing Club", scoreboard[0].ClubName)
	assert.Equal(t, true, scoreboard[0].WithdrawnFromFinals)
	assert.Equal(t, true, scoreboard[0].Disqualified)
	assert.NotNil(t, scoreboard[0].ScoreUpdated)
	assert.Equal(t, future, *scoreboard[0].ScoreUpdated)
	assert.Equal(t, 1234, scoreboard[0].Score)
	assert.Equal(t, 42, scoreboard[0].Placement)
	assert.Equal(t, false, scoreboard[0].Finalist)

	for i := 2; i <= 10; i++ {
		entry := scoreboard[i-1]

		assert.Equal(t, i, entry.ContenderID)
		assert.Equal(t, 1, entry.CompClassID)
		assert.Equal(t, fmt.Sprintf("Climber %d", i), entry.PublicName)
		assert.Equal(t, "Testers' Climbing Club", entry.ClubName)
		assert.Equal(t, true, entry.WithdrawnFromFinals)
		assert.Equal(t, true, entry.Disqualified)
		assert.NotNil(t, entry.ScoreUpdated)
		assert.Equal(t, currentTime, *entry.ScoreUpdated)
		assert.Equal(t, i*10, entry.Score)
		assert.Equal(t, i, entry.Placement)
		assert.Equal(t, true, entry.Finalist)
	}
}
