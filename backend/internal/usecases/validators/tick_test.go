package validators_test

import (
	"testing"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/usecases/validators"
	"github.com/stretchr/testify/assert"
)

func TestTickValidator(t *testing.T) {
	validator := validators.TickValidator{}

	validTick := func() domain.Tick {
		return domain.Tick{
			Timestamp: time.Now(),
			ID:        domain.TickID(1),
			Ownership: domain.OwnershipData{
				OrganizerID: domain.OrganizerID(1),
			},
			ContestID:     domain.ContestID(1),
			ProblemID:     domain.ProblemID(1),
			Zone1:         true,
			AttemptsZone1: 10,
			Zone2:         true,
			AttemptsZone2: 20,
			Top:           true,
			AttemptsTop:   30,
		}
	}

	t.Run("ValidData", func(t *testing.T) {
		err := validator.Validate(validTick())
		assert.NoError(t, err)
	})

	t.Run("AttemptsZone1ExceedsAttemptsZone2", func(t *testing.T) {
		tick := validTick()
		tick.AttemptsZone1 = 21

		err := validator.Validate(tick)

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.True(t, validator.IsValidationError(err))
	})

	t.Run("AttemptsZone2ExceedsAttemptsTop", func(t *testing.T) {
		tick := validTick()
		tick.AttemptsZone2 = 31

		err := validator.Validate(tick)

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.True(t, validator.IsValidationError(err))
	})

	t.Run("TopWithoutZones", func(t *testing.T) {
		tick := validTick()
		tick.Top = true
		tick.Zone2 = false
		tick.Zone1 = false

		err := validator.Validate(tick)

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.True(t, validator.IsValidationError(err))
	})

	t.Run("Zone2WithoutZone1", func(t *testing.T) {
		tick := validTick()
		tick.Top = false
		tick.Zone2 = true
		tick.Zone1 = false

		err := validator.Validate(tick)

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.True(t, validator.IsValidationError(err))
	})
}
