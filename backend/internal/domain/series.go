package domain

import "context"

type Series struct {
	ID        SeriesID
	Ownership OwnershipData
	Name      string
}

type SeriesUseCase interface {
	GetSeriesByOrganizer(ctx context.Context, organizerID OrganizerID) ([]Series, error)
	UpdateSeries(ctx context.Context, seriesID SeriesID, series Series) (Series, error)
	DeleteSeries(ctx context.Context, seriesID SeriesID) error
	CreateSeries(ctx context.Context, organizerID OrganizerID) (Series, error)
}
