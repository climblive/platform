package domain

import "context"

type Contest struct {
	ID                 ResourceID
	Location           string
	OrganizerID        ResourceID
	SeriesID           ResourceID
	Protected          bool
	Name               string
	Description        string
	FinalEnabled       bool
	QualifyingProblems int
	Finalists          int
	Rules              string
	GracePeriod        int
}

type ContestUsecase interface {
	GetContest(ctx context.Context, contestID ResourceID) (Contest, error)
}
