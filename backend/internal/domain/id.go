package domain

import "github.com/google/uuid"

type ResourceID int32

type CompClassID ResourceID
type ContenderID ResourceID
type ContestID ResourceID
type OrganizerID ResourceID
type ProblemID ResourceID
type RaffleID ResourceID
type RaffleWinnerID ResourceID
type SeriesID ResourceID
type UserID ResourceID
type TickID ResourceID

type OrganizerInviteID uuid.UUID

func (id OrganizerInviteID) String() string {
	return uuid.UUID(id).String()
}

type ResourceIDType interface {
	CompClassID |
		ContenderID |
		ContestID |
		OrganizerID |
		ProblemID |
		RaffleID |
		RaffleWinnerID |
		SeriesID |
		UserID |
		TickID
}

type ScoreEngineInstanceID = uuid.UUID
