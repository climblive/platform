package usecases_test

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"testing/synctest"
	"time"

	"uuid"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/testutils"
	"github.com/climblive/platform/backend/internal/usecases"
	"github.com/climblive/platform/backend/internal/usecases/validators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetCompetition(t *testing.T) {
	fakedCompetitionID := testutils.RandomResourceID[domain.CompetitionID]()

	fakedCompetition := domain.Competition{
		ID: fakedCompetitionID,
	}

	mockedRepo := new(repositoryMock)

	mockedRepo.
		On("GetCompetition", mock.Anything, mock.Anything, fakedCompetitionID).
		Return(fakedCompetition, nil)

	ucase := usecases.CompetitionUseCase{
		Repo: mockedRepo,
	}

	competition, err := ucase.GetCompetition(context.Background(), fakedCompetitionID)

	require.NoError(t, err)
	assert.Equal(t, fakedCompetitionID, competition.ID)

	mockedRepo.AssertExpectations(t)
}

func TestGetScoreboard(t *testing.T) {
	fakedCompetitionID := testutils.RandomResourceID[domain.CompetitionID]()
	fakedCategoryID := testutils.RandomResourceID[domain.CategoryID]()

	mockedRepo := new(repositoryMock)
	mockedResultKeeper := new(resultKeeperMock)

	currentTime := time.Now()

	var competitors []domain.Competitor

	for i := 1; i <= 10; i++ {
		competitorID := domain.CompetitorID(i)

		fakedCompetitor := domain.Competitor{
			ID:                  competitorID,
			CategoryID:          fakedCategoryID,
			Name:                fmt.Sprintf("Climber %d", i),
			WithdrawnFromFinals: true,
			Disqualified:        true,
			ScrubbedAt:          currentTime,
			Result: &domain.Result{
				CompetitorID: competitorID,
				Score:        fmt.Sprintf("%dp", i*10),
				Placement:    i,
				Finalist:     true,
				RankOrder:    i - 1,
				Timestamp:    currentTime,
			},
		}

		competitors = append(competitors, fakedCompetitor)
	}

	mockedRepo.
		On("GetCompetitorsByCompetition", mock.Anything, mock.Anything, fakedCompetitionID).
		Return(competitors, nil)

	future := currentTime.Add(time.Minute)

	mockedResultKeeper.On("GetResult", domain.CompetitorID(1)).Return(domain.Result{
		Timestamp:    future,
		CompetitorID: domain.CompetitorID(1),
		Score:        "1234p",
		Placement:    42,
		Finalist:     false,
		RankOrder:    1337,
	}, nil)

	mockedResultKeeper.On("GetResult", mock.Anything).Return(domain.Result{}, errMock)

	ucase := usecases.CompetitionUseCase{
		Repo:         mockedRepo,
		ResultKeeper: mockedResultKeeper,
	}

	scoreboard, err := ucase.GetScoreboard(context.Background(), fakedCompetitionID)

	require.NoError(t, err)

	assert.Len(t, scoreboard, 10)

	assert.Equal(t, domain.CompetitorID(1), scoreboard[0].CompetitorID)
	assert.Equal(t, fakedCategoryID, scoreboard[0].CategoryID)
	assert.Equal(t, "Climber 1", scoreboard[0].Name)
	assert.Equal(t, true, scoreboard[0].WithdrawnFromFinals)
	assert.Equal(t, true, scoreboard[0].Disqualified)
	assert.Equal(t, currentTime, scoreboard[0].ScrubbedAt)
	assert.NotNil(t, scoreboard[0].Result)
	assert.Equal(t, future, scoreboard[0].Result.Timestamp)
	assert.Equal(t, "1234p", scoreboard[0].Result.Score)
	assert.Equal(t, 42, scoreboard[0].Result.Placement)
	assert.False(t, scoreboard[0].Result.Finalist)
	assert.Equal(t, 1337, scoreboard[0].Result.RankOrder)

	for i := 2; i <= 10; i++ {
		entry := scoreboard[i-1]

		assert.Equal(t, domain.CompetitorID(i), entry.CompetitorID)
		assert.Equal(t, fakedCategoryID, entry.CategoryID)
		assert.Equal(t, fmt.Sprintf("Climber %d", i), entry.Name)
		assert.Equal(t, true, entry.WithdrawnFromFinals)
		assert.Equal(t, true, entry.Disqualified)
		assert.Equal(t, currentTime, entry.ScrubbedAt)
		assert.NotNil(t, entry.Result)
		assert.Equal(t, currentTime, entry.Result.Timestamp)
		assert.Equal(t, fmt.Sprintf("%dp", i*10), entry.Result.Score)
		assert.Equal(t, i, entry.Result.Placement)
		assert.Equal(t, i-1, entry.Result.RankOrder)
		assert.True(t, entry.Result.Finalist)
	}

	mockedRepo.AssertExpectations(t)
	mockedResultKeeper.AssertExpectations(t)
}

func TestGetScoreboard_Empty(t *testing.T) {
	fakedCompetitionID := testutils.RandomResourceID[domain.CompetitionID]()

	mockedRepo := new(repositoryMock)

	mockedRepo.
		On("GetCompetitorsByCompetition", mock.Anything, mock.Anything, fakedCompetitionID).
		Return([]domain.Competitor{}, nil)

	ucase := usecases.CompetitionUseCase{
		Repo: mockedRepo,
	}

	scoreboard, err := ucase.GetScoreboard(context.Background(), fakedCompetitionID)

	require.NoError(t, err)
	assert.NotNil(t, scoreboard)
	assert.Len(t, scoreboard, 0)

	mockedRepo.AssertExpectations(t)
}

func TestGetCompetitionsByOrganizer(t *testing.T) {
	fakedOrganizerID := testutils.RandomResourceID[domain.OrganizerID]()
	fakedCompetitionID := testutils.RandomResourceID[domain.CompetitionID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: fakedOrganizerID,
	}

	makeMocks := func() (*repositoryMock, *authorizerMock) {
		mockedRepo := new(repositoryMock)

		mockedRepo.
			On("GetOrganizer", mock.Anything, nil, fakedOrganizerID).
			Return(domain.Organizer{
				ID:        fakedOrganizerID,
				Ownership: fakedOwnership,
			}, nil)

		mockedAuthorizer := new(authorizerMock)

		return mockedRepo, mockedAuthorizer
	}

	t.Run("HappyCase", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("GetCompetitionsByOrganizer", mock.Anything, nil, fakedOrganizerID).
			Return([]domain.Competition{
				{
					ID:        fakedCompetitionID,
					Ownership: fakedOwnership,
				},
			}, nil)

		ucase := usecases.CompetitionUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		competitions, err := ucase.GetCompetitionsByOrganizer(context.Background(), fakedOrganizerID)

		require.NoError(t, err)
		require.Len(t, competitions, 1)
		assert.Equal(t, fakedCompetitionID, competitions[0].ID)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		ucase := usecases.CompetitionUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		_, err := ucase.GetCompetitionsByOrganizer(context.Background(), fakedOrganizerID)

		require.ErrorIs(t, err, domain.ErrNoOwnership)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})
}

