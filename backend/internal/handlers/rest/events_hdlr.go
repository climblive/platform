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

	mux.HandleFunc("GET /contests/{contestID}/events", handler.ListenContestEvents)
}

func (hdlr *eventHandler) ListenContestEvents(w http.ResponseWriter, r *http.Request) {
	contestID := parseResourceID[domain.ContestID](r.PathValue("contestID"))

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("X-Accel-Buffering", "no")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	filter := domain.EventFilter{
		ContestID: contestID,
	}

	slog.Info("starting event subscription", "contest_id", contestID, "remote_addr", r.RemoteAddr)
	subscriptionID, eventReader := hdlr.eventBroker.Subscribe(filter, bufferCapacity)

	defer hdlr.eventBroker.Unsubscribe(subscriptionID)

	w.WriteHeader(http.StatusOK)
	w.(http.Flusher).Flush()

	for {
		event, err := eventReader.AwaitEvent(r.Context())
		if err != nil {
			switch {
			case errors.Is(err, context.Canceled):
			case errors.Is(err, context.DeadlineExceeded):
			default:
				slog.Warn("subscription closed unexpectedly",
					"contest_id", contestID,
					"remote_addr", r.RemoteAddr,
					"error", err)
			}

			return
		}

		json, err := json.Marshal(event.Data)
		if err != nil {
			panic(err)
		}

		w.Write([]byte(fmt.Sprintf("event: %s\n", event.Name)))
		w.Write([]byte(fmt.Sprintf("data: %s\n\n", json)))

		w.(http.Flusher).Flush()
	}
}
