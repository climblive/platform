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
	TimeBegin   time.Time     `json:"timeBegin"`
	TimeEnd     time.Time     `json:"timeEnd"`
}

type CompClassTemplate struct {
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	TimeBegin   time.Time `json:"timeBegin"`
	TimeEnd     time.Time `json:"timeEnd"`
}

type CompClassPatch struct {
	Name        Patch[string]    `json:"name,omitempty" tstype:"string"`
	Description Patch[string]    `json:"description,omitempty" tstype:"string"`
	TimeBegin   Patch[time.Time] `json:"timeBegin,omitempty" tstype:"Date"`
	TimeEnd     Patch[time.Time] `json:"timeEnd,omitempty" tstype:"Date"`
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
	Entered             time.Time     `json:"entered,omitzero"`
	WithdrawnFromFinals bool          `json:"withdrawnFromFinals"`
	Disqualified        bool          `json:"disqualified"`
	Score               *Score        `json:"score,omitempty"`
}

type ContenderPatch struct {
	CompClassID         Patch[CompClassID] `json:"compClassId,omitzero" tstype:"CompClassID"`
	Name                Patch[string]      `json:"name,omitzero" tstype:"string"`
	PublicName          Patch[string]      `json:"publicName,omitzero" tstype:"string"`
	ClubName            Patch[string]      `json:"clubName,omitzero" tstype:"string"`
	WithdrawnFromFinals Patch[bool]        `json:"withdrawnFromFinals,omitzero" tstype:"boolean"`
	Disqualified        Patch[bool]        `json:"disqualified,omitzero" tstype:"boolean"`
}

type Contest struct {
	ID                 ContestID     `json:"id"`
	Ownership          OwnershipData `json:"-"`
	Location           string        `json:"location,omitempty"`
	SeriesID           SeriesID      `json:"seriesId,omitempty"`
	Name               string        `json:"name"`
	Description        string        `json:"description,omitempty"`
	QualifyingProblems int           `json:"qualifyingProblems"`
	Finalists          int           `json:"finalists"`
	Rules              string        `json:"rules,omitempty"`
	GracePeriod        time.Duration `json:"gracePeriod"`
	TimeBegin          time.Time     `json:"timeBegin,omitzero"`
	TimeEnd            time.Time     `json:"timeEnd,omitzero"`
}

type ContestTemplate struct {
	Location           string        `json:"location,omitempty"`
	SeriesID           SeriesID      `json:"seriesId,omitempty"`
	Name               string        `json:"name"`
	Description        string        `json:"description,omitempty"`
	QualifyingProblems int           `json:"qualifyingProblems"`
	Finalists          int           `json:"finalists"`
	Rules              string        `json:"rules,omitempty"`
	GracePeriod        time.Duration `json:"gracePeriod"`
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
	Description        string        `json:"description,omitempty"`
	PointsTop          int           `json:"pointsTop"`
	PointsZone         int           `json:"pointsZone"`
	FlashBonus         int           `json:"flashBonus,omitempty"`
}

type ProblemTemplate struct {
	Number             int    `json:"number"`
	HoldColorPrimary   string `json:"holdColorPrimary"`
	HoldColorSecondary string `json:"holdColorSecondary,omitempty"`
	Description        string `json:"description,omitempty"`
	PointsTop          int    `json:"pointsTop"`
	PointsZone         int    `json:"pointsZone"`
	FlashBonus         int    `json:"flashBonus,omitempty"`
}

type ProblemPatch struct {
	Number             Patch[int]    `json:"number,omitzero" tstype:"number"`
	HoldColorPrimary   Patch[string] `json:"holdColorPrimary,omitzero" tstype:"string"`
	HoldColorSecondary Patch[string] `json:"holdColorSecondary,omitzero" tstype:"string"`
	Description        Patch[string] `json:"description,omitzero" tstype:"string"`
	PointsTop          Patch[int]    `json:"pointsTop,omitzero" tstype:"number"`
	PointsZone         Patch[int]    `json:"pointsZone,omitzero" tstype:"number"`
	FlashBonus         Patch[int]    `json:"flashBonus,omitzero" tstype:"number"`
}

type Raffle struct {
	ID        RaffleID      `json:"id"`
	Ownership OwnershipData `json:"-"`
	ContestID ContestID     `json:"contestId"`
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
