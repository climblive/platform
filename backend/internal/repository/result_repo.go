package repository

import (
	"context"

	"github.com/climblive/platform/backend/internal/database"
	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

func (d *Database) StoreResult(ctx context.Context, tx domain.Transaction, result domain.Result) error {
	params := database.UpsertResultParams{
		CompetitorID: int32(result.CompetitorID),
		Timestamp:    result.Timestamp,
		Score:        result.Score,
		Placement:    int32(result.Placement),
		Finalist:     result.Finalist,
		RankOrder:    int32(result.RankOrder),
	}

	err := d.WithTx(tx).UpsertResult(ctx, params)
	switch {
	case mysqlForeignKeyConstraintViolation.Is(err):
		return errors.New(domain.ErrNotFound)
	case err != nil:
		return errors.Wrap(err, 0)
	}

	return nil
}
