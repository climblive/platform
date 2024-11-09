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

type scoreEngineManagerRepository interface {
	GetContestsCurrentlyRunningOrByStartTime(ctx context.Context, tx domain.Transaction, earliestStartTime, latestStartTime time.Time) ([]domain.Contest, error)
	GetProblemsByContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) ([]domain.Problem, error)
	GetContendersByContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) ([]domain.Contender, error)
	GetTicksByContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) ([]domain.Tick, error)
}

type ScoreEngineManager struct {
	repo        scoreEngineManagerRepository
	eventBroker domain.EventBroker
	handlers    map[domain.ContestID]engineHandler
}

type engineHandler struct {
	engine *ScoreEngine
	cancel func()
	wg     *sync.WaitGroup
}

func NewScoreEngineManager(repo scoreEngineManagerRepository, eventBroker domain.EventBroker) ScoreEngineManager {
	return ScoreEngineManager{
		repo:        repo,
		eventBroker: eventBroker,
		handlers:    make(map[domain.ContestID]engineHandler),
	}
}

func (mngr *ScoreEngineManager) Run(ctx context.Context) *sync.WaitGroup {
	wg := new(sync.WaitGroup)
	ready := make(chan struct{}, 1)

	wg.Add(1)

	go mngr.run(ctx, wg, ready)

	<-ready

	return wg
}

func (mngr *ScoreEngineManager) run(ctx context.Context, wg *sync.WaitGroup, ready chan struct{}) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("score engine manager panicked", "error", r)
		}
	}()

	defer wg.Done()

	signalReady := sync.OnceFunc(func() {
		close(ready)
	})

	for {
		err := mngr.runPeriodicCheck(ctx)
		if err != nil {
			slog.Error("score engine manager failed to complete periodic check", "error", err)
		} else {
			signalReady()
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
		if _, ok := mngr.handlers[contest.ID]; ok {
			continue
		}

		logger := slog.New(slog.Default().Handler()).With("contest_id", contest.ID)

		config := slog.Group("config",
			"qualifying_problems", contest.QualifyingProblems,
			"finalists", contest.Finalists)

		switch {
		case contest.TimeBegin != nil && contest.TimeBegin.After(now):
			logger.Info("detected contest about to start", "starting_in", (*contest.TimeBegin).Sub(now), config)
		default:
			logger.Info("detected contest that is currently running", config)
		}

		handler := engineHandler{
			engine: NewScoreEngine(contest.ID, mngr.eventBroker, &HardestProblems{Number: contest.QualifyingProblems}, NewBasicRanker(contest.Finalists)),
		}

		var cancellableCtx context.Context
		cancellableCtx, handler.cancel = context.WithCancel(ctx)
		handler.wg = handler.engine.Run(cancellableCtx)

		hydrationStartTime := time.Now()
		stats, err := mngr.hydrateEngine(ctx, contest.ID)
		if err != nil {
			logger.Error("failed to hydrate score engine", "error", err)

			handler.cancel()

			continue
		}

		logger.Info("score engine hydration complete",
			"time", time.Since(hydrationStartTime),
			slog.Group("stats",
				"contenders", stats.contenders,
				"problems", stats.problems,
				"ticks", stats.ticks,
			),
		)

		mngr.handlers[contest.ID] = handler
	}

	return nil
}

type hydrationStats struct {
	contenders int
	problems   int
	ticks      int
}

func (mngr *ScoreEngineManager) hydrateEngine(ctx context.Context, contestID domain.ContestID) (hydrationStats, error) {
	problems, err := mngr.repo.GetProblemsByContest(ctx, nil, contestID)
	if err != nil {
		return hydrationStats{}, errors.Wrap(err, 0)
	}

	for problem := range slices.Values(problems) {
		mngr.eventBroker.Dispatch(contestID, domain.ProblemAddedEvent{
			ProblemID:  problem.ID,
			PointsTop:  problem.PointsTop,
			PointsZone: problem.PointsZone,
			FlashBonus: problem.FlashBonus,
		})
	}

	contenders, err := mngr.repo.GetContendersByContest(ctx, nil, contestID)
	if err != nil {
		return hydrationStats{}, errors.Wrap(err, 0)
	}

	for contender := range slices.Values(contenders) {
		mngr.eventBroker.Dispatch(contestID, domain.ContenderEnteredEvent{
			ContenderID: contender.ID,
			CompClassID: contender.CompClassID,
		})

		if contender.WithdrawnFromFinals {
			mngr.eventBroker.Dispatch(contestID, domain.ContenderWithdrewFromFinalsEvent{
				ContenderID: contender.ID,
			})
		}

		if contender.Disqualified {
			mngr.eventBroker.Dispatch(contestID, domain.ContenderDisqualifiedEvent{
				ContenderID: contender.ID,
			})
		}
	}

	ticks, err := mngr.repo.GetTicksByContest(ctx, nil, contestID)
	if err != nil {
		return hydrationStats{}, errors.Wrap(err, 0)
	}

	for tick := range slices.Values(ticks) {
		mngr.eventBroker.Dispatch(contestID, domain.AscentRegisteredEvent{
			ContenderID:  *tick.Ownership.ContenderID,
			ProblemID:    tick.ProblemID,
			Top:          tick.Top,
			AttemptsTop:  tick.AttemptsTop,
			Zone:         tick.Zone,
			AttemptsZone: tick.AttemptsZone,
		})
	}

	return hydrationStats{
		contenders: len(contenders),
		problems:   len(problems),
		ticks:      len(ticks),
	}, nil
}
