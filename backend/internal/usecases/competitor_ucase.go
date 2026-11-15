package usecases

import (
	"context"
	"strings"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

type competitorUseCaseRepository interface {
	domain.Transactor

	GetCompetitor(ctx context.Context, tx domain.Transaction, competitorID domain.CompetitorID) (domain.Competitor, error)
	GetCompetitorByCode(ctx context.Context, tx domain.Transaction, registrationCode string) (domain.Competitor, error)
	GetCompetitorsByCategory(ctx context.Context, tx domain.Transaction, categoryID domain.CategoryID) ([]domain.Competitor, error)
	GetCompetitorsByCompetition(ctx context.Context, tx domain.Transaction, competitionID domain.CompetitionID) ([]domain.Competitor, error)
	StoreCompetitor(ctx context.Context, tx domain.Transaction, competitor domain.Competitor) (domain.Competitor, error)
	DeleteCompetitor(ctx context.Context, tx domain.Transaction, competitorID domain.CompetitorID) error
	GetCompetition(ctx context.Context, tx domain.Transaction, competitionID domain.CompetitionID) (domain.Competition, error)
	GetCategory(ctx context.Context, tx domain.Transaction, categoryID domain.CategoryID) (domain.Category, error)
	GetNumberOfCompetitors(ctx context.Context, tx domain.Transaction, competitionID domain.CompetitionID) (int, error)
	GetScrubEligibleCompetitors(ctx context.Context, deadline time.Time) ([]domain.Competitor, error)
}

type pointValueKeeper interface {
	GetPointValues(competitorID domain.CompetitorID) []domain.PointValue
}

type CompetitorUseCase struct {
	Repo                      competitorUseCaseRepository
	Authorizer                domain.Authorizer
	EventBroker               domain.EventBroker
	ResultKeeper              domain.ResultKeeper
	PointValueKeeper          pointValueKeeper
	RegistrationCodeGenerator domain.CodeGenerator
}

func (uc *CompetitorUseCase) GetCompetitor(ctx context.Context, competitorID domain.CompetitorID) (domain.Competitor, error) {
	competitor, err := uc.Repo.GetCompetitor(ctx, nil, competitorID)
	if err != nil {
		return domain.Competitor{}, errors.Wrap(err, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, competitor.Ownership); err != nil {
		return domain.Competitor{}, errors.Wrap(err, 0)
	}

	return withResult(competitor, uc.ResultKeeper), nil
}

func (uc *CompetitorUseCase) GetCompetitorByCode(ctx context.Context, registrationCode string) (domain.Competitor, error) {
	competitor, err := uc.Repo.GetCompetitorByCode(ctx, nil, registrationCode)
	if err != nil {
		return domain.Competitor{}, errors.Wrap(err, 0)
	}

	return withResult(competitor, uc.ResultKeeper), nil
}

func (uc *CompetitorUseCase) GetPointValues(ctx context.Context, competitorID domain.CompetitorID) ([]domain.PointValue, error) {
	competitor, err := uc.Repo.GetCompetitor(ctx, nil, competitorID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, competitor.Ownership); err != nil {
		return nil, errors.Wrap(err, 0)
	}

	return uc.PointValueKeeper.GetPointValues(competitorID), nil
}

func (uc *CompetitorUseCase) GetCompetitorsByCategory(ctx context.Context, categoryID domain.CategoryID) ([]domain.Competitor, error) {
	category, err := uc.Repo.GetCategory(ctx, nil, categoryID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, category.Ownership); err != nil {
		return nil, errors.Wrap(err, 0)
	}

	competitors, err := uc.Repo.GetCompetitorsByCategory(ctx, nil, categoryID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	for i, competitor := range competitors {
		competitors[i] = withResult(competitor, uc.ResultKeeper)
	}

	return competitors, nil
}

func (uc *CompetitorUseCase) GetCompetitorsByCompetition(ctx context.Context, competitionID domain.CompetitionID) ([]domain.Competitor, error) {
	competition, err := uc.Repo.GetCompetition(ctx, nil, competitionID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, competition.Ownership); err != nil {
		return nil, errors.Wrap(err, 0)
	}

	competitors, err := uc.Repo.GetCompetitorsByCompetition(ctx, nil, competitionID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	for i, competitor := range competitors {
		competitors[i] = withResult(competitor, uc.ResultKeeper)
	}

	return competitors, nil
}

func (uc *CompetitorUseCase) PatchCompetitor(ctx context.Context, competitorID domain.CompetitorID, patch domain.CompetitorPatch) (domain.Competitor, error) {
	var mty domain.Competitor
	var events []any

	competitor, err := uc.Repo.GetCompetitor(ctx, nil, competitorID)
	if err != nil {
		return mty, errors.Wrap(err, 0)
	}

	role, err := uc.Authorizer.HasOwnership(ctx, competitor.Ownership)
	if err != nil {
		return mty, errors.Wrap(err, 0)
	}

	publicInfoEvent := domain.CompetitorPublicInfoUpdatedEvent{
		CompetitorID:        competitorID,
		CategoryID:          competitor.CategoryID,
		Name:                competitor.Name,
		WithdrawnFromFinals: competitor.WithdrawnFromFinals,
		Disqualified:        competitor.Disqualified,
		ScrubbedAt:          competitor.ScrubbedAt,
	}

	publicInfoEventBaseline := publicInfoEvent

	competition, err := uc.Repo.GetCompetition(ctx, nil, competitor.CompetitionID)
	if err != nil {
		return mty, errors.Errorf("%w: %w", domain.ErrRepositoryIntegrityViolation, err)
	}

	if competitor.CategoryID != 0 {
		category, err := uc.Repo.GetCategory(ctx, nil, competitor.CategoryID)
		if err != nil {
			return mty, errors.Errorf("%w: %w", domain.ErrRepositoryIntegrityViolation, err)
		}

		gracePeriodEnd := category.TimeEnd.Add(competition.GracePeriod)
		switch {
		case role.OneOf(domain.AdminRole, domain.OrganizerRole):
			break
		case time.Now().After(gracePeriodEnd):
			return mty, errors.Wrap(domain.ErrCompetitionEnded, 0)
		}
	}

	if patch.CategoryID.Present && competitor.CategoryID != patch.CategoryID.Value {
		if patch.CategoryID.Value == 0 {
			return mty, errors.Wrap(domain.ErrNotAllowed, 0)
		}

		category, err := uc.Repo.GetCategory(ctx, nil, patch.CategoryID.Value)
		if err != nil {
			return mty, errors.Wrap(err, 0)
		}

		if competitor.CategoryID == 0 {
			events = append(events, domain.CompetitorEnteredEvent{
				CompetitorID: competitorID,
				CategoryID:   patch.CategoryID.Value,
			})
		} else {
			events = append(events, domain.CompetitorSwitchedCategoryEvent{
				CompetitorID: competitorID,
				CategoryID:   patch.CategoryID.Value,
			})
		}

		gracePeriodEnd := category.TimeEnd.Add(competition.GracePeriod)

		switch {
		case role.OneOf(domain.AdminRole, domain.OrganizerRole):
			break
		case time.Now().After(gracePeriodEnd):
			return mty, errors.Wrap(domain.ErrCompetitionEnded, 0)
		}

		competitor.CategoryID = patch.CategoryID.Value

		if competitor.Entered.IsZero() {
			competitor.Entered = time.Now()
			competitor.ScrubBefore = category.TimeEnd.Add(competition.NameRetentionTime)
		}
	}

	if competitor.CategoryID == 0 {
		return mty, errors.New(domain.ErrNotRegistered)
	}

	if patch.Name.Present {
		if !competitor.ScrubbedAt.IsZero() {
			return mty, errors.Wrap(domain.ErrNotAllowed, 0)
		}

		competitor.Name = strings.TrimSpace(patch.Name.Value)

		if competitor.Name == "" {
			return mty, errors.Errorf("%w: %w", domain.ErrInvalidData, domain.ErrEmptyName)
		}
	}

	if patch.WithdrawnFromFinals.Present && competitor.WithdrawnFromFinals != patch.WithdrawnFromFinals.Value {
		if patch.WithdrawnFromFinals.Value {
			events = append(events, domain.CompetitorWithdrewFromFinalsEvent{
				CompetitorID: competitorID,
			})
		} else {
			events = append(events, domain.CompetitorReenteredFinalsEvent{
				CompetitorID: competitorID,
			})
		}

		competitor.WithdrawnFromFinals = patch.WithdrawnFromFinals.Value
	}

	if patch.Disqualified.Present && competitor.Disqualified != patch.Disqualified.Value {
		if !role.OneOf(domain.AdminRole, domain.OrganizerRole) {
			return mty, errors.Wrap(domain.ErrInsufficientRole, 0)
		}

		if patch.Disqualified.Value {
			events = append(events, domain.CompetitorDisqualifiedEvent{
				CompetitorID: competitorID,
			})
		} else {
			events = append(events, domain.CompetitorRequalifiedEvent{
				CompetitorID: competitorID,
			})
		}

		competitor.Disqualified = patch.Disqualified.Value
	}

	publicInfoEvent.CategoryID = competitor.CategoryID
	publicInfoEvent.Name = competitor.Name
	publicInfoEvent.WithdrawnFromFinals = competitor.WithdrawnFromFinals
	publicInfoEvent.Disqualified = competitor.Disqualified

	if publicInfoEvent != publicInfoEventBaseline {
		events = append(events, publicInfoEvent)
	}

	if competitor, err = uc.Repo.StoreCompetitor(ctx, nil, competitor); err != nil {
		return mty, errors.Wrap(err, 0)
	}

	for _, event := range events {
		uc.EventBroker.Dispatch(competition.ID, event)
	}

	return withResult(competitor, uc.ResultKeeper), nil
}

func (uc *CompetitorUseCase) ScrubCompetitor(ctx context.Context, competitorID domain.CompetitorID) (domain.Competitor, error) {
	var mty domain.Competitor

	competitor, err := uc.Repo.GetCompetitor(ctx, nil, competitorID)
	if err != nil {
		return mty, errors.Wrap(err, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, competitor.Ownership); err != nil {
		return mty, errors.Wrap(err, 0)
	}

	competitor.Name = ""
	competitor.ScrubbedAt = time.Now()
	competitor.WithdrawnFromFinals = true

	competitor, err = uc.Repo.StoreCompetitor(ctx, nil, competitor)
	if err != nil {
		return mty, errors.Wrap(err, 0)
	}

	uc.EventBroker.Dispatch(competitor.CompetitionID, domain.CompetitorPublicInfoUpdatedEvent{
		CompetitorID:        competitor.ID,
		CategoryID:          competitor.CategoryID,
		Name:                competitor.Name,
		WithdrawnFromFinals: competitor.WithdrawnFromFinals,
		Disqualified:        competitor.Disqualified,
		ScrubbedAt:          competitor.ScrubbedAt,
	})

	uc.EventBroker.Dispatch(competitor.CompetitionID, domain.CompetitorWithdrewFromFinalsEvent{
		CompetitorID: competitorID,
	})

	return withResult(competitor, uc.ResultKeeper), nil
}

func (uc *CompetitorUseCase) DeleteCompetitor(ctx context.Context, competitorID domain.CompetitorID) error {
	competitor, err := uc.Repo.GetCompetitor(ctx, nil, competitorID)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	role, err := uc.Authorizer.HasOwnership(ctx, competitor.Ownership)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	if !role.OneOf(domain.AdminRole, domain.OrganizerRole) {
		return errors.Wrap(domain.ErrInsufficientRole, 0)
	}

	if err := uc.Repo.DeleteCompetitor(ctx, nil, competitorID); err != nil {
		return errors.Wrap(err, 0)
	}

	return nil
}

const registrationCodeLength = 8

func (uc *CompetitorUseCase) CreateCompetitors(ctx context.Context, competitionID domain.CompetitionID, number int) ([]domain.Competitor, error) {
	competition, err := uc.Repo.GetCompetition(ctx, nil, competitionID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, competition.Ownership); err != nil {
		return nil, errors.Wrap(err, 0)
	}

	numberOfCompetitors, err := uc.Repo.GetNumberOfCompetitors(ctx, nil, competitionID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	if numberOfCompetitors+number > 500 {
		return nil, errors.New(domain.ErrLimitExceeded)
	}

	competitors := make([]domain.Competitor, 0)

	tx, err := uc.Repo.Begin()
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	for range number {
		competitor := domain.Competitor{
			ID:                  0,
			CompetitionID:       competitionID,
			Ownership:           competition.Ownership,
			CategoryID:          0,
			RegistrationCode:    uc.RegistrationCodeGenerator.Generate(registrationCodeLength),
			Name:                "",
			Entered:             time.Time{},
			WithdrawnFromFinals: false,
			Disqualified:        false,
			ScrubbedAt:          time.Time{},
			ScrubBefore:         time.Time{},
			Result:              nil,
		}

		if competitor, err = uc.Repo.StoreCompetitor(ctx, tx, competitor); err != nil {
			tx.Rollback()
			return nil, errors.Wrap(err, 0)
		}

		competitors = append(competitors, competitor)
	}

	err = tx.Commit()
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	return competitors, err
}

func (uc *CompetitorUseCase) ScrubCompetitors(ctx context.Context, deadline time.Time) (int, error) {
	competitors, err := uc.Repo.GetScrubEligibleCompetitors(ctx, deadline)
	if err != nil {
		return 0, errors.Wrap(err, 0)
	}

	if len(competitors) == 0 {
		return 0, nil
	}

	tx, err := uc.Repo.Begin()
	if err != nil {
		return 0, errors.Wrap(err, 0)
	}
	defer tx.Rollback()

	for i := range competitors {
		competitors[i].Name = ""
		competitors[i].ScrubbedAt = time.Now()

		if _, err := uc.Repo.StoreCompetitor(ctx, tx, competitors[i]); err != nil {
			return 0, errors.Wrap(err, 0)
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, errors.Wrap(err, 0)
	}

	for _, competitor := range competitors {
		uc.EventBroker.Dispatch(competitor.CompetitionID, domain.CompetitorPublicInfoUpdatedEvent{
			CompetitorID:        competitor.ID,
			CategoryID:          competitor.CategoryID,
			Name:                competitor.Name,
			WithdrawnFromFinals: competitor.WithdrawnFromFinals,
			Disqualified:        competitor.Disqualified,
			ScrubbedAt:          competitor.ScrubbedAt,
		})
	}

	return len(competitors), nil
}
