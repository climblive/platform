package scores

import "github.com/climblive/platform/backend/internal/domain"

func CalculatePoints(value domain.ProblemValue, tick Tick, rules Rules) int {
	calculate := func(points, attempts int) int {
		if rules.MaxAttempts > 0 && attempts > rules.MaxAttempts {
			return 0
		}

		return max(0, points-max(0, attempts-1)*rules.PointDeduction)
	}

	points := [3]int{}

	if tick.Zone1 {
		points[0] = calculate(value.PointsZone1, tick.AttemptsZone1)
	}

	if tick.Zone2 {
		points[1] = calculate(value.PointsZone2, tick.AttemptsZone2)
	}

	if tick.Top {
		top := value.PointsTop

		if tick.AttemptsTop == 1 {
			top += value.FlashBonus
		}

		points[2] = calculate(top, tick.AttemptsTop)
	}

	return max(points[0], points[1], points[2])
}
