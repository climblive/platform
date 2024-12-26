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
	fakedContestID := randomResourceID[domain.ContestID]()
	fakedProblems := []domain.Problem{
		{
			ID:        randomResourceID[domain.ProblemID](),
			ContestID: fakedContestID,
		},
	}

	mockedRepo := new(repositoryMock)

	mockedRepo.
		On("GetProblemsByContest", mock.Anything, mock.Anything, fakedContestID).
		Return(fakedProblems, nil)

	ucase := usecases.ProblemUseCase{
		Repo: mockedRepo,
	}

	problems, err := ucase.GetProblemsByContest(context.Background(), fakedContestID)

	require.NoError(t, err)
	assert.Equal(t, fakedProblems, problems)

	mockedRepo.AssertExpectations(t)
}
