package usecases_test

import (
	"context"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
	"github.com/stretchr/testify/mock"
	"uuid"
)

var errMock = errors.New("mock error")

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

func (m *repositoryMock) GetContendersByCategory(ctx context.Context, tx domain.Transaction, categoryID domain.CategoryID) ([]domain.Contender, error) {
	args := m.Called(ctx, tx, categoryID)
	return args.Get(0).([]domain.Contender), args.Error(1)
}

func (m *repositoryMock) GetContendersByCompetition(ctx context.Context, tx domain.Transaction, competitionID domain.CompetitionID) ([]domain.Contender, error) {
	args := m.Called(ctx, tx, competitionID)
	return args.Get(0).([]domain.Contender), args.Error(1)
}

func (m *repositoryMock) StoreContender(ctx context.Context, tx domain.Transaction, contender domain.Contender) (domain.Contender, error) {
	args := m.Called(ctx, tx, contender)
	return args.Get(0).(domain.Contender), args.Error(1)
}

func (m *repositoryMock) DeleteContender(ctx context.Context, tx domain.Transaction, contenderID domain.ContenderID) error {
	args := m.Called(ctx, tx, contenderID)
	return args.Error(0)
}

func (m *repositoryMock) GetScrubEligibleContenders(ctx context.Context, deadline time.Time) ([]domain.Contender, error) {
	args := m.Called(ctx, deadline)
	return args.Get(0).([]domain.Contender), args.Error(1)
}

func (m *repositoryMock) GetCompetition(ctx context.Context, tx domain.Transaction, competitionID domain.CompetitionID) (domain.Competition, error) {
	args := m.Called(ctx, tx, competitionID)
	return args.Get(0).(domain.Competition), args.Error(1)
}

func (m *repositoryMock) GetAllCompetitions(ctx context.Context, tx domain.Transaction) ([]domain.Competition, error) {
	args := m.Called(ctx, tx)
	return args.Get(0).([]domain.Competition), args.Error(1)
}

func (m *repositoryMock) StoreCompetition(ctx context.Context, tx domain.Transaction, competition domain.Competition) (domain.Competition, error) {
	args := m.Called(ctx, tx, competition)
	return args.Get(0).(domain.Competition), args.Error(1)
}

func (m *repositoryMock) DeleteCompetition(ctx context.Context, tx domain.Transaction, competitionID domain.CompetitionID) error {
	args := m.Called(ctx, tx, competitionID)
	return args.Error(0)
}

func (m *repositoryMock) GetCompetitionsByOrganizer(ctx context.Context, tx domain.Transaction, organizerID domain.OrganizerID) ([]domain.Competition, error) {
	args := m.Called(ctx, tx, organizerID)
	return args.Get(0).([]domain.Competition), args.Error(1)
}

func (m *repositoryMock) GetCategory(ctx context.Context, tx domain.Transaction, categoryID domain.CategoryID) (domain.Category, error) {
	args := m.Called(ctx, tx, categoryID)
	return args.Get(0).(domain.Category), args.Error(1)
}

func (m *repositoryMock) GetCategoriesByCompetition(ctx context.Context, tx domain.Transaction, competitionID domain.CompetitionID) ([]domain.Category, error) {
	args := m.Called(ctx, tx, competitionID)
	return args.Get(0).([]domain.Category), args.Error(1)
}

func (m *repositoryMock) DeleteCategory(ctx context.Context, tx domain.Transaction, categoryID domain.CategoryID) error {
	args := m.Called(ctx, tx, categoryID)
	return args.Error(0)
}

func (m *repositoryMock) StoreCategory(ctx context.Context, tx domain.Transaction, category domain.Category) (domain.Category, error) {
	args := m.Called(ctx, tx, category)
	return args.Get(0).(domain.Category), args.Error(1)
}

func (m *repositoryMock) GetNumberOfContenders(ctx context.Context, tx domain.Transaction, competitionID domain.CompetitionID) (int, error) {
	args := m.Called(ctx, tx, competitionID)
	return args.Get(0).(int), args.Error(1)
}

