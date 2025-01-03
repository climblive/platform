package rest

import (
	"context"
	"net/http"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/google/uuid"
)

type scoreEngineController interface {
	ListScoreEnginesByContest(ctx context.Context, contestID domain.ContestID) ([]domain.ScoreEngineInstanceID, error)
	StopScoreEngine(ctx context.Context, instanceID domain.ScoreEngineInstanceID) error
	StartScoreEngine(ctx context.Context, contestID domain.ContestID) (domain.ScoreEngineInstanceID, error)
}

type engineHandler struct {
	controller scoreEngineController
}

func InstallScoreEngineHandler(mux *Mux, controller scoreEngineController) {
	handler := &engineHandler{
		controller: controller,
	}

	mux.HandleFunc("GET /contests/{contestID}/score-engines", handler.ListScoreEnginesByContest)
	mux.HandleFunc("DELETE /score-engines/{instanceID}", handler.StopScoreEngine)
	mux.HandleFunc("POST /contests/{contestID}/score-engines", handler.StartScoreEngine)
}

func (hdlr *engineHandler) ListScoreEnginesByContest(w http.ResponseWriter, r *http.Request) {
	contestID, err := parseResourceID[domain.ContestID](r.PathValue("contestID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	instances, err := hdlr.controller.ListScoreEnginesByContest(r.Context(), contestID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, instances)
}

func (hdlr *engineHandler) StopScoreEngine(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("instanceID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	instanceID := domain.ScoreEngineInstanceID(id)

	err = hdlr.controller.StopScoreEngine(r.Context(), instanceID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusAccepted, nil)
}

func (hdlr *engineHandler) StartScoreEngine(w http.ResponseWriter, r *http.Request) {
	contestID, err := parseResourceID[domain.ContestID](r.PathValue("contestID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	instanceID, err := hdlr.controller.StartScoreEngine(r.Context(), contestID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, instanceID)
}
