package domain

import "context"

type Problem struct {
	ID                 ResourceID
	ContestID          ResourceID
	Number             int
	HoldColorPrimary   string
	HoldColorSecondary string
	Name               string
	Description        string
	Points             int
	FlashBonus         int
}

type ProblemUsecase interface {
	GetProblem(ctx context.Context, id ResourceID) (Problem, error)
	GetProblemsByContest(ctx context.Context, contestID ResourceID) ([]Problem, error)
	UpdateProblem(ctx context.Context, id ResourceID, problem Problem) (Problem, error)
	DeleteProblem(ctx context.Context, id ResourceID) (error)
	CreateProblem(ctx context.Context, contestID ResourceID, problem Problem) (Problem, error)
}
