package rest

import (
	"context"
	"net/http"

	"github.com/climblive/platform/backend/internal/domain"
)

type healthUseCase interface {
	GetHealth(ctx context.Context) (domain.HealthStatus, error)
}

type healthHandler struct {
	healthUseCase healthUseCase
}

func InstallHealthHandler(mux *Mux, healthUseCase healthUseCase) {
	handler := &healthHandler{
		healthUseCase: healthUseCase,
	}

	mux.HandleFunc("GET /health", handler.GetHealth)
}

func (hdlr *healthHandler) GetHealth(w http.ResponseWriter, r *http.Request) {
	health, err := hdlr.healthUseCase.GetHealth(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	status := http.StatusOK
	if !health.ScoreEngineManager.Healthy || !health.ScoreKeeper.Healthy || !health.Scrubber.Healthy {
		status = http.StatusServiceUnavailable
	}

	writeResponse(w, status, health)
}
