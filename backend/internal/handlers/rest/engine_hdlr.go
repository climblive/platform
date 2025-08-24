package rest

import (
	"context"
	"encoding/json/v2"
	"net/http"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/google/uuid"
)

type scoreEngineUseCase interface {
	ListScoreEnginesByContest(ctx context.Context, contestID domain.ContestID) ([]domain.ScoreEngineInstanceID, error)
	StopScoreEngine(ctx context.Context, instanceID domain.ScoreEngineInstanceID) error
	StartScoreEngine(ctx context.Context, contestID domain.ContestID, terminatedBy time.Time) (domain.ScoreEngineInstanceID, error)
}

type scoreEngineHandler struct {
	scoreEngineUseCase scoreEngineUseCase
}

func InstallScoreEngineHandler(mux *Mux, scoreEngineUseCase scoreEngineUseCase) {
	handler := &scoreEngineHandler{
		scoreEngineUseCase: scoreEngineUseCase,
	}

	mux.HandleFunc("GET /contests/{contestID}/score-engines", handler.ListScoreEnginesByContest)
	mux.HandleFunc("DELETE /score-engines/{instanceID}", handler.StopScoreEngine)
	mux.HandleFunc("POST /contests/{contestID}/score-engines", handler.StartScoreEngine)
}

func (hdlr *scoreEngineHandler) ListScoreEnginesByContest(w http.ResponseWriter, r *http.Request) {
	contestID, err := parseResourceID[domain.ContestID](r.PathValue("contestID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	instances, err := hdlr.scoreEngineUseCase.ListScoreEnginesByContest(r.Context(), contestID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, instances)
}

func (hdlr *scoreEngineHandler) StopScoreEngine(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("instanceID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	instanceID := domain.ScoreEngineInstanceID(id)

	err = hdlr.scoreEngineUseCase.StopScoreEngine(r.Context(), instanceID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusNoContent, nil)
}

func (hdlr *scoreEngineHandler) StartScoreEngine(w http.ResponseWriter, r *http.Request) {
	contestID, err := parseResourceID[domain.ContestID](r.PathValue("contestID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	var arguments StartScoreEngineArguments
	err = json.UnmarshalRead(r.Body, &arguments)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	instanceID, err := hdlr.scoreEngineUseCase.StartScoreEngine(r.Context(), contestID, arguments.TerminatedBy)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusCreated, instanceID)
}
