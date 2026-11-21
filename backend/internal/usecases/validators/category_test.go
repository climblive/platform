package validators_test

import (
	"math"
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

	for _, tc := range []struct {
		name   string
		minAge int
		maxAge int
		valid  bool
	}{
		{name: "adult span", minAge: 18, maxAge: 99, valid: true},
		{name: "single age", minAge: 16, maxAge: 16, valid: true},
		{name: "age zero", minAge: 0, maxAge: 15, valid: true},
		{name: "reversed span", minAge: 18, maxAge: 16},
		{name: "negative minimum", minAge: -1, maxAge: 16},
		{name: "storage overflow", minAge: 18, maxAge: math.MaxInt32 + 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			category := validCategory()
			category.MinAge, category.MaxAge = tc.minAge, tc.maxAge
			err := validator.Validate(category)
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, domain.ErrInvalidData)
			}
		})
	}

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
