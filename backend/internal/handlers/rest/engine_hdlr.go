package rest

import (
	"context"
	"net/http"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/google/uuid"
)

type scoreEngineManager interface {
	ListScoreEnginesByContest(ctx context.Context, contestID domain.ContestID) ([]domain.ScoreEngineInstanceID, error)
	StopScoreEngine(ctx context.Context, instanceID domain.ScoreEngineInstanceID) error
	StartScoreEngine(ctx context.Context, contestID domain.ContestID) (domain.ScoreEngineInstanceID, error)
}

type engineHandler struct {
	scoreEngineManager scoreEngineManager
}

func InstallScoreEngineHandler(mux *Mux, scoreEngineManager scoreEngineManager) {
	handler := &engineHandler{
		scoreEngineManager: scoreEngineManager,
	}

	mux.HandleFunc("GET /contests/{contestID}/score-engines", handler.ListScoreEnginesByContest)
	mux.HandleFunc("DELETE /score-engines/{instanceID}", handler.StopScoreEngine)
	mux.HandleFunc("POST /contests/{contestID}/score-engines", handler.StartScoreEngine)
}

func (hdlr *engineHandler) ListScoreEnginesByContest(w http.ResponseWriter, r *http.Request) {
	contestID := parseResourceID[domain.ContestID](r.PathValue("contestID"))

	instances, err := hdlr.scoreEngineManager.ListScoreEnginesByContest(r.Context(), contestID)
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

	err = hdlr.scoreEngineManager.StopScoreEngine(r.Context(), instanceID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusAccepted, nil)
}

func (hdlr *engineHandler) StartScoreEngine(w http.ResponseWriter, r *http.Request) {
	contestID := parseResourceID[domain.ContestID](r.PathValue("contestID"))

	instanceID, err := hdlr.scoreEngineManager.StartScoreEngine(r.Context(), contestID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, instanceID)
}
