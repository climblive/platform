package scores

import "github.com/climblive/platform/backend/internal/domain"

type keeper struct {
}

func NewScoreKeeper() domain.ScoreKeeper {
	return &keeper{}
}

func (k *keeper) UpdateScore(contenderID domain.ResourceID, score domain.Score) error {
	return nil
}

func (k *keeper) GetScore(contenderID domain.ResourceID) (domain.Score, error) {
	return domain.Score{}, nil
}
