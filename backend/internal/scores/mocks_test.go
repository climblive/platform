package scores_test

import (
	"iter"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/scores"
	"github.com/stretchr/testify/mock"
)

type engineStoreMock struct {
	mock.Mock
}

func (m *engineStoreMock) GetRules() scores.Rules {
	args := m.Called()
	return args.Get(0).(scores.Rules)
}

func (m *engineStoreMock) SaveRules(rules scores.Rules) {
	m.Called(rules)
}

func (m *engineStoreMock) GetContender(contenderID domain.ContenderID) (scores.Contender, bool) {
	args := m.Called(contenderID)
	return args.Get(0).(scores.Contender), args.Bool(1)
}

func (m *engineStoreMock) GetContendersByCompClass(compClassID domain.CompClassID) iter.Seq[scores.Contender] {
	args := m.Called(compClassID)
	return args.Get(0).(iter.Seq[scores.Contender])
}

func (m *engineStoreMock) SaveContender(contender scores.Contender) {
	m.Called(contender)
}

func (m *engineStoreMock) GetCompClassIDs() []domain.CompClassID {
	args := m.Called()
	return args.Get(0).([]domain.CompClassID)
}

func (m *engineStoreMock) GetTicksByContender(contenderID domain.ContenderID) iter.Seq[scores.Tick] {
	args := m.Called(contenderID)
	return args.Get(0).(iter.Seq[scores.Tick])
}

func (m *engineStoreMock) GetTick(contenderID domain.ContenderID, problemID domain.ProblemID) (scores.Tick, bool) {
	args := m.Called(contenderID, problemID)
	return args.Get(0).(scores.Tick), args.Bool(1)
}

func (m *engineStoreMock) GetTicksByProblem(compClassID domain.CompClassID, problemID domain.ProblemID) iter.Seq[scores.Tick] {
	args := m.Called(compClassID, problemID)
	return args.Get(0).(iter.Seq[scores.Tick])
}

func (m *engineStoreMock) SaveTick(contenderID domain.ContenderID, tick scores.Tick) {
	m.Called(contenderID, tick)
}

func (m *engineStoreMock) DeleteTick(contenderID domain.ContenderID, problemID domain.ProblemID) {
	m.Called(contenderID, problemID)
}

func (m *engineStoreMock) GetProblem(problemID domain.ProblemID) (scores.Problem, bool) {
	args := m.Called(problemID)
	return args.Get(0).(scores.Problem), args.Bool(1)
}

func (m *engineStoreMock) SaveProblem(problem scores.Problem) {
	m.Called(problem)
}

func (m *engineStoreMock) GetAllProblems() iter.Seq[scores.Problem] {
	args := m.Called()
	return args.Get(0).(iter.Seq[scores.Problem])
}

func (m *engineStoreMock) GetPointValue(contenderID domain.ContenderID, problemID domain.ProblemID) (domain.PointValue, bool) {
	args := m.Called(contenderID, problemID)
	return args.Get(0).(domain.PointValue), args.Bool(1)
}

func (m *engineStoreMock) SavePointValue(contenderID domain.ContenderID, problemID domain.ProblemID, value domain.PointValue) {
	m.Called(contenderID, problemID, value)
}

func (m *engineStoreMock) GetDirtyPointValues() []domain.PointValue {
	args := m.Called()
	return args.Get(0).([]domain.PointValue)
}

func (m *engineStoreMock) SaveScore(score domain.Score) {
	m.Called(score)
}

func (m *engineStoreMock) GetDirtyScores() []domain.Score {
	args := m.Called()
	return args.Get(0).([]domain.Score)
}
