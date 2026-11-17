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
	favorTopAttempts  bool
	zone1Enabled      bool
	zone2Enabled      bool
}

func (r *BasicRanker) scoresAreEqual(c1, c2 Competitor) bool {
	if r.usePoints {
		return c1.Points == c2.Points
	}

	return c1.Score == c2.Score
}

func NewBasicRanker(numberOfFinalists int, usePoints, zone1Enabled, zone2Enabled, favorTopAttempts bool) *BasicRanker {
	return &BasicRanker{
		numberOfFinalists: numberOfFinalists,
		usePoints:         usePoints,
		favorTopAttempts:  favorTopAttempts,
		zone1Enabled:      zone1Enabled,
		zone2Enabled:      zone2Enabled,
	}
}

func qualifiedCompetitors(competitors iter.Seq[Competitor]) iter.Seq[Competitor] {
	return func(yield func(Competitor) bool) {
		for competitor := range competitors {
			if competitor.Disqualified {
				continue
			}

			if !yield(competitor) {
				return
			}
		}
	}
}

func disqualifiedCompetitors(competitors iter.Seq[Competitor]) iter.Seq[Competitor] {
	return func(yield func(Competitor) bool) {
		for competitor := range competitors {
			if !competitor.Disqualified {
				continue
			}

			if !yield(competitor) {
				return
			}
		}
	}
}

func (r *BasicRanker) RankCompetitors(competitors iter.Seq[Competitor]) []domain.Result {
	var results []domain.Result

	comparator := func(c1, c2 Competitor) int {
		return c1.Compare(c2, r.usePoints || r.favorTopAttempts)
	}

	var previousCompetitor *Competitor
	var placement int
	var gap int
	var numberOfAssignedFinalists int
	var lastFinalistPlacement int

	now := time.Now()

	sortedCompetitors := slices.SortedFunc(qualifiedCompetitors(competitors), comparator)

	for competitor := range disqualifiedCompetitors(competitors) {
		sortedCompetitors = append(sortedCompetitors, competitor)
	}

	for i, competitor := range sortedCompetitors {
		var scoreValue string

		if r.usePoints {
			scoreValue = competitor.Points.String() + "p"
		} else {
			scoreValue = fmt.Sprintf("%dt", competitor.Tops)
			if r.zone2Enabled {
				scoreValue += fmt.Sprintf(" %dz₂", competitor.Zone2s)
			}
			if r.zone1Enabled {
				scoreValue += fmt.Sprintf(" %dz₁", competitor.Zone1s)
			}
		}

		result := domain.Result{
			Timestamp:    now,
			CompetitorID: competitor.ID,
			Score:        scoreValue,
			Placement:    0,
			Finalist:     false,
			RankOrder:    0,
		}

		switch {
		case previousCompetitor == nil:
			placement = 1
			gap = 0
		case r.scoresAreEqual(competitor, *previousCompetitor):
			gap++
		default:
			placement += 1 + gap
			gap = 0
		}

		result.Placement = placement
		result.RankOrder = i

		switch {
		case competitor.Score == Score{}:
			fallthrough
		case competitor.WithdrawnFromFinals:
			fallthrough
		case competitor.Disqualified:
		case numberOfAssignedFinalists < r.numberOfFinalists:
			result.Finalist = true
			numberOfAssignedFinalists++
			lastFinalistPlacement = result.Placement
		case result.Placement == lastFinalistPlacement:
			result.Finalist = true
		}

		results = append(results, result)
		previousCompetitor = &competitor
	}

	return results
}
