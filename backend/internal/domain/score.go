package domain

import "time"

type ScoreKeeper interface {
	GetScore(contenderID ContenderID) (Score, error)
}

type Score struct {
	Timestamp   time.Time   `json:"timestamp"`
	ContenderID ContenderID `json:"contenderId"`
	Score       int         `json:"score"`
	Placement   int         `json:"placement"`
	Finalist    bool        `json:"finalist"`
	RankOrder   int         `json:"rankOrder"`
}
