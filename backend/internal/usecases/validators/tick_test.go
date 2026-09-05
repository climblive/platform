package validators_test

import (
	"math"
	"testing"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/usecases/validators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			Revision:      1,
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

	t.Run("InvalidData", func(t *testing.T) {
		tests := map[string]func(*domain.Tick){
			"RevisionRequired": func(tick *domain.Tick) {
				tick.Revision = 0
			},
			"RevisionExceedsMaximum": func(tick *domain.Tick) {
				tick.Revision = math.MaxUint32 + 1
			},
			"AttemptsTopIsNegative": func(tick *domain.Tick) {
				tick.AttemptsTop = -1
			},
			"AttemptsTopExceedsMaximum": func(tick *domain.Tick) {
				tick.AttemptsTop = 1000
			},
			"AttemptsZone2IsNegative": func(tick *domain.Tick) {
				tick.AttemptsZone2 = -1
			},
			"AttemptsZone2ExceedsMaximum": func(tick *domain.Tick) {
				tick.AttemptsZone2 = 1000
			},
			"AttemptsZone1IsNegative": func(tick *domain.Tick) {
				tick.AttemptsZone1 = -1
			},
			"AttemptsZone1ExceedsMaximum": func(tick *domain.Tick) {
				tick.AttemptsZone1 = 1000
			},
			"AttemptsZone1ExceedsAttemptsZone2": func(tick *domain.Tick) {
				tick.Zone1 = true
				tick.AttemptsZone1 = 2
				tick.Zone2 = true
				tick.AttemptsZone2 = 1
				tick.Top = true
				tick.AttemptsTop = 3
			},
			"AttemptsZone2ExceedsAttemptsTop": func(tick *domain.Tick) {
				tick.Zone1 = true
				tick.AttemptsZone1 = 1
				tick.Zone2 = true
				tick.AttemptsZone2 = 3
				tick.Top = true
				tick.AttemptsTop = 2
			},
			"ReachedFeatureWithoutAttempts/Zone1": func(tick *domain.Tick) {
				tick.AttemptsTop = 0
				tick.AttemptsZone2 = 0
				tick.AttemptsZone1 = 0

				tick.Top = false
				tick.Zone2 = false
				tick.Zone1 = true
			},
			"ReachedFeatureWithoutAttempts/Zone2": func(tick *domain.Tick) {
				tick.AttemptsTop = 0
				tick.AttemptsZone2 = 0
				tick.AttemptsZone1 = 1

				tick.Top = false
				tick.Zone2 = true
				tick.Zone1 = true
			},
			"ReachedFeatureWithoutAttempts/Top": func(tick *domain.Tick) {
				tick.AttemptsTop = 0
				tick.AttemptsZone2 = 1
				tick.AttemptsZone1 = 1

				tick.Top = true
				tick.Zone2 = true
				tick.Zone1 = true
			},
			"TopWithoutZone1": func(tick *domain.Tick) {
				tick.AttemptsTop = 999
				tick.AttemptsZone2 = 999
				tick.AttemptsZone1 = 999

				tick.Top = true
				tick.Zone2 = true
				tick.Zone1 = false
			},
			"TopWithoutZone2": func(tick *domain.Tick) {
				tick.AttemptsTop = 999
				tick.AttemptsZone2 = 999
				tick.AttemptsZone1 = 999

				tick.Top = true
				tick.Zone2 = false
				tick.Zone1 = true
			},
			"Zone2WithoutZone1": func(tick *domain.Tick) {
				tick.AttemptsTop = 999
				tick.AttemptsZone2 = 999
				tick.AttemptsZone1 = 999

				tick.Top = false
				tick.Zone2 = true
				tick.Zone1 = false
			},
		}

		for name, mutate := range tests {
			t.Run(name, func(t *testing.T) {
				require.NotNil(t, mutate)

				tick := validTick()
				mutate(&tick)

				err := validator.Validate(tick)

				assert.ErrorIs(t, err, domain.ErrInvalidData)
				assert.True(t, validator.IsValidationError(err))
			})
		}
	})
}
