package scores

import "github.com/climblive/platform/backend/internal/domain"

func CalculatePoints(value domain.ProblemValue, tick Tick) int {
	current := 0

	if tick.Zone1 {
		current = value.PointsZone1
	}

	if tick.Zone2 {
		current = value.PointsZone2
	}

	if tick.Top {
		current = value.PointsTop

		if tick.AttemptsTop == 1 {
			current += value.FlashBonus
		}
	}

	return current
}
