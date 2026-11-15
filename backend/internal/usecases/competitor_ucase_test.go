package usecases_test

import (
	"context"
	"fmt"
	"slices"
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

func TestGetCompetitor(t *testing.T) {
	fakedCompetitorID := testutils.RandomResourceID[domain.CompetitorID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID:  testutils.RandomResourceID[domain.OrganizerID](),
		CompetitorID: &fakedCompetitorID,
	}

	t.Run("HappyPath", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedRepo := new(repositoryMock)
		mockedResultKeeper := new(resultKeeperMock)

		fakedResult := domain.Result{
			Timestamp:    time.Now(),
			CompetitorID: fakedCompetitorID,
			Score:        "1000p",
			Placement:    5,
			Finalist:     true,
			RankOrder:    6,
		}

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.CompetitorRole, nil)

		mockedRepo.
			On("GetCompetitor", mock.Anything, mock.Anything, fakedCompetitorID).
			Return(domain.Competitor{
				ID:        fakedCompetitorID,
				Ownership: fakedOwnership,
			}, nil)

		mockedResultKeeper.On("GetResult", fakedCompetitorID).Return(fakedResult, nil)

		ucase := usecases.CompetitorUseCase{
			Repo:         mockedRepo,
			Authorizer:   mockedAuthorizer,
			ResultKeeper: mockedResultKeeper,
		}

		competitor, err := ucase.GetCompetitor(context.Background(), fakedCompetitorID)

		require.NoError(t, err)

		assert.Equal(t, fakedCompetitorID, competitor.ID)
		require.NotNil(t, competitor.Result)
		assert.Equal(t, fakedResult, *competitor.Result)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
		mockedResultKeeper.AssertExpectations(t)
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedRepo := new(repositoryMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		mockedRepo.
			On("GetCompetitor", mock.Anything, mock.Anything, fakedCompetitorID).
			Return(domain.Competitor{
				ID:        fakedCompetitorID,
				Ownership: fakedOwnership,
			}, nil)

		ucase := usecases.CompetitorUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		competitor, err := ucase.GetCompetitor(context.Background(), fakedCompetitorID)

		assert.ErrorIs(t, err, domain.ErrNoOwnership)
		assert.Empty(t, competitor)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})
}

func TestGetCompetitorByCode(t *testing.T) {
	fakedCompetitorID := testutils.RandomResourceID[domain.CompetitorID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID:  testutils.RandomResourceID[domain.OrganizerID](),
		CompetitorID: &fakedCompetitorID,
	}

	t.Run("HappyPath", func(t *testing.T) {
		mockedRepo := new(repositoryMock)
		mockedResultKeeper := new(resultKeeperMock)

		fakedResult := domain.Result{
			Timestamp:    time.Now(),
			CompetitorID: fakedCompetitorID,
			Score:        "1000p",
			Placement:    5,
			Finalist:     true,
			RankOrder:    6,
		}

		mockedRepo.
			On("GetCompetitorByCode", mock.Anything, mock.Anything, "ABCD1234").
			Return(domain.Competitor{
				ID:        fakedCompetitorID,
				Ownership: fakedOwnership,
			}, nil)

		mockedResultKeeper.On("GetResult", fakedCompetitorID).Return(fakedResult, nil)

		ucase := usecases.CompetitorUseCase{
			Repo:         mockedRepo,
			ResultKeeper: mockedResultKeeper,
		}

		competitor, err := ucase.GetCompetitorByCode(context.Background(), "ABCD1234")

		require.NoError(t, err)
		assert.Equal(t, fakedCompetitorID, competitor.ID)
		require.NotNil(t, competitor.Result)
		assert.Equal(t, fakedResult, *competitor.Result)

		mockedRepo.AssertExpectations(t)
		mockedResultKeeper.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockedRepo := new(repositoryMock)

		mockedRepo.
			On("GetCompetitorByCode", mock.Anything, mock.Anything, "DEADBEEF").
			Return(domain.Competitor{}, domain.ErrNotFound)

		ucase := usecases.CompetitorUseCase{
			Repo: mockedRepo,
		}

		competitor, err := ucase.GetCompetitorByCode(context.Background(), "DEADBEEF")

		assert.ErrorIs(t, err, domain.ErrNotFound)
		assert.Empty(t, competitor)

		mockedRepo.AssertExpectations(t)
	})
}

func TestGetPointValues(t *testing.T) {
	fakedCompetitorID := testutils.RandomResourceID[domain.CompetitorID]()

	fakedOwnership := domain.OwnershipData{
		OrganizerID:  testutils.RandomResourceID[domain.OrganizerID](),
		CompetitorID: &fakedCompetitorID,
	}

	t.Run("HappyPath", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedRepo := new(repositoryMock)
		mockedPointValueKeeper := new(pointValueKeeperMock)

		fakedPointValues := []domain.PointValue{
			{
				CompetitorID: fakedCompetitorID,
				ProblemID:    testutils.RandomResourceID[domain.ProblemID](),
				Current:      1_000,
				Zone1:        10,
				Zone2:        100,
				Top:          1_000,
				FlashBonus:   5,
			},
			{
				CompetitorID: fakedCompetitorID,
				ProblemID:    testutils.RandomResourceID[domain.ProblemID](),
				Current:      2_000,
				Zone1:        20,
				Zone2:        200,
				Top:          2_000,
				FlashBonus:   10,
			},
			{
				CompetitorID: fakedCompetitorID,
				ProblemID:    testutils.RandomResourceID[domain.ProblemID](),
				Current:      3_000,
				Zone1:        30,
				Zone2:        300,
				Top:          3_000,
				FlashBonus:   15,
			},
		}

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.CompetitorRole, nil)

		mockedRepo.
			On("GetCompetitor", mock.Anything, nil, fakedCompetitorID).
			Return(domain.Competitor{
				ID:        fakedCompetitorID,
				Ownership: fakedOwnership,
			}, nil)

		mockedPointValueKeeper.
			On("GetPointValues", fakedCompetitorID).
			Return(slices.Collect(slices.Values(fakedPointValues)))

		ucase := usecases.CompetitorUseCase{
			Repo:             mockedRepo,
			Authorizer:       mockedAuthorizer,
			PointValueKeeper: mockedPointValueKeeper,
		}

		pointValues, err := ucase.GetPointValues(context.Background(), fakedCompetitorID)

		require.NoError(t, err)
		assert.Equal(t, fakedPointValues, pointValues)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
		mockedPointValueKeeper.AssertExpectations(t)
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedRepo := new(repositoryMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		mockedRepo.
			On("GetCompetitor", mock.Anything, nil, fakedCompetitorID).
			Return(domain.Competitor{
				ID:        fakedCompetitorID,
				Ownership: fakedOwnership,
			}, nil)

		ucase := usecases.CompetitorUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		pointValues, err := ucase.GetPointValues(context.Background(), fakedCompetitorID)

		assert.ErrorIs(t, err, domain.ErrNoOwnership)
		assert.Empty(t, pointValues)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})
}

