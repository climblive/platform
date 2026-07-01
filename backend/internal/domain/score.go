package domain

type ScoreKeeper interface {
	GetScore(contenderID ContenderID) (Score, error)
}

type PointValueKeeper interface {
	GetPointValues(contenderID ContenderID) []PointValue
}
