package usecases

import (
	"context"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/usecases/validators"
	"github.com/go-errors/errors"
)

type tickUseCaseRepository interface {
	domain.Transactor

	GetContender(ctx context.Context, tx domain.Transaction, contenderID domain.ContenderID) (domain.Contender, error)
	GetTicksByContender(ctx context.Context, tx domain.Transaction, contenderID domain.ContenderID) ([]domain.Tick, error)
	GetTicksByContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) ([]domain.Tick, error)
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

func (uc *TickUseCase) GetTicksByContest(ctx context.Context, contestID domain.ContestID) ([]domain.Tick, error) {
	contest, err := uc.Repo.GetContest(ctx, nil, contestID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, contest.Ownership); err != nil {
		return nil, errors.Wrap(err, 0)
	}

	ticks, err := uc.Repo.GetTicksByContest(ctx, nil, contestID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	return ticks, nil
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

	if time.Now().Before(compClass.TimeBegin) {
		return domain.Tick{}, errors.New(domain.ErrContestNotStarted)
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
		ID:            0,
		Ownership:     contender.Ownership,
		Timestamp:     time.Now(),
		ContestID:     contest.ID,
		ProblemID:     problem.ID,
		Top:           tick.Top,
		AttemptsTop:   tick.AttemptsTop,
		Zone1:         tick.Zone1,
		AttemptsZone1: tick.AttemptsZone1,
		Zone2:         tick.Zone2,
		AttemptsZone2: tick.AttemptsZone2,
	}

	if err := (validators.TickValidator{}).Validate(newTick); err != nil {
		return domain.Tick{}, errors.Wrap(err, 0)
	}

	tick, err = uc.Repo.StoreTick(ctx, nil, newTick)
	if err != nil {
		return domain.Tick{}, errors.Wrap(err, 0)
	}

	uc.EventBroker.Dispatch(contest.ID, domain.AscentRegisteredEvent{
		TickID:        tick.ID,
		Timestamp:     tick.Timestamp,
		ContenderID:   contender.ID,
		ProblemID:     problem.ID,
		Top:           tick.Top,
		AttemptsTop:   tick.AttemptsTop,
		Zone1:         tick.Zone1,
		AttemptsZone1: tick.AttemptsZone1,
		Zone2:         tick.Zone2,
		AttemptsZone2: tick.AttemptsZone2,
	})

	return tick, nil
}

func (uc *TickUseCase) UpdateTick(ctx context.Context, tickID domain.TickID, patch domain.TickPatch) (domain.Tick, error) {
	existingTick, err := uc.Repo.GetTick(ctx, nil, tickID)
	if err != nil {
		return domain.Tick{}, errors.Wrap(err, 0)
	}

	role, err := uc.Authorizer.HasOwnership(ctx, existingTick.Ownership)
	if err != nil {
		return domain.Tick{}, errors.Wrap(err, 0)
	}

	contenderID := *existingTick.Ownership.ContenderID

	contender, err := uc.Repo.GetContender(ctx, nil, contenderID)
	if err != nil {
		return domain.Tick{}, errors.Wrap(err, 0)
	}

	contest, err := uc.Repo.GetContest(ctx, nil, existingTick.ContestID)
	if err != nil {
		return domain.Tick{}, errors.Errorf("%w: %w", domain.ErrRepositoryIntegrityViolation, err)
	}

	compClass, err := uc.Repo.GetCompClass(ctx, nil, contender.CompClassID)
	if err != nil {
		return domain.Tick{}, errors.Errorf("%w: %w", domain.ErrRepositoryIntegrityViolation, err)
	}

	if time.Now().Before(compClass.TimeBegin) {
		return domain.Tick{}, errors.New(domain.ErrContestNotStarted)
	}

	gracePeriodEnd := compClass.TimeEnd.Add(contest.GracePeriod)

	switch {
	case role.OneOf(domain.OrganizerRole, domain.AdminRole):
	case time.Now().After(gracePeriodEnd):
		return domain.Tick{}, errors.New(domain.ErrContestEnded)
	}

	updatedTick := existingTick

	changed := false

	if patch.Zone1.PresentAndDistinct(updatedTick.Zone1) {
		updatedTick.Zone1 = patch.Zone1.Value
		changed = true
	}

	if patch.AttemptsZone1.PresentAndDistinct(updatedTick.AttemptsZone1) {
		updatedTick.AttemptsZone1 = patch.AttemptsZone1.Value
		changed = true
	}

	if patch.Zone2.PresentAndDistinct(updatedTick.Zone2) {
		updatedTick.Zone2 = patch.Zone2.Value
		changed = true
	}

	if patch.AttemptsZone2.PresentAndDistinct(updatedTick.AttemptsZone2) {
		updatedTick.AttemptsZone2 = patch.AttemptsZone2.Value
		changed = true
	}

	if patch.Top.PresentAndDistinct(updatedTick.Top) {
		updatedTick.Top = patch.Top.Value
		changed = true
	}

	if patch.AttemptsTop.PresentAndDistinct(updatedTick.AttemptsTop) {
		updatedTick.AttemptsTop = patch.AttemptsTop.Value
		changed = true
	}

	if !changed {
		return existingTick, nil
	}

	updatedTick.Timestamp = time.Now()

	if err := (validators.TickValidator{}).Validate(updatedTick); err != nil {
		return domain.Tick{}, errors.Wrap(err, 0)
	}

	tx, err := uc.Repo.Begin()
	if err != nil {
		return domain.Tick{}, errors.Wrap(err, 0)
	}
	defer tx.Rollback()

	updatedTick, err = uc.Repo.StoreTick(ctx, tx, updatedTick)
	if err != nil {
		return domain.Tick{}, errors.Wrap(err, 0)
	}

	if err := tx.Commit(); err != nil {
		return domain.Tick{}, errors.Wrap(err, 0)
	}

	uc.EventBroker.Dispatch(contest.ID, domain.AscentDeregisteredEvent{
		TickID:      existingTick.ID,
		ContenderID: contender.ID,
		ProblemID:   existingTick.ProblemID,
	})

	uc.EventBroker.Dispatch(contest.ID, domain.AscentRegisteredEvent{
		TickID:        updatedTick.ID,
		Timestamp:     updatedTick.Timestamp,
		ContenderID:   contender.ID,
		ProblemID:     updatedTick.ProblemID,
		Top:           updatedTick.Top,
		AttemptsTop:   updatedTick.AttemptsTop,
		Zone1:         updatedTick.Zone1,
		AttemptsZone1: updatedTick.AttemptsZone1,
		Zone2:         updatedTick.Zone2,
		AttemptsZone2: updatedTick.AttemptsZone2,
	})

	return updatedTick, nil
}
