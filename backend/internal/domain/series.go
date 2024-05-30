package domain

type Series struct {
	ID          ResourceID
	OrganizerID ResourceID
	Name        string
}

type SeriesUsecase interface {
}
