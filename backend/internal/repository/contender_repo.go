package repository

import (
	"context"

	"github.com/climblive/platform/backend/internal/domain"
)

func (d *Database) GetContender(ctx context.Context, tx domain.Transaction, contenderID domain.ResourceID) (domain.Contender, error) {
	return domain.Contender{}, nil
}

func (d *Database) GetContenderByCode(ctx context.Context, tx domain.Transaction, registrationCode string) (domain.Contender, error) {
	return domain.Contender{}, nil
}

func (d *Database) GetContendersByCompClass(ctx context.Context, tx domain.Transaction, compClassID domain.ResourceID) ([]domain.Contender, error) {
	return nil, nil
}

func (d *Database) GetContendersByContest(ctx context.Context, tx domain.Transaction, contestID domain.ResourceID) ([]domain.Contender, error) {
	return nil, nil
}

func (d *Database) StoreContender(ctx context.Context, tx domain.Transaction, contender domain.Contender) error {
	return nil
}

func (d *Database) DeleteContender(ctx context.Context, tx domain.Transaction, contenderID domain.ResourceID) error {
	return nil
}
