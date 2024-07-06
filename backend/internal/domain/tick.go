package domain

import (
	"context"
	"time"
)

type Tick struct {
	ID           ResourceID
	Ownership    OwnershipData
	Timestamp    time.Time
	ProblemID    ResourceID
	Top          bool
	AttemptsTop  int
	Zone         bool
	AttemptsZone int
}

type TickUseCase interface {
	GetTicksByContender(ctx context.Context, contenderID ResourceID) ([]Tick, error)
	GetTicksByProblem(ctx context.Context, problemID ResourceID) ([]Tick, error)
	DeleteTick(ctx context.Context, tickID ResourceID) error
	CreateTick(ctx context.Context, contenderID ResourceID, tick Tick) (Tick, error)
}
