package rest

import (
	"encoding/json"
	"net/http"

	"github.com/climblive/platform/backend/internal/domain"
)

type tickHandler struct {
	tickUseCase domain.TickUseCase
}

func InstallTickHandler(tickUseCase domain.TickUseCase) {
	handler := &tickHandler{
		tickUseCase: tickUseCase,
	}

	http.HandleFunc("GET /contenders/{contenderID}/ticks", handler.GetTicksByContender)
	http.HandleFunc("POST /contenders/{contenderID}/ticks", handler.CreateTick)
}

func (hdlr *tickHandler) GetTicksByContender(w http.ResponseWriter, r *http.Request) {
	contenderID := parseResourceID(r.PathValue("contenderID"))

	ticks, err := hdlr.tickUseCase.GetTicksByContender(r.Context(), contenderID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, ticks)
}

func (hdlr *tickHandler) CreateTick(w http.ResponseWriter, r *http.Request) {
	contenderID := parseResourceID(r.PathValue("contenderID"))

	var tick domain.Tick
	err := json.NewDecoder(r.Body).Decode(&tick)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	createdTick, err := hdlr.tickUseCase.CreateTick(r.Context(), contenderID, tick)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusCreated, createdTick)
}
