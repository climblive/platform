package usecases

import (
	"context"
	"strings"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

type problemUseCaseRepository interface {
	domain.Transactor

	GetProblemsByContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) ([]domain.Problem, error)
	StoreProblem(ctx context.Context, tx domain.Transaction, problem domain.Problem) (domain.Problem, error)
	GetProblem(ctx context.Context, tx domain.Transaction, problemID domain.ProblemID) (domain.Problem, error)
}

type ProblemUseCase struct {
	Authorizer  domain.Authorizer
	Repo        problemUseCaseRepository
	EventBroker domain.EventBroker
}

func (uc *ProblemUseCase) GetProblemsByContest(ctx context.Context, contestID domain.ContestID) ([]domain.Problem, error) {
	problems, err := uc.Repo.GetProblemsByContest(ctx, nil, contestID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	return problems, nil
}

func (uc *ProblemUseCase) PatchProblem(ctx context.Context, problemID domain.ProblemID, patch domain.ProblemPatch) (domain.Problem, error) {
	var mty domain.Problem

	problem, err := uc.Repo.GetProblem(ctx, nil, problemID)
	if err != nil {
		return mty, errors.Wrap(err, 0)
	}

	_, err = uc.Authorizer.HasOwnership(ctx, problem.Ownership)
	if err != nil {
		return mty, errors.Wrap(err, 0)
	}

	problemUpdatedEventBaseline := domain.ProblemUpdatedEvent{
		ProblemID:  problemID,
		PointsTop:  problem.PointsTop,
		PointsZone: problem.PointsZone,
		FlashBonus: problem.FlashBonus,
	}

	if patch.Number.Present {
		problem.Number = patch.Number.Value
	}

	if patch.HoldColorPrimary.Present {
		problem.HoldColorPrimary = strings.TrimSpace(patch.HoldColorPrimary.Value)
	}

	if patch.HoldColorSecondary.Present {
		problem.HoldColorSecondary = strings.TrimSpace(patch.HoldColorSecondary.Value)
	}

	if patch.Name.Present {
		problem.Name = strings.TrimSpace(patch.Name.Value)
	}

	if patch.Description.Present {
		problem.Description = strings.TrimSpace(patch.Description.Value)
	}

	if patch.PointsTop.Present {
		problem.PointsTop = patch.PointsTop.Value
	}

	if patch.PointsZone.Present {
		problem.PointsZone = patch.PointsZone.Value
	}

	if patch.FlashBonus.Present {
		problem.FlashBonus = patch.FlashBonus.Value
	}

	if _, err = uc.Repo.StoreProblem(ctx, nil, problem); err != nil {
		return mty, errors.Wrap(err, 0)
	}

	event := domain.ProblemUpdatedEvent{
		ProblemID:  problemID,
		PointsTop:  problem.PointsTop,
		PointsZone: problem.PointsZone,
		FlashBonus: problem.FlashBonus,
	}

	if event != problemUpdatedEventBaseline {
		uc.EventBroker.Dispatch(problem.ContestID, event)
	}

	return problem, nil
}
