package scores

import (
	"context"
	"log/slog"
	"maps"
	"slices"
	"sync"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
)

const pollInterval = 10 * time.Second

type scoreEngineManagerRepository interface {
	GetContestsRunningOrAboutToStart(ctx context.Context, tx domain.Transaction, earliestStartTime, latestStartTime time.Time) ([]domain.Contest, error)
	GetProblemsByContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) ([]domain.Problem, error)
	GetContendersByContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) ([]domain.Contender, error)
	GetTicksByContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) ([]domain.Tick, error)
}

type ScoreEngineManager struct {
	repo        scoreEngineManagerRepository
	eventBroker domain.EventBroker
	engines     map[domain.ContestID]managedEngine
}

type managedEngine struct {
	engine *ScoreEngine
	cancel func()
	wg     *sync.WaitGroup
}

func NewScoreEngineManager(repo scoreEngineManagerRepository, eventBroker domain.EventBroker) ScoreEngineManager {
	return ScoreEngineManager{
		repo:        repo,
		eventBroker: eventBroker,
		engines:     make(map[domain.ContestID]managedEngine),
	}
}

func (mngr *ScoreEngineManager) Run(ctx context.Context) *sync.WaitGroup {
	wg := new(sync.WaitGroup)
	ready := make(chan struct{}, 1)

	wg.Add(1)

	go func() {
		defer wg.Done()

		mngr.poll(ctx)
		ready <- struct{}{}

		for {
			sleepTimer := time.NewTimer(pollInterval)

			select {
			case <-ctx.Done():
				slog.Info("score engine manager shutting down", "reason", ctx.Err().Error())

				for managedEngine := range maps.Values(mngr.engines) {
					managedEngine.cancel()
				}

				for managedEngine := range maps.Values(mngr.engines) {
					managedEngine.wg.Wait()
				}

				return
			case <-sleepTimer.C:
			}

			mngr.poll(ctx)
		}
	}()

	<-ready

	return wg
}

func (mngr *ScoreEngineManager) poll(ctx context.Context) {
	now := time.Now()
	contests, err := mngr.repo.GetContestsRunningOrAboutToStart(ctx, nil, now.Add(-1*time.Hour), now.Add(time.Hour))
	if err != nil {
		panic(err)
	}

	for contest := range slices.Values(contests) {
		if _, ok := mngr.engines[contest.ID]; ok {
			continue
		}

		logger := slog.New(slog.Default().Handler()).With("contest_id", int(contest.ID))

		startTime := time.Now()

		logger.Info("revving up score engine",
			"qualifying_problems", contest.QualifyingProblems,
			"finalists", contest.Finalists)

		engine := NewScoreEngine(contest.ID, mngr.eventBroker, &HardestProblems{Number: contest.QualifyingProblems}, NewBasicRanker(contest.Finalists))

		cancellableCtx, cancel := context.WithCancel(ctx)
		wg := engine.Run(cancellableCtx)

		mngr.engines[contest.ID] = managedEngine{
			engine: engine,
			cancel: cancel,
			wg:     wg,
		}

		logger.Info("score engine ready for hydration")

		stats := mngr.hydrateEngine(ctx, contest.ID)

		logger.Info("score engine hydration complete",
			"time", time.Since(startTime),
			slog.Group("stats",
				"contenders", stats.contenders,
				"problems", stats.problems,
				"ticks", stats.ticks,
			),
		)
	}
}

type hydrationStats struct {
	contenders int
	problems   int
	ticks      int
}

func (mngr *ScoreEngineManager) hydrateEngine(ctx context.Context, contestID domain.ContestID) hydrationStats {
	problems, err := mngr.repo.GetProblemsByContest(ctx, nil, contestID)
	if err != nil {
		panic(err)
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
		panic(err)
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
		panic(err)
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
	}
}
