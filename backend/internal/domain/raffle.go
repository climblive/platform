package domain

import "context"

type Raffle struct {
	ID        RaffleID
	Ownership OwnershipData
	ContestID ContestID
	Active    bool
}

type RaffleUseCase interface {
	GetRafflesByContest(ctx context.Context, contestID ContestID) ([]Raffle, error)
	UpdateRaffle(ctx context.Context, raffleID RaffleID, raffle Raffle) (Raffle, error)
	DeleteRaffle(ctx context.Context, raffleID RaffleID) error
	CreateRaffle(ctx context.Context, contestID ContestID, raffle Raffle) (Raffle, error)
}
