package usecases

import (
	"context"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

type tickUseCaseRepository interface {
	domain.Transactor

	GetContender(ctx context.Context, tx domain.Transaction, contenderID domain.ContenderID) (domain.Contender, error)
	GetTicksByContender(ctx context.Context, tx domain.Transaction, contenderID domain.ContenderID) ([]domain.Tick, error)
	GetContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) (domain.Contest, error)
	GetCompClass(ctx context.Context, tx domain.Transaction, compClassID domain.CompClassID) (domain.CompClass, error)
	GetProblem(ctx context.Context, tx domain.Transaction, problemID domain.ProblemID) (domain.Problem, error)
	DeleteTick(ctx context.Context, tx domain.Transaction, tickID domain.TickID) error
	StoreTick(ctx context.Context, tx domain.Transaction, tick domain.Tick) (domain.Tick, error)
	GetTick(ctx context.Context, tx domain.Transaction, tickID domain.TickID) (domain.Tick, error)
}

type TickUseCase struct {
	Repo        tickUseCaseRepository
	Authorizer  domain.Authorizer
	EventBroker domain.EventBroker
}

func (uc *TickUseCase) GetTicksByContender(ctx context.Context, contenderID domain.ContenderID) ([]domain.Tick, error) {
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

func (uc *TickUseCase) GetTicksByProblem(ctx context.Context, problemID domain.ProblemID) ([]domain.Tick, error) {
	panic("not implemented")
}

func (uc *TickUseCase) DeleteTick(ctx context.Context, tickID domain.TickID) error {
	tick, err := uc.Repo.GetTick(ctx, nil, tickID)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	role, err := uc.Authorizer.HasOwnership(ctx, tick.Ownership)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	contenderID := *tick.Ownership.ContenderID

	contender, err := uc.Repo.GetContender(ctx, nil, contenderID)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	contest, err := uc.Repo.GetContest(ctx, nil, tick.ContestID)
	if err != nil {
		return errors.Errorf("%w: %w", domain.ErrRepositoryIntegrityViolation, err)
	}

	compClass, err := uc.Repo.GetCompClass(ctx, nil, contender.CompClassID)
	if err != nil {
		return errors.Errorf("%w: %w", domain.ErrRepositoryIntegrityViolation, err)
	}

	gracePeriodEnd := compClass.TimeEnd.Add(contest.GracePeriod)

	switch {
	case role.OneOf(domain.OrganizerRole, domain.AdminRole):
	case time.Now().After(gracePeriodEnd):
		return errors.New(domain.ErrContestEnded)
	}

	err = uc.Repo.DeleteTick(ctx, nil, tickID)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	uc.EventBroker.Dispatch(contest.ID, domain.AscentDeregisteredEvent{
		TickID:      tickID,
		ContenderID: contender.ID,
		ProblemID:   tick.ProblemID,
	})

	return nil
}

func (uc *TickUseCase) CreateTick(ctx context.Context, contenderID domain.ContenderID, tick domain.Tick) (domain.Tick, error) {
	contender, err := uc.Repo.GetContender(ctx, nil, contenderID)
	if err != nil {
		return domain.Tick{}, errors.Wrap(err, 0)
	}

	role, err := uc.Authorizer.HasOwnership(ctx, contender.Ownership)
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

	problem, err := uc.Repo.GetProblem(ctx, nil, tick.ProblemID)
	if err != nil {
		return domain.Tick{}, errors.Wrap(err, 0)
	}

	if problem.ContestID != contest.ID {
		return domain.Tick{}, errors.New(domain.ErrProblemNotInContest)
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

	tick, err = uc.Repo.StoreTick(ctx, nil, newTick)
	if err != nil {
		return domain.Tick{}, errors.Wrap(err, 0)
	}

	uc.EventBroker.Dispatch(contest.ID, domain.AscentRegisteredEvent{
		TickID:       tick.ID,
		Timestamp:    tick.Timestamp,
		ContenderID:  contender.ID,
		ProblemID:    problem.ID,
		Top:          tick.Top,
		AttemptsTop:  tick.AttemptsTop,
		Zone:         tick.Zone,
		AttemptsZone: tick.AttemptsZone,
	})

	return tick, nil
}
