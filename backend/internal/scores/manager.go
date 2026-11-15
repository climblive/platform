package scores

import (
	"context"
	"log/slog"
	"maps"
	"slices"
	"sync"
	"sync/atomic"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
	"uuid"
)

var ErrAlreadyStarted = errors.New("already started")

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
	competitionID domain.CompetitionID
}

type listAllScoreEnginesArguments struct{}

type stopScoreEngineArguments struct {
	instanceID domain.ScoreEngineInstanceID
}

type getScoreEngineArguments struct {
	instanceID domain.ScoreEngineInstanceID
}

type startScoreEngineArguments struct {
	competitionID domain.CompetitionID
	terminatedBy  time.Time
}

const pollInterval = 10 * time.Second

type EngineStoreHydrator interface {
	Hydrate(ctx context.Context, competitionID domain.CompetitionID, store EngineStore) error
}

type scoreEngineManagerRepository interface {
	GetCompetitionsCurrentlyRunningOrByStartTime(ctx context.Context, tx domain.Transaction, earliestStartTime, latestStartTime time.Time) ([]domain.Competition, error)
	GetCompetition(ctx context.Context, tx domain.Transaction, competitionID domain.CompetitionID) (domain.Competition, error)
}

type ScoreEngineManager struct {
	repo                   scoreEngineManagerRepository
	engineStoreHydrator    EngineStoreHydrator
	eventBroker            domain.EventBroker
	handlers               map[domain.CompetitionID]*engineHandler
	requests               chan any
	terminations           chan domain.ScoreEngineInstanceID
	scoreEngineMaxLifetime time.Duration
	running                atomic.Bool
}

type engineHandler struct {
	instanceID domain.ScoreEngineInstanceID
	driver     *ScoreEngineDriver
	stop       func()
	wg         *sync.WaitGroup
}

func NewScoreEngineManager(repo scoreEngineManagerRepository, engineStoreHydrator EngineStoreHydrator, eventBroker domain.EventBroker, scoreEngineMaxLifetime time.Duration) ScoreEngineManager {
	return ScoreEngineManager{
		repo:                   repo,
		engineStoreHydrator:    engineStoreHydrator,
		eventBroker:            eventBroker,
		handlers:               make(map[domain.CompetitionID]*engineHandler),
		requests:               make(chan any),
		terminations:           make(chan domain.ScoreEngineInstanceID),
		scoreEngineMaxLifetime: scoreEngineMaxLifetime,
		running:                atomic.Bool{},
	}
}

func (mngr *ScoreEngineManager) Run(ctx context.Context, options ...func(*runOptions)) *sync.WaitGroup {
	config := &runOptions{}
	for _, opt := range options {
		opt(config)
	}

	wg := new(sync.WaitGroup)

	wg.Add(1)

	go func() {
		defer func() {
			if !config.recoverPanics {
				return
			}

			if r := recover(); r != nil {
				slog.Error("score engine manager panicked", "error", r)
			}
		}()

		defer wg.Done()

		mngr.run(ctx)
	}()

	return wg
}

func (mngr *ScoreEngineManager) ListScoreEnginesByCompetition(
	ctx context.Context,
	competitionID domain.CompetitionID,
) ([]domain.ScoreEngineDescriptor, error) {
	request := Request[listScoreEnginesArguments, []domain.ScoreEngineDescriptor]{
		Args:     listScoreEnginesArguments{competitionID: competitionID},
		Response: nil,
	}
	return request.Do(ctx, mngr.requests)
}

func (mngr *ScoreEngineManager) ListScoreEngines(
	ctx context.Context,
) ([]domain.ScoreEngineDescriptor, error) {
	request := Request[listAllScoreEnginesArguments, []domain.ScoreEngineDescriptor]{
		Args:     listAllScoreEnginesArguments{},
		Response: nil,
	}
	return request.Do(ctx, mngr.requests)
}

