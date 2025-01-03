package domain

import "context"

type OrganizerUseCase interface {
	CreateOrganizer(ctx context.Context, organizer Organizer) (Organizer, error)
	GetOrganizer(ctx context.Context, organizerID OrganizerID) (Organizer, error)
	UpdateOrganizer(ctx context.Context, organizerID OrganizerID, organizer Organizer) (Organizer, error)
	DeleteOrganizer(ctx context.Context, organizerID OrganizerID) error
}
