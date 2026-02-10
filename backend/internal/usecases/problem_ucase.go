package usecases

import (
	"context"
	"strings"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/usecases/validators"
	"github.com/go-errors/errors"
)

type problemUseCaseRepository interface {
	domain.Transactor

	GetProblemsByContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) ([]domain.Problem, error)
	StoreProblem(ctx context.Context, tx domain.Transaction, problem domain.Problem) (domain.Problem, error)
	GetProblem(ctx context.Context, tx domain.Transaction, problemID domain.ProblemID) (domain.Problem, error)
	GetProblemByNumber(ctx context.Context, tx domain.Transaction, contestID domain.ContestID, problemNumber int) (domain.Problem, error)
	GetContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) (domain.Contest, error)
	DeleteProblem(ctx context.Context, tx domain.Transaction, problemID domain.ProblemID) error
	GetTicksByProblem(ctx context.Context, tx domain.Transaction, problemID domain.ProblemID) ([]domain.Tick, error)
	GetCompClass(ctx context.Context, tx domain.Transaction, compClassID domain.CompClassID) (domain.CompClass, error)
}

type ProblemUseCase struct {
	Authorizer         domain.Authorizer
	Repo               problemUseCaseRepository
	EventBroker        domain.EventBroker
	ProblemValueKeeper domain.ProblemValueKeeper
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

func (uc *ProblemUseCase) GetProblemsByCompClass(ctx context.Context, compClassID domain.CompClassID) ([]domain.Problem, error) {
	compClass, err := uc.Repo.GetCompClass(ctx, nil, compClassID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	problems, err := uc.Repo.GetProblemsByContest(ctx, nil, compClass.ContestID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	for i, problem := range problems {
		problems[i] = withProblemValue(problem, compClassID, uc.ProblemValueKeeper)
	}

	return problems, nil
}

func withProblemValue(problem domain.Problem, compClassID domain.CompClassID, keeper domain.ProblemValueKeeper) domain.Problem {
	if keeper == nil {
		return problem
	}

	if value, found := keeper.GetProblemValue(problem.ID, compClassID); found {
		problem.ProblemValue = value
	}

	return problem
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
		ProblemID:    problemID,
		ProblemValue: problem.ProblemValue,
	}

	if patch.Number.PresentAndDistinct(problem.Number) {
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
	}

	if patch.HoldColorSecondary.Present {
		problem.HoldColorSecondary = strings.TrimSpace(patch.HoldColorSecondary.Value)
	}

	if patch.Description.Present {
		problem.Description = strings.TrimSpace(patch.Description.Value)
	}

	if patch.Zone1Enabled.Present {
		problem.Zone1Enabled = patch.Zone1Enabled.Value
	}

	if patch.Zone2Enabled.Present {
		problem.Zone2Enabled = patch.Zone2Enabled.Value
	}

	if patch.PointsTop.Present {
		problem.PointsTop = patch.PointsTop.Value
	}

	if patch.PointsZone1.Present {
		problem.PointsZone1 = patch.PointsZone1.Value
	}

	if patch.PointsZone2.Present {
		problem.PointsZone2 = patch.PointsZone2.Value
	}

	if patch.FlashBonus.Present {
		problem.FlashBonus = patch.FlashBonus.Value
	}

	if err := (validators.ProblemValidator{}).Validate(problem); err != nil {
		return mty, errors.Wrap(err, 0)
	}

	if _, err = uc.Repo.StoreProblem(ctx, nil, problem); err != nil {
		return mty, errors.Wrap(err, 0)
	}

	event := domain.ProblemUpdatedEvent{
		ProblemID:    problemID,
		ProblemValue: problem.ProblemValue,
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
		ID:                 0,
		Ownership:          contest.Ownership,
		ContestID:          contestID,
		Number:             tmpl.Number,
		HoldColorPrimary:   strings.TrimSpace(tmpl.HoldColorPrimary),
		HoldColorSecondary: strings.TrimSpace(tmpl.HoldColorSecondary),
		Description:        strings.TrimSpace(tmpl.Description),
		Zone1Enabled:       tmpl.Zone1Enabled,
		Zone2Enabled:       tmpl.Zone2Enabled,
		ProblemValue: domain.ProblemValue{
			PointsZone1: tmpl.PointsZone1,
			PointsZone2: tmpl.PointsZone2,
			PointsTop:   tmpl.PointsTop,
			FlashBonus:  tmpl.FlashBonus,
		},
	}

	if err := (validators.ProblemValidator{}).Validate(problem); err != nil {
		return domain.Problem{}, errors.Wrap(err, 0)
	}

	createdProblem, err := uc.Repo.StoreProblem(ctx, nil, problem)
	if err != nil {
		return domain.Problem{}, errors.Wrap(err, 0)
	}

	event := domain.ProblemAddedEvent{
		ProblemID: createdProblem.ID,
		ProblemValue: domain.ProblemValue{
			PointsZone1: problem.PointsZone1,
			PointsZone2: problem.PointsZone2,
			PointsTop:   problem.PointsTop,
			FlashBonus:  problem.FlashBonus,
		},
	}

	uc.EventBroker.Dispatch(problem.ContestID, event)

	return createdProblem, nil
}

func (uc *ProblemUseCase) DeleteProblem(ctx context.Context, problemID domain.ProblemID) error {
	problem, err := uc.Repo.GetProblem(ctx, nil, problemID)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, problem.Ownership); err != nil {
		return errors.Wrap(err, 0)
	}

	ticks, err := uc.Repo.GetTicksByProblem(ctx, nil, problemID)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	if len(ticks) > 0 {
		return errors.Wrap(domain.ErrNotAllowed, 0)
	}

	err = uc.Repo.DeleteProblem(ctx, nil, problemID)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	event := domain.ProblemDeletedEvent{
		ProblemID: problem.ID,
	}

	uc.EventBroker.Dispatch(problem.ContestID, event)

	return nil
}
