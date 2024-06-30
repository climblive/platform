package domain

import (
	"context"
	"time"
)

type Contender struct {
	ID                  ResourceID
	Ownership           OwnershipData
	ContestID           ResourceID
	CompClassID         ResourceID
	RegistrationCode    string
	Name                string
	PublicName          string
	ClubName            string
	Entered             time.Time
	WithdrawnFromFinals bool
	Disqualified        bool
	Score               int
	Placement           int
	ScoreTimestamp      time.Time
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
