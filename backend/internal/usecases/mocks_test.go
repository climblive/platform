package usecases_test

import (
	"context"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
	"github.com/stretchr/testify/mock"
)

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

func (m *repositoryMock) GetCompClassesByContest(ctx context.Context, tx domain.Transaction, contestID domain.ResourceID) ([]domain.CompClass, error) {
	args := m.Called(ctx, tx, contestID)
	return args.Get(0).([]domain.CompClass), args.Error(1)
}

func (m *repositoryMock) GetNumberOfContenders(ctx context.Context, tx domain.Transaction, contestID domain.ResourceID) (int, error) {
	args := m.Called(ctx, tx, contestID)
	return args.Get(0).(int), args.Error(1)
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

type eventBrokerMock struct {
	mock.Mock
}

func (m *eventBrokerMock) Dispatch(contestID domain.ResourceID, event any) {
	m.Called(contestID, event)
}
