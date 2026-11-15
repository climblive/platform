package usecases_test

import (
	"context"
	"fmt"
	"testing"
	"testing/synctest"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/testutils"
	"github.com/climblive/platform/backend/internal/usecases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateRaffle(t *testing.T) {
	fakedOrganizerID := testutils.RandomResourceID[domain.OrganizerID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: fakedOrganizerID,
	}
	fakedCompetitionID := testutils.RandomResourceID[domain.CompetitionID]()
	fakedRaffleID := testutils.RandomResourceID[domain.RaffleID]()

	makeMocks := func() (*repositoryMock, *authorizerMock) {
		mockedRepo := new(repositoryMock)

		mockedRepo.
			On("GetCompetition", mock.Anything, nil, fakedCompetitionID).
			Return(domain.Competition{
				ID:        fakedCompetitionID,
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
			On("GetRafflesByCompetition", mock.Anything, nil, fakedCompetitionID).
			Return([]domain.Raffle{}, nil)

		mockedRepo.
			On("StoreRaffle", mock.Anything, nil,
				domain.Raffle{
					Ownership:     fakedOwnership,
					CompetitionID: fakedCompetitionID,
				},
			).
			Return(
				domain.Raffle{
					ID:            fakedRaffleID,
					Ownership:     fakedOwnership,
					CompetitionID: fakedCompetitionID,
				}, nil)

		ucase := usecases.RaffleUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		raffle, err := ucase.CreateRaffle(context.Background(), fakedCompetitionID)

		require.NoError(t, err)
		assert.Equal(t, fakedRaffleID, raffle.ID)
		assert.Equal(t, fakedOwnership, raffle.Ownership)
		assert.Equal(t, fakedCompetitionID, raffle.CompetitionID)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("LimitExceeded", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		raffles := make([]domain.Raffle, 10)
		mockedRepo.
			On("GetRafflesByCompetition", mock.Anything, nil, fakedCompetitionID).
			Return(raffles, nil)

		ucase := usecases.RaffleUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		_, err := ucase.CreateRaffle(context.Background(), fakedCompetitionID)

		require.ErrorIs(t, err, domain.ErrLimitExceeded)

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

		_, err := ucase.CreateRaffle(context.Background(), fakedCompetitionID)

		require.ErrorIs(t, err, domain.ErrNoOwnership)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})
}

func TestGetRaffle(t *testing.T) {
	fakedOrganizerID := testutils.RandomResourceID[domain.OrganizerID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: fakedOrganizerID,
	}
	fakedCompetitionID := testutils.RandomResourceID[domain.CompetitionID]()
	fakedRaffleID := testutils.RandomResourceID[domain.RaffleID]()

	makeMocks := func() (*repositoryMock, *authorizerMock) {
		mockedRepo := new(repositoryMock)

		mockedRepo.
			On("GetRaffle", mock.Anything, nil, fakedRaffleID).
			Return(domain.Raffle{
				ID:            fakedRaffleID,
				Ownership:     fakedOwnership,
				CompetitionID: fakedCompetitionID,
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
		assert.Equal(t, fakedCompetitionID, raffle.CompetitionID)

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

func TestGetRafflesByCompetition(t *testing.T) {
	fakedOrganizerID := testutils.RandomResourceID[domain.OrganizerID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: fakedOrganizerID,
	}
	fakedCompetitionID := testutils.RandomResourceID[domain.CompetitionID]()

	var fakedRaffles []domain.Raffle
	for range 3 {
		fakedRaffles = append(fakedRaffles, domain.Raffle{
			ID:            testutils.RandomResourceID[domain.RaffleID](),
			Ownership:     fakedOwnership,
			CompetitionID: fakedCompetitionID,
		})
	}

	makeMocks := func() (*repositoryMock, *authorizerMock) {
		mockedRepo := new(repositoryMock)

		mockedRepo.
			On("GetCompetition", mock.Anything, nil, fakedCompetitionID).
			Return(domain.Competition{
				ID:        fakedCompetitionID,
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
			On("GetRafflesByCompetition", mock.Anything, nil, fakedCompetitionID).
			Return(fakedRaffles, nil)

		ucase := usecases.RaffleUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		raffles, err := ucase.GetRafflesByCompetition(context.Background(), fakedCompetitionID)

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

		_, err := ucase.GetRafflesByCompetition(context.Background(), fakedCompetitionID)

		require.ErrorIs(t, err, domain.ErrNoOwnership)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})
}

func TestDrawRaffleWinner(t *testing.T) {
	fakedOrganizerID := testutils.RandomResourceID[domain.OrganizerID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: fakedOrganizerID,
	}
	fakedCompetitionID := testutils.RandomResourceID[domain.CompetitionID]()
	fakedRaffleID := testutils.RandomResourceID[domain.RaffleID]()

	makeMocks := func() (*repositoryMock, *eventBrokerMock, *authorizerMock) {
		mockedRepo := new(repositoryMock)

		mockedRepo.
			On("GetRaffle", mock.Anything, nil, fakedRaffleID).
			Return(domain.Raffle{
				ID:            fakedRaffleID,
				Ownership:     fakedOwnership,
				CompetitionID: fakedCompetitionID,
			}, nil)

		mockedAuthorizer := new(authorizerMock)

		mockedEventBroker := new(eventBrokerMock)

		return mockedRepo, mockedEventBroker, mockedAuthorizer
	}

	makeCompetitors := func(count int) []domain.Competitor {
		var competitors []domain.Competitor

		for i := range count {
			competitors = append(competitors, domain.Competitor{
				ID:      domain.CompetitorID(i),
				Entered: time.Now().Add(-time.Hour),
			})
		}

		return competitors
	}

	makeWinners := func(count int) []domain.RaffleWinner {
		var winners []domain.RaffleWinner

		for i := range count {
			winners = append(winners, domain.RaffleWinner{
				ID:             testutils.RandomResourceID[domain.RaffleWinnerID](),
				Ownership:      fakedOwnership,
				RaffleID:       fakedRaffleID,
				CompetitorID:   domain.CompetitorID(i),
				CompetitorName: fmt.Sprintf("Winner %d", i),
				Timestamp:      time.Now(),
			})
		}

		return winners
	}

	t.Run("SingleCompetitorInRaffle", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			mockedRepo, mockedEventBroker, mockedAuthorizer := makeMocks()

			fakedCompetitorID := testutils.RandomResourceID[domain.CompetitorID]()
			fakedRaffleWinnerID := testutils.RandomResourceID[domain.RaffleWinnerID]()

			mockedAuthorizer.
				On("HasOwnership", mock.Anything, fakedOwnership).
				Return(domain.OrganizerRole, nil)

			mockedRepo.
				On("GetCompetitorsByCompetition", mock.Anything, nil, fakedCompetitionID).
				Return([]domain.Competitor{
					{
						ID:      fakedCompetitorID,
						Name:    "John Doe",
						Entered: time.Now().Add(-time.Hour),
					},
				}, nil)

			mockedRepo.
				On("GetRaffleWinners", mock.Anything, nil, fakedRaffleID).
				Return([]domain.RaffleWinner{}, nil)

			mockedRepo.
				On("StoreRaffleWinner", mock.Anything, nil, domain.RaffleWinner{
					Ownership:      fakedOwnership,
					RaffleID:       fakedRaffleID,
					CompetitorID:   fakedCompetitorID,
					CompetitorName: "John Doe",
					Timestamp:      time.Now(),
				}).
				Return(domain.RaffleWinner{
					ID:             fakedRaffleWinnerID,
					Ownership:      fakedOwnership,
					RaffleID:       fakedRaffleID,
					CompetitorID:   fakedCompetitorID,
					CompetitorName: "John Doe",
					Timestamp:      time.Now(),
				}, nil)

			mockedEventBroker.
				On("Dispatch", fakedCompetitionID, domain.RaffleWinnerDrawnEvent{
					RaffleID:     fakedRaffleID,
					CompetitorID: fakedCompetitorID,
					Timestamp:    time.Now(),
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
			assert.Equal(t, fakedCompetitorID, winner.CompetitorID)
			assert.Equal(t, "John Doe", winner.CompetitorName)
			assert.Equal(t, time.Now(), winner.Timestamp)

			mockedRepo.AssertExpectations(t)
			mockedEventBroker.AssertExpectations(t)
			mockedAuthorizer.AssertExpectations(t)
		})
	})

	t.Run("MultipleCompetitorsInRaffle", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			mockedRepo, mockedEventBroker, mockedAuthorizer := makeMocks()

			mockedAuthorizer.
				On("HasOwnership", mock.Anything, fakedOwnership).
				Return(domain.OrganizerRole, nil)

			mockedRepo.
				On("GetCompetitorsByCompetition", mock.Anything, nil, fakedCompetitionID).
				Return(makeCompetitors(5), nil)

			mockedRepo.
				On("GetRaffleWinners", mock.Anything, nil, fakedRaffleID).
				Return([]domain.RaffleWinner{}, nil)

			for k := range 5 {
				mockedRepo.
					On("StoreRaffleWinner", mock.Anything, nil, domain.RaffleWinner{
						Ownership:    fakedOwnership,
						RaffleID:     fakedRaffleID,
						CompetitorID: domain.CompetitorID(k),
						Timestamp:    time.Now(),
					}).
					Return(domain.RaffleWinner{}, nil)
			}

			mockedEventBroker.
				On("Dispatch", fakedCompetitionID, mock.AnythingOfType("domain.RaffleWinnerDrawnEvent")).
				Return()

			ucase := usecases.RaffleUseCase{
				Repo:        mockedRepo,
				Authorizer:  mockedAuthorizer,
				EventBroker: mockedEventBroker,
			}

			for range 1_000 {
				winner, err := ucase.DrawRaffleWinner(context.Background(), fakedRaffleID)
				require.NoError(t, err)

				assert.Contains(t, []domain.CompetitorID{0, 1, 2, 3, 4}, winner.CompetitorID)
			}

			mockedRepo.AssertExpectations(t)
			mockedEventBroker.AssertExpectations(t)
			mockedAuthorizer.AssertExpectations(t)
		})
	})

	t.Run("NoRegisteredCompetitors", func(t *testing.T) {
		mockedRepo, _, mockedAuthorizer := makeMocks()

		mockedRepo.
			On("GetCompetitorsByCompetition", mock.Anything, nil, fakedCompetitionID).
			Return([]domain.Competitor{
				{
					ID: testutils.RandomResourceID[domain.CompetitorID](),
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

	t.Run("DisqualifiedCompetitorsExcluded", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			mockedRepo, mockedEventBroker, mockedAuthorizer := makeMocks()

			fakedRaffleWinnerID := testutils.RandomResourceID[domain.RaffleWinnerID]()

			mockedAuthorizer.
				On("HasOwnership", mock.Anything, fakedOwnership).
				Return(domain.OrganizerRole, nil)

			competitors := makeCompetitors(2)
			competitors[1].Disqualified = true

			mockedRepo.
				On("GetCompetitorsByCompetition", mock.Anything, nil, fakedCompetitionID).
				Return(competitors, nil)

			mockedRepo.
				On("GetRaffleWinners", mock.Anything, nil, fakedRaffleID).
				Return([]domain.RaffleWinner{}, nil)

			mockedRepo.
				On("StoreRaffleWinner", mock.Anything, nil, domain.RaffleWinner{
					Ownership:    fakedOwnership,
					RaffleID:     fakedRaffleID,
					CompetitorID: domain.CompetitorID(0),
					Timestamp:    time.Now(),
				}).
				Return(domain.RaffleWinner{
					ID:           fakedRaffleWinnerID,
					Ownership:    fakedOwnership,
					RaffleID:     fakedRaffleID,
					CompetitorID: domain.CompetitorID(0),
					Timestamp:    time.Now(),
				}, nil)

			mockedEventBroker.
				On("Dispatch", fakedCompetitionID, domain.RaffleWinnerDrawnEvent{
					RaffleID:     fakedRaffleID,
					CompetitorID: domain.CompetitorID(0),
					Timestamp:    time.Now(),
				}).
				Return()

			ucase := usecases.RaffleUseCase{
				Repo:        mockedRepo,
				Authorizer:  mockedAuthorizer,
				EventBroker: mockedEventBroker,
			}

			winner, err := ucase.DrawRaffleWinner(context.Background(), fakedRaffleID)
			require.NoError(t, err)

			assert.Equal(t, domain.CompetitorID(0), winner.CompetitorID)

			mockedRepo.AssertExpectations(t)
			mockedEventBroker.AssertExpectations(t)
			mockedAuthorizer.AssertExpectations(t)
		})
	})

	t.Run("AllEligibleDisqualified", func(t *testing.T) {
		mockedRepo, _, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		competitors := makeCompetitors(3)
		for i := range competitors {
			competitors[i].Disqualified = true
		}

		mockedRepo.
			On("GetCompetitorsByCompetition", mock.Anything, nil, fakedCompetitionID).
			Return(competitors, nil)

		mockedRepo.
			On("GetRaffleWinners", mock.Anything, nil, fakedRaffleID).
			Return([]domain.RaffleWinner{}, nil)

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
			On("GetCompetitorsByCompetition", mock.Anything, nil, fakedCompetitionID).
			Return(makeCompetitors(5), nil)

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

	t.Run("AllCompetitorsScrubbed", func(t *testing.T) {
		mockedRepo, _, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		competitors := makeCompetitors(5)
		for i := range competitors {
			competitors[i].ScrubbedAt = time.Now()
		}

		mockedRepo.
			On("GetCompetitorsByCompetition", mock.Anything, nil, fakedCompetitionID).
			Return(competitors, nil)

		mockedRepo.
			On("GetRaffleWinners", mock.Anything, nil, fakedRaffleID).
			Return([]domain.RaffleWinner{}, nil)

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
	fakedOrganizerID := testutils.RandomResourceID[domain.OrganizerID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: fakedOrganizerID,
	}
	fakedCompetitionID := testutils.RandomResourceID[domain.CompetitionID]()
	fakedRaffleID := testutils.RandomResourceID[domain.RaffleID]()

	makeMocks := func() (*repositoryMock, *authorizerMock) {
		mockedRepo := new(repositoryMock)

		mockedRepo.
			On("GetRaffle", mock.Anything, nil, fakedRaffleID).
			Return(domain.Raffle{
				ID:            fakedRaffleID,
				Ownership:     fakedOwnership,
				CompetitionID: fakedCompetitionID,
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
				ID:             testutils.RandomResourceID[domain.RaffleWinnerID](),
				RaffleID:       fakedRaffleID,
				CompetitorID:   testutils.RandomResourceID[domain.CompetitorID](),
				CompetitorName: "Winner 1",
				Timestamp:      time.Now(),
			},
			{
				ID:             testutils.RandomResourceID[domain.RaffleWinnerID](),
				RaffleID:       fakedRaffleID,
				CompetitorID:   testutils.RandomResourceID[domain.CompetitorID](),
				CompetitorName: "Winner 2",
				Timestamp:      time.Now(),
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

func TestDeleteRaffle(t *testing.T) {
	fakedOrganizerID := testutils.RandomResourceID[domain.OrganizerID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: fakedOrganizerID,
	}
	fakedCompetitionID := testutils.RandomResourceID[domain.CompetitionID]()
	fakedRaffleID := testutils.RandomResourceID[domain.RaffleID]()
	fakedRaffleWinnerID := testutils.RandomResourceID[domain.RaffleWinnerID]()

	makeMocks := func() (*repositoryMock, *transactionMock, *authorizerMock) {
		mockedRepo := new(repositoryMock)

		mockedRepo.
			On("GetRaffle", mock.Anything, nil, fakedRaffleID).
			Return(domain.Raffle{
				ID:            fakedRaffleID,
				Ownership:     fakedOwnership,
				CompetitionID: fakedCompetitionID,
			}, nil)

		mockedAuthorizer := new(authorizerMock)

		mockedTx := new(transactionMock)

		return mockedRepo, mockedTx, mockedAuthorizer
	}

	t.Run("HappyCase", func(t *testing.T) {
		mockedRepo, mockedTx, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("GetRaffleWinners", mock.Anything, nil, fakedRaffleID).
			Return([]domain.RaffleWinner{
				{ID: fakedRaffleWinnerID},
			}, nil)

		mockedRepo.On("Begin").Return(mockedTx, nil)

		mockedRepo.
			On("DeleteRaffleWinner", mock.Anything, mockedTx, fakedRaffleWinnerID).
			Return(nil)

		mockedRepo.
			On("DeleteRaffle", mock.Anything, mockedTx, fakedRaffleID).
			Return(nil)

		mockedTx.On("Commit").Return(nil)
		mockedTx.On("Rollback").Return()

		ucase := usecases.RaffleUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		err := ucase.DeleteRaffle(context.Background(), fakedRaffleID)

		require.NoError(t, err)

		mockedRepo.AssertExpectations(t)
		mockedTx.AssertExpectations(t)
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

		err := ucase.DeleteRaffle(context.Background(), fakedRaffleID)

		require.ErrorIs(t, err, domain.ErrNoOwnership)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})
}
