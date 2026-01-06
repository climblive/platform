package usecases_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/usecases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetContender(t *testing.T) {
	fakedContenderID := randomResourceID[domain.ContenderID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: randomResourceID[domain.OrganizerID](),
		ContenderID: &fakedContenderID,
	}

	t.Run("HappyPath", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedRepo := new(repositoryMock)
		mockedScoreKeeper := new(scoreKeeperMock)

		fakedScore := domain.Score{
			Timestamp:   time.Now(),
			ContenderID: fakedContenderID,
			Score:       1000,
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

		mockedScoreKeeper.On("GetScore", fakedContenderID).Return(fakedScore, nil)

		ucase := usecases.ContenderUseCase{
			Repo:        mockedRepo,
			Authorizer:  mockedAuthorizer,
			ScoreKeeper: mockedScoreKeeper,
		}

		contender, err := ucase.GetContender(context.Background(), fakedContenderID)

		require.NoError(t, err)

		assert.Equal(t, fakedContenderID, contender.ID)
		require.NotNil(t, contender.Score)
		assert.Equal(t, fakedScore, *contender.Score)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
		mockedScoreKeeper.AssertExpectations(t)
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
	fakedContenderID := randomResourceID[domain.ContenderID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: randomResourceID[domain.OrganizerID](),
		ContenderID: &fakedContenderID,
	}

	t.Run("HappyPath", func(t *testing.T) {
		mockedRepo := new(repositoryMock)
		mockedScoreKeeper := new(scoreKeeperMock)

		fakedScore := domain.Score{
			Timestamp:   time.Now(),
			ContenderID: fakedContenderID,
			Score:       1000,
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

		mockedScoreKeeper.On("GetScore", fakedContenderID).Return(fakedScore, nil)

		ucase := usecases.ContenderUseCase{
			Repo:        mockedRepo,
			ScoreKeeper: mockedScoreKeeper,
		}

		contender, err := ucase.GetContenderByCode(context.Background(), "ABCD1234")

		require.NoError(t, err)
		assert.Equal(t, fakedContenderID, contender.ID)
		require.NotNil(t, contender.Score)
		assert.Equal(t, fakedScore, *contender.Score)

		mockedRepo.AssertExpectations(t)
		mockedScoreKeeper.AssertExpectations(t)
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

func TestGetContendersByCompClass(t *testing.T) {
	fakedCompClassID := randomResourceID[domain.CompClassID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: randomResourceID[domain.OrganizerID](),
	}
	currentTime := time.Now()

	t.Run("HappyPath", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedRepo := new(repositoryMock)
		mockedScoreKeeper := new(scoreKeeperMock)

		var contenders []domain.Contender

		for k := 1; k <= 10; k++ {
			contenderID := domain.ContenderID(k)

			contenders = append(contenders, domain.Contender{
				ID: contenderID,
			})

			mockedScoreKeeper.On("GetScore", contenderID).Return(domain.Score{
				Timestamp:   currentTime,
				ContenderID: contenderID,
				Score:       k * 10,
				Placement:   k,
				RankOrder:   k - 1,
				Finalist:    true,
			}, nil)
		}

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("GetCompClass", mock.Anything, mock.Anything, fakedCompClassID).
			Return(domain.CompClass{
				ID:        fakedCompClassID,
				Ownership: fakedOwnership,
			}, nil)

		mockedRepo.
			On("GetContendersByCompClass", mock.Anything, mock.Anything, fakedCompClassID).
			Return(contenders, nil)

		ucase := usecases.ContenderUseCase{
			Repo:        mockedRepo,
			Authorizer:  mockedAuthorizer,
			ScoreKeeper: mockedScoreKeeper,
		}

		contenders, err := ucase.GetContendersByCompClass(context.Background(), fakedCompClassID)

		require.NoError(t, err)
		assert.Len(t, contenders, 10)

		for i, contender := range contenders {
			assert.Equal(t, domain.ContenderID(i+1), contender.ID)
			require.NotNil(t, contender.Score)
			assert.Equal(t, (i+1)*10, contender.Score.Score)
			assert.Equal(t, i+1, contender.Score.Placement)
			assert.Equal(t, i, contender.Score.RankOrder)
			assert.True(t, contender.Score.Finalist)
			assert.Equal(t, currentTime, contender.Score.Timestamp)
		}

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
		mockedScoreKeeper.AssertExpectations(t)
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedRepo := new(repositoryMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		mockedRepo.
			On("GetCompClass", mock.Anything, mock.Anything, fakedCompClassID).
			Return(domain.CompClass{
				ID:        fakedCompClassID,
				Ownership: fakedOwnership,
			}, nil)

		ucase := usecases.ContenderUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		contenders, err := ucase.GetContendersByCompClass(context.Background(), fakedCompClassID)

		assert.ErrorIs(t, err, domain.ErrNoOwnership)
		assert.Nil(t, contenders)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})
}

func TestGetContendersByContest(t *testing.T) {
	fakedContestID := randomResourceID[domain.ContestID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: randomResourceID[domain.OrganizerID](),
	}
	currentTime := time.Now()

	t.Run("HappyPath", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedRepo := new(repositoryMock)
		mockedScoreKeeper := new(scoreKeeperMock)

		var contenders []domain.Contender

		for k := 1; k <= 10; k++ {
			contenderID := domain.ContenderID(k)

			contenders = append(contenders, domain.Contender{
				ID: contenderID,
			})

			mockedScoreKeeper.On("GetScore", contenderID).Return(domain.Score{
				Timestamp:   currentTime,
				ContenderID: contenderID,
				Score:       k * 10,
				Placement:   k,
				RankOrder:   k - 1,
				Finalist:    true,
			}, nil)
		}

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("GetContest", mock.Anything, mock.Anything, fakedContestID).
			Return(domain.Contest{
				ID:        fakedContestID,
				Ownership: fakedOwnership,
			}, nil)

		mockedRepo.
			On("GetContendersByContest", mock.Anything, mock.Anything, fakedContestID).
			Return(contenders, nil)

		ucase := usecases.ContenderUseCase{
			Repo:        mockedRepo,
			Authorizer:  mockedAuthorizer,
			ScoreKeeper: mockedScoreKeeper,
		}

		contenders, err := ucase.GetContendersByContest(context.Background(), fakedContestID)

		require.NoError(t, err)
		assert.Len(t, contenders, 10)

		for i, contender := range contenders {
			assert.Equal(t, domain.ContenderID(i+1), contender.ID)
			require.NotNil(t, contender.Score)
			assert.Equal(t, (i+1)*10, contender.Score.Score)
			assert.Equal(t, i+1, contender.Score.Placement)
			assert.Equal(t, i, contender.Score.RankOrder)
			assert.True(t, contender.Score.Finalist)
			assert.Equal(t, currentTime, contender.Score.Timestamp)
		}

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
		mockedScoreKeeper.AssertExpectations(t)
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedRepo := new(repositoryMock)

		mockedRepo.
			On("GetContest", mock.Anything, mock.Anything, fakedContestID).
			Return(domain.Contest{
				ID:        fakedContestID,
				Ownership: fakedOwnership,
			}, nil)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		ucase := usecases.ContenderUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		contenders, err := ucase.GetContendersByContest(context.Background(), fakedContestID)

		assert.ErrorIs(t, err, domain.ErrNoOwnership)
		assert.Nil(t, contenders)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})
}

func TestDeleteContender(t *testing.T) {
	fakedContenderID := randomResourceID[domain.ContenderID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: randomResourceID[domain.OrganizerID](),
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
	fakedContestID := randomResourceID[domain.ContestID]()
	fakedOrganizerID := randomResourceID[domain.OrganizerID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: fakedOrganizerID,
	}

	makeMocks := func() (*repositoryMock, *transactionMock, *codeGeneratorMock) {
		mockedRepo := new(repositoryMock)
		mockedTx := new(transactionMock)
		mockedCodeGenerator := new(codeGeneratorMock)

		mockedRepo.
			On("GetContest", mock.Anything, mock.Anything, fakedContestID).
			Return(domain.Contest{
				ID:        fakedContestID,
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
			On("GetNumberOfContenders", mock.Anything, mock.Anything, fakedContestID).
			Return(400, nil)

		for n := range 100 {
			code := fmt.Sprintf("%08d", n)
			contender := domain.Contender{
				ContestID: fakedContestID,
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

		contenders, err := ucase.CreateContenders(context.Background(), fakedContestID, 100)

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
			On("GetNumberOfContenders", mock.Anything, mock.Anything, fakedContestID).
			Return(400, nil)

		ucase := usecases.ContenderUseCase{
			Repo:                      mockedRepo,
			Authorizer:                mockedAuthorizer,
			RegistrationCodeGenerator: mockedCodeGenerator,
		}

		contenders, err := ucase.CreateContenders(context.Background(), fakedContestID, 101)

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

		contenders, err := ucase.CreateContenders(context.Background(), fakedContestID, 100)

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
			On("GetNumberOfContenders", mock.Anything, mock.Anything, fakedContestID).
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

		contenders, err := ucase.CreateContenders(context.Background(), fakedContestID, 1)

		assert.Error(t, err)
		assert.Nil(t, contenders)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
		mockedTx.AssertExpectations(t)
		mockedCodeGenerator.AssertExpectations(t)
	})
}

func TestPatchContender(t *testing.T) {
	fakedContenderID := randomResourceID[domain.ContenderID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: randomResourceID[domain.OrganizerID](),
		ContenderID: &fakedContenderID,
	}
	fakedContestID := randomResourceID[domain.ContestID]()
	fakedCompClassID := randomResourceID[domain.CompClassID]()

	currentTime := time.Now()
	gracePeriod := 15 * time.Minute

	makeMockedRepo := func(contender domain.Contender) *repositoryMock {
		fakedContest := domain.Contest{
			ID:          fakedContestID,
			GracePeriod: gracePeriod,
		}

		mockedRepo := new(repositoryMock)

		mockedRepo.
			On("GetContender", mock.Anything, mock.Anything, contender.ID).
			Return(contender, nil)

		mockedRepo.
			On("GetContest", mock.Anything, mock.Anything, fakedContestID).
			Return(fakedContest, nil)

		return mockedRepo
	}

	t.Run("UpdateWithoutChanges", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedScoreKeeper := new(scoreKeeperMock)

		fakedContender := domain.Contender{
			ID:                  fakedContenderID,
			Ownership:           fakedOwnership,
			ContestID:           fakedContestID,
			CompClassID:         fakedCompClassID,
			RegistrationCode:    "ABCD1234",
			Name:                "John Doe",
			Entered:             currentTime,
			WithdrawnFromFinals: false,
			Disqualified:        false,
		}

		mockedRepo := makeMockedRepo(fakedContender)

		mockedRepo.
			On("GetCompClass", mock.Anything, mock.Anything, fakedCompClassID).
			Return(domain.CompClass{
				ID:        fakedCompClassID,
				TimeBegin: currentTime.Add(-1 * time.Hour),
				TimeEnd:   currentTime.Add(time.Hour),
			}, nil)

		mockedRepo.
			On("StoreContender", mock.Anything, mock.Anything, mock.AnythingOfType("domain.Contender")).
			Return(mirrorInstruction{}, nil)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.ContenderRole, nil)

		fakedScore := domain.Score{
			Timestamp:   currentTime,
			ContenderID: fakedContenderID,
			Score:       1000,
			Placement:   5,
			Finalist:    true,
			RankOrder:   6,
		}

		mockedScoreKeeper.On("GetScore", fakedContenderID).Return(fakedScore, nil)

		ucase := usecases.ContenderUseCase{
			Repo:        mockedRepo,
			Authorizer:  mockedAuthorizer,
			ScoreKeeper: mockedScoreKeeper,
		}

		contender, err := ucase.PatchContender(context.Background(), fakedContenderID, domain.ContenderPatch{})

		require.NoError(t, err)

		assert.Equal(t, fakedContender.ID, contender.ID)
		assert.Equal(t, fakedContender.Ownership, contender.Ownership)
		assert.Equal(t, fakedContender.ContestID, contender.ContestID)
		assert.Equal(t, fakedContender.CompClassID, contender.CompClassID)
		assert.Equal(t, fakedContender.RegistrationCode, contender.RegistrationCode)
		assert.Equal(t, fakedContender.Name, contender.Name)
		assert.Equal(t, fakedContender.Entered, contender.Entered)
		assert.Equal(t, fakedContender.WithdrawnFromFinals, contender.WithdrawnFromFinals)
		assert.Equal(t, fakedContender.Disqualified, contender.Disqualified)

		require.NotNil(t, contender.Score)
		assert.Equal(t, fakedScore, *contender.Score)

		mockedAuthorizer.AssertExpectations(t)
		mockedScoreKeeper.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("ContenderCannotAlterDisqualifiedState", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)

		fakedContender := domain.Contender{
			ID:                  fakedContenderID,
			Ownership:           fakedOwnership,
			ContestID:           fakedContestID,
			CompClassID:         fakedCompClassID,
			RegistrationCode:    "ABCD1234",
			Name:                "John Doe",
			Entered:             currentTime,
			WithdrawnFromFinals: false,
			Disqualified:        true,
		}

		mockedRepo := makeMockedRepo(fakedContender)

		mockedRepo.
			On("GetCompClass", mock.Anything, mock.Anything, fakedCompClassID).
			Return(domain.CompClass{
				ID:        fakedCompClassID,
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

	t.Run("EnterContest", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedScoreKeeper := new(scoreKeeperMock)
		mockedEventBroker := new(eventBrokerMock)

		fakedContender := domain.Contender{
			ID:                  fakedContenderID,
			Ownership:           fakedOwnership,
			ContestID:           fakedContestID,
			CompClassID:         0,
			RegistrationCode:    "ABCD1234",
			Name:                "",
			Entered:             time.Time{},
			WithdrawnFromFinals: false,
			Disqualified:        false,
		}

		mockedRepo := makeMockedRepo(fakedContender)

		mockedRepo.
			On("GetCompClass", mock.Anything, mock.Anything, fakedCompClassID).
			Return(domain.CompClass{
				ID:        fakedCompClassID,
				TimeBegin: currentTime.Add(-1 * time.Hour),
				TimeEnd:   currentTime.Add(time.Hour),
			}, nil)

		mockedRepo.
			On("StoreContender", mock.Anything, mock.Anything, mock.AnythingOfType("domain.Contender")).
			Return(mirrorInstruction{}, nil)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.ContenderRole, nil)

		mockedScoreKeeper.On("GetScore", fakedContenderID).Return(domain.Score{}, errMock)

		mockedEventBroker.On("Dispatch", fakedContestID, mock.Anything).Return()

		ucase := usecases.ContenderUseCase{
			Repo:        mockedRepo,
			Authorizer:  mockedAuthorizer,
			ScoreKeeper: mockedScoreKeeper,
			EventBroker: mockedEventBroker,
		}

		contender, err := ucase.PatchContender(context.Background(), fakedContenderID, domain.ContenderPatch{
			CompClassID: domain.NewPatch(fakedCompClassID),
			Name:        domain.NewPatch("John Doe"),
		})

		require.NoError(t, err)

		assert.Equal(t, fakedCompClassID, contender.CompClassID)
		assert.Equal(t, "John Doe", contender.Name)
		assert.WithinDuration(t, time.Now(), contender.Entered, time.Minute)

		mockedEventBroker.AssertCalled(t, "Dispatch", fakedContestID, domain.ContenderEnteredEvent{
			ContenderID: fakedContenderID,
			CompClassID: fakedCompClassID,
		})

		mockedEventBroker.AssertCalled(t, "Dispatch", fakedContestID, domain.ContenderPublicInfoUpdatedEvent{
			ContenderID:         fakedContenderID,
			CompClassID:         fakedCompClassID,
			Name:                "John Doe",
			WithdrawnFromFinals: false,
			Disqualified:        false,
		})

		mockedAuthorizer.AssertExpectations(t)
		mockedScoreKeeper.AssertExpectations(t)
		mockedEventBroker.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("CannotMakeChangesToAnUnregisteredContender", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)

		fakedContender := domain.Contender{
			ID:                  fakedContenderID,
			Ownership:           fakedOwnership,
			ContestID:           fakedContestID,
			CompClassID:         0,
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

	t.Run("CannotLeaveContest", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)

		fakedContender := domain.Contender{
			ID:                  fakedContenderID,
			Ownership:           fakedOwnership,
			ContestID:           fakedContestID,
			CompClassID:         fakedCompClassID,
			RegistrationCode:    "ABCD1234",
			Name:                "John Doe",
			Entered:             currentTime,
			WithdrawnFromFinals: false,
			Disqualified:        false,
		}

		mockedRepo := makeMockedRepo(fakedContender)

		mockedRepo.
			On("GetCompClass", mock.Anything, mock.Anything, fakedCompClassID).
			Return(domain.CompClass{
				ID:        fakedCompClassID,
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
			CompClassID: domain.NewPatch(domain.CompClassID(0)),
		})

		assert.ErrorIs(t, err, domain.ErrNotAllowed)
		assert.Empty(t, contender)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("BatchUpdate", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedScoreKeeper := new(scoreKeeperMock)
		mockedEventBroker := new(eventBrokerMock)

		fakedContender := domain.Contender{
			ID:                  fakedContenderID,
			Ownership:           fakedOwnership,
			ContestID:           fakedContestID,
			CompClassID:         fakedCompClassID,
			RegistrationCode:    "ABCD1234",
			Name:                "John Doe",
			Entered:             currentTime,
			WithdrawnFromFinals: false,
			Disqualified:        false,
		}

		mockedRepo := makeMockedRepo(fakedContender)

		mockedRepo.
			On("GetCompClass", mock.Anything, mock.Anything, fakedCompClassID).
			Return(domain.CompClass{
				ID:        fakedCompClassID,
				TimeBegin: currentTime.Add(-1 * time.Hour),
				TimeEnd:   currentTime.Add(time.Hour),
			}, nil)

		mockedRepo.
			On("StoreContender", mock.Anything, mock.Anything, mock.AnythingOfType("domain.Contender")).
			Return(mirrorInstruction{}, nil)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedScoreKeeper.On("GetScore", fakedContenderID).Return(domain.Score{}, errMock)

		mockedEventBroker.On("Dispatch", fakedContestID, mock.Anything).Return()

		fakedOtherCompClass := domain.CompClass{
			ID:        randomResourceID[domain.CompClassID](),
			TimeBegin: currentTime.Add(-1 * time.Hour),
			TimeEnd:   currentTime,
		}

		mockedRepo.
			On("GetCompClass", mock.Anything, mock.Anything, fakedOtherCompClass.ID).
			Return(fakedOtherCompClass, nil)

		ucase := usecases.ContenderUseCase{
			Repo:        mockedRepo,
			Authorizer:  mockedAuthorizer,
			ScoreKeeper: mockedScoreKeeper,
			EventBroker: mockedEventBroker,
		}

		contender, err := ucase.PatchContender(context.Background(), fakedContenderID, domain.ContenderPatch{
			CompClassID:         domain.NewPatch(fakedOtherCompClass.ID),
			Name:                domain.NewPatch("Jane Doe"),
			WithdrawnFromFinals: domain.NewPatch(true),
			Disqualified:        domain.NewPatch(true),
		})

		require.NoError(t, err)

		assert.Equal(t, fakedOtherCompClass.ID, contender.CompClassID)
		assert.Equal(t, "Jane Doe", contender.Name)
		assert.Equal(t, true, contender.WithdrawnFromFinals)
		assert.Equal(t, true, contender.Disqualified)
		assert.Equal(t, currentTime, contender.Entered)

		mockedEventBroker.AssertCalled(t, "Dispatch", fakedContestID, domain.ContenderSwitchedClassEvent{
			ContenderID: fakedContenderID,
			CompClassID: fakedOtherCompClass.ID,
		})

		mockedEventBroker.AssertCalled(t, "Dispatch", fakedContestID, domain.ContenderPublicInfoUpdatedEvent{
			ContenderID:         fakedContenderID,
			CompClassID:         fakedOtherCompClass.ID,
			Name:                "Jane Doe",
			WithdrawnFromFinals: true,
			Disqualified:        true,
		})

		mockedEventBroker.AssertCalled(t, "Dispatch", fakedContestID, domain.ContenderWithdrewFromFinalsEvent{
			ContenderID: fakedContenderID,
		})

		mockedEventBroker.AssertCalled(t, "Dispatch", fakedContestID, domain.ContenderDisqualifiedEvent{
			ContenderID: fakedContenderID,
		})

		mockedAuthorizer.AssertExpectations(t)
		mockedScoreKeeper.AssertExpectations(t)
		mockedEventBroker.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("NameCannotBeEmpty", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)

		fakedContender := domain.Contender{
			ID:                  fakedContenderID,
			Ownership:           fakedOwnership,
			ContestID:           fakedContestID,
			CompClassID:         fakedCompClassID,
			RegistrationCode:    "ABCD1234",
			Name:                "John Doe",
			Entered:             currentTime,
			WithdrawnFromFinals: false,
			Disqualified:        false,
		}

		mockedRepo := makeMockedRepo(fakedContender)

		mockedRepo.
			On("GetCompClass", mock.Anything, mock.Anything, fakedCompClassID).
			Return(domain.CompClass{
				ID:        fakedCompClassID,
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
		mockedScoreKeeper := new(scoreKeeperMock)
		mockedEventBroker := new(eventBrokerMock)

		fakedContender := domain.Contender{
			ID:                  fakedContenderID,
			Ownership:           fakedOwnership,
			ContestID:           fakedContestID,
			CompClassID:         fakedCompClassID,
			RegistrationCode:    "ABCD1234",
			Name:                "John Doe",
			Entered:             currentTime,
			WithdrawnFromFinals: true,
			Disqualified:        false,
		}

		mockedRepo := makeMockedRepo(fakedContender)

		mockedRepo.
			On("GetCompClass", mock.Anything, mock.Anything, fakedCompClassID).
			Return(domain.CompClass{
				ID:        fakedCompClassID,
				TimeBegin: currentTime.Add(-1 * time.Hour),
				TimeEnd:   currentTime.Add(time.Hour),
			}, nil)

		mockedRepo.
			On("StoreContender", mock.Anything, mock.Anything, mock.AnythingOfType("domain.Contender")).
			Return(mirrorInstruction{}, nil)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.ContenderRole, nil)

		mockedScoreKeeper.On("GetScore", fakedContenderID).Return(domain.Score{}, errMock)

		mockedEventBroker.On("Dispatch", fakedContestID, mock.Anything).Return()

		ucase := usecases.ContenderUseCase{
			Repo:        mockedRepo,
			Authorizer:  mockedAuthorizer,
			ScoreKeeper: mockedScoreKeeper,
			EventBroker: mockedEventBroker,
		}

		contender, err := ucase.PatchContender(context.Background(), fakedContenderID, domain.ContenderPatch{
			WithdrawnFromFinals: domain.NewPatch(false),
		})

		require.NoError(t, err)
		assert.Equal(t, false, contender.WithdrawnFromFinals)

		mockedEventBroker.AssertCalled(t, "Dispatch", fakedContestID, domain.ContenderReenteredFinalsEvent{
			ContenderID: fakedContenderID,
		})

		mockedAuthorizer.AssertExpectations(t)
		mockedScoreKeeper.AssertExpectations(t)
		mockedEventBroker.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("Requalify", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedScoreKeeper := new(scoreKeeperMock)
		mockedEventBroker := new(eventBrokerMock)

		fakedContender := domain.Contender{
			ID:                  fakedContenderID,
			Ownership:           fakedOwnership,
			ContestID:           fakedContestID,
			CompClassID:         fakedCompClassID,
			RegistrationCode:    "ABCD1234",
			Name:                "John Doe",
			Entered:             currentTime,
			WithdrawnFromFinals: false,
			Disqualified:        true,
		}

		mockedRepo := makeMockedRepo(fakedContender)

		mockedRepo.
			On("GetCompClass", mock.Anything, mock.Anything, fakedCompClassID).
			Return(domain.CompClass{
				ID:        fakedCompClassID,
				TimeBegin: currentTime.Add(-1 * time.Hour),
				TimeEnd:   currentTime.Add(time.Hour),
			}, nil)

		mockedRepo.
			On("StoreContender", mock.Anything, mock.Anything, mock.AnythingOfType("domain.Contender")).
			Return(mirrorInstruction{}, nil)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedScoreKeeper.On("GetScore", fakedContenderID).Return(domain.Score{}, errMock)

		mockedEventBroker.On("Dispatch", fakedContestID, mock.Anything).Return()

		ucase := usecases.ContenderUseCase{
			Repo:        mockedRepo,
			Authorizer:  mockedAuthorizer,
			ScoreKeeper: mockedScoreKeeper,
			EventBroker: mockedEventBroker,
		}

		contender, err := ucase.PatchContender(context.Background(), fakedContenderID, domain.ContenderPatch{
			Disqualified: domain.NewPatch(false),
		})

		require.NoError(t, err)
		assert.Equal(t, false, contender.Disqualified)

		mockedEventBroker.AssertCalled(t, "Dispatch", fakedContestID, domain.ContenderRequalifiedEvent{
			ContenderID: fakedContenderID,
		})

		mockedAuthorizer.AssertExpectations(t)
		mockedScoreKeeper.AssertExpectations(t)
		mockedEventBroker.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("CannotSwitchToAnEndedCompClass", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.ContenderRole, nil)

		fakedContender := domain.Contender{
			ID:                  fakedContenderID,
			Ownership:           fakedOwnership,
			ContestID:           fakedContestID,
			CompClassID:         fakedCompClassID,
			RegistrationCode:    "ABCD1234",
			Name:                "John Doe",
			Entered:             currentTime,
			WithdrawnFromFinals: false,
			Disqualified:        false,
		}

		fakedOtherCompClass := domain.CompClass{
			ID:        randomResourceID[domain.CompClassID](),
			TimeBegin: currentTime.Add(-1 * time.Hour),
			TimeEnd:   currentTime.Add(-1 * gracePeriod),
		}

		mockedRepo := makeMockedRepo(fakedContender)

		mockedRepo.
			On("GetCompClass", mock.Anything, mock.Anything, fakedCompClassID).
			Return(domain.CompClass{
				ID:        fakedCompClassID,
				TimeBegin: currentTime.Add(-1 * time.Hour),
				TimeEnd:   currentTime.Add(time.Hour),
			}, nil)

		mockedRepo.
			On("GetCompClass", mock.Anything, mock.Anything, fakedOtherCompClass.ID).
			Return(fakedOtherCompClass, nil)

		ucase := usecases.ContenderUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		contender, err := ucase.PatchContender(context.Background(), fakedContenderID, domain.ContenderPatch{
			CompClassID: domain.NewPatch(fakedOtherCompClass.ID),
		})

		assert.ErrorIs(t, err, domain.ErrContestEnded)
		assert.Empty(t, contender)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("ContenderCannotMakeChangesAfterGracePeriod", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.ContenderRole, nil)

		fakedOtherCompClass := domain.CompClass{
			ID:        randomResourceID[domain.CompClassID](),
			TimeBegin: currentTime.Add(-1 * time.Hour),
			TimeEnd:   currentTime.Add(-1 * gracePeriod),
		}

		fakedContender := domain.Contender{
			ID:                  fakedContenderID,
			Ownership:           fakedOwnership,
			ContestID:           fakedContestID,
			CompClassID:         fakedOtherCompClass.ID,
			RegistrationCode:    "ABCD1234",
			Name:                "John Doe",
			Entered:             currentTime,
			WithdrawnFromFinals: false,
			Disqualified:        false,
		}

		mockedRepo := makeMockedRepo(fakedContender)

		mockedRepo.
			On("GetCompClass", mock.Anything, mock.Anything, fakedOtherCompClass.ID).
			Return(fakedOtherCompClass, nil)

		ucase := usecases.ContenderUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		contender, err := ucase.PatchContender(context.Background(), fakedContenderID, domain.ContenderPatch{})
		assert.ErrorIs(t, err, domain.ErrContestEnded)
		assert.Empty(t, contender)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("OrganizerCanMakeChangesAfterGracePeriod", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedScoreKeeper := new(scoreKeeperMock)
		mockedEventBroker := new(eventBrokerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedScoreKeeper.On("GetScore", fakedContenderID).Return(domain.Score{}, errMock)

		mockedEventBroker.On("Dispatch", fakedContestID, mock.Anything).Return()

		fakedSecondCompClass := domain.CompClass{
			ID:        randomResourceID[domain.CompClassID](),
			TimeBegin: currentTime.Add(-1 * time.Hour),
			TimeEnd:   currentTime.Add(-1 * gracePeriod),
		}

		fakedThirdCompClass := domain.CompClass{
			ID:        randomResourceID[domain.CompClassID](),
			TimeBegin: currentTime.Add(-1 * time.Hour),
			TimeEnd:   currentTime.Add(-1 * gracePeriod),
		}

		fakedContender := domain.Contender{
			ID:                  fakedContenderID,
			Ownership:           fakedOwnership,
			ContestID:           fakedContestID,
			CompClassID:         fakedSecondCompClass.ID,
			RegistrationCode:    "ABCD1234",
			Name:                "John Doe",
			Entered:             currentTime,
			WithdrawnFromFinals: false,
			Disqualified:        false,
		}

		mockedRepo := makeMockedRepo(fakedContender)

		mockedRepo.
			On("GetCompClass", mock.Anything, mock.Anything, fakedSecondCompClass.ID).
			Return(fakedSecondCompClass, nil)

		mockedRepo.
			On("GetCompClass", mock.Anything, mock.Anything, fakedThirdCompClass.ID).
			Return(fakedThirdCompClass, nil)

		mockedRepo.
			On("StoreContender", mock.Anything, mock.Anything, mock.AnythingOfType("domain.Contender")).
			Return(mirrorInstruction{}, nil)

		ucase := usecases.ContenderUseCase{
			Repo:        mockedRepo,
			Authorizer:  mockedAuthorizer,
			ScoreKeeper: mockedScoreKeeper,
			EventBroker: mockedEventBroker,
		}

		contender, err := ucase.PatchContender(context.Background(), fakedContenderID, domain.ContenderPatch{
			CompClassID: domain.NewPatch(fakedThirdCompClass.ID),
		})

		require.NoError(t, err)
		assert.Equal(t, fakedThirdCompClass.ID, contender.CompClassID)

		mockedAuthorizer.AssertExpectations(t)
		mockedScoreKeeper.AssertExpectations(t)
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
