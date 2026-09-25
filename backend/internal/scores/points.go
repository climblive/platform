package scores

import "github.com/climblive/platform/backend/internal/domain"

func CalculatePoints(value domain.ProblemValue, tick Tick, rules Rules) int {
	current := 0
	attempts := 0

	if tick.Zone1 {
		current = value.PointsZone1
		attempts = tick.AttemptsZone1
	}

	if tick.Zone2 {
		current = value.PointsZone2
		attempts = tick.AttemptsZone2
	}

	if tick.Top {
		current = value.PointsTop
		attempts = tick.AttemptsTop

		if tick.AttemptsTop == 1 {
			current += value.FlashBonus
		}
	}

	if rules.MaxAttempts > 0 && attempts > rules.MaxAttempts {
		return 0
	}

	return max(0, current-max(0, attempts-1)*rules.PointDeduction)
}
