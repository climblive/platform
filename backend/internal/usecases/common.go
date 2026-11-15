package usecases

import (
	"github.com/climblive/platform/backend/internal/domain"
)

func withResult(competitor domain.Competitor, resultKeeper domain.ResultKeeper) domain.Competitor {
	result, err := resultKeeper.GetResult(competitor.ID)
	if err == nil {
		competitor.Result = &result
	}

	return competitor
}