func TestGetCompetitorsByCategory(t *testing.T) {
	fakedCategoryID := testutils.RandomResourceID[domain.CategoryID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: testutils.RandomResourceID[domain.OrganizerID](),
	}
	currentTime := time.Now()

	t.Run("HappyPath", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedRepo := new(repositoryMock)
		mockedResultKeeper := new(resultKeeperMock)

		var competitors []domain.Competitor

		for k := 1; k <= 10; k++ {
			competitorID := domain.CompetitorID(k)

			competitors = append(competitors, domain.Competitor{
				ID: competitorID,
			})

			mockedResultKeeper.On("GetResult", competitorID).Return(domain.Result{
				Timestamp:    currentTime,
				CompetitorID: competitorID,
				Score:        fmt.Sprintf("%dp", k*10),
				Placement:    k,
				RankOrder:    k - 1,
				Finalist:     true,
			}, nil)
		}

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("GetCategory", mock.Anything, mock.Anything, fakedCategoryID).
			Return(domain.Category{
				ID:        fakedCategoryID,
				Ownership: fakedOwnership,
			}, nil)

		mockedRepo.
			On("GetCompetitorsByCategory", mock.Anything, mock.Anything, fakedCategoryID).
			Return(competitors, nil)

		ucase := usecases.CompetitorUseCase{
			Repo:         mockedRepo,
			Authorizer:   mockedAuthorizer,
			ResultKeeper: mockedResultKeeper,
		}

		competitors, err := ucase.GetCompetitorsByCategory(context.Background(), fakedCategoryID)

		require.NoError(t, err)
		assert.Len(t, competitors, 10)

		for i, competitor := range competitors {
			assert.Equal(t, domain.CompetitorID(i+1), competitor.ID)
			require.NotNil(t, competitor.Result)
			assert.Equal(t, fmt.Sprintf("%dp", (i+1)*10), competitor.Result.Score)
			assert.Equal(t, i+1, competitor.Result.Placement)
			assert.Equal(t, i, competitor.Result.RankOrder)
			assert.True(t, competitor.Result.Finalist)
			assert.Equal(t, currentTime, competitor.Result.Timestamp)
		}

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
		mockedResultKeeper.AssertExpectations(t)
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedRepo := new(repositoryMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		mockedRepo.
			On("GetCategory", mock.Anything, mock.Anything, fakedCategoryID).
			Return(domain.Category{
				ID:        fakedCategoryID,
				Ownership: fakedOwnership,
			}, nil)

		ucase := usecases.CompetitorUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		competitors, err := ucase.GetCompetitorsByCategory(context.Background(), fakedCategoryID)

		assert.ErrorIs(t, err, domain.ErrNoOwnership)
		assert.Nil(t, competitors)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})
}

func TestGetCompetitorsByCompetition(t *testing.T) {
	fakedCompetitionID := testutils.RandomCompetitionID()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: testutils.RandomResourceID[domain.OrganizerID](),
	}
	currentTime := time.Now()

	t.Run("HappyPath", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedRepo := new(repositoryMock)
		mockedResultKeeper := new(resultKeeperMock)

		var competitors []domain.Competitor

		for k := 1; k <= 10; k++ {
			competitorID := domain.CompetitorID(k)

			competitors = append(competitors, domain.Competitor{
				ID: competitorID,
			})

			mockedResultKeeper.On("GetResult", competitorID).Return(domain.Result{
				Timestamp:    currentTime,
				CompetitorID: competitorID,
				Score:        fmt.Sprintf("%dp", k*10),
				Placement:    k,
				RankOrder:    k - 1,
				Finalist:     true,
			}, nil)
		}

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("GetCompetition", mock.Anything, mock.Anything, fakedCompetitionID).
			Return(domain.Competition{
				ID:        fakedCompetitionID,
				Ownership: fakedOwnership,
			}, nil)

		mockedRepo.
			On("GetCompetitorsByCompetition", mock.Anything, mock.Anything, fakedCompetitionID).
			Return(competitors, nil)

		ucase := usecases.CompetitorUseCase{
			Repo:         mockedRepo,
			Authorizer:   mockedAuthorizer,
			ResultKeeper: mockedResultKeeper,
		}

		competitors, err := ucase.GetCompetitorsByCompetition(context.Background(), fakedCompetitionID)

		require.NoError(t, err)
		assert.Len(t, competitors, 10)

		for i, competitor := range competitors {
			assert.Equal(t, domain.CompetitorID(i+1), competitor.ID)
			require.NotNil(t, competitor.Result)
			assert.Equal(t, fmt.Sprintf("%dp", (i+1)*10), competitor.Result.Score)
			assert.Equal(t, i+1, competitor.Result.Placement)
			assert.Equal(t, i, competitor.Result.RankOrder)
			assert.True(t, competitor.Result.Finalist)
			assert.Equal(t, currentTime, competitor.Result.Timestamp)
		}

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
		mockedResultKeeper.AssertExpectations(t)
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedRepo := new(repositoryMock)

		mockedRepo.
			On("GetCompetition", mock.Anything, mock.Anything, fakedCompetitionID).
			Return(domain.Competition{
				ID:        fakedCompetitionID,
				Ownership: fakedOwnership,
			}, nil)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		ucase := usecases.CompetitorUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		competitors, err := ucase.GetCompetitorsByCompetition(context.Background(), fakedCompetitionID)

		assert.ErrorIs(t, err, domain.ErrNoOwnership)
		assert.Nil(t, competitors)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})
}

