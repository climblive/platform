package usecases

import (
	"context"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

type contestUseCaseRepository interface {
	domain.Transactor

	GetContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) (domain.Contest, error)
	GetContendersByContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) ([]domain.Contender, error)
	GetContestsByOrganizer(ctx context.Context, tx domain.Transaction, organizerID domain.OrganizerID) ([]domain.Contest, error)
	GetOrganizer(ctx context.Context, tx domain.Transaction, organizerID domain.OrganizerID) (domain.Organizer, error)
	CreateContest(ctx context.Context, tx domain.Transaction, contest domain.Contest) (domain.Contest, error)
}

type ContestUseCase struct {
	Authorizer  domain.Authorizer
	Repo        contestUseCaseRepository
	ScoreKeeper domain.ScoreKeeper
}

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

func (uc *ContestUseCase) CreateContest(ctx context.Context, organizerID domain.OrganizerID, tmpl domain.ContestTemplate) (domain.Contest, error) {
	organizer, err := uc.Repo.GetOrganizer(ctx, nil, organizerID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, organizer.Ownership); err != nil {
		return nil, errors.Wrap(err, 0)
	}

	contest := domain.Contest{
		Ownership: domain.OwnershipData{
			OrganizerID: organizerID,
		},
		Location:           tmpl.Location,
		SeriesID:           tmpl.SeriesID,
		Protected:          false,
		Name:               tmpl.Name,
		Description:        tmpl.Description,
		FinalsEnabled:      tmpl.FinalsEnabled,
		QualifyingProblems: tmpl.QualifyingProblems,
		Finalists:          tmpl.Finalists,
		Rules:              tmpl.Rules,
		GracePeriod:        tmpl.GracePeriod,
	}

	contest, err = uc.Repo.CreateContest(ctx, nil, contest)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	return contest, nil
}
