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
	mockContenderID := domain.ResourceID(1)
	mockOwnership := domain.OwnershipData{
		OrganizerID: 1,
		ContenderID: &mockContenderID,
	}

	mockContender := domain.Contender{
		ID:        mockContenderID,
		Ownership: mockOwnership,
	}

	mockRepo := new(mockRepository)
	mockScoreKeeper := new(mockScoreKeeper)

	mockRepo.
		On("GetContender", mock.Anything, mock.Anything, mockContenderID).
		Return(mockContender, nil)

	mockScoreKeeper.On("GetScore", mockContenderID).Return(domain.Score{
		Timestamp:   time.Now().Truncate(time.Hour),
		ContenderID: mockContenderID,
		Score:       1000,
		Placement:   5,
	}, nil)

	t.Run("HappyPath", func(t *testing.T) {
		mockAuthorizer := new(mockAuthorizer)

		mockAuthorizer.
			On("HasOwnership", mock.Anything, mockOwnership).
			Return(domain.ContenderRole, nil)

		ucase := usecases.ContenderUseCase{
			Repo:        mockRepo,
			Authorizer:  mockAuthorizer,
			ScoreKeeper: mockScoreKeeper,
		}

		contender, err := ucase.GetContender(context.Background(), mockContenderID)

		assert.NoError(t, err)

		assert.Equal(t, mockContenderID, contender.ID)
		assert.Equal(t, 1000, contender.Score)
		assert.Equal(t, 5, contender.Placement)
		assert.Equal(t, time.Now().Truncate(time.Hour), *contender.ScoreUpdated)

	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockAuthorizer := new(mockAuthorizer)

		mockAuthorizer.
			On("HasOwnership", mock.Anything, mockOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		ucase := usecases.ContenderUseCase{
			Repo:        mockRepo,
			Authorizer:  mockAuthorizer,
			ScoreKeeper: mockScoreKeeper,
		}

		contender, err := ucase.GetContender(context.Background(), mockContenderID)

		assert.ErrorIs(t, err, domain.ErrNoOwnership)
		assert.Empty(t, contender)
	})
}

