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

var ErrAlreadyStarted = errors.New("already started")

type ScoreEngineDescriptor struct {
	InstanceID domain.ScoreEngineInstanceID
	ContestID  domain.ContestID
}

type Request[A any, R any] struct {
	Args     A
	Response chan<- Response[R]
}

type Response[R any] struct {
	Value R
	Err   error
}

func (r Request[A, R]) Do(ctx context.Context, requests chan<- any) (R, error) {
	response := make(chan Response[R], 1)
	r.Response = response

	requests <- r

	select {
	case r := <-response:
		return r.Value, r.Err
	case <-ctx.Done():
		var empty R
		return empty, ctx.Err()
	}
}

type listScoreEnginesArguments struct {
	contestID domain.ContestID
}

type stopScoreEngineArguments struct {
	instanceID domain.ScoreEngineInstanceID
}

type getScoreEngineArguments struct {
	instanceID domain.ScoreEngineInstanceID
}

type startScoreEngineArguments struct {
	contestID    domain.ContestID
	terminatedBy time.Time
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
	terminations        chan domain.ScoreEngineInstanceID
}

type engineHandler struct {
	instanceID         domain.ScoreEngineInstanceID
	driver             *ScoreEngineDriver
	stop               func()
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
		requests:            make(chan any),
		terminations:        make(chan domain.ScoreEngineInstanceID),
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
) ([]ScoreEngineDescriptor, error) {
	request := Request[listScoreEnginesArguments, []ScoreEngineDescriptor]{Args: listScoreEnginesArguments{contestID: contestID}}
	return request.Do(ctx, mngr.requests)
}

func (mngr *ScoreEngineManager) StopScoreEngine(
	ctx context.Context,
	instanceID domain.ScoreEngineInstanceID,
) error {
	request := Request[stopScoreEngineArguments, struct{}]{Args: stopScoreEngineArguments{instanceID: instanceID}}
	_, err := request.Do(ctx, mngr.requests)

	return err
}

func (mngr *ScoreEngineManager) StartScoreEngine(
	ctx context.Context,
	contestID domain.ContestID,
	terminatedBy time.Time,
) (domain.ScoreEngineInstanceID, error) {
	request := Request[startScoreEngineArguments, domain.ScoreEngineInstanceID]{Args: startScoreEngineArguments{
		contestID:    contestID,
		terminatedBy: terminatedBy,
	}}
	return request.Do(ctx, mngr.requests)
}

func (mngr *ScoreEngineManager) GetScoreEngine(ctx context.Context, instanceID domain.ScoreEngineInstanceID) (ScoreEngineDescriptor, error) {
	request := Request[getScoreEngineArguments, ScoreEngineDescriptor]{Args: getScoreEngineArguments{instanceID: instanceID}}
	return request.Do(ctx, mngr.requests)
}

func (mngr *ScoreEngineManager) run(ctx context.Context, wg *sync.WaitGroup) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("score engine manager panicked", "error", r)
		}
	}()

	defer wg.Done()

	ticker := time.Tick(pollInterval)

	mngr.runPeriodicCheck(ctx)

	for {

		select {
		case <-ctx.Done():
			slog.Info("score engine manager shutting down", "reason", ctx.Err().Error())

			for handler := range maps.Values(mngr.handlers) {
				handler.stop()
			}

			for handler := range maps.Values(mngr.handlers) {
				handler.wg.Wait()
			}

			return
		case <-ticker:
			mngr.runPeriodicCheck(ctx)
		case request := <-mngr.requests:
			mngr.handleRequest(request)
		case terminatedInstanceID := <-mngr.terminations:
			for contestID, handler := range mngr.handlers {
				if handler.instanceID == terminatedInstanceID {
					slog.Warn("removing crashed score engine", "instance_id", terminatedInstanceID)
					delete(mngr.handlers, contestID)

					break
				}
			}
		}
	}
}

func (mngr *ScoreEngineManager) handleRequest(request any) {
	switch req := request.(type) {
	case Request[listScoreEnginesArguments, []ScoreEngineDescriptor]:
		req.Response <- Response[[]ScoreEngineDescriptor]{Value: mngr.listScoreEnginesByContest(req.Args.contestID)}

		close(req.Response)
	case Request[startScoreEngineArguments, domain.ScoreEngineInstanceID]:
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		instanceID, err := mngr.startScoreEngine(ctx, req.Args.contestID, req.Args.terminatedBy)

		req.Response <- Response[domain.ScoreEngineInstanceID]{
			Value: instanceID,
			Err:   err,
		}

		close(req.Response)
	case Request[stopScoreEngineArguments, struct{}]:
		mngr.stopScoreEngine(req.Args.instanceID)

		close(req.Response)
	case Request[getScoreEngineArguments, ScoreEngineDescriptor]:
		descriptor, err := mngr.getScoreEngine(req.Args.instanceID)
		req.Response <- Response[ScoreEngineDescriptor]{
			Value: descriptor,
			Err:   err,
		}

		close(req.Response)
	}
}

