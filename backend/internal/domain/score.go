package domain

import "github.com/google/uuid"

type ScoreKeeper interface {
	GetScore(contenderID ContenderID) (Score, error)
}

type ScoreEngineInstanceID = uuid.UUID
