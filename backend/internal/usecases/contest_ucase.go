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

type contestUseCaseRepository interface {
	domain.Transactor

	GetContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) (domain.Contest, error)
	GetAllContests(ctx context.Context, tx domain.Transaction) ([]domain.Contest, error)
	GetContendersByContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) ([]domain.Contender, error)
	GetContestsByOrganizer(ctx context.Context, tx domain.Transaction, organizerID domain.OrganizerID) ([]domain.Contest, error)
	GetOrganizer(ctx context.Context, tx domain.Transaction, organizerID domain.OrganizerID) (domain.Organizer, error)
	StoreContest(ctx context.Context, tx domain.Transaction, contest domain.Contest) (domain.Contest, error)
	GetCompClassesByContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) ([]domain.CompClass, error)
	StoreCompClass(ctx context.Context, tx domain.Transaction, compClass domain.CompClass) (domain.CompClass, error)
	GetProblemsByContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) ([]domain.Problem, error)
	StoreProblem(ctx context.Context, tx domain.Transaction, problem domain.Problem) (domain.Problem, error)
	DeleteContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) error
	DeleteProblem(ctx context.Context, tx domain.Transaction, problemID domain.ProblemID) error
	DeleteCompClass(ctx context.Context, tx domain.Transaction, compClassID domain.CompClassID) error
	DeleteContender(ctx context.Context, tx domain.Transaction, contenderID domain.ContenderID) error
	StoreContender(ctx context.Context, tx domain.Transaction, contender domain.Contender) (domain.Contender, error)
	GetRafflesByContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) ([]domain.Raffle, error)
	DeleteRaffle(ctx context.Context, tx domain.Transaction, raffleID domain.RaffleID) error
	StoreRaffle(ctx context.Context, tx domain.Transaction, raffle domain.Raffle) (domain.Raffle, error)
	GetRaffleWinners(ctx context.Context, tx domain.Transaction, raffleID domain.RaffleID) ([]domain.RaffleWinner, error)
	DeleteRaffleWinner(ctx context.Context, tx domain.Transaction, raffleWinnerID domain.RaffleWinnerID) error
	StoreRaffleWinner(ctx context.Context, tx domain.Transaction, winner domain.RaffleWinner) (domain.RaffleWinner, error)
	GetTicksByContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) ([]domain.Tick, error)
	DeleteTick(ctx context.Context, tx domain.Transaction, tickID domain.TickID) error
	StoreTick(ctx context.Context, tx domain.Transaction, tick domain.Tick) (domain.Tick, error)
	StoreScore(ctx context.Context, tx domain.Transaction, score domain.Score) error
}

type ContestUseCase struct {
	Authorizer         domain.Authorizer
	Repo               contestUseCaseRepository
	ScoreKeeper        domain.ScoreKeeper
	ScoreEngineManager scoreEngineManager
	EventBroker        domain.EventBroker
}

var sanitizationPolicy = bluemonday.UGCPolicy()

func (uc *ContestUseCase) GetContest(ctx context.Context, contestID domain.ContestID) (domain.Contest, error) {
	contest, err := uc.Repo.GetContest(ctx, nil, contestID)
	if err != nil {
		return domain.Contest{}, errors.Wrap(err, 0)
	}

	return contest, nil
}

