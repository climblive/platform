package domain

type Raffle struct {
	ID        ResourceID
	ContestID ResourceID
	Active    bool
}

type RaffleUsecase interface {
}
