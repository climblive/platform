package domain

import "time"

type EventBroker interface {
	Dispatch(contestID ResourceID, event any)
}

type ContenderEnterEvent struct {
	ContenderID ResourceID
	CompClassID ResourceID
}

type ContenderSwitchClassEvent struct {
	ContenderID ResourceID
	CompClassID ResourceID
}

type ContenderWithdrawFromFinalsEvent struct {
	ContenderID ResourceID
}

type ContenderReenterFinalsEvent struct {
	ContenderID ResourceID
}

type ContenderDisqualifyEvent struct {
	ContenderID ResourceID
}

type ContenderRequalifyEvent struct {
	ContenderID ResourceID
}

type RegisterAscentEvent struct {
	ContenderID  ResourceID
	ProblemID    ResourceID
	Top          bool
	AttemptsTop  int
	Zone         bool
	AttemptsZone int
}

type DeregisterAscentEvent struct {
	ContenderID ResourceID
	ProblemID   ResourceID
}

type AddProblemEvent struct {
	ProblemID  ResourceID
	PointsTop  int
	PointsZone int
	FlashBonus int
}

type UpdateProblemEvent struct {
	ProblemID  ResourceID
	PointsTop  int
	PointsZone int
	FlashBonus int
}

type DeleteProblemEvent struct {
	ProblemID ResourceID
}

type ContenderPublicInfoUpdateEvent struct {
	ContenderID         ResourceID
	CompClassID         ResourceID
	PublicName          string
	ClubName            string
	WithdrawnFromFinals bool
	Disqualified        bool
}

type ContenderScoreUpdateEvent struct {
	Timestamp   time.Time
	ContenderID ResourceID
	Score       int
	Placement   int
}
