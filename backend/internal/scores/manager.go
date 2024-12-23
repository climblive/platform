package scores

import (
	"context"
	"log/slog"
	"maps"
	"slices"
	"sync"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

const pollInterval = 10 * time.Second

type EngineStoreHydrator interface {
	Hydrate(ctx context.Context, contestID domain.ContestID, store EngineStore) error
}

type scoreEngineManagerRepository interface {
	GetContestsCurrentlyRunningOrByStartTime(ctx context.Context, tx domain.Transaction, earliestStartTime, latestStartTime time.Time) ([]domain.Contest, error)
}

type ScoreEngineManager struct {
	repo                scoreEngineManagerRepository
	engineStoreHydrator EngineStoreHydrator
	eventBroker         domain.EventBroker
	handlers            map[domain.ContestID]*engineHandler
}

type engineHandler struct {
	engine             *ScoreEngine
	cancel             func()
	wg                 *sync.WaitGroup
	finalists          int
	qualifyingProblems int
}

func NewScoreEngineManager(repo scoreEngineManagerRepository, engineStoreHydrator EngineStoreHydrator, eventBroker domain.EventBroker) ScoreEngineManager {
	return ScoreEngineManager{
		repo:                repo,
		engineStoreHydrator: engineStoreHydrator,
		eventBroker:         eventBroker,
		handlers:            make(map[domain.ContestID]*engineHandler),
	}
}

func (mngr *ScoreEngineManager) Run(ctx context.Context) *sync.WaitGroup {
	wg := new(sync.WaitGroup)

	wg.Add(1)

	go mngr.run(ctx, wg)

	return wg
}

func (mngr *ScoreEngineManager) run(ctx context.Context, wg *sync.WaitGroup) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("score engine manager panicked", "error", r)
		}
	}()

	defer wg.Done()

	for {
		err := mngr.runPeriodicCheck(ctx)
		if err != nil {
			slog.Error("score engine manager failed to complete periodic check", "error", err)
		}

		sleepTimer := time.NewTimer(pollInterval)

		select {
		case <-ctx.Done():
			slog.Info("score engine manager shutting down", "reason", ctx.Err().Error())

			for handler := range maps.Values(mngr.handlers) {
				handler.cancel()
			}

			for handler := range maps.Values(mngr.handlers) {
				handler.wg.Wait()
			}

			return
		case <-sleepTimer.C:
		}
	}
}

func (mngr *ScoreEngineManager) runPeriodicCheck(ctx context.Context) error {
	now := time.Now()
	contests, err := mngr.repo.GetContestsCurrentlyRunningOrByStartTime(ctx, nil, now, now.Add(time.Hour))
	if err != nil {
		return errors.Wrap(err, 0)
	}

	for contest := range slices.Values(contests) {
		logger := slog.New(slog.Default().Handler()).With("contest_id", contest.ID)

		if handler, ok := mngr.handlers[contest.ID]; ok {
			if contest.QualifyingProblems != handler.qualifyingProblems {
				logger.Info("updating scoring rules", "qualifying_problems", contest.QualifyingProblems)
				handler.engine.SetScoringRules(&HardestProblems{Number: contest.QualifyingProblems})
				handler.qualifyingProblems = contest.QualifyingProblems
			}

			if contest.Finalists != handler.finalists {
				logger.Info("updating ranker", "finalists", contest.Finalists)
				handler.engine.SetRanker(NewBasicRanker(contest.Finalists))
				handler.finalists = contest.Finalists
			}

			continue
		}

		config := slog.Group("config",
			"qualifying_problems", contest.QualifyingProblems,
			"finalists", contest.Finalists)

		switch {
		case contest.TimeBegin != nil && contest.TimeBegin.After(now):
			logger.Info("detected contest about to start", "starting_in", time.Until(*contest.TimeBegin), config)
		default:
			logger.Info("detected contest that is currently running", config)
		}

		store := NewMemoryStore()

		handler := engineHandler{
			engine:             NewScoreEngine(contest.ID, mngr.eventBroker, &HardestProblems{Number: contest.QualifyingProblems}, NewBasicRanker(contest.Finalists), store),
			finalists:          contest.Finalists,
			qualifyingProblems: contest.QualifyingProblems,
		}

		var cancellableCtx context.Context
		cancellableCtx, handler.cancel = context.WithCancel(ctx)
		handler.wg = handler.engine.Run(cancellableCtx)

		hydrationStartTime := time.Now()
		err := mngr.engineStoreHydrator.Hydrate(ctx, contest.ID, store)
		if err != nil {
			logger.Error("failed to hydrate score engine", "error", err)

			handler.cancel()

			continue
		}

		logger.Info("score engine store hydration complete", "time", time.Since(hydrationStartTime))

		mngr.handlers[contest.ID] = &handler
	}

	return nil
}
