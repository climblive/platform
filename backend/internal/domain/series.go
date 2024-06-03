package domain

import "context"

type Series struct {
	ID          ResourceID
	OrganizerID ResourceID
	Name        string
}

type SeriesUsecase interface {
	GetSeries(ctx context.Context, organizerID ResourceID) (Series, error)
	GetSeriesByOrganizer(ctx context.Context, organizerID ResourceID) ([]Series, error)
	DeleteSeries(ctx context.Context, id ResourceID) (error)
	CreateSeries(ctx context.Context, organizerID ResourceID) ([]Series, error)
}
