package domain

import (
	"context"
	"time"
)

type Tick struct {
	ID           ResourceID    `json:"id,omitempty"`
	Ownership    OwnershipData `json:"-"`
	Timestamp    time.Time     `json:"timestamp"`
	ContestID    ResourceID    `json:"contestId"`
	ProblemID    ResourceID    `json:"problemId"`
	Top          bool          `json:"top"`
	AttemptsTop  int           `json:"attemptsTop"`
	Zone         bool          `json:"zone"`
	AttemptsZone int           `json:"attemptsZone"`
}

type TickUseCase interface {
	GetTicksByContender(ctx context.Context, contenderID ResourceID) ([]Tick, error)
	GetTicksByProblem(ctx context.Context, problemID ResourceID) ([]Tick, error)
	DeleteTick(ctx context.Context, tickID ResourceID) error
	CreateTick(ctx context.Context, contenderID ResourceID, tick Tick) (Tick, error)
}
