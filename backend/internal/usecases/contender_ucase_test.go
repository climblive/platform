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
	mockedContenderID := randomResourceID[domain.ContenderID]()
	mockedOwnership := domain.OwnershipData{
		OrganizerID: randomResourceID[domain.OrganizerID](),
		ContenderID: &mockedContenderID,
	}

	t.Run("HappyPath", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedRepo := new(repositoryMock)
		mockedScoreKeeper := new(scoreKeeperMock)

		mockedScore := domain.Score{
			Timestamp:   time.Now(),
			ContenderID: mockedContenderID,
			Score:       1000,
			Placement:   5,
			Finalist:    true,
			RankOrder:   6,
		}

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mockedOwnership).
			Return(domain.ContenderRole, nil)

		mockedRepo.
			On("GetContender", mock.Anything, mock.Anything, mockedContenderID).
			Return(domain.Contender{
				ID:        mockedContenderID,
				Ownership: mockedOwnership,
			}, nil)

		mockedScoreKeeper.On("GetScore", mockedContenderID).Return(mockedScore, nil)

		ucase := usecases.ContenderUseCase{
			Repo:        mockedRepo,
			Authorizer:  mockedAuthorizer,
			ScoreKeeper: mockedScoreKeeper,
		}

		contender, err := ucase.GetContender(context.Background(), mockedContenderID)

		require.NoError(t, err)

		assert.Equal(t, mockedContenderID, contender.ID)
		require.NotNil(t, contender.Score)
		assert.Equal(t, mockedScore, *contender.Score)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
		mockedScoreKeeper.AssertExpectations(t)
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedRepo := new(repositoryMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mockedOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		mockedRepo.
			On("GetContender", mock.Anything, mock.Anything, mockedContenderID).
			Return(domain.Contender{
				ID:        mockedContenderID,
				Ownership: mockedOwnership,
			}, nil)

		ucase := usecases.ContenderUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		contender, err := ucase.GetContender(context.Background(), mockedContenderID)

		assert.ErrorIs(t, err, domain.ErrNoOwnership)
		assert.Empty(t, contender)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})
}

