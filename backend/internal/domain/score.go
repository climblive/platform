package domain

import "time"

type ScoreKeeper interface {
	UpdateScore(contenderID ResourceID, score Score) error
	GetScore(contenderID ResourceID) (Score, error)
}

type Score struct {
	Timestamp   *time.Time `json:"timestamp"`
	ContenderID ResourceID `json:"contenderId"`
	Score       int        `json:"score"`
	Placement   int        `json:"placement"`
	Finalist    bool       `json:"finalist"`
}
