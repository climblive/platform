package usecases_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/usecases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetContender(t *testing.T) {
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
		On("GetContender", mock.Anything, mock.Anything, mockedContenderID).
		Return(mockedContender, nil)

	mockedScoreKeeper.On("GetScore", mockedContenderID).Return(domain.Score{
		Timestamp:   time.Now().Truncate(time.Hour),
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

		assert.NoError(t, err)

		assert.Equal(t, mockedContenderID, contender.ID)
		assert.Equal(t, 1000, contender.Score)
		assert.Equal(t, 5, contender.Placement)
		assert.Equal(t, time.Now().Truncate(time.Hour), *contender.ScoreUpdated)

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

		assert.NoError(t, err)
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
		mockedAuthorizer := new(authorizerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mockedOwnership).
			Return(domain.ContenderRole, nil)

		ucase := usecases.ContenderUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		err := ucase.DeleteContender(context.Background(), mockedContenderID)

		assert.ErrorIs(t, err, domain.ErrInsufficientRole)
	})
}

func TestCreateContenders(t *testing.T) {
	mockedContestID := domain.ResourceID(1)
	mockedOwnership := domain.OwnershipData{
		OrganizerID: 1,
	}

	mockedRepo := new(repositoryMock)
	mockedTx := new(transactionMock)
	mockedCodeGenerator := new(codeGeneratorMock)

	mockedRepo.
		On("GetContest", mock.Anything, mock.Anything, mockedContestID).
		Return(domain.Contest{
			ID:        mockedContestID,
			Ownership: mockedOwnership,
		}, nil)

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

	t.Run("HappyPath", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mockedOwnership).
			Return(domain.OrganizerRole, nil)

		ucase := usecases.ContenderUseCase{
			Repo:                      mockedRepo,
			Authorizer:                mockedAuthorizer,
			RegistrationCodeGenerator: mockedCodeGenerator,
		}

		contenders, err := ucase.CreateContenders(context.Background(), mockedContestID, 100)

		assert.NoError(t, err)
		assert.Len(t, contenders, 100)

		mockedRepo.AssertExpectations(t)
		mockedTx.AssertNumberOfCalls(t, "Commit", 1)
		mockedTx.AssertNotCalled(t, "Rollback")

		for idx, contender := range contenders {
			assert.Equal(t, fmt.Sprintf("%08d", idx), contender.RegistrationCode)
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

		contender, err := ucase.CreateContenders(context.Background(), mockedContestID, 100)

		assert.ErrorIs(t, err, domain.ErrNoOwnership)
		assert.Empty(t, contender)
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

	mockedContender := domain.Contender{
		ID:                  mockedContenderID,
		Ownership:           mockedOwnership,
		ContestID:           1,
		CompClassID:         1,
		RegistrationCode:    "ABCD1234",
		Name:                "John Doe",
		PublicName:          "John Doe",
		ClubName:            "Testers' Climbing Club",
		Entered:             &currentTime,
		WithdrawnFromFinals: false,
		Disqualified:        false,
		Score:               100,
		Placement:           1,
		ScoreUpdated:        &currentTime,
	}

	mockedContest := domain.Contest{
		ID:          1,
		GracePeriod: 15 * time.Minute,
	}

	mockedRepo := new(repositoryMock)
	mockedScoreKeeper := new(scoreKeeperMock)

	mockedRepo.
		On("GetContender", mock.Anything, mock.Anything, mockedContenderID).
		Return(mockedContender, nil)

	mockedRepo.
		On("GetContest", mock.Anything, mock.Anything, 1).
		Return(mockedContest, nil)

	mockedRepo.
		On("StoreContender", mock.Anything, mock.Anything, mock.AnythingOfType("domain.Contender")).
		Return(mirrorInstruction{}, nil)

	mockedScoreKeeper.On("GetScore", mockedContenderID).Return(domain.Score{
		Timestamp:   currentTime.Truncate(time.Hour),
		ContenderID: mockedContenderID,
		Score:       1000,
		Placement:   5,
	}, nil)

	t.Run("Identity", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mock.Anything).
			Return(domain.ContenderRole, nil)

		ucase := usecases.ContenderUseCase{
			Repo:        mockedRepo,
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
		assert.Equal(t, currentTime.Truncate(time.Hour), *contender.ScoreUpdated)
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedAuthorizer := new(authorizerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mock.Anything).
			Return(domain.NilRole, domain.ErrNoOwnership)

		ucase := usecases.ContenderUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		contender, err := ucase.UpdateContender(context.Background(), mockedContenderID, domain.Contender{})

		assert.ErrorIs(t, err, domain.ErrNoOwnership)
		assert.Empty(t, contender)
	})
}

