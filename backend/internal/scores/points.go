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

func HypotheticalZone1(tick Tick) Tick {
	return Tick{
		ContenderID:   tick.ContenderID,
		ProblemID:     tick.ProblemID,
		Zone1:         true,
		AttemptsZone1: 999,
		AttemptsZone2: 999,
		AttemptsTop:   999,
	}
}

func HypotheticalZone2(tick Tick) Tick {
	tick = HypotheticalZone1(tick)
	tick.Zone2 = true
	tick.AttemptsZone2 = 999

	return tick
}

func HypotheticalTop(tick Tick) Tick {
	tick = HypotheticalZone2(tick)
	tick.Top = true
	tick.AttemptsTop = 999

	return tick
}

func HypotheticalFlash(tick Tick) Tick {
	tick = HypotheticalTop(tick)
	tick.AttemptsZone1 = 1
	tick.AttemptsZone2 = 1
	tick.AttemptsTop = 1

	return tick
}
