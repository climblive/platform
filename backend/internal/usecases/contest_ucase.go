package usecases

import (
	"context"
	"strings"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/usecases/validators"
	"github.com/go-errors/errors"
	"github.com/microcosm-cc/bluemonday"
)

type contestUseCaseRepository interface {
	domain.Transactor

	GetContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) (domain.Contest, error)
	GetContendersByContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) ([]domain.Contender, error)
	GetContestsByOrganizer(ctx context.Context, tx domain.Transaction, organizerID domain.OrganizerID) ([]domain.Contest, error)
	GetOrganizer(ctx context.Context, tx domain.Transaction, organizerID domain.OrganizerID) (domain.Organizer, error)
	StoreContest(ctx context.Context, tx domain.Transaction, contest domain.Contest) (domain.Contest, error)
	GetCompClassesByContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) ([]domain.CompClass, error)
	StoreCompClass(ctx context.Context, tx domain.Transaction, compClass domain.CompClass) (domain.CompClass, error)
	GetProblemsByContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) ([]domain.Problem, error)
	StoreProblem(ctx context.Context, tx domain.Transaction, problem domain.Problem) (domain.Problem, error)
}

type ContestUseCase struct {
	Authorizer  domain.Authorizer
	Repo        contestUseCaseRepository
	ScoreKeeper domain.ScoreKeeper
}

var sanitizationPolicy = bluemonday.UGCPolicy()

func (uc *ContestUseCase) GetContest(ctx context.Context, contestID domain.ContestID) (domain.Contest, error) {
	contest, err := uc.Repo.GetContest(ctx, nil, contestID)
	if err != nil {
		return domain.Contest{}, errors.Wrap(err, 0)
	}

	return contest, nil
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
			PublicName:          contender.PublicName,
			ClubName:            contender.ClubName,
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

	if patch.Rules.Present {
		contest.Rules = sanitizationPolicy.Sanitize(patch.Rules.Value)
	}

	if patch.GracePeriod.Present {
		contest.GracePeriod = patch.GracePeriod.Value
	}

	if err := (validators.ContestValidator{}).Validate(contest); err != nil {
		return mty, errors.Wrap(err, 0)
	}

	if _, err = uc.Repo.StoreContest(ctx, nil, contest); err != nil {
		return mty, errors.Wrap(err, 0)
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
		Rules:              sanitizationPolicy.Sanitize(tmpl.Rules),
		GracePeriod:        tmpl.GracePeriod,
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
