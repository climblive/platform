package repository

import (
	"context"
	"time"

	"github.com/climblive/platform/backend/internal/database"
	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

func (d *Database) GetScrubEligibleContenders(ctx context.Context, deadline time.Time) ([]domain.Contender, error) {
	records, err := d.queries.GetScrubEligibleContenders(ctx, makeNullTime(deadline))
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	contenders := make([]domain.Contender, 0)

	for _, record := range records {
		contender := contenderToDomain(database.GetContenderRow(record))

		contenders = append(contenders, contender)
	}

	return contenders, nil
}

func (d *Database) UpdateContenderScrubbed(ctx context.Context, arg database.UpdateContenderScrubbedParams) error {
	if err := d.queries.UpdateContenderScrubbed(ctx, arg); err != nil {
		return errors.Wrap(err, 0)
	}

	return nil
}
