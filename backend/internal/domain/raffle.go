package domain

import "context"

type Raffle struct {
	ID        ResourceID
	Ownership OwnershipData
	ContestID ResourceID
	Active    bool
}

type RaffleUseCase interface {
	GetRafflesByContest(ctx context.Context, contestID ResourceID) ([]Raffle, error)
	UpdateRaffle(ctx context.Context, raffleID ResourceID, raffle Raffle) (Raffle, error)
	DeleteRaffle(ctx context.Context, raffleID ResourceID) error
	CreateRaffle(ctx context.Context, contestID ResourceID, raffle Raffle) (Raffle, error)
}
