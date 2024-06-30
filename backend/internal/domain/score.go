package domain

import "time"

type ScoreKeeper interface {
	UpdateScore(contenderID ResourceID, score Score) error
	GetScore(contenderID ResourceID) (Score, error)
}

type Score struct {
	Timestamp   time.Time
	ContenderID ResourceID
	Score       int
	Placement   int
}