func (uc *ContestUseCase) GetAllContests(ctx context.Context) ([]domain.Contest, error) {
	var role domain.AuthRole
	var err error

	if role, err = uc.Authorizer.HasOwnership(ctx, domain.OwnershipData{}); err != nil {
		return nil, errors.Wrap(err, 0)
	}

	if role != domain.AdminRole {
		return nil, domain.ErrNotAuthorized
	}

	contests, err := uc.Repo.GetAllContests(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	return contests, nil
}

func (uc *ContestUseCase) GetContestsByOrganizer(ctx context.Context, organizerID domain.OrganizerID) ([]domain.Contest, error) {
	organizer, err := uc.Repo.GetOrganizer(ctx, nil, organizerID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, organizer.Ownership); err != nil {
		return nil, errors.Wrap(err, 0)
	}

	contests, err := uc.Repo.GetContestsByOrganizer(ctx, nil, organizerID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	return contests, nil
}

func (uc *ContestUseCase) GetScoreboard(ctx context.Context, contestID domain.ContestID) ([]domain.ScoreboardEntry, error) {
	contenders, err := uc.Repo.GetContendersByContest(ctx, nil, contestID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	entries := make([]domain.ScoreboardEntry, 0)

	for _, contender := range contenders {
		if contender.CompClassID == 0 {
			continue
		}

		entry := domain.ScoreboardEntry{
			ContenderID:         contender.ID,
			CompClassID:         contender.CompClassID,
			Name:                contender.Name,
			WithdrawnFromFinals: contender.WithdrawnFromFinals,
			Disqualified:        contender.Disqualified,
			Score:               contender.Score,
		}

		if score, err := uc.ScoreKeeper.GetScore(contender.ID); err == nil {
			entry.Score = &score
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

func (uc *ContestUseCase) PatchContest(ctx context.Context, contestID domain.ContestID, patch domain.ContestPatch) (domain.Contest, error) {
	var mty domain.Contest

	contest, err := uc.Repo.GetContest(ctx, nil, contestID)
	if err != nil {
		return mty, errors.Wrap(err, 0)
	}

	_, err = uc.Authorizer.HasOwnership(ctx, contest.Ownership)
	if err != nil {
		return mty, errors.Wrap(err, 0)
	}

	rulesUpdateEventBaseline := domain.RulesUpdatedEvent{
		QualifyingProblems: contest.QualifyingProblems,
		Finalists:          contest.Finalists,
	}

	if patch.Archived.Present {
		contest.Archived = patch.Archived.Value

		if contest.Archived {
			engines, err := uc.ScoreEngineManager.ListScoreEnginesByContest(ctx, contestID)
			if err != nil {
				return mty, errors.Wrap(err, 0)
			}

			for _, engine := range engines {
				err = uc.ScoreEngineManager.StopScoreEngine(ctx, engine.InstanceID)
				if err != nil {
					return mty, errors.Wrap(err, 0)
				}
			}
		}
	}

	patchAnythingOtherThanArchive := patch != (domain.ContestPatch{}) && patch != domain.ContestPatch{
		Archived: patch.Archived,
	}

	if contest.Archived && patchAnythingOtherThanArchive {
		return mty, errors.Wrap(domain.ErrArchived, 0)
	}

	if patch.Location.Present {
		contest.Location = strings.TrimSpace(patch.Location.Value)
	}

	if patch.SeriesID.Present {
		contest.SeriesID = patch.SeriesID.Value
	}

	if patch.Name.Present {
		contest.Name = strings.TrimSpace(patch.Name.Value)
	}

	if patch.Description.Present {
		contest.Description = strings.TrimSpace(patch.Description.Value)
	}

	if patch.QualifyingProblems.Present {
		contest.QualifyingProblems = patch.QualifyingProblems.Value
	}

	if patch.Finalists.Present {
		contest.Finalists = patch.Finalists.Value
	}

	if patch.Info.Present {
		contest.Info = sanitizationPolicy.Sanitize(patch.Info.Value)
	}

	if patch.GracePeriod.Present {
		contest.GracePeriod = patch.GracePeriod.Value
	}

	if patch.EvaluationMode.Present {
		contest.EvaluationMode = patch.EvaluationMode.Value
	}

	if err := (validators.ContestValidator{}).Validate(contest); err != nil {
		return mty, errors.Wrap(err, 0)
	}

	if _, err = uc.Repo.StoreContest(ctx, nil, contest); err != nil {
		return mty, errors.Wrap(err, 0)
	}

	event := domain.RulesUpdatedEvent{
		QualifyingProblems: contest.QualifyingProblems,
		Finalists:          contest.Finalists,
	}

	if event != rulesUpdateEventBaseline {
		uc.EventBroker.Dispatch(contestID, event)
	}

	return contest, nil
}

func (uc *ContestUseCase) CreateContest(ctx context.Context, organizerID domain.OrganizerID, tmpl domain.ContestTemplate) (domain.Contest, error) {
	organizer, err := uc.Repo.GetOrganizer(ctx, nil, organizerID)
	if err != nil {
		return domain.Contest{}, errors.Wrap(err, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, organizer.Ownership); err != nil {
		return domain.Contest{}, errors.Wrap(err, 0)
	}

	contest := domain.Contest{
		Ownership: domain.OwnershipData{
			OrganizerID: organizerID,
		},
		Location:           strings.TrimSpace(tmpl.Location),
		Name:               strings.TrimSpace(tmpl.Name),
		Description:        strings.TrimSpace(tmpl.Description),
		QualifyingProblems: tmpl.QualifyingProblems,
		Finalists:          tmpl.Finalists,
		Info:               sanitizationPolicy.Sanitize(tmpl.Info),
		GracePeriod:        tmpl.GracePeriod,
		Created:            time.Now(),
		EvaluationMode:     true,
	}

	if err := (validators.ContestValidator{}).Validate(contest); err != nil {
		return domain.Contest{}, errors.Wrap(err, 0)
	}

	contest, err = uc.Repo.StoreContest(ctx, nil, contest)
	if err != nil {
		return domain.Contest{}, errors.Wrap(err, 0)
	}

	return contest, nil
}

func (uc *ContestUseCase) DuplicateContest(ctx context.Context, contestID domain.ContestID) (domain.Contest, error) {
	contest, err := uc.Repo.GetContest(ctx, nil, contestID)
	if err != nil {
		return domain.Contest{}, errors.Wrap(err, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, contest.Ownership); err != nil {
		return domain.Contest{}, errors.Wrap(err, 0)
	}

	if contest.Archived {
		return domain.Contest{}, errors.Wrap(domain.ErrArchived, 0)
	}

	compClasses, err := uc.Repo.GetCompClassesByContest(ctx, nil, contestID)
	if err != nil {
		return domain.Contest{}, errors.Wrap(err, 0)
	}

	problems, err := uc.Repo.GetProblemsByContest(ctx, nil, contestID)
	if err != nil {
		return domain.Contest{}, errors.Wrap(err, 0)
	}

	duplicatedContest := contest
	duplicatedContest.ID = 0
	duplicatedContest.Name += " (Copy)"

	tx, err := uc.Repo.Begin()
	if err != nil {
		return domain.Contest{}, errors.Wrap(err, 0)
	}

	duplicate := func() (domain.Contest, error) {
		createdContest, err := uc.Repo.StoreContest(ctx, tx, duplicatedContest)
		if err != nil {
			return domain.Contest{}, err
		}

		for _, compClass := range compClasses {
			compClass.ID = 0
			compClass.ContestID = createdContest.ID

			_, err = uc.Repo.StoreCompClass(ctx, tx, compClass)
			if err != nil {
				return domain.Contest{}, err
			}
		}

		for _, problem := range problems {
			problem.ID = 0
			problem.ContestID = createdContest.ID

			_, err = uc.Repo.StoreProblem(ctx, tx, problem)
			if err != nil {
				return domain.Contest{}, err
			}
		}

		return createdContest, nil
	}

	createdContest, err := duplicate()
	if err != nil {
		tx.Rollback()
		return domain.Contest{}, errors.Wrap(err, 0)
	}

	err = tx.Commit()
	if err != nil {
		return domain.Contest{}, errors.Wrap(err, 0)
	}

	return createdContest, nil
}

func (uc *ContestUseCase) TransferContest(ctx context.Context, contestID domain.ContestID, newOrganizerID domain.OrganizerID) (domain.Contest, error) {
	contest, err := uc.Repo.GetContest(ctx, nil, contestID)
	if err != nil {
		return domain.Contest{}, errors.Wrap(err, 0)
	}

	if contest.Archived {
		return domain.Contest{}, errors.Wrap(domain.ErrArchived, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, contest.Ownership); err != nil {
		return domain.Contest{}, errors.Wrap(err, 0)
	}

	newOrganizer, err := uc.Repo.GetOrganizer(ctx, nil, newOrganizerID)
	if err != nil {
		return domain.Contest{}, errors.Wrap(err, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, newOrganizer.Ownership); err != nil {
		return domain.Contest{}, errors.Wrap(err, 0)
	}

	compClasses, err := uc.Repo.GetCompClassesByContest(ctx, nil, contestID)
	if err != nil {
		return domain.Contest{}, errors.Wrap(err, 0)
	}

	problems, err := uc.Repo.GetProblemsByContest(ctx, nil, contestID)
	if err != nil {
		return domain.Contest{}, errors.Wrap(err, 0)
	}

	contenders, err := uc.Repo.GetContendersByContest(ctx, nil, contestID)
	if err != nil {
		return domain.Contest{}, errors.Wrap(err, 0)
	}

	raffles, err := uc.Repo.GetRafflesByContest(ctx, nil, contestID)
	if err != nil {
		return domain.Contest{}, errors.Wrap(err, 0)
	}

	allRaffleWinners := make([]domain.RaffleWinner, 0)

	for _, raffle := range raffles {
		winners, err := uc.Repo.GetRaffleWinners(ctx, nil, raffle.ID)
		if err != nil {
			return domain.Contest{}, errors.Wrap(err, 0)
		}

		allRaffleWinners = append(allRaffleWinners, winners...)
	}

	ticks, err := uc.Repo.GetTicksByContest(ctx, nil, contestID)
	if err != nil {
		return domain.Contest{}, errors.Wrap(err, 0)
	}

	tx, err := uc.Repo.Begin()
	if err != nil {
		return domain.Contest{}, errors.Wrap(err, 0)
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

		for _, contender := range contenders {
			err = uc.Repo.DeleteContender(ctx, tx, contender.ID)
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

		for _, compClass := range compClasses {
			err = uc.Repo.DeleteCompClass(ctx, tx, compClass.ID)
			if err != nil {
				return err
			}
		}

		err := uc.Repo.DeleteContest(ctx, tx, contestID)
		if err != nil {
			return err
		}

		return nil
	}

	recreate := func() error {
		_, err := uc.Repo.StoreContest(ctx, tx, contest)
		if err != nil {
			return err
		}

		for _, compClass := range compClasses {
			_, err = uc.Repo.StoreCompClass(ctx, tx, compClass)
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

		for _, contender := range contenders {
			_, err = uc.Repo.StoreContender(ctx, tx, contender)
			if err != nil {
				return err
			}

			score := contender.Score
			if score == nil {
				continue
			}

			err := uc.Repo.StoreScore(ctx, tx, *score)
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
		return domain.Contest{}, errors.Wrap(err, 0)
	}

	contest.Ownership.OrganizerID = newOrganizerID

	for index := range compClasses {
		compClasses[index].Ownership.OrganizerID = newOrganizerID
	}

	for index := range problems {
		problems[index].Ownership.OrganizerID = newOrganizerID
	}

	for index := range contenders {
		contenders[index].Ownership.OrganizerID = newOrganizerID
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
		return domain.Contest{}, errors.Wrap(err, 0)
	}

	err = tx.Commit()
	if err != nil {
		return domain.Contest{}, errors.Wrap(err, 0)
	}

	return contest, nil
}
