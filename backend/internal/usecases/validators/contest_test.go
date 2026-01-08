package validators_test

import (
	"testing"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/usecases/validators"
	"github.com/stretchr/testify/assert"
)

func TestContestValidator(t *testing.T) {
	validator := validators.ContestValidator{}

	validContest := func() domain.Contest {
		return domain.Contest{
			Name:               "Swedish Championships",
			QualifyingProblems: 10,
			Finalists:          7,
			GracePeriod:        time.Minute * 15,
		}
	}

	t.Run("ValidData", func(t *testing.T) {
		err := validator.Validate(validContest())
		assert.NoError(t, err)
	})

	t.Run("EmptyName", func(t *testing.T) {
		contest := validContest()
		contest.Name = whitespaceCharacters

		err := validator.Validate(contest)

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.True(t, validator.IsValidationError(err))
	})

	t.Run("NegativeFinalists", func(t *testing.T) {
		contest := validContest()
		contest.Finalists = -1

		err := validator.Validate(contest)

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.True(t, validator.IsValidationError(err))
	})

	t.Run("FinalistsTooLarge", func(t *testing.T) {
		contest := validContest()
		contest.Finalists = 65536 + 1

		err := validator.Validate(contest)

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.True(t, validator.IsValidationError(err))
	})

	t.Run("NegativeQualifyingProblems", func(t *testing.T) {
		contest := validContest()
		contest.QualifyingProblems = -1

		err := validator.Validate(contest)

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.True(t, validator.IsValidationError(err))
	})

	t.Run("QualifyingProblemsTooLarge", func(t *testing.T) {
		contest := validContest()
		contest.QualifyingProblems = 65536 + 1

		err := validator.Validate(contest)

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.True(t, validator.IsValidationError(err))
	})

	t.Run("NegativeGracePeriod", func(t *testing.T) {
		contest := validContest()
		contest.GracePeriod = -1

		err := validator.Validate(contest)

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.True(t, validator.IsValidationError(err))
	})

	t.Run("GracePeriodLongerThanOneHour", func(t *testing.T) {
		contest := validContest()
		contest.GracePeriod = time.Hour + time.Nanosecond

		err := validator.Validate(contest)

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.True(t, validator.IsValidationError(err))
	})
}
