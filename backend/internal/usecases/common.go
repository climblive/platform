package usecases

import (
	"github.com/climblive/platform/backend/internal/domain"
)

func withScore(contender domain.Contender, scoreKeeper domain.ScoreKeeper) domain.Contender {
	score, err := scoreKeeper.GetScore(contender.ID)
	if err == nil {
		contender.ScoreUpdated = score.Timestamp
		contender.Score = score.Score
		contender.Placement = score.Placement
		contender.Finalist = score.Finalist
		contender.RankOrder = score.RankOrder
	}

	return contender
}
