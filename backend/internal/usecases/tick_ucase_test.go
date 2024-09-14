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

func TestGetTicksByContender(t *testing.T) {
	mockedContenderID := domain.ResourceID(1)
	mockedOwnership := domain.OwnershipData{
		OrganizerID: 1,
		ContenderID: &mockedContenderID,
	}

	mockedContender := domain.Contender{
		ID:        mockedContenderID,
		Ownership: mockedOwnership,
	}

	mockedTicks := []domain.Tick{
		{
			ID: 1,
		},
	}

	mockedRepo := new(repositoryMock)

	mockedRepo.
		On("GetContender", mock.Anything, mock.Anything, mockedContenderID).
		Return(mockedContender, nil)

	mockedRepo.
		On("GetTicksByContender", mock.Anything, mock.Anything, mockedContenderID).
		Return(mockedTicks, nil)

	t.Run("HappyPath", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mockedOwnership).
			Return(domain.ContenderRole, nil)

		ucase := usecases.TickUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		ticks, err := ucase.GetTicksByContender(context.Background(), mockedContenderID)

		require.NoError(t, err)
		assert.Equal(t, mockedTicks, ticks)
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mockedOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		ucase := usecases.TickUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		ticks, err := ucase.GetTicksByContender(context.Background(), mockedContenderID)

		assert.ErrorIs(t, err, domain.ErrNoOwnership)
		assert.Nil(t, ticks)
	})
}
