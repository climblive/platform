package validators

import (
	"regexp"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

var (
	errProblemConstraintViolation = errors.New("constraint violation")
	validHexColor                 = regexp.MustCompile(`^#([0-9a-fA-F]{3}){1,2}$`)
	maxAllowedPointValue          = 2_147_483_647
)

type ProblemValidator struct {
}

func (v ProblemValidator) Validate(problem domain.Problem) error {
	switch {
	case problem.Number < 0:
		fallthrough
	case !validHexColor.MatchString(problem.HoldColorPrimary):
		fallthrough
	case len(problem.HoldColorSecondary) > 0 && !validHexColor.MatchString(problem.HoldColorSecondary):
		fallthrough
	case problem.PointsTop < 0:
		fallthrough
	case problem.FlashBonus < 0:
		fallthrough
	case problem.PointsZone1 < 0:
		fallthrough
	case problem.PointsZone2 < 0:
		fallthrough
	case problem.PointsTop > maxAllowedPointValue:
		fallthrough
	case problem.FlashBonus > maxAllowedPointValue:
		fallthrough
	case problem.PointsZone1 > maxAllowedPointValue:
		fallthrough
	case problem.PointsZone2 > maxAllowedPointValue:
		fallthrough
	case problem.Zone2Enabled && !problem.Zone1Enabled:
		return errors.Errorf("%w: %w", domain.ErrInvalidData, errProblemConstraintViolation)
	}

	return nil
}

func (v ProblemValidator) IsValidationError(err error) bool {
	return errors.Is(err, errProblemConstraintViolation)
}
