package usecases

import (
	"context"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/usecases/validators"
	"github.com/go-errors/errors"
)

const maxCategoriesPerCompetition = 20

type categoryUseCaseRepository interface {
	domain.Transactor

	GetCategoriesByCompetition(ctx context.Context, tx domain.Transaction, competitionID domain.CompetitionID) ([]domain.Category, error)
	GetCompetition(ctx context.Context, tx domain.Transaction, competitionID domain.CompetitionID) (domain.Competition, error)
	StoreCategory(ctx context.Context, tx domain.Transaction, category domain.Category) (domain.Category, error)
	DeleteCategory(ctx context.Context, tx domain.Transaction, categoryID domain.CategoryID) error
	GetCategory(ctx context.Context, tx domain.Transaction, categoryID domain.CategoryID) (domain.Category, error)
	GetCompetitorsByCategory(ctx context.Context, tx domain.Transaction, categoryID domain.CategoryID) ([]domain.Competitor, error)
}

type CategoryUseCase struct {
	Repo       categoryUseCaseRepository
	Authorizer domain.Authorizer
}

func (uc *CategoryUseCase) GetCategory(ctx context.Context, categoryID domain.CategoryID) (domain.Category, error) {
	category, err := uc.Repo.GetCategory(ctx, nil, categoryID)
	if err != nil {
		return domain.Category{}, errors.Wrap(err, 0)
	}

	return category, nil
}

func (uc *CategoryUseCase) GetCategoriesByCompetition(ctx context.Context, competitionID domain.CompetitionID) ([]domain.Category, error) {
	categories, err := uc.Repo.GetCategoriesByCompetition(ctx, nil, competitionID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	return categories, nil
}

func (uc *CategoryUseCase) CreateCategory(ctx context.Context, competitionID domain.CompetitionID, tmpl domain.CategoryTemplate) (domain.Category, error) {
	competition, err := uc.Repo.GetCompetition(ctx, nil, competitionID)
	if err != nil {
		return domain.Category{}, errors.Wrap(err, 0)
	}

	_, err = uc.Authorizer.HasOwnership(ctx, competition.Ownership)
	if err != nil {
		return domain.Category{}, errors.Wrap(err, 0)
	}

	categories, err := uc.Repo.GetCategoriesByCompetition(ctx, nil, competitionID)
	if err != nil {
		return domain.Category{}, errors.Wrap(err, 0)
	}

	if len(categories) >= maxCategoriesPerCompetition {
		return domain.Category{}, errors.New(domain.ErrLimitExceeded)
	}

	category := domain.Category{
		ID:            0,
		Ownership:     competition.Ownership,
		CompetitionID: competitionID,
		Name:          tmpl.Name,
		Description:   tmpl.Description,
		TimeBegin:     tmpl.TimeBegin,
		TimeEnd:       tmpl.TimeEnd,
	}

	if err := (validators.CategoryValidator{}).Validate(category); err != nil {
		return domain.Category{}, errors.Wrap(err, 0)
	}

	createdCategory, err := uc.Repo.StoreCategory(ctx, nil, category)
	if err != nil {
		return domain.Category{}, errors.Wrap(err, 0)
	}

	return createdCategory, nil
}

func (uc *CategoryUseCase) DeleteCategory(ctx context.Context, categoryID domain.CategoryID) error {
	category, err := uc.Repo.GetCategory(ctx, nil, categoryID)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, category.Ownership); err != nil {
		return errors.Wrap(err, 0)
	}

	competitors, err := uc.Repo.GetCompetitorsByCategory(ctx, nil, categoryID)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	if len(competitors) > 0 {
		return errors.Wrap(domain.ErrNotAllowed, 0)
	}

	err = uc.Repo.DeleteCategory(ctx, nil, categoryID)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	return nil
}

func (uc *CategoryUseCase) PatchCategory(ctx context.Context, categoryID domain.CategoryID, patch domain.CategoryPatch) (domain.Category, error) {
	category, err := uc.Repo.GetCategory(ctx, nil, categoryID)
	if err != nil {
		return domain.Category{}, errors.Wrap(err, 0)
	}

	if _, err = uc.Authorizer.HasOwnership(ctx, category.Ownership); err != nil {
		return domain.Category{}, errors.Wrap(err, 0)
	}

	if patch.Name.Present {
		category.Name = patch.Name.Value
	}

	if patch.Description.Present {
		category.Description = patch.Description.Value
	}

	if patch.TimeBegin.Present {
		category.TimeBegin = patch.TimeBegin.Value
	}

	if patch.TimeEnd.Present {
		category.TimeEnd = patch.TimeEnd.Value
	}

	if err := (validators.CategoryValidator{}).Validate(category); err != nil {
		return domain.Category{}, errors.Wrap(err, 0)
	}

	if _, err = uc.Repo.StoreCategory(ctx, nil, category); err != nil {
		return domain.Category{}, errors.Wrap(err, 0)
	}

	return category, nil
}
