package repository

type transaction struct {
}

func (tx *transaction) Commit() error {
	return nil
}

func (tx *transaction) Rollback() {

}
