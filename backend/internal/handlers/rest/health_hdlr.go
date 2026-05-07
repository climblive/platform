package rest

import (
	"context"
	"net/http"

	"github.com/climblive/platform/backend/internal/domain"
)

type healthUseCase interface {
	GetHealth(ctx context.Context) ([]domain.ServiceStatus, error)
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

	issues := 0
	for _, service := range health {
		if !service.Healthy {
			issues++
		}
	}

	if issues > 0 {
		status = http.StatusServiceUnavailable
	}

	writeResponse(w, status, health)
}
