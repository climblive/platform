package rest

import (
	"context"
	"net/http"

	"github.com/climblive/platform/backend/internal/domain"
)

type problemUseCase interface {
	GetProblem(ctx context.Context, problemID domain.ProblemID) (domain.Problem, error)
	GetProblemsByCompetition(ctx context.Context, competitionID domain.CompetitionID) ([]domain.Problem, error)
	PatchProblem(ctx context.Context, problemID domain.ProblemID, patch domain.ProblemPatch) (domain.Problem, error)
	CreateProblem(ctx context.Context, competitionID domain.CompetitionID, tmpl domain.ProblemTemplate) (domain.Problem, error)
	DeleteProblem(ctx context.Context, problemID domain.ProblemID) error
}

type problemHandler struct {
	problemUseCase problemUseCase
}

func InstallProblemHandler(mux *Mux, problemUseCase problemUseCase) {
	handler := &problemHandler{
		problemUseCase: problemUseCase,
	}

	mux.HandleFunc("GET /problems/{problemID}", handler.GetProblem)
	mux.HandleFunc("GET /competitions/{competitionID}/problems", handler.GetProblemsByCompetition)
	mux.HandleFunc("PATCH /problems/{problemID}", handler.PatchProblem)
	mux.HandleFunc("POST /competitions/{competitionID}/problems", handler.CreateProblem)
	mux.HandleFunc("DELETE /problems/{problemID}", handler.DeleteProblem)
}

func (hdlr *problemHandler) GetProblem(w http.ResponseWriter, r *http.Request) {
	problemID, err := parseResourceID[domain.ProblemID](r.PathValue("problemID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	problem, err := hdlr.problemUseCase.GetProblem(r.Context(), problemID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, problem)
}

func (hdlr *problemHandler) GetProblemsByCompetition(w http.ResponseWriter, r *http.Request) {
	competitionID, err := parseCompetitionID(r.PathValue("competitionID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	problems, err := hdlr.problemUseCase.GetProblemsByCompetition(r.Context(), competitionID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, problems)
}

func (hdlr *problemHandler) PatchProblem(w http.ResponseWriter, r *http.Request) {
	problemID, err := parseResourceID[domain.ProblemID](r.PathValue("problemID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var patch domain.ProblemPatch
	if !readJSON(w, r, &patch) {
		return
	}

	updatedProblem, err := hdlr.problemUseCase.PatchProblem(r.Context(), problemID, patch)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, updatedProblem)
}

func (hdlr *problemHandler) CreateProblem(w http.ResponseWriter, r *http.Request) {
	competitionID, err := parseCompetitionID(r.PathValue("competitionID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var tmpl domain.ProblemTemplate
	if !readJSON(w, r, &tmpl) {
		return
	}

	createdProblem, err := hdlr.problemUseCase.CreateProblem(r.Context(), competitionID, tmpl)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusCreated, createdProblem)
}

func (hdlr *problemHandler) DeleteProblem(w http.ResponseWriter, r *http.Request) {
	problemID, err := parseResourceID[domain.ProblemID](r.PathValue("problemID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = hdlr.problemUseCase.DeleteProblem(r.Context(), problemID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusNoContent, nil)
}
