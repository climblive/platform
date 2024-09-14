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

func TestGetContest(t *testing.T) {
	mockedContestID := domain.ResourceID(1)
	mockedOwnership := domain.OwnershipData{
		OrganizerID: 1,
	}

	mockedContest := domain.Contest{
		ID:        mockedContestID,
		Ownership: mockedOwnership,
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
