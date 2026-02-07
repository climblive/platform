package domain

type ScoreKeeper interface {
	GetScore(contenderID ContenderID) (Score, error)
}

type ProblemValueKeeper interface {
	GetProblemValue(problemID ProblemID, compClassID CompClassID) (ProblemValue, bool)
}
