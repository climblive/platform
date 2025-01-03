package domain

import "context"

type ProblemUseCase interface {
	GetProblemsByContest(ctx context.Context, contestID ContestID) ([]Problem, error)
	UpdateProblem(ctx context.Context, problemID ProblemID, problem Problem) (Problem, error)
	DeleteProblem(ctx context.Context, problemID ProblemID) error
	CreateProblem(ctx context.Context, contestID ContestID, problem Problem) (Problem, error)
}
