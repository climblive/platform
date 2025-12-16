package validators

import (
	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

var errTickConstraintViolation = errors.New("constraint violation")

type TickValidator struct {
}

func (v TickValidator) Validate(tick domain.Tick) error {
	switch {
	case tick.AttemptsZone1 > tick.AttemptsZone2:
		fallthrough
	case tick.AttemptsZone2 > tick.AttemptsTop:
		fallthrough
	case tick.Top && !tick.Zone1:
		fallthrough
	case tick.Top && !tick.Zone2:
		fallthrough
	case tick.Zone2 && !tick.Zone1:
		return errors.Errorf("%w: %w", domain.ErrInvalidData, errTickConstraintViolation)
	}

	return nil
}

func (v TickValidator) IsValidationError(err error) bool {
	return errors.Is(err, errTickConstraintViolation)
}
