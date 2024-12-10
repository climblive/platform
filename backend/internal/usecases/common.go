package usecases

import (
	"github.com/climblive/platform/backend/internal/domain"
)

func withScore(contender domain.Contender, scoreKeeper domain.ScoreKeeper) domain.Contender {
	score, err := scoreKeeper.GetScore(contender.ID)
	if err == nil {
		contender.Score = &score
	}

	return contender
}
