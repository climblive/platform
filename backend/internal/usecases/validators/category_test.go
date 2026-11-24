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
			EligibilityBasis: domain.EligibilityBasisAge,
			Name:             "Females",
			Description:      "Female climbers",
			TimeBegin:        now,
			TimeEnd:          now.Add(time.Hour),
		}
	}

	t.Run("ValidData", func(t *testing.T) {
		err := validator.Validate(validCategory())
		assert.NoError(t, err)
	})

	for _, tc := range []struct {
		name         string
		minAge       *int
		maxAge       *int
		basis        domain.EligibilityBasis
		minBirthYear *int
		maxBirthYear *int
		valid        bool
	}{
		{name: "AdultSpan", basis: domain.EligibilityBasisAge, minAge: new(18), maxAge: new(99), valid: true},
		{name: "SingleAge", basis: domain.EligibilityBasisAge, minAge: new(16), maxAge: new(16), valid: true},
		{name: "AgeZero", basis: domain.EligibilityBasisAge, minAge: new(0), maxAge: new(15), valid: true},
		{name: "NoAgeLimits", basis: domain.EligibilityBasisAge, valid: true},
		{name: "NoMinimumAge", basis: domain.EligibilityBasisAge, maxAge: new(15), valid: true},
		{name: "NoMaximumAge", basis: domain.EligibilityBasisAge, minAge: new(18), valid: true},
		{name: "ReversedSpan", basis: domain.EligibilityBasisAge, minAge: new(18), maxAge: new(16)},
		{name: "NegativeMinimum", basis: domain.EligibilityBasisAge, minAge: new(-1), maxAge: new(16)},
		{name: "NegativeMaximum", basis: domain.EligibilityBasisAge, maxAge: new(-1)},
		{name: "MinimumStorageOverflow", basis: domain.EligibilityBasisAge, minAge: new(math.MaxInt32 + 1)},
		{name: "MaximumStorageOverflow", basis: domain.EligibilityBasisAge, maxAge: new(math.MaxInt32 + 1)},
		{name: "AgeWithBirthYears", basis: domain.EligibilityBasisAge, minBirthYear: new(2010)},
		{name: "BirthYears", basis: domain.EligibilityBasisBirthYear, minBirthYear: new(2010), maxBirthYear: new(2014), valid: true},
		{name: "SingleBirthYear", basis: domain.EligibilityBasisBirthYear, minBirthYear: new(2010), maxBirthYear: new(2010), valid: true},
		{name: "NoBirthYearLimits", basis: domain.EligibilityBasisBirthYear, valid: true},
		{name: "NoMinimumBirthYear", basis: domain.EligibilityBasisBirthYear, maxBirthYear: new(2009), valid: true},
		{name: "NoMaximumBirthYear", basis: domain.EligibilityBasisBirthYear, minBirthYear: new(2012), valid: true},
		{name: "CalendarYearLimits", basis: domain.EligibilityBasisBirthYear, minBirthYear: new(1), maxBirthYear: new(9999), valid: true},
		{name: "ReversedBirthYears", basis: domain.EligibilityBasisBirthYear, minBirthYear: new(2014), maxBirthYear: new(2010)},
		{name: "BirthYearZero", basis: domain.EligibilityBasisBirthYear, minBirthYear: new(0)},
		{name: "NegativeBirthYear", basis: domain.EligibilityBasisBirthYear, maxBirthYear: new(-1)},
		{name: "MinimumBirthYearTooLarge", basis: domain.EligibilityBasisBirthYear, minBirthYear: new(10000)},
		{name: "MaximumBirthYearTooLarge", basis: domain.EligibilityBasisBirthYear, maxBirthYear: new(10000)},
		{name: "BirthYearsWithAges", basis: domain.EligibilityBasisBirthYear, minAge: new(0)},
		{name: "MissingBasis"},
		{name: "UnknownBasis", basis: "unknown"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			category := validCategory()
			category.MinAge, category.MaxAge = tc.minAge, tc.maxAge
			category.EligibilityBasis = tc.basis
			category.MinBirthYear, category.MaxBirthYear = tc.minBirthYear, tc.maxBirthYear
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