func (mngr *ScoreEngineManager) StopScoreEngine(
	ctx context.Context,
	instanceID domain.ScoreEngineInstanceID,
) error {
	request := Request[stopScoreEngineArguments, struct{}]{
		Args:     stopScoreEngineArguments{instanceID: instanceID},
		Response: nil,
	}
	_, err := request.Do(ctx, mngr.requests)

	return err
}

func (mngr *ScoreEngineManager) StartScoreEngine(
	ctx context.Context,
	competitionID domain.CompetitionID,
	terminatedBy time.Time,
) (domain.ScoreEngineInstanceID, error) {
	request := Request[startScoreEngineArguments, domain.ScoreEngineInstanceID]{
		Args: startScoreEngineArguments{
			competitionID: competitionID,
			terminatedBy:  terminatedBy,
		},
		Response: nil,
	}
	return request.Do(ctx, mngr.requests)
}

func (mngr *ScoreEngineManager) GetScoreEngine(ctx context.Context, instanceID domain.ScoreEngineInstanceID) (domain.ScoreEngineDescriptor, error) {
	request := Request[getScoreEngineArguments, domain.ScoreEngineDescriptor]{
		Args:     getScoreEngineArguments{instanceID: instanceID},
		Response: nil,
	}
	return request.Do(ctx, mngr.requests)
}

func (mngr *ScoreEngineManager) run(ctx context.Context) {
	mngr.running.Store(true)
	defer mngr.running.Store(false)

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
			for competitionID, handler := range mngr.handlers {
				if handler.instanceID == terminatedInstanceID {
					slog.Info("removing terminated score engine", "instance_id", terminatedInstanceID)
					delete(mngr.handlers, competitionID)

					break
				}
			}
		}
	}
}

func (mngr *ScoreEngineManager) handleRequest(request any) {
	switch req := request.(type) {
	case Request[listScoreEnginesArguments, []domain.ScoreEngineDescriptor]:
		req.Response <- Response[[]domain.ScoreEngineDescriptor]{
			Value: mngr.listScoreEnginesByCompetition(req.Args.competitionID),
			Err:   nil,
		}

		close(req.Response)
	case Request[listAllScoreEnginesArguments, []domain.ScoreEngineDescriptor]:
		req.Response <- Response[[]domain.ScoreEngineDescriptor]{
			Value: mngr.listScoreEngines(),
			Err:   nil,
		}

		close(req.Response)
	case Request[startScoreEngineArguments, domain.ScoreEngineInstanceID]:
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		instanceID, err := mngr.startScoreEngine(ctx, req.Args.competitionID, req.Args.terminatedBy)

		req.Response <- Response[domain.ScoreEngineInstanceID]{
			Value: instanceID,
			Err:   err,
		}

		close(req.Response)
	case Request[stopScoreEngineArguments, struct{}]:
		mngr.stopScoreEngine(req.Args.instanceID)

		close(req.Response)
	case Request[getScoreEngineArguments, domain.ScoreEngineDescriptor]:
		descriptor, err := mngr.getScoreEngine(req.Args.instanceID)
		req.Response <- Response[domain.ScoreEngineDescriptor]{
			Value: descriptor,
			Err:   err,
		}

		close(req.Response)
	}
}

func (mngr *ScoreEngineManager) runPeriodicCheck(ctx context.Context) {
	now := time.Now()
	competitions, err := mngr.repo.GetCompetitionsCurrentlyRunningOrByStartTime(ctx, nil, now, now.Add(5*time.Minute))
	if err != nil {
		slog.Error("score engine manager failed to complete periodic check", "error", err)

		return
	}

	for competition := range slices.Values(competitions) {
		if competition.TimeEnd.IsZero() {
			continue
		}

		if _, found := mngr.handlers[competition.ID]; found {
			continue
		}

		_, _ = mngr.startScoreEngine(ctx, competition.ID, competition.TimeEnd.Add(12*time.Hour))
	}
}

