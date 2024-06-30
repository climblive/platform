package domain

import (
	"context"
	"time"
)

type RaffleWinner struct {
	ID            ResourceID
	Ownership     OwnershipData
	RaffleID      ResourceID
	ContenderID   ResourceID
	ContenderName string
	Timestamp     time.Time
}

type RaffleWinnerUseCase interface {
	GetRaffleWinners(ctx context.Context, raffleID ResourceID) ([]RaffleWinner, error)
	DrawRaffleWinner(ctx context.Context, raffleID ResourceID) (RaffleWinner, error)
}
