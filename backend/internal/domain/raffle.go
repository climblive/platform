package domain

import "context"

type Raffle struct {
	ID        ResourceID
	ContestID ResourceID
	Active    bool
}

type RaffleUsecase interface {
	GetRaffle(ctx context.Context, id ResourceID) (Raffle, error)
	UpdateRaffle(ctx context.Context, id ResourceID, raffle Raffle) (Raffle, error)
	DeleteRaffle(ctx context.Context, id ResourceID) (error)
	GetRafflesByContest(ctx context.Context, contestID ResourceID) ([]Raffle, error)
	CreateRaffle(ctx context.Context, contestID ResourceID, raffle Raffle) (Raffle, error)
}
