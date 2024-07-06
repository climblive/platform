package domain

import (
	"context"
	"time"
)

type Contender struct {
	ID                  ResourceID    `json:"id"`
	Ownership           OwnershipData `json:"-"`
	ContestID           ResourceID    `json:"contestId"`
	CompClassID         ResourceID    `json:"compClassId,omitempty"`
	RegistrationCode    string        `json:"registrationCode"`
	Name                string        `json:"name"`
	PublicName          string        `json:"publicName"`
	ClubName            string        `json:"clubName"`
	Entered             time.Time     `json:"entered,omitempty"`
	WithdrawnFromFinals bool          `json:"withdrawnFromFinals"`
	Disqualified        bool          `json:"disqualified"`
	Score               int           `json:"score"`
	Placement           int           `json:"placement"`
	ScoreUpdated        time.Time     `json:"scoreUpdated"`
}

type ContenderUseCase interface {
	GetContender(ctx context.Context, contenderID ResourceID) (Contender, error)
	GetContenderByCode(ctx context.Context, registrationCode string) (Contender, error)
	GetContendersByCompClass(ctx context.Context, compClassID ResourceID) ([]Contender, error)
	GetContendersByContest(ctx context.Context, contestID ResourceID) ([]Contender, error)
	UpdateContender(ctx context.Context, contenderID ResourceID, contender Contender) (Contender, error)
	DeleteContender(ctx context.Context, contenderID ResourceID) error
	CreateContenders(ctx context.Context, contestID ResourceID, number int) ([]Contender, error)
}
