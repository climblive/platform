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

func TestGetContender(t *testing.T) {
	fakedContenderID := testutils.RandomResourceID[domain.ContenderID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: testutils.RandomResourceID[domain.OrganizerID](),
		ContenderID: &fakedContenderID,
	}

	t.Run("HappyPath", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedRepo := new(repositoryMock)
		mockedResultKeeper := new(resultKeeperMock)

		fakedResult := domain.Result{
			Timestamp:   time.Now(),
			ContenderID: fakedContenderID,
			Score:       "1000p",
			Placement:   5,
			Finalist:    true,
			RankOrder:   6,
		}

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.ContenderRole, nil)

		mockedRepo.
			On("GetContender", mock.Anything, mock.Anything, fakedContenderID).
			Return(domain.Contender{
				ID:        fakedContenderID,
				Ownership: fakedOwnership,
			}, nil)

		mockedResultKeeper.On("GetResult", fakedContenderID).Return(fakedResult, nil)

		ucase := usecases.ContenderUseCase{
			Repo:         mockedRepo,
			Authorizer:   mockedAuthorizer,
			ResultKeeper: mockedResultKeeper,
		}

		contender, err := ucase.GetContender(context.Background(), fakedContenderID)

		require.NoError(t, err)

		assert.Equal(t, fakedContenderID, contender.ID)
		require.NotNil(t, contender.Result)
		assert.Equal(t, fakedResult, *contender.Result)

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
			On("GetContender", mock.Anything, mock.Anything, fakedContenderID).
			Return(domain.Contender{
				ID:        fakedContenderID,
				Ownership: fakedOwnership,
			}, nil)

		ucase := usecases.ContenderUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		contender, err := ucase.GetContender(context.Background(), fakedContenderID)

		assert.ErrorIs(t, err, domain.ErrNoOwnership)
		assert.Empty(t, contender)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})
}

func TestGetContenderByCode(t *testing.T) {
	fakedContenderID := testutils.RandomResourceID[domain.ContenderID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: testutils.RandomResourceID[domain.OrganizerID](),
		ContenderID: &fakedContenderID,
	}

	t.Run("HappyPath", func(t *testing.T) {
		mockedRepo := new(repositoryMock)
		mockedResultKeeper := new(resultKeeperMock)

		fakedResult := domain.Result{
			Timestamp:   time.Now(),
			ContenderID: fakedContenderID,
			Score:       "1000p",
			Placement:   5,
			Finalist:    true,
			RankOrder:   6,
		}

		mockedRepo.
			On("GetContenderByCode", mock.Anything, mock.Anything, "ABCD1234").
			Return(domain.Contender{
				ID:        fakedContenderID,
				Ownership: fakedOwnership,
			}, nil)

		mockedResultKeeper.On("GetResult", fakedContenderID).Return(fakedResult, nil)

		ucase := usecases.ContenderUseCase{
			Repo:         mockedRepo,
			ResultKeeper: mockedResultKeeper,
		}

		contender, err := ucase.GetContenderByCode(context.Background(), "ABCD1234")

		require.NoError(t, err)
		assert.Equal(t, fakedContenderID, contender.ID)
		require.NotNil(t, contender.Result)
		assert.Equal(t, fakedResult, *contender.Result)

		mockedRepo.AssertExpectations(t)
		mockedResultKeeper.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockedRepo := new(repositoryMock)

		mockedRepo.
			On("GetContenderByCode", mock.Anything, mock.Anything, "DEADBEEF").
			Return(domain.Contender{}, domain.ErrNotFound)

		ucase := usecases.ContenderUseCase{
			Repo: mockedRepo,
		}

		contender, err := ucase.GetContenderByCode(context.Background(), "DEADBEEF")

		assert.ErrorIs(t, err, domain.ErrNotFound)
		assert.Empty(t, contender)

		mockedRepo.AssertExpectations(t)
	})
}

func TestGetPointValues(t *testing.T) {
	fakedContenderID := testutils.RandomResourceID[domain.ContenderID]()

	fakedOwnership := domain.OwnershipData{
		OrganizerID: testutils.RandomResourceID[domain.OrganizerID](),
		ContenderID: &fakedContenderID,
	}

	t.Run("HappyPath", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedRepo := new(repositoryMock)
		mockedPointValueKeeper := new(pointValueKeeperMock)

		fakedPointValues := []domain.PointValue{
			{
				ContenderID: fakedContenderID,
				ProblemID:   testutils.RandomResourceID[domain.ProblemID](),
				Current:     1_000,
				Zone1:       10,
				Zone2:       100,
				Top:         1_000,
				FlashBonus:  5,
			},
			{
				ContenderID: fakedContenderID,
				ProblemID:   testutils.RandomResourceID[domain.ProblemID](),
				Current:     2_000,
				Zone1:       20,
				Zone2:       200,
				Top:         2_000,
				FlashBonus:  10,
			},
			{
				ContenderID: fakedContenderID,
				ProblemID:   testutils.RandomResourceID[domain.ProblemID](),
				Current:     3_000,
				Zone1:       30,
				Zone2:       300,
				Top:         3_000,
				FlashBonus:  15,
			},
		}

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.ContenderRole, nil)

		mockedRepo.
			On("GetContender", mock.Anything, nil, fakedContenderID).
			Return(domain.Contender{
				ID:        fakedContenderID,
				Ownership: fakedOwnership,
			}, nil)

		mockedPointValueKeeper.
			On("GetPointValues", fakedContenderID).
			Return(slices.Collect(slices.Values(fakedPointValues)))

		ucase := usecases.ContenderUseCase{
			Repo:             mockedRepo,
			Authorizer:       mockedAuthorizer,
			PointValueKeeper: mockedPointValueKeeper,
		}

		pointValues, err := ucase.GetPointValues(context.Background(), fakedContenderID)

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
			On("GetContender", mock.Anything, nil, fakedContenderID).
			Return(domain.Contender{
				ID:        fakedContenderID,
				Ownership: fakedOwnership,
			}, nil)

		ucase := usecases.ContenderUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		pointValues, err := ucase.GetPointValues(context.Background(), fakedContenderID)

		assert.ErrorIs(t, err, domain.ErrNoOwnership)
		assert.Empty(t, pointValues)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})
}

