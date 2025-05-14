package validators_test

import (
	"testing"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/usecases/validators"
	"github.com/stretchr/testify/assert"
)

func TestCompClassValidator(t *testing.T) {
	now := time.Now()

	validator := validators.CompClassValidator{}

	validCompClass := func() domain.CompClass {
		return domain.CompClass{
			Name:        "Females",
			Description: "Female climbers",
			TimeBegin:   now,
			TimeEnd:     now.Add(time.Hour),
		}
	}

	t.Run("ValidData", func(t *testing.T) {
		err := validator.Validate(validCompClass())
		assert.NoError(t, err)
	})

	t.Run("EmptyName", func(t *testing.T) {
		compClass := validCompClass()
		compClass.Name = ""

		err := validator.Validate(compClass)

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.True(t, validator.IsValidationError(err))
	})

	t.Run("TimeEndBeforeTimeBegin", func(t *testing.T) {
		compClass := validCompClass()
		compClass.TimeEnd = compClass.TimeBegin.Add(-time.Nanosecond)

		err := validator.Validate(compClass)

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.True(t, validator.IsValidationError(err))
	})

	t.Run("TotalDurationExceedingTwelveHours", func(t *testing.T) {
		compClass := validCompClass()
		compClass.TimeEnd = compClass.TimeBegin.Add(12*time.Hour + time.Nanosecond)

		err := validator.Validate(compClass)

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.True(t, validator.IsValidationError(err))
	})
}