func TestDeleteCompetitor(t *testing.T) {
	fakedCompetitorID := testutils.RandomResourceID[domain.CompetitorID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID:  testutils.RandomResourceID[domain.OrganizerID](),
		CompetitorID: &fakedCompetitorID,
	}

	t.Run("HappyPath", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedRepo := new(repositoryMock)

		mockedRepo.
			On("GetCompetitor", mock.Anything, mock.Anything, fakedCompetitorID).
			Return(domain.Competitor{Ownership: fakedOwnership}, nil)

		mockedRepo.
			On("DeleteCompetitor", mock.Anything, mock.Anything, fakedCompetitorID).
			Return(nil)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		ucase := usecases.CompetitorUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		err := ucase.DeleteCompetitor(context.Background(), fakedCompetitorID)

		require.NoError(t, err)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedRepo := new(repositoryMock)

		mockedRepo.
			On("GetCompetitor", mock.Anything, mock.Anything, fakedCompetitorID).
			Return(domain.Competitor{Ownership: fakedOwnership}, nil)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		ucase := usecases.CompetitorUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		err := ucase.DeleteCompetitor(context.Background(), fakedCompetitorID)

		assert.ErrorIs(t, err, domain.ErrNoOwnership)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("InsufficientRole", func(t *testing.T) {
		for _, insufficientRole := range []domain.AuthRole{domain.NilRole, domain.CompetitorRole, domain.JudgeRole} {
			mockedAuthorizer := new(authorizerMock)
			mockedRepo := new(repositoryMock)

			mockedAuthorizer.
				On("HasOwnership", mock.Anything, fakedOwnership).
				Return(insufficientRole, nil)

			mockedRepo.
				On("GetCompetitor", mock.Anything, mock.Anything, fakedCompetitorID).
				Return(domain.Competitor{Ownership: fakedOwnership}, nil)

			ucase := usecases.CompetitorUseCase{
				Repo:       mockedRepo,
				Authorizer: mockedAuthorizer,
			}

			err := ucase.DeleteCompetitor(context.Background(), fakedCompetitorID)

			assert.ErrorIs(t, err, domain.ErrInsufficientRole)

			mockedAuthorizer.AssertExpectations(t)
			mockedRepo.AssertExpectations(t)
		}
	})
}

func TestCreateCompetitors(t *testing.T) {
	fakedCompetitionID := testutils.RandomCompetitionID()
	fakedOrganizerID := testutils.RandomResourceID[domain.OrganizerID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: fakedOrganizerID,
	}

	makeMocks := func() (*repositoryMock, *transactionMock, *codeGeneratorMock) {
		mockedRepo := new(repositoryMock)
		mockedTx := new(transactionMock)
		mockedCodeGenerator := new(codeGeneratorMock)

		mockedRepo.
			On("GetCompetition", mock.Anything, mock.Anything, fakedCompetitionID).
			Return(domain.Competition{
				ID:        fakedCompetitionID,
				Ownership: fakedOwnership,
			}, nil)

		return mockedRepo, mockedTx, mockedCodeGenerator
	}

	t.Run("HappyPath", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedRepo, mockedTx, mockedCodeGenerator := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("GetNumberOfCompetitors", mock.Anything, mock.Anything, fakedCompetitionID).
			Return(400, nil)

		for n := range 100 {
			code := fmt.Sprintf("%08d", n)
			competitor := domain.Competitor{
				CompetitionID: fakedCompetitionID,
				Ownership: domain.OwnershipData{
					OrganizerID: fakedOrganizerID,
				},
				RegistrationCode: code,
			}

			mockedCodeGenerator.
				On("Generate", 8).
				Return(code).Once()

			mockedRepo.
				On("StoreCompetitor", mock.Anything, mock.Anything, competitor).
				Return(competitor, nil)
		}

		mockedRepo.
			On("Begin").
			Return(mockedTx, nil)

		mockedTx.On("Commit").Return(nil)

		ucase := usecases.CompetitorUseCase{
			Repo:                      mockedRepo,
			Authorizer:                mockedAuthorizer,
			RegistrationCodeGenerator: mockedCodeGenerator,
		}

		competitors, err := ucase.CreateCompetitors(context.Background(), fakedCompetitionID, 100)

		require.NoError(t, err)
		assert.Len(t, competitors, 100)

		for i, competitor := range competitors {
			assert.Equal(t, fmt.Sprintf("%08d", i), competitor.RegistrationCode)
		}

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
		mockedTx.AssertExpectations(t)
		mockedCodeGenerator.AssertExpectations(t)
	})

	t.Run("CannotExceed500Competitors", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedRepo, _, mockedCodeGenerator := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("GetNumberOfCompetitors", mock.Anything, mock.Anything, fakedCompetitionID).
			Return(400, nil)

		ucase := usecases.CompetitorUseCase{
			Repo:                      mockedRepo,
			Authorizer:                mockedAuthorizer,
			RegistrationCodeGenerator: mockedCodeGenerator,
		}

		competitors, err := ucase.CreateCompetitors(context.Background(), fakedCompetitionID, 101)

		assert.ErrorIs(t, err, domain.ErrLimitExceeded)
		assert.Nil(t, competitors)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
		mockedCodeGenerator.AssertExpectations(t)
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedRepo, _, _ := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		ucase := usecases.CompetitorUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		competitors, err := ucase.CreateCompetitors(context.Background(), fakedCompetitionID, 100)

		assert.ErrorIs(t, err, domain.ErrNoOwnership)
		assert.Empty(t, competitors)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("Rollback", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedRepo, mockedTx, mockedCodeGenerator := makeMocks()

		mockedRepo.
			On("Begin").
			Return(mockedTx, nil)

		mockedTx.On("Rollback").Return()

		mockedRepo.
			On("GetNumberOfCompetitors", mock.Anything, mock.Anything, fakedCompetitionID).
			Return(0, nil)

		mockedRepo.
			On("StoreCompetitor", mock.Anything, mock.Anything, mock.Anything).
			Return(domain.Competitor{}, errMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.AdminRole, nil)

		mockedCodeGenerator.
			On("Generate", 8).
			Return("ABCD1234")

		ucase := usecases.CompetitorUseCase{
			Repo:                      mockedRepo,
			Authorizer:                mockedAuthorizer,
			RegistrationCodeGenerator: mockedCodeGenerator,
		}

		competitors, err := ucase.CreateCompetitors(context.Background(), fakedCompetitionID, 1)

		assert.Error(t, err)
		assert.Nil(t, competitors)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
		mockedTx.AssertExpectations(t)
		mockedCodeGenerator.AssertExpectations(t)
	})
}

