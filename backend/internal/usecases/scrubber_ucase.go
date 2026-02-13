package usecases

import (
	"context"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

type scrubberRepository interface {
	domain.Transactor

	GetScrubEligibleContenders(ctx context.Context, deadline time.Time) ([]domain.Contender, error)
	StoreContender(ctx context.Context, tx domain.Transaction, contender domain.Contender) (domain.Contender, error)
}

type ScrubberUseCase struct {
	Repo        scrubberRepository
	EventBroker domain.EventBroker
}

func (uc *ScrubberUseCase) ScrubContenders(ctx context.Context, deadline time.Time) (int, error) {
	contenders, err := uc.Repo.GetScrubEligibleContenders(ctx, deadline)
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

	for i := range contenders {
		contenders[i].Name = ""
		contenders[i].ScrubbedAt = time.Now()

		if _, err := uc.Repo.StoreContender(ctx, tx, contenders[i]); err != nil {
			return 0, errors.Wrap(err, 0)
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, errors.Wrap(err, 0)
	}

	for _, contender := range contenders {
		uc.EventBroker.Dispatch(contender.ContestID, domain.ContenderPublicInfoUpdatedEvent{
			ContenderID:         contender.ID,
			CompClassID:         contender.CompClassID,
			Name:                "",
			WithdrawnFromFinals: contender.WithdrawnFromFinals,
			Disqualified:        contender.Disqualified,
		})
	}

	return len(contenders), nil
}
