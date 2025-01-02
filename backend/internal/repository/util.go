package repository

import (
	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-sql-driver/mysql"
)

func nillableIntToResourceID[T domain.ResourceIDType](value *int32) *T {
	if value == nil {
		return nil
	}

	var out T = T(*value)
	return &out
}

var mysqlForeignKeyConstraintViolation = mysql.MySQLError{Number: 1452}
