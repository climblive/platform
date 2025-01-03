package domain

import (
	"context"
)

type RaffleWinnerUseCase interface {
	GetRaffleWinners(ctx context.Context, raffleID RaffleID) ([]RaffleWinner, error)
	DrawRaffleWinner(ctx context.Context, raffleID RaffleID) (RaffleWinner, error)
}
