package repository

import (
	"github.com/climblive/platform/backend/internal/domain"
)

type Database struct {
}

func NewDatabase() *Database {
	return &Database{}
}

func (d *Database) Begin() domain.Transaction {
	return &transaction{}
}
