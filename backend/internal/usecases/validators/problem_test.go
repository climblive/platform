package validators_test

import (
	"testing"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/usecases/validators"
	"github.com/stretchr/testify/assert"
)

func TestProblemValidator(t *testing.T) {
	validator := validators.ProblemValidator{}

	validProblem := func() domain.Problem {
		return domain.Problem{
			ID:                 domain.ProblemID(1),
			ContestID:          domain.ContestID(1),
			Number:             1,
			HoldColorPrimary:   "#ffffff",
			HoldColorSecondary: "#000000",
			Description:        "First boulder",
			Zone1Enabled:       true,
			Zone2Enabled:       true,
			ProblemValue: domain.ProblemValue{
				PointsZone1: 50,
				PointsZone2: 75,
				PointsTop:   100,
				FlashBonus:  10,
			},
		}
	}

	t.Run("ValidData", func(t *testing.T) {
		err := validator.Validate(validProblem())
		assert.NoError(t, err)
	})

	t.Run("NumberNegative", func(t *testing.T) {
		problem := validProblem()
		problem.Number = -5

		err := validator.Validate(problem)

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.True(t, validator.IsValidationError(err))
	})

	t.Run("InvalidPrimaryColorFormat", func(t *testing.T) {
		problem := validProblem()
		problem.HoldColorPrimary = "invalid"

		err := validator.Validate(problem)

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.True(t, validator.IsValidationError(err))
	})

	t.Run("InvalidSecondaryColor", func(t *testing.T) {
		problem := validProblem()
		problem.HoldColorSecondary = "invalid"

		err := validator.Validate(problem)

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.True(t, validator.IsValidationError(err))
	})

	t.Run("NegativePointsTop", func(t *testing.T) {
		problem := validProblem()
		problem.PointsTop = -1

		err := validator.Validate(problem)

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.True(t, validator.IsValidationError(err))
	})

	t.Run("PointsTopExceedsMax", func(t *testing.T) {
		problem := validProblem()
		problem.PointsTop = 2_147_483_648

		err := validator.Validate(problem)

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.True(t, validator.IsValidationError(err))
	})

	t.Run("NegativeFlashBonus", func(t *testing.T) {
		problem := validProblem()
		problem.FlashBonus = -1

		err := validator.Validate(problem)

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.True(t, validator.IsValidationError(err))
	})

	t.Run("FlashBonusExceedsMax", func(t *testing.T) {
		problem := validProblem()
		problem.FlashBonus = 2_147_483_648

		err := validator.Validate(problem)

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.True(t, validator.IsValidationError(err))
	})

	t.Run("NegativePointsZone1", func(t *testing.T) {
		problem := validProblem()
		problem.PointsZone1 = -1

		err := validator.Validate(problem)

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.True(t, validator.IsValidationError(err))
	})

	t.Run("PointsZone1ExceedsMax", func(t *testing.T) {
		problem := validProblem()
		problem.PointsZone1 = 2_147_483_648

		err := validator.Validate(problem)

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.True(t, validator.IsValidationError(err))
	})

	t.Run("NegativePointsZone2", func(t *testing.T) {
		problem := validProblem()
		problem.PointsZone2 = -1

		err := validator.Validate(problem)

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.True(t, validator.IsValidationError(err))
	})

	t.Run("PointsZone2ExceedsMax", func(t *testing.T) {
		problem := validProblem()
		problem.PointsZone2 = 2_147_483_648

		err := validator.Validate(problem)

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.True(t, validator.IsValidationError(err))
	})

	t.Run("Zone2EnabledWithoutZone1", func(t *testing.T) {
		problem := validProblem()
		problem.Zone1Enabled = false
		problem.Zone2Enabled = true

		err := validator.Validate(problem)

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.True(t, validator.IsValidationError(err))
	})
}