func (mngr *ScoreEngineManager) runPeriodicCheck(ctx context.Context) {
	now := time.Now()
	contests, err := mngr.repo.GetContestsCurrentlyRunningOrByStartTime(ctx, nil, now, now.Add(time.Hour))
	if err != nil {
		slog.Error("score engine manager failed to complete periodic check", "error", err)

		return
	}

	for contest := range slices.Values(contests) {
		if contest.TimeEnd == nil {
			continue
		}

		if handler, ok := mngr.handlers[contest.ID]; ok {
			logger := slog.New(slog.Default().Handler()).
				With("contest_id", contest.ID).
				With("instance_id", handler.instanceID)

			if contest.QualifyingProblems != handler.qualifyingProblems {
				logger.Info("updating scoring rules", "qualifying_problems", contest.QualifyingProblems)
				handler.driver.SetScoringRules(&HardestProblems{Number: contest.QualifyingProblems})
				handler.qualifyingProblems = contest.QualifyingProblems
			}

			if contest.Finalists != handler.finalists {
				logger.Info("updating ranker", "finalists", contest.Finalists)
				handler.driver.SetRanker(NewBasicRanker(contest.Finalists))
				handler.finalists = contest.Finalists
			}

			continue
		}

		_, _ = mngr.startScoreEngine(ctx, contest.ID, (*contest.TimeEnd).Add(contest.GracePeriod).Add(time.Hour))
	}
}

func (mngr *ScoreEngineManager) startScoreEngine(ctx context.Context, contestID domain.ContestID, terminateBy time.Time) (domain.ScoreEngineInstanceID, error) {
	if _, ok := mngr.handlers[contestID]; ok {
		return uuid.Nil, errors.New(ErrAlreadyStarted)
	}

	now := time.Now()

	contest, err := mngr.repo.GetContest(ctx, nil, contestID)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, 0)
	}

	logger := slog.New(slog.Default().Handler()).With("contest_id", contestID)

	logger = logger.With(slog.Group("config",
		"qualifying_problems", contest.QualifyingProblems,
		"finalists", contest.Finalists)).
		With("terminated_by", terminateBy)

	if contest.TimeBegin != nil && contest.TimeBegin.After(now) {
		logger = logger.With("starting_in", time.Until(*contest.TimeBegin))
	}

	logger.Info("spinning up score engine")

	instanceID := uuid.New()
	store := NewMemoryStore()
	rules := &HardestProblems{Number: contest.QualifyingProblems}
	ranker := NewBasicRanker(contest.Finalists)
	driver := NewScoreEngineDriver(contest.ID, instanceID, mngr.eventBroker)
	engine := NewDefaultScoreEngine(ranker, rules, store)

	cancellableCtx, stop := context.WithDeadline(context.Background(), terminateBy)
	wg, installEngine := driver.Run(cancellableCtx)

	hydrationStartTime := time.Now()
	err = mngr.engineStoreHydrator.Hydrate(ctx, contestID, store)
	if err != nil {
		logger.Error("hydration failed", "error", err)

		stop()

		return uuid.Nil, errors.Wrap(err, 0)
	}

	logger.Debug("score engine store hydration complete", "time", time.Since(hydrationStartTime))

	installEngine(engine)

	mngr.handlers[contestID] = &engineHandler{
		instanceID:         instanceID,
		driver:             driver,
		stop:               stop,
		wg:                 wg,
		finalists:          contest.Finalists,
		qualifyingProblems: contest.QualifyingProblems,
	}

	go func() {
		wg.Wait()

		mngr.terminations <- instanceID
	}()

	logger.Info("score engine started", "instance_id", instanceID)

	return instanceID, nil
}

func (mngr *ScoreEngineManager) listScoreEnginesByContest(needle domain.ContestID) []ScoreEngineDescriptor {
	instances := make([]ScoreEngineDescriptor, 0)

	for contestID, handler := range mngr.handlers {
		if contestID == needle {
			instances = append(instances, ScoreEngineDescriptor{
				InstanceID: handler.instanceID,
				ContestID:  contestID,
			})
		}
	}

	return instances
}

func (mngr *ScoreEngineManager) stopScoreEngine(instanceID domain.ScoreEngineInstanceID) {
	for contestID, handler := range mngr.handlers {
		if handler.instanceID == instanceID {
			handler.stop()
			handler.wg.Wait()

			delete(mngr.handlers, contestID)

			return
		}
	}
}

func (mngr *ScoreEngineManager) getScoreEngine(instanceID domain.ScoreEngineInstanceID) (ScoreEngineDescriptor, error) {
	for contestID, handler := range mngr.handlers {
		if handler.instanceID == instanceID {
			return ScoreEngineDescriptor{
				InstanceID: handler.instanceID,
				ContestID:  contestID,
			}, nil
		}
	}

	return ScoreEngineDescriptor{}, errors.Wrap(domain.ErrNotFound, 0)
}