var errMock = errors.New("mock error")

type mirrorInstruction struct{}

type transactionMock struct {
	mock.Mock
}

func (m *transactionMock) Commit() error {
	args := m.Called()
	return args.Error(0)
}

func (m *transactionMock) Rollback() {
	m.Called()
}

type codeGeneratorMock struct {
	mock.Mock
}

func (m *codeGeneratorMock) Generate(length int) string {
	args := m.Called(length)
	return args.Get(0).(string)
}

type repositoryMock struct {
	mock.Mock
}

func (m *repositoryMock) Begin() domain.Transaction {
	args := m.Called()
	return args.Get(0).(domain.Transaction)
}

func (m *repositoryMock) GetContender(ctx context.Context, tx domain.Transaction, contenderID domain.ResourceID) (domain.Contender, error) {
	args := m.Called(ctx, tx, contenderID)
	return args.Get(0).(domain.Contender), args.Error(1)
}

func (m *repositoryMock) GetContenderByCode(ctx context.Context, tx domain.Transaction, registrationCode string) (domain.Contender, error) {
	args := m.Called(ctx, tx, registrationCode)
	return args.Get(0).(domain.Contender), args.Error(1)
}

func (m *repositoryMock) GetContendersByCompClass(ctx context.Context, tx domain.Transaction, compClassID domain.ResourceID) ([]domain.Contender, error) {
	args := m.Called(ctx, tx, compClassID)
	return args.Get(0).([]domain.Contender), args.Error(1)
}

func (m *repositoryMock) GetContendersByContest(ctx context.Context, tx domain.Transaction, contestID domain.ResourceID) ([]domain.Contender, error) {
	args := m.Called(ctx, tx, contestID)
	return args.Get(0).([]domain.Contender), args.Error(1)
}

func (m *repositoryMock) StoreContender(ctx context.Context, tx domain.Transaction, contender domain.Contender) (domain.Contender, error) {
	args := m.Called(ctx, tx, contender)

	if _, ok := args.Get(0).(mirrorInstruction); ok {
		return contender, nil
	} else {
		return args.Get(0).(domain.Contender), args.Error(1)
	}
}

func (m *repositoryMock) DeleteContender(ctx context.Context, tx domain.Transaction, contenderID domain.ResourceID) error {
	args := m.Called(ctx, tx, contenderID)
	return args.Error(0)
}

func (m *repositoryMock) GetContest(ctx context.Context, tx domain.Transaction, contestID domain.ResourceID) (domain.Contest, error) {
	args := m.Called(ctx, tx, contestID)
	return args.Get(0).(domain.Contest), args.Error(1)
}

func (m *repositoryMock) GetCompClass(ctx context.Context, tx domain.Transaction, compClassID domain.ResourceID) (domain.CompClass, error) {
	args := m.Called(ctx, tx, compClassID)
	return args.Get(0).(domain.CompClass), args.Error(1)
}

type authorizerMock struct {
	mock.Mock
}

func (m *authorizerMock) HasOwnership(ctx context.Context, resourceOwnership domain.OwnershipData) (domain.AuthRole, error) {
	args := m.Called(ctx, resourceOwnership)
	return args.Get(0).(domain.AuthRole), args.Error(1)
}

type scoreKeeperMock struct {
	mock.Mock
}

func (m *scoreKeeperMock) UpdateScore(contenderID domain.ResourceID, score domain.Score) error {
	args := m.Called(contenderID, score)
	return args.Error(0)
}

func (m *scoreKeeperMock) GetScore(contenderID domain.ResourceID) (domain.Score, error) {
	args := m.Called(contenderID)
	return args.Get(0).(domain.Score), args.Error(1)
}
