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

func TestCreateRaffle(t *testing.T) {
	fakedOrganizerID := randomResourceID[domain.OrganizerID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: fakedOrganizerID,
	}
	fakedContestID := randomResourceID[domain.ContestID]()
	fakedRaffleID := randomResourceID[domain.RaffleID]()

	makeMocks := func() (*repositoryMock, *authorizerMock) {
		mockedRepo := new(repositoryMock)

		mockedRepo.
			On("GetContest", mock.Anything, nil, fakedContestID).
			Return(domain.Contest{
				ID:        fakedContestID,
				Ownership: fakedOwnership,
			}, nil)

		mockedAuthorizer := new(authorizerMock)

		return mockedRepo, mockedAuthorizer
	}

	t.Run("HappyCase", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

		mockedEventBroker := new(eventBrokerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("StoreRaffle", mock.Anything, nil,
				domain.Raffle{
					Ownership: fakedOwnership,
					ContestID: fakedContestID,
				},
			).
			Return(
				domain.Raffle{
					ID:        fakedRaffleID,
					Ownership: fakedOwnership,
					ContestID: fakedContestID,
				}, nil)

		ucase := usecases.RaffleUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		problem, err := ucase.CreateRaffle(context.Background(), fakedContestID, domain.RaffleTemplate{
			ContestID: fakedContestID,
		})

		require.NoError(t, err)
		assert.Equal(t, fakedRaffleID, problem.ID)
		assert.Equal(t, fakedOwnership, problem.Ownership)
		assert.Equal(t, fakedContestID, problem.ContestID)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
		mockedEventBroker.AssertExpectations(t)
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		ucase := usecases.RaffleUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		_, err := ucase.CreateRaffle(context.Background(), fakedContestID, domain.RaffleTemplate{})

		require.ErrorIs(t, err, domain.ErrNoOwnership)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})
}
