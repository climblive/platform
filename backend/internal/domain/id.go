package domain

import "uuid"

type ResourceID int32

type CategoryID ResourceID
type ContenderID ResourceID
type CompetitionID ResourceID
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
		ContenderID |
		CompetitionID |
		OrganizerID |
		ProblemID |
		RaffleID |
		RaffleWinnerID |
		SeriesID |
		UserID |
		TickID
}

type ScoreEngineInstanceID = uuid.UUID
