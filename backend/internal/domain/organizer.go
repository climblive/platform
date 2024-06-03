package domain

import "context"

type Organizer struct {
	ID       ResourceID
	Name     string
	Homepage string
}

type OrganizerUsecase interface {
	CreateOrganizer(ctx context.Context, organizer Organizer) (Organizer, error)
	GetOrganizer(ctx context.Context, id ResourceID) (Organizer, error)
	UpdateOrganizer(ctx context.Context, id ResourceID, organizer Organizer) (Organizer, error)
	DeleteOrganizer(ctx context.Context, id ResourceID) error
}
