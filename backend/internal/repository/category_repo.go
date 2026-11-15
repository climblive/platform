package repository

import (
	"context"
	"database/sql"

	"github.com/climblive/platform/backend/internal/database"
	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

func (d *Database) GetCategory(ctx context.Context, tx domain.Transaction, categoryID domain.CategoryID) (domain.Category, error) {
	record, err := d.WithTx(tx).GetCategory(ctx, int32(categoryID))
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return domain.Category{}, errors.Wrap(domain.ErrNotFound, 0)
	case err != nil:
		return domain.Category{}, errors.Wrap(err, 0)
	}

	return categoryToDomain(record.Category), nil
}

func (d *Database) GetCategoriesByCompetition(ctx context.Context, tx domain.Transaction, competitionID domain.CompetitionID) ([]domain.Category, error) {
	records, err := d.WithTx(tx).GetCategoriesByCompetition(ctx, int32(competitionID))
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	categories := make([]domain.Category, 0)

	for _, record := range records {
		categories = append(categories, categoryToDomain(record.Category))
	}

	return categories, nil
}

func (d *Database) StoreCategory(ctx context.Context, tx domain.Transaction, category domain.Category) (domain.Category, error) {
	params := database.UpsertCategoryParams{
		ID:            int32(category.ID),
		OrganizerID:   int32(category.Ownership.OrganizerID),
		CompetitionID: int32(category.CompetitionID),
		Name:          category.Name,
		Description:   makeNullString(category.Description),
		Color:         sql.NullString{String: "", Valid: false},
		TimeBegin:     category.TimeBegin,
		TimeEnd:       category.TimeEnd,
	}

	insertID, err := d.WithTx(tx).UpsertCategory(ctx, params)
	if err != nil {
		return domain.Category{}, errors.Wrap(err, 0)
	}

	if insertID != 0 {
		category.ID = domain.CategoryID(insertID)
	}

	return category, err
}

func (d *Database) DeleteCategory(ctx context.Context, tx domain.Transaction, categoryID domain.CategoryID) error {
	err := d.WithTx(tx).DeleteCategory(ctx, int32(categoryID))
	if err != nil {
		return errors.Wrap(err, 0)
	}

	return nil
}
