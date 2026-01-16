package rest

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/utils"
	"github.com/go-errors/errors"
	"github.com/google/uuid"
)

func parseResourceID[T domain.ResourceIDType](id string) (T, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		var empty T
		return empty, err
	}

	return T(parsed), nil
}

func writeResponse(w http.ResponseWriter, status int, data any) {
	if data == nil {
		w.WriteHeader(status)
		return
	}

	json, err := json.Marshal(data)
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
	case errors.Is(err, domain.ErrInvalidData):
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)

		if stack := utils.GetErrorStack(err); stack != "" {
			fmt.Println(stack)
		}
	}
}
