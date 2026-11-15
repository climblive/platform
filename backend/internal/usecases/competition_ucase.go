package usecases

import (
	"context"
	"strings"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/usecases/validators"
	"github.com/go-errors/errors"
	"github.com/microcosm-cc/bluemonday"
)

const maxCompetitionsPerWeek = 10

type competitionUseCaseRepository interface {
	domain.Transactor

	GetCompetition(ctx context.Context, tx domain.Transaction, competitionID domain.CompetitionID) (domain.Competition, error)
	GetAllCompetitions(ctx context.Context, tx domain.Transaction) ([]domain.Competition, error)
	GetCompetitorsByCompetition(ctx context.Context, tx domain.Transaction, competitionID domain.CompetitionID) ([]domain.Competitor, error)
	GetCompetitionsByOrganizer(ctx context.Context, tx domain.Transaction, organizerID domain.OrganizerID) ([]domain.Competition, error)
	GetOrganizer(ctx context.Context, tx domain.Transaction, organizerID domain.OrganizerID) (domain.Organizer, error)
	StoreCompetition(ctx context.Context, tx domain.Transaction, competition domain.Competition) (domain.Competition, error)
	GetCategoriesByCompetition(ctx context.Context, tx domain.Transaction, competitionID domain.CompetitionID) ([]domain.Category, error)
	StoreCategory(ctx context.Context, tx domain.Transaction, category domain.Category) (domain.Category, error)
	GetProblemsByCompetition(ctx context.Context, tx domain.Transaction, competitionID domain.CompetitionID) ([]domain.Problem, error)
	StoreProblem(ctx context.Context, tx domain.Transaction, problem domain.Problem) (domain.Problem, error)
	DeleteCompetition(ctx context.Context, tx domain.Transaction, competitionID domain.CompetitionID) error
	DeleteProblem(ctx context.Context, tx domain.Transaction, problemID domain.ProblemID) error
	DeleteCategory(ctx context.Context, tx domain.Transaction, categoryID domain.CategoryID) error
	DeleteCompetitor(ctx context.Context, tx domain.Transaction, competitorID domain.CompetitorID) error
	StoreCompetitor(ctx context.Context, tx domain.Transaction, competitor domain.Competitor) (domain.Competitor, error)
	GetRafflesByCompetition(ctx context.Context, tx domain.Transaction, competitionID domain.CompetitionID) ([]domain.Raffle, error)
	DeleteRaffle(ctx context.Context, tx domain.Transaction, raffleID domain.RaffleID) error
	StoreRaffle(ctx context.Context, tx domain.Transaction, raffle domain.Raffle) (domain.Raffle, error)
	GetRaffleWinners(ctx context.Context, tx domain.Transaction, raffleID domain.RaffleID) ([]domain.RaffleWinner, error)
	DeleteRaffleWinner(ctx context.Context, tx domain.Transaction, raffleWinnerID domain.RaffleWinnerID) error
	StoreRaffleWinner(ctx context.Context, tx domain.Transaction, winner domain.RaffleWinner) (domain.RaffleWinner, error)
	GetTicksByCompetition(ctx context.Context, tx domain.Transaction, competitionID domain.CompetitionID) ([]domain.Tick, error)
	DeleteTick(ctx context.Context, tx domain.Transaction, tickID domain.TickID) error
	StoreTick(ctx context.Context, tx domain.Transaction, tick domain.Tick) (domain.Tick, error)
	StoreResult(ctx context.Context, tx domain.Transaction, result domain.Result) error
}

type CompetitionUseCase struct {
	Authorizer         domain.Authorizer
	Repo               competitionUseCaseRepository
	ResultKeeper       domain.ResultKeeper
	ScoreEngineManager scoreEngineManager
	EventBroker        domain.EventBroker
}

var sanitizationPolicy = bluemonday.UGCPolicy()

func (uc *CompetitionUseCase) GetCompetition(ctx context.Context, competitionID domain.CompetitionID) (domain.Competition, error) {
	competition, err := uc.Repo.GetCompetition(ctx, nil, competitionID)
	if err != nil {
		return domain.Competition{}, errors.Wrap(err, 0)
	}

	return competition, nil
}

