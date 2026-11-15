package domain

import "uuid"

type ResourceID int32

type CategoryID ResourceID
type CompetitorID ResourceID
type CompetitionID uuid.UUID
type OrganizerID ResourceID
type ProblemID ResourceID
type RaffleID ResourceID
type RaffleWinnerID ResourceID
type SeriesID ResourceID
type UserID ResourceID
type TickID ResourceID

type OrganizerInviteID = uuid.UUID

type ResourceIDType interface {
	CategoryID |
		CompetitorID |
		OrganizerID |
		ProblemID |
		RaffleID |
		RaffleWinnerID |
		SeriesID |
		UserID |
		TickID
}

type ScoreEngineInstanceID = uuid.UUID
