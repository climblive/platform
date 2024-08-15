package usecases_test

import (
	"context"
	"testing"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/usecases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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
	return args.Error(1)
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

func TestGetContender_BadCredentials(t *testing.T) {
	mockRepo := new(mockRepository)
	mockAuthorizer := new(mockAuthorizer)
	contenderOwnership := domain.OwnershipData{}

	mockAuthorizer.
		On("HasOwnership", mock.Anything, contenderOwnership).
		Return(domain.NilRole, domain.ErrPermissionDenied)

	mockRepo.
		On("GetContender", mock.Anything, mock.Anything, 1).
		Return(domain.Contender{Ownership: contenderOwnership}, nil)

	ucase := usecases.ContenderUseCase{
		Repo:        mockRepo,
		Authorizer:  mockAuthorizer,
		ScoreKeeper: &mockScoreKeeper{},
	}

	contender, err := ucase.GetContender(context.Background(), 1)

	assert.ErrorIs(t, err, domain.ErrPermissionDenied)
	assert.Empty(t, contender)
}
