package usecases

import (
	"context"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/usecases/validators"
	"github.com/go-errors/errors"
)

const maxCompClassesPerCompetition = 20

type compClassUseCaseRepository interface {
	domain.Transactor

	GetCompClassesByCompetition(ctx context.Context, tx domain.Transaction, competitionID domain.CompetitionID) ([]domain.CompClass, error)
	GetCompetition(ctx context.Context, tx domain.Transaction, competitionID domain.CompetitionID) (domain.Competition, error)
	StoreCompClass(ctx context.Context, tx domain.Transaction, compClass domain.CompClass) (domain.CompClass, error)
	DeleteCompClass(ctx context.Context, tx domain.Transaction, compClassID domain.CompClassID) error
	GetCompClass(ctx context.Context, tx domain.Transaction, compClassID domain.CompClassID) (domain.CompClass, error)
	GetContendersByCompClass(ctx context.Context, tx domain.Transaction, compClassID domain.CompClassID) ([]domain.Contender, error)
}

type CompClassUseCase struct {
	Repo       compClassUseCaseRepository
	Authorizer domain.Authorizer
}

func (uc *CompClassUseCase) GetCompClass(ctx context.Context, compClassID domain.CompClassID) (domain.CompClass, error) {
	compClass, err := uc.Repo.GetCompClass(ctx, nil, compClassID)
	if err != nil {
		return domain.CompClass{}, errors.Wrap(err, 0)
	}

	return compClass, nil
}

func (uc *CompClassUseCase) GetCompClassesByCompetition(ctx context.Context, competitionID domain.CompetitionID) ([]domain.CompClass, error) {
	compClasses, err := uc.Repo.GetCompClassesByCompetition(ctx, nil, competitionID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	return compClasses, nil
}

func (uc *CompClassUseCase) CreateCompClass(ctx context.Context, competitionID domain.CompetitionID, tmpl domain.CompClassTemplate) (domain.CompClass, error) {
	competition, err := uc.Repo.GetCompetition(ctx, nil, competitionID)
	if err != nil {
		return domain.CompClass{}, errors.Wrap(err, 0)
	}

	_, err = uc.Authorizer.HasOwnership(ctx, competition.Ownership)
	if err != nil {
		return domain.CompClass{}, errors.Wrap(err, 0)
	}

	compClasses, err := uc.Repo.GetCompClassesByCompetition(ctx, nil, competitionID)
	if err != nil {
		return domain.CompClass{}, errors.Wrap(err, 0)
	}

	if len(compClasses) >= maxCompClassesPerCompetition {
		return domain.CompClass{}, errors.New(domain.ErrLimitExceeded)
	}

	compClass := domain.CompClass{
		ID:            0,
		Ownership:     competition.Ownership,
		CompetitionID: competitionID,
		Name:          tmpl.Name,
		Description:   tmpl.Description,
		TimeBegin:     tmpl.TimeBegin,
		TimeEnd:       tmpl.TimeEnd,
	}

	if err := (validators.CompClassValidator{}).Validate(compClass); err != nil {
		return domain.CompClass{}, errors.Wrap(err, 0)
	}

	createdCompClass, err := uc.Repo.StoreCompClass(ctx, nil, compClass)
	if err != nil {
		return domain.CompClass{}, errors.Wrap(err, 0)
	}

	return createdCompClass, nil
}

func (uc *CompClassUseCase) DeleteCompClass(ctx context.Context, compClassID domain.CompClassID) error {
	compClass, err := uc.Repo.GetCompClass(ctx, nil, compClassID)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, compClass.Ownership); err != nil {
		return errors.Wrap(err, 0)
	}

	contenders, err := uc.Repo.GetContendersByCompClass(ctx, nil, compClassID)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	if len(contenders) > 0 {
		return errors.Wrap(domain.ErrNotAllowed, 0)
	}

	err = uc.Repo.DeleteCompClass(ctx, nil, compClassID)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	return nil
}

func (uc *CompClassUseCase) PatchCompClass(ctx context.Context, compClassID domain.CompClassID, patch domain.CompClassPatch) (domain.CompClass, error) {
	compClass, err := uc.Repo.GetCompClass(ctx, nil, compClassID)
	if err != nil {
		return domain.CompClass{}, errors.Wrap(err, 0)
	}

	if _, err = uc.Authorizer.HasOwnership(ctx, compClass.Ownership); err != nil {
		return domain.CompClass{}, errors.Wrap(err, 0)
	}

	if patch.Name.Present {
		compClass.Name = patch.Name.Value
	}

	if patch.Description.Present {
		compClass.Description = patch.Description.Value
	}

	if patch.TimeBegin.Present {
		compClass.TimeBegin = patch.TimeBegin.Value
	}

	if patch.TimeEnd.Present {
		compClass.TimeEnd = patch.TimeEnd.Value
	}

	if err := (validators.CompClassValidator{}).Validate(compClass); err != nil {
		return domain.CompClass{}, errors.Wrap(err, 0)
	}

	if _, err = uc.Repo.StoreCompClass(ctx, nil, compClass); err != nil {
		return domain.CompClass{}, errors.Wrap(err, 0)
	}

	return compClass, nil
}
