package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/climblive/platform/backend/internal/domain"
)

const bufferCapacity = 1_000

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

	logger.Info("starting event subscription")
	subscriptionID, eventReader := hdlr.eventBroker.Subscribe(filter, bufferCapacity)

	defer hdlr.eventBroker.Unsubscribe(subscriptionID)

	w.WriteHeader(http.StatusOK)
	w.(http.Flusher).Flush()

	for {
		event, err := eventReader.AwaitEvent(r.Context())
		if err != nil {
			switch {
			case errors.Is(err, context.Canceled):
				fallthrough
			case errors.Is(err, context.DeadlineExceeded):
				logger.Info("subscription closed")
			default:
				logger.Warn("subscription closed unexpectedly", "error", err)
			}

			return
		}

		json, err := json.Marshal(event.Data)
		if err != nil {
			panic(err)
		}

		_, err = w.Write([]byte(fmt.Sprintf("event: %s\ndata: %s\n\n", event.Name, json)))
		if err != nil {
			slog.Error("failed to write server-sent event", "error", err)
			return
		}

		w.(http.Flusher).Flush()
	}
}
