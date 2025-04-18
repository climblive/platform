package usecases

import (
	"context"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

type compClassUseCaseRepository interface {
	domain.Transactor

	GetCompClassesByContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) ([]domain.CompClass, error)
	GetContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) (domain.Contest, error)
	StoreCompClass(ctx context.Context, tx domain.Transaction, compClass domain.CompClass) (domain.CompClass, error)
}

type CompClassUseCase struct {
	Repo       compClassUseCaseRepository
	Authorizer domain.Authorizer
}

func (uc *CompClassUseCase) GetCompClassesByContest(ctx context.Context, contestID domain.ContestID) ([]domain.CompClass, error) {
	compClasses, err := uc.Repo.GetCompClassesByContest(ctx, nil, contestID)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	return compClasses, nil
}

func (uc *CompClassUseCase) CreateCompClass(ctx context.Context, contestID domain.ContestID, tmpl domain.CompClassTemplate) (domain.CompClass, error) {
	contest, err := uc.Repo.GetContest(ctx, nil, contestID)
	if err != nil {
		return domain.CompClass{}, errors.Wrap(err, 0)
	}

	if _, err := uc.Authorizer.HasOwnership(ctx, contest.Ownership); err != nil {
		return domain.CompClass{}, errors.Wrap(err, 0)
	}

	switch {
	case len(tmpl.Name) < 1:
		fallthrough
	case tmpl.TimeEnd.Before(tmpl.TimeBegin):
		fallthrough
	case tmpl.TimeEnd.Sub(tmpl.TimeBegin) > 12*time.Hour:
		return domain.CompClass{}, errors.Wrap(domain.ErrInvalidData, 0)
	}

	compClass := domain.CompClass{
		Ownership:   contest.Ownership,
		ContestID:   contestID,
		Name:        tmpl.Name,
		Description: tmpl.Description,
		TimeBegin:   tmpl.TimeBegin,
		TimeEnd:     tmpl.TimeEnd,
	}

	createdCompClass, err := uc.Repo.StoreCompClass(ctx, nil, compClass)
	if err != nil {
		return domain.CompClass{}, errors.Wrap(err, 0)
	}

	return createdCompClass, nil
}
