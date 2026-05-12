package rest

import (
	"context"
	"net/http"
	"runtime"
	"runtime/debug"

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
	mux.HandleFunc("GET /health/ok", handler.GetHealthOk)
	mux.HandleFunc("GET /version", handler.GetVersion)
}

func (hdlr *healthHandler) GetHealth(w http.ResponseWriter, r *http.Request) {
	health, err := hdlr.healthUseCase.GetHealth(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, health)
}

func (hdlr *healthHandler) GetHealthOk(w http.ResponseWriter, r *http.Request) {
	health, err := hdlr.healthUseCase.GetHealth(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	status := http.StatusOK

	for _, service := range health {
		if !service.Healthy {
			status = http.StatusServiceUnavailable
			break
		}
	}

	writeResponse(w, status, nil)
}

func (hdlr *healthHandler) GetVersion(w http.ResponseWriter, _ *http.Request) {
	if buildInfo, ok := debug.ReadBuildInfo(); ok {
		writeResponse(w, http.StatusOK, buildInfo.GoVersion)
		return
	}

	writeResponse(w, http.StatusOK, runtime.Version())
}
