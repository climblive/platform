package domain

import (
	"context"
	"time"
)

type Contest struct {
	ID                 ContestID     `json:"id,omitempty"`
	Ownership          OwnershipData `json:"-"`
	Location           string        `json:"location,omitempty"`
	SeriesID           SeriesID      `json:"seriesId,omitempty"`
	Protected          bool          `json:"protected"`
	Name               string        `json:"name"`
	Description        string        `json:"description,omitempty"`
	FinalsEnabled      bool          `json:"finalsEnabled"`
	QualifyingProblems int           `json:"qualifyingProblems"`
	Finalists          int           `json:"finalists"`
	Rules              string        `json:"rules,omitempty"`
	GracePeriod        time.Duration `json:"gracePeriod"`
}

type ContestUseCase interface {
	GetContest(ctx context.Context, contestID ContestID) (Contest, error)
	GetContestsByOrganizer(ctx context.Context, organizerID OrganizerID) ([]Contest, error)
	UpdateContest(ctx context.Context, contestID ContestID, contest Contest) (Contest, error)
	DeleteContest(ctx context.Context, contestID ContestID) error
	DuplicateContest(ctx context.Context, contestID ContestID) (Contest, error)
	CreateContest(ctx context.Context, organizerID OrganizerID, contest Contest) (Contest, error)
	GetScoreboard(ctx context.Context, contestID ContestID) ([]ScoreboardEntry, error)
}
