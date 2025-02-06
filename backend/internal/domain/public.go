package domain

import (
	"time"
)

type ColorRGB string

type CompClass struct {
	ID          CompClassID   `json:"id"`
	Ownership   OwnershipData `json:"-"`
	ContestID   ContestID     `json:"contestId"`
	Name        string        `json:"name"`
	Description string        `json:"description,omitempty"`
	Color       ColorRGB      `json:"color,omitempty"`
	TimeBegin   time.Time     `json:"timeBegin"`
	TimeEnd     time.Time     `json:"timeEnd"`
}

type Contender struct {
	ID                  ContenderID   `json:"id"`
	Ownership           OwnershipData `json:"-"`
	ContestID           ContestID     `json:"contestId"`
	CompClassID         CompClassID   `json:"compClassId,omitempty"`
	RegistrationCode    string        `json:"registrationCode"`
	Name                string        `json:"name,omitempty"`
	PublicName          string        `json:"publicName,omitempty"`
	ClubName            string        `json:"clubName,omitempty"`
	Entered             *time.Time    `json:"entered,omitempty"`
	WithdrawnFromFinals bool          `json:"withdrawnFromFinals"`
	Disqualified        bool          `json:"disqualified"`
	Score               *Score        `json:"score,omitempty"`
}

type ContenderPatch struct {
	CompClassID         *Patch[CompClassID] `json:"compClassId" tstype:"CompClassID"`
	Name                *Patch[string]      `json:"name" tstype:"string"`
	PublicName          *Patch[string]      `json:"publicName" tstype:"string"`
	ClubName            *Patch[string]      `json:"clubName" tstype:"string"`
	WithdrawnFromFinals *Patch[bool]        `json:"withdrawnFromFinals" tstype:"boolean"`
	Disqualified        *Patch[bool]        `json:"disqualified" tstype:"boolean"`
}

type Contest struct {
	ID                 ContestID     `json:"id"`
	Ownership          OwnershipData `json:"-"`
	Location           string        `json:"location,omitempty"`
	SeriesID           SeriesID      `json:"seriesId,omitempty"`
	Protected          bool          `json:"protected"`
	Name               string        `json:"name"`
	Description        string        `json:"description,omitempty"`
	FinalsEnabled      bool          `json:"finalsEnabled"`
	QualifyingProblems int           `json:"qualifyingProblems"`
	Finalists          int           `json:"finalists"`
	Rules              string        `json:"rules,omitempty"`
	GracePeriod        time.Duration `json:"gracePeriod"`
	TimeBegin          *time.Time    `json:"timeBegin,omitempty"`
	TimeEnd            *time.Time    `json:"timeEnd,omitempty"`
}

type Organizer struct {
	ID        OrganizerID   `json:"id"`
	Ownership OwnershipData `json:"-"`
	Name      string        `json:"name"`
	Homepage  string        `json:"homepage,omitempty"`
}

type Problem struct {
	ID                 ProblemID     `json:"id"`
	Ownership          OwnershipData `json:"-"`
	ContestID          ContestID     `json:"contestId"`
	Number             int           `json:"number"`
	HoldColorPrimary   string        `json:"holdColorPrimary"`
	HoldColorSecondary string        `json:"holdColorSecondary,omitempty"`
	Name               string        `json:"name,omitempty"`
	Description        string        `json:"description,omitempty"`
	PointsTop          int           `json:"pointsTop"`
	PointsZone         int           `json:"pointsZone"`
	FlashBonus         int           `json:"flashBonus,omitempty"`
}

type Raffle struct {
	ID        RaffleID      `json:"id"`
	Ownership OwnershipData `json:"-"`
	ContestID ContestID     `json:"contestId"`
	Active    bool          `json:"active"`
}

type RaffleWinner struct {
	ID            RaffleWinnerID `json:"id"`
	Ownership     OwnershipData  `json:"-"`
	RaffleID      RaffleID       `json:"raffleId"`
	ContenderID   ContenderID    `json:"contenderId"`
	ContenderName string         `json:"contenderName" tstype:"string,readonly"`
	Timestamp     time.Time      `json:"timestamp"`
}

type Score struct {
	Timestamp   time.Time   `json:"timestamp"`
	ContenderID ContenderID `json:"contenderId"`
	Score       int         `json:"score"`
	Placement   int         `json:"placement"`
	Finalist    bool        `json:"finalist"`
	RankOrder   int         `json:"rankOrder"`
}

type Series struct {
	ID        SeriesID      `json:"id"`
	Ownership OwnershipData `json:"-"`
	Name      string        `json:"name"`
}

type ScoreboardEntry struct {
	ContenderID         ContenderID `json:"contenderId"`
	CompClassID         CompClassID `json:"compClassId"`
	PublicName          string      `json:"publicName"`
	ClubName            string      `json:"clubName,omitempty"`
	WithdrawnFromFinals bool        `json:"withdrawnFromFinals"`
	Disqualified        bool        `json:"disqualified"`
	Score               *Score      `json:"score,omitempty"`
}

type Tick struct {
	ID           TickID        `json:"id"`
	Ownership    OwnershipData `json:"-"`
	Timestamp    time.Time     `json:"timestamp"`
	ContestID    ContestID     `json:"-"`
	ProblemID    ProblemID     `json:"problemId"`
	Top          bool          `json:"top"`
	AttemptsTop  int           `json:"attemptsTop"`
	Zone         bool          `json:"zone"`
	AttemptsZone int           `json:"attemptsZone"`
}

type User struct {
	ID         UserID        `json:"id"`
	Name       string        `json:"name"`
	Username   string        `json:"username"`
	Admin      bool          `json:"admin"`
	Organizers []OrganizerID `json:"organizers"`
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
	TickID       TickID      `json:"tickId"`
	Timestamp    time.Time   `json:"timestamp"`
	ContenderID  ContenderID `json:"contenderId"`
	ProblemID    ProblemID   `json:"problemId"`
	Top          bool        `json:"top"`
	AttemptsTop  int         `json:"attemptsTop"`
	Zone         bool        `json:"zone"`
	AttemptsZone int         `json:"attemptsZone"`
}

type AscentDeregisteredEvent struct {
	TickID      TickID      `json:"tickId"`
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
	ClubName            string      `json:"clubName,omitempty"`
	WithdrawnFromFinals bool        `json:"withdrawnFromFinals"`
	Disqualified        bool        `json:"disqualified"`
}

type ContenderScoreUpdatedEvent struct {
	Timestamp   time.Time   `json:"timestamp"`
	ContenderID ContenderID `json:"contenderId"`
	Score       int         `json:"score"`
	Placement   int         `json:"placement"`
	Finalist    bool        `json:"finalist"`
	RankOrder   int         `json:"rankOrder"`
}

type ScoreEngineStartedEvent struct {
	InstanceID ScoreEngineInstanceID `json:"instanceId"`
}

type ScoreEngineStoppedEvent struct {
	InstanceID ScoreEngineInstanceID `json:"instanceId"`
}
