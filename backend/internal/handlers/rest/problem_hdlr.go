package rest

import (
	"context"
	"encoding/json/v2"
	"net/http"

	"github.com/climblive/platform/backend/internal/domain"
)

type problemUseCase interface {
	GetProblem(ctx context.Context, problemID domain.ProblemID) (domain.Problem, error)
	GetProblemsByContest(ctx context.Context, contestID domain.ContestID) ([]domain.Problem, error)
	PatchProblem(ctx context.Context, problemID domain.ProblemID, patch domain.ProblemPatch) (domain.Problem, error)
	CreateProblem(ctx context.Context, contestID domain.ContestID, tmpl domain.ProblemTemplate) (domain.Problem, error)
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
	mux.HandleFunc("GET /contests/{contestID}/problems", handler.GetProblemsByContest)
	mux.HandleFunc("PATCH /problems/{problemID}", handler.PatchProblem)
	mux.HandleFunc("POST /contests/{contestID}/problems", handler.CreateProblem)
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

func (hdlr *problemHandler) GetProblemsByContest(w http.ResponseWriter, r *http.Request) {
	contestID, err := parseResourceID[domain.ContestID](r.PathValue("contestID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	problems, err := hdlr.problemUseCase.GetProblemsByContest(r.Context(), contestID)
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
	err = json.UnmarshalRead(r.Body, &patch)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
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
	contestID, err := parseResourceID[domain.ContestID](r.PathValue("contestID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var tmpl domain.ProblemTemplate
	err = json.UnmarshalRead(r.Body, &tmpl)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	createdProblem, err := hdlr.problemUseCase.CreateProblem(r.Context(), contestID, tmpl)
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
