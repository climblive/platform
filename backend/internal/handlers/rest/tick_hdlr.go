package rest

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/climblive/platform/backend/internal/domain"
)

type tickUseCase interface {
	GetTicksByContender(ctx context.Context, contenderID domain.ContenderID) ([]domain.Tick, error)
	DeleteTick(ctx context.Context, tickID domain.TickID) error
	CreateTick(ctx context.Context, contenderID domain.ContenderID, tick domain.Tick) (domain.Tick, error)
}

type tickHandler struct {
	tickUseCase tickUseCase
}

func InstallTickHandler(mux *Mux, tickUseCase tickUseCase) {
	handler := &tickHandler{
		tickUseCase: tickUseCase,
	}

	mux.HandleFunc("GET /contenders/{contenderID}/ticks", handler.GetTicksByContender)
	mux.HandleFunc("POST /contenders/{contenderID}/ticks", handler.CreateTick)
	mux.HandleFunc("DELETE /ticks/{tickID}", handler.DeleteTick)
}

func (hdlr *tickHandler) GetTicksByContender(w http.ResponseWriter, r *http.Request) {
	contenderID, err := parseResourceID[domain.ContenderID](r.PathValue("contenderID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ticks, err := hdlr.tickUseCase.GetTicksByContender(r.Context(), contenderID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, ticks)
}

func (hdlr *tickHandler) CreateTick(w http.ResponseWriter, r *http.Request) {
	contenderID, err := parseResourceID[domain.ContenderID](r.PathValue("contenderID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var tick domain.Tick
	err = json.NewDecoder(r.Body).Decode(&tick)
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
	tickID, err := parseResourceID[domain.TickID](r.PathValue("tickID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = hdlr.tickUseCase.DeleteTick(r.Context(), tickID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusNoContent, nil)
}