func TestGetAllCompetitions(t *testing.T) {
	fakedOrganizerID := testutils.RandomResourceID[domain.OrganizerID]()
	fakedCompetitionID := testutils.RandomResourceID[domain.CompetitionID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: fakedOrganizerID,
	}

	makeMocks := func() (*repositoryMock, *authorizerMock) {
		mockedRepo := new(repositoryMock)

		mockedAuthorizer := new(authorizerMock)

		return mockedRepo, mockedAuthorizer
	}

	t.Run("HappyCase", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mock.AnythingOfType("domain.OwnershipData")).
			Return(domain.AdminRole, nil)

		mockedRepo.
			On("GetAllCompetitions", mock.Anything, nil).
			Return([]domain.Competition{
				{
					ID:        fakedCompetitionID,
					Ownership: fakedOwnership,
				},
			}, nil)

		ucase := usecases.CompetitionUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		competitions, err := ucase.GetAllCompetitions(context.Background())

		require.NoError(t, err)
		require.Len(t, competitions, 1)
		assert.Equal(t, fakedCompetitionID, competitions[0].ID)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mock.AnythingOfType("domain.OwnershipData")).
			Return(domain.NilRole, domain.ErrNoOwnership)

		ucase := usecases.CompetitionUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		_, err := ucase.GetAllCompetitions(context.Background())

		require.ErrorIs(t, err, domain.ErrNoOwnership)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("NotAdmin", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, mock.AnythingOfType("domain.OwnershipData")).
			Return(domain.OrganizerRole, nil)

		ucase := usecases.CompetitionUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		_, err := ucase.GetAllCompetitions(context.Background())

		require.ErrorIs(t, err, domain.ErrNotAuthorized)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})
}
func TestCreateCompetition(t *testing.T) {
	fakedOrganizerID := testutils.RandomResourceID[domain.OrganizerID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: fakedOrganizerID,
	}
	fakedCompetitionID := testutils.RandomResourceID[domain.CompetitionID]()

	makeMocks := func() (*repositoryMock, *authorizerMock) {
		mockedRepo := new(repositoryMock)

		mockedRepo.
			On("GetOrganizer", mock.Anything, nil, fakedOrganizerID).
			Return(domain.Organizer{
				ID:        fakedOrganizerID,
				Ownership: fakedOwnership,
			}, nil)

		mockedAuthorizer := new(authorizerMock)

		return mockedRepo, mockedAuthorizer
	}

	t.Run("HappyCase", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			mockedRepo, mockedAuthorizer := makeMocks()

			mockedAuthorizer.
				On("HasOwnership", mock.Anything, fakedOwnership).
				Return(domain.OrganizerRole, nil)

			mockedRepo.
				On("GetCompetitionsByOrganizer", mock.Anything, nil, fakedOrganizerID).
				Return([]domain.Competition{}, nil)

			mockedRepo.
				On("StoreCompetition", mock.Anything, nil,
					domain.Competition{
						Ownership:          fakedOwnership,
						Location:           "The garage",
						Country:            "SE",
						SeriesID:           0,
						Name:               "Swedish Championships",
						Description:        "Who is the best climber in Sweden?",
						QualifyingProblems: 10,
						Finalists:          7,
						UsePoints:          true,
						PooledPoints:       true,
						MaxAttempts:        5,
						PointDeduction:     2,
						Info:               "No rules!",
						GracePeriod:        time.Hour,
						Created:            time.Now(),
						NameRetentionTime:  14 * 24 * time.Hour,
					},
				).
				Return(domain.Competition{
					ID:                 fakedCompetitionID,
					Ownership:          fakedOwnership,
					Location:           "The garage",
					Country:            "SE",
					SeriesID:           0,
					Name:               "Swedish Championships",
					Description:        "Who is the best climber in Sweden?",
					QualifyingProblems: 10,
					Finalists:          7,
					UsePoints:          true,
					PooledPoints:       true,
					MaxAttempts:        5,
					PointDeduction:     2,
					Info:               "No rules!",
					GracePeriod:        time.Hour,
					Created:            time.Now(),
					NameRetentionTime:  14 * 24 * time.Hour,
				}, nil)

			ucase := usecases.CompetitionUseCase{
				Repo:       mockedRepo,
				Authorizer: mockedAuthorizer,
			}

			competition, err := ucase.CreateCompetition(context.Background(), fakedOrganizerID, domain.CompetitionTemplate{
				Location:           "The garage",
				Country:            "SE",
				Name:               "Swedish Championships",
				Description:        "Who is the best climber in Sweden?",
				QualifyingProblems: 10,
				Finalists:          7,
				UsePoints:          true,
				PooledPoints:       true,
				MaxAttempts:        5,
				PointDeduction:     2,
				Info:               "No rules!",
				GracePeriod:        time.Hour,
				NameRetentionTime:  14 * 24 * time.Hour,
			})

			require.NoError(t, err)
			assert.Equal(t, fakedCompetitionID, competition.ID)
			assert.Equal(t, fakedOwnership, competition.Ownership)
			assert.Zero(t, competition.ArchivedAt)
			assert.Equal(t, "The garage", competition.Location)
			assert.Equal(t, "SE", competition.Country)
			assert.Equal(t, "Swedish Championships", competition.Name)
			assert.Equal(t, "Who is the best climber in Sweden?", competition.Description)
			assert.Equal(t, 10, competition.QualifyingProblems)
			assert.Equal(t, 7, competition.Finalists)
			assert.True(t, competition.UsePoints)
			assert.True(t, competition.PooledPoints)
			assert.Equal(t, 5, competition.MaxAttempts)
			assert.Equal(t, domain.Tenths(2), competition.PointDeduction)
			assert.Equal(t, "No rules!", competition.Info)
			assert.Equal(t, time.Hour, competition.GracePeriod)
			assert.Empty(t, competition.TimeBegin)
			assert.Empty(t, competition.TimeEnd)
			assert.Equal(t, time.Now(), competition.Created)
			assert.Equal(t, 14*24*time.Hour, competition.NameRetentionTime)

			mockedRepo.AssertExpectations(t)
			mockedAuthorizer.AssertExpectations(t)
		})
	})

	t.Run("ValidatorIsInvoked", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("GetCompetitionsByOrganizer", mock.Anything, nil, fakedOrganizerID).
			Return([]domain.Competition{}, nil)

		ucase := usecases.CompetitionUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		_, err := ucase.CreateCompetition(context.Background(), fakedOrganizerID, domain.CompetitionTemplate{})

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.True(t, validators.CompetitionValidator{}.IsValidationError(err))

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("InfoIsSanitized", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			mockedRepo, mockedAuthorizer := makeMocks()

			mockedAuthorizer.
				On("HasOwnership", mock.Anything, fakedOwnership).
				Return(domain.OrganizerRole, nil)

			mockedRepo.
				On("GetCompetitionsByOrganizer", mock.Anything, nil, fakedOrganizerID).
				Return([]domain.Competition{}, nil)

			mockedRepo.
				On("StoreCompetition", mock.Anything, nil,
					domain.Competition{
						Ownership:          fakedOwnership,
						Location:           "The garage",
						Country:            "SE",
						Name:               "Swedish Championships",
						Description:        "Who is the best climber in Sweden?",
						QualifyingProblems: 10,
						UsePoints:          true,
						Finalists:          7,
						Info:               "XSS",
						GracePeriod:        time.Hour,
						NameRetentionTime:  14 * 24 * time.Hour,
						Created:            time.Now(),
					},
				).
				Return(domain.Competition{
					Ownership:          fakedOwnership,
					Location:           "The garage",
					Country:            "SE",
					Name:               "Swedish Championships",
					Description:        "Who is the best climber in Sweden?",
					QualifyingProblems: 10,
					UsePoints:          true,
					Finalists:          7,
					Info:               "XSS",
					GracePeriod:        time.Hour,
					NameRetentionTime:  14 * 24 * time.Hour,
					Created:            time.Now(),
				}, nil)

			ucase := usecases.CompetitionUseCase{
				Repo:       mockedRepo,
				Authorizer: mockedAuthorizer,
			}

			competition, err := ucase.CreateCompetition(context.Background(), fakedOrganizerID, domain.CompetitionTemplate{
				Location:           "The garage",
				Country:            "SE",
				Name:               "Swedish Championships",
				Description:        "Who is the best climber in Sweden?",
				QualifyingProblems: 10,
				UsePoints:          true,
				Finalists:          7,
				Info:               `<a href="javascript:alert('XSS1')" onmouseover="alert('XSS2')">XSS<a>`,
				GracePeriod:        time.Hour,
				NameRetentionTime:  14 * 24 * time.Hour,
			})

			require.NoError(t, err)
			assert.Equal(t, "XSS", competition.Info)

			mockedRepo.AssertExpectations(t)
			mockedAuthorizer.AssertExpectations(t)
		})
	})

	t.Run("LimitExceeded", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			mockedRepo, mockedAuthorizer := makeMocks()

			mockedAuthorizer.
				On("HasOwnership", mock.Anything, fakedOwnership).
				Return(domain.OrganizerRole, nil)

			recentCompetitions := make([]domain.Competition, 10)
			for i := range recentCompetitions {
				recentCompetitions[i] = domain.Competition{
					Created: time.Now().Add(-7 * 24 * time.Hour).Add(time.Nanosecond),
				}
			}
			mockedRepo.
				On("GetCompetitionsByOrganizer", mock.Anything, nil, fakedOrganizerID).
				Return(recentCompetitions, nil)

			ucase := usecases.CompetitionUseCase{
				Repo:       mockedRepo,
				Authorizer: mockedAuthorizer,
			}

			_, err := ucase.CreateCompetition(context.Background(), fakedOrganizerID, domain.CompetitionTemplate{})

			require.ErrorIs(t, err, domain.ErrLimitExceeded)

			mockedRepo.AssertExpectations(t)
			mockedAuthorizer.AssertExpectations(t)
		})
	})

	t.Run("ArchivedCompetitionsDoNotCountTowardsLimit", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			mockedRepo, mockedAuthorizer := makeMocks()

			mockedAuthorizer.
				On("HasOwnership", mock.Anything, fakedOwnership).
				Return(domain.OrganizerRole, nil)

			competitions := make([]domain.Competition, 10)
			for i := range competitions {
				competitions[i] = domain.Competition{
					ArchivedAt: time.Now(),
					Created:    time.Now(),
				}
			}
			mockedRepo.
				On("GetCompetitionsByOrganizer", mock.Anything, nil, fakedOrganizerID).
				Return(competitions, nil)

			storedCompetition := domain.Competition{
				Ownership: domain.OwnershipData{
					OrganizerID: fakedOrganizerID,
				},
				Location:           "The garage",
				Country:            "SE",
				Name:               "Swedish Championships",
				QualifyingProblems: 10,
				UsePoints:          true,
				Finalists:          7,
				GracePeriod:        time.Hour,
				NameRetentionTime:  14 * 24 * time.Hour,
				Created:            time.Now(),
			}

			mockedRepo.
				On("StoreCompetition", mock.Anything, nil, storedCompetition).
				Return(storedCompetition, nil)

			ucase := usecases.CompetitionUseCase{
				Repo:       mockedRepo,
				Authorizer: mockedAuthorizer,
			}

			_, err := ucase.CreateCompetition(context.Background(), fakedOrganizerID, domain.CompetitionTemplate{
				Location:           "The garage",
				Country:            "SE",
				Name:               "Swedish Championships",
				QualifyingProblems: 10,
				UsePoints:          true,
				Finalists:          7,
				GracePeriod:        time.Hour,
				NameRetentionTime:  14 * 24 * time.Hour,
			})

			require.NoError(t, err)

			mockedRepo.AssertExpectations(t)
			mockedAuthorizer.AssertExpectations(t)
		})
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		ucase := usecases.CompetitionUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		_, err := ucase.CreateCompetition(context.Background(), fakedOrganizerID, domain.CompetitionTemplate{})

		require.ErrorIs(t, err, domain.ErrNoOwnership)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})
}

