package domain

import (
	"context"
	"time"
)

type ColorRGB string

type CompClass struct {
	ID          ResourceID
	ContestID   ResourceID
	Name        string
	Description string
	Color       ColorRGB
	TimeBegin   time.Time
	TimeEnd     time.Time
}

type CompClassUsecase interface {
	GetCompClassesByContest(ctx context.Context, contestID ResourceID) ([]CompClass, error)
	UpdateCompClass(ctx context.Context, compClassID ResourceID, compClass CompClass) (CompClass, error)
	DeleteCompClass(ctx context.Context, compClassID ResourceID) error
	CreateCompClass(ctx context.Context, contestID ResourceID, compClass CompClass) (CompClass, error)
}
