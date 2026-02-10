package usecases_test

import (
	"context"
	"fmt"
	"testing"
	"testing/synctest"
	"time"

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

func TestDrawRaffleWinner(t *testing.T) {
	fakedOrganizerID := randomResourceID[domain.OrganizerID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: fakedOrganizerID,
	}
	fakedContestID := randomResourceID[domain.ContestID]()
	fakedRaffleID := randomResourceID[domain.RaffleID]()

	makeMocks := func() (*repositoryMock, *eventBrokerMock, *authorizerMock) {
		mockedRepo := new(repositoryMock)

		mockedRepo.
			On("GetRaffle", mock.Anything, nil, fakedRaffleID).
			Return(domain.Raffle{
				ID:        fakedRaffleID,
				Ownership: fakedOwnership,
				ContestID: fakedContestID,
			}, nil)

		mockedAuthorizer := new(authorizerMock)

		mockedEventBroker := new(eventBrokerMock)

		return mockedRepo, mockedEventBroker, mockedAuthorizer
	}

	makeContenders := func(count int) []domain.Contender {
		var contenders []domain.Contender

		for i := range count {
			contenders = append(contenders, domain.Contender{
				ID:      domain.ContenderID(i),
				Entered: time.Now().Add(-time.Hour),
			})
		}

		return contenders
	}

	makeWinners := func(count int) []domain.RaffleWinner {
		var winners []domain.RaffleWinner

		for i := range count {
			winners = append(winners, domain.RaffleWinner{
				Ownership:     fakedOwnership,
				RaffleID:      fakedRaffleID,
				ContenderID:   domain.ContenderID(i),
				ContenderName: fmt.Sprintf("Winner %d", i),
				Timestamp:     time.Now(),
			})
		}

		return winners
	}

	t.Run("SingleContenderInRaffle", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			mockedRepo, mockedEventBroker, mockedAuthorizer := makeMocks()

			fakedContenderID := randomResourceID[domain.ContenderID]()
			fakedRaffleWinnerID := randomResourceID[domain.RaffleWinnerID]()

			mockedAuthorizer.
				On("HasOwnership", mock.Anything, fakedOwnership).
				Return(domain.OrganizerRole, nil)

			mockedRepo.
				On("GetContendersByContest", mock.Anything, nil, fakedContestID).
				Return([]domain.Contender{
					{
						ID:      fakedContenderID,
						Name:    "John Doe",
						Entered: time.Now().Add(-time.Hour),
					},
				}, nil)

			mockedRepo.
				On("GetRaffleWinners", mock.Anything, nil, fakedRaffleID).
				Return([]domain.RaffleWinner{}, nil)

			mockedRepo.
				On("StoreRaffleWinner", mock.Anything, nil, domain.RaffleWinner{
					Ownership:     fakedOwnership,
					RaffleID:      fakedRaffleID,
					ContenderID:   fakedContenderID,
					ContenderName: "John Doe",
					Timestamp:     time.Now(),
				}).
				Return(domain.RaffleWinner{
					ID:            fakedRaffleWinnerID,
					Ownership:     fakedOwnership,
					RaffleID:      fakedRaffleID,
					ContenderID:   fakedContenderID,
					ContenderName: "John Doe",
					Timestamp:     time.Now(),
				}, nil)

			mockedEventBroker.
				On("Dispatch", fakedContestID, domain.RaffleWinnerDrawnEvent{
					RaffleID:      fakedRaffleID,
					ContenderID:   fakedContenderID,
					ContenderName: "John Doe",
					Timestamp:     time.Now(),
				}).
				Return()

			ucase := usecases.RaffleUseCase{
				Repo:        mockedRepo,
				Authorizer:  mockedAuthorizer,
				EventBroker: mockedEventBroker,
			}

			winner, err := ucase.DrawRaffleWinner(context.Background(), fakedRaffleID)

			require.NoError(t, err)
			assert.Equal(t, fakedRaffleID, winner.RaffleID)
			assert.Equal(t, fakedOwnership, winner.Ownership)
			assert.Equal(t, fakedContenderID, winner.ContenderID)
			assert.Equal(t, "John Doe", winner.ContenderName)
			assert.Equal(t, time.Now(), winner.Timestamp)

			mockedRepo.AssertExpectations(t)
			mockedEventBroker.AssertExpectations(t)
			mockedAuthorizer.AssertExpectations(t)
		})
	})

	t.Run("MultipleContendersInRaffle", func(t *testing.T) {
		mockedRepo, mockedEventBroker, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("GetContendersByContest", mock.Anything, nil, fakedContestID).
			Return(makeContenders(5), nil)

		mockedRepo.
			On("GetRaffleWinners", mock.Anything, nil, fakedRaffleID).
			Return([]domain.RaffleWinner{}, nil)

		mockedRepo.
			On("StoreRaffleWinner", mock.Anything, nil, mock.AnythingOfType("domain.RaffleWinner")).
			Return(mirrorInstruction{}, nil)

		mockedEventBroker.
			On("Dispatch", fakedContestID, mock.AnythingOfType("domain.RaffleWinnerDrawnEvent")).
			Return()

		ucase := usecases.RaffleUseCase{
			Repo:        mockedRepo,
			Authorizer:  mockedAuthorizer,
			EventBroker: mockedEventBroker,
		}

		for range 100 {
			winner, err := ucase.DrawRaffleWinner(context.Background(), fakedRaffleID)
			require.NoError(t, err)

			assert.Contains(t, []domain.ContenderID{0, 1, 2, 3, 4}, winner.ContenderID)
		}

		mockedRepo.AssertExpectations(t)
		mockedEventBroker.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("NoRegisteredContenders", func(t *testing.T) {
		mockedRepo, _, mockedAuthorizer := makeMocks()

		mockedRepo.
			On("GetContendersByContest", mock.Anything, nil, fakedContestID).
			Return([]domain.Contender{
				{
					ID: randomResourceID[domain.ContenderID](),
				},
			}, nil)

		mockedRepo.
			On("GetRaffleWinners", mock.Anything, nil, fakedRaffleID).
			Return([]domain.RaffleWinner{}, nil)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		ucase := usecases.RaffleUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		_, err := ucase.DrawRaffleWinner(context.Background(), fakedRaffleID)

		require.ErrorIs(t, err, domain.ErrAllWinnersDrawn)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("AllWinnersDrawn", func(t *testing.T) {
		mockedRepo, _, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("GetContendersByContest", mock.Anything, nil, fakedContestID).
			Return(makeContenders(5), nil)

		mockedRepo.
			On("GetRaffleWinners", mock.Anything, nil, fakedRaffleID).
			Return(makeWinners(5), nil)

		ucase := usecases.RaffleUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		_, err := ucase.DrawRaffleWinner(context.Background(), fakedRaffleID)
		require.ErrorIs(t, err, domain.ErrAllWinnersDrawn)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedRepo, _, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		ucase := usecases.RaffleUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		_, err := ucase.DrawRaffleWinner(context.Background(), fakedRaffleID)

		require.ErrorIs(t, err, domain.ErrNoOwnership)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})
}

func TestGetRaffleWinners(t *testing.T) {
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

		fakedWinners := []domain.RaffleWinner{
			{
				ID:            randomResourceID[domain.RaffleWinnerID](),
				RaffleID:      fakedRaffleID,
				ContenderID:   randomResourceID[domain.ContenderID](),
				ContenderName: "Winner 1",
				Timestamp:     time.Now(),
			},
			{
				ID:            randomResourceID[domain.RaffleWinnerID](),
				RaffleID:      fakedRaffleID,
				ContenderID:   randomResourceID[domain.ContenderID](),
				ContenderName: "Winner 2",
				Timestamp:     time.Now(),
			},
		}

		mockedRepo.
			On("GetRaffleWinners", mock.Anything, nil, fakedRaffleID).
			Return(fakedWinners, nil)

		ucase := usecases.RaffleUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		winners, err := ucase.GetRaffleWinners(context.Background(), fakedRaffleID)

		require.NoError(t, err)
		assert.Equal(t, fakedWinners, winners)

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

		_, err := ucase.GetRaffleWinners(context.Background(), fakedRaffleID)

		require.ErrorIs(t, err, domain.ErrNoOwnership)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})
}
