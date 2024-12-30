package repository

import (
	"context"
	"database/sql"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

func (d *Database) GetCompClass(ctx context.Context, tx domain.Transaction, compClassID domain.CompClassID) (domain.CompClass, error) {
	record, err := d.WithTx(tx).GetCompClass(ctx, int32(compClassID))
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return domain.CompClass{}, errors.Wrap(domain.ErrNotFound, 0)
	case err != nil:
		return domain.CompClass{}, errors.Wrap(err, 0)
	}

	return compClassToDomain(record.CompClass), nil
}

func (d *Database) GetCompClassesByContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) ([]domain.CompClass, error) {
	records, err := d.WithTx(tx).GetCompClassesByContest(ctx, int32(contestID))
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	compClasses := make([]domain.CompClass, 0)

	for _, record := range records {
		compClasses = append(compClasses, compClassToDomain(record.CompClass))
	}

	return compClasses, nil
}
