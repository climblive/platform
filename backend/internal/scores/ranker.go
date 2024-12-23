package scores

import (
	"iter"
	"slices"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
)

type BasicRanker struct {
	numberOfFinalists int
}

func NewBasicRanker(numberOfFinalists int) *BasicRanker {
	return &BasicRanker{
		numberOfFinalists: numberOfFinalists,
	}
}

func (r *BasicRanker) RankContenders(contenders iter.Seq[Contender]) []domain.Score {
	var scores []domain.Score

	comparator := func(c1, c2 Contender) int {
		return c1.Compare(c2)
	}

	var previousContender *Contender
	var placement int
	var gap int
	var numberOfAssignedFinalists int
	var lastFinalistPlacement int

	now := time.Now()

	for i, contender := range slices.SortedFunc(contenders, comparator) {
		score := domain.Score{
			Timestamp:   now,
			ContenderID: contender.ID,
			Score:       contender.Score,
		}

		switch {
		case previousContender == nil:
			placement = 1
			gap = 0
		case contender.Score == previousContender.Score:
			gap++
		case contender.Score != previousContender.Score:
			placement += 1 + gap
			gap = 0
		}

		score.Placement = placement
		score.RankOrder = i

		switch {
		case contender.Score == 0:
			fallthrough
		case contender.WithdrawnFromFinals:
			fallthrough
		case contender.Disqualified:
		case numberOfAssignedFinalists < r.numberOfFinalists:
			score.Finalist = true
			numberOfAssignedFinalists++
			lastFinalistPlacement = score.Placement
		case score.Placement == lastFinalistPlacement:
			score.Finalist = true
		}

		scores = append(scores, score)
		previousContender = &contender
	}

	return scores
}
