package usecases_test

import (
	"context"
	"testing"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/testutils"
	"github.com/climblive/platform/backend/internal/usecases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"uuid"
)

func TestScoreEngineUseCase(t *testing.T) {
	fakedCompetitionID := testutils.RandomCompetitionID()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: testutils.RandomResourceID[domain.OrganizerID](),
	}

	now := time.Now()

	makeMocks := func() (*repositoryMock, *authorizerMock, *scoreEngineManagerMock) {
		mockedRepo := new(repositoryMock)
		mockedAuthorizer := new(authorizerMock)
		mockedScoreEngineManager := new(scoreEngineManagerMock)

		return mockedRepo, mockedAuthorizer, mockedScoreEngineManager
	}

	t.Run("ListScoreEngines", func(t *testing.T) {
		t.Run("HappyCase", func(t *testing.T) {
			mockedRepo, mockedAuthorizer, mockedScoreEngineManager := makeMocks()

			secondfakedCompetitionID := testutils.RandomCompetitionID()

			fakedScoreEngines := []domain.ScoreEngineDescriptor{
				{
					InstanceID:    uuid.New(),
					CompetitionID: fakedCompetitionID,
				},
				{
					InstanceID:    uuid.New(),
					CompetitionID: secondfakedCompetitionID,
				},
			}

			mockedScoreEngineManager.
				On("ListScoreEngines", mock.Anything).
				Return(fakedScoreEngines, nil)

			mockedAuthorizer.On("HasOwnership", mock.Anything, domain.OwnershipData{}).Return(domain.AdminRole, nil)

			ucase := usecases.ScoreEngineUseCase{
				Repo:               mockedRepo,
				Authorizer:         mockedAuthorizer,
				ScoreEngineManager: mockedScoreEngineManager,
			}

			engines, err := ucase.ListScoreEngines(context.Background())

			require.NoError(t, err)
			assert.Equal(t, fakedScoreEngines, engines)
		})

		t.Run("BadCredentials", func(t *testing.T) {
			mockedRepo, mockedAuthorizer, _ := makeMocks()

			mockedAuthorizer.On("HasOwnership", mock.Anything, domain.OwnershipData{}).Return(domain.NilRole, domain.ErrNoOwnership)

			ucase := usecases.ScoreEngineUseCase{
				Repo:       mockedRepo,
				Authorizer: mockedAuthorizer,
			}

			engines, err := ucase.ListScoreEngines(context.Background())

			require.ErrorIs(t, err, domain.ErrNoOwnership)
			assert.Empty(t, engines)
		})

		t.Run("NotAdmin", func(t *testing.T) {
			mockedRepo, mockedAuthorizer, _ := makeMocks()

			mockedAuthorizer.On("HasOwnership", mock.Anything, domain.OwnershipData{}).Return(domain.OrganizerRole, nil)

			ucase := usecases.ScoreEngineUseCase{
				Repo:       mockedRepo,
				Authorizer: mockedAuthorizer,
			}

			engines, err := ucase.ListScoreEngines(context.Background())

			require.ErrorIs(t, err, domain.ErrNotAuthorized)
			assert.Empty(t, engines)
		})
	})

	t.Run("ListScoreEnginesByCompetition", func(t *testing.T) {
		t.Run("HappyCase", func(t *testing.T) {
			mockedRepo, mockedAuthorizer, mockedScoreEngineManager := makeMocks()

			mockedRepo.
				On("GetCompetition", mock.Anything, nil, fakedCompetitionID).
				Return(domain.Competition{ID: fakedCompetitionID, Ownership: fakedOwnership}, nil)

			mockedAuthorizer.On("HasOwnership", mock.Anything, fakedOwnership).Return(domain.OrganizerRole, nil)

			fakedScoreEngines := []domain.ScoreEngineDescriptor{
				{
					InstanceID:    uuid.New(),
					CompetitionID: fakedCompetitionID,
				},
				{
					InstanceID:    uuid.New(),
					CompetitionID: fakedCompetitionID,
				},
				{
					InstanceID:    uuid.New(),
					CompetitionID: fakedCompetitionID,
				},
			}

			mockedScoreEngineManager.
				On("ListScoreEnginesByCompetition", mock.Anything, fakedCompetitionID).
				Return(fakedScoreEngines, nil)

			ucase := usecases.ScoreEngineUseCase{
				Repo:               mockedRepo,
				Authorizer:         mockedAuthorizer,
				ScoreEngineManager: mockedScoreEngineManager,
			}

			instances, err := ucase.ListScoreEnginesByCompetition(context.Background(), fakedCompetitionID)

			require.NoError(t, err)
			assert.ElementsMatch(t, fakedScoreEngines, instances)
		})

		t.Run("BadCredentials", func(t *testing.T) {
			mockedRepo, mockedAuthorizer, _ := makeMocks()

			mockedRepo.
				On("GetCompetition", mock.Anything, nil, fakedCompetitionID).
				Return(domain.Competition{ID: fakedCompetitionID, Ownership: fakedOwnership}, nil)

			mockedAuthorizer.On("HasOwnership", mock.Anything, fakedOwnership).Return(domain.NilRole, domain.ErrNoOwnership)

			ucase := usecases.ScoreEngineUseCase{
				Repo:       mockedRepo,
				Authorizer: mockedAuthorizer,
			}

			scoreEngines, err := ucase.ListScoreEnginesByCompetition(context.Background(), fakedCompetitionID)

			require.ErrorIs(t, err, domain.ErrNoOwnership)
			assert.Nil(t, scoreEngines)
		})
	})

	t.Run("StartScoreEngine", func(t *testing.T) {
		t.Run("HappyCase", func(t *testing.T) {
			mockedRepo, mockedAuthorizer, mockedScoreEngineManager := makeMocks()

			endTime := now.Add(time.Hour)

			mockedRepo.
				On("GetCompetition", mock.Anything, nil, fakedCompetitionID).
				Return(domain.Competition{
					ID:        fakedCompetitionID,
					Ownership: fakedOwnership,
					TimeBegin: now,
					TimeEnd:   endTime,
				}, nil)

			mockedAuthorizer.On("HasOwnership", mock.Anything, fakedOwnership).Return(domain.OrganizerRole, nil)

			fakedInstanceID := domain.ScoreEngineInstanceID(uuid.New())

			mockedScoreEngineManager.
				On("StartScoreEngine", mock.Anything, fakedCompetitionID, endTime.Add(time.Hour)).
				Return(fakedInstanceID, nil)

			ucase := usecases.ScoreEngineUseCase{
				Repo:               mockedRepo,
				Authorizer:         mockedAuthorizer,
				ScoreEngineManager: mockedScoreEngineManager,
			}

			instanceID, err := ucase.StartScoreEngine(context.Background(), fakedCompetitionID, endTime.Add(time.Hour))

			require.NoError(t, err)
			assert.Equal(t, fakedInstanceID, instanceID)
		})

		t.Run("CannotStartEngineForCompetitionWithoutStartOrEndTime", func(t *testing.T) {
			mockedRepo, mockedAuthorizer, _ := makeMocks()

			mockedRepo.
				On("GetCompetition", mock.Anything, nil, fakedCompetitionID).
				Return(domain.Competition{
					ID:        fakedCompetitionID,
					Ownership: fakedOwnership,
				}, nil)

			mockedAuthorizer.On("HasOwnership", mock.Anything, fakedOwnership).Return(domain.OrganizerRole, nil)

			ucase := usecases.ScoreEngineUseCase{
				Repo:       mockedRepo,
				Authorizer: mockedAuthorizer,
			}

			_, err := ucase.StartScoreEngine(context.Background(), fakedCompetitionID, time.Now().Add(time.Hour))

			require.ErrorIs(t, err, domain.ErrNotAllowed)
		})

		t.Run("WithinPermittedTimeIntervals", func(t *testing.T) {
			type scenario struct {
				name         string
				timeBegin    time.Time
				timeEnd      time.Time
				terminatedBy time.Time
				expected     error
			}

			scenarios := []scenario{
				{
					name:         "StartScoreEngineMoreThanOneHourBeforeCompetitionStart",
					timeBegin:    now.Add(time.Hour).Add(time.Second),
					timeEnd:      now.Add(2 * time.Hour),
					terminatedBy: now.Add(3 * time.Hour),
					expected:     domain.ErrNotAllowed,
				},
				{
					name:         "StartScoreEngineWithinOneHourOfCompetitionStart",
					timeBegin:    now.Add(time.Hour),
					timeEnd:      now.Add(2 * time.Hour),
					terminatedBy: now.Add(3 * time.Hour),
					expected:     nil,
				},
				{
					name:         "TerminationTimeBeforeCurrentTime",
					timeBegin:    now,
					timeEnd:      now.Add(time.Hour),
					terminatedBy: now,
					expected:     domain.ErrNotAllowed,
				},
				{
					name:         "TerminationTimeWithin12HoursFromNow",
					timeBegin:    now,
					timeEnd:      now.Add(time.Hour),
					terminatedBy: now.Add(12 * time.Hour),
					expected:     nil,
				},
				{
					name:         "TerminationTimePast12HoursFromNow",
					timeBegin:    now,
					timeEnd:      now.Add(time.Hour),
					terminatedBy: now.Add(12 * time.Hour).Add(time.Second),
					expected:     domain.ErrNotAllowed,
				},
			}

			for _, scenario := range scenarios {
				mockedRepo, mockedAuthorizer, mockedScoreEngineManager := makeMocks()

				t.Run(scenario.name, func(t *testing.T) {
					mockedRepo.
						On("GetCompetition", mock.Anything, nil, fakedCompetitionID).
						Return(domain.Competition{
							ID:        fakedCompetitionID,
							Ownership: fakedOwnership,
							TimeBegin: scenario.timeBegin,
							TimeEnd:   scenario.timeEnd,
						}, nil)

					mockedAuthorizer.On("HasOwnership", mock.Anything, fakedOwnership).Return(domain.OrganizerRole, nil)

					fakedInstanceID := domain.ScoreEngineInstanceID(uuid.New())

					mockedScoreEngineManager.
						On("StartScoreEngine", mock.Anything, fakedCompetitionID, scenario.terminatedBy).
						Return(fakedInstanceID, nil)

					ucase := usecases.ScoreEngineUseCase{
						Repo:               mockedRepo,
						Authorizer:         mockedAuthorizer,
						ScoreEngineManager: mockedScoreEngineManager,
					}

					instanceID, err := ucase.StartScoreEngine(context.Background(), fakedCompetitionID, scenario.terminatedBy)

					if scenario.expected == nil {
						require.NoError(t, err)
						assert.Equal(t, fakedInstanceID, instanceID)
					} else {
						assert.ErrorIs(t, err, scenario.expected)
					}
				})
			}
		})

		t.Run("BadCredentials", func(t *testing.T) {
			mockedRepo, mockedAuthorizer, _ := makeMocks()

			mockedRepo.
				On("GetCompetition", mock.Anything, nil, fakedCompetitionID).
				Return(domain.Competition{ID: fakedCompetitionID, Ownership: fakedOwnership}, nil)

			mockedAuthorizer.On("HasOwnership", mock.Anything, fakedOwnership).Return(domain.NilRole, domain.ErrNoOwnership)

			ucase := usecases.ScoreEngineUseCase{
				Repo:       mockedRepo,
				Authorizer: mockedAuthorizer,
			}

			instanceID, err := ucase.StartScoreEngine(context.Background(), fakedCompetitionID, time.Now().Add(time.Hour))

			require.ErrorIs(t, err, domain.ErrNoOwnership)
			assert.Empty(t, instanceID)
		})
	})

	t.Run("StopScoreEngine", func(t *testing.T) {
		t.Run("HappyCase", func(t *testing.T) {
			mockedRepo, mockedAuthorizer, mockedScoreEngineManager := makeMocks()

			mockedRepo.
				On("GetCompetition", mock.Anything, nil, fakedCompetitionID).
				Return(domain.Competition{ID: fakedCompetitionID, Ownership: fakedOwnership}, nil)

			mockedAuthorizer.On("HasOwnership", mock.Anything, fakedOwnership).Return(domain.OrganizerRole, nil)

			fakedInstanceID := domain.ScoreEngineInstanceID(uuid.New())

			mockedScoreEngineManager.
				On("GetScoreEngine", mock.Anything, fakedInstanceID).
				Return(domain.ScoreEngineDescriptor{
					InstanceID:    fakedInstanceID,
					CompetitionID: fakedCompetitionID,
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

			mockedRepo.
				On("GetCompetition", mock.Anything, nil, fakedCompetitionID).
				Return(domain.Competition{ID: fakedCompetitionID, Ownership: fakedOwnership}, nil)

			mockedAuthorizer.On("HasOwnership", mock.Anything, fakedOwnership).Return(domain.NilRole, domain.ErrNoOwnership)

			fakedInstanceID := domain.ScoreEngineInstanceID(uuid.New())

			mockedScoreEngineManager.
				On("GetScoreEngine", mock.Anything, fakedInstanceID).
				Return(domain.ScoreEngineDescriptor{
					InstanceID:    fakedInstanceID,
					CompetitionID: fakedCompetitionID,
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

func (m *scoreEngineManagerMock) GetScoreEngine(ctx context.Context, instanceID domain.ScoreEngineInstanceID) (domain.ScoreEngineDescriptor, error) {
	args := m.Called(ctx, instanceID)
	return args.Get(0).(domain.ScoreEngineDescriptor), args.Error(1)
}

func (m *scoreEngineManagerMock) ListScoreEngines(ctx context.Context) ([]domain.ScoreEngineDescriptor, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.ScoreEngineDescriptor), args.Error(1)
}

func (m *scoreEngineManagerMock) ListScoreEnginesByCompetition(ctx context.Context, competitionID domain.CompetitionID) ([]domain.ScoreEngineDescriptor, error) {
	args := m.Called(ctx, competitionID)
	return args.Get(0).([]domain.ScoreEngineDescriptor), args.Error(1)
}

func (m *scoreEngineManagerMock) StopScoreEngine(ctx context.Context, instanceID domain.ScoreEngineInstanceID) error {
	args := m.Called(ctx, instanceID)
	return args.Error(0)
}

func (m *scoreEngineManagerMock) StartScoreEngine(ctx context.Context, competitionID domain.CompetitionID, terminatedBy time.Time) (domain.ScoreEngineInstanceID, error) {
	args := m.Called(ctx, competitionID, terminatedBy)
	return args.Get(0).(domain.ScoreEngineInstanceID), args.Error(1)
}
