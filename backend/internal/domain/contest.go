package domain

import (
	"context"
	"time"
)

type Contest struct {
	ID                 ResourceID    `json:"id,omitempty"`
	Ownership          OwnershipData `json:"-"`
	Location           string        `json:"location"`
	SeriesID           ResourceID    `json:"seriesId"`
	Protected          bool          `json:"protected"`
	Name               string        `json:"name"`
	Description        string        `json:"description"`
	FinalsEnabled      bool          `json:"finalsEnabled"`
	QualifyingProblems int           `json:"qualifyingProblems"`
	Finalists          int           `json:"finalists"`
	Rules              string        `json:"rules"`
	GracePeriod        time.Duration `json:"gracePeriod"`
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
