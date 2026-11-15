package rest

import (
	"context"
	"net/http"

	"github.com/climblive/platform/backend/internal/domain"
)

type tickUseCase interface {
	GetTicksByCompetitor(ctx context.Context, competitorID domain.CompetitorID) ([]domain.Tick, error)
	GetTicksByCompetition(ctx context.Context, competitionID domain.CompetitionID) ([]domain.Tick, error)
	DeleteTick(ctx context.Context, tickID domain.TickID) error
	PutTick(ctx context.Context, competitorID domain.CompetitorID, tick domain.Tick) (domain.Tick, error)
}

type tickHandler struct {
	tickUseCase tickUseCase
}

func InstallTickHandler(mux *Mux, tickUseCase tickUseCase) {
	handler := &tickHandler{
		tickUseCase: tickUseCase,
	}

	mux.HandleFunc("GET /competitors/{competitorID}/ticks", handler.GetTicksByCompetitor)
	mux.HandleFunc("GET /competitions/{competitionID}/ticks", handler.GetTicksByCompetition)
	mux.HandleFunc("PUT /competitors/{competitorID}/ticks", handler.PutTick)
	mux.HandleFunc("DELETE /ticks/{tickID}", handler.DeleteTick)
}

func (hdlr *tickHandler) GetTicksByCompetitor(w http.ResponseWriter, r *http.Request) {
	competitorID, err := parseResourceID[domain.CompetitorID](r.PathValue("competitorID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ticks, err := hdlr.tickUseCase.GetTicksByCompetitor(r.Context(), competitorID)
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
	competitorID, err := parseResourceID[domain.CompetitorID](r.PathValue("competitorID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var tick domain.Tick
	if !readJSON(w, r, &tick) {
		return
	}

	updatedTick, err := hdlr.tickUseCase.PutTick(r.Context(), competitorID, tick)
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
