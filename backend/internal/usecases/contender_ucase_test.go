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
	mockedContenderID := domain.ResourceID(1)
	mockedOwnership := domain.OwnershipData{
		OrganizerID: 1,
		ContenderID: &mockedContenderID,
	}
	currentTime := time.Now()

	mockedContender := domain.Contender{
		ID:        mockedContenderID,
		Ownership: mockedOwnership,
	}

	mockedRepo := new(repositoryMock)
	mockedScoreKeeper := new(scoreKeeperMock)

	mockedRepo.
		On("GetContender", mock.Anything, mock.Anything, mockedContenderID).
		Return(mockedContender, nil)

	mockedScoreKeeper.On("GetScore", mockedContenderID).Return(domain.Score{
		Timestamp:   currentTime,
		ContenderID: mockedContenderID,
		Score:       1000,
		Placement:   5,
	}, nil)

	t.Run("HappyPath", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mockedOwnership).
			Return(domain.ContenderRole, nil)

		ucase := usecases.ContenderUseCase{
			Repo:        mockedRepo,
			Authorizer:  mockedAuthorizer,
			ScoreKeeper: mockedScoreKeeper,
		}

		contender, err := ucase.GetContender(context.Background(), mockedContenderID)

		require.NoError(t, err)

		assert.Equal(t, mockedContenderID, contender.ID)
		assert.Equal(t, 1000, contender.Score)
		assert.Equal(t, 5, contender.Placement)
		assert.Equal(t, currentTime, *contender.ScoreUpdated)

	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mockedOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		ucase := usecases.ContenderUseCase{
			Repo:        mockedRepo,
			Authorizer:  mockedAuthorizer,
			ScoreKeeper: mockedScoreKeeper,
		}

		contender, err := ucase.GetContender(context.Background(), mockedContenderID)

		assert.ErrorIs(t, err, domain.ErrNoOwnership)
		assert.Empty(t, contender)
	})
}

func TestGetContenderByCode(t *testing.T) {
	mockedContenderID := domain.ResourceID(1)
	mockedOwnership := domain.OwnershipData{
		OrganizerID: 1,
		ContenderID: &mockedContenderID,
	}

	mockedContender := domain.Contender{
		ID:        mockedContenderID,
		Ownership: mockedOwnership,
	}

	mockedRepo := new(repositoryMock)
	mockedScoreKeeper := new(scoreKeeperMock)

	mockedRepo.
		On("GetContenderByCode", mock.Anything, mock.Anything, "ABCD1234").
		Return(mockedContender, nil)

	mockedRepo.
		On("GetContenderByCode", mock.Anything, mock.Anything, mock.AnythingOfType("string")).
		Return(domain.Contender{}, domain.ErrNotFound)

	mockedScoreKeeper.On("GetScore", mockedContenderID).Return(domain.Score{}, errMock)

	t.Run("HappyPath", func(t *testing.T) {
		ucase := usecases.ContenderUseCase{
			Repo:        mockedRepo,
			ScoreKeeper: mockedScoreKeeper,
		}

		contender, err := ucase.GetContenderByCode(context.Background(), "ABCD1234")

		require.NoError(t, err)
		assert.Equal(t, mockedContenderID, contender.ID)
	})

	t.Run("NotFound", func(t *testing.T) {
		ucase := usecases.ContenderUseCase{
			Repo:        mockedRepo,
			ScoreKeeper: mockedScoreKeeper,
		}

		contender, err := ucase.GetContenderByCode(context.Background(), "DEADBEEF")

		assert.ErrorIs(t, err, domain.ErrNotFound)
		assert.Empty(t, contender)
	})
}

