package rest

import (
	"context"
	"net/http"

	"github.com/climblive/platform/backend/internal/domain"
)

type tickUseCase interface {
	GetTicksByContender(ctx context.Context, contenderID domain.ContenderID) ([]domain.Tick, error)
	GetTicksByCompetition(ctx context.Context, competitionID domain.CompetitionID) ([]domain.Tick, error)
	DeleteTick(ctx context.Context, tickID domain.TickID) error
	PutTick(ctx context.Context, contenderID domain.ContenderID, tick domain.Tick) (domain.Tick, error)
}

type tickHandler struct {
	tickUseCase tickUseCase
}

func InstallTickHandler(mux *Mux, tickUseCase tickUseCase) {
	handler := &tickHandler{
		tickUseCase: tickUseCase,
	}

	mux.HandleFunc("GET /contenders/{contenderID}/ticks", handler.GetTicksByContender)
	mux.HandleFunc("GET /competitions/{competitionID}/ticks", handler.GetTicksByCompetition)
	mux.HandleFunc("PUT /contenders/{contenderID}/ticks", handler.PutTick)
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

func (hdlr *tickHandler) GetTicksByCompetition(w http.ResponseWriter, r *http.Request) {
	competitionID, err := parseResourceID[domain.CompetitionID](r.PathValue("competitionID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ticks, err := hdlr.tickUseCase.GetTicksByCompetition(r.Context(), competitionID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, ticks)
}

func (hdlr *tickHandler) PutTick(w http.ResponseWriter, r *http.Request) {
	contenderID, err := parseResourceID[domain.ContenderID](r.PathValue("contenderID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var tick domain.Tick
	if !readJSON(w, r, &tick) {
		return
	}

	updatedTick, err := hdlr.tickUseCase.PutTick(r.Context(), contenderID, tick)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, updatedTick)
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
