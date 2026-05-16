package rest

import (
	"context"
	"net/http"
	"runtime/debug"

	"github.com/climblive/platform/backend/internal/domain"
)

var Version = ""

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
	version, found := getVersion()
	if !found {
		version = "dev"
	}

	writeResponse(w, http.StatusOK, version)
}

func getVersion() (string, bool) {
	if Version != "" {
		return Version, true
	}

	buildInfo, found := debug.ReadBuildInfo()
	if !found {
		return "", false
	}

	for _, setting := range buildInfo.Settings {
		if setting.Key == "vcs.revision" {
			if len(setting.Value) <= 7 {
				return setting.Value, true
			}

			return setting.Value[:7], true
		}
	}

	return "", false
}
