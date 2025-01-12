package scores

import (
	"context"
	"slices"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

type standardEngineStoreHydratorRepository interface {
	GetContendersByContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) ([]domain.Contender, error)
	GetProblemsByContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) ([]domain.Problem, error)
	GetTicksByContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) ([]domain.Tick, error)
}

type StandardEngineStoreHydrator struct {
	Repo standardEngineStoreHydratorRepository
}

func (h *StandardEngineStoreHydrator) Hydrate(ctx context.Context, contestID domain.ContestID, store EngineStore) error {
	problems, err := h.Repo.GetProblemsByContest(ctx, nil, contestID)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	for problem := range slices.Values(problems) {
		store.SaveProblem(Problem{
			ID:         problem.ID,
			PointsTop:  problem.PointsTop,
			PointsZone: problem.PointsZone,
			FlashBonus: problem.FlashBonus,
		})
	}

	contenders, err := h.Repo.GetContendersByContest(ctx, nil, contestID)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	for contender := range slices.Values(contenders) {
		if contender.CompClassID == 0 {
			continue
		}

		store.SaveContender(Contender{
			ID:                  contender.ID,
			CompClassID:         contender.CompClassID,
			WithdrawnFromFinals: contender.WithdrawnFromFinals,
			Disqualified:        contender.Disqualified,
		})
	}

	ticks, err := h.Repo.GetTicksByContest(ctx, nil, contestID)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	for tick := range slices.Values(ticks) {
		store.SaveTick(*tick.Ownership.ContenderID, Tick{
			ProblemID:    tick.ProblemID,
			Top:          tick.Top,
			AttemptsTop:  tick.AttemptsTop,
			Zone:         tick.Zone,
			AttemptsZone: tick.AttemptsZone,
		})
	}

	return nil
}
