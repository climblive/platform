package repository

import (
	"context"

	"github.com/climblive/platform/backend/internal/database"
	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

func (d *Database) StoreScore(ctx context.Context, tx domain.Transaction, score domain.Score) error {
	params := database.UpsertScoreParams{
		ContenderID: int32(score.ContenderID),
		Timestamp:   score.Timestamp,
		Score:       int32(score.Score),
		Placement:   int32(score.Placement),
		Finalist:    score.Finalist,
		RankOrder:   int32(score.RankOrder),
	}

	err := d.WithTx(tx).UpsertScore(ctx, params)
	switch {
	case mysqlForeignKeyConstraintViolation.Is(err):
		return errors.New(domain.ErrNotFound)
	case err != nil:
		return errors.Wrap(err, 0)
	}

	return nil
}
