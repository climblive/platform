package domain

type ResultKeeper interface {
	GetResult(contenderID ContenderID) (Result, error)
}
