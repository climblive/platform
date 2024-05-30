package domain

import "time"

type Tick struct {
	ID          ResourceID
	Timestamp   time.Time
	ContenderID ResourceID
	ProblemID   ResourceID
	Flash       bool
}

type TickUsecase interface {
}
