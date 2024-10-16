package domain

import "time"

type ScoreKeeper interface {
	GetScore(contenderID ResourceID) (Score, error)
}

type Score struct {
	Timestamp   time.Time  `json:"timestamp"`
	ContenderID ResourceID `json:"contenderId"`
	Score       int        `json:"score"`
	Placement   int        `json:"placement,omitempty"`
	Finalist    bool       `json:"finalist"`
	RankOrder   int        `json:"rankOrder"`
}
