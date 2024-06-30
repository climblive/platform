package repository

import (
	"context"

	"github.com/climblive/platform/backend/internal/domain"
)

func (d *Database) GetCompClass(ctx context.Context, tx domain.Transaction, compClassID domain.ResourceID) (domain.CompClass, error) {
	return domain.CompClass{}, nil
}
