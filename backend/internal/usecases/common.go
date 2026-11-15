package usecases

import (
	"github.com/climblive/platform/backend/internal/domain"
)

func withResult(contender domain.Contender, resultKeeper domain.ResultKeeper) domain.Contender {
	result, err := resultKeeper.GetResult(contender.ID)
	if err == nil {
		contender.Result = &result
	}

	return contender
}
