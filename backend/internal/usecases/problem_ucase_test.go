package usecases_test

import (
	"context"
	"testing"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/usecases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetProblemsByContest(t *testing.T) {
	mockedContestID := domain.ContestID(1)
	mockedProblems := []domain.Problem{
		{
			ID:        1,
			ContestID: mockedContestID,
		},
	}

	mockedRepo := new(repositoryMock)

	mockedRepo.
		On("GetProblemsByContest", mock.Anything, mock.Anything, mockedContestID).
		Return(mockedProblems, nil)

	ucase := usecases.ProblemUseCase{
		Repo: mockedRepo,
	}

	problems, err := ucase.GetProblemsByContest(context.Background(), mockedContestID)

	require.NoError(t, err)
	assert.Equal(t, mockedProblems, problems)
}
