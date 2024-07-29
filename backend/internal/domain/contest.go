package domain

import (
	"context"
	"time"
)

type Contest struct {
	ID                 ResourceID
	Ownership          OwnershipData
	Location           string
	SeriesID           ResourceID
	Protected          bool
	Name               string
	Description        string
	FinalsEnabled      bool
	QualifyingProblems int
	Finalists          int
	Rules              string
	GracePeriod        time.Duration
}

type ContestUseCase interface {
	GetContest(ctx context.Context, contestID ResourceID) (Contest, error)
	GetContestsByOrganizer(ctx context.Context, organizerID ResourceID) ([]Contest, error)
	UpdateContest(ctx context.Context, contestID ResourceID, contest Contest) (Contest, error)
	DeleteContest(ctx context.Context, contestID ResourceID) error
	DuplicateContest(ctx context.Context, contestID ResourceID) (Contest, error)
	CreateContest(ctx context.Context, organizerID ResourceID, contest Contest) (Contest, error)
	GetScores(ctx context.Context, contestID ResourceID) ([]Score, error)
}