func TestPatchCompetitor(t *testing.T) {
	fakedCompetitorID := testutils.RandomResourceID[domain.CompetitorID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID:  testutils.RandomResourceID[domain.OrganizerID](),
		CompetitorID: &fakedCompetitorID,
	}
	fakedCompetitionID := testutils.RandomCompetitionID()
	fakedCategoryID := testutils.RandomResourceID[domain.CategoryID]()

	currentTime := time.Now()
	gracePeriod := 15 * time.Minute
	fakedNameRetentionTime := 14 * 24 * time.Hour

	makeMockedRepo := func(competitor domain.Competitor) *repositoryMock {
		fakedCompetition := domain.Competition{
			ID:                fakedCompetitionID,
			GracePeriod:       gracePeriod,
			NameRetentionTime: fakedNameRetentionTime,
		}

		mockedRepo := new(repositoryMock)

		mockedRepo.
			On("GetCompetitor", mock.Anything, mock.Anything, competitor.ID).
			Return(competitor, nil)

		mockedRepo.
			On("GetCompetition", mock.Anything, mock.Anything, fakedCompetitionID).
			Return(fakedCompetition, nil)

		return mockedRepo
	}

	t.Run("UpdateWithoutChanges", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedResultKeeper := new(resultKeeperMock)

		fakedCompetitor := domain.Competitor{
			ID:                  fakedCompetitorID,
			Ownership:           fakedOwnership,
			CompetitionID:       fakedCompetitionID,
			CategoryID:          fakedCategoryID,
			RegistrationCode:    "ABCD1234",
			Name:                "John Doe",
			Entered:             currentTime,
			WithdrawnFromFinals: false,
			Disqualified:        false,
		}

		mockedRepo := makeMockedRepo(fakedCompetitor)

		mockedRepo.
			On("GetCategory", mock.Anything, mock.Anything, fakedCategoryID).
			Return(domain.Category{
				ID:        fakedCategoryID,
				TimeBegin: currentTime.Add(-1 * time.Hour),
				TimeEnd:   currentTime.Add(time.Hour),
			}, nil)

		mockedRepo.
			On("StoreCompetitor", mock.Anything, mock.Anything,
				domain.Competitor{
					ID:                  fakedCompetitorID,
					Ownership:           fakedOwnership,
					CompetitionID:       fakedCompetitionID,
					CategoryID:          fakedCategoryID,
					RegistrationCode:    "ABCD1234",
					Name:                "John Doe",
					Entered:             currentTime,
					WithdrawnFromFinals: false,
					Disqualified:        false,
				},
			).
			Return(domain.Competitor{
				ID:                  fakedCompetitorID,
				Ownership:           fakedOwnership,
				CompetitionID:       fakedCompetitionID,
				CategoryID:          fakedCategoryID,
				RegistrationCode:    "ABCD1234",
				Name:                "John Doe",
				Entered:             currentTime,
				WithdrawnFromFinals: false,
				Disqualified:        false,
			}, nil)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.CompetitorRole, nil)

		fakedResult := domain.Result{
			Timestamp:    currentTime,
			CompetitorID: fakedCompetitorID,
			Score:        "1000p",
			Placement:    5,
			Finalist:     true,
			RankOrder:    6,
		}

		mockedResultKeeper.On("GetResult", fakedCompetitorID).Return(fakedResult, nil)

		ucase := usecases.CompetitorUseCase{
			Repo:         mockedRepo,
			Authorizer:   mockedAuthorizer,
			ResultKeeper: mockedResultKeeper,
		}

		competitor, err := ucase.PatchCompetitor(context.Background(), fakedCompetitorID, domain.CompetitorPatch{})

		require.NoError(t, err)

		assert.Equal(t, fakedCompetitor.ID, competitor.ID)
		assert.Equal(t, fakedCompetitor.Ownership, competitor.Ownership)
		assert.Equal(t, fakedCompetitor.CompetitionID, competitor.CompetitionID)
		assert.Equal(t, fakedCompetitor.CategoryID, competitor.CategoryID)
		assert.Equal(t, fakedCompetitor.RegistrationCode, competitor.RegistrationCode)
		assert.Equal(t, fakedCompetitor.Name, competitor.Name)
		assert.Equal(t, fakedCompetitor.Entered, competitor.Entered)
		assert.Equal(t, fakedCompetitor.WithdrawnFromFinals, competitor.WithdrawnFromFinals)
		assert.Equal(t, fakedCompetitor.Disqualified, competitor.Disqualified)

		require.NotNil(t, competitor.Result)
		assert.Equal(t, fakedResult, *competitor.Result)

		mockedAuthorizer.AssertExpectations(t)
		mockedResultKeeper.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("CompetitorCannotAlterDisqualifiedState", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)

		fakedCompetitor := domain.Competitor{
			ID:                  fakedCompetitorID,
			Ownership:           fakedOwnership,
			CompetitionID:       fakedCompetitionID,
			CategoryID:          fakedCategoryID,
			RegistrationCode:    "ABCD1234",
			Name:                "John Doe",
			Entered:             currentTime,
			WithdrawnFromFinals: false,
			Disqualified:        true,
		}

		mockedRepo := makeMockedRepo(fakedCompetitor)

		mockedRepo.
			On("GetCategory", mock.Anything, mock.Anything, fakedCategoryID).
			Return(domain.Category{
				ID:        fakedCategoryID,
				TimeBegin: currentTime.Add(-1 * time.Hour),
				TimeEnd:   currentTime.Add(time.Hour),
			}, nil)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.CompetitorRole, nil)

		ucase := usecases.CompetitorUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		competitor, err := ucase.PatchCompetitor(context.Background(), fakedCompetitorID, domain.CompetitorPatch{
			Disqualified: domain.NewPatch(false),
		})

		assert.ErrorIs(t, err, domain.ErrInsufficientRole)
		assert.Empty(t, competitor)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("EnterCompetition", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			mockedAuthorizer := new(authorizerMock)
			mockedResultKeeper := new(resultKeeperMock)
			mockedEventBroker := new(eventBrokerMock)

			fakedCompetitor := domain.Competitor{
				ID:                  fakedCompetitorID,
				Ownership:           fakedOwnership,
				CompetitionID:       fakedCompetitionID,
				CategoryID:          0,
				RegistrationCode:    "ABCD1234",
				Name:                "",
				Entered:             time.Time{},
				WithdrawnFromFinals: false,
				Disqualified:        false,
			}

			mockedRepo := makeMockedRepo(fakedCompetitor)

			mockedRepo.
				On("GetCategory", mock.Anything, mock.Anything, fakedCategoryID).
				Return(domain.Category{
					ID:        fakedCategoryID,
					TimeBegin: currentTime.Add(-1 * time.Hour),
					TimeEnd:   currentTime.Add(time.Hour),
				}, nil)

			mockedRepo.
				On("StoreCompetitor", mock.Anything, mock.Anything,
					domain.Competitor{
						ID:                  fakedCompetitorID,
						Ownership:           fakedOwnership,
						CompetitionID:       fakedCompetitionID,
						CategoryID:          fakedCategoryID,
						RegistrationCode:    "ABCD1234",
						Name:                "John Doe",
						Entered:             time.Now(),
						WithdrawnFromFinals: false,
						Disqualified:        false,
						ScrubBefore:         currentTime.Add(time.Hour).Add(fakedNameRetentionTime),
					},
				).
				Return(domain.Competitor{
					ID:                  fakedCompetitorID,
					Ownership:           fakedOwnership,
					CompetitionID:       fakedCompetitionID,
					CategoryID:          fakedCategoryID,
					RegistrationCode:    "ABCD1234",
					Name:                "John Doe",
					Entered:             time.Now(),
					WithdrawnFromFinals: false,
					Disqualified:        false,
					ScrubBefore:         currentTime.Add(time.Hour).Add(fakedNameRetentionTime),
				}, nil)

			mockedAuthorizer.
				On("HasOwnership", mock.Anything, fakedOwnership).
				Return(domain.CompetitorRole, nil)

			mockedResultKeeper.On("GetResult", fakedCompetitorID).Return(domain.Result{}, errMock)

			mockedEventBroker.On("Dispatch", fakedCompetitionID, mock.Anything).Return()

			ucase := usecases.CompetitorUseCase{
				Repo:         mockedRepo,
				Authorizer:   mockedAuthorizer,
				ResultKeeper: mockedResultKeeper,
				EventBroker:  mockedEventBroker,
			}

			competitor, err := ucase.PatchCompetitor(context.Background(), fakedCompetitorID, domain.CompetitorPatch{
				CategoryID: domain.NewPatch(fakedCategoryID),
				Name:       domain.NewPatch("John Doe"),
			})

			require.NoError(t, err)

			assert.Equal(t, fakedCategoryID, competitor.CategoryID)
			assert.Equal(t, "John Doe", competitor.Name)
			assert.Equal(t, time.Now(), competitor.Entered)
			assert.Equal(t, currentTime.Add(time.Hour).Add(fakedNameRetentionTime), competitor.ScrubBefore)

			mockedEventBroker.AssertCalled(t, "Dispatch", fakedCompetitionID, domain.CompetitorEnteredEvent{
				CompetitorID: fakedCompetitorID,
				CategoryID:   fakedCategoryID,
			})

			mockedEventBroker.AssertCalled(t, "Dispatch", fakedCompetitionID, domain.CompetitorPublicInfoUpdatedEvent{
				CompetitorID:        fakedCompetitorID,
				CategoryID:          fakedCategoryID,
				Name:                "John Doe",
				WithdrawnFromFinals: false,
				Disqualified:        false,
			})

			mockedAuthorizer.AssertExpectations(t)
			mockedResultKeeper.AssertExpectations(t)
			mockedEventBroker.AssertExpectations(t)
			mockedRepo.AssertExpectations(t)
		})
	})

	t.Run("CannotMakeChangesToAnUnregisteredCompetitor", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)

		fakedCompetitor := domain.Competitor{
			ID:                  fakedCompetitorID,
			Ownership:           fakedOwnership,
			CompetitionID:       fakedCompetitionID,
			CategoryID:          0,
			RegistrationCode:    "ABCD1234",
			Name:                "",
			Entered:             time.Time{},
			WithdrawnFromFinals: false,
			Disqualified:        false,
		}

		mockedRepo := makeMockedRepo(fakedCompetitor)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.AdminRole, nil)

		ucase := usecases.CompetitorUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		competitor, err := ucase.PatchCompetitor(context.Background(), fakedCompetitorID, domain.CompetitorPatch{})

		assert.ErrorIs(t, err, domain.ErrNotRegistered)
		assert.Empty(t, competitor)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("CannotLeaveCompetition", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)

		fakedCompetitor := domain.Competitor{
			ID:                  fakedCompetitorID,
			Ownership:           fakedOwnership,
			CompetitionID:       fakedCompetitionID,
			CategoryID:          fakedCategoryID,
			RegistrationCode:    "ABCD1234",
			Name:                "John Doe",
			Entered:             currentTime,
			WithdrawnFromFinals: false,
			Disqualified:        false,
		}

		mockedRepo := makeMockedRepo(fakedCompetitor)

		mockedRepo.
			On("GetCategory", mock.Anything, mock.Anything, fakedCategoryID).
			Return(domain.Category{
				ID:        fakedCategoryID,
				TimeBegin: currentTime.Add(-1 * time.Hour),
				TimeEnd:   currentTime.Add(time.Hour),
			}, nil)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		ucase := usecases.CompetitorUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		competitor, err := ucase.PatchCompetitor(context.Background(), fakedCompetitorID, domain.CompetitorPatch{
			CategoryID: domain.NewPatch(domain.CategoryID(0)),
		})

		assert.ErrorIs(t, err, domain.ErrNotAllowed)
		assert.Empty(t, competitor)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("BatchUpdate", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedResultKeeper := new(resultKeeperMock)
		mockedEventBroker := new(eventBrokerMock)

		fakedCompetitor := domain.Competitor{
			ID:                  fakedCompetitorID,
			Ownership:           fakedOwnership,
			CompetitionID:       fakedCompetitionID,
			CategoryID:          fakedCategoryID,
			RegistrationCode:    "ABCD1234",
			Name:                "John Doe",
			Entered:             currentTime,
			WithdrawnFromFinals: false,
			Disqualified:        false,
			ScrubBefore:         currentTime.Add(42 * time.Hour),
		}

		mockedRepo := makeMockedRepo(fakedCompetitor)

		mockedRepo.
			On("GetCategory", mock.Anything, mock.Anything, fakedCategoryID).
			Return(domain.Category{
				ID:        fakedCategoryID,
				TimeBegin: currentTime.Add(-1 * time.Hour),
				TimeEnd:   currentTime.Add(time.Hour),
			}, nil)

		fakedOtherCategory := domain.Category{
			ID:        testutils.RandomResourceID[domain.CategoryID](),
			TimeBegin: currentTime.Add(-1 * time.Hour),
			TimeEnd:   currentTime,
		}

		mockedRepo.
			On("StoreCompetitor", mock.Anything, mock.Anything,
				domain.Competitor{
					ID:                  fakedCompetitorID,
					Ownership:           fakedOwnership,
					CompetitionID:       fakedCompetitionID,
					CategoryID:          fakedOtherCategory.ID,
					RegistrationCode:    "ABCD1234",
					Name:                "Jane Doe",
					Entered:             currentTime,
					WithdrawnFromFinals: true,
					Disqualified:        true,
					ScrubBefore:         currentTime.Add(42 * time.Hour),
				},
			).
			Return(domain.Competitor{
				ID:                  fakedCompetitorID,
				Ownership:           fakedOwnership,
				CompetitionID:       fakedCompetitionID,
				CategoryID:          fakedOtherCategory.ID,
				RegistrationCode:    "ABCD1234",
				Name:                "Jane Doe",
				Entered:             currentTime,
				WithdrawnFromFinals: true,
				Disqualified:        true,
				ScrubBefore:         currentTime.Add(42 * time.Hour),
			}, nil)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedResultKeeper.On("GetResult", fakedCompetitorID).Return(domain.Result{}, errMock)

		mockedEventBroker.On("Dispatch", fakedCompetitionID, mock.Anything).Return()

		mockedRepo.
			On("GetCategory", mock.Anything, mock.Anything, fakedOtherCategory.ID).
			Return(fakedOtherCategory, nil)

		ucase := usecases.CompetitorUseCase{
			Repo:         mockedRepo,
			Authorizer:   mockedAuthorizer,
			ResultKeeper: mockedResultKeeper,
			EventBroker:  mockedEventBroker,
		}

		competitor, err := ucase.PatchCompetitor(context.Background(), fakedCompetitorID, domain.CompetitorPatch{
			CategoryID:          domain.NewPatch(fakedOtherCategory.ID),
			Name:                domain.NewPatch("Jane Doe"),
			WithdrawnFromFinals: domain.NewPatch(true),
			Disqualified:        domain.NewPatch(true),
		})

		require.NoError(t, err)

		assert.Equal(t, fakedOtherCategory.ID, competitor.CategoryID)
		assert.Equal(t, "Jane Doe", competitor.Name)
		assert.Equal(t, true, competitor.WithdrawnFromFinals)
		assert.Equal(t, true, competitor.Disqualified)
		assert.Equal(t, currentTime, competitor.Entered)
		assert.Equal(t, currentTime.Add(42*time.Hour), competitor.ScrubBefore)

		mockedEventBroker.AssertCalled(t, "Dispatch", fakedCompetitionID, domain.CompetitorSwitchedCategoryEvent{
			CompetitorID: fakedCompetitorID,
			CategoryID:   fakedOtherCategory.ID,
		})

		mockedEventBroker.AssertCalled(t, "Dispatch", fakedCompetitionID, domain.CompetitorPublicInfoUpdatedEvent{
			CompetitorID:        fakedCompetitorID,
			CategoryID:          fakedOtherCategory.ID,
			Name:                "Jane Doe",
			WithdrawnFromFinals: true,
			Disqualified:        true,
		})

		mockedEventBroker.AssertCalled(t, "Dispatch", fakedCompetitionID, domain.CompetitorWithdrewFromFinalsEvent{
			CompetitorID: fakedCompetitorID,
		})

		mockedEventBroker.AssertCalled(t, "Dispatch", fakedCompetitionID, domain.CompetitorDisqualifiedEvent{
			CompetitorID: fakedCompetitorID,
		})

		mockedAuthorizer.AssertExpectations(t)
		mockedResultKeeper.AssertExpectations(t)
		mockedEventBroker.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("NameCannotBeChangedAfterScrubbed", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)

		scrubbedAt := currentTime.Add(-2 * time.Hour)
		scrubBefore := currentTime.Add(42 * time.Hour)

		fakedCompetitor := domain.Competitor{
			ID:                  fakedCompetitorID,
			Ownership:           fakedOwnership,
			CompetitionID:       fakedCompetitionID,
			CategoryID:          fakedCategoryID,
			RegistrationCode:    "ABCD1234",
			Name:                "",
			Entered:             currentTime,
			WithdrawnFromFinals: false,
			Disqualified:        false,
			ScrubbedAt:          scrubbedAt,
			ScrubBefore:         scrubBefore,
		}

		mockedRepo := makeMockedRepo(fakedCompetitor)

		mockedRepo.
			On("GetCategory", mock.Anything, mock.Anything, fakedCategoryID).
			Return(domain.Category{
				ID:        fakedCategoryID,
				TimeBegin: currentTime.Add(-1 * time.Hour),
				TimeEnd:   currentTime.Add(time.Hour),
			}, nil)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		ucase := usecases.CompetitorUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		competitor, err := ucase.PatchCompetitor(context.Background(), fakedCompetitorID, domain.CompetitorPatch{
			Name: domain.NewPatch("John Doe"),
		})

		assert.ErrorIs(t, err, domain.ErrNotAllowed)
		assert.Empty(t, competitor)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("NameCannotBeEmpty", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)

		fakedCompetitor := domain.Competitor{
			ID:                  fakedCompetitorID,
			Ownership:           fakedOwnership,
			CompetitionID:       fakedCompetitionID,
			CategoryID:          fakedCategoryID,
			RegistrationCode:    "ABCD1234",
			Name:                "John Doe",
			Entered:             currentTime,
			WithdrawnFromFinals: false,
			Disqualified:        false,
		}

		mockedRepo := makeMockedRepo(fakedCompetitor)

		mockedRepo.
			On("GetCategory", mock.Anything, mock.Anything, fakedCategoryID).
			Return(domain.Category{
				ID:        fakedCategoryID,
				TimeBegin: currentTime.Add(-1 * time.Hour),
				TimeEnd:   currentTime.Add(time.Hour),
			}, nil)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.CompetitorRole, nil)

		ucase := usecases.CompetitorUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		competitor, err := ucase.PatchCompetitor(context.Background(), fakedCompetitorID, domain.CompetitorPatch{
			Name: domain.NewPatch(string(whitespaceCharacters)),
		})

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.ErrorIs(t, err, domain.ErrEmptyName)
		assert.Empty(t, competitor)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("ReenterFinals", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedResultKeeper := new(resultKeeperMock)
		mockedEventBroker := new(eventBrokerMock)

		fakedCompetitor := domain.Competitor{
			ID:                  fakedCompetitorID,
			Ownership:           fakedOwnership,
			CompetitionID:       fakedCompetitionID,
			CategoryID:          fakedCategoryID,
			RegistrationCode:    "ABCD1234",
			Name:                "John Doe",
			Entered:             currentTime,
			WithdrawnFromFinals: true,
			Disqualified:        false,
		}

		mockedRepo := makeMockedRepo(fakedCompetitor)

		mockedRepo.
			On("GetCategory", mock.Anything, mock.Anything, fakedCategoryID).
			Return(domain.Category{
				ID:        fakedCategoryID,
				TimeBegin: currentTime.Add(-1 * time.Hour),
				TimeEnd:   currentTime.Add(time.Hour),
			}, nil)

		mockedRepo.
			On("StoreCompetitor", mock.Anything, mock.Anything,
				domain.Competitor{
					ID:                  fakedCompetitorID,
					Ownership:           fakedOwnership,
					CompetitionID:       fakedCompetitionID,
					CategoryID:          fakedCategoryID,
					RegistrationCode:    "ABCD1234",
					Name:                "John Doe",
					Entered:             currentTime,
					WithdrawnFromFinals: false,
					Disqualified:        false,
				},
			).
			Return(domain.Competitor{
				ID:                  fakedCompetitorID,
				Ownership:           fakedOwnership,
				CompetitionID:       fakedCompetitionID,
				CategoryID:          fakedCategoryID,
				RegistrationCode:    "ABCD1234",
				Name:                "John Doe",
				Entered:             currentTime,
				WithdrawnFromFinals: false,
				Disqualified:        false,
			}, nil)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.CompetitorRole, nil)

		mockedResultKeeper.On("GetResult", fakedCompetitorID).Return(domain.Result{}, errMock)

		mockedEventBroker.On("Dispatch", fakedCompetitionID, mock.Anything).Return()

		ucase := usecases.CompetitorUseCase{
			Repo:         mockedRepo,
			Authorizer:   mockedAuthorizer,
			ResultKeeper: mockedResultKeeper,
			EventBroker:  mockedEventBroker,
		}

		competitor, err := ucase.PatchCompetitor(context.Background(), fakedCompetitorID, domain.CompetitorPatch{
			WithdrawnFromFinals: domain.NewPatch(false),
		})

		require.NoError(t, err)
		assert.Equal(t, false, competitor.WithdrawnFromFinals)

		mockedEventBroker.AssertCalled(t, "Dispatch", fakedCompetitionID, domain.CompetitorReenteredFinalsEvent{
			CompetitorID: fakedCompetitorID,
		})

		mockedAuthorizer.AssertExpectations(t)
		mockedResultKeeper.AssertExpectations(t)
		mockedEventBroker.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("Requalify", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedResultKeeper := new(resultKeeperMock)
		mockedEventBroker := new(eventBrokerMock)

		fakedCompetitor := domain.Competitor{
			ID:                  fakedCompetitorID,
			Ownership:           fakedOwnership,
			CompetitionID:       fakedCompetitionID,
			CategoryID:          fakedCategoryID,
			RegistrationCode:    "ABCD1234",
			Name:                "John Doe",
			Entered:             currentTime,
			WithdrawnFromFinals: false,
			Disqualified:        true,
		}

		mockedRepo := makeMockedRepo(fakedCompetitor)

		mockedRepo.
			On("GetCategory", mock.Anything, mock.Anything, fakedCategoryID).
			Return(domain.Category{
				ID:        fakedCategoryID,
				TimeBegin: currentTime.Add(-1 * time.Hour),
				TimeEnd:   currentTime.Add(time.Hour),
			}, nil)

		mockedRepo.
			On("StoreCompetitor", mock.Anything, mock.Anything,
				domain.Competitor{
					ID:                  fakedCompetitorID,
					Ownership:           fakedOwnership,
					CompetitionID:       fakedCompetitionID,
					CategoryID:          fakedCategoryID,
					RegistrationCode:    "ABCD1234",
					Name:                "John Doe",
					Entered:             currentTime,
					WithdrawnFromFinals: false,
					Disqualified:        false,
				},
			).
			Return(domain.Competitor{
				ID:                  fakedCompetitorID,
				Ownership:           fakedOwnership,
				CompetitionID:       fakedCompetitionID,
				CategoryID:          fakedCategoryID,
				RegistrationCode:    "ABCD1234",
				Name:                "John Doe",
				Entered:             currentTime,
				WithdrawnFromFinals: false,
				Disqualified:        false,
			}, nil)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedResultKeeper.On("GetResult", fakedCompetitorID).Return(domain.Result{}, errMock)

		mockedEventBroker.On("Dispatch", fakedCompetitionID, mock.Anything).Return()

		ucase := usecases.CompetitorUseCase{
			Repo:         mockedRepo,
			Authorizer:   mockedAuthorizer,
			ResultKeeper: mockedResultKeeper,
			EventBroker:  mockedEventBroker,
		}

		competitor, err := ucase.PatchCompetitor(context.Background(), fakedCompetitorID, domain.CompetitorPatch{
			Disqualified: domain.NewPatch(false),
		})

		require.NoError(t, err)
		assert.Equal(t, false, competitor.Disqualified)

		mockedEventBroker.AssertCalled(t, "Dispatch", fakedCompetitionID, domain.CompetitorRequalifiedEvent{
			CompetitorID: fakedCompetitorID,
		})

		mockedAuthorizer.AssertExpectations(t)
		mockedResultKeeper.AssertExpectations(t)
		mockedEventBroker.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("CannotSwitchToAnEndedCategory", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.CompetitorRole, nil)

		fakedCompetitor := domain.Competitor{
			ID:                  fakedCompetitorID,
			Ownership:           fakedOwnership,
			CompetitionID:       fakedCompetitionID,
			CategoryID:          fakedCategoryID,
			RegistrationCode:    "ABCD1234",
			Name:                "John Doe",
			Entered:             currentTime,
			WithdrawnFromFinals: false,
			Disqualified:        false,
		}

		fakedOtherCategory := domain.Category{
			ID:        testutils.RandomResourceID[domain.CategoryID](),
			TimeBegin: currentTime.Add(-1 * time.Hour),
			TimeEnd:   currentTime.Add(-1 * gracePeriod),
		}

		mockedRepo := makeMockedRepo(fakedCompetitor)

		mockedRepo.
			On("GetCategory", mock.Anything, mock.Anything, fakedCategoryID).
			Return(domain.Category{
				ID:        fakedCategoryID,
				TimeBegin: currentTime.Add(-1 * time.Hour),
				TimeEnd:   currentTime.Add(time.Hour),
			}, nil)

		mockedRepo.
			On("GetCategory", mock.Anything, mock.Anything, fakedOtherCategory.ID).
			Return(fakedOtherCategory, nil)

		ucase := usecases.CompetitorUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		competitor, err := ucase.PatchCompetitor(context.Background(), fakedCompetitorID, domain.CompetitorPatch{
			CategoryID: domain.NewPatch(fakedOtherCategory.ID),
		})

		assert.ErrorIs(t, err, domain.ErrCompetitionEnded)
		assert.Empty(t, competitor)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("CompetitorCannotMakeChangesAfterGracePeriod", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.CompetitorRole, nil)

		fakedOtherCategory := domain.Category{
			ID:        testutils.RandomResourceID[domain.CategoryID](),
			TimeBegin: currentTime.Add(-1 * time.Hour),
			TimeEnd:   currentTime.Add(-1 * gracePeriod),
		}

		fakedCompetitor := domain.Competitor{
			ID:                  fakedCompetitorID,
			Ownership:           fakedOwnership,
			CompetitionID:       fakedCompetitionID,
			CategoryID:          fakedOtherCategory.ID,
			RegistrationCode:    "ABCD1234",
			Name:                "John Doe",
			Entered:             currentTime,
			WithdrawnFromFinals: false,
			Disqualified:        false,
		}

		mockedRepo := makeMockedRepo(fakedCompetitor)

		mockedRepo.
			On("GetCategory", mock.Anything, mock.Anything, fakedOtherCategory.ID).
			Return(fakedOtherCategory, nil)

		ucase := usecases.CompetitorUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		competitor, err := ucase.PatchCompetitor(context.Background(), fakedCompetitorID, domain.CompetitorPatch{})
		assert.ErrorIs(t, err, domain.ErrCompetitionEnded)
		assert.Empty(t, competitor)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("OrganizerCanMakeChangesAfterGracePeriod", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedResultKeeper := new(resultKeeperMock)
		mockedEventBroker := new(eventBrokerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedResultKeeper.On("GetResult", fakedCompetitorID).Return(domain.Result{}, errMock)

		mockedEventBroker.On("Dispatch", fakedCompetitionID, mock.Anything).Return()

		fakedSecondCategory := domain.Category{
			ID:        testutils.RandomResourceID[domain.CategoryID](),
			TimeBegin: currentTime.Add(-1 * time.Hour),
			TimeEnd:   currentTime.Add(-1 * gracePeriod),
		}

		fakedThirdCategory := domain.Category{
			ID:        testutils.RandomResourceID[domain.CategoryID](),
			TimeBegin: currentTime.Add(-1 * time.Hour),
			TimeEnd:   currentTime.Add(-1 * gracePeriod),
		}

		fakedCompetitor := domain.Competitor{
			ID:                  fakedCompetitorID,
			Ownership:           fakedOwnership,
			CompetitionID:       fakedCompetitionID,
			CategoryID:          fakedSecondCategory.ID,
			RegistrationCode:    "ABCD1234",
			Name:                "John Doe",
			Entered:             currentTime,
			WithdrawnFromFinals: false,
			Disqualified:        false,
		}

		mockedRepo := makeMockedRepo(fakedCompetitor)

		mockedRepo.
			On("GetCategory", mock.Anything, mock.Anything, fakedSecondCategory.ID).
			Return(fakedSecondCategory, nil)

		mockedRepo.
			On("GetCategory", mock.Anything, mock.Anything, fakedThirdCategory.ID).
			Return(fakedThirdCategory, nil)

		mockedRepo.
			On("StoreCompetitor", mock.Anything, mock.Anything,
				domain.Competitor{
					ID:                  fakedCompetitorID,
					Ownership:           fakedOwnership,
					CompetitionID:       fakedCompetitionID,
					CategoryID:          fakedThirdCategory.ID,
					RegistrationCode:    "ABCD1234",
					Name:                "John Doe",
					Entered:             currentTime,
					WithdrawnFromFinals: false,
					Disqualified:        false,
				},
			).
			Return(domain.Competitor{
				ID:                  fakedCompetitorID,
				Ownership:           fakedOwnership,
				CompetitionID:       fakedCompetitionID,
				CategoryID:          fakedThirdCategory.ID,
				RegistrationCode:    "ABCD1234",
				Name:                "John Doe",
				Entered:             currentTime,
				WithdrawnFromFinals: false,
				Disqualified:        false,
			}, nil)

		ucase := usecases.CompetitorUseCase{
			Repo:         mockedRepo,
			Authorizer:   mockedAuthorizer,
			ResultKeeper: mockedResultKeeper,
			EventBroker:  mockedEventBroker,
		}

		competitor, err := ucase.PatchCompetitor(context.Background(), fakedCompetitorID, domain.CompetitorPatch{
			CategoryID: domain.NewPatch(fakedThirdCategory.ID),
		})

		require.NoError(t, err)
		assert.Equal(t, fakedThirdCategory.ID, competitor.CategoryID)

		mockedAuthorizer.AssertExpectations(t)
		mockedResultKeeper.AssertExpectations(t)
		mockedEventBroker.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedRepo := new(repositoryMock)

		mockedRepo.
			On("GetCompetitor", mock.Anything, mock.Anything, fakedCompetitorID).
			Return(domain.Competitor{
				ID:        fakedCompetitorID,
				Ownership: fakedOwnership,
			}, nil)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		ucase := usecases.CompetitorUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		competitor, err := ucase.PatchCompetitor(context.Background(), fakedCompetitorID, domain.CompetitorPatch{})

		assert.ErrorIs(t, err, domain.ErrNoOwnership)
		assert.Empty(t, competitor)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})
}