func (m *repositoryMock) GetProblemsByCompetition(ctx context.Context, tx domain.Transaction, competitionID domain.CompetitionID) ([]domain.Problem, error) {
	args := m.Called(ctx, tx, competitionID)
	return args.Get(0).([]domain.Problem), args.Error(1)
}

func (m *repositoryMock) GetTicksByContender(ctx context.Context, tx domain.Transaction, contenderID domain.ContenderID) ([]domain.Tick, error) {
	args := m.Called(ctx, tx, contenderID)
	return args.Get(0).([]domain.Tick), args.Error(1)
}

func (m *repositoryMock) GetTicksByCompetition(ctx context.Context, tx domain.Transaction, competitionID domain.CompetitionID) ([]domain.Tick, error) {
	args := m.Called(ctx, tx, competitionID)
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

func (m *repositoryMock) GetTickByContenderAndProblem(ctx context.Context, tx domain.Transaction, contenderID domain.ContenderID, problemID domain.ProblemID) (domain.Tick, error) {
	args := m.Called(ctx, tx, contenderID, problemID)
	return args.Get(0).(domain.Tick), args.Error(1)
}

func (m *repositoryMock) StoreTick(ctx context.Context, tx domain.Transaction, tick domain.Tick) (domain.Tick, error) {
	args := m.Called(ctx, tx, tick)
	return args.Get(0).(domain.Tick), args.Error(1)
}

func (m *repositoryMock) GetTicksByProblem(ctx context.Context, tx domain.Transaction, problemID domain.ProblemID) ([]domain.Tick, error) {
	args := m.Called(ctx, tx, problemID)
	return args.Get(0).([]domain.Tick), args.Error(1)
}

func (m *repositoryMock) GetProblem(ctx context.Context, tx domain.Transaction, problemID domain.ProblemID) (domain.Problem, error) {
	args := m.Called(ctx, tx, problemID)
	return args.Get(0).(domain.Problem), args.Error(1)
}

func (m *repositoryMock) GetProblemByNumber(ctx context.Context, tx domain.Transaction, competitionID domain.CompetitionID, problemNumber int) (domain.Problem, error) {
	args := m.Called(ctx, tx, competitionID, problemNumber)
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

func (m *repositoryMock) StoreOrganizer(ctx context.Context, tx domain.Transaction, organizer domain.Organizer) (domain.Organizer, error) {
	args := m.Called(ctx, tx, organizer)
	return args.Get(0).(domain.Organizer), args.Error(1)
}

func (m *repositoryMock) GetRaffle(ctx context.Context, tx domain.Transaction, raffleID domain.RaffleID) (domain.Raffle, error) {
	args := m.Called(ctx, tx, raffleID)
	return args.Get(0).(domain.Raffle), args.Error(1)
}

func (m *repositoryMock) GetRafflesByCompetition(ctx context.Context, tx domain.Transaction, competitionID domain.CompetitionID) ([]domain.Raffle, error) {
	args := m.Called(ctx, tx, competitionID)
	return args.Get(0).([]domain.Raffle), args.Error(1)
}

func (m *repositoryMock) StoreRaffle(ctx context.Context, tx domain.Transaction, raffle domain.Raffle) (domain.Raffle, error) {
	args := m.Called(ctx, tx, raffle)
	return args.Get(0).(domain.Raffle), args.Error(1)
}

func (m *repositoryMock) DeleteRaffle(ctx context.Context, tx domain.Transaction, raffleID domain.RaffleID) error {
	args := m.Called(ctx, tx, raffleID)
	return args.Error(0)
}

func (m *repositoryMock) GetRaffleWinners(ctx context.Context, tx domain.Transaction, raffleID domain.RaffleID) ([]domain.RaffleWinner, error) {
	args := m.Called(ctx, tx, raffleID)
	return args.Get(0).([]domain.RaffleWinner), args.Error(1)
}

func (m *repositoryMock) StoreRaffleWinner(ctx context.Context, tx domain.Transaction, winner domain.RaffleWinner) (domain.RaffleWinner, error) {
	args := m.Called(ctx, tx, winner)
	return args.Get(0).(domain.RaffleWinner), args.Error(1)
}

func (m *repositoryMock) DeleteRaffleWinner(ctx context.Context, tx domain.Transaction, raffleWinnerID domain.RaffleWinnerID) error {
	args := m.Called(ctx, tx, raffleWinnerID)
	return args.Error(0)
}

func (m *repositoryMock) GetUserByUsername(ctx context.Context, tx domain.Transaction, username string) (domain.User, error) {
	args := m.Called(ctx, tx, username)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *repositoryMock) GetUsersByOrganizer(ctx context.Context, tx domain.Transaction, organizerID domain.OrganizerID) ([]domain.User, error) {
	args := m.Called(ctx, tx, organizerID)
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *repositoryMock) GetAllOrganizers(ctx context.Context, tx domain.Transaction) ([]domain.Organizer, error) {
	args := m.Called(ctx, tx)
	return args.Get(0).([]domain.Organizer), args.Error(1)
}

func (m *repositoryMock) AddUserToOrganizer(ctx context.Context, tx domain.Transaction, userID domain.UserID, organizerID domain.OrganizerID) error {
	args := m.Called(ctx, tx, userID, organizerID)
	return args.Error(0)
}

func (m *repositoryMock) DeleteOrganizerInvite(ctx context.Context, tx domain.Transaction, inviteID domain.OrganizerInviteID) error {
	args := m.Called(ctx, tx, inviteID)
	return args.Error(0)
}

func (m *repositoryMock) GetOrganizerInvite(ctx context.Context, tx domain.Transaction, inviteID domain.OrganizerInviteID) (domain.OrganizerInvite, error) {
	args := m.Called(ctx, tx, inviteID)
	return args.Get(0).(domain.OrganizerInvite), args.Error(1)
}

func (m *repositoryMock) GetOrganizerInvitesByOrganizer(ctx context.Context, tx domain.Transaction, organizerID domain.OrganizerID) ([]domain.OrganizerInvite, error) {
	args := m.Called(ctx, tx, organizerID)
	return args.Get(0).([]domain.OrganizerInvite), args.Error(1)
}

func (m *repositoryMock) StoreOrganizerInvite(ctx context.Context, tx domain.Transaction, invite domain.OrganizerInvite) error {
	args := m.Called(ctx, tx, invite)
	return args.Error(0)
}

func (m *repositoryMock) StoreResult(ctx context.Context, tx domain.Transaction, result domain.Result) error {
	args := m.Called(ctx, tx, result)
	return args.Error(0)
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

type resultKeeperMock struct {
	mock.Mock
}

func (m *resultKeeperMock) GetResult(contenderID domain.ContenderID) (domain.Result, error) {
	args := m.Called(contenderID)
	return args.Get(0).(domain.Result), args.Error(1)
}

type eventBrokerMock struct {
	mock.Mock
}

func (m *eventBrokerMock) Dispatch(competitionID domain.CompetitionID, event any) {
	m.Called(competitionID, event)
}

func (m *eventBrokerMock) Subscribe(filters []domain.EventFilter, bufferCapacity int) (domain.SubscriptionID, domain.EventReader) {
	args := m.Called(filters, bufferCapacity)
	return args.Get(0).(domain.SubscriptionID), args.Get(1).(domain.EventReader)
}

func (m *eventBrokerMock) Unsubscribe(subscriptionID domain.SubscriptionID) {
	m.Called(subscriptionID)
}

type uuidGeneratorMock struct {
	mock.Mock
}

func (m *uuidGeneratorMock) Generate() uuid.UUID {
	args := m.Called()
	return args.Get(0).(uuid.UUID)
}

type pointValueKeeperMock struct {
	mock.Mock
}

func (m *pointValueKeeperMock) GetPointValues(contenderID domain.ContenderID) []domain.PointValue {
	args := m.Called(contenderID)
	return args.Get(0).([]domain.PointValue)
}
