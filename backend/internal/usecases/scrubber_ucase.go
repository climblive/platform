package usecases

import (
	"context"
	"database/sql"
	"time"

	"github.com/climblive/platform/backend/internal/database"
	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

type scrubberRepository interface {
	domain.Transactor

	GetScrubEligibleContenders(ctx context.Context) ([]database.GetScrubEligibleContendersRow, error)
	UpdateContenderScrubbed(ctx context.Context, arg database.UpdateContenderScrubbedParams) error
}

type ScrubberUseCase struct {
	Repo scrubberRepository
}

func (uc *ScrubberUseCase) ScrubOldContenders(ctx context.Context) (int, error) {
	contenders, err := uc.Repo.GetScrubEligibleContenders(ctx)
	if err != nil {
		return 0, errors.Wrap(err, 0)
	}

	if len(contenders) == 0 {
		return 0, nil
	}

	tx, err := uc.Repo.Begin()
	if err != nil {
		return 0, errors.Wrap(err, 0)
	}
	defer tx.Rollback()

	now := sql.NullTime{Time: time.Now(), Valid: true}
	for _, contender := range contenders {
		params := database.UpdateContenderScrubbedParams{
			ScrubbedAt: now,
			ID:         contender.Contender.ID,
		}

		if err := uc.Repo.UpdateContenderScrubbed(ctx, params); err != nil {
			return 0, errors.Wrap(err, 0)
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, errors.Wrap(err, 0)
	}

	return len(contenders), nil
}
