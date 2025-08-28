package usecases_test

import (
	"context"
	"math/rand"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
	"github.com/stretchr/testify/mock"
)

var errMock = errors.New("mock error")

func randomResourceID[T domain.ResourceIDType]() T {
	return T(rand.Int())
}

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

func (m *repositoryMock) Begin() (domain.Transaction, error) {
	args := m.Called()
	return args.Get(0).(domain.Transaction), args.Error(1)
}

func (m *repositoryMock) GetContender(ctx context.Context, tx domain.Transaction, contenderID domain.ContenderID) (domain.Contender, error) {
	args := m.Called(ctx, tx, contenderID)
	return args.Get(0).(domain.Contender), args.Error(1)
}

func (m *repositoryMock) GetContenderByCode(ctx context.Context, tx domain.Transaction, registrationCode string) (domain.Contender, error) {
	args := m.Called(ctx, tx, registrationCode)
	return args.Get(0).(domain.Contender), args.Error(1)
}

func (m *repositoryMock) GetContendersByCompClass(ctx context.Context, tx domain.Transaction, compClassID domain.CompClassID) ([]domain.Contender, error) {
	args := m.Called(ctx, tx, compClassID)
	return args.Get(0).([]domain.Contender), args.Error(1)
}

func (m *repositoryMock) GetContendersByContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) ([]domain.Contender, error) {
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

func (m *repositoryMock) DeleteContender(ctx context.Context, tx domain.Transaction, contenderID domain.ContenderID) error {
	args := m.Called(ctx, tx, contenderID)
	return args.Error(0)
}

func (m *repositoryMock) GetContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) (domain.Contest, error) {
	args := m.Called(ctx, tx, contestID)
	return args.Get(0).(domain.Contest), args.Error(1)
}

func (m *repositoryMock) StoreContest(ctx context.Context, tx domain.Transaction, contest domain.Contest) (domain.Contest, error) {
	args := m.Called(ctx, tx, contest)

	if _, ok := args.Get(0).(mirrorInstruction); ok {
		return contest, nil
	} else {
		return args.Get(0).(domain.Contest), args.Error(1)
	}
}

func (m *repositoryMock) GetContestsByOrganizer(ctx context.Context, tx domain.Transaction, organizerID domain.OrganizerID) ([]domain.Contest, error) {
	args := m.Called(ctx, tx, organizerID)
	return args.Get(0).([]domain.Contest), args.Error(1)
}

func (m *repositoryMock) GetCompClass(ctx context.Context, tx domain.Transaction, compClassID domain.CompClassID) (domain.CompClass, error) {
	args := m.Called(ctx, tx, compClassID)
	return args.Get(0).(domain.CompClass), args.Error(1)
}

func (m *repositoryMock) GetCompClassesByContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) ([]domain.CompClass, error) {
	args := m.Called(ctx, tx, contestID)
	return args.Get(0).([]domain.CompClass), args.Error(1)
}

func (m *repositoryMock) DeleteCompClass(ctx context.Context, tx domain.Transaction, compClassID domain.CompClassID) error {
	args := m.Called(ctx, tx, compClassID)
	return args.Error(0)
}

func (m *repositoryMock) StoreCompClass(ctx context.Context, tx domain.Transaction, compClass domain.CompClass) (domain.CompClass, error) {
	args := m.Called(ctx, tx, compClass)
	return args.Get(0).(domain.CompClass), args.Error(1)
}

func (m *repositoryMock) GetNumberOfContenders(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) (int, error) {
	args := m.Called(ctx, tx, contestID)
	return args.Get(0).(int), args.Error(1)
}

func (m *repositoryMock) GetProblemsByContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) ([]domain.Problem, error) {
	args := m.Called(ctx, tx, contestID)
	return args.Get(0).([]domain.Problem), args.Error(1)
}

func (m *repositoryMock) GetTicksByContender(ctx context.Context, tx domain.Transaction, contenderID domain.ContenderID) ([]domain.Tick, error) {
	args := m.Called(ctx, tx, contenderID)
	return args.Get(0).([]domain.Tick), args.Error(1)
}

func (m *repositoryMock) GetTicksByContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) ([]domain.Tick, error) {
	args := m.Called(ctx, tx, contestID)
	return args.Get(0).([]domain.Tick), args.Error(1)
}

func (m *repositoryMock) DeleteTick(ctx context.Context, tx domain.Transaction, tickID domain.TickID) error {
	args := m.Called(ctx, tx, tickID)
	return args.Error(0)
}

func (m *repositoryMock) GetTick(ctx context.Context, tx domain.Transaction, tickID domain.TickID) (domain.Tick, error) {
	args := m.Called(ctx, tx, tickID)
	return args.Get(0).(domain.Tick), args.Error(1)
}

func (m *repositoryMock) StoreTick(ctx context.Context, tx domain.Transaction, tick domain.Tick) (domain.Tick, error) {
	args := m.Called(ctx, tx, tick)

	if _, ok := args.Get(0).(mirrorInstruction); ok {
		return tick, nil
	} else {
		return args.Get(0).(domain.Tick), args.Error(1)
	}
}

