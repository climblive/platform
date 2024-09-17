package usecases

import (
	"context"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

type tickUseCaseRepository interface {
	domain.Transactor

	GetContender(ctx context.Context, tx domain.Transaction, contenderID domain.ResourceID) (domain.Contender, error)
	GetTicksByContender(ctx context.Context, tx domain.Transaction, contenderID domain.ResourceID) ([]domain.Tick, error)
	GetContest(ctx context.Context, tx domain.Transaction, contestID domain.ResourceID) (domain.Contest, error)
	GetCompClass(ctx context.Context, tx domain.Transaction, compClassID domain.ResourceID) (domain.CompClass, error)
	GetProblem(ctx context.Context, tx domain.Transaction, problemID domain.ResourceID) (domain.Problem, error)
	DeleteTick(ctx context.Context, tx domain.Transaction, tickID domain.ResourceID) error
	CreateTick(ctx context.Context, tx domain.Transaction, tick domain.Tick) (domain.Tick, error)
}

type TickUseCase struct {
	Repo        tickUseCaseRepository
	Authorizer  domain.Authorizer
	EventBroker domain.EventBroker
}

func (uc *TickUseCase) GetTicksByContender(ctx context.Context, contenderID domain.ResourceID) ([]domain.Tick, error) {
	contender, err := uc.Repo.GetContender(ctx, nil, contenderID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, contender.Ownership); err != nil {
		return nil, errors.Wrap(err, 0)
	}

	ticks, err := uc.Repo.GetTicksByContender(ctx, nil, contenderID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	return ticks, nil
}

func (uc *TickUseCase) GetTicksByProblem(ctx context.Context, problemID domain.ResourceID) ([]domain.Tick, error) {
	panic("not implemented")
}

func (uc *TickUseCase) DeleteTick(ctx context.Context, tickID domain.ResourceID) error {
	panic("not implemented")
}

func (uc *TickUseCase) CreateTick(ctx context.Context, contenderID domain.ResourceID, tick domain.Tick) (domain.Tick, error) {
	contender, err := uc.Repo.GetContender(ctx, nil, contenderID)
	if err != nil {
		return domain.Tick{}, errors.Wrap(err, 0)
	}

	role, err := uc.Authorizer.HasOwnership(ctx, contender.Ownership)
	if err != nil {
		return domain.Tick{}, errors.Wrap(err, 0)
	}

	problem, err := uc.Repo.GetProblem(ctx, nil, tick.ProblemID)
	if err != nil {
		return domain.Tick{}, errors.Wrap(err, 0)
	}

	contest, err := uc.Repo.GetContest(ctx, nil, contender.ContestID)
	if err != nil {
		return domain.Tick{}, errors.Errorf("%w: %w", domain.ErrRepositoryIntegrityViolation, err)
	}

	compClass, err := uc.Repo.GetCompClass(ctx, nil, contender.CompClassID)
	if err != nil {
		return domain.Tick{}, errors.Errorf("%w: %w", domain.ErrRepositoryIntegrityViolation, err)
	}

	gracePeriodEnd := compClass.TimeEnd.Add(contest.GracePeriod)

	switch {
	case role.OneOf(domain.OrganizerRole, domain.AdminRole):
	case time.Now().After(gracePeriodEnd):
		return domain.Tick{}, errors.New(domain.ErrContestEnded)
	}

	newTick := domain.Tick{
		Ownership:    contender.Ownership,
		Timestamp:    time.Now(),
		ContestID:    contest.ID,
		ProblemID:    problem.ID,
		Top:          tick.Top,
		AttemptsTop:  tick.AttemptsTop,
		Zone:         tick.Zone,
		AttemptsZone: tick.AttemptsZone,
	}

	tick, err = uc.Repo.CreateTick(ctx, nil, newTick)
	if err != nil {
		return domain.Tick{}, errors.Wrap(err, 0)
	}

	uc.EventBroker.Dispatch(contest.ID, domain.AscentRegisteredEvent{
		ContenderID:  contender.ID,
		ProblemID:    problem.ID,
		Top:          tick.Top,
		AttemptsTop:  tick.AttemptsTop,
		Zone:         tick.Zone,
		AttemptsZone: tick.AttemptsTop,
	})

	return tick, nil
}
