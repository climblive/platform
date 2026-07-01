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

func HypotheticalBestZone1(tick Tick) Tick {
	if !tick.Zone1 {
		tick.Zone1 = true
		tick.AttemptsZone1 += 1
	}

	tick.Zone2 = false
	tick.Top = false

	return tick
}

func HypotheticalBestZone2(tick Tick) Tick {
	if !tick.Zone1 {
		tick.Zone1 = true
		tick.AttemptsZone1 += 1
	}

	if !tick.Zone2 {
		tick.Zone2 = true
		tick.AttemptsZone2 += 1
	}

	tick.Top = false

	return tick
}

func HypotheticalBestTop(tick Tick) Tick {
	if !tick.Zone1 {
		tick.Zone1 = true
		tick.AttemptsZone1 += 1
	}

	if !tick.Zone2 {
		tick.Zone2 = true
		tick.AttemptsZone2 += 1
	}

	if !tick.Top {
		tick.Top = true
		tick.AttemptsTop += 1
	}

	return tick
}

func HypotheticalSecondBestTop(tick Tick) Tick {
	tick = HypotheticalBestTop(tick)

	if tick.AttemptsTop == 1 {
		tick.AttemptsZone1 += 1
		tick.AttemptsZone2 += 1
		tick.AttemptsTop += 1
	}

	return tick
}
