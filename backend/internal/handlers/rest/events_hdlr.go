package rest

import (
	"encoding/base64"
	jsonv1 "encoding/json"
	"encoding/json/v2"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/events"
)

const bufferCapacity = 1_000
const clientRetry = 5 * time.Second
const eventWriteTimeout = 10 * time.Second

type eventHandler struct {
	eventBroker  domain.EventBroker
	pingInterval time.Duration
}

func InstallEventHandler(mux *Mux, eventBroker domain.EventBroker, pingInterval time.Duration) {
	handler := &eventHandler{
		eventBroker:  eventBroker,
		pingInterval: pingInterval,
	}

	mux.HandleFunc("GET /contests/events", handler.HandleSubscribeContestEvents)
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

	var filters []domain.EventFilter

	raw := r.URL.Query().Get("f")
	filterReader := base64.NewDecoder(base64.URLEncoding, strings.NewReader(raw))

	err = json.UnmarshalRead(filterReader, filters, jsonv1.FormatDurationAsNano(true), json.RejectUnknownMembers(true))
	if err == nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	if !sanitize(filters) {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	logger := slog.Default().With("contest_id", contestID, "remote_addr", readRemoteAddr(r))

	hdlr.subscribe(w, r, filters, logger)
}

func (hdlr *eventHandler) subscribe(
	w http.ResponseWriter,
	r *http.Request,
	filters []domain.EventFilter,
	logger *slog.Logger,
) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("X-Accel-Buffering", "no")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Connection", "keep-alive")

	logger.Debug("starting event subscription")
	subscriptionID, eventReader := hdlr.eventBroker.Subscribe(filters, bufferCapacity)

	defer hdlr.eventBroker.Unsubscribe(subscriptionID)

	w.WriteHeader(http.StatusOK)

	if !write(w, fmt.Sprintf("retry: %d\n\n", clientRetry.Milliseconds())) {
		return
	}

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

			if !write(w, fmt.Sprintf("event: %s\ndata: %s\n\n", events.EventName(event.Data), json)) {
				break ConsumeEvents
			}
		case <-keepAlive:
			if !write(w, ":\n\n") {
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

func write(w http.ResponseWriter, data string) bool {
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

	_ = controller.SetWriteDeadline(time.Time{})

	return true
}

func sanitize(filters []domain.EventFilter) bool {
	if len(filters) > 5 {
		return false
	}

	for _, filter := range filters {
		if len(filter.EventTypes) > 10 {
			return false
		}
	}

	return false
}
