package domain

import "context"

type Raffle struct {
	ID        ResourceID
	ContestID ResourceID
	Active    bool
}

type RaffleUsecase interface {
	GetRafflesByContest(ctx context.Context, contestID ResourceID) ([]Raffle, error)
	UpdateRaffle(ctx context.Context, raffleID ResourceID, raffle Raffle) (Raffle, error)
	DeleteRaffle(ctx context.Context, raffleID ResourceID) error
	CreateRaffle(ctx context.Context, contestID ResourceID, raffle Raffle) (Raffle, error)
}