func TestDuplicateCompetition(t *testing.T) {
	fakedCompetitionID := testutils.RandomResourceID[domain.CompetitionID]()
	fakedDuplicatedCompetitionID := testutils.RandomResourceID[domain.CompetitionID]()
	fakedCategoryID := testutils.RandomResourceID[domain.CategoryID]()
	fakedProblemID := testutils.RandomResourceID[domain.ProblemID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: testutils.RandomResourceID[domain.OrganizerID](),
	}

	timeBegin := time.Now()
	timeEnd := timeBegin.Add(time.Hour)

	fakedCompetition := domain.Competition{
		ID:                 fakedCompetitionID,
		Ownership:          fakedOwnership,
		Location:           "The garage",
		SeriesID:           testutils.RandomResourceID[domain.SeriesID](),
		Name:               "Original Competition",
		Description:        "Who is the best climber in Sweden?",
		QualifyingProblems: 10,
		Finalists:          7,
		UsePoints:          true,
		PooledPoints:       true,
		MaxAttempts:        5,
		PointDeduction:     2,
		Info:               "No rules!",
		GracePeriod:        time.Hour,
		TimeBegin:          timeBegin,
		TimeEnd:            timeEnd,
	}

	fakedCategory := domain.Category{
		ID:            fakedCategoryID,
		Ownership:     fakedOwnership,
		CompetitionID: fakedCompetitionID,
		Name:          "Males",
		Description:   "Male climbers",
		TimeBegin:     timeBegin,
		TimeEnd:       timeEnd,
	}

	fakedProblem := domain.Problem{
		ID:                 fakedProblemID,
		Ownership:          fakedOwnership,
		CompetitionID:      fakedCompetitionID,
		Number:             42,
		HoldColorPrimary:   "#FF0000",
		HoldColorSecondary: "#00FF00",
		Description:        "Test Problem",
		Zone1Enabled:       true,
		Zone2Enabled:       true,
		ProblemValue: domain.ProblemValue{
			PointsTop:   100,
			PointsZone1: 50,
			PointsZone2: 75,
			FlashBonus:  20,
		},
	}

	makeMocks := func() (*repositoryMock, *authorizerMock) {
		mockedRepo := new(repositoryMock)
		mockedAuthorizer := new(authorizerMock)

		return mockedRepo, mockedAuthorizer
	}

	t.Run("HappyCase", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

		fakedDuplicatedCompetition := fakedCompetition
		fakedDuplicatedCompetition.ID = fakedDuplicatedCompetitionID
		fakedDuplicatedCompetition.Name = "Original Competition (Copy)"

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("GetCompetition", mock.Anything, nil, fakedCompetitionID).
			Return(fakedCompetition, nil)

		mockedRepo.
			On("GetCategoriesByCompetition", mock.Anything, nil, fakedCompetitionID).
			Return([]domain.Category{fakedCategory}, nil)

		mockedRepo.
			On("GetProblemsByCompetition", mock.Anything, nil, fakedCompetitionID).
			Return([]domain.Problem{fakedProblem}, nil)

		mockedTx := new(transactionMock)

		mockedRepo.
			On("Begin").
			Return(mockedTx, nil)

		mockedRepo.
			On("StoreCompetition", mock.Anything, mockedTx, mock.MatchedBy(func(competition domain.Competition) bool {
				expected := fakedDuplicatedCompetition
				expected.ID = 0

				return competition == expected
			})).
			Return(fakedDuplicatedCompetition, nil)

		mockedRepo.
			On("StoreCategory", mock.Anything, mockedTx, mock.MatchedBy(func(category domain.Category) bool {
				expected := fakedCategory
				expected.ID = 0
				expected.CompetitionID = fakedDuplicatedCompetitionID

				return category == expected
			})).
			Return(domain.Category{}, nil)

		mockedRepo.
			On("StoreProblem", mock.Anything, mockedTx, mock.MatchedBy(func(problem domain.Problem) bool {
				expected := fakedProblem
				expected.ID = 0
				expected.CompetitionID = fakedDuplicatedCompetitionID

				return problem == expected
			})).
			Return(domain.Problem{}, nil)

		mockedTx.
			On("Commit").
			Return(nil)

		ucase := usecases.CompetitionUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		duplicatedCompetition, err := ucase.DuplicateCompetition(context.Background(), fakedCompetitionID)

		require.NoError(t, err)
		assert.Equal(t, fakedDuplicatedCompetitionID, duplicatedCompetition.ID)
		assert.Equal(t, "Original Competition (Copy)", duplicatedCompetition.Name)
		assert.Equal(t, fakedOwnership, duplicatedCompetition.Ownership)
		assert.Equal(t, "The garage", duplicatedCompetition.Location)
		assert.Equal(t, fakedCompetition.SeriesID, duplicatedCompetition.SeriesID)
		assert.Equal(t, "Who is the best climber in Sweden?", duplicatedCompetition.Description)
		assert.Equal(t, 10, duplicatedCompetition.QualifyingProblems)
		assert.Equal(t, 7, duplicatedCompetition.Finalists)
		assert.True(t, duplicatedCompetition.UsePoints)
		assert.True(t, duplicatedCompetition.PooledPoints)
		assert.Equal(t, 5, duplicatedCompetition.MaxAttempts)
		assert.Equal(t, domain.Tenths(2), duplicatedCompetition.PointDeduction)
		assert.Equal(t, "No rules!", duplicatedCompetition.Info)
		assert.Equal(t, time.Hour, duplicatedCompetition.GracePeriod)
		assert.Equal(t, timeBegin, duplicatedCompetition.TimeBegin)
		assert.NotZero(t, timeEnd, duplicatedCompetition.TimeEnd)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		mockedRepo.
			On("GetCompetition", mock.Anything, nil, fakedCompetitionID).
			Return(fakedCompetition, nil)

		ucase := usecases.CompetitionUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		_, err := ucase.DuplicateCompetition(context.Background(), fakedCompetitionID)

		require.ErrorIs(t, err, domain.ErrNoOwnership)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("CompetitionArchived", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("GetCompetition", mock.Anything, nil, fakedCompetitionID).
			Return(domain.Competition{
				ID:         fakedCompetitionID,
				Ownership:  fakedOwnership,
				ArchivedAt: time.Now(),
			}, nil)

		ucase := usecases.CompetitionUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		_, err := ucase.DuplicateCompetition(context.Background(), fakedCompetitionID)

		require.ErrorIs(t, err, domain.ErrArchived)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})
}

