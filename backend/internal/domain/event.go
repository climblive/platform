package domain

import (
	"time"

	"github.com/google/uuid"
)

type SubscriptionID = uuid.UUID

type EventBroker interface {
	Dispatch(contestID ResourceID, event any)
	Subscribe(contestID ResourceID, contenderID *ResourceID, ch chan EventEnclosure) SubscriptionID
	Unsubscribe(subscriptionID SubscriptionID)
}

type EventEnclosure struct {
	Name string
	Data any
}

type ContenderEnteredEvent struct {
	ContenderID ResourceID `json:"contenderId"`
	CompClassID ResourceID `json:"compClassId"`
}

type ContenderSwitchedClassEvent struct {
	ContenderID ResourceID `json:"contenderId"`
	CompClassID ResourceID `json:"compClassId"`
}

type ContenderWithdrewFromFinalsEvent struct {
	ContenderID ResourceID `json:"contenderId"`
}

type ContenderReenteredFinalsEvent struct {
	ContenderID ResourceID `json:"contenderId"`
}

type ContenderDisqualifiedEvent struct {
	ContenderID ResourceID `json:"contenderId"`
}

type ContenderRequalifiedEvent struct {
	ContenderID ResourceID `json:"contenderId"`
}

type AscentRegisteredEvent struct {
	ContenderID  ResourceID `json:"contenderId"`
	ProblemID    ResourceID `json:"problemId"`
	Top          bool
	AttemptsTop  int
	Zone         bool
	AttemptsZone int
}

type AscentDeregisteredEvent struct {
	ContenderID ResourceID `json:"contenderId"`
	ProblemID   ResourceID `json:"problemId"`
}

type ProblemAddedEvent struct {
	ProblemID  ResourceID `json:"problemId"`
	PointsTop  int        `json:"pointsTop"`
	PointsZone int        `json:"pointsZone"`
	FlashBonus int        `json:"flashBonus"`
}

type ProblemUpdatedEvent struct {
	ProblemID  ResourceID `json:"problemId"`
	PointsTop  int        `json:"pointsTop"`
	PointsZone int        `json:"pointsZone"`
	FlashBonus int        `json:"flashBonus"`
}

type ProblemDeletedEvent struct {
	ProblemID ResourceID `json:"problemId"`
}

type ContenderPublicInfoUpdatedEvent struct {
	ContenderID         ResourceID `json:"contenderId"`
	CompClassID         ResourceID `json:"compClassId"`
	PublicName          string     `json:"publicName"`
	ClubName            string     `json:"clubName"`
	WithdrawnFromFinals bool       `json:"withdrawnFromFinals"`
	Disqualified        bool       `json:"disqualified"`
}

type ContenderScoreUpdatedEvent struct {
	Timestamp   time.Time  `json:"timestamp"`
	ContenderID ResourceID `json:"contenderId"`
	Score       int        `json:"score"`
	Placement   int        `json:"placement"`
}