func TestGetContendersByCompClass(t *testing.T) {
	mockedCompClassID := domain.ResourceID(1)
	mockedOwnership := domain.OwnershipData{
		OrganizerID: 1,
	}
	currentTime := time.Now()

	mockedCompClass := domain.CompClass{
		ID:        mockedCompClassID,
		Ownership: mockedOwnership,
	}

	mockedRepo := new(repositoryMock)
	mockedScoreKeeper := new(scoreKeeperMock)

	mockedRepo.
		On("GetCompClass", mock.Anything, mock.Anything, mockedCompClassID).
		Return(mockedCompClass, nil)

	var contenders []domain.Contender

	for k := 1; k <= 10; k++ {
		contenderID := k

		contenders = append(contenders, domain.Contender{
			ID: contenderID,
		})

		mockedScoreKeeper.On("GetScore", contenderID).Return(domain.Score{
			Timestamp:   currentTime,
			ContenderID: contenderID,
			Score:       k * 10,
			Placement:   k,
		}, nil)
	}

	mockedRepo.
		On("GetContendersByCompClass", mock.Anything, mock.Anything, mockedCompClassID).
		Return(contenders, nil)

	t.Run("HappyPath", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mockedOwnership).
			Return(domain.OrganizerRole, nil)

		ucase := usecases.ContenderUseCase{
			Repo:        mockedRepo,
			Authorizer:  mockedAuthorizer,
			ScoreKeeper: mockedScoreKeeper,
		}

		contenders, err := ucase.GetContendersByCompClass(context.Background(), mockedCompClassID)

		require.NoError(t, err)
		assert.Len(t, contenders, 10)

		for i, contender := range contenders {
			assert.Equal(t, i+1, contender.ID)
			assert.Equal(t, (i+1)*10, contender.Score)
			assert.Equal(t, i+1, contender.Placement)
			assert.Equal(t, currentTime, *contender.ScoreUpdated)
		}
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mockedOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		ucase := usecases.ContenderUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		contenders, err := ucase.GetContendersByCompClass(context.Background(), mockedCompClassID)

		assert.ErrorIs(t, err, domain.ErrNoOwnership)
		assert.Nil(t, contenders)
	})
}

func TestGetContendersByContest(t *testing.T) {
	mockedContestID := domain.ResourceID(1)
	mockedOwnership := domain.OwnershipData{
		OrganizerID: 1,
	}
	currentTime := time.Now()

	mockedContest := domain.Contest{
		ID:        mockedContestID,
		Ownership: mockedOwnership,
	}

	mockedRepo := new(repositoryMock)
	mockedScoreKeeper := new(scoreKeeperMock)

	mockedRepo.
		On("GetContest", mock.Anything, mock.Anything, mockedContestID).
		Return(mockedContest, nil)

	var contenders []domain.Contender

	for k := 1; k <= 10; k++ {
		contenderID := k

		contenders = append(contenders, domain.Contender{
			ID: contenderID,
		})

		mockedScoreKeeper.On("GetScore", contenderID).Return(domain.Score{
			Timestamp:   currentTime,
			ContenderID: contenderID,
			Score:       k * 10,
			Placement:   k,
		}, nil)
	}

	mockedRepo.
		On("GetContendersByContest", mock.Anything, mock.Anything, mockedContestID).
		Return(contenders, nil)

	t.Run("HappyPath", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mockedOwnership).
			Return(domain.OrganizerRole, nil)

		ucase := usecases.ContenderUseCase{
			Repo:        mockedRepo,
			Authorizer:  mockedAuthorizer,
			ScoreKeeper: mockedScoreKeeper,
		}

		contenders, err := ucase.GetContendersByContest(context.Background(), mockedContestID)

		require.NoError(t, err)
		assert.Len(t, contenders, 10)

		for i, contender := range contenders {
			assert.Equal(t, i+1, contender.ID)
			assert.Equal(t, (i+1)*10, contender.Score)
			assert.Equal(t, i+1, contender.Placement)
			assert.Equal(t, currentTime, *contender.ScoreUpdated)
		}
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)

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
	})
}

func TestDeleteContender(t *testing.T) {
	mockedContenderID := domain.ResourceID(1)
	mockedOwnership := domain.OwnershipData{
		OrganizerID: 1,
		ContenderID: &mockedContenderID,
	}

	mockedRepo := new(repositoryMock)

	mockedRepo.
		On("GetContender", mock.Anything, mock.Anything, mockedContenderID).
		Return(domain.Contender{Ownership: mockedOwnership}, nil)

	mockedRepo.
		On("DeleteContender", mock.Anything, mock.Anything, mockedContenderID).
		Return(nil)

	t.Run("HappyPath", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mockedOwnership).
			Return(domain.OrganizerRole, nil)

		ucase := usecases.ContenderUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		err := ucase.DeleteContender(context.Background(), mockedContenderID)

		require.NoError(t, err)
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mockedOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		ucase := usecases.ContenderUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		err := ucase.DeleteContender(context.Background(), mockedContenderID)

		assert.ErrorIs(t, err, domain.ErrNoOwnership)
	})

	t.Run("InsufficientRole", func(t *testing.T) {
		for _, insufficientRole := range []domain.AuthRole{domain.NilRole, domain.ContenderRole, domain.JudgeRole} {
			mockedAuthorizer := new(authorizerMock)

			mockedAuthorizer.
				On("HasOwnership", mock.Anything, mockedOwnership).
				Return(insufficientRole, nil)

			ucase := usecases.ContenderUseCase{
				Repo:       mockedRepo,
				Authorizer: mockedAuthorizer,
			}

			err := ucase.DeleteContender(context.Background(), mockedContenderID)

			assert.ErrorIs(t, err, domain.ErrInsufficientRole)
		}
	})
}

