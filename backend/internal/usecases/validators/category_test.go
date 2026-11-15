package validators_test

import (
	"testing"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/usecases/validators"
	"github.com/stretchr/testify/assert"
)

func TestCategoryValidator(t *testing.T) {
	now := time.Now()

	validator := validators.CategoryValidator{}

	validCategory := func() domain.Category {
		return domain.Category{
			Name:        "Females",
			Description: "Female climbers",
			TimeBegin:   now,
			TimeEnd:     now.Add(time.Hour),
		}
	}

	t.Run("ValidData", func(t *testing.T) {
		err := validator.Validate(validCategory())
		assert.NoError(t, err)
	})

	t.Run("EmptyName", func(t *testing.T) {
		category := validCategory()
		category.Name = whitespaceCharacters

		err := validator.Validate(category)

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.True(t, validator.IsValidationError(err))
	})

	t.Run("TimeEndBeforeTimeBegin", func(t *testing.T) {
		category := validCategory()
		category.TimeEnd = category.TimeBegin.Add(-time.Nanosecond)

		err := validator.Validate(category)

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.True(t, validator.IsValidationError(err))
	})

	t.Run("TotalDurationExceeding31Days", func(t *testing.T) {
		category := validCategory()
		category.TimeEnd = category.TimeBegin.Add(31*24*time.Hour + time.Nanosecond)

		err := validator.Validate(category)

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.True(t, validator.IsValidationError(err))
	})
}
