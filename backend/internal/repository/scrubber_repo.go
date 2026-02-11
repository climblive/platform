package repository

import (
	"context"

	"github.com/climblive/platform/backend/internal/database"
	"github.com/go-errors/errors"
)

func (d *Database) GetScrubEligibleContenders(ctx context.Context) ([]database.GetScrubEligibleContendersRow, error) {
	contenders, err := d.queries.GetScrubEligibleContenders(ctx)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	return contenders, nil
}

func (d *Database) UpdateContenderScrubbed(ctx context.Context, arg database.UpdateContenderScrubbedParams) error {
	if err := d.queries.UpdateContenderScrubbed(ctx, arg); err != nil {
		return errors.Wrap(err, 0)
	}

	return nil
}