func TestCreateContenders(t *testing.T) {
	mockedContestID := domain.ResourceID(1)
	mockedOwnership := domain.OwnershipData{
		OrganizerID: 1,
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

		mockedRepo.
			On("GetNumberOfContenders", mock.Anything, mock.Anything, mockedContestID).
			Return(400, nil)

		for n := range 100 {
			code := fmt.Sprintf("%08d", n)
			contender := domain.Contender{
				ContestID: 1,
				Ownership: domain.OwnershipData{
					OrganizerID: 1,
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
		mockedTx.On("Rollback").Return()

		return mockedRepo, mockedTx, mockedCodeGenerator
	}

	t.Run("HappyPath", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedRepo, mockedTx, mockedCodeGenerator := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mockedOwnership).
			Return(domain.OrganizerRole, nil)

		ucase := usecases.ContenderUseCase{
			Repo:                      mockedRepo,
			Authorizer:                mockedAuthorizer,
			RegistrationCodeGenerator: mockedCodeGenerator,
		}

		contenders, err := ucase.CreateContenders(context.Background(), mockedContestID, 100)

		require.NoError(t, err)
		assert.Len(t, contenders, 100)

		mockedRepo.AssertExpectations(t)
		mockedTx.AssertNumberOfCalls(t, "Commit", 1)
		mockedTx.AssertNotCalled(t, "Rollback")

		for i, contender := range contenders {
			assert.Equal(t, fmt.Sprintf("%08d", i), contender.RegistrationCode)
		}
	})

	t.Run("CannotExceed500Contenders", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedRepo, _, mockedCodeGenerator := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mockedOwnership).
			Return(domain.OrganizerRole, nil)

		ucase := usecases.ContenderUseCase{
			Repo:                      mockedRepo,
			Authorizer:                mockedAuthorizer,
			RegistrationCodeGenerator: mockedCodeGenerator,
		}

		contenders, err := ucase.CreateContenders(context.Background(), mockedContestID, 101)

		assert.ErrorIs(t, err, domain.ErrLimitExceeded)
		assert.Nil(t, contenders)
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
	})
}

func TestCreateContenders_Rollback(t *testing.T) {
	mockedContestID := domain.ResourceID(1)

	mockedRepo := new(repositoryMock)
	mockedTx := new(transactionMock)
	mockedAuthorizer := new(authorizerMock)
	mockedCodeGenerator := new(codeGeneratorMock)

	mockedRepo.
		On("GetContest", mock.Anything, mock.Anything, mockedContestID).
		Return(domain.Contest{}, nil)

	mockedRepo.
		On("Begin").
		Return(mockedTx, nil)

	mockedRepo.
		On("StoreContender", mock.Anything, mock.Anything, mock.Anything).
		Return(domain.Contender{}, errMock)

	mockedRepo.
		On("GetNumberOfContenders", mock.Anything, mock.Anything, mockedContestID).
		Return(0, nil)

	mockedTx.On("Rollback").Return()

	mockedAuthorizer.
		On("HasOwnership", mock.Anything, domain.OwnershipData{}).
		Return(domain.AdminRole, nil)

	mockedCodeGenerator.
		On("Generate", 8).
		Return("DEAFBEEF")

	ucase := usecases.ContenderUseCase{
		Repo:                      mockedRepo,
		Authorizer:                mockedAuthorizer,
		RegistrationCodeGenerator: mockedCodeGenerator,
	}

	contenders, err := ucase.CreateContenders(context.Background(), mockedContestID, 1)

	assert.Error(t, err)
	assert.Nil(t, contenders)

	mockedRepo.AssertExpectations(t)
	mockedTx.AssertExpectations(t)
}