func TestDeleteContender(t *testing.T) {
	mockContenderID := domain.ResourceID(1)
	mockOwnership := domain.OwnershipData{
		OrganizerID: 1,
		ContenderID: &mockContenderID,
	}

	mockRepo := new(mockRepository)

	mockRepo.
		On("GetContender", mock.Anything, mock.Anything, mockContenderID).
		Return(domain.Contender{Ownership: mockOwnership}, nil)

	mockRepo.
		On("DeleteContender", mock.Anything, mock.Anything, mockContenderID).
		Return(nil)

	t.Run("HappyPath", func(t *testing.T) {
		mockAuthorizer := new(mockAuthorizer)

		mockAuthorizer.
			On("HasOwnership", mock.Anything, mockOwnership).
			Return(domain.OrganizerRole, nil)

		ucase := usecases.ContenderUseCase{
			Repo:       mockRepo,
			Authorizer: mockAuthorizer,
		}

		err := ucase.DeleteContender(context.Background(), mockContenderID)

		assert.NoError(t, err)
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockAuthorizer := new(mockAuthorizer)

		mockAuthorizer.
			On("HasOwnership", mock.Anything, mockOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		ucase := usecases.ContenderUseCase{
			Repo:       mockRepo,
			Authorizer: mockAuthorizer,
		}

		err := ucase.DeleteContender(context.Background(), mockContenderID)

		assert.ErrorIs(t, err, domain.ErrNoOwnership)
	})

	t.Run("InsufficientRole", func(t *testing.T) {
		mockAuthorizer := new(mockAuthorizer)

		mockAuthorizer.
			On("HasOwnership", mock.Anything, mockOwnership).
			Return(domain.ContenderRole, nil)

		ucase := usecases.ContenderUseCase{
			Repo:       mockRepo,
			Authorizer: mockAuthorizer,
		}

		err := ucase.DeleteContender(context.Background(), mockContenderID)

		assert.ErrorIs(t, err, domain.ErrInsufficientRole)
	})
}

func TestCreateContenders(t *testing.T) {
	mockContestID := domain.ResourceID(1)
	mockOwnership := domain.OwnershipData{
		OrganizerID: 1,
	}

	mockRepo := new(mockRepository)
	mockTx := new(mockTransaction)
	mockedCodeGenerator := new(codeGeneratorMock)

	mockRepo.
		On("GetContest", mock.Anything, mock.Anything, mockContestID).
		Return(domain.Contest{
			ID:        mockContestID,
			Ownership: mockOwnership,
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

		mockRepo.
			On("StoreContender", mock.Anything, mock.Anything, contender).
			Return(contender, nil)
	}

	mockRepo.
		On("Begin").
		Return(mockTx, nil)

	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return()

	t.Run("HappyPath", func(t *testing.T) {
		mockAuthorizer := new(mockAuthorizer)

		mockAuthorizer.
			On("HasOwnership", mock.Anything, mockOwnership).
			Return(domain.OrganizerRole, nil)

		ucase := usecases.ContenderUseCase{
			Repo:                      mockRepo,
			Authorizer:                mockAuthorizer,
			RegistrationCodeGenerator: mockedCodeGenerator,
		}

		contenders, err := ucase.CreateContenders(context.Background(), mockContestID, 100)

		assert.NoError(t, err)
		assert.Len(t, contenders, 100)

		mockRepo.AssertExpectations(t)
		mockTx.AssertNumberOfCalls(t, "Commit", 1)
		mockTx.AssertNotCalled(t, "Rollback")

		for idx, contender := range contenders {
			assert.Equal(t, fmt.Sprintf("%08d", idx), contender.RegistrationCode)
		}
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockAuthorizer := new(mockAuthorizer)

		mockAuthorizer.
			On("HasOwnership", mock.Anything, mockOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		ucase := usecases.ContenderUseCase{
			Repo:       mockRepo,
			Authorizer: mockAuthorizer,
		}

		contender, err := ucase.CreateContenders(context.Background(), mockContestID, 100)

		assert.ErrorIs(t, err, domain.ErrNoOwnership)
		assert.Empty(t, contender)
	})
}

func TestCreateContenders_Rollback(t *testing.T) {
	mockContestID := domain.ResourceID(1)

	mockRepo := new(mockRepository)
	mockTx := new(mockTransaction)
	mockAuthorizer := new(mockAuthorizer)
	mockedCodeGenerator := new(codeGeneratorMock)

	mockRepo.
		On("GetContest", mock.Anything, mock.Anything, mockContestID).
		Return(domain.Contest{}, nil)

	mockRepo.
		On("Begin").
		Return(mockTx, nil)

	mockRepo.
		On("StoreContender", mock.Anything, mock.Anything, mock.Anything).
		Return(domain.Contender{}, errMock)

	mockTx.On("Rollback").Return()

	mockAuthorizer.
		On("HasOwnership", mock.Anything, domain.OwnershipData{}).
		Return(domain.AdminRole, nil)

	mockedCodeGenerator.
		On("Generate", 8).
		Return("DEAFBEEF")

	ucase := usecases.ContenderUseCase{
		Repo:                      mockRepo,
		Authorizer:                mockAuthorizer,
		RegistrationCodeGenerator: mockedCodeGenerator,
	}

	contenders, err := ucase.CreateContenders(context.Background(), mockContestID, 1)

	assert.Error(t, err)
	assert.Nil(t, contenders)

	mockRepo.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

var errMock = errors.New("mock error")

type mockTransaction struct {
	mock.Mock
}

func (m *mockTransaction) Commit() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockTransaction) Rollback() {
	m.Called()
}

type codeGeneratorMock struct {
	mock.Mock
}

func (m *codeGeneratorMock) Generate(length int) string {
	args := m.Called(length)
	return args.Get(0).(string)
}

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) Begin() domain.Transaction {
	args := m.Called()
	return args.Get(0).(domain.Transaction)
}

func (m *mockRepository) GetContender(ctx context.Context, tx domain.Transaction, contenderID domain.ResourceID) (domain.Contender, error) {
	args := m.Called(ctx, tx, contenderID)
	return args.Get(0).(domain.Contender), args.Error(1)
}

func (m *mockRepository) GetContenderByCode(ctx context.Context, tx domain.Transaction, registrationCode string) (domain.Contender, error) {
	args := m.Called(ctx, tx, registrationCode)
	return args.Get(0).(domain.Contender), args.Error(1)
}

func (m *mockRepository) GetContendersByCompClass(ctx context.Context, tx domain.Transaction, compClassID domain.ResourceID) ([]domain.Contender, error) {
	args := m.Called(ctx, tx, compClassID)
	return args.Get(0).([]domain.Contender), args.Error(1)
}

func (m *mockRepository) GetContendersByContest(ctx context.Context, tx domain.Transaction, contestID domain.ResourceID) ([]domain.Contender, error) {
	args := m.Called(ctx, tx, contestID)
	return args.Get(0).([]domain.Contender), args.Error(1)
}

func (m *mockRepository) StoreContender(ctx context.Context, tx domain.Transaction, contender domain.Contender) (domain.Contender, error) {
	args := m.Called(ctx, tx, contender)
	return args.Get(0).(domain.Contender), args.Error(1)
}

func (m *mockRepository) DeleteContender(ctx context.Context, tx domain.Transaction, contenderID domain.ResourceID) error {
	args := m.Called(ctx, tx, contenderID)
	return args.Error(0)
}

func (m *mockRepository) GetContest(ctx context.Context, tx domain.Transaction, contestID domain.ResourceID) (domain.Contest, error) {
	args := m.Called(ctx, tx, contestID)
	return args.Get(0).(domain.Contest), args.Error(1)
}

func (m *mockRepository) GetCompClass(ctx context.Context, tx domain.Transaction, compClassID domain.ResourceID) (domain.CompClass, error) {
	args := m.Called(ctx, tx, compClassID)
	return args.Get(0).(domain.CompClass), args.Error(1)
}

type mockAuthorizer struct {
	mock.Mock
}

func (m *mockAuthorizer) HasOwnership(ctx context.Context, resourceOwnership domain.OwnershipData) (domain.AuthRole, error) {
	args := m.Called(ctx, resourceOwnership)
	return args.Get(0).(domain.AuthRole), args.Error(1)
}

type mockScoreKeeper struct {
	mock.Mock
}

func (m *mockScoreKeeper) UpdateScore(contenderID domain.ResourceID, score domain.Score) error {
	args := m.Called(contenderID, score)
	return args.Error(0)
}

func (m *mockScoreKeeper) GetScore(contenderID domain.ResourceID) (domain.Score, error) {
	args := m.Called(contenderID)
	return args.Get(0).(domain.Score), args.Error(1)
}
