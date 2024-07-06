package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/climblive/platform/backend/internal/domain"
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

	w.WriteHeader(status)
	w.Write(json)
}

func writeError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
}
