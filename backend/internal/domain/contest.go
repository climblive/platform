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

type Score struct {
	ContenderID         ResourceID
	ContenderPublicName string
	CompClassID         ResourceID
	Score               int
	Placement           int
}

type ContestUsecase interface {
	GetContest(ctx context.Context, contestID ResourceID) (Contest, error)
	GetContestsByOrganizer(ctx context.Context, organizerID ResourceID) ([]Contest, error)
	UpdateContest(ctx context.Context, contestID ResourceID, contest Contest) (Contest, error)
	DeleteContest(ctx context.Context, contestID ResourceID) error
	DuplicateContest(ctx context.Context, contestID ResourceID) (Contest, error)
	CreateContest(ctx context.Context, organizerID ResourceID, contest Contest) (Contest, error)
	GetScores(ctx context.Context, contestID ResourceID) ([]Score, error)
}
