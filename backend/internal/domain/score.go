package domain

import "time"

type ScoreKeeper interface {
	GetScore(contenderID ResourceID) (Score, error)
}

type Score struct {
	Timestamp   *time.Time `json:"timestamp"`
	ContenderID ResourceID `json:"contenderId"`
	Score       int        `json:"score,omitempty"`
	Placement   int        `json:"placement"`
	Finalist    bool       `json:"finalist"`
	RankOrder   int        `json:"rankOrder"`
}
