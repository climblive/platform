package domain

import "time"

type RaffleWinner struct {
	ID            ResourceID
	RaffleID      ResourceID
	ContenderID   ResourceID
	ContenderName string
	Timestamp     time.Time
}

type RaffleWinnerUsecase interface {
}
