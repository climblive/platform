package domain

import "context"

type Series struct {
	ID          ResourceID
	OrganizerID ResourceID
	Name        string
}

type SeriesUsecase interface {
	GetSeriesByOrganizer(ctx context.Context, organizerID ResourceID) ([]Series, error)
	UpdateSeries(ctx context.Context, seriesID ResourceID, series Series) (Series, error)
	DeleteSeries(ctx context.Context, seriesID ResourceID) error
	CreateSeries(ctx context.Context, organizerID ResourceID) (Series, error)
}
