package domain

import (
	"context"
	"time"
)

type RaffleWinner struct {
	ID            ResourceID
	RaffleID      ResourceID
	ContenderID   ResourceID
	ContenderName string
	Timestamp     time.Time
}

type RaffleWinnerUsecase interface {
	GetRaffleWinners(ctx context.Context, raffleID ResourceID) ([]RaffleWinner, error)
	DrawRaffleWinner(ctx context.Context, raffleID ResourceID) (RaffleWinner, error)
}
