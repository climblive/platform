package rest

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
)

const bufferCapacity = 1_000
const clientRetry = 5 * time.Second

type eventHandler struct {
	eventBroker domain.EventBroker
}

func InstallEventHandler(mux *Mux, eventBroker domain.EventBroker) {
	handler := &eventHandler{
		eventBroker: eventBroker,
	}

	mux.HandleFunc("GET /contests/{contestID}/events", handler.HandleSubscribeContestEvents)
	mux.HandleFunc("GET /contenders/{contenderID}/events", handler.HandleSubscribeContenderEvents)
}

func (hdlr *eventHandler) HandleSubscribeContestEvents(w http.ResponseWriter, r *http.Request) {
	contestID := parseResourceID[domain.ContestID](r.PathValue("contestID"))

	logger := slog.Default().With("contest_id", contestID, "remote_addr", r.RemoteAddr)

	filter := domain.NewEventFilter(
		contestID,
		0,
		"CONTENDER_PUBLIC_INFO_UPDATED",
		"[]CONTENDER_SCORE_UPDATED",
	)

	hdlr.subscribe(w, r, filter, logger)
}

func (hdlr *eventHandler) HandleSubscribeContenderEvents(w http.ResponseWriter, r *http.Request) {
	contenderID := parseResourceID[domain.ContenderID](r.PathValue("contenderID"))

	logger := slog.Default().With("contender_id", contenderID, "remote_addr", r.RemoteAddr)

	filter := domain.NewEventFilter(
		0,
		contenderID,
		"CONTENDER_PUBLIC_INFO_UPDATED",
		"CONTENDER_SCORE_UPDATED",
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
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	write(w, fmt.Sprintf("retry: %d\n\n", clientRetry.Milliseconds()))

	logger.Info("starting event subscription")
	subscriptionID, eventReader := hdlr.eventBroker.Subscribe(filter, bufferCapacity)

	defer hdlr.eventBroker.Unsubscribe(subscriptionID)

	w.WriteHeader(http.StatusOK)
	w.(http.Flusher).Flush()

	keepAlive := time.Tick(10 * time.Second)
	events := eventReader.EventsChan(r.Context())

ConsumeEvents:
	for {
		select {
		case event, open := <-events:
			if !open {
				break ConsumeEvents
			}

			json, err := json.Marshal(event.Data)
			if err != nil {
				panic(err)
			}

			write(w, fmt.Sprintf("event: %s\ndata: %s\n\n", event.Name, json))
		case <-keepAlive:
			write(w, ":\n\n")
		case <-r.Context().Done():
			logger.Info("subscription closed", "reason", r.Context().Err())
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
