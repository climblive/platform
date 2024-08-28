package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

func parseResourceID(id string) domain.ResourceID {
	number, _ := strconv.Atoi(id)
	return domain.ResourceID(number)
}

func writeResponse(w http.ResponseWriter, status int, data any) {
	if data == nil {
		w.WriteHeader(status)
		return
	}

	json, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	w.Write(json)
}

func handleError(w http.ResponseWriter, err error) {
	if err, ok := err.(*errors.Error); ok {
		fmt.Println(err.ErrorStack())
	}

	switch {
	case errors.Is(err, domain.ErrNotFound):
		w.WriteHeader(http.StatusNotFound)
	case errors.Is(err, domain.ErrPermissionDenied):
		fallthrough
	case errors.Is(err, domain.ErrContestEnded):
		fallthrough
	case errors.Is(err, domain.ErrNotAllowed):
		w.WriteHeader(http.StatusForbidden)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}
