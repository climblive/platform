package rest

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/scores"
	"github.com/google/uuid"
)

type scoreEngineUseCase interface {
	ListScoreEngines(ctx context.Context) ([]scores.ScoreEngineDescriptor, error)
	ListScoreEnginesByContest(ctx context.Context, contestID domain.ContestID) ([]domain.ScoreEngineInstanceID, error)
	StopScoreEngine(ctx context.Context, instanceID domain.ScoreEngineInstanceID) error
	StartScoreEngine(ctx context.Context, contestID domain.ContestID, terminatedBy time.Time) (domain.ScoreEngineInstanceID, error)
}

type scoreEngineHandler struct {
	scoreEngineUseCase scoreEngineUseCase
}

type runningScoreEngine struct {
	ContestID  domain.ContestID             `json:"contestId"`
	InstanceID domain.ScoreEngineInstanceID `json:"instanceId"`
}

func InstallScoreEngineHandler(mux *Mux, scoreEngineUseCase scoreEngineUseCase) {
	handler := &scoreEngineHandler{
		scoreEngineUseCase: scoreEngineUseCase,
	}

	mux.HandleFunc("GET /score-engines", handler.ListScoreEngines)
	mux.HandleFunc("GET /contests/{contestID}/score-engines", handler.ListScoreEnginesByContest)
	mux.HandleFunc("DELETE /score-engines/{instanceID}", handler.StopScoreEngine)
	mux.HandleFunc("POST /contests/{contestID}/score-engines", handler.StartScoreEngine)
}

func (hdlr *scoreEngineHandler) ListScoreEngines(w http.ResponseWriter, r *http.Request) {
	instances, err := hdlr.scoreEngineUseCase.ListScoreEngines(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	response := make([]runningScoreEngine, 0, len(instances))

	for _, instance := range instances {
		response = append(response, runningScoreEngine{
			ContestID:  instance.ContestID,
			InstanceID: instance.InstanceID,
		})
	}

	writeResponse(w, http.StatusOK, response)
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
	err = json.NewDecoder(r.Body).Decode(&arguments)
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
