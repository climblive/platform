package rest

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/events"
)

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

	events := make(chan domain.EventContainer, events.EventChannelBufferSize)

	subscriptionID := hdlr.eventBroker.Subscribe(domain.EventFilter{
		ContestID: contestID,
	}, events)
	slog.Info("start event subscription", "contest_id", contestID, "remote_addr", r.RemoteAddr)

	defer hdlr.eventBroker.Unsubscribe(subscriptionID)

	w.WriteHeader(http.StatusOK)
	w.(http.Flusher).Flush()

	for {
		select {
		case event := <-events:
			json, err := json.Marshal(event.Data)
			if err != nil {
				return
			}

			w.Write([]byte(fmt.Sprintf("event: %s\n", event.Name)))
			w.Write([]byte(fmt.Sprintf("data: %s\n\n", json)))

			w.(http.Flusher).Flush()
		case <-r.Context().Done():
			return
		}
	}
}