func TestGetContenderByCode(t *testing.T) {
	mockedContenderID := randomResourceID[domain.ContenderID]()
	mockedOwnership := domain.OwnershipData{
		OrganizerID: randomResourceID[domain.OrganizerID](),
		ContenderID: &mockedContenderID,
	}

	t.Run("HappyPath", func(t *testing.T) {
		mockedRepo := new(repositoryMock)
		mockedScoreKeeper := new(scoreKeeperMock)

		mockedScore := domain.Score{
			Timestamp:   time.Now(),
			ContenderID: mockedContenderID,
			Score:       1000,
			Placement:   5,
			Finalist:    true,
			RankOrder:   6,
		}

		mockedRepo.
			On("GetContenderByCode", mock.Anything, mock.Anything, "ABCD1234").
			Return(domain.Contender{
				ID:        mockedContenderID,
				Ownership: mockedOwnership,
			}, nil)

		mockedScoreKeeper.On("GetScore", mockedContenderID).Return(mockedScore, nil)

		ucase := usecases.ContenderUseCase{
			Repo:        mockedRepo,
			ScoreKeeper: mockedScoreKeeper,
		}

		contender, err := ucase.GetContenderByCode(context.Background(), "ABCD1234")

		require.NoError(t, err)
		assert.Equal(t, mockedContenderID, contender.ID)
		require.NotNil(t, contender.Score)
		assert.Equal(t, mockedScore, *contender.Score)

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
	mockedCompClassID := randomResourceID[domain.CompClassID]()
	mockedOwnership := domain.OwnershipData{
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
			On("HasOwnership", mock.Anything, mockedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("GetCompClass", mock.Anything, mock.Anything, mockedCompClassID).
			Return(domain.CompClass{
				ID:        mockedCompClassID,
				Ownership: mockedOwnership,
			}, nil)

		mockedRepo.
			On("GetContendersByCompClass", mock.Anything, mock.Anything, mockedCompClassID).
			Return(contenders, nil)

		ucase := usecases.ContenderUseCase{
			Repo:        mockedRepo,
			Authorizer:  mockedAuthorizer,
			ScoreKeeper: mockedScoreKeeper,
		}

		contenders, err := ucase.GetContendersByCompClass(context.Background(), mockedCompClassID)

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
			On("HasOwnership", mock.Anything, mockedOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		mockedRepo.
			On("GetCompClass", mock.Anything, mock.Anything, mockedCompClassID).
			Return(domain.CompClass{
				ID:        mockedCompClassID,
				Ownership: mockedOwnership,
			}, nil)

		ucase := usecases.ContenderUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		contenders, err := ucase.GetContendersByCompClass(context.Background(), mockedCompClassID)

		assert.ErrorIs(t, err, domain.ErrNoOwnership)
		assert.Nil(t, contenders)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})
}

func TestGetContendersByContest(t *testing.T) {
	mockedContestID := randomResourceID[domain.ContestID]()
	mockedOwnership := domain.OwnershipData{
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
			On("HasOwnership", mock.Anything, mockedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("GetContest", mock.Anything, mock.Anything, mockedContestID).
			Return(domain.Contest{
				ID:        mockedContestID,
				Ownership: mockedOwnership,
			}, nil)

		mockedRepo.
			On("GetContendersByContest", mock.Anything, mock.Anything, mockedContestID).
			Return(contenders, nil)

		ucase := usecases.ContenderUseCase{
			Repo:        mockedRepo,
			Authorizer:  mockedAuthorizer,
			ScoreKeeper: mockedScoreKeeper,
		}

		contenders, err := ucase.GetContendersByContest(context.Background(), mockedContestID)

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
			On("GetContest", mock.Anything, mock.Anything, mockedContestID).
			Return(domain.Contest{
				ID:        mockedContestID,
				Ownership: mockedOwnership,
			}, nil)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mockedOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		ucase := usecases.ContenderUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		contenders, err := ucase.GetContendersByContest(context.Background(), mockedContestID)

		assert.ErrorIs(t, err, domain.ErrNoOwnership)
		assert.Nil(t, contenders)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})
}

func TestDeleteContender(t *testing.T) {
	mockedContenderID := randomResourceID[domain.ContenderID]()
	mockedOwnership := domain.OwnershipData{
		OrganizerID: randomResourceID[domain.OrganizerID](),
		ContenderID: &mockedContenderID,
	}

	t.Run("HappyPath", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedRepo := new(repositoryMock)

		mockedRepo.
			On("GetContender", mock.Anything, mock.Anything, mockedContenderID).
			Return(domain.Contender{Ownership: mockedOwnership}, nil)

		mockedRepo.
			On("DeleteContender", mock.Anything, mock.Anything, mockedContenderID).
			Return(nil)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mockedOwnership).
			Return(domain.OrganizerRole, nil)

		ucase := usecases.ContenderUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		err := ucase.DeleteContender(context.Background(), mockedContenderID)

		require.NoError(t, err)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedRepo := new(repositoryMock)

		mockedRepo.
			On("GetContender", mock.Anything, mock.Anything, mockedContenderID).
			Return(domain.Contender{Ownership: mockedOwnership}, nil)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mockedOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		ucase := usecases.ContenderUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		err := ucase.DeleteContender(context.Background(), mockedContenderID)

		assert.ErrorIs(t, err, domain.ErrNoOwnership)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("InsufficientRole", func(t *testing.T) {
		for _, insufficientRole := range []domain.AuthRole{domain.NilRole, domain.ContenderRole, domain.JudgeRole} {
			mockedAuthorizer := new(authorizerMock)
			mockedRepo := new(repositoryMock)

			mockedAuthorizer.
				On("HasOwnership", mock.Anything, mockedOwnership).
				Return(insufficientRole, nil)

			mockedRepo.
				On("GetContender", mock.Anything, mock.Anything, mockedContenderID).
				Return(domain.Contender{Ownership: mockedOwnership}, nil)

			ucase := usecases.ContenderUseCase{
				Repo:       mockedRepo,
				Authorizer: mockedAuthorizer,
			}

			err := ucase.DeleteContender(context.Background(), mockedContenderID)

			assert.ErrorIs(t, err, domain.ErrInsufficientRole)

			mockedAuthorizer.AssertExpectations(t)
			mockedRepo.AssertExpectations(t)
		}
	})
}

func TestCreateContenders(t *testing.T) {
	mockedContestID := randomResourceID[domain.ContestID]()
	mockedOrganizerID := randomResourceID[domain.OrganizerID]()
	mockedOwnership := domain.OwnershipData{
		OrganizerID: mockedOrganizerID,
	}

	makeMocks := func() (*repositoryMock, *transactionMock, *codeGeneratorMock) {
		mockedRepo := new(repositoryMock)
		mockedTx := new(transactionMock)
		mockedCodeGenerator := new(codeGeneratorMock)

		mockedRepo.
			On("GetContest", mock.Anything, mock.Anything, mockedContestID).
			Return(domain.Contest{
				ID:        mockedContestID,
				Ownership: mockedOwnership,
			}, nil)

		return mockedRepo, mockedTx, mockedCodeGenerator
	}

	t.Run("HappyPath", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedRepo, mockedTx, mockedCodeGenerator := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mockedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("GetNumberOfContenders", mock.Anything, mock.Anything, mockedContestID).
			Return(400, nil)

		for n := range 100 {
			code := fmt.Sprintf("%08d", n)
			contender := domain.Contender{
				ContestID: mockedContestID,
				Ownership: domain.OwnershipData{
					OrganizerID: mockedOrganizerID,
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

		contenders, err := ucase.CreateContenders(context.Background(), mockedContestID, 100)

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
			On("HasOwnership", mock.Anything, mockedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("GetNumberOfContenders", mock.Anything, mock.Anything, mockedContestID).
			Return(400, nil)

		ucase := usecases.ContenderUseCase{
			Repo:                      mockedRepo,
			Authorizer:                mockedAuthorizer,
			RegistrationCodeGenerator: mockedCodeGenerator,
		}

		contenders, err := ucase.CreateContenders(context.Background(), mockedContestID, 101)

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
			On("HasOwnership", mock.Anything, mockedOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		ucase := usecases.ContenderUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		contenders, err := ucase.CreateContenders(context.Background(), mockedContestID, 100)

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
			On("GetNumberOfContenders", mock.Anything, mock.Anything, mockedContestID).
			Return(0, nil)

		mockedRepo.
			On("StoreContender", mock.Anything, mock.Anything, mock.Anything).
			Return(domain.Contender{}, errMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mockedOwnership).
			Return(domain.AdminRole, nil)

		mockedCodeGenerator.
			On("Generate", 8).
			Return("ABCD1234")

		ucase := usecases.ContenderUseCase{
			Repo:                      mockedRepo,
			Authorizer:                mockedAuthorizer,
			RegistrationCodeGenerator: mockedCodeGenerator,
		}

		contenders, err := ucase.CreateContenders(context.Background(), mockedContestID, 1)

		assert.Error(t, err)
		assert.Nil(t, contenders)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
		mockedTx.AssertExpectations(t)
		mockedCodeGenerator.AssertExpectations(t)
	})
}

func TestUpdateContender(t *testing.T) {
	mockedContenderID := randomResourceID[domain.ContenderID]()
	mockedOwnership := domain.OwnershipData{
		OrganizerID: randomResourceID[domain.OrganizerID](),
		ContenderID: &mockedContenderID,
	}
	mockedContestID := randomResourceID[domain.ContestID]()
	mockedCompClassID := randomResourceID[domain.CompClassID]()

	currentTime := time.Now()
	gracePeriod := 15 * time.Minute

	makeMockedRepo := func(contender domain.Contender) *repositoryMock {
		mockedContest := domain.Contest{
			ID:          mockedContestID,
			GracePeriod: gracePeriod,
		}

		mockedRepo := new(repositoryMock)

		mockedRepo.
			On("GetContender", mock.Anything, mock.Anything, contender.ID).
			Return(contender, nil)

		mockedRepo.
			On("GetContest", mock.Anything, mock.Anything, mockedContestID).
			Return(mockedContest, nil)

		return mockedRepo
	}

	t.Run("UpdateWithoutChanges", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedScoreKeeper := new(scoreKeeperMock)

		mockedContender := domain.Contender{
			ID:                  mockedContenderID,
			Ownership:           mockedOwnership,
			ContestID:           mockedContestID,
			CompClassID:         mockedCompClassID,
			RegistrationCode:    "ABCD1234",
			Name:                "John Doe",
			PublicName:          "John",
			ClubName:            "Testers' Climbing Club",
			Entered:             &currentTime,
			WithdrawnFromFinals: false,
			Disqualified:        false,
		}

		mockedRepo := makeMockedRepo(mockedContender)

		mockedRepo.
			On("GetCompClass", mock.Anything, mock.Anything, mockedCompClassID).
			Return(domain.CompClass{
				ID:        mockedCompClassID,
				TimeBegin: currentTime.Add(-1 * time.Hour),
				TimeEnd:   currentTime.Add(time.Hour),
			}, nil)

		mockedRepo.
			On("StoreContender", mock.Anything, mock.Anything, mock.AnythingOfType("domain.Contender")).
			Return(mirrorInstruction{}, nil)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mockedOwnership).
			Return(domain.ContenderRole, nil)

		mockedScore := domain.Score{
			Timestamp:   currentTime,
			ContenderID: mockedContenderID,
			Score:       1000,
			Placement:   5,
			Finalist:    true,
			RankOrder:   6,
		}

		mockedScoreKeeper.On("GetScore", mockedContenderID).Return(mockedScore, nil)

		ucase := usecases.ContenderUseCase{
			Repo:        mockedRepo,
			Authorizer:  mockedAuthorizer,
			ScoreKeeper: mockedScoreKeeper,
		}

		contender, err := ucase.UpdateContender(context.Background(), mockedContenderID, mockedContender)

		require.NoError(t, err)

		assert.Equal(t, mockedContender.ID, contender.ID)
		assert.Equal(t, mockedContender.Ownership, contender.Ownership)
		assert.Equal(t, mockedContender.ContestID, contender.ContestID)
		assert.Equal(t, mockedContender.CompClassID, contender.CompClassID)
		assert.Equal(t, mockedContender.RegistrationCode, contender.RegistrationCode)
		assert.Equal(t, mockedContender.Name, contender.Name)
		assert.Equal(t, mockedContender.PublicName, contender.PublicName)
		assert.Equal(t, mockedContender.ClubName, contender.ClubName)
		assert.Equal(t, mockedContender.Entered, contender.Entered)
		assert.Equal(t, mockedContender.WithdrawnFromFinals, contender.WithdrawnFromFinals)
		assert.Equal(t, mockedContender.Disqualified, contender.Disqualified)

		require.NotNil(t, contender.Score)
		assert.Equal(t, mockedScore, *contender.Score)

		mockedAuthorizer.AssertExpectations(t)
		mockedScoreKeeper.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("ReadOnlyFieldsAreUnaltered", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedScoreKeeper := new(scoreKeeperMock)

		mockedContender := domain.Contender{
			ID:                  mockedContenderID,
			Ownership:           mockedOwnership,
			ContestID:           mockedContestID,
			CompClassID:         mockedCompClassID,
			RegistrationCode:    "ABCD1234",
			Name:                "John Doe",
			PublicName:          "John",
			ClubName:            "Testers' Climbing Club",
			Entered:             &currentTime,
			WithdrawnFromFinals: false,
			Disqualified:        false,
		}

		mockedRepo := makeMockedRepo(mockedContender)

		mockedRepo.
			On("GetCompClass", mock.Anything, mock.Anything, mockedCompClassID).
			Return(domain.CompClass{
				ID:        mockedCompClassID,
				TimeBegin: currentTime.Add(-1 * time.Hour),
				TimeEnd:   currentTime.Add(time.Hour),
			}, nil)

		mockedRepo.
			On("StoreContender", mock.Anything, mock.Anything, mock.AnythingOfType("domain.Contender")).
			Return(mirrorInstruction{}, nil)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mockedOwnership).
			Return(domain.ContenderRole, nil)

		mockedScoreKeeper.On("GetScore", mockedContenderID).Return(domain.Score{}, errMock)

		ucase := usecases.ContenderUseCase{
			Repo:        mockedRepo,
			Authorizer:  mockedAuthorizer,
			ScoreKeeper: mockedScoreKeeper,
		}

		modifiers := []func(domain.Contender) domain.Contender{
			func(contender domain.Contender) domain.Contender {
				contender.ID += 1
				return contender
			},
			func(contender domain.Contender) domain.Contender {
				contender.Ownership.OrganizerID += 1
				return contender
			},
			func(contender domain.Contender) domain.Contender {
				contender.ContestID += 1
				return contender
			},
			func(contender domain.Contender) domain.Contender {
				contender.RegistrationCode = "DEADBEEF"
				return contender
			},
			func(contender domain.Contender) domain.Contender {
				contender.Score = &domain.Score{}
				return contender
			},
			func(contender domain.Contender) domain.Contender {
				soon := currentTime.Add(time.Hour)
				contender.Entered = &soon
				return contender
			},
		}

		for _, fn := range modifiers {
			contender, err := ucase.UpdateContender(context.Background(), mockedContenderID, fn(mockedContender))

			require.NoError(t, err)
			assert.Equal(t, mockedContender, contender)
		}

		mockedAuthorizer.AssertExpectations(t)
		mockedScoreKeeper.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("ContenderCannotAlterDisqualifiedState", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)

		mockedContender := domain.Contender{
			ID:                  mockedContenderID,
			Ownership:           mockedOwnership,
			ContestID:           mockedContestID,
			CompClassID:         mockedCompClassID,
			RegistrationCode:    "ABCD1234",
			Name:                "John Doe",
			PublicName:          "John",
			ClubName:            "Testers' Climbing Club",
			Entered:             &currentTime,
			WithdrawnFromFinals: false,
			Disqualified:        true,
		}

		mockedRepo := makeMockedRepo(mockedContender)

		mockedRepo.
			On("GetCompClass", mock.Anything, mock.Anything, mockedCompClassID).
			Return(domain.CompClass{
				ID:        mockedCompClassID,
				TimeBegin: currentTime.Add(-1 * time.Hour),
				TimeEnd:   currentTime.Add(time.Hour),
			}, nil)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mockedOwnership).
			Return(domain.ContenderRole, nil)

		ucase := usecases.ContenderUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		updatedContender := mockedContender
		updatedContender.Disqualified = false

		contender, err := ucase.UpdateContender(context.Background(), mockedContenderID, updatedContender)

		assert.ErrorIs(t, err, domain.ErrInsufficientRole)
		assert.Empty(t, contender)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("EnterContest", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedScoreKeeper := new(scoreKeeperMock)
		mockedEventBroker := new(eventBrokerMock)

		mockedContender := domain.Contender{
			ID:                  mockedContenderID,
			Ownership:           mockedOwnership,
			ContestID:           mockedContestID,
			CompClassID:         0,
			RegistrationCode:    "ABCD1234",
			Name:                "",
			PublicName:          "",
			ClubName:            "",
			Entered:             nil,
			WithdrawnFromFinals: false,
			Disqualified:        false,
		}

		mockedRepo := makeMockedRepo(mockedContender)

		mockedRepo.
			On("GetCompClass", mock.Anything, mock.Anything, mockedCompClassID).
			Return(domain.CompClass{
				ID:        mockedCompClassID,
				TimeBegin: currentTime.Add(-1 * time.Hour),
				TimeEnd:   currentTime.Add(time.Hour),
			}, nil)

		mockedRepo.
			On("StoreContender", mock.Anything, mock.Anything, mock.AnythingOfType("domain.Contender")).
			Return(mirrorInstruction{}, nil)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mockedOwnership).
			Return(domain.ContenderRole, nil)

		mockedScoreKeeper.On("GetScore", mockedContenderID).Return(domain.Score{}, errMock)

		mockedEventBroker.On("Dispatch", mockedContestID, mock.Anything).Return()

		ucase := usecases.ContenderUseCase{
			Repo:        mockedRepo,
			Authorizer:  mockedAuthorizer,
			ScoreKeeper: mockedScoreKeeper,
			EventBroker: mockedEventBroker,
		}

		updatedContender := mockedContender
		updatedContender.CompClassID = mockedCompClassID
		updatedContender.Name = "John Doe"
		updatedContender.PublicName = "John"
		updatedContender.ClubName = "Testers' Climbing Club"

		contender, err := ucase.UpdateContender(context.Background(), mockedContenderID, updatedContender)

		require.NoError(t, err)

		assert.Equal(t, mockedCompClassID, contender.CompClassID)
		assert.Equal(t, "John Doe", contender.Name)
		assert.Equal(t, "John", contender.PublicName)
		assert.Equal(t, "Testers' Climbing Club", contender.ClubName)
		require.NotNil(t, contender.Entered)
		assert.WithinDuration(t, time.Now(), *contender.Entered, time.Minute)

		mockedEventBroker.AssertCalled(t, "Dispatch", mockedContestID, domain.ContenderEnteredEvent{
			ContenderID: mockedContenderID,
			CompClassID: mockedCompClassID,
		})

		mockedEventBroker.AssertCalled(t, "Dispatch", mockedContestID, domain.ContenderPublicInfoUpdatedEvent{
			ContenderID:         mockedContenderID,
			CompClassID:         mockedCompClassID,
			PublicName:          "John",
			ClubName:            "Testers' Climbing Club",
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

		mockedContender := domain.Contender{
			ID:                  mockedContenderID,
			Ownership:           mockedOwnership,
			ContestID:           mockedContestID,
			CompClassID:         0,
			RegistrationCode:    "ABCD1234",
			Name:                "",
			PublicName:          "",
			ClubName:            "",
			Entered:             nil,
			WithdrawnFromFinals: false,
			Disqualified:        false,
		}

		mockedRepo := makeMockedRepo(mockedContender)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mockedOwnership).
			Return(domain.AdminRole, nil)

		ucase := usecases.ContenderUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		contender, err := ucase.UpdateContender(context.Background(), mockedContenderID, mockedContender)

		assert.ErrorIs(t, err, domain.ErrNotRegistered)
		assert.Empty(t, contender)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("CannotLeaveContest", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)

		mockedContender := domain.Contender{
			ID:                  mockedContenderID,
			Ownership:           mockedOwnership,
			ContestID:           mockedContestID,
			CompClassID:         mockedCompClassID,
			RegistrationCode:    "ABCD1234",
			Name:                "John Doe",
			PublicName:          "John",
			ClubName:            "Testers' Climbing Club",
			Entered:             &currentTime,
			WithdrawnFromFinals: false,
			Disqualified:        false,
		}

		mockedRepo := makeMockedRepo(mockedContender)

		mockedRepo.
			On("GetCompClass", mock.Anything, mock.Anything, mockedCompClassID).
			Return(domain.CompClass{
				ID:        mockedCompClassID,
				TimeBegin: currentTime.Add(-1 * time.Hour),
				TimeEnd:   currentTime.Add(time.Hour),
			}, nil)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mockedOwnership).
			Return(domain.OrganizerRole, nil)

		ucase := usecases.ContenderUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		updatedContender := mockedContender
		updatedContender.CompClassID = 0

		contender, err := ucase.UpdateContender(context.Background(), mockedContenderID, updatedContender)

		assert.ErrorIs(t, err, domain.ErrNotAllowed)
		assert.Empty(t, contender)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("BatchUpdate", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedScoreKeeper := new(scoreKeeperMock)
		mockedEventBroker := new(eventBrokerMock)

		mockedContender := domain.Contender{
			ID:                  mockedContenderID,
			Ownership:           mockedOwnership,
			ContestID:           mockedContestID,
			CompClassID:         mockedCompClassID,
			RegistrationCode:    "ABCD1234",
			Name:                "John Doe",
			PublicName:          "John",
			ClubName:            "Testers' Climbing Club",
			Entered:             &currentTime,
			WithdrawnFromFinals: false,
			Disqualified:        false,
		}

		mockedRepo := makeMockedRepo(mockedContender)

		mockedRepo.
			On("GetCompClass", mock.Anything, mock.Anything, mockedCompClassID).
			Return(domain.CompClass{
				ID:        mockedCompClassID,
				TimeBegin: currentTime.Add(-1 * time.Hour),
				TimeEnd:   currentTime.Add(time.Hour),
			}, nil)

		mockedRepo.
			On("StoreContender", mock.Anything, mock.Anything, mock.AnythingOfType("domain.Contender")).
			Return(mirrorInstruction{}, nil)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mockedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedScoreKeeper.On("GetScore", mockedContenderID).Return(domain.Score{}, errMock)

		mockedEventBroker.On("Dispatch", mockedContestID, mock.Anything).Return()

		mockedOtherCompClass := domain.CompClass{
			ID:        randomResourceID[domain.CompClassID](),
			TimeBegin: currentTime.Add(-1 * time.Hour),
			TimeEnd:   currentTime,
		}

		mockedRepo.
			On("GetCompClass", mock.Anything, mock.Anything, mockedOtherCompClass.ID).
			Return(mockedOtherCompClass, nil)

		ucase := usecases.ContenderUseCase{
			Repo:        mockedRepo,
			Authorizer:  mockedAuthorizer,
			ScoreKeeper: mockedScoreKeeper,
			EventBroker: mockedEventBroker,
		}

		updatedContender := mockedContender
		updatedContender.CompClassID = mockedOtherCompClass.ID
		updatedContender.Name = "Jane Doe"
		updatedContender.PublicName = "Jane"
		updatedContender.ClubName = "Space Climbers"
		updatedContender.WithdrawnFromFinals = true
		updatedContender.Disqualified = true

		contender, err := ucase.UpdateContender(context.Background(), mockedContenderID, updatedContender)

		require.NoError(t, err)

		assert.Equal(t, mockedOtherCompClass.ID, contender.CompClassID)
		assert.Equal(t, "Jane Doe", contender.Name)
		assert.Equal(t, "Jane", contender.PublicName)
		assert.Equal(t, "Space Climbers", contender.ClubName)
		assert.Equal(t, true, contender.WithdrawnFromFinals)
		assert.Equal(t, true, contender.Disqualified)
		assert.Equal(t, currentTime, *contender.Entered)

		mockedEventBroker.AssertCalled(t, "Dispatch", mockedContestID, domain.ContenderSwitchedClassEvent{
			ContenderID: mockedContenderID,
			CompClassID: mockedOtherCompClass.ID,
		})

		mockedEventBroker.AssertCalled(t, "Dispatch", mockedContestID, domain.ContenderPublicInfoUpdatedEvent{
			ContenderID:         mockedContenderID,
			CompClassID:         mockedOtherCompClass.ID,
			PublicName:          "Jane",
			ClubName:            "Space Climbers",
			WithdrawnFromFinals: true,
			Disqualified:        true,
		})

		mockedEventBroker.AssertCalled(t, "Dispatch", mockedContestID, domain.ContenderWithdrewFromFinalsEvent{
			ContenderID: mockedContenderID,
		})

		mockedEventBroker.AssertCalled(t, "Dispatch", mockedContestID, domain.ContenderDisqualifiedEvent{
			ContenderID: mockedContenderID,
		})

		mockedAuthorizer.AssertExpectations(t)
		mockedScoreKeeper.AssertExpectations(t)
		mockedEventBroker.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("NameCannotBeEmpty", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)

		mockedContender := domain.Contender{
			ID:                  mockedContenderID,
			Ownership:           mockedOwnership,
			ContestID:           mockedContestID,
			CompClassID:         mockedCompClassID,
			RegistrationCode:    "ABCD1234",
			Name:                "John Doe",
			PublicName:          "John",
			ClubName:            "Testers' Climbing Club",
			Entered:             &currentTime,
			WithdrawnFromFinals: false,
			Disqualified:        false,
		}

		mockedRepo := makeMockedRepo(mockedContender)

		mockedRepo.
			On("GetCompClass", mock.Anything, mock.Anything, mockedCompClassID).
			Return(domain.CompClass{
				ID:        mockedCompClassID,
				TimeBegin: currentTime.Add(-1 * time.Hour),
				TimeEnd:   currentTime.Add(time.Hour),
			}, nil)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mockedOwnership).
			Return(domain.ContenderRole, nil)

		ucase := usecases.ContenderUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		updatedContender := mockedContender
		updatedContender.Name = ""

		contender, err := ucase.UpdateContender(context.Background(), mockedContenderID, updatedContender)

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

		mockedContender := domain.Contender{
			ID:                  mockedContenderID,
			Ownership:           mockedOwnership,
			ContestID:           mockedContestID,
			CompClassID:         mockedCompClassID,
			RegistrationCode:    "ABCD1234",
			Name:                "John Doe",
			PublicName:          "John",
			ClubName:            "Testers' Climbing Club",
			Entered:             &currentTime,
			WithdrawnFromFinals: true,
			Disqualified:        false,
		}

		mockedRepo := makeMockedRepo(mockedContender)

		mockedRepo.
			On("GetCompClass", mock.Anything, mock.Anything, mockedCompClassID).
			Return(domain.CompClass{
				ID:        mockedCompClassID,
				TimeBegin: currentTime.Add(-1 * time.Hour),
				TimeEnd:   currentTime.Add(time.Hour),
			}, nil)

		mockedRepo.
			On("StoreContender", mock.Anything, mock.Anything, mock.AnythingOfType("domain.Contender")).
			Return(mirrorInstruction{}, nil)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mockedOwnership).
			Return(domain.ContenderRole, nil)

		mockedScoreKeeper.On("GetScore", mockedContenderID).Return(domain.Score{}, errMock)

		mockedEventBroker.On("Dispatch", mockedContestID, mock.Anything).Return()

		ucase := usecases.ContenderUseCase{
			Repo:        mockedRepo,
			Authorizer:  mockedAuthorizer,
			ScoreKeeper: mockedScoreKeeper,
			EventBroker: mockedEventBroker,
		}

		updatedContender := mockedContender
		updatedContender.WithdrawnFromFinals = false

		contender, err := ucase.UpdateContender(context.Background(), mockedContenderID, updatedContender)

		require.NoError(t, err)
		assert.Equal(t, false, contender.WithdrawnFromFinals)

		mockedEventBroker.AssertCalled(t, "Dispatch", mockedContestID, domain.ContenderReenteredFinalsEvent{
			ContenderID: mockedContenderID,
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

		mockedContender := domain.Contender{
			ID:                  mockedContenderID,
			Ownership:           mockedOwnership,
			ContestID:           mockedContestID,
			CompClassID:         mockedCompClassID,
			RegistrationCode:    "ABCD1234",
			Name:                "John Doe",
			PublicName:          "John",
			ClubName:            "Testers' Climbing Club",
			Entered:             &currentTime,
			WithdrawnFromFinals: false,
			Disqualified:        true,
		}

		mockedRepo := makeMockedRepo(mockedContender)

		mockedRepo.
			On("GetCompClass", mock.Anything, mock.Anything, mockedCompClassID).
			Return(domain.CompClass{
				ID:        mockedCompClassID,
				TimeBegin: currentTime.Add(-1 * time.Hour),
				TimeEnd:   currentTime.Add(time.Hour),
			}, nil)

		mockedRepo.
			On("StoreContender", mock.Anything, mock.Anything, mock.AnythingOfType("domain.Contender")).
			Return(mirrorInstruction{}, nil)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mockedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedScoreKeeper.On("GetScore", mockedContenderID).Return(domain.Score{}, errMock)

		mockedEventBroker.On("Dispatch", mockedContestID, mock.Anything).Return()

		ucase := usecases.ContenderUseCase{
			Repo:        mockedRepo,
			Authorizer:  mockedAuthorizer,
			ScoreKeeper: mockedScoreKeeper,
			EventBroker: mockedEventBroker,
		}

		updatedContender := mockedContender
		updatedContender.Disqualified = false

		contender, err := ucase.UpdateContender(context.Background(), mockedContenderID, updatedContender)

		require.NoError(t, err)
		assert.Equal(t, false, contender.Disqualified)

		mockedEventBroker.AssertCalled(t, "Dispatch", mockedContestID, domain.ContenderRequalifiedEvent{
			ContenderID: mockedContenderID,
		})

		mockedAuthorizer.AssertExpectations(t)
		mockedScoreKeeper.AssertExpectations(t)
		mockedEventBroker.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("CannotSwitchToAnEndedCompClass", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mockedOwnership).
			Return(domain.ContenderRole, nil)

		mockedContender := domain.Contender{
			ID:                  mockedContenderID,
			Ownership:           mockedOwnership,
			ContestID:           mockedContestID,
			CompClassID:         mockedCompClassID,
			RegistrationCode:    "ABCD1234",
			Name:                "John Doe",
			PublicName:          "John",
			ClubName:            "Testers' Climbing Club",
			Entered:             &currentTime,
			WithdrawnFromFinals: false,
			Disqualified:        false,
		}

		mockedOtherCompClass := domain.CompClass{
			ID:        randomResourceID[domain.CompClassID](),
			TimeBegin: currentTime.Add(-1 * time.Hour),
			TimeEnd:   currentTime.Add(-1 * gracePeriod),
		}

		mockedRepo := makeMockedRepo(mockedContender)

		mockedRepo.
			On("GetCompClass", mock.Anything, mock.Anything, mockedCompClassID).
			Return(domain.CompClass{
				ID:        mockedCompClassID,
				TimeBegin: currentTime.Add(-1 * time.Hour),
				TimeEnd:   currentTime.Add(time.Hour),
			}, nil)

		mockedRepo.
			On("GetCompClass", mock.Anything, mock.Anything, mockedOtherCompClass.ID).
			Return(mockedOtherCompClass, nil)

		ucase := usecases.ContenderUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		updatedContender := mockedContender
		updatedContender.CompClassID = mockedOtherCompClass.ID

		contender, err := ucase.UpdateContender(context.Background(), mockedContenderID, updatedContender)

		assert.ErrorIs(t, err, domain.ErrContestEnded)
		assert.Empty(t, contender)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("ContenderCannotMakeChangesAfterGracePeriod", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mockedOwnership).
			Return(domain.ContenderRole, nil)

		mockedOtherCompClass := domain.CompClass{
			ID:        randomResourceID[domain.CompClassID](),
			TimeBegin: currentTime.Add(-1 * time.Hour),
			TimeEnd:   currentTime.Add(-1 * gracePeriod),
		}

		mockedContender := domain.Contender{
			ID:                  mockedContenderID,
			Ownership:           mockedOwnership,
			ContestID:           mockedContestID,
			CompClassID:         mockedOtherCompClass.ID,
			RegistrationCode:    "ABCD1234",
			Name:                "John Doe",
			PublicName:          "John",
			ClubName:            "Testers' Climbing Club",
			Entered:             &currentTime,
			WithdrawnFromFinals: false,
			Disqualified:        false,
		}

		mockedRepo := makeMockedRepo(mockedContender)

		mockedRepo.
			On("GetCompClass", mock.Anything, mock.Anything, mockedOtherCompClass.ID).
			Return(mockedOtherCompClass, nil)

		ucase := usecases.ContenderUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		contender, err := ucase.UpdateContender(context.Background(), mockedContenderID, mockedContender)
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
			On("HasOwnership", mock.Anything, mockedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedScoreKeeper.On("GetScore", mockedContenderID).Return(domain.Score{}, errMock)

		mockedEventBroker.On("Dispatch", mockedContestID, mock.Anything).Return()

		mockedSecondCompClass := domain.CompClass{
			ID:        randomResourceID[domain.CompClassID](),
			TimeBegin: currentTime.Add(-1 * time.Hour),
			TimeEnd:   currentTime.Add(-1 * gracePeriod),
		}

		mockedThirdCompClass := domain.CompClass{
			ID:        randomResourceID[domain.CompClassID](),
			TimeBegin: currentTime.Add(-1 * time.Hour),
			TimeEnd:   currentTime.Add(-1 * gracePeriod),
		}

		mockedContender := domain.Contender{
			ID:                  mockedContenderID,
			Ownership:           mockedOwnership,
			ContestID:           mockedContestID,
			CompClassID:         mockedSecondCompClass.ID,
			RegistrationCode:    "ABCD1234",
			Name:                "John Doe",
			PublicName:          "John",
			ClubName:            "Testers' Climbing Club",
			Entered:             &currentTime,
			WithdrawnFromFinals: false,
			Disqualified:        false,
		}

		mockedRepo := makeMockedRepo(mockedContender)

		mockedRepo.
			On("GetCompClass", mock.Anything, mock.Anything, mockedSecondCompClass.ID).
			Return(mockedSecondCompClass, nil)

		mockedRepo.
			On("GetCompClass", mock.Anything, mock.Anything, mockedThirdCompClass.ID).
			Return(mockedThirdCompClass, nil)

		mockedRepo.
			On("StoreContender", mock.Anything, mock.Anything, mock.AnythingOfType("domain.Contender")).
			Return(mirrorInstruction{}, nil)

		ucase := usecases.ContenderUseCase{
			Repo:        mockedRepo,
			Authorizer:  mockedAuthorizer,
			ScoreKeeper: mockedScoreKeeper,
			EventBroker: mockedEventBroker,
		}

		updatedContender := mockedContender
		updatedContender.CompClassID = mockedThirdCompClass.ID

		contender, err := ucase.UpdateContender(context.Background(), mockedContenderID, updatedContender)

		require.NoError(t, err)
		assert.Equal(t, mockedThirdCompClass.ID, contender.CompClassID)

		mockedAuthorizer.AssertExpectations(t)
		mockedScoreKeeper.AssertExpectations(t)
		mockedEventBroker.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedRepo := new(repositoryMock)

		mockedRepo.
			On("GetContender", mock.Anything, mock.Anything, mockedContenderID).
			Return(domain.Contender{
				ID:        mockedContenderID,
				Ownership: mockedOwnership,
			}, nil)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mockedOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		ucase := usecases.ContenderUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		contender, err := ucase.UpdateContender(context.Background(), mockedContenderID, domain.Contender{})

		assert.ErrorIs(t, err, domain.ErrNoOwnership)
		assert.Empty(t, contender)

		mockedAuthorizer.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})
}
