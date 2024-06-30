package repository

import (
	"context"

	"github.com/climblive/platform/backend/internal/domain"
)

func (d *Database) GetContest(ctx context.Context, tx domain.Transaction, contestID domain.ResourceID) (domain.Contest, error) {
	return domain.Contest{}, nil
}
