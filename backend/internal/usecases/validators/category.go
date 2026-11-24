package validators

import (
	"math"
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
	case !validCategoryEligibility(category):
		fallthrough
	case len(strings.TrimSpace(category.Name)) < 1:
		fallthrough
	case category.TimeEnd.Before(category.TimeBegin):
		fallthrough
	case category.TimeEnd.Sub(category.TimeBegin) > 31*24*time.Hour:
		return errors.Errorf("%w: %w", domain.ErrInvalidData, errCategoryConstraintViolation)
	}

	return nil
}

func validCategoryEligibility(category domain.Category) bool {
	switch category.EligibilityBasis {
	case domain.EligibilityBasisAge:
		return category.MinBirthYear == nil && category.MaxBirthYear == nil && validCategoryRange(category.MinAge, category.MaxAge, 0, math.MaxInt32)
	case domain.EligibilityBasisBirthYear:
		return category.MinAge == nil && category.MaxAge == nil && validCategoryRange(category.MinBirthYear, category.MaxBirthYear, 1, 9999)
	default:
		return false
	}
}

func validCategoryRange(minimum, maximum *int, lowerLimit, upperLimit int) bool {
	if minimum != nil && (*minimum < lowerLimit || *minimum > upperLimit) {
		return false
	}
	if maximum != nil && (*maximum < lowerLimit || *maximum > upperLimit) {
		return false
	}
	return minimum == nil || maximum == nil || *minimum <= *maximum
}

func (v CategoryValidator) IsValidationError(err error) bool {
	return errors.Is(err, errCategoryConstraintViolation)
}
