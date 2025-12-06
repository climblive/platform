package domain

import (
	"time"
)

type ColorRGB string

type OwnershipData struct {
	OrganizerID OrganizerID  `json:"organizerId"`
	ContenderID *ContenderID `json:"-"`
}

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
	Entered             time.Time     `json:"entered,omitzero"`
	WithdrawnFromFinals bool          `json:"withdrawnFromFinals"`
	Disqualified        bool          `json:"disqualified"`
	Score               *Score        `json:"score,omitempty"`
}

type ContenderPatch struct {
	CompClassID         Patch[CompClassID] `json:"compClassId,omitzero" tstype:"CompClassID"`
	Name                Patch[string]      `json:"name,omitzero" tstype:"string"`
	WithdrawnFromFinals Patch[bool]        `json:"withdrawnFromFinals,omitzero" tstype:"boolean"`
	Disqualified        Patch[bool]        `json:"disqualified,omitzero" tstype:"boolean"`
}

type Contest struct {
	ID                 ContestID     `json:"id"`
	Ownership          OwnershipData `json:"ownership"`
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

type ContestPatch struct {
	Location           Patch[string]        `json:"location,omitzero" tstype:"string"`
	SeriesID           Patch[SeriesID]      `json:"seriesId,omitzero" tstype:"number"`
	Name               Patch[string]        `json:"name,omitzero" tstype:"string"`
	Description        Patch[string]        `json:"description,omitzero" tstype:"string"`
	QualifyingProblems Patch[int]           `json:"qualifyingProblems,omitzero" tstype:"number"`
	Finalists          Patch[int]           `json:"finalists,omitzero" tstype:"number"`
	Rules              Patch[string]        `json:"rules,omitzero" tstype:"string"`
	GracePeriod        Patch[time.Duration] `json:"gracePeriod,omitzero" tstype:"number"`
}

type Organizer struct {
	ID        OrganizerID   `json:"id"`
	Ownership OwnershipData `json:"-"`
	Name      string        `json:"name"`
}

type Problem struct {
	ID                 ProblemID     `json:"id"`
	Ownership          OwnershipData `json:"-"`
	ContestID          ContestID     `json:"contestId"`
	Number             int           `json:"number"`
	HoldColorPrimary   string        `json:"holdColorPrimary"`
	HoldColorSecondary string        `json:"holdColorSecondary,omitempty"`
	Description        string        `json:"description,omitempty"`
	Zone1Enabled       bool          `json:"zone1Enabled"`
	Zone2Enabled       bool          `json:"zone2Enabled"`
	PointsZone1        int           `json:"pointsZone1"`
	PointsZone2        int           `json:"pointsZone2"`
	PointsTop          int           `json:"pointsTop"`
	FlashBonus         int           `json:"flashBonus,omitempty"`
}

type ProblemTemplate struct {
	Number             int    `json:"number"`
	HoldColorPrimary   string `json:"holdColorPrimary"`
	HoldColorSecondary string `json:"holdColorSecondary,omitempty"`
	Description        string `json:"description,omitempty"`
	Zone1Enabled       bool   `json:"zone1Enabled"`
	Zone2Enabled       bool   `json:"zone2Enabled"`
	PointsZone1        int    `json:"pointsZone1"`
	PointsZone2        int    `json:"pointsZone2"`
	PointsTop          int    `json:"pointsTop"`
	FlashBonus         int    `json:"flashBonus,omitempty"`
}

type ProblemPatch struct {
	Number             Patch[int]    `json:"number,omitzero" tstype:"number"`
	HoldColorPrimary   Patch[string] `json:"holdColorPrimary,omitzero" tstype:"string"`
	HoldColorSecondary Patch[string] `json:"holdColorSecondary,omitzero" tstype:"string"`
	Description        Patch[string] `json:"description,omitzero" tstype:"string"`
	Zone1Enabled       Patch[bool]   `json:"zone1Enabled,omitzero" tstype:"boolean"`
	Zone2Enabled       Patch[bool]   `json:"zone2Enabled,omitzero" tstype:"boolean"`
	PointsZone1        Patch[int]    `json:"pointsZone1,omitzero" tstype:"number"`
	PointsZone2        Patch[int]    `json:"pointsZone2,omitzero" tstype:"number"`
	PointsTop          Patch[int]    `json:"pointsTop,omitzero" tstype:"number"`
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
	Name                string      `json:"name"`
	WithdrawnFromFinals bool        `json:"withdrawnFromFinals"`
	Disqualified        bool        `json:"disqualified"`
	Score               *Score      `json:"score,omitempty"`
}

type Tick struct {
	ID            TickID        `json:"id"`
	Ownership     OwnershipData `json:"-"`
	Timestamp     time.Time     `json:"timestamp"`
	ContestID     ContestID     `json:"-"`
	ProblemID     ProblemID     `json:"problemId"`
	Zone1         bool          `json:"zone1"`
	AttemptsZone1 int           `json:"attemptsZone1"`
	Zone2         bool          `json:"zone2"`
	AttemptsZone2 int           `json:"attemptsZone2"`
	Top           bool          `json:"top"`
	AttemptsTop   int           `json:"attemptsTop"`
}

type User struct {
	ID         UserID      `json:"id"`
	Username   string      `json:"username"`
	Admin      bool        `json:"admin"`
	Organizers []Organizer `json:"organizers"`
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
	TickID        TickID      `json:"tickId"`
	Timestamp     time.Time   `json:"timestamp"`
	ContenderID   ContenderID `json:"contenderId"`
	ProblemID     ProblemID   `json:"problemId"`
	Zone1         bool        `json:"zone1"`
	AttemptsZone1 int         `json:"attemptsZone1"`
	Zone2         bool        `json:"zone2"`
	AttemptsZone2 int         `json:"attemptsZone2"`
	Top           bool        `json:"top"`
	AttemptsTop   int         `json:"attemptsTop"`
}

type AscentDeregisteredEvent struct {
	TickID      TickID      `json:"tickId"`
	ContenderID ContenderID `json:"contenderId"`
	ProblemID   ProblemID   `json:"problemId"`
}

type ProblemAddedEvent struct {
	ProblemID   ProblemID `json:"problemId"`
	PointsZone1 int       `json:"pointsZone1"`
	PointsZone2 int       `json:"pointsZone2"`
	PointsTop   int       `json:"pointsTop"`
	FlashBonus  int       `json:"flashBonus"`
}

type ProblemUpdatedEvent struct {
	ProblemID   ProblemID `json:"problemId"`
	PointsZone1 int       `json:"pointsZone1"`
	PointsZone2 int       `json:"pointsZone2"`
	PointsTop   int       `json:"pointsTop"`
	FlashBonus  int       `json:"flashBonus"`
}

type ProblemDeletedEvent struct {
	ProblemID ProblemID `json:"problemId"`
}

type ContenderPublicInfoUpdatedEvent struct {
	ContenderID         ContenderID `json:"contenderId"`
	CompClassID         CompClassID `json:"compClassId"`
	Name                string      `json:"name"`
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
