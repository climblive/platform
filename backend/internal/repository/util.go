package repository

import (
	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

func nillableUUIDToResourceID[T domain.ResourceIDType](value uuid.UUID) *T {
	if value == uuid.Nil {
		return nil
	}

	out := T(value)
	return &out
}

var mysqlForeignKeyConstraintViolation = mysql.MySQLError{Number: 1452}
var mysqlDuplicateKeyConstraintViolation = mysql.MySQLError{Number: 1062}