func TestGetContendersByCategory(t *testing.T) {
	fakedCategoryID := testutils.RandomResourceID[domain.CategoryID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: testutils.RandomResourceID[domain.OrganizerID](),
	}
	currentTime := time.Now()

	t.Run("HappyPath", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedRepo := new(repositoryMock)
		mockedResultKeeper := new(resultKeeperMock)

		var contenders []domain.Contender

		for k := 1; k <= 10; k++ {
			contenderID := domain.ContenderID(k)

			contenders = append(contenders, domain.Contender{
				ID: contenderID,
			})

			mockedResultKeeper.On("GetResult", contenderID).Return(domain.Result{
				Timestamp:   currentTime,
				ContenderID: contenderID,
				Score:       fmt.Sprintf("%dp", k*10),
				Placement:   k,
				RankOrder:   k - 1,
				Finalist:    true,
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
			On("GetContendersByCategory", mock.Anything, mock.Anything, fakedCategoryID).
			Return(contenders, nil)

		ucase := usecases.ContenderUseCase{
			Repo:         mockedRepo,
			Authorizer:   mockedAuthorizer,
			ResultKeeper: mockedResultKeeper,
		}

		contenders, err := ucase.GetContendersByCategory(context.Background(), fakedCategoryID)

		require.NoError(t, err)
		assert.Len(t, contenders, 10)

		for i, contender := range contenders {
			assert.Equal(t, domain.ContenderID(i+1), contender.ID)
			require.NotNil(t, contender.Result)
			assert.Equal(t, fmt.Sprintf("%dp", (i+1)*10), contender.Result.Score)
			assert.Equal(t, i+1, contender.Result.Placement)
			assert.Equal(t, i, contender.Result.RankOrder)
			assert.True(t, contender.Result.Finalist)
			assert.Equal(t, currentTime, contender.Result.Timestamp)
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

		ucase := usecases.ContenderUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		contenders, err := ucase.GetContendersByCategory(context.Background(), fakedCategoryID)

		assert.ErrorIs(t, err, domain.ErrNoOwnership)
		assert.Nil(t, contenders)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})
}

func TestGetContendersByCompetition(t *testing.T) {
	fakedCompetitionID := testutils.RandomResourceID[domain.CompetitionID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: testutils.RandomResourceID[domain.OrganizerID](),
	}
	currentTime := time.Now()

	t.Run("HappyPath", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedRepo := new(repositoryMock)
		mockedResultKeeper := new(resultKeeperMock)

		var contenders []domain.Contender

		for k := 1; k <= 10; k++ {
			contenderID := domain.ContenderID(k)

			contenders = append(contenders, domain.Contender{
				ID: contenderID,
			})

			mockedResultKeeper.On("GetResult", contenderID).Return(domain.Result{
				Timestamp:   currentTime,
				ContenderID: contenderID,
				Score:       fmt.Sprintf("%dp", k*10),
				Placement:   k,
				RankOrder:   k - 1,
				Finalist:    true,
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
			On("GetContendersByCompetition", mock.Anything, mock.Anything, fakedCompetitionID).
			Return(contenders, nil)

		ucase := usecases.ContenderUseCase{
			Repo:         mockedRepo,
			Authorizer:   mockedAuthorizer,
			ResultKeeper: mockedResultKeeper,
		}

		contenders, err := ucase.GetContendersByCompetition(context.Background(), fakedCompetitionID)

		require.NoError(t, err)
		assert.Len(t, contenders, 10)

		for i, contender := range contenders {
			assert.Equal(t, domain.ContenderID(i+1), contender.ID)
			require.NotNil(t, contender.Result)
			assert.Equal(t, fmt.Sprintf("%dp", (i+1)*10), contender.Result.Score)
			assert.Equal(t, i+1, contender.Result.Placement)
			assert.Equal(t, i, contender.Result.RankOrder)
			assert.True(t, contender.Result.Finalist)
			assert.Equal(t, currentTime, contender.Result.Timestamp)
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

		ucase := usecases.ContenderUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		contenders, err := ucase.GetContendersByCompetition(context.Background(), fakedCompetitionID)

		assert.ErrorIs(t, err, domain.ErrNoOwnership)
		assert.Nil(t, contenders)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})
}

func TestDeleteContender(t *testing.T) {
	fakedContenderID := testutils.RandomResourceID[domain.ContenderID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: testutils.RandomResourceID[domain.OrganizerID](),
		ContenderID: &fakedContenderID,
	}

	t.Run("HappyPath", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedRepo := new(repositoryMock)

		mockedRepo.
			On("GetContender", mock.Anything, mock.Anything, fakedContenderID).
			Return(domain.Contender{Ownership: fakedOwnership}, nil)

		mockedRepo.
			On("DeleteContender", mock.Anything, mock.Anything, fakedContenderID).
			Return(nil)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		ucase := usecases.ContenderUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		err := ucase.DeleteContender(context.Background(), fakedContenderID)

		require.NoError(t, err)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedRepo := new(repositoryMock)

		mockedRepo.
			On("GetContender", mock.Anything, mock.Anything, fakedContenderID).
			Return(domain.Contender{Ownership: fakedOwnership}, nil)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		ucase := usecases.ContenderUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		err := ucase.DeleteContender(context.Background(), fakedContenderID)

		assert.ErrorIs(t, err, domain.ErrNoOwnership)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("InsufficientRole", func(t *testing.T) {
		for _, insufficientRole := range []domain.AuthRole{domain.NilRole, domain.ContenderRole, domain.JudgeRole} {
			mockedAuthorizer := new(authorizerMock)
			mockedRepo := new(repositoryMock)

			mockedAuthorizer.
				On("HasOwnership", mock.Anything, fakedOwnership).
				Return(insufficientRole, nil)

			mockedRepo.
				On("GetContender", mock.Anything, mock.Anything, fakedContenderID).
				Return(domain.Contender{Ownership: fakedOwnership}, nil)

			ucase := usecases.ContenderUseCase{
				Repo:       mockedRepo,
				Authorizer: mockedAuthorizer,
			}

			err := ucase.DeleteContender(context.Background(), fakedContenderID)

			assert.ErrorIs(t, err, domain.ErrInsufficientRole)

			mockedAuthorizer.AssertExpectations(t)
			mockedRepo.AssertExpectations(t)
		}
	})
}

func TestCreateContenders(t *testing.T) {
	fakedCompetitionID := testutils.RandomResourceID[domain.CompetitionID]()
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
			On("GetNumberOfContenders", mock.Anything, mock.Anything, fakedCompetitionID).
			Return(400, nil)

		for n := range 100 {
			code := fmt.Sprintf("%08d", n)
			contender := domain.Contender{
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
				On("StoreContender", mock.Anything, mock.Anything, contender).
				Return(contender, nil)
		}

		mockedRepo.
			On("Begin").
			Return(mockedTx, nil)

		mockedTx.On("Commit").Return(nil)

		ucase := usecases.ContenderUseCase{
			Repo:                      mockedRepo,
			Authorizer:                mockedAuthorizer,
			RegistrationCodeGenerator: mockedCodeGenerator,
		}

		contenders, err := ucase.CreateContenders(context.Background(), fakedCompetitionID, 100)

		require.NoError(t, err)
		assert.Len(t, contenders, 100)

		for i, contender := range contenders {
			assert.Equal(t, fmt.Sprintf("%08d", i), contender.RegistrationCode)
		}

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
		mockedTx.AssertExpectations(t)
		mockedCodeGenerator.AssertExpectations(t)
	})

	t.Run("CannotExceed500Contenders", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedRepo, _, mockedCodeGenerator := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("GetNumberOfContenders", mock.Anything, mock.Anything, fakedCompetitionID).
			Return(400, nil)

		ucase := usecases.ContenderUseCase{
			Repo:                      mockedRepo,
			Authorizer:                mockedAuthorizer,
			RegistrationCodeGenerator: mockedCodeGenerator,
		}

		contenders, err := ucase.CreateContenders(context.Background(), fakedCompetitionID, 101)

		assert.ErrorIs(t, err, domain.ErrLimitExceeded)
		assert.Nil(t, contenders)

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

		ucase := usecases.ContenderUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		contenders, err := ucase.CreateContenders(context.Background(), fakedCompetitionID, 100)

		assert.ErrorIs(t, err, domain.ErrNoOwnership)
		assert.Empty(t, contenders)

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
			On("GetNumberOfContenders", mock.Anything, mock.Anything, fakedCompetitionID).
			Return(0, nil)

		mockedRepo.
			On("StoreContender", mock.Anything, mock.Anything, mock.Anything).
			Return(domain.Contender{}, errMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.AdminRole, nil)

		mockedCodeGenerator.
			On("Generate", 8).
			Return("ABCD1234")

		ucase := usecases.ContenderUseCase{
			Repo:                      mockedRepo,
			Authorizer:                mockedAuthorizer,
			RegistrationCodeGenerator: mockedCodeGenerator,
		}

		contenders, err := ucase.CreateContenders(context.Background(), fakedCompetitionID, 1)

		assert.Error(t, err)
		assert.Nil(t, contenders)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
		mockedTx.AssertExpectations(t)
		mockedCodeGenerator.AssertExpectations(t)
	})
}

func TestPatchContender(t *testing.T) {
	fakedContenderID := testutils.RandomResourceID[domain.ContenderID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: testutils.RandomResourceID[domain.OrganizerID](),
		ContenderID: &fakedContenderID,
	}
	fakedCompetitionID := testutils.RandomResourceID[domain.CompetitionID]()
	fakedCategoryID := testutils.RandomResourceID[domain.CategoryID]()

	currentTime := time.Now()
	gracePeriod := 15 * time.Minute
	fakedNameRetentionTime := 14 * 24 * time.Hour

	makeMockedRepo := func(contender domain.Contender) *repositoryMock {
		fakedCompetition := domain.Competition{
			ID:                fakedCompetitionID,
			GracePeriod:       gracePeriod,
			NameRetentionTime: fakedNameRetentionTime,
		}

		mockedRepo := new(repositoryMock)

		mockedRepo.
			On("GetContender", mock.Anything, mock.Anything, contender.ID).
			Return(contender, nil)

		mockedRepo.
			On("GetCompetition", mock.Anything, mock.Anything, fakedCompetitionID).
			Return(fakedCompetition, nil)

		return mockedRepo
	}

	t.Run("UpdateWithoutChanges", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedResultKeeper := new(resultKeeperMock)

		fakedContender := domain.Contender{
			ID:                  fakedContenderID,
			Ownership:           fakedOwnership,
			CompetitionID:       fakedCompetitionID,
			CategoryID:          fakedCategoryID,
			RegistrationCode:    "ABCD1234",
			Name:                "John Doe",
			Entered:             currentTime,
			WithdrawnFromFinals: false,
			Disqualified:        false,
		}

		mockedRepo := makeMockedRepo(fakedContender)

		mockedRepo.
			On("GetCategory", mock.Anything, mock.Anything, fakedCategoryID).
			Return(domain.Category{
				ID:        fakedCategoryID,
				TimeBegin: currentTime.Add(-1 * time.Hour),
				TimeEnd:   currentTime.Add(time.Hour),
			}, nil)

		mockedRepo.
			On("StoreContender", mock.Anything, mock.Anything,
				domain.Contender{
					ID:                  fakedContenderID,
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
			Return(domain.Contender{
				ID:                  fakedContenderID,
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
			Return(domain.ContenderRole, nil)

		fakedResult := domain.Result{
			Timestamp:   currentTime,
			ContenderID: fakedContenderID,
			Score:       "1000p",
			Placement:   5,
			Finalist:    true,
			RankOrder:   6,
		}

		mockedResultKeeper.On("GetResult", fakedContenderID).Return(fakedResult, nil)

		ucase := usecases.ContenderUseCase{
			Repo:         mockedRepo,
			Authorizer:   mockedAuthorizer,
			ResultKeeper: mockedResultKeeper,
		}

		contender, err := ucase.PatchContender(context.Background(), fakedContenderID, domain.ContenderPatch{})

		require.NoError(t, err)

		assert.Equal(t, fakedContender.ID, contender.ID)
		assert.Equal(t, fakedContender.Ownership, contender.Ownership)
		assert.Equal(t, fakedContender.CompetitionID, contender.CompetitionID)
		assert.Equal(t, fakedContender.CategoryID, contender.CategoryID)
		assert.Equal(t, fakedContender.RegistrationCode, contender.RegistrationCode)
		assert.Equal(t, fakedContender.Name, contender.Name)
		assert.Equal(t, fakedContender.Entered, contender.Entered)
		assert.Equal(t, fakedContender.WithdrawnFromFinals, contender.WithdrawnFromFinals)
		assert.Equal(t, fakedContender.Disqualified, contender.Disqualified)

		require.NotNil(t, contender.Result)
		assert.Equal(t, fakedResult, *contender.Result)

		mockedAuthorizer.AssertExpectations(t)
		mockedResultKeeper.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("ContenderCannotAlterDisqualifiedState", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)

		fakedContender := domain.Contender{
			ID:                  fakedContenderID,
			Ownership:           fakedOwnership,
			CompetitionID:       fakedCompetitionID,
			CategoryID:          fakedCategoryID,
			RegistrationCode:    "ABCD1234",
			Name:                "John Doe",
			Entered:             currentTime,
			WithdrawnFromFinals: false,
			Disqualified:        true,
		}

		mockedRepo := makeMockedRepo(fakedContender)

		mockedRepo.
			On("GetCategory", mock.Anything, mock.Anything, fakedCategoryID).
			Return(domain.Category{
				ID:        fakedCategoryID,
				TimeBegin: currentTime.Add(-1 * time.Hour),
				TimeEnd:   currentTime.Add(time.Hour),
			}, nil)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.ContenderRole, nil)

		ucase := usecases.ContenderUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		contender, err := ucase.PatchContender(context.Background(), fakedContenderID, domain.ContenderPatch{
			Disqualified: domain.NewPatch(false),
		})

		assert.ErrorIs(t, err, domain.ErrInsufficientRole)
		assert.Empty(t, contender)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("EnterCompetition", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			mockedAuthorizer := new(authorizerMock)
			mockedResultKeeper := new(resultKeeperMock)
			mockedEventBroker := new(eventBrokerMock)

			fakedContender := domain.Contender{
				ID:                  fakedContenderID,
				Ownership:           fakedOwnership,
				CompetitionID:       fakedCompetitionID,
				CategoryID:          0,
				RegistrationCode:    "ABCD1234",
				Name:                "",
				Entered:             time.Time{},
				WithdrawnFromFinals: false,
				Disqualified:        false,
			}

			mockedRepo := makeMockedRepo(fakedContender)

			mockedRepo.
				On("GetCategory", mock.Anything, mock.Anything, fakedCategoryID).
				Return(domain.Category{
					ID:        fakedCategoryID,
					TimeBegin: currentTime.Add(-1 * time.Hour),
					TimeEnd:   currentTime.Add(time.Hour),
				}, nil)

			mockedRepo.
				On("StoreContender", mock.Anything, mock.Anything,
					domain.Contender{
						ID:                  fakedContenderID,
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
				Return(domain.Contender{
					ID:                  fakedContenderID,
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
				Return(domain.ContenderRole, nil)

			mockedResultKeeper.On("GetResult", fakedContenderID).Return(domain.Result{}, errMock)

			mockedEventBroker.On("Dispatch", fakedCompetitionID, mock.Anything).Return()

			ucase := usecases.ContenderUseCase{
				Repo:         mockedRepo,
				Authorizer:   mockedAuthorizer,
				ResultKeeper: mockedResultKeeper,
				EventBroker:  mockedEventBroker,
			}

			contender, err := ucase.PatchContender(context.Background(), fakedContenderID, domain.ContenderPatch{
				CategoryID: domain.NewPatch(fakedCategoryID),
				Name:       domain.NewPatch("John Doe"),
			})

			require.NoError(t, err)

			assert.Equal(t, fakedCategoryID, contender.CategoryID)
			assert.Equal(t, "John Doe", contender.Name)
			assert.Equal(t, time.Now(), contender.Entered)
			assert.Equal(t, currentTime.Add(time.Hour).Add(fakedNameRetentionTime), contender.ScrubBefore)

			mockedEventBroker.AssertCalled(t, "Dispatch", fakedCompetitionID, domain.ContenderEnteredEvent{
				ContenderID: fakedContenderID,
				CategoryID:  fakedCategoryID,
			})

			mockedEventBroker.AssertCalled(t, "Dispatch", fakedCompetitionID, domain.ContenderPublicInfoUpdatedEvent{
				ContenderID:         fakedContenderID,
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

	t.Run("CannotMakeChangesToAnUnregisteredContender", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)

		fakedContender := domain.Contender{
			ID:                  fakedContenderID,
			Ownership:           fakedOwnership,
			CompetitionID:       fakedCompetitionID,
			CategoryID:          0,
			RegistrationCode:    "ABCD1234",
			Name:                "",
			Entered:             time.Time{},
			WithdrawnFromFinals: false,
			Disqualified:        false,
		}

		mockedRepo := makeMockedRepo(fakedContender)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.AdminRole, nil)

		ucase := usecases.ContenderUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		contender, err := ucase.PatchContender(context.Background(), fakedContenderID, domain.ContenderPatch{})

		assert.ErrorIs(t, err, domain.ErrNotRegistered)
		assert.Empty(t, contender)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("CannotLeaveCompetition", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)

		fakedContender := domain.Contender{
			ID:                  fakedContenderID,
			Ownership:           fakedOwnership,
			CompetitionID:       fakedCompetitionID,
			CategoryID:          fakedCategoryID,
			RegistrationCode:    "ABCD1234",
			Name:                "John Doe",
			Entered:             currentTime,
			WithdrawnFromFinals: false,
			Disqualified:        false,
		}

		mockedRepo := makeMockedRepo(fakedContender)

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

		ucase := usecases.ContenderUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		contender, err := ucase.PatchContender(context.Background(), fakedContenderID, domain.ContenderPatch{
			CategoryID: domain.NewPatch(domain.CategoryID(0)),
		})

		assert.ErrorIs(t, err, domain.ErrNotAllowed)
		assert.Empty(t, contender)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("BatchUpdate", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedResultKeeper := new(resultKeeperMock)
		mockedEventBroker := new(eventBrokerMock)

		fakedContender := domain.Contender{
			ID:                  fakedContenderID,
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

		mockedRepo := makeMockedRepo(fakedContender)

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
			On("StoreContender", mock.Anything, mock.Anything,
				domain.Contender{
					ID:                  fakedContenderID,
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
			Return(domain.Contender{
				ID:                  fakedContenderID,
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

		mockedResultKeeper.On("GetResult", fakedContenderID).Return(domain.Result{}, errMock)

		mockedEventBroker.On("Dispatch", fakedCompetitionID, mock.Anything).Return()

		mockedRepo.
			On("GetCategory", mock.Anything, mock.Anything, fakedOtherCategory.ID).
			Return(fakedOtherCategory, nil)

		ucase := usecases.ContenderUseCase{
			Repo:         mockedRepo,
			Authorizer:   mockedAuthorizer,
			ResultKeeper: mockedResultKeeper,
			EventBroker:  mockedEventBroker,
		}

		contender, err := ucase.PatchContender(context.Background(), fakedContenderID, domain.ContenderPatch{
			CategoryID:          domain.NewPatch(fakedOtherCategory.ID),
			Name:                domain.NewPatch("Jane Doe"),
			WithdrawnFromFinals: domain.NewPatch(true),
			Disqualified:        domain.NewPatch(true),
		})

		require.NoError(t, err)

		assert.Equal(t, fakedOtherCategory.ID, contender.CategoryID)
		assert.Equal(t, "Jane Doe", contender.Name)
		assert.Equal(t, true, contender.WithdrawnFromFinals)
		assert.Equal(t, true, contender.Disqualified)
		assert.Equal(t, currentTime, contender.Entered)
		assert.Equal(t, currentTime.Add(42*time.Hour), contender.ScrubBefore)

		mockedEventBroker.AssertCalled(t, "Dispatch", fakedCompetitionID, domain.ContenderSwitchedCategoryEvent{
			ContenderID: fakedContenderID,
			CategoryID:  fakedOtherCategory.ID,
		})

		mockedEventBroker.AssertCalled(t, "Dispatch", fakedCompetitionID, domain.ContenderPublicInfoUpdatedEvent{
			ContenderID:         fakedContenderID,
			CategoryID:          fakedOtherCategory.ID,
			Name:                "Jane Doe",
			WithdrawnFromFinals: true,
			Disqualified:        true,
		})

		mockedEventBroker.AssertCalled(t, "Dispatch", fakedCompetitionID, domain.ContenderWithdrewFromFinalsEvent{
			ContenderID: fakedContenderID,
		})

		mockedEventBroker.AssertCalled(t, "Dispatch", fakedCompetitionID, domain.ContenderDisqualifiedEvent{
			ContenderID: fakedContenderID,
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

		fakedContender := domain.Contender{
			ID:                  fakedContenderID,
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

		mockedRepo := makeMockedRepo(fakedContender)

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

		ucase := usecases.ContenderUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		contender, err := ucase.PatchContender(context.Background(), fakedContenderID, domain.ContenderPatch{
			Name: domain.NewPatch("John Doe"),
		})

		assert.ErrorIs(t, err, domain.ErrNotAllowed)
		assert.Empty(t, contender)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("NameCannotBeEmpty", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)

		fakedContender := domain.Contender{
			ID:                  fakedContenderID,
			Ownership:           fakedOwnership,
			CompetitionID:       fakedCompetitionID,
			CategoryID:          fakedCategoryID,
			RegistrationCode:    "ABCD1234",
			Name:                "John Doe",
			Entered:             currentTime,
			WithdrawnFromFinals: false,
			Disqualified:        false,
		}

		mockedRepo := makeMockedRepo(fakedContender)

		mockedRepo.
			On("GetCategory", mock.Anything, mock.Anything, fakedCategoryID).
			Return(domain.Category{
				ID:        fakedCategoryID,
				TimeBegin: currentTime.Add(-1 * time.Hour),
				TimeEnd:   currentTime.Add(time.Hour),
			}, nil)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.ContenderRole, nil)

		ucase := usecases.ContenderUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		contender, err := ucase.PatchContender(context.Background(), fakedContenderID, domain.ContenderPatch{
			Name: domain.NewPatch(string(whitespaceCharacters)),
		})

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.ErrorIs(t, err, domain.ErrEmptyName)
		assert.Empty(t, contender)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("ReenterFinals", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedResultKeeper := new(resultKeeperMock)
		mockedEventBroker := new(eventBrokerMock)

		fakedContender := domain.Contender{
			ID:                  fakedContenderID,
			Ownership:           fakedOwnership,
			CompetitionID:       fakedCompetitionID,
			CategoryID:          fakedCategoryID,
			RegistrationCode:    "ABCD1234",
			Name:                "John Doe",
			Entered:             currentTime,
			WithdrawnFromFinals: true,
			Disqualified:        false,
		}

		mockedRepo := makeMockedRepo(fakedContender)

		mockedRepo.
			On("GetCategory", mock.Anything, mock.Anything, fakedCategoryID).
			Return(domain.Category{
				ID:        fakedCategoryID,
				TimeBegin: currentTime.Add(-1 * time.Hour),
				TimeEnd:   currentTime.Add(time.Hour),
			}, nil)

		mockedRepo.
			On("StoreContender", mock.Anything, mock.Anything,
				domain.Contender{
					ID:                  fakedContenderID,
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
			Return(domain.Contender{
				ID:                  fakedContenderID,
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
			Return(domain.ContenderRole, nil)

		mockedResultKeeper.On("GetResult", fakedContenderID).Return(domain.Result{}, errMock)

		mockedEventBroker.On("Dispatch", fakedCompetitionID, mock.Anything).Return()

		ucase := usecases.ContenderUseCase{
			Repo:         mockedRepo,
			Authorizer:   mockedAuthorizer,
			ResultKeeper: mockedResultKeeper,
			EventBroker:  mockedEventBroker,
		}

		contender, err := ucase.PatchContender(context.Background(), fakedContenderID, domain.ContenderPatch{
			WithdrawnFromFinals: domain.NewPatch(false),
		})

		require.NoError(t, err)
		assert.Equal(t, false, contender.WithdrawnFromFinals)

		mockedEventBroker.AssertCalled(t, "Dispatch", fakedCompetitionID, domain.ContenderReenteredFinalsEvent{
			ContenderID: fakedContenderID,
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

		fakedContender := domain.Contender{
			ID:                  fakedContenderID,
			Ownership:           fakedOwnership,
			CompetitionID:       fakedCompetitionID,
			CategoryID:          fakedCategoryID,
			RegistrationCode:    "ABCD1234",
			Name:                "John Doe",
			Entered:             currentTime,
			WithdrawnFromFinals: false,
			Disqualified:        true,
		}

		mockedRepo := makeMockedRepo(fakedContender)

		mockedRepo.
			On("GetCategory", mock.Anything, mock.Anything, fakedCategoryID).
			Return(domain.Category{
				ID:        fakedCategoryID,
				TimeBegin: currentTime.Add(-1 * time.Hour),
				TimeEnd:   currentTime.Add(time.Hour),
			}, nil)

		mockedRepo.
			On("StoreContender", mock.Anything, mock.Anything,
				domain.Contender{
					ID:                  fakedContenderID,
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
			Return(domain.Contender{
				ID:                  fakedContenderID,
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

		mockedResultKeeper.On("GetResult", fakedContenderID).Return(domain.Result{}, errMock)

		mockedEventBroker.On("Dispatch", fakedCompetitionID, mock.Anything).Return()

		ucase := usecases.ContenderUseCase{
			Repo:         mockedRepo,
			Authorizer:   mockedAuthorizer,
			ResultKeeper: mockedResultKeeper,
			EventBroker:  mockedEventBroker,
		}

		contender, err := ucase.PatchContender(context.Background(), fakedContenderID, domain.ContenderPatch{
			Disqualified: domain.NewPatch(false),
		})

		require.NoError(t, err)
		assert.Equal(t, false, contender.Disqualified)

		mockedEventBroker.AssertCalled(t, "Dispatch", fakedCompetitionID, domain.ContenderRequalifiedEvent{
			ContenderID: fakedContenderID,
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
			Return(domain.ContenderRole, nil)

		fakedContender := domain.Contender{
			ID:                  fakedContenderID,
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

		mockedRepo := makeMockedRepo(fakedContender)

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

		ucase := usecases.ContenderUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		contender, err := ucase.PatchContender(context.Background(), fakedContenderID, domain.ContenderPatch{
			CategoryID: domain.NewPatch(fakedOtherCategory.ID),
		})

		assert.ErrorIs(t, err, domain.ErrCompetitionEnded)
		assert.Empty(t, contender)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("ContenderCannotMakeChangesAfterGracePeriod", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.ContenderRole, nil)

		fakedOtherCategory := domain.Category{
			ID:        testutils.RandomResourceID[domain.CategoryID](),
			TimeBegin: currentTime.Add(-1 * time.Hour),
			TimeEnd:   currentTime.Add(-1 * gracePeriod),
		}

		fakedContender := domain.Contender{
			ID:                  fakedContenderID,
			Ownership:           fakedOwnership,
			CompetitionID:       fakedCompetitionID,
			CategoryID:          fakedOtherCategory.ID,
			RegistrationCode:    "ABCD1234",
			Name:                "John Doe",
			Entered:             currentTime,
			WithdrawnFromFinals: false,
			Disqualified:        false,
		}

		mockedRepo := makeMockedRepo(fakedContender)

		mockedRepo.
			On("GetCategory", mock.Anything, mock.Anything, fakedOtherCategory.ID).
			Return(fakedOtherCategory, nil)

		ucase := usecases.ContenderUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		contender, err := ucase.PatchContender(context.Background(), fakedContenderID, domain.ContenderPatch{})
		assert.ErrorIs(t, err, domain.ErrCompetitionEnded)
		assert.Empty(t, contender)

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

		mockedResultKeeper.On("GetResult", fakedContenderID).Return(domain.Result{}, errMock)

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

		fakedContender := domain.Contender{
			ID:                  fakedContenderID,
			Ownership:           fakedOwnership,
			CompetitionID:       fakedCompetitionID,
			CategoryID:          fakedSecondCategory.ID,
			RegistrationCode:    "ABCD1234",
			Name:                "John Doe",
			Entered:             currentTime,
			WithdrawnFromFinals: false,
			Disqualified:        false,
		}

		mockedRepo := makeMockedRepo(fakedContender)

		mockedRepo.
			On("GetCategory", mock.Anything, mock.Anything, fakedSecondCategory.ID).
			Return(fakedSecondCategory, nil)

		mockedRepo.
			On("GetCategory", mock.Anything, mock.Anything, fakedThirdCategory.ID).
			Return(fakedThirdCategory, nil)

		mockedRepo.
			On("StoreContender", mock.Anything, mock.Anything,
				domain.Contender{
					ID:                  fakedContenderID,
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
			Return(domain.Contender{
				ID:                  fakedContenderID,
				Ownership:           fakedOwnership,
				CompetitionID:       fakedCompetitionID,
				CategoryID:          fakedThirdCategory.ID,
				RegistrationCode:    "ABCD1234",
				Name:                "John Doe",
				Entered:             currentTime,
				WithdrawnFromFinals: false,
				Disqualified:        false,
			}, nil)

		ucase := usecases.ContenderUseCase{
			Repo:         mockedRepo,
			Authorizer:   mockedAuthorizer,
			ResultKeeper: mockedResultKeeper,
			EventBroker:  mockedEventBroker,
		}

		contender, err := ucase.PatchContender(context.Background(), fakedContenderID, domain.ContenderPatch{
			CategoryID: domain.NewPatch(fakedThirdCategory.ID),
		})

		require.NoError(t, err)
		assert.Equal(t, fakedThirdCategory.ID, contender.CategoryID)

		mockedAuthorizer.AssertExpectations(t)
		mockedResultKeeper.AssertExpectations(t)
		mockedEventBroker.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedRepo := new(repositoryMock)

		mockedRepo.
			On("GetContender", mock.Anything, mock.Anything, fakedContenderID).
			Return(domain.Contender{
				ID:        fakedContenderID,
				Ownership: fakedOwnership,
			}, nil)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		ucase := usecases.ContenderUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		contender, err := ucase.PatchContender(context.Background(), fakedContenderID, domain.ContenderPatch{})

		assert.ErrorIs(t, err, domain.ErrNoOwnership)
		assert.Empty(t, contender)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})
}

func TestScrubContender(t *testing.T) {
	fakedContenderID := testutils.RandomResourceID[domain.ContenderID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: testutils.RandomResourceID[domain.OrganizerID](),
		ContenderID: &fakedContenderID,
	}
	fakedCompetitionID := testutils.RandomResourceID[domain.CompetitionID]()
	fakedCategoryID := testutils.RandomResourceID[domain.CategoryID]()

	t.Run("HappyPath", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			mockedAuthorizer := new(authorizerMock)
			mockedRepo := new(repositoryMock)
			mockedResultKeeper := new(resultKeeperMock)
			mockedEventBroker := new(eventBrokerMock)

			fakedContender := domain.Contender{
				ID:               fakedContenderID,
				Ownership:        fakedOwnership,
				CompetitionID:    fakedCompetitionID,
				CategoryID:       fakedCategoryID,
				RegistrationCode: "ABCD1234",
				Name:             "John Doe",
				Entered:          time.Now().Add(-1 * time.Hour),
			}

			mockedRepo.
				On("GetContender", mock.Anything, mock.Anything, fakedContenderID).
				Return(fakedContender, nil)

			mockedAuthorizer.
				On("HasOwnership", mock.Anything, fakedOwnership).
				Return(domain.ContenderRole, nil)

			mockedRepo.
				On("StoreContender", mock.Anything, mock.Anything, domain.Contender{
					ID:                  fakedContenderID,
					Ownership:           fakedOwnership,
					CompetitionID:       fakedCompetitionID,
					CategoryID:          fakedCategoryID,
					RegistrationCode:    "ABCD1234",
					Name:                "",
					Entered:             fakedContender.Entered,
					WithdrawnFromFinals: true,
					Disqualified:        false,
					ScrubbedAt:          time.Now(),
				}).
				Return(domain.Contender{
					ID:                  fakedContenderID,
					Ownership:           fakedOwnership,
					CompetitionID:       fakedCompetitionID,
					CategoryID:          fakedCategoryID,
					RegistrationCode:    "ABCD1234",
					Name:                "",
					Entered:             fakedContender.Entered,
					WithdrawnFromFinals: true,
					Disqualified:        false,
					ScrubbedAt:          time.Now(),
				}, nil)

			mockedEventBroker.
				On("Dispatch", fakedCompetitionID, domain.ContenderPublicInfoUpdatedEvent{
					ContenderID:         fakedContenderID,
					CategoryID:          fakedCategoryID,
					Name:                "",
					WithdrawnFromFinals: true,
					Disqualified:        false,
					ScrubbedAt:          time.Now(),
				}).
				Return()

			mockedEventBroker.
				On("Dispatch", fakedCompetitionID, domain.ContenderWithdrewFromFinalsEvent{
					ContenderID: fakedContenderID,
				}).
				Return()

			mockedResultKeeper.On("GetResult", fakedContenderID).Return(domain.Result{}, domain.ErrNotFound)

			ucase := usecases.ContenderUseCase{
				Repo:         mockedRepo,
				Authorizer:   mockedAuthorizer,
				EventBroker:  mockedEventBroker,
				ResultKeeper: mockedResultKeeper,
			}

			contender, err := ucase.ScrubContender(context.Background(), fakedContenderID)

			require.NoError(t, err)
			assert.Equal(t, fakedContenderID, contender.ID)
			assert.Equal(t, "", contender.Name)
			assert.True(t, contender.WithdrawnFromFinals)
			assert.Equal(t, time.Now(), contender.ScrubbedAt)

			mockedAuthorizer.AssertExpectations(t)
			mockedRepo.AssertExpectations(t)
			mockedResultKeeper.AssertExpectations(t)
			mockedEventBroker.AssertExpectations(t)
		})
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedRepo := new(repositoryMock)

		fakedContender := domain.Contender{
			ID:        fakedContenderID,
			Ownership: fakedOwnership,
		}

		mockedRepo.
			On("GetContender", mock.Anything, mock.Anything, fakedContenderID).
			Return(fakedContender, nil)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		ucase := usecases.ContenderUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		contender, err := ucase.ScrubContender(context.Background(), fakedContenderID)

		assert.ErrorIs(t, err, domain.ErrNoOwnership)
		assert.Empty(t, contender)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})
}

func TestScrubContenders(t *testing.T) {
	fakedCompetitionID := testutils.RandomResourceID[domain.CompetitionID]()
	fakedDeadline := time.Date(2025, 6, 15, 12, 0, 0, 0, time.UTC)

	t.Run("HappyPath", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			mockedRepo := new(repositoryMock)
			mockedEventBroker := new(eventBrokerMock)
			mockedTx := new(transactionMock)

			fakedContenders := []domain.Contender{
				{
					ID:                  testutils.RandomResourceID[domain.ContenderID](),
					CompetitionID:       fakedCompetitionID,
					CategoryID:          testutils.RandomResourceID[domain.CategoryID](),
					Name:                "Alice",
					WithdrawnFromFinals: false,
					Disqualified:        false,
				},
				{
					ID:                  testutils.RandomResourceID[domain.ContenderID](),
					CompetitionID:       fakedCompetitionID,
					CategoryID:          testutils.RandomResourceID[domain.CategoryID](),
					Name:                "Bob",
					WithdrawnFromFinals: true,
					Disqualified:        false,
				},
			}

			mockedRepo.
				On("GetScrubEligibleContenders", mock.Anything, fakedDeadline).
				Return(fakedContenders, nil)

			mockedRepo.On("Begin").Return(mockedTx, nil)

			mockedRepo.
				On("StoreContender", mock.Anything, mockedTx, domain.Contender{
					ID:            fakedContenders[0].ID,
					CompetitionID: fakedCompetitionID,
					CategoryID:    fakedContenders[0].CategoryID,
					ScrubbedAt:    time.Now(),
				}).
				Return(domain.Contender{}, nil).Once()

			mockedRepo.
				On("StoreContender", mock.Anything, mockedTx, domain.Contender{
					ID:                  fakedContenders[1].ID,
					CompetitionID:       fakedCompetitionID,
					CategoryID:          fakedContenders[1].CategoryID,
					WithdrawnFromFinals: true,
					ScrubbedAt:          time.Now(),
				}).
				Return(domain.Contender{}, nil).Once()

			mockedTx.On("Commit").Return(nil)
			mockedTx.On("Rollback").Return()

			mockedEventBroker.
				On("Dispatch", fakedCompetitionID, domain.ContenderPublicInfoUpdatedEvent{
					ContenderID: fakedContenders[0].ID,
					CategoryID:  fakedContenders[0].CategoryID,
					ScrubbedAt:  time.Now(),
				}).Return()

			mockedEventBroker.
				On("Dispatch", fakedCompetitionID, domain.ContenderPublicInfoUpdatedEvent{
					ContenderID:         fakedContenders[1].ID,
					CategoryID:          fakedContenders[1].CategoryID,
					WithdrawnFromFinals: true,
					ScrubbedAt:          time.Now(),
				}).Return()

			ucase := usecases.ContenderUseCase{
				Repo:        mockedRepo,
				EventBroker: mockedEventBroker,
			}

			count, err := ucase.ScrubContenders(context.Background(), fakedDeadline)

			require.NoError(t, err)
			assert.Equal(t, 2, count)

			mockedRepo.AssertExpectations(t)
			mockedEventBroker.AssertExpectations(t)
			mockedTx.AssertExpectations(t)
		})
	})

	t.Run("NoEligibleContenders", func(t *testing.T) {
		mockedRepo := new(repositoryMock)

		mockedRepo.
			On("GetScrubEligibleContenders", mock.Anything, fakedDeadline).
			Return([]domain.Contender{}, nil)

		ucase := usecases.ContenderUseCase{
			Repo: mockedRepo,
		}

		count, err := ucase.ScrubContenders(context.Background(), fakedDeadline)

		require.NoError(t, err)
		assert.Equal(t, 0, count)

		mockedRepo.AssertExpectations(t)
	})
}
