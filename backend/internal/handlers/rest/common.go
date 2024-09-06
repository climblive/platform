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
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	w.Write(json)
}

func handleError(w http.ResponseWriter, err error) {
	if stack := utils.GetErrorStack(err); stack != "" {
		fmt.Println(stack)
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
