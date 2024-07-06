package domain

import "context"

type Series struct {
	ID        ResourceID
	Ownership OwnershipData
	Name      string
}

type SeriesUseCase interface {
	GetSeriesByOrganizer(ctx context.Context, organizerID ResourceID) ([]Series, error)
	UpdateSeries(ctx context.Context, seriesID ResourceID, series Series) (Series, error)
	DeleteSeries(ctx context.Context, seriesID ResourceID) error
	CreateSeries(ctx context.Context, organizerID ResourceID) (Series, error)
}
