package domain

import "time"

type ScoreboardEntry struct {
	ContenderID         ResourceID `json:"contenderId"`
	CompClassID         ResourceID `json:"compClassId"`
	PublicName          string     `json:"publicName"`
	ClubName            string     `json:"clubName"`
	WithdrawnFromFinals bool       `json:"withdrawnFromFinals"`
	Disqualified        bool       `json:"disqualified"`
	ScoreUpdated        *time.Time `json:"scoreUpdated,omitempty"`
	Score               int        `json:"score"`
	Placement           int        `json:"placement,omitempty"`
	Finalist            bool       `json:"finalist"`
	RankOrder           int        `json:"rankOrder"`
}
