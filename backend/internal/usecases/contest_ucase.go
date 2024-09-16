package usecases

import (
	"context"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

type contestUseCaseRepository interface {
	domain.Transactor

	GetContest(ctx context.Context, tx domain.Transaction, contestID domain.ResourceID) (domain.Contest, error)
	GetContendersByContest(ctx context.Context, tx domain.Transaction, contestID domain.ResourceID) ([]domain.Contender, error)
}

type ContestUseCase struct {
	Repo        contestUseCaseRepository
	ScoreKeeper domain.ScoreKeeper
}

func (uc *ContestUseCase) GetContest(ctx context.Context, contestID domain.ResourceID) (domain.Contest, error) {
	contest, err := uc.Repo.GetContest(ctx, nil, contestID)
	if err != nil {
		return domain.Contest{}, errors.Wrap(err, 0)
	}

	return contest, nil
}

func (uc *ContestUseCase) GetContestsByOrganizer(ctx context.Context, organizerID domain.ResourceID) ([]domain.Contest, error) {
	panic("not implemented")
}

func (uc *ContestUseCase) UpdateContest(ctx context.Context, contestID domain.ResourceID, contest domain.Contest) (domain.Contest, error) {
	panic("not implemented")
}

func (uc *ContestUseCase) DeleteContest(ctx context.Context, contestID domain.ResourceID) error {
	panic("not implemented")
}

func (uc *ContestUseCase) DuplicateContest(ctx context.Context, contestID domain.ResourceID) (domain.Contest, error) {
	panic("not implemented")
}

func (uc *ContestUseCase) CreateContest(ctx context.Context, organizerID domain.ResourceID, contest domain.Contest) (domain.Contest, error) {
	panic("not implemented")
}

func (uc *ContestUseCase) GetScoreboard(ctx context.Context, contestID domain.ResourceID) ([]domain.ScoreboardEntry, error) {
	contenders, err := uc.Repo.GetContendersByContest(ctx, nil, contestID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	var entries []domain.ScoreboardEntry

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
			ScoreUpdated:        contender.ScoreUpdated,
			Score:               contender.Score,
			Placement:           contender.Placement,
			Finalist:            contender.Finalist,
		}

		if score, err := uc.ScoreKeeper.GetScore(contender.ID); err == nil {
			entry.Score = score.Score
			entry.Placement = score.Placement
			entry.Finalist = score.Finalist
			entry.ScoreUpdated = score.Timestamp
		}

		entries = append(entries, entry)
	}

	return entries, nil
}