func TestTransferCompetition(t *testing.T) {
	fakedCompetitionID := testutils.RandomResourceID[domain.CompetitionID]()
	fakedOldOrganizerID := testutils.RandomResourceID[domain.OrganizerID]()
	fakedNewOrganizerID := testutils.RandomResourceID[domain.OrganizerID]()
	fakedCategoryID := testutils.RandomResourceID[domain.CategoryID]()
	fakedProblemID := testutils.RandomResourceID[domain.ProblemID]()
	fakedCompetitorID := testutils.RandomResourceID[domain.CompetitorID]()
	fakedRaffleID := testutils.RandomResourceID[domain.RaffleID]()
	fakedRaffleWinnerID := testutils.RandomResourceID[domain.RaffleWinnerID]()
	fakedTickID := testutils.RandomResourceID[domain.TickID]()
	fakedSeriesID := testutils.RandomResourceID[domain.SeriesID]()

	now := time.Now()

	fakedOldOwnership := domain.OwnershipData{OrganizerID: fakedOldOrganizerID}
	fakedNewOwnership := domain.OwnershipData{OrganizerID: fakedNewOrganizerID}

	timeBegin := time.Now()
	timeEnd := timeBegin.Add(2 * time.Hour)

	fakedCompetition := domain.Competition{
		ID:                 fakedCompetitionID,
		Ownership:          fakedOldOwnership,
		Location:           "Stockholm Climbing Center",
		SeriesID:           fakedSeriesID,
		Name:               "Swedish Open 2025",
		Description:        "National bouldering championship",
		QualifyingProblems: 8,
		Finalists:          6,
		UsePoints:          true,
		PooledPoints:       true,
		MaxAttempts:        5,
		PointDeduction:     2,
		Info:               "Standard IFSC rules apply",
		GracePeriod:        30 * time.Minute,
		TimeBegin:          timeBegin,
		TimeEnd:            timeEnd,
		Created:            now.Add(time.Duration(rand.Int())),
	}

	fakedCategory := domain.Category{
		ID:            fakedCategoryID,
		Ownership:     fakedOldOwnership,
		CompetitionID: fakedCompetitionID,
		Name:          "Males",
		Description:   "Male competitors",
		TimeBegin:     timeBegin,
		TimeEnd:       timeEnd,
	}

	fakedProblem := domain.Problem{
		ID:                 fakedProblemID,
		Ownership:          fakedOldOwnership,
		CompetitionID:      fakedCompetitionID,
		Number:             15,
		HoldColorPrimary:   "#0000FF",
		HoldColorSecondary: "#FFFFFF",
		Description:        "Technical slab problem",
		Zone1Enabled:       true,
		Zone2Enabled:       true,
		ProblemValue: domain.ProblemValue{
			PointsTop:   100,
			PointsZone1: 25,
			PointsZone2: 50,
			FlashBonus:  10,
		},
	}

	fakedResult := domain.Result{
		Timestamp:    now.Add(time.Duration(rand.Int())),
		CompetitorID: fakedCompetitorID,
		Score:        "10p",
		Placement:    3,
		Finalist:     true,
		RankOrder:    2,
	}

	fakedCompetitor := domain.Competitor{
		ID:                  fakedCompetitorID,
		Ownership:           domain.OwnershipData{OrganizerID: fakedOldOrganizerID, CompetitorID: &fakedCompetitorID},
		CompetitionID:       fakedCompetitionID,
		CategoryID:          fakedCategoryID,
		RegistrationCode:    "ABCD1234",
		Name:                "John Doe",
		Entered:             now.Add(time.Duration(rand.Int())),
		WithdrawnFromFinals: true,
		Disqualified:        true,
		Result:              &fakedResult,
	}

	fakedRaffle := domain.Raffle{
		ID:            fakedRaffleID,
		Ownership:     fakedOldOwnership,
		CompetitionID: fakedCompetitionID,
	}

	fakedWinner := domain.RaffleWinner{
		ID:             fakedRaffleWinnerID,
		Ownership:      fakedOldOwnership,
		RaffleID:       fakedRaffleID,
		CompetitorID:   fakedCompetitorID,
		CompetitorName: "John Doe",
		Timestamp:      now.Add(time.Duration(rand.Int())),
	}

	fakedTick := testutils.MustTick(domain.Tick{
		Revision:      1,
		ID:            fakedTickID,
		Ownership:     domain.OwnershipData{OrganizerID: fakedOldOrganizerID, CompetitorID: &fakedCompetitorID},
		Timestamp:     now.Add(time.Duration(rand.Int())),
		CompetitionID: fakedCompetitionID,
		ProblemID:     fakedProblemID,
		Zone1:         true,
		AttemptsZone1: 2,
		Zone2:         true,
		AttemptsZone2: 3,
		Top:           true,
		AttemptsTop:   5,
	})

	makeMocks := func() (*repositoryMock, *authorizerMock, *transactionMock) {
		mockedRepo := new(repositoryMock)
		mockedAuthorizer := new(authorizerMock)
		mockedTx := new(transactionMock)

		return mockedRepo, mockedAuthorizer, mockedTx
	}

	t.Run("HappyCase", func(t *testing.T) {
		mockedRepo, mockedAuthorizer, mockedTx := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOldOwnership).
			Return(domain.OrganizerRole, nil).
			On("HasOwnership", mock.Anything, fakedNewOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("GetCompetition", mock.Anything, nil, fakedCompetitionID).
			Return(fakedCompetition, nil)
		mockedRepo.
			On("GetOrganizer", mock.Anything, nil, fakedNewOrganizerID).
			Return(domain.Organizer{
				ID:        fakedNewOrganizerID,
				Ownership: fakedNewOwnership,
			}, nil)
		mockedRepo.On("GetCategoriesByCompetition", mock.Anything, nil, fakedCompetitionID).
			Return([]domain.Category{fakedCategory}, nil)
		mockedRepo.
			On("GetProblemsByCompetition", mock.Anything, nil, fakedCompetitionID).
			Return([]domain.Problem{fakedProblem}, nil)
		mockedRepo.
			On("GetCompetitorsByCompetition", mock.Anything, nil, fakedCompetitionID).
			Return([]domain.Competitor{fakedCompetitor}, nil)
		mockedRepo.
			On("GetRafflesByCompetition", mock.Anything, nil, fakedCompetitionID).
			Return([]domain.Raffle{fakedRaffle}, nil)
		mockedRepo.
			On("GetRaffleWinners", mock.Anything, nil, fakedRaffleID).
			Return([]domain.RaffleWinner{fakedWinner}, nil)
		mockedRepo.
			On("GetTicksByCompetition", mock.Anything, nil, fakedCompetitionID).
			Return([]domain.Tick{fakedTick}, nil)

		mockedRepo.On("Begin").Return(mockedTx, nil)

		mockedRepo.On("DeleteTick", mock.Anything, mockedTx, fakedTickID).Return(nil)
		mockedRepo.On("DeleteRaffleWinner", mock.Anything, mockedTx, fakedRaffleWinnerID).Return(nil)
		mockedRepo.On("DeleteRaffle", mock.Anything, mockedTx, fakedRaffleID).Return(nil)
		mockedRepo.On("DeleteCompetitor", mock.Anything, mockedTx, fakedCompetitorID).Return(nil)
		mockedRepo.On("DeleteProblem", mock.Anything, mockedTx, fakedProblemID).Return(nil)
		mockedRepo.On("DeleteCategory", mock.Anything, mockedTx, fakedCategoryID).Return(nil)
		mockedRepo.On("DeleteCompetition", mock.Anything, mockedTx, fakedCompetitionID).Return(nil)

		mockedRepo.
			On("StoreCompetition", mock.Anything, mockedTx, domain.Competition{
				ID:                 fakedCompetitionID,
				Ownership:          fakedNewOwnership,
				Location:           "Stockholm Climbing Center",
				SeriesID:           fakedSeriesID,
				Name:               "Swedish Open 2025",
				Description:        "National bouldering championship",
				QualifyingProblems: 8,
				Finalists:          6,
				UsePoints:          true,
				PooledPoints:       true,
				MaxAttempts:        5,
				PointDeduction:     2,
				Info:               "Standard IFSC rules apply",
				GracePeriod:        30 * time.Minute,
				TimeBegin:          fakedCompetition.TimeBegin,
				TimeEnd:            fakedCompetition.TimeEnd,
				Created:            fakedCompetition.Created,
			}).
			Return(domain.Competition{}, nil)

		mockedRepo.
			On("StoreCategory", mock.Anything, mockedTx, domain.Category{
				ID:            fakedCategoryID,
				Ownership:     fakedNewOwnership,
				CompetitionID: fakedCompetitionID,
				Name:          "Males",
				Description:   "Male competitors",
				TimeBegin:     fakedCategory.TimeBegin,
				TimeEnd:       fakedCategory.TimeEnd,
			}).
			Return(domain.Category{}, nil)

		mockedRepo.
			On("StoreProblem", mock.Anything, mockedTx, domain.Problem{
				ID:                 fakedProblemID,
				Ownership:          fakedNewOwnership,
				CompetitionID:      fakedCompetitionID,
				Number:             15,
				HoldColorPrimary:   "#0000FF",
				HoldColorSecondary: "#FFFFFF",
				Description:        "Technical slab problem",
				Zone1Enabled:       true,
				Zone2Enabled:       true,
				ProblemValue: domain.ProblemValue{
					PointsTop:   100,
					PointsZone1: 25,
					PointsZone2: 50,
					FlashBonus:  10,
				},
			}).Return(domain.Problem{}, nil)

		mockedRepo.
			On("StoreCompetitor", mock.Anything, mockedTx, domain.Competitor{
				ID:                  fakedCompetitorID,
				Ownership:           domain.OwnershipData{OrganizerID: fakedNewOrganizerID, CompetitorID: &fakedCompetitorID},
				CompetitionID:       fakedCompetitionID,
				CategoryID:          fakedCategoryID,
				RegistrationCode:    "ABCD1234",
				Name:                "John Doe",
				Entered:             fakedCompetitor.Entered,
				WithdrawnFromFinals: true,
				Disqualified:        true,
				Result:              &fakedResult,
			}).
			Return(domain.Competitor{}, nil)

		mockedRepo.
			On("StoreResult", mock.Anything, mockedTx, fakedResult).
			Return(nil)

		mockedRepo.
			On("StoreRaffle", mock.Anything, mockedTx, domain.Raffle{
				ID:            fakedRaffleID,
				Ownership:     fakedNewOwnership,
				CompetitionID: fakedCompetitionID,
			}).
			Return(domain.Raffle{}, nil)

		mockedRepo.
			On("StoreRaffleWinner", mock.Anything, mockedTx, domain.RaffleWinner{
				ID:             fakedRaffleWinnerID,
				Ownership:      fakedNewOwnership,
				RaffleID:       fakedRaffleID,
				CompetitorID:   fakedCompetitorID,
				CompetitorName: "John Doe",
				Timestamp:      fakedWinner.Timestamp,
			}).
			Return(domain.RaffleWinner{}, nil)

		mockedRepo.
			On("StoreTick", mock.Anything, mockedTx, testutils.MustTick(domain.Tick{
				ID:            fakedTickID,
				Ownership:     domain.OwnershipData{OrganizerID: fakedNewOrganizerID, CompetitorID: &fakedCompetitorID},
				Timestamp:     fakedTick.Timestamp,
				CompetitionID: fakedCompetitionID,
				ProblemID:     fakedProblemID,
				Revision:      2,
				Zone1:         true,
				AttemptsZone1: 2,
				Zone2:         true,
				AttemptsZone2: 3,
				Top:           true,
				AttemptsTop:   5,
			})).
			Return(testutils.MustTick(domain.Tick{Revision: 2}), nil)

		mockedTx.On("Commit").Return(nil)

		ucase := usecases.CompetitionUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		competition, err := ucase.TransferCompetition(context.Background(), fakedCompetitionID, fakedNewOrganizerID)

		require.NoError(t, err)
		assert.Equal(t, fakedNewOrganizerID, competition.Ownership.OrganizerID)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
		mockedTx.AssertExpectations(t)
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedRepo, mockedAuthorizer, _ := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOldOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		mockedRepo.
			On("GetCompetition", mock.Anything, nil, fakedCompetitionID).
			Return(fakedCompetition, nil)

		ucase := usecases.CompetitionUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		_, err := ucase.TransferCompetition(context.Background(), fakedCompetitionID, fakedNewOrganizerID)

		assert.ErrorIs(t, err, domain.ErrNoOwnership)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("NewOrganizerUnauthorized", func(t *testing.T) {
		mockedRepo, mockedAuthorizer, _ := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOldOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("GetCompetition", mock.Anything, nil, fakedCompetitionID).
			Return(fakedCompetition, nil)

		mockedRepo.
			On("GetOrganizer", mock.Anything, nil, fakedNewOrganizerID).
			Return(domain.Organizer{
				ID:        fakedNewOrganizerID,
				Ownership: fakedNewOwnership,
			}, nil)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedNewOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		ucase := usecases.CompetitionUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		_, err := ucase.TransferCompetition(context.Background(), fakedCompetitionID, fakedNewOrganizerID)

		assert.ErrorIs(t, err, domain.ErrNoOwnership)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("CompetitionArchived", func(t *testing.T) {
		mockedRepo, _, _ := makeMocks()

		archivedCompetition := fakedCompetition
		archivedCompetition.ArchivedAt = time.Now()

		mockedRepo.
			On("GetCompetition", mock.Anything, nil, fakedCompetitionID).
			Return(archivedCompetition, nil)

		ucase := usecases.CompetitionUseCase{
			Repo: mockedRepo,
		}

		_, err := ucase.TransferCompetition(context.Background(), fakedCompetitionID, fakedNewOrganizerID)

		assert.ErrorIs(t, err, domain.ErrArchived)

		mockedRepo.AssertExpectations(t)
	})
}

func TestPatchCompetition(t *testing.T) {
	t.Parallel()

	fakedCompetitionID := testutils.RandomResourceID[domain.CompetitionID]()
	fakedOrganizerID := testutils.RandomResourceID[domain.OrganizerID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: fakedOrganizerID,
	}

	makeMocks := func() (*repositoryMock, *authorizerMock, *eventBrokerMock) {
		mockedRepo := new(repositoryMock)

		mockedAuthorizer := new(authorizerMock)

		mockedEventBroker := new(eventBrokerMock)

		return mockedRepo, mockedAuthorizer, mockedEventBroker
	}

	t.Run("HappyCase", func(t *testing.T) {
		mockedRepo, mockedAuthorizer, mockedEventBroker := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("GetCompetition", mock.Anything, nil, fakedCompetitionID).
			Return(domain.Competition{
				ID:                fakedCompetitionID,
				Ownership:         fakedOwnership,
				NameRetentionTime: 14 * 24 * time.Hour,
			}, nil)

		mockedRepo.
			On("StoreCompetition", mock.Anything, nil,
				domain.Competition{
					ID:                 fakedCompetitionID,
					Ownership:          fakedOwnership,
					Location:           "The garage",
					Country:            "SE",
					SeriesID:           domain.SeriesID(1),
					Name:               "Swedish Championships",
					Description:        "Who is the best climber in Sweden?",
					QualifyingProblems: 20,
					Finalists:          5,
					UsePoints:          true,
					PooledPoints:       true,
					MaxAttempts:        5,
					PointDeduction:     2,
					Info:               "No rules!",
					GracePeriod:        time.Hour,
					NameRetentionTime:  14 * 24 * time.Hour,
				},
			).
			Return(domain.Competition{
				ID:                 fakedCompetitionID,
				Ownership:          fakedOwnership,
				Location:           "The garage",
				Country:            "SE",
				SeriesID:           domain.SeriesID(1),
				Name:               "Swedish Championships",
				Description:        "Who is the best climber in Sweden?",
				QualifyingProblems: 20,
				Finalists:          5,
				UsePoints:          true,
				PooledPoints:       true,
				MaxAttempts:        5,
				PointDeduction:     2,
				Info:               "No rules!",
				GracePeriod:        time.Hour,
				NameRetentionTime:  14 * 24 * time.Hour,
			}, nil)

		mockedEventBroker.
			On("Dispatch", fakedCompetitionID, domain.RulesUpdatedEvent{
				CompetitionID:      fakedCompetitionID,
				QualifyingProblems: 20,
				Finalists:          5,
				UsePoints:          true,
				PooledPoints:       true,
				MaxAttempts:        5,
				PointDeduction:     2,
			}).
			Return(nil)

		ucase := usecases.CompetitionUseCase{
			Repo:        mockedRepo,
			Authorizer:  mockedAuthorizer,
			EventBroker: mockedEventBroker,
		}

		patch := domain.CompetitionPatch{
			Location:           domain.NewPatch("The garage"),
			Country:            domain.NewPatch("SE"),
			SeriesID:           domain.NewPatch(domain.SeriesID(1)),
			Name:               domain.NewPatch("Swedish Championships"),
			Description:        domain.NewPatch("Who is the best climber in Sweden?"),
			QualifyingProblems: domain.NewPatch(20),
			Finalists:          domain.NewPatch(5),
			UsePoints:          domain.NewPatch(true),
			PooledPoints:       domain.NewPatch(true),
			MaxAttempts:        domain.NewPatch(5),
			PointDeduction:     domain.NewPatch(domain.Tenths(2)),
			Info:               domain.NewPatch("No rules!"),
			GracePeriod:        domain.NewPatch(time.Hour),
		}

		competition, err := ucase.PatchCompetition(context.Background(), fakedCompetitionID, patch)

		require.NoError(t, err)
		assert.Equal(t, "The garage", competition.Location)
		assert.Equal(t, "SE", competition.Country)
		assert.Equal(t, domain.SeriesID(1), competition.SeriesID)
		assert.Equal(t, "Swedish Championships", competition.Name)
		assert.Equal(t, "Who is the best climber in Sweden?", competition.Description)
		assert.Equal(t, 20, competition.QualifyingProblems)
		assert.Equal(t, 5, competition.Finalists)
		assert.True(t, competition.UsePoints)
		assert.True(t, competition.PooledPoints)
		assert.Equal(t, 5, competition.MaxAttempts)
		assert.Equal(t, domain.Tenths(2), competition.PointDeduction)
		assert.Equal(t, "No rules!", competition.Info)
		assert.Equal(t, time.Hour, competition.GracePeriod)
		assert.Equal(t, 14*24*time.Hour, competition.NameRetentionTime)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
		mockedEventBroker.AssertExpectations(t)
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedRepo, mockedAuthorizer, _ := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		mockedRepo.
			On("GetCompetition", mock.Anything, nil, fakedCompetitionID).
			Return(domain.Competition{
				ID:        fakedCompetitionID,
				Ownership: fakedOwnership,
			}, nil)

		ucase := usecases.CompetitionUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		_, err := ucase.PatchCompetition(context.Background(), fakedCompetitionID, domain.CompetitionPatch{})

		assert.ErrorIs(t, err, domain.ErrNoOwnership)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("ValidatorIsInvoked", func(t *testing.T) {
		mockedRepo, mockedAuthorizer, _ := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("GetCompetition", mock.Anything, nil, fakedCompetitionID).
			Return(domain.Competition{
				ID:        fakedCompetitionID,
				Ownership: fakedOwnership,
			}, nil)

		ucase := usecases.CompetitionUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		_, err := ucase.PatchCompetition(context.Background(), fakedCompetitionID, domain.CompetitionPatch{})

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.True(t, validators.CompetitionValidator{}.IsValidationError(err))

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("CompetitionIsArchived", func(t *testing.T) {
		mockedRepo, mockedAuthorizer, _ := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("GetCompetition", mock.Anything, nil, fakedCompetitionID).
			Return(domain.Competition{
				ID:         fakedCompetitionID,
				Ownership:  fakedOwnership,
				ArchivedAt: time.Now(),
			}, nil)

		ucase := usecases.CompetitionUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		_, err := ucase.PatchCompetition(context.Background(), fakedCompetitionID, domain.CompetitionPatch{
			Name: domain.NewPatch("Norweigan Championships"),
		})

		assert.ErrorIs(t, err, domain.ErrArchived)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})
}

func TestArchiveCompetition(t *testing.T) {
	fakedCompetitionID := testutils.RandomResourceID[domain.CompetitionID]()
	fakedOrganizerID := testutils.RandomResourceID[domain.OrganizerID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: fakedOrganizerID,
	}

	t.Run("HappyCase", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			fakedScoreEngineInstanceID := domain.ScoreEngineInstanceID(uuid.New())
			mockedRepo := new(repositoryMock)
			mockedAuthorizer := new(authorizerMock)
			mockedScoreEngineManager := new(scoreEngineManagerMock)

			mockedAuthorizer.
				On("HasOwnership", mock.Anything, fakedOwnership).
				Return(domain.OrganizerRole, nil)

			mockedRepo.
				On("GetCompetition", mock.Anything, nil, fakedCompetitionID).
				Return(domain.Competition{
					ID:                fakedCompetitionID,
					Ownership:         fakedOwnership,
					Name:              "Swedish Championships",
					Country:           "SE",
					NameRetentionTime: 14 * 24 * time.Hour,
				}, nil)

			mockedRepo.
				On("StoreCompetition", mock.Anything, nil,
					domain.Competition{
						ID:                fakedCompetitionID,
						Ownership:         fakedOwnership,
						Name:              "Swedish Championships",
						Country:           "SE",
						ArchivedAt:        time.Now(),
						NameRetentionTime: 14 * 24 * time.Hour,
					},
				).
				Return(domain.Competition{
					ID:                fakedCompetitionID,
					Ownership:         fakedOwnership,
					ArchivedAt:        time.Now(),
					Name:              "Swedish Championships",
					Country:           "SE",
					NameRetentionTime: 14 * 24 * time.Hour,
				}, nil)

			mockedScoreEngineManager.
				On("ListScoreEnginesByCompetition", mock.Anything, fakedCompetitionID).
				Return([]domain.ScoreEngineDescriptor{
					{
						InstanceID:    fakedScoreEngineInstanceID,
						CompetitionID: fakedCompetitionID,
					},
				}, nil)

			mockedScoreEngineManager.
				On("StopScoreEngine", mock.Anything, fakedScoreEngineInstanceID).
				Return(nil)

			ucase := usecases.CompetitionUseCase{
				Repo:               mockedRepo,
				Authorizer:         mockedAuthorizer,
				ScoreEngineManager: mockedScoreEngineManager,
			}

			competition, err := ucase.ArchiveCompetition(context.Background(), fakedCompetitionID)

			require.NoError(t, err)
			assert.Equal(t, time.Now(), competition.ArchivedAt)

			mockedRepo.AssertExpectations(t)
			mockedAuthorizer.AssertExpectations(t)
			mockedScoreEngineManager.AssertExpectations(t)
		})
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedRepo := new(repositoryMock)
		mockedAuthorizer := new(authorizerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		mockedRepo.
			On("GetCompetition", mock.Anything, nil, fakedCompetitionID).
			Return(domain.Competition{
				ID:        fakedCompetitionID,
				Ownership: fakedOwnership,
			}, nil)

		ucase := usecases.CompetitionUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		_, err := ucase.ArchiveCompetition(context.Background(), fakedCompetitionID)

		require.ErrorIs(t, err, domain.ErrNoOwnership)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})
}

func TestRestoreCompetition(t *testing.T) {
	fakedCompetitionID := testutils.RandomResourceID[domain.CompetitionID]()
	fakedOrganizerID := testutils.RandomResourceID[domain.OrganizerID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: fakedOrganizerID,
	}

	t.Run("HappyCase", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			mockedRepo := new(repositoryMock)
			mockedAuthorizer := new(authorizerMock)

			mockedAuthorizer.
				On("HasOwnership", mock.Anything, fakedOwnership).
				Return(domain.OrganizerRole, nil)

			mockedRepo.
				On("GetCompetition", mock.Anything, nil, fakedCompetitionID).
				Return(domain.Competition{
					ID:                fakedCompetitionID,
					Ownership:         fakedOwnership,
					Name:              "Swedish Championships",
					Country:           "SE",
					ArchivedAt:        time.Now(),
					NameRetentionTime: 14 * 24 * time.Hour,
				}, nil)

			mockedRepo.
				On("StoreCompetition", mock.Anything, nil,
					domain.Competition{
						ID:                fakedCompetitionID,
						Ownership:         fakedOwnership,
						Name:              "Swedish Championships",
						Country:           "SE",
						NameRetentionTime: 14 * 24 * time.Hour,
					},
				).
				Return(domain.Competition{
					ID:                fakedCompetitionID,
					Ownership:         fakedOwnership,
					Name:              "Swedish Championships",
					Country:           "SE",
					NameRetentionTime: 14 * 24 * time.Hour,
				}, nil)

			ucase := usecases.CompetitionUseCase{
				Repo:       mockedRepo,
				Authorizer: mockedAuthorizer,
			}

			competition, err := ucase.RestoreCompetition(context.Background(), fakedCompetitionID)

			require.NoError(t, err)
			assert.Zero(t, competition.ArchivedAt)

			mockedRepo.AssertExpectations(t)
			mockedAuthorizer.AssertExpectations(t)
		})
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedRepo := new(repositoryMock)
		mockedAuthorizer := new(authorizerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		mockedRepo.
			On("GetCompetition", mock.Anything, nil, fakedCompetitionID).
			Return(domain.Competition{
				ID:        fakedCompetitionID,
				Ownership: fakedOwnership,
			}, nil)

		ucase := usecases.CompetitionUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		_, err := ucase.RestoreCompetition(context.Background(), fakedCompetitionID)

		require.ErrorIs(t, err, domain.ErrNoOwnership)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})
}