func (uc *CompetitionUseCase) GetAllCompetitions(ctx context.Context) ([]domain.Competition, error) {
	var role domain.AuthRole
	var err error

	if role, err = uc.Authorizer.HasOwnership(ctx, domain.OwnershipData{}); err != nil {
		return nil, errors.Wrap(err, 0)
	}

	if role != domain.AdminRole {
		return nil, domain.ErrNotAuthorized
	}

	competitions, err := uc.Repo.GetAllCompetitions(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	return competitions, nil
}

func (uc *CompetitionUseCase) GetCompetitionsByOrganizer(ctx context.Context, organizerID domain.OrganizerID) ([]domain.Competition, error) {
	organizer, err := uc.Repo.GetOrganizer(ctx, nil, organizerID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, organizer.Ownership); err != nil {
		return nil, errors.Wrap(err, 0)
	}

	competitions, err := uc.Repo.GetCompetitionsByOrganizer(ctx, nil, organizerID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	return competitions, nil
}

func (uc *CompetitionUseCase) GetScoreboard(ctx context.Context, competitionID domain.CompetitionID) ([]domain.ScoreboardEntry, error) {
	competitors, err := uc.Repo.GetCompetitorsByCompetition(ctx, nil, competitionID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	entries := make([]domain.ScoreboardEntry, 0)

	for _, competitor := range competitors {
		if competitor.CategoryID == 0 {
			continue
		}

		entry := domain.ScoreboardEntry{
			CompetitorID:        competitor.ID,
			CategoryID:          competitor.CategoryID,
			Name:                competitor.Name,
			WithdrawnFromFinals: competitor.WithdrawnFromFinals,
			Disqualified:        competitor.Disqualified,
			ScrubbedAt:          competitor.ScrubbedAt,
			Result:              competitor.Result,
		}

		if result, err := uc.ResultKeeper.GetResult(competitor.ID); err == nil {
			entry.Result = &result
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

func (uc *CompetitionUseCase) PatchCompetition(ctx context.Context, competitionID domain.CompetitionID, patch domain.CompetitionPatch) (domain.Competition, error) {
	var mty domain.Competition

	competition, err := uc.Repo.GetCompetition(ctx, nil, competitionID)
	if err != nil {
		return mty, errors.Wrap(err, 0)
	}

	_, err = uc.Authorizer.HasOwnership(ctx, competition.Ownership)
	if err != nil {
		return mty, errors.Wrap(err, 0)
	}

	if !competition.ArchivedAt.IsZero() {
		return mty, errors.Wrap(domain.ErrArchived, 0)
	}

	rulesUpdateEventBaseline := domain.RulesUpdatedEvent{
		CompetitionID:      competition.ID,
		QualifyingProblems: competition.QualifyingProblems,
		Finalists:          competition.Finalists,
		UsePoints:          competition.UsePoints,
		PooledPoints:       competition.PooledPoints,
		MaxAttempts:        competition.MaxAttempts,
		PointDeduction:     competition.PointDeduction,
	}

	if patch.Location.Present {
		competition.Location = strings.TrimSpace(patch.Location.Value)
	}

	if patch.Country.Present {
		competition.Country = strings.TrimSpace(patch.Country.Value)
	}

	if patch.SeriesID.Present {
		competition.SeriesID = patch.SeriesID.Value
	}

	if patch.Name.Present {
		competition.Name = strings.TrimSpace(patch.Name.Value)
	}

	if patch.Description.Present {
		competition.Description = strings.TrimSpace(patch.Description.Value)
	}

	if patch.QualifyingProblems.Present {
		competition.QualifyingProblems = patch.QualifyingProblems.Value
	}

	if patch.Finalists.Present {
		competition.Finalists = patch.Finalists.Value
	}

	if patch.UsePoints.Present {
		competition.UsePoints = patch.UsePoints.Value
	}

	if patch.PooledPoints.Present {
		competition.PooledPoints = patch.PooledPoints.Value
	}

	if patch.MaxAttempts.Present {
		competition.MaxAttempts = patch.MaxAttempts.Value
	}

	if patch.PointDeduction.Present {
		competition.PointDeduction = patch.PointDeduction.Value
	}

	if patch.Info.Present {
		competition.Info = sanitizationPolicy.Sanitize(patch.Info.Value)
	}

	if patch.GracePeriod.Present {
		competition.GracePeriod = patch.GracePeriod.Value
	}

	if err := (validators.CompetitionValidator{}).Validate(competition); err != nil {
		return mty, errors.Wrap(err, 0)
	}

	if _, err = uc.Repo.StoreCompetition(ctx, nil, competition); err != nil {
		return mty, errors.Wrap(err, 0)
	}

	event := domain.RulesUpdatedEvent{
		CompetitionID:      competition.ID,
		QualifyingProblems: competition.QualifyingProblems,
		Finalists:          competition.Finalists,
		UsePoints:          competition.UsePoints,
		PooledPoints:       competition.PooledPoints,
		MaxAttempts:        competition.MaxAttempts,
		PointDeduction:     competition.PointDeduction,
	}

	if event != rulesUpdateEventBaseline {
		uc.EventBroker.Dispatch(competitionID, event)
	}

	return competition, nil
}

func (uc *CompetitionUseCase) ArchiveCompetition(ctx context.Context, competitionID domain.CompetitionID) (domain.Competition, error) {
	competition, err := uc.Repo.GetCompetition(ctx, nil, competitionID)
	if err != nil {
		return domain.Competition{}, errors.Wrap(err, 0)
	}

	if _, err = uc.Authorizer.HasOwnership(ctx, competition.Ownership); err != nil {
		return domain.Competition{}, errors.Wrap(err, 0)
	}

	competition.ArchivedAt = time.Now()

	engines, err := uc.ScoreEngineManager.ListScoreEnginesByCompetition(ctx, competitionID)
	if err != nil {
		return domain.Competition{}, errors.Wrap(err, 0)
	}

	for _, engine := range engines {
		if err = uc.ScoreEngineManager.StopScoreEngine(ctx, engine.InstanceID); err != nil {
			return domain.Competition{}, errors.Wrap(err, 0)
		}
	}

	competition, err = uc.Repo.StoreCompetition(ctx, nil, competition)
	if err != nil {
		return domain.Competition{}, errors.Wrap(err, 0)
	}

	return competition, nil
}

func (uc *CompetitionUseCase) RestoreCompetition(ctx context.Context, competitionID domain.CompetitionID) (domain.Competition, error) {
	competition, err := uc.Repo.GetCompetition(ctx, nil, competitionID)
	if err != nil {
		return domain.Competition{}, errors.Wrap(err, 0)
	}

	if _, err = uc.Authorizer.HasOwnership(ctx, competition.Ownership); err != nil {
		return domain.Competition{}, errors.Wrap(err, 0)
	}

	competition.ArchivedAt = time.Time{}

	competition, err = uc.Repo.StoreCompetition(ctx, nil, competition)
	if err != nil {
		return domain.Competition{}, errors.Wrap(err, 0)
	}

	return competition, nil
}

func (uc *CompetitionUseCase) CreateCompetition(ctx context.Context, organizerID domain.OrganizerID, tmpl domain.CompetitionTemplate) (domain.Competition, error) {
	organizer, err := uc.Repo.GetOrganizer(ctx, nil, organizerID)
	if err != nil {
		return domain.Competition{}, errors.Wrap(err, 0)
	}

	role, err := uc.Authorizer.HasOwnership(ctx, organizer.Ownership)
	if err != nil {
		return domain.Competition{}, errors.Wrap(err, 0)
	}

	if !role.OneOf(domain.AdminRole) {
		competitions, err := uc.Repo.GetCompetitionsByOrganizer(ctx, nil, organizerID)
		if err != nil {
			return domain.Competition{}, errors.Wrap(err, 0)
		}

		oneWeekAgo := time.Now().Add(-7 * 24 * time.Hour)
		recentCount := 0
		for _, c := range competitions {
			if c.ArchivedAt.IsZero() && c.Created.After(oneWeekAgo) {
				recentCount++
			}
		}

		if recentCount >= maxCompetitionsPerWeek {
			return domain.Competition{}, errors.New(domain.ErrLimitExceeded)
		}
	}

	competition := domain.Competition{
		ID: domain.CompetitionID{},
		Ownership: domain.OwnershipData{
			OrganizerID:  organizerID,
			CompetitorID: nil,
		},
		ArchivedAt:            time.Time{},
		SeriesID:              0,
		TimeBegin:             time.Time{},
		TimeEnd:               time.Time{},
		RegisteredCompetitors: 0,
		Location:              strings.TrimSpace(tmpl.Location),
		Country:               strings.TrimSpace(tmpl.Country),
		Name:                  strings.TrimSpace(tmpl.Name),
		Description:           strings.TrimSpace(tmpl.Description),
		QualifyingProblems:    tmpl.QualifyingProblems,
		Finalists:             tmpl.Finalists,
		UsePoints:             tmpl.UsePoints,
		PooledPoints:          tmpl.PooledPoints,
		MaxAttempts:           tmpl.MaxAttempts,
		PointDeduction:        tmpl.PointDeduction,
		Info:                  sanitizationPolicy.Sanitize(tmpl.Info),
		GracePeriod:           tmpl.GracePeriod,
		NameRetentionTime:     tmpl.NameRetentionTime,
		Created:               time.Now(),
	}

	if err := (validators.CompetitionValidator{}).Validate(competition); err != nil {
		return domain.Competition{}, errors.Wrap(err, 0)
	}

	competition, err = uc.Repo.StoreCompetition(ctx, nil, competition)
	if err != nil {
		return domain.Competition{}, errors.Wrap(err, 0)
	}

	return competition, nil
}

func (uc *CompetitionUseCase) DuplicateCompetition(ctx context.Context, competitionID domain.CompetitionID) (domain.Competition, error) {
	competition, err := uc.Repo.GetCompetition(ctx, nil, competitionID)
	if err != nil {
		return domain.Competition{}, errors.Wrap(err, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, competition.Ownership); err != nil {
		return domain.Competition{}, errors.Wrap(err, 0)
	}

	if !competition.ArchivedAt.IsZero() {
		return domain.Competition{}, errors.Wrap(domain.ErrArchived, 0)
	}

	categories, err := uc.Repo.GetCategoriesByCompetition(ctx, nil, competitionID)
	if err != nil {
		return domain.Competition{}, errors.Wrap(err, 0)
	}

	problems, err := uc.Repo.GetProblemsByCompetition(ctx, nil, competitionID)
	if err != nil {
		return domain.Competition{}, errors.Wrap(err, 0)
	}

	duplicatedCompetition := competition
	duplicatedCompetition.ID = domain.CompetitionID{}
	duplicatedCompetition.Name += " (Copy)"

	tx, err := uc.Repo.Begin()
	if err != nil {
		return domain.Competition{}, errors.Wrap(err, 0)
	}

	duplicate := func() (domain.Competition, error) {
		createdCompetition, err := uc.Repo.StoreCompetition(ctx, tx, duplicatedCompetition)
		if err != nil {
			return domain.Competition{}, err
		}

		for _, category := range categories {
			category.ID = 0
			category.CompetitionID = createdCompetition.ID

			_, err = uc.Repo.StoreCategory(ctx, tx, category)
			if err != nil {
				return domain.Competition{}, err
			}
		}

		for _, problem := range problems {
			problem.ID = 0
			problem.CompetitionID = createdCompetition.ID

			_, err = uc.Repo.StoreProblem(ctx, tx, problem)
			if err != nil {
				return domain.Competition{}, err
			}
		}

		return createdCompetition, nil
	}

	createdCompetition, err := duplicate()
	if err != nil {
		tx.Rollback()
		return domain.Competition{}, errors.Wrap(err, 0)
	}

	err = tx.Commit()
	if err != nil {
		return domain.Competition{}, errors.Wrap(err, 0)
	}

	return createdCompetition, nil
}

func (uc *CompetitionUseCase) TransferCompetition(ctx context.Context, competitionID domain.CompetitionID, newOrganizerID domain.OrganizerID) (domain.Competition, error) {
	competition, err := uc.Repo.GetCompetition(ctx, nil, competitionID)
	if err != nil {
		return domain.Competition{}, errors.Wrap(err, 0)
	}

	if !competition.ArchivedAt.IsZero() {
		return domain.Competition{}, errors.Wrap(domain.ErrArchived, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, competition.Ownership); err != nil {
		return domain.Competition{}, errors.Wrap(err, 0)
	}

	newOrganizer, err := uc.Repo.GetOrganizer(ctx, nil, newOrganizerID)
	if err != nil {
		return domain.Competition{}, errors.Wrap(err, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, newOrganizer.Ownership); err != nil {
		return domain.Competition{}, errors.Wrap(err, 0)
	}

	categories, err := uc.Repo.GetCategoriesByCompetition(ctx, nil, competitionID)
	if err != nil {
		return domain.Competition{}, errors.Wrap(err, 0)
	}

	problems, err := uc.Repo.GetProblemsByCompetition(ctx, nil, competitionID)
	if err != nil {
		return domain.Competition{}, errors.Wrap(err, 0)
	}

	competitors, err := uc.Repo.GetCompetitorsByCompetition(ctx, nil, competitionID)
	if err != nil {
		return domain.Competition{}, errors.Wrap(err, 0)
	}

	raffles, err := uc.Repo.GetRafflesByCompetition(ctx, nil, competitionID)
	if err != nil {
		return domain.Competition{}, errors.Wrap(err, 0)
	}

	allRaffleWinners := make([]domain.RaffleWinner, 0)

	for _, raffle := range raffles {
		winners, err := uc.Repo.GetRaffleWinners(ctx, nil, raffle.ID)
		if err != nil {
			return domain.Competition{}, errors.Wrap(err, 0)
		}

		allRaffleWinners = append(allRaffleWinners, winners...)
	}

	ticks, err := uc.Repo.GetTicksByCompetition(ctx, nil, competitionID)
	if err != nil {
		return domain.Competition{}, errors.Wrap(err, 0)
	}

	tx, err := uc.Repo.Begin()
	if err != nil {
		return domain.Competition{}, errors.Wrap(err, 0)
	}

	purge := func() error {
		for _, tick := range ticks {
			err = uc.Repo.DeleteTick(ctx, tx, tick.ID)
			if err != nil {
				return err
			}
		}

		for _, winner := range allRaffleWinners {
			err = uc.Repo.DeleteRaffleWinner(ctx, tx, winner.ID)
			if err != nil {
				return err
			}
		}

		for _, raffle := range raffles {
			err = uc.Repo.DeleteRaffle(ctx, tx, raffle.ID)
			if err != nil {
				return err
			}
		}

		for _, competitor := range competitors {
			err = uc.Repo.DeleteCompetitor(ctx, tx, competitor.ID)
			if err != nil {
				return err
			}
		}

		for _, problem := range problems {
			err = uc.Repo.DeleteProblem(ctx, tx, problem.ID)
			if err != nil {
				return err
			}
		}

		for _, category := range categories {
			err = uc.Repo.DeleteCategory(ctx, tx, category.ID)
			if err != nil {
				return err
			}
		}

		err := uc.Repo.DeleteCompetition(ctx, tx, competitionID)
		if err != nil {
			return err
		}

		return nil
	}

	recreate := func() error {
		_, err := uc.Repo.StoreCompetition(ctx, tx, competition)
		if err != nil {
			return err
		}

		for _, category := range categories {
			_, err = uc.Repo.StoreCategory(ctx, tx, category)
			if err != nil {
				return err
			}
		}

		for _, problem := range problems {
			_, err = uc.Repo.StoreProblem(ctx, tx, problem)
			if err != nil {
				return err
			}
		}

		for _, competitor := range competitors {
			_, err = uc.Repo.StoreCompetitor(ctx, tx, competitor)
			if err != nil {
				return err
			}

			result := competitor.Result
			if result == nil {
				continue
			}

			err := uc.Repo.StoreResult(ctx, tx, *result)
			if err != nil {
				return err
			}
		}

		for _, raffle := range raffles {
			_, err = uc.Repo.StoreRaffle(ctx, tx, raffle)
			if err != nil {
				return err
			}
		}

		for _, winner := range allRaffleWinners {
			_, err = uc.Repo.StoreRaffleWinner(ctx, tx, winner)
			if err != nil {
				return err
			}
		}

		for _, tick := range ticks {
			tick.Revision++
			_, err = uc.Repo.StoreTick(ctx, tx, tick)
			if err != nil {
				return err
			}
		}

		return nil
	}

	err = purge()
	if err != nil {
		tx.Rollback()
		return domain.Competition{}, errors.Wrap(err, 0)
	}

	competition.Ownership.OrganizerID = newOrganizerID

	for index := range categories {
		categories[index].Ownership.OrganizerID = newOrganizerID
	}

	for index := range problems {
		problems[index].Ownership.OrganizerID = newOrganizerID
	}

	for index := range competitors {
		competitors[index].Ownership.OrganizerID = newOrganizerID
	}

	for index := range raffles {
		raffles[index].Ownership.OrganizerID = newOrganizerID
	}

	for index := range allRaffleWinners {
		allRaffleWinners[index].Ownership.OrganizerID = newOrganizerID
	}

	for index := range ticks {
		ticks[index].Ownership.OrganizerID = newOrganizerID
	}

	err = recreate()
	if err != nil {
		tx.Rollback()
		return domain.Competition{}, errors.Wrap(err, 0)
	}

	err = tx.Commit()
	if err != nil {
		return domain.Competition{}, errors.Wrap(err, 0)
	}

	return competition, nil
}
