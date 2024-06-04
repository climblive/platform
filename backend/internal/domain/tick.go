package domain

import (
	"context"
	"time"
)

type Tick struct {
	ID          ResourceID
	Timestamp   time.Time
	ContenderID ResourceID
	ProblemID   ResourceID
	Flash       bool
}

type TickUsecase interface {
	GetTicksByContender(ctx context.Context, contenderID ResourceID) ([]Tick, error)
	GetTicksByProblem(ctx context.Context, problemID ResourceID) ([]Tick, error)
	DeleteTick(ctx context.Context, tickID ResourceID) error
	CreateTick(ctx context.Context, contenderID ResourceID, tick Tick) (Tick, error)
}
