package domain

type ScoreboardEntry struct {
	ContenderID         ContenderID `json:"contenderId"`
	CompClassID         CompClassID `json:"compClassId"`
	PublicName          string      `json:"publicName,omitempty"`
	ClubName            string      `json:"clubName,omitempty"`
	WithdrawnFromFinals bool        `json:"withdrawnFromFinals"`
	Disqualified        bool        `json:"disqualified"`
	Score               *Score      `json:"score,omitempty"`
}
