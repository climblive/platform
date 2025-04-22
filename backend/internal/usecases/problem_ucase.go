package usecases

import (
	"context"
	"regexp"
	"strings"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

var validHexColor *regexp.Regexp = regexp.MustCompile(`^#([0-9a-fA-F]{3}){1,2}$`)

type problemUseCaseRepository interface {
	domain.Transactor

	GetProblemsByContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) ([]domain.Problem, error)
	StoreProblem(ctx context.Context, tx domain.Transaction, problem domain.Problem) (domain.Problem, error)
	GetProblem(ctx context.Context, tx domain.Transaction, problemID domain.ProblemID) (domain.Problem, error)
	GetProblemByNumber(ctx context.Context, tx domain.Transaction, contestID domain.ContestID, problemNumber int) (domain.Problem, error)
	GetContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) (domain.Contest, error)
}

type ProblemUseCase struct {
	Authorizer  domain.Authorizer
	Repo        problemUseCaseRepository
	EventBroker domain.EventBroker
}

func (uc *ProblemUseCase) GetProblem(ctx context.Context, problemID domain.ProblemID) (domain.Problem, error) {
	problem, err := uc.Repo.GetProblem(ctx, nil, problemID)
	if err != nil {
		return domain.Problem{}, errors.Wrap(err, 0)
	}

	return problem, nil
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
		_, err = uc.Repo.GetProblemByNumber(ctx, nil, problem.ContestID, patch.Number.Value)
		switch {
		case err == nil:
			return domain.Problem{}, errors.Wrap(domain.ErrDuplicate, 0)
		case errors.Is(err, domain.ErrNotFound):
		default:
			return domain.Problem{}, errors.Wrap(err, 0)
		}

		problem.Number = patch.Number.Value
	}

	if patch.HoldColorPrimary.Present {
		problem.HoldColorPrimary = strings.TrimSpace(patch.HoldColorPrimary.Value)

		if !validHexColor.MatchString(problem.HoldColorPrimary) {
			return domain.Problem{}, errors.Wrap(domain.ErrInvalidData, 0)
		}
	}

	if patch.HoldColorSecondary.Present {
		problem.HoldColorSecondary = strings.TrimSpace(patch.HoldColorSecondary.Value)

		if len(problem.HoldColorSecondary) > 0 && !validHexColor.MatchString(problem.HoldColorSecondary) {
			return domain.Problem{}, errors.Wrap(domain.ErrInvalidData, 0)
		}
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

func (uc *ProblemUseCase) CreateProblem(ctx context.Context, contestID domain.ContestID, tmpl domain.ProblemTemplate) (domain.Problem, error) {
	contest, err := uc.Repo.GetContest(ctx, nil, contestID)
	if err != nil {
		return domain.Problem{}, errors.Wrap(err, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, contest.Ownership); err != nil {
		return domain.Problem{}, errors.Wrap(err, 0)
	}

	_, err = uc.Repo.GetProblemByNumber(ctx, nil, contestID, tmpl.Number)
	switch {
	case err == nil:
		return domain.Problem{}, errors.Wrap(domain.ErrDuplicate, 0)
	case errors.Is(err, domain.ErrNotFound):
	default:
		return domain.Problem{}, errors.Wrap(err, 0)
	}

	problem := domain.Problem{
		Ownership:          contest.Ownership,
		ContestID:          contestID,
		Number:             tmpl.Number,
		HoldColorPrimary:   strings.TrimSpace(tmpl.HoldColorPrimary),
		HoldColorSecondary: strings.TrimSpace(tmpl.HoldColorSecondary),
		Description:        strings.TrimSpace(tmpl.Description),
		PointsTop:          tmpl.PointsTop,
		PointsZone:         tmpl.PointsZone,
		FlashBonus:         tmpl.FlashBonus,
	}

	switch {
	case !validHexColor.MatchString(problem.HoldColorPrimary):
		fallthrough
	case len(problem.HoldColorSecondary) > 0 && !validHexColor.MatchString(tmpl.HoldColorSecondary):
		return domain.Problem{}, errors.Wrap(domain.ErrInvalidData, 0)
	}

	createdProblem, err := uc.Repo.StoreProblem(ctx, nil, problem)
	if err != nil {
		return domain.Problem{}, errors.Wrap(err, 0)
	}

	event := domain.ProblemAddedEvent{
		ProblemID:  createdProblem.ID,
		PointsTop:  problem.PointsTop,
		PointsZone: problem.PointsZone,
		FlashBonus: problem.FlashBonus,
	}

	uc.EventBroker.Dispatch(problem.ContestID, event)

	return createdProblem, nil
}
