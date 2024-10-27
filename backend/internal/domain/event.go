package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type SubscriptionID = uuid.UUID

type EventFilter struct {
	ContestID   ContestID
	ContenderID ContenderID
}

type EventBroker interface {
	Dispatch(contestID ContestID, event any)
	Subscribe(filter EventFilter, bufferCapacity int) (SubscriptionID, EventReader)
	Unsubscribe(subscriptionID SubscriptionID)
}

type EventReader interface {
	AwaitEvent(ctx context.Context) (EventEnvelope, error)
}

type EventEnvelope struct {
	Name string
	Data any
}

type ContenderEnteredEvent struct {
	ContenderID ContenderID `json:"contenderId"`
	CompClassID CompClassID `json:"compClassId"`
}

type ContenderSwitchedClassEvent struct {
	ContenderID ContenderID `json:"contenderId"`
	CompClassID CompClassID `json:"compClassId"`
}

type ContenderWithdrewFromFinalsEvent struct {
	ContenderID ContenderID `json:"contenderId"`
}

type ContenderReenteredFinalsEvent struct {
	ContenderID ContenderID `json:"contenderId"`
}

type ContenderDisqualifiedEvent struct {
	ContenderID ContenderID `json:"contenderId"`
}

type ContenderRequalifiedEvent struct {
	ContenderID ContenderID `json:"contenderId"`
}

type AscentRegisteredEvent struct {
	ContenderID  ContenderID `json:"contenderId"`
	ProblemID    ProblemID   `json:"problemId"`
	Top          bool        `json:"top"`
	AttemptsTop  int         `json:"attemptsTop"`
	Zone         bool        `json:"zone"`
	AttemptsZone int         `json:"attemptsZone"`
}

type AscentDeregisteredEvent struct {
	ContenderID ContenderID `json:"contenderId"`
	ProblemID   ProblemID   `json:"problemId"`
}

type ProblemAddedEvent struct {
	ProblemID  ProblemID `json:"problemId"`
	PointsTop  int       `json:"pointsTop"`
	PointsZone int       `json:"pointsZone"`
	FlashBonus int       `json:"flashBonus"`
}

type ProblemUpdatedEvent struct {
	ProblemID  ProblemID `json:"problemId"`
	PointsTop  int       `json:"pointsTop"`
	PointsZone int       `json:"pointsZone"`
	FlashBonus int       `json:"flashBonus"`
}

type ProblemDeletedEvent struct {
	ProblemID ProblemID `json:"problemId"`
}

type ContenderPublicInfoUpdatedEvent struct {
	ContenderID         ContenderID `json:"contenderId"`
	CompClassID         CompClassID `json:"compClassId"`
	PublicName          string      `json:"publicName"`
	ClubName            string      `json:"clubName"`
	WithdrawnFromFinals bool        `json:"withdrawnFromFinals"`
	Disqualified        bool        `json:"disqualified"`
}

type ContenderScoreUpdatedEvent struct {
	Timestamp   time.Time   `json:"timestamp"`
	ContenderID ContenderID `json:"contenderId"`
	Score       int         `json:"score"`
	Placement   int         `json:"placement,omitempty"`
	Finalist    bool        `json:"finalist"`
	RankOrder   int         `json:"rankOrder"`
}
