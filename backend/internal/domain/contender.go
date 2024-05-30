package domain

import "time"

type Contender struct {
	ID               ResourceID
	CompClassID      ResourceID
	ContestID        ResourceID
	RegistrationCode string
	Name             string
	Club             string
	Entered          time.Time
	Disqualified     bool
}

type ContenderUsecase interface {
}
