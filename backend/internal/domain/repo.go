package domain

type TransactionController interface {
	Commit() error
	Rollback()
}

type Transactor[R any] interface {
	Autocommit() R
	Begin() (TransactionController, R)
}
