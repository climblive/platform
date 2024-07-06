package domain

type Transaction interface {
	Commit() error
	Rollback()
}

type Transactor interface {
	Begin() Transaction
}
