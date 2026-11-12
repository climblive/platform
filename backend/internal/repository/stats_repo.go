package repository

import (
	"context"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

func (d *Database) GetGlobalStats(ctx context.Context) (domain.GlobalStats, error) {
	record, err := d.queries.GetGlobalStats(ctx)
	if err != nil {
		return domain.GlobalStats{}, errors.Wrap(err, 0)
	}

	return domain.GlobalStats{
		Contests:             record.Contests,
		RegisteredContenders: record.RegisteredContenders,
		Ticks:                record.Ticks,
	}, nil
}
