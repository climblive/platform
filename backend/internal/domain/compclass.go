package domain

import "time"

type ColorRGB string

type CompClass struct {
	ID          ResourceID
	ContestID   ResourceID
	Name        string
	Description string
	Color       ColorRGB
	TimeBegin   time.Time
	TimeEnd     time.Time
}

type CompClassUsecase interface {
}