func (m *repositoryMock) GetTicksByProblem(ctx context.Context, tx domain.Transaction, problemID domain.ProblemID) ([]domain.Tick, error) {
	args := m.Called(ctx, tx, problemID)
	return args.Get(0).([]domain.Tick), args.Error(1)
}

func (m *repositoryMock) GetProblem(ctx context.Context, tx domain.Transaction, problemID domain.ProblemID) (domain.Problem, error) {
	args := m.Called(ctx, tx, problemID)
	return args.Get(0).(domain.Problem), args.Error(1)
}

func (m *repositoryMock) GetProblemByNumber(ctx context.Context, tx domain.Transaction, contestID domain.ContestID, problemNumber int) (domain.Problem, error) {
	args := m.Called(ctx, tx, contestID, problemNumber)
	return args.Get(0).(domain.Problem), args.Error(1)
}

func (m *repositoryMock) StoreProblem(ctx context.Context, tx domain.Transaction, problem domain.Problem) (domain.Problem, error) {
	args := m.Called(ctx, tx, problem)
	return args.Get(0).(domain.Problem), args.Error(1)
}

func (m *repositoryMock) DeleteProblem(ctx context.Context, tx domain.Transaction, problemID domain.ProblemID) error {
	args := m.Called(ctx, tx, problemID)
	return args.Error(0)
}

func (m *repositoryMock) GetOrganizer(ctx context.Context, tx domain.Transaction, organizerID domain.OrganizerID) (domain.Organizer, error) {
	args := m.Called(ctx, tx, organizerID)
	return args.Get(0).(domain.Organizer), args.Error(1)
}

func (m *repositoryMock) GetRaffle(ctx context.Context, tx domain.Transaction, raffleID domain.RaffleID) (domain.Raffle, error) {
	args := m.Called(ctx, tx, raffleID)
	return args.Get(0).(domain.Raffle), args.Error(1)
}

func (m *repositoryMock) GetRafflesByContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) ([]domain.Raffle, error) {
	args := m.Called(ctx, tx, contestID)
	return args.Get(0).([]domain.Raffle), args.Error(1)
}

func (m *repositoryMock) StoreRaffle(ctx context.Context, tx domain.Transaction, raffle domain.Raffle) (domain.Raffle, error) {
	args := m.Called(ctx, tx, raffle)
	return args.Get(0).(domain.Raffle), args.Error(1)
}

func (m *repositoryMock) GetRaffleWinners(ctx context.Context, tx domain.Transaction, raffleID domain.RaffleID) ([]domain.RaffleWinner, error) {
	args := m.Called(ctx, tx, raffleID)
	return args.Get(0).([]domain.RaffleWinner), args.Error(1)
}

func (m *repositoryMock) StoreRaffleWinner(ctx context.Context, tx domain.Transaction, winner domain.RaffleWinner) (domain.RaffleWinner, error) {
	args := m.Called(ctx, tx, winner)

	if _, ok := args.Get(0).(mirrorInstruction); ok {
		return winner, nil
	} else {
		return args.Get(0).(domain.RaffleWinner), args.Error(1)
	}
}

func (m *repositoryMock) GetUserByUsername(ctx context.Context, tx domain.Transaction, username string) (domain.User, error) {
	args := m.Called(ctx, tx, username)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *repositoryMock) GetAllOrganizers(ctx context.Context, tx domain.Transaction) ([]domain.Organizer, error) {
	args := m.Called(ctx, tx)
	return args.Get(0).([]domain.Organizer), args.Error(1)
}

type authorizerMock struct {
	mock.Mock
}

func (m *authorizerMock) HasOwnership(ctx context.Context, resourceOwnership domain.OwnershipData) (domain.AuthRole, error) {
	args := m.Called(ctx, resourceOwnership)
	return args.Get(0).(domain.AuthRole), args.Error(1)
}

func (m *authorizerMock) GetAuthentication(ctx context.Context) (domain.Authentication, error) {
	args := m.Called(ctx)
	return args.Get(0).(domain.Authentication), args.Error(1)
}

type scoreKeeperMock struct {
	mock.Mock
}

func (m *scoreKeeperMock) GetScore(contenderID domain.ContenderID) (domain.Score, error) {
	args := m.Called(contenderID)
	return args.Get(0).(domain.Score), args.Error(1)
}

type eventBrokerMock struct {
	mock.Mock
}

func (m *eventBrokerMock) Dispatch(contestID domain.ContestID, event any) {
	m.Called(contestID, event)
}

func (m *eventBrokerMock) Subscribe(filter domain.EventFilter, bufferCapacity int) (domain.SubscriptionID, domain.EventReader) {
	args := m.Called(filter, bufferCapacity)
	return args.Get(0).(domain.SubscriptionID), args.Get(1).(domain.EventReader)
}

func (m *eventBrokerMock) Unsubscribe(subscriptionID domain.SubscriptionID) {
	m.Called(subscriptionID)
}
