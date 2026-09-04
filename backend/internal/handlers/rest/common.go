package rest

import (
	jsonv1 "encoding/json"
	"encoding/json/v2"
	"fmt"
	"log/slog"
	"mime"
	"net/http"
	"strconv"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/utils"
	"github.com/go-errors/errors"
)

const maxJSONBodySize = 1 << 20

func parseResourceID[T domain.ResourceIDType](id string) (T, error) {
	number, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		var empty T
		return empty, err
	}

	return T(number), nil
}

func writeResponse(w http.ResponseWriter, status int, data any) {
	if data == nil {
		w.WriteHeader(status)
		return
	}

	json, err := json.Marshal(data, jsonv1.FormatDurationAsNano(true))
	if err != nil {
		handleError(w, errors.Wrap(err, 0))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	_, err = w.Write(json)
	if err != nil {
		slog.Error("failed to write http response", "error", err)
	}
}

func readJSON(w http.ResponseWriter, r *http.Request, out any, opts ...json.Options) bool {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return false
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxJSONBodySize)
	err = json.UnmarshalRead(r.Body, out, append(opts, json.RejectUnknownMembers(true))...)
	if err == nil {
		return true
	}

	var maxBytesError *http.MaxBytesError
	if errors.As(err, &maxBytesError) {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return false
	}

	w.WriteHeader(http.StatusBadRequest)
	return false
}

func handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrAllWinnersDrawn):
		fallthrough
	case errors.Is(err, domain.ErrNotFound):
		fallthrough
	case errors.Is(err, domain.ErrArchived):
		w.WriteHeader(http.StatusNotFound)
	case errors.Is(err, domain.ErrDuplicate):
		w.WriteHeader(http.StatusConflict)
	case errors.Is(err, domain.ErrNotAuthenticated):
		fallthrough
	case errors.Is(err, domain.ErrNotAuthorized):
		fallthrough
	case errors.Is(err, domain.ErrNoOwnership):
		fallthrough
	case errors.Is(err, domain.ErrContestEnded):
		fallthrough
	case errors.Is(err, domain.ErrInsufficientRole):
		fallthrough
	case errors.Is(err, domain.ErrNotAllowed):
		w.WriteHeader(http.StatusForbidden)
	case errors.Is(err, domain.ErrLimitExceeded):
		w.WriteHeader(http.StatusConflict)
	case errors.Is(err, domain.ErrInvalidData):
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)

		if stack := utils.GetErrorStack(err); stack != "" {
			fmt.Println(stack)
		}
	}
}
