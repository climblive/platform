package domain

import "time"

type ScoreboardEntry struct {
	ContenderID         ContenderID `json:"contenderId"`
	CompClassID         CompClassID `json:"compClassId"`
	PublicName          string      `json:"publicName,omitempty"`
	ClubName            string      `json:"clubName,omitempty"`
	WithdrawnFromFinals bool        `json:"withdrawnFromFinals"`
	Disqualified        bool        `json:"disqualified"`
	ScoreUpdated        *time.Time  `json:"scoreUpdated,omitempty"`
	Score               int         `json:"score"`
	Placement           int         `json:"placement,omitempty"`
	Finalist            bool        `json:"finalist"`
	RankOrder           int         `json:"rankOrder"`
}