func (mngr *ScoreEngineManager) startScoreEngine(ctx context.Context, competitionID domain.CompetitionID, terminatedBy time.Time) (domain.ScoreEngineInstanceID, error) {
	if _, ok := mngr.handlers[competitionID]; ok {
		return uuid.Nil(), errors.New(ErrAlreadyStarted)
	}

	now := time.Now()

	competition, err := mngr.repo.GetCompetition(ctx, nil, competitionID)
	if err != nil {
		return uuid.Nil(), errors.Wrap(err, 0)
	}

	logger := slog.New(slog.Default().Handler()).With("competition_id", competitionID)

	latestPermittedTerminationTime := time.Now().Add(mngr.scoreEngineMaxLifetime)

	if terminatedBy.After(latestPermittedTerminationTime) {
		logger.Warn("capping score engine lifetime", "limit", mngr.scoreEngineMaxLifetime, "orig_terminated_by", terminatedBy, "new_terminated_by", latestPermittedTerminationTime)

		terminatedBy = latestPermittedTerminationTime
	}

	logger = logger.With("terminated_by", terminatedBy)

	if competition.TimeBegin.After(now) {
		logger = logger.With("starting_in", time.Until(competition.TimeBegin))
	}

	logger.Info("spinning up score engine")

	instanceID := uuid.New()
	store := NewMemoryStore()
	driver := NewScoreEngineDriver(competition.ID, instanceID, mngr.eventBroker)
	engine := NewDefaultScoreEngine(store)

	cancellableCtx, stop := context.WithDeadline(context.Background(), terminatedBy)
	wg, installEngine := driver.Run(cancellableCtx, WithPanicRecovery())

	hydrationStartTime := time.Now()
	err = mngr.engineStoreHydrator.Hydrate(ctx, competitionID, store)
	if err != nil {
		logger.Error("hydration failed", "error", err)

		stop()

		return uuid.Nil(), errors.Wrap(err, 0)
	}

	logger.Debug("score engine store hydration complete", "time", time.Since(hydrationStartTime))

	installEngine(engine)

	mngr.handlers[competitionID] = &engineHandler{
		instanceID: instanceID,
		driver:     driver,
		stop:       stop,
		wg:         wg,
	}

	go func() {
		wg.Wait()

		mngr.terminations <- instanceID
	}()

	logger.Info("score engine started", "instance_id", instanceID)

	return instanceID, nil
}

func (mngr *ScoreEngineManager) listScoreEngines() []domain.ScoreEngineDescriptor {
	instances := make([]domain.ScoreEngineDescriptor, 0, len(mngr.handlers))

	for competitionID, handler := range mngr.handlers {
		instances = append(instances, domain.ScoreEngineDescriptor{
			InstanceID:    handler.instanceID,
			CompetitionID: competitionID,
		})
	}

	return instances
}

func (mngr *ScoreEngineManager) listScoreEnginesByCompetition(needle domain.CompetitionID) []domain.ScoreEngineDescriptor {
	instances := make([]domain.ScoreEngineDescriptor, 0)

	for competitionID, handler := range mngr.handlers {
		if competitionID == needle {
			instances = append(instances, domain.ScoreEngineDescriptor{
				InstanceID:    handler.instanceID,
				CompetitionID: competitionID,
			})
		}
	}

	return instances
}

func (mngr *ScoreEngineManager) stopScoreEngine(instanceID domain.ScoreEngineInstanceID) {
	for competitionID, handler := range mngr.handlers {
		if handler.instanceID == instanceID {
			handler.stop()
			handler.wg.Wait()

			delete(mngr.handlers, competitionID)

			return
		}
	}
}

func (mngr *ScoreEngineManager) getScoreEngine(instanceID domain.ScoreEngineInstanceID) (domain.ScoreEngineDescriptor, error) {
	for competitionID, handler := range mngr.handlers {
		if handler.instanceID == instanceID {
			return domain.ScoreEngineDescriptor{
				InstanceID:    handler.instanceID,
				CompetitionID: competitionID,
			}, nil
		}
	}

	return domain.ScoreEngineDescriptor{}, errors.Wrap(domain.ErrNotFound, 0)
}

func (mngr *ScoreEngineManager) GetStatus() domain.ServiceStatus {
	return domain.ServiceStatus{Name: "ScoreEngineManager", Healthy: mngr.running.Load(), CheckedAt: time.Now()}
}
