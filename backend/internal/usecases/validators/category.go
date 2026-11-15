package validators

import (
	"strings"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

var errCategoryConstraintViolation = errors.New("constraint violation")

type CategoryValidator struct {
}

func (v CategoryValidator) Validate(category domain.Category) error {
	switch {
	case len(strings.TrimSpace(category.Name)) < 1:
		fallthrough
	case category.TimeEnd.Before(category.TimeBegin):
		fallthrough
	case category.TimeEnd.Sub(category.TimeBegin) > 31*24*time.Hour:
		return errors.Errorf("%w: %w", domain.ErrInvalidData, errCategoryConstraintViolation)
	}

	return nil
}

func (v CategoryValidator) IsValidationError(err error) bool {
	return errors.Is(err, errCategoryConstraintViolation)
}
