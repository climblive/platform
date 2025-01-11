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
	response  chan startResponse
}

type startResponse struct {
	instanceID domain.ScoreEngineInstanceID
	err        error
}

type reverseLookupRequest struct {
	instanceID domain.ScoreEngineInstanceID
	response   chan reverseLookupResponse
}

type reverseLookupResponse struct {
	contestID domain.ContestID
	err       error
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
	response := make(chan startResponse, 1)

	mngr.requests <- startRequest{
		contestID: contestID,
		response:  response,
	}

	select {
	case response := <-response:
		return response.instanceID, response.err
	case <-ctx.Done():
		return uuid.UUID{}, ctx.Err()
	}
}

func (mngr *ScoreEngineManager) ReverseLoopupScoreEngine(ctx context.Context, instanceID domain.ScoreEngineInstanceID) (domain.ContestID, error) {
	response := make(chan reverseLookupResponse, 1)

	mngr.requests <- reverseLookupRequest{
		instanceID: instanceID,
		response:   response,
	}

	select {
	case response := <-response:
		return response.contestID, response.err
	case <-ctx.Done():
		return 0, ctx.Err()
	}
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
					slog.Warn("garbage collecting terminated score engine", "instance_id", terminatedInstanceID)
					delete(mngr.handlers, contestID)

					break
				}
			}
		}
	}
}

func (mngr *ScoreEngineManager) handleRequest(request any) {
	switch req := request.(type) {
	case listRequest:
		req.response <- mngr.listScoreEnginesByContest(req.contestID)

		close(req.response)
	case startRequest:
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		instanceID, err := mngr.startScoreEngine(ctx, req.contestID)

		req.response <- startResponse{
			instanceID: instanceID,
			err:        err,
		}

		close(req.response)
	case stopRequest:
		mngr.stopScoreEngine(req.instanceID)

		close(req.response)
	case reverseLookupRequest:
		contestID, err := mngr.reverseLookupInstance(req.instanceID)
		req.response <- reverseLookupResponse{
			contestID: contestID,
			err:       err,
		}

		close(req.response)
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

		_, _ = mngr.startScoreEngine(ctx, contest.ID)
	}
}

func (mngr *ScoreEngineManager) startScoreEngine(ctx context.Context, contestID domain.ContestID) (domain.ScoreEngineInstanceID, error) {
	if _, ok := mngr.handlers[contestID]; ok {
		return uuid.UUID{}, errors.New(ErrAlreadyStarted)
	}

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
	driver := NewScoreEngineDriver(contest.ID, instanceID, mngr.eventBroker)
	engine := NewDefaultScoreEngine(ranker, rules, store)

	cancellableCtx, stop := context.WithCancel(context.Background())
	wg, installEngine := driver.Run(cancellableCtx)

	hydrationStartTime := time.Now()
	err = mngr.engineStoreHydrator.Hydrate(ctx, contestID, store)
	if err != nil {
		logger.Error("hydration failed", "error", err)

		stop()

		return uuid.UUID{}, errors.Wrap(err, 0)
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

func (mngr *ScoreEngineManager) listScoreEnginesByContest(contestID domain.ContestID) []domain.ScoreEngineInstanceID {
	instances := make([]domain.ScoreEngineInstanceID, 0)

	for id, handler := range mngr.handlers {
		if id == contestID {
			instances = append(instances, handler.instanceID)
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

func (mngr *ScoreEngineManager) reverseLookupInstance(instanceID domain.ScoreEngineInstanceID) (domain.ContestID, error) {
	for contestID, handler := range mngr.handlers {
		if handler.instanceID == instanceID {
			return contestID, nil
		}
	}

	return 0, errors.Wrap(domain.ErrNotFound, 0)
}
