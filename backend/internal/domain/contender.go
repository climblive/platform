package domain

import (
	"context"
	"time"
)

type Contender struct {
	ID                  ContenderID   `json:"id,omitempty"`
	Ownership           OwnershipData `json:"-"`
	ContestID           ContestID     `json:"contestId"`
	CompClassID         CompClassID   `json:"compClassId,omitempty"`
	RegistrationCode    string        `json:"registrationCode"`
	Name                string        `json:"name,omitempty"`
	PublicName          string        `json:"publicName,omitempty"`
	ClubName            string        `json:"clubName,omitempty"`
	Entered             *time.Time    `json:"entered,omitempty"`
	WithdrawnFromFinals bool          `json:"withdrawnFromFinals"`
	Disqualified        bool          `json:"disqualified"`
	Score               *Score        `json:"score,omitempty"`
}

type ContenderUseCase interface {
	GetContender(ctx context.Context, contenderID ContenderID) (Contender, error)
	GetContenderByCode(ctx context.Context, registrationCode string) (Contender, error)
	GetContendersByCompClass(ctx context.Context, compClassID CompClassID) ([]Contender, error)
	GetContendersByContest(ctx context.Context, contestID ContestID) ([]Contender, error)
	UpdateContender(ctx context.Context, contenderID ContenderID, contender Contender) (Contender, error)
	DeleteContender(ctx context.Context, contenderID ContenderID) error
	CreateContenders(ctx context.Context, contestID ContestID, number int) ([]Contender, error)
}

type CodeGenerator interface {
	Generate(length int) string
}
