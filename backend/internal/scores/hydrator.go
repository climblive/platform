package scores

import (
	"context"
	"slices"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

type standardEngineStoreHydratorRepository interface {
	GetCompetition(ctx context.Context, tx domain.Transaction, competitionID domain.CompetitionID) (domain.Competition, error)
	GetCompetitorsByCompetition(ctx context.Context, tx domain.Transaction, competitionID domain.CompetitionID) ([]domain.Competitor, error)
	GetProblemsByCompetition(ctx context.Context, tx domain.Transaction, competitionID domain.CompetitionID) ([]domain.Problem, error)
	GetTicksByCompetition(ctx context.Context, tx domain.Transaction, competitionID domain.CompetitionID) ([]domain.Tick, error)
}

type StandardEngineStoreHydrator struct {
	Repo standardEngineStoreHydratorRepository
}

func (h *StandardEngineStoreHydrator) Hydrate(ctx context.Context, competitionID domain.CompetitionID, store EngineStore) error {
	competition, err := h.Repo.GetCompetition(ctx, nil, competitionID)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	store.SaveRules(Rules{
		QualifyingProblems: competition.QualifyingProblems,
		Finalists:          competition.Finalists,
		UsePoints:          competition.UsePoints,
		FavorTopAttempts:   competition.FavorTopAttempts,
		PooledPoints:       competition.PooledPoints,
		MaxAttempts:        competition.MaxAttempts,
		PointDeduction:     competition.PointDeduction,
	})

	problems, err := h.Repo.GetProblemsByCompetition(ctx, nil, competitionID)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	for problem := range slices.Values(problems) {
		store.SaveProblem(Problem{
			ID: problem.ID,
			ProblemValue: ProblemValue{
				PointsZone1: domain.Tenths(problem.PointsZone1) * 10,
				PointsZone2: domain.Tenths(problem.PointsZone2) * 10,
				PointsTop:   domain.Tenths(problem.PointsTop) * 10,
				FlashBonus:  domain.Tenths(problem.FlashBonus) * 10,
			},
			Zone1Enabled: problem.Zone1Enabled,
			Zone2Enabled: problem.Zone2Enabled,
		})
	}

	competitors, err := h.Repo.GetCompetitorsByCompetition(ctx, nil, competitionID)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	for competitor := range slices.Values(competitors) {
		if competitor.CategoryID == 0 {
			continue
		}

		store.SaveCompetitor(Competitor{
			ID:                  competitor.ID,
			CategoryID:          competitor.CategoryID,
			WithdrawnFromFinals: competitor.WithdrawnFromFinals,
			Disqualified:        competitor.Disqualified,
			Score: Score{
				Tops:           0,
				Points:         0,
				AttemptsTops:   0,
				Zone1s:         0,
				AttemptsZone1s: 0,
				Zone2s:         0,
				AttemptsZone2s: 0,
			},
		})
	}

	ticks, err := h.Repo.GetTicksByCompetition(ctx, nil, competitionID)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	for tick := range slices.Values(ticks) {
		store.SaveTick(*tick.Ownership.CompetitorID, Tick{
			ID:            tick.ID,
			Revision:      tick.Revision,
			CompetitorID:  *tick.Ownership.CompetitorID,
			ProblemID:     tick.ProblemID,
			Zone1:         tick.Zone1,
			AttemptsZone1: tick.AttemptsZone1,
			Zone2:         tick.Zone2,
			AttemptsZone2: tick.AttemptsZone2,
			Top:           tick.Top,
			AttemptsTop:   tick.AttemptsTop,
		})
	}

	return nil
}
