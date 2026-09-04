package rest

import (
	"context"
	"encoding/json/v2"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/events"
)

const eventBufferCapacity = 64
const clientRetry = 5 * time.Second
const eventWriteTimeout = 10 * time.Second

type EventStreamLimits struct {
	Global    int
	PerClient int
}

func defaultEventStreamLimits() EventStreamLimits {
	return EventStreamLimits{Global: 200, PerClient: 5}
}

type eventConnectionLimiter struct {
	mu       sync.Mutex
	limits   EventStreamLimits
	total    int
	byClient map[string]int
}

type eventHandler struct {
	eventBroker  domain.EventBroker
	pingInterval time.Duration
	repo         eventHandlerRepository
	limiter      *eventConnectionLimiter
}

type eventHandlerRepository interface {
	GetContender(ctx context.Context, tx domain.Transaction, contenderID domain.ContenderID) (domain.Contender, error)
}

func InstallEventHandler(mux *Mux, eventBroker domain.EventBroker, repo eventHandlerRepository, pingInterval time.Duration, configuredLimits ...EventStreamLimits) {
	limits := defaultEventStreamLimits()
	if len(configuredLimits) > 0 {
		limits = configuredLimits[0]
	}

	handler := &eventHandler{
		eventBroker:  eventBroker,
		pingInterval: pingInterval,
		repo:         repo,
		limiter:      newEventConnectionLimiter(limits),
	}

	mux.HandleFunc("GET /contests/{contestID}/events", handler.HandleSubscribeContestEvents)
	mux.HandleFunc("GET /contenders/{contenderID}/events", handler.HandleSubscribeContenderEvents)
}

func newEventConnectionLimiter(limits EventStreamLimits) *eventConnectionLimiter {
	return &eventConnectionLimiter{
		mu:       sync.Mutex{},
		limits:   limits,
		total:    0,
		byClient: make(map[string]int),
	}
}

func (l *eventConnectionLimiter) acquire(client string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.total >= l.limits.Global || l.byClient[client] >= l.limits.PerClient {
		return false
	}

	l.total++
	l.byClient[client]++
	return true
}

func (l *eventConnectionLimiter) release(client string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.total--
	l.byClient[client]--
	if l.byClient[client] == 0 {
		delete(l.byClient, client)
	}
}

func readRemoteAddr(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}

	remoteIP := net.ParseIP(host)
	if remoteIP != nil && remoteIP.IsLoopback() {
		realIP := net.ParseIP(r.Header.Get("X-Real-IP"))
		if realIP != nil {
			return realIP.String()
		}
	}

	return host
}

func (hdlr *eventHandler) HandleSubscribeContestEvents(w http.ResponseWriter, r *http.Request) {
	contestID, err := parseResourceID[domain.ContestID](r.PathValue("contestID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	logger := slog.Default().With("contest_id", contestID, "remote_addr", readRemoteAddr(r))

	filter := domain.NewEventFilter(
		contestID,
		0,
		"CONTENDER_PUBLIC_INFO_UPDATED",
		"[]CONTENDER_SCORE_UPDATED",
		"SCORE_ENGINE_STARTED",
		"SCORE_ENGINE_STOPPED",
	)

	hdlr.subscribe(w, r, []domain.EventFilter{filter}, logger)
}

func (hdlr *eventHandler) HandleSubscribeContenderEvents(w http.ResponseWriter, r *http.Request) {
	contenderID, err := parseResourceID[domain.ContenderID](r.PathValue("contenderID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	contender, err := hdlr.repo.GetContender(r.Context(), nil, contenderID)
	if err != nil {
		handleError(w, err)
		return
	}

	logger := slog.Default().With("contender_id", contenderID, "remote_addr", readRemoteAddr(r))

	filters := make([]domain.EventFilter, 0, 2)

	filters = append(filters, domain.NewEventFilter(
		contender.ContestID,
		contenderID,
		"CONTENDER_PUBLIC_INFO_UPDATED",
		"CONTENDER_SCORE_UPDATED",
		"ASCENT_REGISTERED",
		"ASCENT_DEREGISTERED",
		"POINT_VALUE_UPDATED",
		"RAFFLE_WINNER_DRAWN",
	))

	filters = append(filters, domain.NewEventFilter(
		contender.ContestID,
		0,
		"RULES_UPDATED",
		"PROBLEM_ADDED",
		"PROBLEM_UPDATED",
		"PROBLEM_DELETED",
	))

	hdlr.subscribe(w, r, filters, logger)
}

func (hdlr *eventHandler) subscribe(
	w http.ResponseWriter,
	r *http.Request,
	filters []domain.EventFilter,
	logger *slog.Logger,
) {
	client := readRemoteAddr(r)
	if !hdlr.limiter.acquire(client) {
		w.Header().Set("Retry-After", "5")
		w.WriteHeader(http.StatusTooManyRequests)
		return
	}
	defer hdlr.limiter.release(client)

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("X-Accel-Buffering", "no")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Connection", "keep-alive")

	logger.Debug("starting event subscription")
	subscriptionID, eventReader := hdlr.eventBroker.Subscribe(filters, eventBufferCapacity)

	defer hdlr.eventBroker.Unsubscribe(subscriptionID)

	w.WriteHeader(http.StatusOK)

	if !writeEvent(w, fmt.Sprintf("retry: %d\n\n", clientRetry.Milliseconds())) {
		return
	}

	var keepAlive <-chan time.Time
	var keepAliveTicker *time.Ticker
	if hdlr.pingInterval > 0 {
		keepAliveTicker = time.NewTicker(hdlr.pingInterval)
		keepAlive = keepAliveTicker.C
		defer keepAliveTicker.Stop()
	}
	eventsCh := eventReader.EventsChan(r.Context())

ConsumeEvents:
	for {
		select {
		case event, open := <-eventsCh:
			if !open {
				break ConsumeEvents
			}

			json, err := json.Marshal(event.Data)
			if err != nil {
				panic(err)
			}

			if !writeEvent(w, fmt.Sprintf("event: %s\ndata: %s\n\n", events.EventName(event.Data), json)) {
				break ConsumeEvents
			}
		case <-keepAlive:
			if !writeEvent(w, ":\n\n") {
				break ConsumeEvents
			}
		case <-r.Context().Done():
			logger.Debug("subscription closed", "reason", r.Context().Err())
			break ConsumeEvents
		}
	}

	if r.Context().Err() == nil {
		logger.Warn("subscription closed unexpectedly")
	}
}

func writeEvent(w http.ResponseWriter, data string) bool {
	controller := http.NewResponseController(w)
	_ = controller.SetWriteDeadline(time.Now().Add(eventWriteTimeout))

	if _, err := w.Write([]byte(data)); err != nil {
		slog.Error("failed to write server-sent event", "error", err)
		return false
	}

	if err := controller.Flush(); err != nil {
		slog.Error("failed to flush server-sent event", "error", err)
		return false
	}

	return true
}
