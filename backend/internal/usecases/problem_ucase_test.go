package usecases_test

import (
	"context"
	"testing"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/usecases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetProblemsByContest(t *testing.T) {
	fakedContestID := randomResourceID[domain.ContestID]()
	fakedProblems := []domain.Problem{
		{
			ID:        randomResourceID[domain.ProblemID](),
			ContestID: fakedContestID,
		},
	}

	mockedRepo := new(repositoryMock)

	mockedRepo.
		On("GetProblemsByContest", mock.Anything, mock.Anything, fakedContestID).
		Return(fakedProblems, nil)

	ucase := usecases.ProblemUseCase{
		Repo: mockedRepo,
	}

	problems, err := ucase.GetProblemsByContest(context.Background(), fakedContestID)

	require.NoError(t, err)
	assert.Equal(t, fakedProblems, problems)

	mockedRepo.AssertExpectations(t)
}

func TestPatchProblem(t *testing.T) {
	fakedProblemID := randomResourceID[domain.ProblemID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: randomResourceID[domain.OrganizerID](),
	}
	fakedContestID := randomResourceID[domain.ContestID]()

	fakedProblem := domain.Problem{
		ID:                 fakedProblemID,
		Ownership:          fakedOwnership,
		ContestID:          fakedContestID,
		Number:             10,
		HoldColorPrimary:   "#ffffff",
		HoldColorSecondary: "#000000",
		Description:        "The tenth boulder",
		PointsTop:          100,
		PointsZone:         50,
		FlashBonus:         10,
	}

	makeMocks := func() (*repositoryMock, *eventBrokerMock, *authorizerMock) {
		mockedRepo := new(repositoryMock)

		mockedRepo.
			On("GetProblem", mock.Anything, nil, fakedProblemID).
			Return(fakedProblem, nil)

		mockedAuthorizer := new(authorizerMock)

		mockedEventBroker := new(eventBrokerMock)

		return mockedRepo, mockedEventBroker, mockedAuthorizer
	}

	t.Run("HappyCase", func(t *testing.T) {
		mockedRepo, mockedEventBroker, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("GetProblemByNumber", mock.Anything, nil, fakedContestID, 20).
			Return(domain.Problem{}, domain.ErrNotFound)

		mockedRepo.
			On("StoreProblem", mock.Anything, nil, domain.Problem{
				ID:                 fakedProblemID,
				Ownership:          fakedOwnership,
				ContestID:          fakedContestID,
				Number:             20,
				HoldColorPrimary:   "#ff0000",
				HoldColorSecondary: "",
				Description:        "The twentieth boulder",
				PointsTop:          1000,
				PointsZone:         500,
				FlashBonus:         25,
			}).
			Return(domain.Problem{
				ID:                 fakedProblemID,
				Ownership:          fakedOwnership,
				ContestID:          fakedContestID,
				Number:             20,
				HoldColorPrimary:   "#ff0000",
				HoldColorSecondary: "",
				Description:        "The twentieth boulder",
				PointsTop:          1000,
				PointsZone:         500,
				FlashBonus:         25,
			}, nil)

		mockedEventBroker.
			On("Dispatch", fakedContestID, domain.ProblemUpdatedEvent{
				ProblemID:  fakedProblemID,
				PointsTop:  1000,
				PointsZone: 500,
				FlashBonus: 25,
			}).Return()

		ucase := usecases.ProblemUseCase{
			Repo:        mockedRepo,
			Authorizer:  mockedAuthorizer,
			EventBroker: mockedEventBroker,
		}

		problem, err := ucase.PatchProblem(context.Background(), fakedProblemID, domain.ProblemPatch{
			Number:             domain.NewPatch(20),
			HoldColorPrimary:   domain.NewPatch("#ff0000"),
			HoldColorSecondary: domain.NewPatch(""),
			Name:               domain.NewPatch("Boulder #20"),
			Description:        domain.NewPatch("The twentieth boulder"),
			PointsTop:          domain.NewPatch(1000),
			PointsZone:         domain.NewPatch(500),
			FlashBonus:         domain.NewPatch(25),
		})

		require.NoError(t, err)
		assert.Equal(t, 20, problem.Number)
		assert.Equal(t, "#ff0000", problem.HoldColorPrimary)
		assert.Equal(t, "", problem.HoldColorSecondary)
		assert.Equal(t, "The twentieth boulder", problem.Description)
		assert.Equal(t, 1000, problem.PointsTop)
		assert.Equal(t, 500, problem.PointsZone)
		assert.Equal(t, 25, problem.FlashBonus)

		mockedRepo.AssertExpectations(t)
		mockedEventBroker.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("NumberAlreadyTaken", func(t *testing.T) {
		mockedRepo, mockedEventBroker, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("GetProblemByNumber", mock.Anything, nil, fakedContestID, 20).
			Return(domain.Problem{}, nil)

		ucase := usecases.ProblemUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		_, err := ucase.PatchProblem(context.Background(), fakedProblemID, domain.ProblemPatch{
			Number: domain.NewPatch(20),
		})

		require.ErrorIs(t, err, domain.ErrDuplicate)

		mockedRepo.AssertExpectations(t)
		mockedEventBroker.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("NoEventDispatchedWhenPointsUnchanged", func(t *testing.T) {
		mockedRepo, mockedEventBroker, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("StoreProblem", mock.Anything, nil, mock.AnythingOfType("domain.Problem")).
			Return(domain.Problem{}, nil)

		ucase := usecases.ProblemUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		_, err := ucase.PatchProblem(context.Background(), fakedProblemID, domain.ProblemPatch{
			PointsTop:  domain.NewPatch(100),
			PointsZone: domain.NewPatch(50),
			FlashBonus: domain.NewPatch(10),
		})

		require.NoError(t, err)

		mockedRepo.AssertExpectations(t)
		mockedEventBroker.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedRepo, _, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		ucase := usecases.ProblemUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		_, err := ucase.PatchProblem(context.Background(), fakedProblemID, domain.ProblemPatch{})

		require.ErrorIs(t, err, domain.ErrNoOwnership)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})
}

func TestCreateProblem(t *testing.T) {
	fakedOrganizerID := randomResourceID[domain.OrganizerID]()
	fakedOwnership := domain.OwnershipData{
		OrganizerID: fakedOrganizerID,
	}
	fakedContestID := randomResourceID[domain.ContestID]()
	fakedProblemID := randomResourceID[domain.ProblemID]()

	makeMocks := func() (*repositoryMock, *authorizerMock) {
		mockedRepo := new(repositoryMock)

		mockedRepo.
			On("GetContest", mock.Anything, nil, fakedContestID).
			Return(domain.Contest{
				ID:        fakedContestID,
				Ownership: fakedOwnership,
			}, nil)

		mockedAuthorizer := new(authorizerMock)

		return mockedRepo, mockedAuthorizer
	}

	t.Run("HappyCase", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

		mockedEventBroker := new(eventBrokerMock)

		mockedEventBroker.
			On("Dispatch", fakedContestID, domain.ProblemAddedEvent{
				ProblemID:  fakedProblemID,
				PointsTop:  100,
				PointsZone: 50,
				FlashBonus: 15,
			}).
			Return()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("GetProblemByNumber", mock.Anything, nil, fakedContestID, 10).
			Return(domain.Problem{}, domain.ErrNotFound)

		mockedRepo.
			On("StoreProblem", mock.Anything, nil,
				domain.Problem{
					Ownership:          fakedOwnership,
					ContestID:          fakedContestID,
					Number:             10,
					HoldColorPrimary:   "#ffffff",
					HoldColorSecondary: "#000",
					Description:        "Crack volumes are included",
					PointsTop:          100,
					PointsZone:         50,
					FlashBonus:         15,
				},
			).
			Return(
				domain.Problem{
					ID:                 fakedProblemID,
					Ownership:          fakedOwnership,
					ContestID:          fakedContestID,
					Number:             10,
					HoldColorPrimary:   "#ffffff",
					HoldColorSecondary: "#000",
					Description:        "Crack volumes are included",
					PointsTop:          100,
					PointsZone:         50,
					FlashBonus:         15,
				}, nil)

		ucase := usecases.ProblemUseCase{
			Repo:        mockedRepo,
			Authorizer:  mockedAuthorizer,
			EventBroker: mockedEventBroker,
		}

		problem, err := ucase.CreateProblem(context.Background(), fakedContestID, domain.ProblemTemplate{
			Number:             10,
			HoldColorPrimary:   "#ffffff",
			HoldColorSecondary: "#000",
			Description:        "Crack volumes are included",
			PointsTop:          100,
			PointsZone:         50,
			FlashBonus:         15,
		})

		require.NoError(t, err)
		assert.Equal(t, fakedProblemID, problem.ID)
		assert.Equal(t, fakedOwnership, problem.Ownership)
		assert.Equal(t, fakedContestID, problem.ContestID)
		assert.Equal(t, 10, problem.Number)
		assert.Equal(t, "#ffffff", problem.HoldColorPrimary)
		assert.Equal(t, "#000", problem.HoldColorSecondary)
		assert.Equal(t, "Crack volumes are included", problem.Description)
		assert.Equal(t, 100, problem.PointsTop)
		assert.Equal(t, 50, problem.PointsZone)
		assert.Equal(t, 15, problem.FlashBonus)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
		mockedEventBroker.AssertExpectations(t)
	})

	t.Run("NumberAlreadyUsed", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("GetProblemByNumber", mock.Anything, nil, fakedContestID, 10).
			Return(domain.Problem{}, nil)

		ucase := usecases.ProblemUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		_, err := ucase.CreateProblem(context.Background(), fakedContestID, domain.ProblemTemplate{
			Number:             10,
			HoldColorPrimary:   "#ffffff",
			HoldColorSecondary: "#000",
			Description:        "Crack volumes are included",
			PointsTop:          100,
			PointsZone:         50,
			FlashBonus:         15,
		})

		require.ErrorIs(t, err, domain.ErrDuplicate)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})

	t.Run("InvalidData", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()
		mockedEventBroker := new(eventBrokerMock)

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.OrganizerRole, nil)

		mockedRepo.
			On("GetProblemByNumber", mock.Anything, nil, fakedContestID, 10).
			Return(domain.Problem{}, domain.ErrNotFound)

		mockedRepo.
			On("StoreProblem", mock.Anything, nil, mock.AnythingOfType("domain.Problem")).
			Return(domain.Problem{}, nil)

		mockedEventBroker.
			On("Dispatch", fakedContestID, mock.AnythingOfType("domain.ProblemAddedEvent")).
			Return()

		ucase := usecases.ProblemUseCase{
			Repo:        mockedRepo,
			Authorizer:  mockedAuthorizer,
			EventBroker: mockedEventBroker,
		}

		validTemplate := func() domain.ProblemTemplate {
			return domain.ProblemTemplate{
				Number:             10,
				HoldColorPrimary:   "#ffffff",
				HoldColorSecondary: "#000",
				Description:        "Crack volumes are included",
				PointsTop:          100,
				PointsZone:         50,
				FlashBonus:         15,
			}
		}

		_, err := ucase.CreateProblem(context.Background(), fakedContestID, validTemplate())

		require.NoError(t, err)

		t.Run("BadPrimaryColor", func(t *testing.T) {
			tmpl := validTemplate()
			tmpl.HoldColorPrimary = "invalid"

			_, err := ucase.CreateProblem(context.Background(), fakedContestID, tmpl)

			assert.ErrorIs(t, err, domain.ErrInvalidData)
		})

		t.Run("BadSecondaryColor", func(t *testing.T) {
			tmpl := validTemplate()
			tmpl.HoldColorSecondary = "invalid"

			_, err := ucase.CreateProblem(context.Background(), fakedContestID, tmpl)

			assert.ErrorIs(t, err, domain.ErrInvalidData)
		})

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
		mockedEventBroker.AssertExpectations(t)
	})

	t.Run("BadCredentials", func(t *testing.T) {
		mockedRepo, mockedAuthorizer := makeMocks()

		mockedAuthorizer.
			On("HasOwnership", mock.Anything, fakedOwnership).
			Return(domain.NilRole, domain.ErrNoOwnership)

		ucase := usecases.ProblemUseCase{
			Repo:       mockedRepo,
			Authorizer: mockedAuthorizer,
		}

		_, err := ucase.CreateProblem(context.Background(), fakedContestID, domain.ProblemTemplate{})

		require.ErrorIs(t, err, domain.ErrNoOwnership)

		mockedRepo.AssertExpectations(t)
		mockedAuthorizer.AssertExpectations(t)
	})
}