func TestUpdateContender(t *testing.T) {
	mockedContenderID := domain.ResourceID(1)
	mockedOwnership := domain.OwnershipData{
		OrganizerID: 1,
		ContenderID: &mockedContenderID,
	}
	currentTime := time.Now()
	gracePeriod := 15 * time.Minute

	makeMockedRepo := func(contender domain.Contender) *repositoryMock {
		mockedContest := domain.Contest{
			ID:          1,
			GracePeriod: gracePeriod,
		}

		mockedCompClass := domain.CompClass{
			ID:        1,
			TimeBegin: currentTime.Add(-1 * time.Hour),
			TimeEnd:   currentTime.Add(time.Hour),
		}

		mockedRepo := new(repositoryMock)

		mockedRepo.
			On("GetContender", mock.Anything, mock.Anything, contender.ID).
			Return(contender, nil)

		mockedRepo.
			On("GetContest", mock.Anything, mock.Anything, 1).
			Return(mockedContest, nil)

		mockedRepo.
			On("GetCompClass", mock.Anything, mock.Anything, 1).
			Return(mockedCompClass, nil)

		mockedRepo.
			On("StoreContender", mock.Anything, mock.Anything, mock.AnythingOfType("domain.Contender")).
			Return(mirrorInstruction{}, nil)

		return mockedRepo
	}

	t.Run("UpdateWithoutChanges", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedScoreKeeper := new(scoreKeeperMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mock.Anything).
			Return(domain.ContenderRole, nil)

		mockedScoreKeeper.On("GetScore", mockedContenderID).Return(domain.Score{
			Timestamp:   currentTime,
			ContenderID: mockedContenderID,
			Score:       1000,
			Placement:   5,
		}, nil)

		mockedContender := domain.Contender{
			ID:                  mockedContenderID,
			Ownership:           mockedOwnership,
			ContestID:           1,
			CompClassID:         1,
			RegistrationCode:    "ABCD1234",
			Name:                "John Doe",
			PublicName:          "John",
			ClubName:            "Testers' Climbing Club",
			Entered:             &currentTime,
			WithdrawnFromFinals: false,
			Disqualified:        false,
			Score:               100,
			Placement:           1,
			ScoreUpdated:        &currentTime,
		}

		ucase := usecases.ContenderUseCase{
			Repo:        makeMockedRepo(mockedContender),
			Authorizer:  mockedAuthorizer,
			ScoreKeeper: mockedScoreKeeper,
		}

		contender, err := ucase.UpdateContender(context.Background(), mockedContenderID, mockedContender)

		assert.NoError(t, err)

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

		assert.Equal(t, 1000, contender.Score)
		assert.Equal(t, 5, contender.Placement)
		assert.Equal(t, currentTime, *contender.ScoreUpdated)
	})

	t.Run("ReadOnlyFieldsAreUnaltered", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedScoreKeeper := new(scoreKeeperMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mock.Anything).
			Return(domain.ContenderRole, nil)

		mockedScoreKeeper.On("GetScore", mockedContenderID).Return(domain.Score{}, errMock)

		mockedContender := domain.Contender{
			ID:                  mockedContenderID,
			Ownership:           mockedOwnership,
			ContestID:           1,
			CompClassID:         1,
			RegistrationCode:    "ABCD1234",
			Name:                "John Doe",
			PublicName:          "John",
			ClubName:            "Testers' Climbing Club",
			Entered:             &currentTime,
			WithdrawnFromFinals: false,
			Disqualified:        false,
			Score:               100,
			Placement:           1,
			ScoreUpdated:        &currentTime,
		}

		ucase := usecases.ContenderUseCase{
			Repo:        makeMockedRepo(mockedContender),
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
				contender.Score += 1
				return contender
			},
			func(contender domain.Contender) domain.Contender {
				contender.Placement += 1
				return contender
			},
			func(contender domain.Contender) domain.Contender {
				soon := currentTime.Add(time.Hour)
				contender.Entered = &soon
				return contender
			},
			func(contender domain.Contender) domain.Contender {
				soon := currentTime.Add(time.Hour)
				contender.ScoreUpdated = &soon
				return contender
			},
		}

		for _, fn := range modifiers {
			contender, err := ucase.UpdateContender(context.Background(), mockedContenderID, fn(mockedContender))

			assert.NoError(t, err)
			assert.Equal(t, mockedContender, contender)
		}

	})

	t.Run("ContenderCannotAlterDisqualifiedState", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedScoreKeeper := new(scoreKeeperMock)
		mockedEventBroker := new(eventBrokerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mock.Anything).
			Return(domain.ContenderRole, nil)

		mockedScoreKeeper.On("GetScore", mockedContenderID).Return(domain.Score{}, errMock)

		mockedEventBroker.On("Dispatch", 1, mock.Anything).Return()

		mockedContender := domain.Contender{
			ID:                  mockedContenderID,
			Ownership:           mockedOwnership,
			ContestID:           1,
			CompClassID:         1,
			RegistrationCode:    "ABCD1234",
			Name:                "John Doe",
			PublicName:          "John",
			ClubName:            "Testers' Climbing Club",
			Entered:             &currentTime,
			WithdrawnFromFinals: false,
			Disqualified:        true,
			Score:               100,
			Placement:           1,
			ScoreUpdated:        &currentTime,
		}

		ucase := usecases.ContenderUseCase{
			Repo:        makeMockedRepo(mockedContender),
			Authorizer:  mockedAuthorizer,
			ScoreKeeper: mockedScoreKeeper,
			EventBroker: mockedEventBroker,
		}

		_, err := ucase.UpdateContender(context.Background(), mockedContenderID, mockedContender)
		assert.NoError(t, err)

		updatedContender := mockedContender
		updatedContender.Disqualified = false

		contender, err := ucase.UpdateContender(context.Background(), mockedContenderID, updatedContender)

		assert.ErrorIs(t, err, domain.ErrInsufficientRole)
		assert.Empty(t, contender)
	})

	t.Run("EnterContest", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedScoreKeeper := new(scoreKeeperMock)
		mockedEventBroker := new(eventBrokerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mock.Anything).
			Return(domain.ContenderRole, nil)

		mockedScoreKeeper.On("GetScore", mockedContenderID).Return(domain.Score{}, errMock)

		mockedEventBroker.On("Dispatch", 1, mock.Anything).Return()

		mockedContender := domain.Contender{
			ID:                  mockedContenderID,
			Ownership:           mockedOwnership,
			ContestID:           1,
			CompClassID:         0,
			RegistrationCode:    "ABCD1234",
			Name:                "",
			PublicName:          "",
			ClubName:            "",
			Entered:             nil,
			WithdrawnFromFinals: false,
			Disqualified:        false,
			Score:               0,
			Placement:           0,
			ScoreUpdated:        nil,
		}

		ucase := usecases.ContenderUseCase{
			Repo:        makeMockedRepo(mockedContender),
			Authorizer:  mockedAuthorizer,
			ScoreKeeper: mockedScoreKeeper,
			EventBroker: mockedEventBroker,
		}

		updatedContender := mockedContender
		updatedContender.CompClassID = 1
		updatedContender.Name = "John Doe"
		updatedContender.PublicName = "John"
		updatedContender.ClubName = "Testers' Climbing Club"

		contender, err := ucase.UpdateContender(context.Background(), mockedContenderID, updatedContender)

		require.NoError(t, err)
		assert.Equal(t, 1, contender.CompClassID)
		assert.Equal(t, "John Doe", contender.Name)
		assert.Equal(t, "John", contender.PublicName)
		assert.Equal(t, "Testers' Climbing Club", contender.ClubName)
		require.NotNil(t, contender.Entered)
		assert.WithinDuration(t, time.Now(), *contender.Entered, time.Minute)

		mockedEventBroker.AssertCalled(t, "Dispatch", 1, domain.ContenderEnterEvent{
			ContenderID: mockedContenderID,
			CompClassID: 1,
		})

		mockedEventBroker.AssertCalled(t, "Dispatch", 1, domain.ContenderPublicInfoUpdateEvent{
			ContenderID:         mockedContenderID,
			CompClassID:         1,
			PublicName:          "John",
			ClubName:            "Testers' Climbing Club",
			WithdrawnFromFinals: false,
			Disqualified:        false,
		})
	})

	t.Run("CannotChangeUnregisteredContender", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mock.Anything).
			Return(domain.AdminRole, nil)

		mockedContender := domain.Contender{
			ID:                  mockedContenderID,
			Ownership:           mockedOwnership,
			ContestID:           1,
			CompClassID:         0,
			RegistrationCode:    "ABCD1234",
			Name:                "",
			PublicName:          "",
			ClubName:            "",
			Entered:             nil,
			WithdrawnFromFinals: false,
			Disqualified:        false,
			Score:               0,
			Placement:           0,
			ScoreUpdated:        nil,
		}

		ucase := usecases.ContenderUseCase{
			Repo:       makeMockedRepo(mockedContender),
			Authorizer: mockedAuthorizer,
		}

		contender, err := ucase.UpdateContender(context.Background(), mockedContenderID, mockedContender)

		assert.ErrorIs(t, err, domain.ErrNotRegistered)
		assert.Empty(t, contender)
	})

	t.Run("CannotLeaveContest", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mock.Anything).
			Return(domain.OrganizerRole, nil)

		mockedContender := domain.Contender{
			ID:                  mockedContenderID,
			Ownership:           mockedOwnership,
			ContestID:           1,
			CompClassID:         1,
			RegistrationCode:    "ABCD1234",
			Name:                "John Doe",
			PublicName:          "John",
			ClubName:            "Testers' Climbing Club",
			Entered:             &currentTime,
			WithdrawnFromFinals: false,
			Disqualified:        false,
			Score:               100,
			Placement:           1,
			ScoreUpdated:        &currentTime,
		}

		ucase := usecases.ContenderUseCase{
			Repo:       makeMockedRepo(mockedContender),
			Authorizer: mockedAuthorizer,
		}

		updatedContender := mockedContender
		updatedContender.CompClassID = 0

		contender, err := ucase.UpdateContender(context.Background(), mockedContenderID, updatedContender)

		assert.ErrorIs(t, err, domain.ErrNotAllowed)
		assert.Empty(t, contender)
	})

	t.Run("BatchUpdate", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedScoreKeeper := new(scoreKeeperMock)
		mockedEventBroker := new(eventBrokerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mock.Anything).
			Return(domain.OrganizerRole, nil)

		mockedScoreKeeper.On("GetScore", mockedContenderID).Return(domain.Score{}, errMock)

		mockedEventBroker.On("Dispatch", 1, mock.Anything).Return()

		mockedContender := domain.Contender{
			ID:                  mockedContenderID,
			Ownership:           mockedOwnership,
			ContestID:           1,
			CompClassID:         1,
			RegistrationCode:    "ABCD1234",
			Name:                "John Doe",
			PublicName:          "John",
			ClubName:            "Testers' Climbing Club",
			Entered:             &currentTime,
			WithdrawnFromFinals: false,
			Disqualified:        false,
			Score:               100,
			Placement:           1,
			ScoreUpdated:        &currentTime,
		}

		mockedCompClass := domain.CompClass{
			ID:        2,
			TimeBegin: currentTime.Add(-1 * time.Hour),
			TimeEnd:   currentTime,
		}

		mockedRepo := makeMockedRepo(mockedContender)

		mockedRepo.
			On("GetCompClass", mock.Anything, mock.Anything, 2).
			Return(mockedCompClass, nil)

		ucase := usecases.ContenderUseCase{
			Repo:        mockedRepo,
			Authorizer:  mockedAuthorizer,
			ScoreKeeper: mockedScoreKeeper,
			EventBroker: mockedEventBroker,
		}

		updatedContender := mockedContender
		updatedContender.CompClassID = 2
		updatedContender.Name = "Jane Doe"
		updatedContender.PublicName = "Jane"
		updatedContender.ClubName = "Space Climbers"
		updatedContender.WithdrawnFromFinals = true
		updatedContender.Disqualified = true

		contender, err := ucase.UpdateContender(context.Background(), mockedContenderID, updatedContender)

		assert.NoError(t, err)
		assert.Equal(t, 2, contender.CompClassID)
		assert.Equal(t, "Jane Doe", contender.Name)
		assert.Equal(t, "Jane", contender.PublicName)
		assert.Equal(t, "Space Climbers", contender.ClubName)
		assert.Equal(t, true, contender.WithdrawnFromFinals)
		assert.Equal(t, true, contender.Disqualified)
		assert.Equal(t, currentTime, *contender.Entered)

		mockedEventBroker.AssertCalled(t, "Dispatch", 1, domain.ContenderSwitchClassEvent{
			ContenderID: mockedContenderID,
			CompClassID: 2,
		})

		mockedEventBroker.AssertCalled(t, "Dispatch", 1, domain.ContenderPublicInfoUpdateEvent{
			ContenderID:         mockedContenderID,
			CompClassID:         2,
			PublicName:          "Jane",
			ClubName:            "Space Climbers",
			WithdrawnFromFinals: true,
			Disqualified:        true,
		})

		mockedEventBroker.AssertCalled(t, "Dispatch", 1, domain.ContenderWithdrawFromFinalsEvent{
			ContenderID: mockedContenderID,
		})

		mockedEventBroker.AssertCalled(t, "Dispatch", 1, domain.ContenderDisqualifyEvent{
			ContenderID: mockedContenderID,
		})
	})

	t.Run("NameCannotBeEmpty", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mock.Anything).
			Return(domain.ContenderRole, nil)

		mockedContender := domain.Contender{
			ID:                  mockedContenderID,
			Ownership:           mockedOwnership,
			ContestID:           1,
			CompClassID:         1,
			RegistrationCode:    "ABCD1234",
			Name:                "John Doe",
			PublicName:          "John",
			ClubName:            "Testers' Climbing Club",
			Entered:             &currentTime,
			WithdrawnFromFinals: false,
			Disqualified:        false,
			Score:               100,
			Placement:           1,
			ScoreUpdated:        &currentTime,
		}

		ucase := usecases.ContenderUseCase{
			Repo:       makeMockedRepo(mockedContender),
			Authorizer: mockedAuthorizer,
		}

		updatedContender := mockedContender
		updatedContender.Name = ""

		contender, err := ucase.UpdateContender(context.Background(), mockedContenderID, updatedContender)

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.ErrorIs(t, err, domain.ErrEmptyName)
		assert.Empty(t, contender)
	})

	t.Run("ReenterFinals", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedScoreKeeper := new(scoreKeeperMock)
		mockedEventBroker := new(eventBrokerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mock.Anything).
			Return(domain.OrganizerRole, nil)

		mockedScoreKeeper.On("GetScore", mockedContenderID).Return(domain.Score{}, errMock)

		mockedEventBroker.On("Dispatch", 1, mock.Anything).Return()

		mockedContender := domain.Contender{
			ID:                  mockedContenderID,
			Ownership:           mockedOwnership,
			ContestID:           1,
			CompClassID:         1,
			RegistrationCode:    "ABCD1234",
			Name:                "John Doe",
			PublicName:          "John",
			ClubName:            "Testers' Climbing Club",
			Entered:             &currentTime,
			WithdrawnFromFinals: true,
			Disqualified:        false,
			Score:               100,
			Placement:           1,
			ScoreUpdated:        &currentTime,
		}

		ucase := usecases.ContenderUseCase{
			Repo:        makeMockedRepo(mockedContender),
			Authorizer:  mockedAuthorizer,
			ScoreKeeper: mockedScoreKeeper,
			EventBroker: mockedEventBroker,
		}

		updatedContender := mockedContender
		updatedContender.WithdrawnFromFinals = false

		contender, err := ucase.UpdateContender(context.Background(), mockedContenderID, updatedContender)

		assert.NoError(t, err)
		assert.Equal(t, false, contender.WithdrawnFromFinals)

		mockedEventBroker.AssertCalled(t, "Dispatch", 1, domain.ContenderReenterFinalsEvent{
			ContenderID: mockedContenderID,
		})
	})

	t.Run("Requalify", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedScoreKeeper := new(scoreKeeperMock)
		mockedEventBroker := new(eventBrokerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mock.Anything).
			Return(domain.OrganizerRole, nil)

		mockedScoreKeeper.On("GetScore", mockedContenderID).Return(domain.Score{}, errMock)

		mockedEventBroker.On("Dispatch", 1, mock.Anything).Return()

		mockedContender := domain.Contender{
			ID:                  mockedContenderID,
			Ownership:           mockedOwnership,
			ContestID:           1,
			CompClassID:         1,
			RegistrationCode:    "ABCD1234",
			Name:                "John Doe",
			PublicName:          "John",
			ClubName:            "Testers' Climbing Club",
			Entered:             &currentTime,
			WithdrawnFromFinals: false,
			Disqualified:        true,
			Score:               100,
			Placement:           1,
			ScoreUpdated:        &currentTime,
		}

		ucase := usecases.ContenderUseCase{
			Repo:        makeMockedRepo(mockedContender),
			Authorizer:  mockedAuthorizer,
			ScoreKeeper: mockedScoreKeeper,
			EventBroker: mockedEventBroker,
		}

		updatedContender := mockedContender
		updatedContender.Disqualified = false

		contender, err := ucase.UpdateContender(context.Background(), mockedContenderID, updatedContender)

		assert.NoError(t, err)
		assert.Equal(t, false, contender.Disqualified)

		mockedEventBroker.AssertCalled(t, "Dispatch", 1, domain.ContenderRequalifyEvent{
			ContenderID: mockedContenderID,
		})
	})

	t.Run("CannotSwitchToAnEndedCompClass", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mock.Anything).
			Return(domain.ContenderRole, nil)

		mockedContender := domain.Contender{
			ID:                  mockedContenderID,
			Ownership:           mockedOwnership,
			ContestID:           1,
			CompClassID:         1,
			RegistrationCode:    "ABCD1234",
			Name:                "John Doe",
			PublicName:          "John",
			ClubName:            "Testers' Climbing Club",
			Entered:             &currentTime,
			WithdrawnFromFinals: false,
			Disqualified:        true,
			Score:               100,
			Placement:           1,
			ScoreUpdated:        &currentTime,
		}

		mockedCompClass := domain.CompClass{
			ID:        2,
			TimeBegin: currentTime.Add(-1 * time.Hour),
			TimeEnd:   currentTime.Add(-1 * gracePeriod),
		}

		mockedRepo := makeMockedRepo(mockedContender)

		mockedRepo.
			On("GetCompClass", mock.Anything, mock.Anything, 2).
			Return(mockedCompClass, nil)

		ucase := usecases.ContenderUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		updatedContender := mockedContender
		updatedContender.CompClassID = 2

		contender, err := ucase.UpdateContender(context.Background(), mockedContenderID, updatedContender)

		assert.ErrorIs(t, err, domain.ErrContestEnded)
		assert.Empty(t, contender)
	})

	t.Run("ContenderCannotMakeChangesAfterGracePeriod", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mock.Anything).
			Return(domain.ContenderRole, nil)

		mockedContender := domain.Contender{
			ID:                  mockedContenderID,
			Ownership:           mockedOwnership,
			ContestID:           1,
			CompClassID:         2,
			RegistrationCode:    "ABCD1234",
			Name:                "John Doe",
			PublicName:          "John",
			ClubName:            "Testers' Climbing Club",
			Entered:             &currentTime,
			WithdrawnFromFinals: false,
			Disqualified:        true,
			Score:               100,
			Placement:           1,
			ScoreUpdated:        &currentTime,
		}

		mockedCompClass := domain.CompClass{
			ID:        2,
			TimeBegin: currentTime.Add(-1 * time.Hour),
			TimeEnd:   currentTime.Add(-1 * gracePeriod),
		}

		mockedRepo := makeMockedRepo(mockedContender)

		mockedRepo.
			On("GetCompClass", mock.Anything, mock.Anything, 2).
			Return(mockedCompClass, nil)

		ucase := usecases.ContenderUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		contender, err := ucase.UpdateContender(context.Background(), mockedContenderID, mockedContender)
		assert.ErrorIs(t, err, domain.ErrContestEnded)
		assert.Empty(t, contender)
	})

	t.Run("OrganizerCanMakeChangesAfterGracePeriod", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)
		mockedScoreKeeper := new(scoreKeeperMock)
		mockedEventBroker := new(eventBrokerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mock.Anything).
			Return(domain.OrganizerRole, nil)

		mockedScoreKeeper.On("GetScore", mockedContenderID).Return(domain.Score{}, errMock)

		mockedEventBroker.On("Dispatch", 1, mock.Anything).Return()

		mockedContender := domain.Contender{
			ID:                  mockedContenderID,
			Ownership:           mockedOwnership,
			ContestID:           1,
			CompClassID:         2,
			RegistrationCode:    "ABCD1234",
			Name:                "John Doe",
			PublicName:          "John",
			ClubName:            "Testers' Climbing Club",
			Entered:             &currentTime,
			WithdrawnFromFinals: false,
			Disqualified:        false,
			Score:               100,
			Placement:           1,
			ScoreUpdated:        &currentTime,
		}

		mockedRepo := makeMockedRepo(mockedContender)

		mockedRepo.
			On("GetCompClass", mock.Anything, mock.Anything, 2).
			Return(domain.CompClass{
				ID:        2,
				TimeBegin: currentTime.Add(-1 * time.Hour),
				TimeEnd:   currentTime.Add(-1 * gracePeriod),
			}, nil)

		mockedRepo.
			On("GetCompClass", mock.Anything, mock.Anything, 3).
			Return(domain.CompClass{
				ID:        3,
				TimeBegin: currentTime.Add(-1 * time.Hour),
				TimeEnd:   currentTime.Add(-1 * gracePeriod),
			}, nil)

		ucase := usecases.ContenderUseCase{
			Repo:        mockedRepo,
			Authorizer:  mockedAuthorizer,
			ScoreKeeper: mockedScoreKeeper,
			EventBroker: mockedEventBroker,
		}

		updatedContender := mockedContender
		updatedContender.CompClassID = 3

		contender, err := ucase.UpdateContender(context.Background(), mockedContenderID, updatedContender)

		assert.NoError(t, err)
		assert.Equal(t, 3, contender.CompClassID)
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mock.Anything).
			Return(domain.NilRole, domain.ErrNoOwnership)

		ucase := usecases.ContenderUseCase{
			Repo: makeMockedRepo(domain.Contender{
				ID:        mockedContenderID,
				Ownership: mockedOwnership,
			}),
			Authorizer: mockedAuthorizer,
		}

		contender, err := ucase.UpdateContender(context.Background(), mockedContenderID, domain.Contender{})

		assert.ErrorIs(t, err, domain.ErrNoOwnership)
		assert.Empty(t, contender)
	})
}
