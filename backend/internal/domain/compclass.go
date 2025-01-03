package domain

import (
	"context"
)

type CompClassUseCase interface {
	GetCompClassesByContest(ctx context.Context, contestID ContestID) ([]CompClass, error)
	UpdateCompClass(ctx context.Context, compClassID CompClassID, compClass CompClass) (CompClass, error)
	DeleteCompClass(ctx context.Context, compClassID CompClassID) error
	CreateCompClass(ctx context.Context, contestID ContestID, compClass CompClass) (CompClass, error)
}
