package utils

import "github.com/go-errors/errors"

func GetErrorStack(err error) string {
	if err, ok := err.(*errors.Error); ok {
		return err.ErrorStack()
	}

	return ""
}
