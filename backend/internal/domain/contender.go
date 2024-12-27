package domain

import (
	"context"
	"encoding/json"
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

type Patch[T any] struct {
	Value T
}

func NewPatch[T any](v T) *Patch[T] {
	return &Patch[T]{Value: v}
}

func (p *Patch[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Value)
}

func (p *Patch[T]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &p.Value)
}

type ContenderPatch struct {
	CompClassID         *Patch[CompClassID] `json:"compClassId"`
	Name                *Patch[string]      `json:"name"`
	PublicName          *Patch[string]      `json:"publicName"`
	ClubName            *Patch[string]      `json:"clubName"`
	WithdrawnFromFinals *Patch[bool]        `json:"withdrawnFromFinals"`
	Disqualified        *Patch[bool]        `json:"disqualified"`
}

type ContenderUseCase interface {
	GetContender(ctx context.Context, contenderID ContenderID) (Contender, error)
	GetContenderByCode(ctx context.Context, registrationCode string) (Contender, error)
	GetContendersByCompClass(ctx context.Context, compClassID CompClassID) ([]Contender, error)
	GetContendersByContest(ctx context.Context, contestID ContestID) ([]Contender, error)
	PatchContender(ctx context.Context, contenderID ContenderID, patch ContenderPatch) (Contender, error)
	DeleteContender(ctx context.Context, contenderID ContenderID) error
	CreateContenders(ctx context.Context, contestID ContestID, number int) ([]Contender, error)
}

type CodeGenerator interface {
	Generate(length int) string
}
