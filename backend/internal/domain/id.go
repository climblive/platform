package domain

type resourceID int

type CompClassID resourceID
type ContenderID resourceID
type ContestID resourceID
type OrganizerID resourceID
type ProblemID resourceID
type RaffleID resourceID
type RaffleWinnerID resourceID
type SeriesID resourceID
type UserID resourceID
type TickID resourceID

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
