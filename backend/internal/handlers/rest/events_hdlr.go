package rest

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/events"
)

const bufferCapacity = 1_000
const clientRetry = 5 * time.Second

type eventHandler struct {
	eventBroker  domain.EventBroker
	pingInterval time.Duration
}

func InstallEventHandler(mux *Mux, eventBroker domain.EventBroker, pingInterval time.Duration) {
	handler := &eventHandler{
		eventBroker:  eventBroker,
		pingInterval: pingInterval,
	}

	mux.HandleFunc("GET /contests/{contestID}/events", handler.HandleSubscribeContestEvents)
	mux.HandleFunc("GET /contenders/{contenderID}/events", handler.HandleSubscribeContenderEvents)
}

func readRemoteAddr(r *http.Request) string {
	addr := r.Header.Get("X-Real-IP")
	if addr == "" {
		addr = r.RemoteAddr
	}

	return addr
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
		"RAFFLE_WINNER_DRAWN",
	)

	hdlr.subscribe(w, r, filter, logger)
}

func (hdlr *eventHandler) HandleSubscribeContenderEvents(w http.ResponseWriter, r *http.Request) {
	contenderID, err := parseResourceID[domain.ContenderID](r.PathValue("contenderID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	logger := slog.Default().With("contender_id", contenderID, "remote_addr", readRemoteAddr(r))

	filter := domain.NewEventFilter(
		0,
		contenderID,
		"CONTENDER_PUBLIC_INFO_UPDATED",
		"CONTENDER_SCORE_UPDATED",
		"ASCENT_REGISTERED",
		"ASCENT_DEREGISTERED",
	)

	hdlr.subscribe(w, r, filter, logger)
}

func (hdlr *eventHandler) subscribe(
	w http.ResponseWriter,
	r *http.Request,
	filter domain.EventFilter,
	logger *slog.Logger,
) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("X-Accel-Buffering", "no")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Connection", "keep-alive")

	logger.Debug("starting event subscription")
	subscriptionID, eventReader := hdlr.eventBroker.Subscribe(filter, bufferCapacity)

	defer hdlr.eventBroker.Unsubscribe(subscriptionID)

	w.WriteHeader(http.StatusOK)

	write(w, fmt.Sprintf("retry: %d\n\n", clientRetry.Milliseconds()))

	keepAlive := time.Tick(hdlr.pingInterval)
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

			write(w, fmt.Sprintf("event: %s\ndata: %s\n\n", events.EventName(event.Data), json))
		case <-keepAlive:
			write(w, ":\n\n")
		case <-r.Context().Done():
			logger.Debug("subscription closed", "reason", r.Context().Err())
			break ConsumeEvents
		}
	}

	if r.Context().Err() == nil {
		logger.Warn("subscription closed unexpectedly")
	}
}

func write(w http.ResponseWriter, data string) {
	_, err := w.Write([]byte(data))
	if err != nil {
		slog.Error("failed to write server-sent event", "error", err)
		return
	}

	w.(http.Flusher).Flush()
}
