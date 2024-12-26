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

	fakedContestID := randomResourceID[domain.ContestID]()

	fakedCompClasses := []domain.CompClass{{
		ID:        randomResourceID[domain.CompClassID](),
		ContestID: fakedContestID,
	},
	}

	mockedRepo := new(repositoryMock)

	mockedRepo.
		On("GetCompClassesByContest", mock.Anything, mock.Anything, fakedContestID).
		Return(fakedCompClasses, nil)

	ucase := usecases.CompClassUseCase{
		Repo: mockedRepo,
	}

	compClasses, err := ucase.GetCompClassesByContest(context.Background(), fakedContestID)

	require.NoError(t, err)
	assert.Equal(t, fakedCompClasses, compClasses)

	mockedRepo.AssertExpectations(t)
}