func TestScrubCompetitor(t *testing.T) {
	fakedCompetitorID := testutils.RandomResourceID[domain.CompetitorID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID:  testutils.RandomResourceID[domain.OrganizerID](),
		CompetitorID: &fakedCompetitorID,
	}
	fakedCompetitionID := testutils.RandomCompetitionID()
	fakedCategoryID := testutils.RandomResourceID[domain.CategoryID]()

	t.Run("HappyPath", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			mockedAuthorizer := new(authorizerMock)
			mockedRepo := new(repositoryMock)
			mockedResultKeeper := new(resultKeeperMock)
			mockedEventBroker := new(eventBrokerMock)

			fakedCompetitor := domain.Competitor{
				ID:               fakedCompetitorID,
				Ownership:        fakedOwnership,
				CompetitionID:    fakedCompetitionID,
				CategoryID:       fakedCategoryID,
				RegistrationCode: "ABCD1234",
				Name:             "John Doe",
				Entered:          time.Now().Add(-1 * time.Hour),
			}

			mockedRepo.
				On("GetCompetitor", mock.Anything, mock.Anything, fakedCompetitorID).
				Return(fakedCompetitor, nil)

			mockedAuthorizer.
				On("HasOwnership", mock.Anything, fakedOwnership).
				Return(domain.CompetitorRole, nil)

			mockedRepo.
				On("StoreCompetitor", mock.Anything, mock.Anything, domain.Competitor{
					ID:                  fakedCompetitorID,
					Ownership:           fakedOwnership,
					CompetitionID:       fakedCompetitionID,
					CategoryID:          fakedCategoryID,
					RegistrationCode:    "ABCD1234",
					Name:                "",
					Entered:             fakedCompetitor.Entered,
					WithdrawnFromFinals: true,
					Disqualified:        false,
					ScrubbedAt:          time.Now(),
				}).
				Return(domain.Competitor{
					ID:                  fakedCompetitorID,
					Ownership:           fakedOwnership,
					CompetitionID:       fakedCompetitionID,
					CategoryID:          fakedCategoryID,
					RegistrationCode:    "ABCD1234",
					Name:                "",
					Entered:             fakedCompetitor.Entered,
					WithdrawnFromFinals: true,
					Disqualified:        false,
					ScrubbedAt:          time.Now(),
				}, nil)

			mockedEventBroker.
				On("Dispatch", fakedCompetitionID, domain.CompetitorPublicInfoUpdatedEvent{
					CompetitorID:        fakedCompetitorID,
					CategoryID:          fakedCategoryID,
					Name:                "",
					WithdrawnFromFinals: true,
					Disqualified:        false,
					ScrubbedAt:          time.Now(),
				}).
				Return()

			mockedEventBroker.
				On("Dispatch", fakedCompetitionID, domain.CompetitorWithdrewFromFinalsEvent{
					CompetitorID: fakedCompetitorID,
				}).
				Return()

			mockedResultKeeper.On("GetResult", fakedCompetitorID).Return(domain.Result{}, domain.ErrNotFound)

			ucase := usecases.CompetitorUseCase{
				Repo:         mockedRepo,
				Authorizer:   mockedAuthorizer,
				EventBroker:  mockedEventBroker,
				ResultKeeper: mockedResultKeeper,
			}

			competitor, err := ucase.ScrubCompetitor(context.Background(), fakedCompetitorID)

			require.NoError(t, err)
			assert.Equal(t, fakedCompetitorID, competitor.ID)
			assert.Equal(t, "", competitor.Name)
			assert.True(t, competitor.WithdrawnFromFinals)
			assert.Equal(t, time.Now(), competitor.ScrubbedAt)

			mockedAuthorizer.AssertExpectations(t)
			mockedRepo.AssertExpectations(t)
			mockedResultKeeper.AssertExpectations(t)
			mockedEventBroker.AssertExpectations(t)
		})
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedRepo := new(repositoryMock)

		fakedCompetitor := domain.Competitor{
			ID:        fakedCompetitorID,
			Ownership: fakedOwnership,
		}

		mockedRepo.
			On("GetCompetitor", mock.Anything, mock.Anything, fakedCompetitorID).
			Return(fakedCompetitor, nil)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		ucase := usecases.CompetitorUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		competitor, err := ucase.ScrubCompetitor(context.Background(), fakedCompetitorID)

		assert.ErrorIs(t, err, domain.ErrNoOwnership)
		assert.Empty(t, competitor)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})
}

