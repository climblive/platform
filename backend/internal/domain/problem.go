package domain

import "context"

type Problem struct {
	ID                 ProblemID     `json:"id,omitempty"`
	Ownership          OwnershipData `json:"-"`
	ContestID          ContestID     `json:"contestId"`
	Number             int           `json:"number"`
	HoldColorPrimary   string        `json:"holdColorPrimary"`
	HoldColorSecondary string        `json:"holdColorSecondary,omitempty"`
	Name               string        `json:"name,omitempty"`
	Description        string        `json:"description,omitempty"`
	PointsTop          int           `json:"pointsTop"`
	PointsZone         int           `json:"pointsZone"`
	FlashBonus         int           `json:"flashBonus,omitempty"`
}

type ProblemUseCase interface {
	GetProblemsByContest(ctx context.Context, contestID ContestID) ([]Problem, error)
	UpdateProblem(ctx context.Context, problemID ProblemID, problem Problem) (Problem, error)
	DeleteProblem(ctx context.Context, problemID ProblemID) error
	CreateProblem(ctx context.Context, contestID ContestID, problem Problem) (Problem, error)
}
