package usecases

import (
	"context"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/usecases/validators"
	"github.com/go-errors/errors"
)

const maxTickRevisionJump = 1_000

type tickUseCaseRepository interface {
	domain.Transactor

	GetCompetitor(ctx context.Context, tx domain.Transaction, competitorID domain.CompetitorID) (domain.Competitor, error)
	GetTicksByCompetitor(ctx context.Context, tx domain.Transaction, competitorID domain.CompetitorID) ([]domain.Tick, error)
	GetTicksByCompetition(ctx context.Context, tx domain.Transaction, competitionID domain.CompetitionID) ([]domain.Tick, error)
	GetCompetition(ctx context.Context, tx domain.Transaction, competitionID domain.CompetitionID) (domain.Competition, error)
	GetCategory(ctx context.Context, tx domain.Transaction, categoryID domain.CategoryID) (domain.Category, error)
	GetProblem(ctx context.Context, tx domain.Transaction, problemID domain.ProblemID) (domain.Problem, error)
	DeleteTick(ctx context.Context, tx domain.Transaction, tickID domain.TickID) error
	StoreTick(ctx context.Context, tx domain.Transaction, tick domain.Tick) (domain.Tick, error)
	GetTick(ctx context.Context, tx domain.Transaction, tickID domain.TickID) (domain.Tick, error)
	GetTickByCompetitorAndProblem(ctx context.Context, tx domain.Transaction, competitorID domain.CompetitorID, problemID domain.ProblemID) (domain.Tick, error)
}

type TickUseCase struct {
	Repo        tickUseCaseRepository
	Authorizer  domain.Authorizer
	EventBroker domain.EventBroker
}

func (uc *TickUseCase) GetTicksByCompetitor(ctx context.Context, competitorID domain.CompetitorID) ([]domain.Tick, error) {
	competitor, err := uc.Repo.GetCompetitor(ctx, nil, competitorID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, competitor.Ownership); err != nil {
		return nil, errors.Wrap(err, 0)
	}

	ticks, err := uc.Repo.GetTicksByCompetitor(ctx, nil, competitorID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	return ticks, nil
}

func (uc *TickUseCase) GetTicksByCompetition(ctx context.Context, competitionID domain.CompetitionID) ([]domain.Tick, error) {
	competition, err := uc.Repo.GetCompetition(ctx, nil, competitionID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, competition.Ownership); err != nil {
		return nil, errors.Wrap(err, 0)
	}

	ticks, err := uc.Repo.GetTicksByCompetition(ctx, nil, competitionID)
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

	competitorID := *tick.Ownership.CompetitorID

	competitor, err := uc.Repo.GetCompetitor(ctx, nil, competitorID)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	competition, err := uc.Repo.GetCompetition(ctx, nil, tick.CompetitionID)
	if err != nil {
		return errors.Errorf("%w: %w", domain.ErrRepositoryIntegrityViolation, err)
	}

	category, err := uc.Repo.GetCategory(ctx, nil, competitor.CategoryID)
	if err != nil {
		return errors.Errorf("%w: %w", domain.ErrRepositoryIntegrityViolation, err)
	}

	gracePeriodEnd := category.TimeEnd.Add(competition.GracePeriod)

	switch {
	case role.OneOf(domain.OrganizerRole, domain.AdminRole):
	case time.Now().After(gracePeriodEnd):
		return errors.New(domain.ErrCompetitionEnded)
	}

	err = uc.Repo.DeleteTick(ctx, nil, tickID)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	uc.EventBroker.Dispatch(competition.ID, domain.AscentDeregisteredEvent{
		TickID:       tickID,
		CompetitorID: competitor.ID,
		ProblemID:    tick.ProblemID,
	})

	return nil
}

func (uc *TickUseCase) PutTick(ctx context.Context, competitorID domain.CompetitorID, tick domain.Tick) (domain.Tick, error) {
	competitor, err := uc.Repo.GetCompetitor(ctx, nil, competitorID)
	if err != nil {
		return domain.Tick{}, errors.Wrap(err, 0)
	}

	role, err := uc.Authorizer.HasOwnership(ctx, competitor.Ownership)
	if err != nil {
		return domain.Tick{}, errors.Wrap(err, 0)
	}

	competition, err := uc.Repo.GetCompetition(ctx, nil, competitor.CompetitionID)
	if err != nil {
		return domain.Tick{}, errors.Errorf("%w: %w", domain.ErrRepositoryIntegrityViolation, err)
	}

	category, err := uc.Repo.GetCategory(ctx, nil, competitor.CategoryID)
	if err != nil {
		return domain.Tick{}, errors.Errorf("%w: %w", domain.ErrRepositoryIntegrityViolation, err)
	}

	if time.Now().Before(category.TimeBegin) {
		return domain.Tick{}, errors.New(domain.ErrCompetitionNotStarted)
	}

	gracePeriodEnd := category.TimeEnd.Add(competition.GracePeriod)

	switch {
	case role.OneOf(domain.OrganizerRole, domain.AdminRole):
	case time.Now().After(gracePeriodEnd):
		return domain.Tick{}, errors.New(domain.ErrCompetitionEnded)
	}

	problem, err := uc.Repo.GetProblem(ctx, nil, tick.ProblemID)
	if err != nil {
		return domain.Tick{}, errors.Wrap(err, 0)
	}

	if problem.CompetitionID != competition.ID {
		return domain.Tick{}, errors.New(domain.ErrProblemNotInCompetition)
	}

	existingTick, err := uc.Repo.GetTickByCompetitorAndProblem(ctx, nil, competitorID, problem.ID)
	switch {
	case err == nil:
	case errors.Is(err, domain.ErrNotFound):
	default:
		return domain.Tick{}, errors.Wrap(err, 0)
	}

	if tick.Revision <= existingTick.Revision {
		return domain.Tick{}, errors.Wrap(domain.ErrSuperseded, 0)
	}

	if tick.Revision-existingTick.Revision > maxTickRevisionJump {
		return domain.Tick{}, errors.Wrap(domain.ErrInvalidData, 0)
	}

	existingTick.Ownership = competitor.Ownership
	existingTick.CompetitionID = competition.ID
	existingTick.ProblemID = problem.ID
	existingTick.Timestamp = time.Now()
	existingTick.Revision = tick.Revision

	existingTick.Top = tick.Top
	existingTick.AttemptsTop = tick.AttemptsTop
	existingTick.Zone1 = tick.Zone1
	existingTick.AttemptsZone1 = tick.AttemptsZone1
	existingTick.Zone2 = tick.Zone2
	existingTick.AttemptsZone2 = tick.AttemptsZone2

	if err := (validators.TickValidator{}).Validate(existingTick); err != nil {
		return domain.Tick{}, errors.Wrap(err, 0)
	}

	existingTick, err = uc.Repo.StoreTick(ctx, nil, existingTick)
	if err != nil {
		return domain.Tick{}, errors.Wrap(err, 0)
	}

	uc.EventBroker.Dispatch(competition.ID, domain.AscentRegisteredEvent{
		TickID:        existingTick.ID,
		Timestamp:     existingTick.Timestamp,
		CompetitorID:  *existingTick.Ownership.CompetitorID,
		ProblemID:     existingTick.ProblemID,
		Top:           existingTick.Top,
		Revision:      existingTick.Revision,
		AttemptsTop:   existingTick.AttemptsTop,
		Zone1:         existingTick.Zone1,
		AttemptsZone1: existingTick.AttemptsZone1,
		Zone2:         existingTick.Zone2,
		AttemptsZone2: existingTick.AttemptsZone2,
	})

	return existingTick, nil
}
