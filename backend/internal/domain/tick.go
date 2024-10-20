package domain

import (
	"context"
	"time"
)

type Tick struct {
	ID           TickID        `json:"id,omitempty"`
	Ownership    OwnershipData `json:"-"`
	Timestamp    time.Time     `json:"timestamp"`
	ContestID    ContestID     `json:"contestId"`
	ProblemID    ProblemID     `json:"problemId"`
	Top          bool          `json:"top"`
	AttemptsTop  int           `json:"attemptsTop"`
	Zone         bool          `json:"zone"`
	AttemptsZone int           `json:"attemptsZone"`
}

type TickUseCase interface {
	GetTicksByContender(ctx context.Context, contenderID ContenderID) ([]Tick, error)
	GetTicksByProblem(ctx context.Context, problemID ProblemID) ([]Tick, error)
	DeleteTick(ctx context.Context, tickID TickID) error
	CreateTick(ctx context.Context, contenderID ContenderID, tick Tick) (Tick, error)
}
