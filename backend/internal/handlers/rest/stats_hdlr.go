package rest

import (
	"context"
	"net/http"

	"github.com/climblive/platform/backend/internal/domain"
)

type statsUseCase interface {
	GetGlobalStats(ctx context.Context) (domain.GlobalStats, error)
}

type statsHandler struct {
	statsUseCase statsUseCase
}

func InstallStatsHandler(mux *Mux, statsUseCase statsUseCase) {
	handler := &statsHandler{statsUseCase: statsUseCase}
	mux.HandleFunc("GET /stats/global", handler.GetGlobalStats)
}

func (hdlr *statsHandler) GetGlobalStats(w http.ResponseWriter, r *http.Request) {
	stats, err := hdlr.statsUseCase.GetGlobalStats(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, stats)
}
