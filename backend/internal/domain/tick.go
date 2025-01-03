package domain

import (
	"context"
)

type TickUseCase interface {
	GetTicksByContender(ctx context.Context, contenderID ContenderID) ([]Tick, error)
	GetTicksByProblem(ctx context.Context, problemID ProblemID) ([]Tick, error)
	DeleteTick(ctx context.Context, tickID TickID) error
	CreateTick(ctx context.Context, contenderID ContenderID, tick Tick) (Tick, error)
}
