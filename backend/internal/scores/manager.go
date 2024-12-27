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
	"github.com/google/uuid"
)

type listRequest struct {
	contestID domain.ContestID
	response  chan []domain.ScoreEngineInstanceID
}

type stopRequest struct {
	instanceID domain.ScoreEngineInstanceID
	response   chan struct{}
}

type startRequest struct {
	contestID domain.ContestID
	response  chan domain.ScoreEngineInstanceID
}

const pollInterval = 10 * time.Second

type EngineStoreHydrator interface {
	Hydrate(ctx context.Context, contestID domain.ContestID, store EngineStore) error
}

type scoreEngineManagerRepository interface {
	GetContestsCurrentlyRunningOrByStartTime(ctx context.Context, tx domain.Transaction, earliestStartTime, latestStartTime time.Time) ([]domain.Contest, error)
	GetContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) (domain.Contest, error)
}

type ScoreEngineManager struct {
	repo                scoreEngineManagerRepository
	engineStoreHydrator EngineStoreHydrator
	eventBroker         domain.EventBroker
	handlers            map[domain.ContestID]*engineHandler
	requests            chan any
}

type engineHandler struct {
	instanceID         domain.ScoreEngineInstanceID
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

func (mngr *ScoreEngineManager) ListScoreEnginesByContest(
	ctx context.Context,
	contestID domain.ContestID,
) ([]domain.ScoreEngineInstanceID, error) {
	response := make(chan []domain.ScoreEngineInstanceID, 1)

	mngr.requests <- listRequest{
		contestID: contestID,
		response:  response,
	}

	select {
	case instances := <-response:
		return instances, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (mngr *ScoreEngineManager) StopScoreEngine(
	ctx context.Context,
	instanceID domain.ScoreEngineInstanceID,
) error {
	response := make(chan struct{}, 1)

	mngr.requests <- stopRequest{
		instanceID: instanceID,
		response:   response,
	}

	select {
	case <-response:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (mngr *ScoreEngineManager) StartScoreEngine(
	ctx context.Context,
	contestID domain.ContestID,
) (domain.ScoreEngineInstanceID, error) {
	response := make(chan domain.ScoreEngineInstanceID, 1)

	mngr.requests <- startRequest{
		contestID: contestID,
		response:  response,
	}

	select {
	case instanceID := <-response:
		return instanceID, nil
	case <-ctx.Done():
		return uuid.UUID{}, ctx.Err()
	}
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
		case request := <-mngr.requests:
			switch req := request.(type) {
			case listRequest:
				req.response <- mngr.listScoreEnginesByContest(req.contestID)
				close(req.response)
			case stopRequest:
				mngr.stopScoreEngine(req.instanceID)
				close(req.response)
			case startRequest:
				instanceID, err := mngr.startScoreEngine(ctx, req.contestID)
				if err != nil {
					panic("not implemented")
				}

				req.response <- instanceID
				close(req.response)
			}
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
		if handler, ok := mngr.handlers[contest.ID]; ok {
			logger := slog.New(slog.Default().Handler()).
				With("contest_id", contest.ID).
				With("instance_id", handler.instanceID)

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

		_, _ = mngr.startScoreEngine(ctx, contest.ID)
	}

	return nil
}

func (mngr *ScoreEngineManager) startScoreEngine(ctx context.Context, contestID domain.ContestID) (domain.ScoreEngineInstanceID, error) {
	now := time.Now()

	contest, err := mngr.repo.GetContest(ctx, nil, contestID)
	if err != nil {
		return uuid.UUID{}, errors.Wrap(err, 0)
	}

	logger := slog.New(slog.Default().Handler()).With("contest_id", contestID)

	logger = logger.With(slog.Group("config",
		"qualifying_problems", contest.QualifyingProblems,
		"finalists", contest.Finalists))

	if contest.TimeBegin != nil && contest.TimeBegin.After(now) {
		logger = logger.With("starting_in", time.Until(*contest.TimeBegin))
	}

	logger.Info("spinning up score engine")

	instanceID := uuid.New()
	store := NewMemoryStore()
	rules := &HardestProblems{Number: contest.QualifyingProblems}
	ranker := NewBasicRanker(contest.Finalists)
	engine := NewScoreEngine(contestID, mngr.eventBroker, rules, ranker, store)

	hydrationStartTime := time.Now()
	err = mngr.engineStoreHydrator.Hydrate(ctx, contestID, store)
	if err != nil {
		logger.Error("hydration failed", "error", err)

		return uuid.UUID{}, errors.Wrap(err, 0)
	}

	logger.Debug("score engine store hydration complete", "time", time.Since(hydrationStartTime))

	cancellableCtx, cancel := context.WithCancel(ctx)
	wg := engine.Run(cancellableCtx)

	engine.ScoreAll()

	mngr.handlers[contestID] = &engineHandler{
		instanceID:         instanceID,
		engine:             engine,
		cancel:             cancel,
		wg:                 wg,
		finalists:          contest.Finalists,
		qualifyingProblems: contest.QualifyingProblems,
	}

	logger.Info("score engine started", "instance_id", instanceID)

	return instanceID, nil
}
