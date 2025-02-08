package usecases_test

import (
	"context"
	"testing"
	"time"

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

			mockedRepo.
				On("GetContest", mock.Anything, nil, fakedContestID).
				Return(domain.Contest{ID: fakedContestID, Ownership: fakedOwnership}, nil)

			mockedAuthorizer.On("HasOwnership", mock.Anything, fakedOwnership).Return(domain.OrganizerRole, nil)

			fakedScoreEngines := []scores.ScoreEngineDescriptor{
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

			mockedRepo.
				On("GetContest", mock.Anything, nil, fakedContestID).
				Return(domain.Contest{ID: fakedContestID, Ownership: fakedOwnership}, nil)

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

			endTime := now.Add(time.Hour)

			mockedRepo.
				On("GetContest", mock.Anything, nil, fakedContestID).
				Return(domain.Contest{
					ID:        fakedContestID,
					Ownership: fakedOwnership,
					TimeBegin: &now,
					TimeEnd:   &endTime,
				}, nil)

			mockedAuthorizer.On("HasOwnership", mock.Anything, fakedOwnership).Return(domain.OrganizerRole, nil)

			fakedInstanceID := domain.ScoreEngineInstanceID(uuid.New())

			mockedScoreEngineManager.
				On("StartScoreEngine", mock.Anything, fakedContestID, endTime.Add(time.Hour)).
				Return(fakedInstanceID, nil)

			ucase := usecases.ScoreEngineUseCase{
				Repo:               mockedRepo,
				Authorizer:         mockedAuthorizer,
				ScoreEngineManager: mockedScoreEngineManager,
			}

			instanceID, err := ucase.StartScoreEngine(context.Background(), fakedContestID, endTime.Add(time.Hour))

			require.NoError(t, err)
			assert.Equal(t, fakedInstanceID, instanceID)
		})

		t.Run("CannotStartEngineForContestWithoutStartOrEndTime", func(t *testing.T) {
			mockedRepo, mockedAuthorizer, _ := makeMocks()

			mockedRepo.
				On("GetContest", mock.Anything, nil, fakedContestID).
				Return(domain.Contest{
					ID:        fakedContestID,
					Ownership: fakedOwnership,
				}, nil)

			mockedAuthorizer.On("HasOwnership", mock.Anything, fakedOwnership).Return(domain.OrganizerRole, nil)

			ucase := usecases.ScoreEngineUseCase{
				Repo:       mockedRepo,
				Authorizer: mockedAuthorizer,
			}

			_, err := ucase.StartScoreEngine(context.Background(), fakedContestID, time.Now().Add(time.Hour))

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
					name:         "TerminationTimeBeforeContestEndTime",
					timeBegin:    now,
					timeEnd:      now.Add(24 * time.Hour),
					terminatedBy: now.Add(24 * time.Hour).Add(-1 * time.Nanosecond),
					expected:     domain.ErrNotAllowed,
				},
				{
					name:         "StartScoreEngineMoreThanOneHourBeforeContestStart",
					timeBegin:    now.Add(time.Hour).Add(time.Second),
					timeEnd:      now.Add(2 * time.Hour),
					terminatedBy: now.Add(3 * time.Hour),
					expected:     domain.ErrNotAllowed,
				},
				{
					name:         "StartScoreEngineWithinOneHourOfContestStart",
					timeBegin:    now.Add(time.Hour),
					timeEnd:      now.Add(2 * time.Hour),
					terminatedBy: now.Add(3 * time.Hour),
					expected:     nil,
				},
				{
					name:         "RunningContest_TerminationTimeWithin12HoursPastEndTime",
					timeBegin:    now.Add(-1 * time.Hour),
					timeEnd:      now.Add(12 * time.Hour),
					terminatedBy: now.Add(24 * time.Hour),
					expected:     nil,
				},
				{
					name:         "RunningContest_TerminationTimeMoreThan12HoursPastEndTime",
					timeBegin:    now.Add(-1 * time.Hour),
					timeEnd:      now.Add(12 * time.Hour),
					terminatedBy: now.Add(24 * time.Hour).Add(time.Nanosecond),
					expected:     domain.ErrNotAllowed,
				},
				{
					name:         "EndedContest_TerminationTimeWithinOneHourFromNow",
					timeBegin:    now.Add(-24 * time.Hour),
					timeEnd:      now.Add(-12 * time.Hour),
					terminatedBy: now.Add(time.Hour),
					expected:     nil,
				},
				{
					name:         "EndedContest_TerminationTimePastOneHourFromNow",
					timeBegin:    now.Add(-24 * time.Hour),
					timeEnd:      now.Add(-12 * time.Hour),
					terminatedBy: now.Add(time.Hour).Add(time.Second),
					expected:     domain.ErrNotAllowed,
				},
			}

			for _, scenario := range scenarios {
				mockedRepo, mockedAuthorizer, mockedScoreEngineManager := makeMocks()

				t.Run(scenario.name, func(t *testing.T) {
					mockedRepo.
						On("GetContest", mock.Anything, nil, fakedContestID).
						Return(domain.Contest{
							ID:        fakedContestID,
							Ownership: fakedOwnership,
							TimeBegin: &scenario.timeBegin,
							TimeEnd:   &scenario.timeEnd,
						}, nil)

					mockedAuthorizer.On("HasOwnership", mock.Anything, fakedOwnership).Return(domain.OrganizerRole, nil)

					fakedInstanceID := domain.ScoreEngineInstanceID(uuid.New())

					mockedScoreEngineManager.
						On("StartScoreEngine", mock.Anything, fakedContestID, scenario.terminatedBy).
						Return(fakedInstanceID, nil)

					ucase := usecases.ScoreEngineUseCase{
						Repo:               mockedRepo,
						Authorizer:         mockedAuthorizer,
						ScoreEngineManager: mockedScoreEngineManager,
					}

					instanceID, err := ucase.StartScoreEngine(context.Background(), fakedContestID, scenario.terminatedBy)

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
				On("GetContest", mock.Anything, nil, fakedContestID).
				Return(domain.Contest{ID: fakedContestID, Ownership: fakedOwnership}, nil)

			mockedAuthorizer.On("HasOwnership", mock.Anything, fakedOwnership).Return(domain.NilRole, domain.ErrNoOwnership)

			ucase := usecases.ScoreEngineUseCase{
				Repo:       mockedRepo,
				Authorizer: mockedAuthorizer,
			}

			instanceID, err := ucase.StartScoreEngine(context.Background(), fakedContestID, time.Now().Add(time.Hour))

			require.ErrorIs(t, err, domain.ErrNoOwnership)
			assert.Empty(t, instanceID)
		})
	})

	t.Run("StopScoreEngine", func(t *testing.T) {
		t.Run("HappyCase", func(t *testing.T) {
			mockedRepo, mockedAuthorizer, mockedScoreEngineManager := makeMocks()

			mockedRepo.
				On("GetContest", mock.Anything, nil, fakedContestID).
				Return(domain.Contest{ID: fakedContestID, Ownership: fakedOwnership}, nil)

			mockedAuthorizer.On("HasOwnership", mock.Anything, fakedOwnership).Return(domain.OrganizerRole, nil)

			fakedInstanceID := domain.ScoreEngineInstanceID(uuid.New())

			mockedScoreEngineManager.
				On("GetScoreEngine", mock.Anything, fakedInstanceID).
				Return(scores.ScoreEngineDescriptor{
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

			mockedRepo.
				On("GetContest", mock.Anything, nil, fakedContestID).
				Return(domain.Contest{ID: fakedContestID, Ownership: fakedOwnership}, nil)

			mockedAuthorizer.On("HasOwnership", mock.Anything, fakedOwnership).Return(domain.NilRole, domain.ErrNoOwnership)

			fakedInstanceID := domain.ScoreEngineInstanceID(uuid.New())

			mockedScoreEngineManager.
				On("GetScoreEngine", mock.Anything, fakedInstanceID).
				Return(scores.ScoreEngineDescriptor{
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

func (m *scoreEngineManagerMock) GetScoreEngine(ctx context.Context, instanceID domain.ScoreEngineInstanceID) (scores.ScoreEngineDescriptor, error) {
	args := m.Called(ctx, instanceID)
	return args.Get(0).(scores.ScoreEngineDescriptor), args.Error(1)
}

func (m *scoreEngineManagerMock) ListScoreEnginesByContest(ctx context.Context, contestID domain.ContestID) ([]scores.ScoreEngineDescriptor, error) {
	args := m.Called(ctx, contestID)
	return args.Get(0).([]scores.ScoreEngineDescriptor), args.Error(1)
}

func (m *scoreEngineManagerMock) StopScoreEngine(ctx context.Context, instanceID domain.ScoreEngineInstanceID) error {
	args := m.Called(ctx, instanceID)
	return args.Error(0)
}

func (m *scoreEngineManagerMock) StartScoreEngine(ctx context.Context, contestID domain.ContestID, terminatedBy time.Time) (domain.ScoreEngineInstanceID, error) {
	args := m.Called(ctx, contestID, terminatedBy)
	return args.Get(0).(domain.ScoreEngineInstanceID), args.Error(1)
}
