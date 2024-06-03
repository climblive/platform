package domain

import (
	"context"
	"time"
)

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
	GetContender(ctx context.Context, id ResourceID) (Contender, error)
	GetContenderByCode(ctx context.Context, registrationCode string) (Contender, error)
	GetContendersByContest(ctx context.Context, contestID ResourceID) ([]Contender, error)
	UpdateContender(ctx context.Context, id ResourceID, contender Contender) (Contender, error)
	DeleteContender(ctx context.Context, id ResourceID) error
	CreateContender(ctx context.Context, contestID ResourceID, template Contender) (Contender, error)
	AddTickets(ctx context.Context, contestID ResourceID, count int) ([]Contender, error)
	ExportContestResults(ctx context.Context, contestID ResourceID) ([]byte, error)
	ExportTicketsPDF(ctx context.Context, contestID ResourceID) ([]byte, error)
}
