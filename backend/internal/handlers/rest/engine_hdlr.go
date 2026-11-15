package rest

import (
	"context"
	"net/http"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"uuid"
)

type scoreEngineUseCase interface {
	ListScoreEngines(ctx context.Context) ([]domain.ScoreEngineDescriptor, error)
	ListScoreEnginesByCompetition(ctx context.Context, competitionID domain.CompetitionID) ([]domain.ScoreEngineDescriptor, error)
	StopScoreEngine(ctx context.Context, instanceID domain.ScoreEngineInstanceID) error
	StartScoreEngine(ctx context.Context, competitionID domain.CompetitionID, terminatedBy time.Time) (domain.ScoreEngineInstanceID, error)
}

type scoreEngineHandler struct {
	scoreEngineUseCase scoreEngineUseCase
}

func InstallScoreEngineHandler(mux *Mux, scoreEngineUseCase scoreEngineUseCase) {
	handler := &scoreEngineHandler{
		scoreEngineUseCase: scoreEngineUseCase,
	}

	mux.HandleFunc("GET /score-engines", handler.ListScoreEngines)
	mux.HandleFunc("GET /competitions/{competitionID}/score-engines", handler.ListScoreEnginesByCompetition)
	mux.HandleFunc("DELETE /score-engines/{instanceID}", handler.StopScoreEngine)
	mux.HandleFunc("POST /competitions/{competitionID}/score-engines", handler.StartScoreEngine)
}

func (hdlr *scoreEngineHandler) ListScoreEngines(w http.ResponseWriter, r *http.Request) {
	instances, err := hdlr.scoreEngineUseCase.ListScoreEngines(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, instances)
}

func (hdlr *scoreEngineHandler) ListScoreEnginesByCompetition(w http.ResponseWriter, r *http.Request) {
	competitionID, err := parseResourceID[domain.CompetitionID](r.PathValue("competitionID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	instances, err := hdlr.scoreEngineUseCase.ListScoreEnginesByCompetition(r.Context(), competitionID)
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
	competitionID, err := parseResourceID[domain.CompetitionID](r.PathValue("competitionID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	var arguments StartScoreEngineArguments
	if !readJSON(w, r, &arguments) {
		return
	}

	instanceID, err := hdlr.scoreEngineUseCase.StartScoreEngine(r.Context(), competitionID, arguments.TerminatedBy)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusCreated, instanceID)
}
