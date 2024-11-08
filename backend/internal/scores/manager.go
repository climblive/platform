package scores

import (
	"context"
	"log/slog"
	"slices"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
)

type scoreEngineManagerRepository interface {
	GetContestsRunningOrAboutToStart(ctx context.Context, tx domain.Transaction, earliestStartTime, latestStartTime time.Time) ([]domain.Contest, error)
	GetProblemsByContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) ([]domain.Problem, error)
	GetContendersByContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) ([]domain.Contender, error)
	GetTicksByContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) ([]domain.Tick, error)
}

type ScoreEngineManager struct {
	repo        scoreEngineManagerRepository
	eventBroker domain.EventBroker
	running     map[domain.ContestID]struct{}
}

func NewScoreEngineManager(repo scoreEngineManagerRepository, eventBroker domain.EventBroker) ScoreEngineManager {
	return ScoreEngineManager{
		repo:        repo,
		eventBroker: eventBroker,
		running:     make(map[domain.ContestID]struct{}),
	}
}

func (mngr *ScoreEngineManager) Run(ctx context.Context) {
	go func() {
		for {
			mngr.pollContests(ctx)

			time.Sleep(10 * time.Second)
		}
	}()
}

func (mngr *ScoreEngineManager) pollContests(ctx context.Context) {
	now := time.Now()
	contests, err := mngr.repo.GetContestsRunningOrAboutToStart(ctx, nil, now.Add(-1*time.Hour), now.Add(time.Hour))
	if err != nil {
		panic(err)
	}

	for contest := range slices.Values(contests) {
		if _, ok := mngr.running[contest.ID]; ok {
			continue
		}

		logger := slog.New(slog.Default().Handler()).With("contest_id", int(contest.ID))

		startTime := time.Now()

		engine := NewScoreEngine(contest.ID, mngr.eventBroker, &HardestProblems{Number: contest.QualifyingProblems}, NewBasicRanker(contest.Finalists))

		logger.Info("launching score engine",
			"contest_id", contest.ID,
			"qualifying_problems", contest.QualifyingProblems,
			"finalists", contest.Finalists)
		controlChannel := engine.Run(context.Background())
		mngr.running[contest.ID] = struct{}{}

		<-controlChannel

		problems, err := mngr.repo.GetProblemsByContest(ctx, nil, contest.ID)
		if err != nil {
			panic(err)
		}

		for problem := range slices.Values(problems) {
			mngr.eventBroker.Dispatch(1, domain.ProblemAddedEvent{
				ProblemID:  problem.ID,
				PointsTop:  problem.PointsTop,
				PointsZone: problem.PointsZone,
				FlashBonus: problem.FlashBonus,
			})
		}

		contenders, err := mngr.repo.GetContendersByContest(ctx, nil, contest.ID)
		if err != nil {
			panic(err)
		}

		for contender := range slices.Values(contenders) {
			mngr.eventBroker.Dispatch(1, domain.ContenderEnteredEvent{
				ContenderID: contender.ID,
				CompClassID: contender.CompClassID,
			})

			if contender.WithdrawnFromFinals {
				mngr.eventBroker.Dispatch(1, domain.ContenderWithdrewFromFinalsEvent{
					ContenderID: contender.ID,
				})
			}

			if contender.Disqualified {
				mngr.eventBroker.Dispatch(1, domain.ContenderDisqualifiedEvent{
					ContenderID: contender.ID,
				})
			}
		}

		ticks, err := mngr.repo.GetTicksByContest(ctx, nil, contest.ID)
		if err != nil {
			panic(err)
		}

		for tick := range slices.Values(ticks) {
			mngr.eventBroker.Dispatch(1, domain.AscentRegisteredEvent{
				ContenderID:  *tick.Ownership.ContenderID,
				ProblemID:    tick.ProblemID,
				Top:          tick.Top,
				AttemptsTop:  tick.AttemptsTop,
				Zone:         tick.Zone,
				AttemptsZone: tick.AttemptsTop,
			})
		}

		logger.Info("score engine hydration complete",
			"time", time.Since(startTime),
			"contenders", len(contenders),
			"problems", len(problems),
			"ticks", len(ticks),
		)
	}
}
