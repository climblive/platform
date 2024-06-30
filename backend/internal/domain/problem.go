package domain

import "context"

type Problem struct {
	ID                 ResourceID
	Ownership          OwnershipData
	ContestID          ResourceID
	Number             int
	HoldColorPrimary   string
	HoldColorSecondary string
	Name               string
	Description        string
	PointsTop          int
	PointsZone         int
	FlashBonus         int
}

type ProblemUseCase interface {
	GetProblemsByContest(ctx context.Context, contestID ResourceID) ([]Problem, error)
	UpdateProblem(ctx context.Context, problemID ResourceID, problem Problem) (Problem, error)
	DeleteProblem(ctx context.Context, problemID ResourceID) error
	CreateProblem(ctx context.Context, contestID ResourceID, problem Problem) (Problem, error)
}
