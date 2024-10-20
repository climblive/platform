package domain

import (
	"context"
	"time"
)

type RaffleWinner struct {
	ID            RaffleWinnerID
	Ownership     OwnershipData
	RaffleID      RaffleID
	ContenderID   ContenderID
	ContenderName string
	Timestamp     time.Time
}

type RaffleWinnerUseCase interface {
	GetRaffleWinners(ctx context.Context, raffleID RaffleID) ([]RaffleWinner, error)
	DrawRaffleWinner(ctx context.Context, raffleID RaffleID) (RaffleWinner, error)
}
