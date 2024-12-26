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

func TestGetCompClassesByContest(t *testing.T) {
	t.Parallel()

	mockedContestID := randomResourceID[domain.ContestID]()

	mockedCompClasses := []domain.CompClass{{
		ID:        randomResourceID[domain.CompClassID](),
		ContestID: mockedContestID,
	},
	}

	mockedRepo := new(repositoryMock)

	mockedRepo.
		On("GetCompClassesByContest", mock.Anything, mock.Anything, mockedContestID).
		Return(mockedCompClasses, nil)

	ucase := usecases.CompClassUseCase{
		Repo: mockedRepo,
	}

	compClasses, err := ucase.GetCompClassesByContest(context.Background(), mockedContestID)

	require.NoError(t, err)
	assert.Equal(t, mockedCompClasses, compClasses)

	mockedRepo.AssertExpectations(t)
}