func TestScrubCompetitors(t *testing.T) {
	fakedCompetitionID := testutils.RandomCompetitionID()
	fakedDeadline := time.Date(2025, 6, 15, 12, 0, 0, 0, time.UTC)

	t.Run("HappyPath", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			mockedRepo := new(repositoryMock)
			mockedEventBroker := new(eventBrokerMock)
			mockedTx := new(transactionMock)

			fakedCompetitors := []domain.Competitor{
				{
					ID:                  testutils.RandomResourceID[domain.CompetitorID](),
					CompetitionID:       fakedCompetitionID,
					CategoryID:          testutils.RandomResourceID[domain.CategoryID](),
					Name:                "Alice",
					WithdrawnFromFinals: false,
					Disqualified:        false,
				},
				{
					ID:                  testutils.RandomResourceID[domain.CompetitorID](),
					CompetitionID:       fakedCompetitionID,
					CategoryID:          testutils.RandomResourceID[domain.CategoryID](),
					Name:                "Bob",
					WithdrawnFromFinals: true,
					Disqualified:        false,
				},
			}

			mockedRepo.
				On("GetScrubEligibleCompetitors", mock.Anything, fakedDeadline).
				Return(fakedCompetitors, nil)

			mockedRepo.On("Begin").Return(mockedTx, nil)

			mockedRepo.
				On("StoreCompetitor", mock.Anything, mockedTx, domain.Competitor{
					ID:            fakedCompetitors[0].ID,
					CompetitionID: fakedCompetitionID,
					CategoryID:    fakedCompetitors[0].CategoryID,
					ScrubbedAt:    time.Now(),
				}).
				Return(domain.Competitor{}, nil).Once()

			mockedRepo.
				On("StoreCompetitor", mock.Anything, mockedTx, domain.Competitor{
					ID:                  fakedCompetitors[1].ID,
					CompetitionID:       fakedCompetitionID,
					CategoryID:          fakedCompetitors[1].CategoryID,
					WithdrawnFromFinals: true,
					ScrubbedAt:          time.Now(),
				}).
				Return(domain.Competitor{}, nil).Once()

			mockedTx.On("Commit").Return(nil)
			mockedTx.On("Rollback").Return()

			mockedEventBroker.
				On("Dispatch", fakedCompetitionID, domain.CompetitorPublicInfoUpdatedEvent{
					CompetitorID: fakedCompetitors[0].ID,
					CategoryID:   fakedCompetitors[0].CategoryID,
					ScrubbedAt:   time.Now(),
				}).Return()

			mockedEventBroker.
				On("Dispatch", fakedCompetitionID, domain.CompetitorPublicInfoUpdatedEvent{
					CompetitorID:        fakedCompetitors[1].ID,
					CategoryID:          fakedCompetitors[1].CategoryID,
					WithdrawnFromFinals: true,
					ScrubbedAt:          time.Now(),
				}).Return()

			ucase := usecases.CompetitorUseCase{
				Repo:        mockedRepo,
				EventBroker: mockedEventBroker,
			}

			count, err := ucase.ScrubCompetitors(context.Background(), fakedDeadline)

			require.NoError(t, err)
			assert.Equal(t, 2, count)

			mockedRepo.AssertExpectations(t)
			mockedEventBroker.AssertExpectations(t)
			mockedTx.AssertExpectations(t)
		})
	})

	t.Run("NoEligibleCompetitors", func(t *testing.T) {
		mockedRepo := new(repositoryMock)

		mockedRepo.
			On("GetScrubEligibleCompetitors", mock.Anything, fakedDeadline).
			Return([]domain.Competitor{}, nil)

		ucase := usecases.CompetitorUseCase{
			Repo: mockedRepo,
		}

		count, err := ucase.ScrubCompetitors(context.Background(), fakedDeadline)

		require.NoError(t, err)
		assert.Equal(t, 0, count)

		mockedRepo.AssertExpectations(t)
	})
}
