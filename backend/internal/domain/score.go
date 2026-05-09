package domain

type ScoreKeeper interface {
	GetScore(contenderID ContenderID) (Score, error)
}
