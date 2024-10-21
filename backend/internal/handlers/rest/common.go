package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/utils"
	"github.com/go-errors/errors"
)

func parseResourceID[T domain.ResourceIDType](id string) T {
	number, _ := strconv.Atoi(id)
	return T(number)
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
	w.Write(json)
}

func handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		w.WriteHeader(http.StatusNotFound)
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
