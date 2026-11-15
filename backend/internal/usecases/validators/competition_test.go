package validators_test

import (
	"testing"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/usecases/validators"
	"github.com/stretchr/testify/assert"
)

func TestCompetitionValidator(t *testing.T) {
	validator := validators.CompetitionValidator{}

	validCompetition := func() domain.Competition {
		return domain.Competition{
			Name:              "Swedish Championships",
			Country:           "SE",
			GracePeriod:       time.Minute * 15,
			NameRetentionTime: 14 * 24 * time.Hour,
		}
	}

	t.Run("ValidData", func(t *testing.T) {
		err := validator.Validate(validCompetition())
		assert.NoError(t, err)
	})

	t.Run("EmptyName", func(t *testing.T) {
		competition := validCompetition()
		competition.Name = whitespaceCharacters

		err := validator.Validate(competition)

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.True(t, validator.IsValidationError(err))
	})

	t.Run("NegativeFinalists", func(t *testing.T) {
		competition := validCompetition()
		competition.Finalists = -1

		err := validator.Validate(competition)

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.True(t, validator.IsValidationError(err))
	})

	t.Run("FinalistsTooLarge", func(t *testing.T) {
		competition := validCompetition()
		competition.Finalists = 65536 + 1

		err := validator.Validate(competition)

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.True(t, validator.IsValidationError(err))
	})

	t.Run("QualifyingProblemsWithoutPoints", func(t *testing.T) {
		competition := validCompetition()
		competition.UsePoints = false
		competition.QualifyingProblems = 10

		err := validator.Validate(competition)

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.True(t, validator.IsValidationError(err))
	})

	t.Run("PooledPointsWithoutPoints", func(t *testing.T) {
		competition := validCompetition()
		competition.UsePoints = false
		competition.PooledPoints = true

		err := validator.Validate(competition)

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.True(t, validator.IsValidationError(err))
	})

	t.Run("NegativeGracePeriod", func(t *testing.T) {
		competition := validCompetition()
		competition.GracePeriod = -1

		err := validator.Validate(competition)

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.True(t, validator.IsValidationError(err))
	})

	t.Run("GracePeriodLongerThanOneHour", func(t *testing.T) {
		competition := validCompetition()
		competition.GracePeriod = time.Hour + time.Nanosecond

		err := validator.Validate(competition)

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.True(t, validator.IsValidationError(err))
	})

	t.Run("InvalidCountryCode", func(t *testing.T) {
		competition := validCompetition()
		competition.Country = "XYZ"

		err := validator.Validate(competition)

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.True(t, validator.IsValidationError(err))
	})

	t.Run("ValidNameRetentionTimes", func(t *testing.T) {
		for _, nrt := range []time.Duration{14 * 24 * time.Hour, 30 * 24 * time.Hour, 90 * 24 * time.Hour} {
			competition := validCompetition()
			competition.NameRetentionTime = nrt

			err := validator.Validate(competition)

			assert.NoError(t, err)
		}
	})

	t.Run("NameRetentionTimeInvalid", func(t *testing.T) {
		competition := validCompetition()
		competition.NameRetentionTime = 0

		err := validator.Validate(competition)

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.True(t, validator.IsValidationError(err))
	})

	t.Run("MaxAttemptsWithoutPoints", func(t *testing.T) {
		competition := validCompetition()
		competition.UsePoints = false
		competition.MaxAttempts = 1

		err := validator.Validate(competition)

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.True(t, validator.IsValidationError(err))
	})

	t.Run("PointDeductionWithoutPoints", func(t *testing.T) {
		competition := validCompetition()
		competition.UsePoints = false
		competition.PointDeduction = 1

		err := validator.Validate(competition)

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.True(t, validator.IsValidationError(err))
	})

	t.Run("WithPoints", func(t *testing.T) {
		validCompetition := func() domain.Competition {
			return domain.Competition{
				Name:              "Swedish Championships",
				Country:           "SE",
				UsePoints:         true,
				GracePeriod:       time.Minute * 15,
				NameRetentionTime: 14 * 24 * time.Hour,
			}
		}

		t.Run("ValidData", func(t *testing.T) {
			err := validator.Validate(validCompetition())
			assert.NoError(t, err)
		})

		t.Run("NegativeQualifyingProblems", func(t *testing.T) {
			competition := validCompetition()
			competition.QualifyingProblems = -1

			err := validator.Validate(competition)

			assert.ErrorIs(t, err, domain.ErrInvalidData)
			assert.True(t, validator.IsValidationError(err))
		})

		t.Run("QualifyingProblemsTooLarge", func(t *testing.T) {
			competition := validCompetition()
			competition.QualifyingProblems = 65536 + 1

			err := validator.Validate(competition)

			assert.ErrorIs(t, err, domain.ErrInvalidData)
			assert.True(t, validator.IsValidationError(err))
		})

		t.Run("NegativeMaxAttempts", func(t *testing.T) {
			competition := validCompetition()
			competition.MaxAttempts = -1

			err := validator.Validate(competition)

			assert.ErrorIs(t, err, domain.ErrInvalidData)
			assert.True(t, validator.IsValidationError(err))
		})

		t.Run("MaxAttemptsTooLarge", func(t *testing.T) {
			competition := validCompetition()
			competition.MaxAttempts = 1000

			err := validator.Validate(competition)

			assert.ErrorIs(t, err, domain.ErrInvalidData)
			assert.True(t, validator.IsValidationError(err))
		})

		t.Run("NegativePointDeduction", func(t *testing.T) {
			competition := validCompetition()
			competition.PointDeduction = -1

			err := validator.Validate(competition)

			assert.ErrorIs(t, err, domain.ErrInvalidData)
			assert.True(t, validator.IsValidationError(err))
		})

		t.Run("PointDeductionTooLarge", func(t *testing.T) {
			competition := validCompetition()
			competition.PointDeduction = 21_474_836_471

			err := validator.Validate(competition)

			assert.ErrorIs(t, err, domain.ErrInvalidData)
			assert.True(t, validator.IsValidationError(err))
		})
	})
}
