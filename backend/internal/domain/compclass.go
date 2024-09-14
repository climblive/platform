package domain

import (
	"context"
	"time"
)

type ColorRGB string

type CompClass struct {
	ID          ResourceID    `json:"id,omitempty"`
	Ownership   OwnershipData `json:"-"`
	ContestID   ResourceID    `json:"contestId"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Color       ColorRGB      `json:"color"`
	TimeBegin   time.Time     `json:"timeBegin"`
	TimeEnd     time.Time     `json:"timeEnd"`
}

type CompClassUseCase interface {
	GetCompClassesByContest(ctx context.Context, contestID ResourceID) ([]CompClass, error)
	UpdateCompClass(ctx context.Context, compClassID ResourceID, compClass CompClass) (CompClass, error)
	DeleteCompClass(ctx context.Context, compClassID ResourceID) error
	CreateCompClass(ctx context.Context, contestID ResourceID, compClass CompClass) (CompClass, error)
}
