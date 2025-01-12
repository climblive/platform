package usecases_test

import (
	"context"
	"testing"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/scores"
	"github.com/climblive/platform/backend/internal/usecases"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestScoreEngineUseCase(t *testing.T) {
	fakedContestID := randomResourceID[domain.ContestID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: randomResourceID[domain.OrganizerID](),
	}

	makeMocks := func() (*repositoryMock, *authorizerMock, *scoreEngineManagerMock) {
		mockedRepo := new(repositoryMock)
		mockedAuthorizer := new(authorizerMock)
		mockedScoreEngineManager := new(scoreEngineManagerMock)

		mockedRepo.
			On("GetContest", mock.Anything, nil, fakedContestID).
			Return(domain.Contest{ID: fakedContestID, Ownership: fakedOwnership}, nil)

		return mockedRepo, mockedAuthorizer, mockedScoreEngineManager
	}

	t.Run("ListScoreEngines", func(t *testing.T) {
		t.Run("HappyCase", func(t *testing.T) {
			mockedRepo, mockedAuthorizer, mockedScoreEngineManager := makeMocks()

			mockedAuthorizer.On("HasOwnership", mock.Anything, fakedOwnership).Return(domain.OrganizerRole, nil)

			fakedScoreEngines := []scores.ScoreEngineMeta{
				{
					InstanceID: uuid.New(),
					ContestID:  fakedContestID,
				},
				{
					InstanceID: uuid.New(),
					ContestID:  fakedContestID,
				},
				{
					InstanceID: uuid.New(),
					ContestID:  fakedContestID,
				},
			}

			mockedScoreEngineManager.
				On("ListScoreEnginesByContest", mock.Anything, fakedContestID).
				Return(fakedScoreEngines, nil)

			ucase := usecases.ScoreEngineUseCase{
				Repo:               mockedRepo,
				Authorizer:         mockedAuthorizer,
				ScoreEngineManager: mockedScoreEngineManager,
			}

			instances, err := ucase.ListScoreEnginesByContest(context.Background(), fakedContestID)

			require.NoError(t, err)
			assert.ElementsMatch(t, []domain.ScoreEngineInstanceID{
				fakedScoreEngines[0].InstanceID,
				fakedScoreEngines[1].InstanceID,
				fakedScoreEngines[2].InstanceID,
			}, instances)
		})

		t.Run("BadCredentials", func(t *testing.T) {
			mockedRepo, mockedAuthorizer, _ := makeMocks()

			mockedAuthorizer.On("HasOwnership", mock.Anything, fakedOwnership).Return(domain.NilRole, domain.ErrNoOwnership)

			ucase := usecases.ScoreEngineUseCase{
				Repo:       mockedRepo,
				Authorizer: mockedAuthorizer,
			}

			scoreEngines, err := ucase.ListScoreEnginesByContest(context.Background(), fakedContestID)

			require.ErrorIs(t, err, domain.ErrNoOwnership)
			assert.Nil(t, scoreEngines)
		})
	})

	t.Run("StartScoreEngine", func(t *testing.T) {
		t.Run("HappyCase", func(t *testing.T) {
			mockedRepo, mockedAuthorizer, mockedScoreEngineManager := makeMocks()

			mockedAuthorizer.On("HasOwnership", mock.Anything, fakedOwnership).Return(domain.OrganizerRole, nil)

			fakedInstanceID := domain.ScoreEngineInstanceID(uuid.New())

			mockedScoreEngineManager.
				On("StartScoreEngine", mock.Anything, fakedContestID).
				Return(fakedInstanceID, nil)

			ucase := usecases.ScoreEngineUseCase{
				Repo:               mockedRepo,
				Authorizer:         mockedAuthorizer,
				ScoreEngineManager: mockedScoreEngineManager,
			}

			instanceID, err := ucase.StartScoreEngine(context.Background(), fakedContestID)

			require.NoError(t, err)
			assert.Equal(t, fakedInstanceID, instanceID)
		})

		t.Run("BadCredentials", func(t *testing.T) {
			mockedRepo, mockedAuthorizer, _ := makeMocks()

			mockedAuthorizer.On("HasOwnership", mock.Anything, fakedOwnership).Return(domain.NilRole, domain.ErrNoOwnership)

			ucase := usecases.ScoreEngineUseCase{
				Repo:       mockedRepo,
				Authorizer: mockedAuthorizer,
			}

			instanceID, err := ucase.StartScoreEngine(context.Background(), fakedContestID)

			require.ErrorIs(t, err, domain.ErrNoOwnership)
			assert.Empty(t, instanceID)
		})
	})

	t.Run("StopScoreEngine", func(t *testing.T) {
		t.Run("HappyCase", func(t *testing.T) {
			mockedRepo, mockedAuthorizer, mockedScoreEngineManager := makeMocks()

			mockedAuthorizer.On("HasOwnership", mock.Anything, fakedOwnership).Return(domain.OrganizerRole, nil)

			fakedInstanceID := domain.ScoreEngineInstanceID(uuid.New())

			mockedScoreEngineManager.
				On("GetScoreEngine", mock.Anything, fakedInstanceID).
				Return(scores.ScoreEngineMeta{
					InstanceID: fakedInstanceID,
					ContestID:  fakedContestID,
				}, nil)

			mockedScoreEngineManager.
				On("StopScoreEngine", mock.Anything, fakedInstanceID).
				Return(nil)

			ucase := usecases.ScoreEngineUseCase{
				Repo:               mockedRepo,
				Authorizer:         mockedAuthorizer,
				ScoreEngineManager: mockedScoreEngineManager,
			}

			err := ucase.StopScoreEngine(context.Background(), fakedInstanceID)

			require.NoError(t, err)
		})

		t.Run("BadCredentials", func(t *testing.T) {
			mockedRepo, mockedAuthorizer, mockedScoreEngineManager := makeMocks()

			mockedAuthorizer.On("HasOwnership", mock.Anything, fakedOwnership).Return(domain.NilRole, domain.ErrNoOwnership)

			fakedInstanceID := domain.ScoreEngineInstanceID(uuid.New())

			mockedScoreEngineManager.
				On("GetScoreEngine", mock.Anything, fakedInstanceID).
				Return(scores.ScoreEngineMeta{
					InstanceID: fakedInstanceID,
					ContestID:  fakedContestID,
				}, nil)

			ucase := usecases.ScoreEngineUseCase{
				Repo:               mockedRepo,
				Authorizer:         mockedAuthorizer,
				ScoreEngineManager: mockedScoreEngineManager,
			}

			err := ucase.StopScoreEngine(context.Background(), fakedInstanceID)

			require.ErrorIs(t, err, domain.ErrNoOwnership)
		})
	})
}

type scoreEngineManagerMock struct {
	mock.Mock
}

func (m *scoreEngineManagerMock) GetScoreEngine(ctx context.Context, instanceID domain.ScoreEngineInstanceID) (scores.ScoreEngineMeta, error) {
	args := m.Called(ctx, instanceID)
	return args.Get(0).(scores.ScoreEngineMeta), args.Error(1)
}

func (m *scoreEngineManagerMock) ListScoreEnginesByContest(ctx context.Context, contestID domain.ContestID) ([]scores.ScoreEngineMeta, error) {
	args := m.Called(ctx, contestID)
	return args.Get(0).([]scores.ScoreEngineMeta), args.Error(1)
}

func (m *scoreEngineManagerMock) StopScoreEngine(ctx context.Context, instanceID domain.ScoreEngineInstanceID) error {
	args := m.Called(ctx, instanceID)
	return args.Error(0)
}

func (m *scoreEngineManagerMock) StartScoreEngine(ctx context.Context, contestID domain.ContestID) (domain.ScoreEngineInstanceID, error) {
	args := m.Called(ctx, contestID)
	return args.Get(0).(domain.ScoreEngineInstanceID), args.Error(1)
}
