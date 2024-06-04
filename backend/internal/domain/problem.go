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
	GetProblemsByContest(ctx context.Context, contestID ResourceID) ([]Problem, error)
	UpdateProblem(ctx context.Context, problemID ResourceID, problem Problem) (Problem, error)
	DeleteProblem(ctx context.Context, problemID ResourceID) error
	CreateProblem(ctx context.Context, contestID ResourceID, problem Problem) (Problem, error)
}
