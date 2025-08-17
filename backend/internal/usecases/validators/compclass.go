package validators

import (
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

var errCompClassConstraintViolation = errors.New("constraint violation")

type CompClassValidator struct {
}

func (v CompClassValidator) Validate(compClass domain.CompClass) error {
	switch {
	case len(compClass.Name) < 1:
		fallthrough
	case compClass.TimeEnd.Before(compClass.TimeBegin):
		fallthrough
	case compClass.TimeEnd.Sub(compClass.TimeBegin) > 31*24*time.Hour:
		return errors.Errorf("%w: %w", domain.ErrInvalidData, errCompClassConstraintViolation)
	}

	return nil
}

func (v CompClassValidator) IsValidationError(err error) bool {
	return errors.Is(err, errCompClassConstraintViolation)
}
