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
}

func NewBasicRanker(numberOfFinalists int, usePoints bool) *BasicRanker {
	return &BasicRanker{
		numberOfFinalists: numberOfFinalists,
		usePoints:         usePoints,
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

	sortedContenders := slices.SortedFunc(qualifiedContenders(contenders), comparator)

	for contender := range disqualifiedContenders(contenders) {
		sortedContenders = append(sortedContenders, contender)
	}

	for i, contender := range sortedContenders {
		var scoreValue string

		if r.usePoints {
			scoreValue = fmt.Sprintf("%dp", contender.Points)
		} else {
			scoreValue = fmt.Sprintf("%dt %dz₂ %dz₁", contender.Tops, contender.Zone2s, contender.Zone1s)
		}

		score := domain.Score{
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
		case contender.Score == previousContender.Score:
			gap++
		case contender.Score != previousContender.Score:
			placement += 1 + gap
			gap = 0
		}

		score.Placement = placement
		score.RankOrder = i

		switch {
		case contender.Score == Score{}:
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
