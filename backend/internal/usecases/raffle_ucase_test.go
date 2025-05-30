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

		raffle, err := ucase.CreateRaffle(context.Background(), fakedContestID)

		require.NoError(t, err)
		assert.Equal(t, fakedRaffleID, raffle.ID)
		assert.Equal(t, fakedOwnership, raffle.Ownership)
		assert.Equal(t, fakedContestID, raffle.ContestID)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
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

		_, err := ucase.CreateRaffle(context.Background(), fakedContestID)

		require.ErrorIs(t, err, domain.ErrNoOwnership)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})
}

func TestGetRaffle(t *testing.T) {
	fakedOrganizerID := randomResourceID[domain.OrganizerID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: fakedOrganizerID,
	}
	fakedContestID := randomResourceID[domain.ContestID]()
	fakedRaffleID := randomResourceID[domain.RaffleID]()

	makeMocks := func() (*repositoryMock, *authorizerMock) {
		mockedRepo := new(repositoryMock)

		mockedRepo.
			On("GetRaffle", mock.Anything, nil, fakedRaffleID).
			Return(domain.Raffle{
				ID:        fakedRaffleID,
				Ownership: fakedOwnership,
				ContestID: fakedContestID,
			}, nil)

		mockedAuthorizer := new(authorizerMock)

		return mockedRepo, mockedAuthorizer
	}

	t.Run("HappyCase", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		ucase := usecases.RaffleUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		raffle, err := ucase.GetRaffle(context.Background(), fakedRaffleID)

		require.NoError(t, err)
		assert.Equal(t, fakedRaffleID, raffle.ID)
		assert.Equal(t, fakedOwnership, raffle.Ownership)
		assert.Equal(t, fakedContestID, raffle.ContestID)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
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

		_, err := ucase.GetRaffle(context.Background(), fakedRaffleID)

		require.ErrorIs(t, err, domain.ErrNoOwnership)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})
}

func TestGetRafflesByContest(t *testing.T) {
	fakedOrganizerID := randomResourceID[domain.OrganizerID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: fakedOrganizerID,
	}
	fakedContestID := randomResourceID[domain.ContestID]()

	var fakedRaffles []domain.Raffle
	for range 3 {
		fakedRaffles = append(fakedRaffles, domain.Raffle{
			ID:        randomResourceID[domain.RaffleID](),
			Ownership: fakedOwnership,
			ContestID: fakedContestID,
		})
	}

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

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("GetRafflesByContest", mock.Anything, nil, fakedContestID).
			Return(fakedRaffles, nil)

		ucase := usecases.RaffleUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		raffles, err := ucase.GetRafflesByContest(context.Background(), fakedContestID)

		require.NoError(t, err)
		assert.Equal(t, fakedRaffles, raffles)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
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

		_, err := ucase.GetRafflesByContest(context.Background(), fakedContestID)

		require.ErrorIs(t, err, domain.ErrNoOwnership)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})
}
