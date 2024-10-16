package domain

import (
	"context"
	"time"
)

type Contender struct {
	ID                  ResourceID    `json:"id,omitempty"`
	Ownership           OwnershipData `json:"-"`
	ContestID           ResourceID    `json:"contestId"`
	CompClassID         ResourceID    `json:"compClassId,omitempty"`
	RegistrationCode    string        `json:"registrationCode"`
	Name                string        `json:"name,omitempty"`
	PublicName          string        `json:"publicName,omitempty"`
	ClubName            string        `json:"clubName,omitempty"`
	Entered             *time.Time    `json:"entered,omitempty"`
	WithdrawnFromFinals bool          `json:"withdrawnFromFinals"`
	Disqualified        bool          `json:"disqualified"`
	Score               int           `json:"score"`
	Placement           int           `json:"placement,omitempty"`
	Finalist            bool          `json:"finalist"`
	RankOrder           int           `json:"rankOrder"`
	ScoreUpdated        *time.Time    `json:"scoreUpdated,omitempty"`
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

type CodeGenerator interface {
	Generate(length int) string
}
