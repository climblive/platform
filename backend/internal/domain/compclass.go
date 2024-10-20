package domain

import (
	"context"
	"time"
)

type ColorRGB string

type CompClass struct {
	ID          CompClassID   `json:"id,omitempty"`
	Ownership   OwnershipData `json:"-"`
	ContestID   ContestID     `json:"contestId"`
	Name        string        `json:"name"`
	Description string        `json:"description,omitempty"`
	Color       ColorRGB      `json:"color,omitempty"`
	TimeBegin   time.Time     `json:"timeBegin"`
	TimeEnd     time.Time     `json:"timeEnd"`
}

type CompClassUseCase interface {
	GetCompClassesByContest(ctx context.Context, contestID ContestID) ([]CompClass, error)
	UpdateCompClass(ctx context.Context, compClassID CompClassID, compClass CompClass) (CompClass, error)
	DeleteCompClass(ctx context.Context, compClassID CompClassID) error
	CreateCompClass(ctx context.Context, contestID ContestID, compClass CompClass) (CompClass, error)
}
