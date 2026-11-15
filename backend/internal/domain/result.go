package domain

type ResultKeeper interface {
	GetResult(competitorID CompetitorID) (Result, error)
}
