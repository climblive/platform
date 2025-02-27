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
		Name:               "Boulder #10",
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
			On("StoreProblem", mock.Anything, nil, domain.Problem{
				ID:                 fakedProblemID,
				Ownership:          fakedOwnership,
				ContestID:          fakedContestID,
				Number:             20,
				HoldColorPrimary:   "#ff0000",
				HoldColorSecondary: "",
				Name:               "Boulder #20",
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
				Name:               "Boulder #20",
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
		assert.Equal(t, "Boulder #20", problem.Name)
		assert.Equal(t, "The twentieth boulder", problem.Description)
		assert.Equal(t, 1000, problem.PointsTop)
		assert.Equal(t, 500, problem.PointsZone)
		assert.Equal(t, 25, problem.FlashBonus)

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
