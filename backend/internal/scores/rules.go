package scores

import (
	"iter"
	"slices"
)

type HardestProblems struct {
	Number int
}

func (r *HardestProblems) CalculateScore(points iter.Seq[int]) int {
	score := 0

	n := 0
	for _, p := range slices.Backward(slices.Sorted(points)) {
		if n >= r.Number {
			break
		}

		score += p
		n++
	}

	return score
}

func (r *HardestProblems) CalculatePoints(tick Tick, problem Problem, allTicks iter.Seq[Tick]) int {
	points := 0

	if tick.Zone {
		tick.Points = problem.PointsZone
	}

	if tick.Top {
		tick.Points = problem.PointsTop
	}

	if tick.IsFlash() {
		tick.Points += problem.FlashBonus
	}

	return points
}

type PooledPoints struct {
}

func (r *PooledPoints) CalculateScore(points iter.Seq[int]) int {
	score := 0

	n := 0
	for p := range points {
		score += p
		n++
	}

	return score
}

func (r *PooledPoints) CalculatePoints(tick Tick, problem Problem, poolTicks iter.Seq[Tick]) int {
	points := 0

	totalTops := 0
	totalZones := 0
	totalFlashes := 0

	for tick := range poolTicks {
		if tick.IsFlash() {
			totalFlashes += 1
		}

		if tick.Top {
			totalTops += 1
		}

		if tick.Zone {
			totalZones += 1
		}
	}

	if totalTops == 0 {
		totalTops = 1
	}

	if totalZones == 0 {
		totalZones = 1
	}

	if totalFlashes == 0 {
		totalFlashes = 1
	}

	if tick.Zone {
		points = problem.PointsZone / totalZones
	}

	if tick.Top {
		points = problem.PointsTop / totalTops
	}

	if tick.IsFlash() {
		points += problem.FlashBonus / totalFlashes
	}

	return points
}
