package scores

import (
	"fmt"
	"iter"
	"slices"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
)

type BasicRanker struct {
	numberOfFinalists int
	usePoints         bool
	zone1Enabled      bool
	zone2Enabled      bool
}

func (r *BasicRanker) scoresAreEqual(c1, c2 Contender) bool {
	if r.usePoints {
		return c1.Points == c2.Points
	}

	return c1.Score == c2.Score
}

func NewBasicRanker(numberOfFinalists int, usePoints, zone1Enabled, zone2Enabled bool) *BasicRanker {
	return &BasicRanker{
		numberOfFinalists: numberOfFinalists,
		usePoints:         usePoints,
		zone1Enabled:      zone1Enabled,
		zone2Enabled:      zone2Enabled,
	}
}

func qualifiedContenders(contenders iter.Seq[Contender]) iter.Seq[Contender] {
	return func(yield func(Contender) bool) {
		for contender := range contenders {
			if contender.Disqualified {
				continue
			}

			if !yield(contender) {
				return
			}
		}
	}
}

func disqualifiedContenders(contenders iter.Seq[Contender]) iter.Seq[Contender] {
	return func(yield func(Contender) bool) {
		for contender := range contenders {
			if !contender.Disqualified {
				continue
			}

			if !yield(contender) {
				return
			}
		}
	}
}

func (r *BasicRanker) RankContenders(contenders iter.Seq[Contender]) []domain.Result {
	var results []domain.Result

	comparator := func(c1, c2 Contender) int {
		return c1.Compare(c2)
	}

	var previousContender *Contender
	var placement int
	var gap int
	var numberOfAssignedFinalists int
	var lastFinalistPlacement int

	now := time.Now()

	sortedContenders := slices.SortedFunc(qualifiedContenders(contenders), comparator)

	for contender := range disqualifiedContenders(contenders) {
		sortedContenders = append(sortedContenders, contender)
	}

	for i, contender := range sortedContenders {
		var scoreValue string

		if r.usePoints {
			scoreValue = contender.Points.String() + "p"
		} else {
			scoreValue = fmt.Sprintf("%dt", contender.Tops)
			if r.zone2Enabled {
				scoreValue += fmt.Sprintf(" %dz₂", contender.Zone2s)
			}
			if r.zone1Enabled {
				scoreValue += fmt.Sprintf(" %dz₁", contender.Zone1s)
			}
		}

		result := domain.Result{
			Timestamp:   now,
			ContenderID: contender.ID,
			Score:       scoreValue,
			Placement:   0,
			Finalist:    false,
			RankOrder:   0,
		}

		switch {
		case previousContender == nil:
			placement = 1
			gap = 0
		case r.scoresAreEqual(contender, *previousContender):
			gap++
		default:
			placement += 1 + gap
			gap = 0
		}

		result.Placement = placement
		result.RankOrder = i

		switch {
		case contender.Score == Score{}:
			fallthrough
		case contender.WithdrawnFromFinals:
			fallthrough
		case contender.Disqualified:
		case numberOfAssignedFinalists < r.numberOfFinalists:
			result.Finalist = true
			numberOfAssignedFinalists++
			lastFinalistPlacement = result.Placement
		case result.Placement == lastFinalistPlacement:
			result.Finalist = true
		}

		results = append(results, result)
		previousContender = &contender
	}

	return results
}
