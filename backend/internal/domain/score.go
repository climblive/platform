package domain

import "time"

type Score struct {
	Timestamp   time.Time
	ContenderID ResourceID
	Score       int
	Placement   int
}
