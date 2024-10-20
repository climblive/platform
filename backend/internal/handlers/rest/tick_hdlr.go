package rest

import (
	"encoding/json"
	"net/http"

	"github.com/climblive/platform/backend/internal/domain"
)

type tickHandler struct {
	tickUseCase domain.TickUseCase
}

func InstallTickHandler(mux *Mux, tickUseCase domain.TickUseCase) {
	handler := &tickHandler{
		tickUseCase: tickUseCase,
	}

	mux.HandleFunc("GET /contenders/{contenderID}/ticks", handler.GetTicksByContender)
	mux.HandleFunc("POST /contenders/{contenderID}/ticks", handler.CreateTick)
	mux.HandleFunc("DELETE /ticks/{tickID}", handler.DeleteTick)
}

func (hdlr *tickHandler) GetTicksByContender(w http.ResponseWriter, r *http.Request) {
	contenderID := parseResourceID[domain.ContenderID](r.PathValue("contenderID"))

	ticks, err := hdlr.tickUseCase.GetTicksByContender(r.Context(), contenderID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, ticks)
}

func (hdlr *tickHandler) CreateTick(w http.ResponseWriter, r *http.Request) {
	contenderID := parseResourceID[domain.ContenderID](r.PathValue("contenderID"))

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

func (hdlr *tickHandler) DeleteTick(w http.ResponseWriter, r *http.Request) {
	tickID := parseResourceID[domain.TickID](r.PathValue("tickID"))

	err := hdlr.tickUseCase.DeleteTick(r.Context(), tickID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusNoContent, nil)
}
