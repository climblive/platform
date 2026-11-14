package rest

import (
	"context"
	"net/http"

	"github.com/climblive/platform/backend/internal/domain"
)

type compClassUseCase interface {
	GetCompClassesByCompetition(ctx context.Context, competitionID domain.CompetitionID) ([]domain.CompClass, error)
	CreateCompClass(ctx context.Context, competitionID domain.CompetitionID, tmpl domain.CompClassTemplate) (domain.CompClass, error)
	DeleteCompClass(ctx context.Context, compClassID domain.CompClassID) error
	GetCompClass(ctx context.Context, compClassID domain.CompClassID) (domain.CompClass, error)
	PatchCompClass(ctx context.Context, compClassID domain.CompClassID, patch domain.CompClassPatch) (domain.CompClass, error)
}

type compClassHandler struct {
	compClassUseCase compClassUseCase
}

func InstallCompClassHandler(mux *Mux, compClassUseCase compClassUseCase) {
	handler := &compClassHandler{
		compClassUseCase: compClassUseCase,
	}

	mux.HandleFunc("GET /comp-classes/{compClassID}", handler.GetCompClass)
	mux.HandleFunc("GET /competitions/{competitionID}/comp-classes", handler.GetCompClassesByCompetition)
	mux.HandleFunc("POST /competitions/{competitionID}/comp-classes", handler.CreateCompClass)
	mux.HandleFunc("DELETE /comp-classes/{compClassID}", handler.DeleteCompClass)
	mux.HandleFunc("PATCH /comp-classes/{compClassID}", handler.PatchCompClass)
}

func (hdlr *compClassHandler) GetCompClass(w http.ResponseWriter, r *http.Request) {
	compClassID, err := parseResourceID[domain.CompClassID](r.PathValue("compClassID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	compClass, err := hdlr.compClassUseCase.GetCompClass(r.Context(), compClassID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, compClass)
}

func (hdlr *compClassHandler) GetCompClassesByCompetition(w http.ResponseWriter, r *http.Request) {
	competitionID, err := parseResourceID[domain.CompetitionID](r.PathValue("competitionID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	compClasses, err := hdlr.compClassUseCase.GetCompClassesByCompetition(r.Context(), competitionID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, compClasses)
}

func (hdlr *compClassHandler) CreateCompClass(w http.ResponseWriter, r *http.Request) {
	competitionID, err := parseResourceID[domain.CompetitionID](r.PathValue("competitionID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var tmpl domain.CompClassTemplate
	if !readJSON(w, r, &tmpl) {
		return
	}

	createdCompClass, err := hdlr.compClassUseCase.CreateCompClass(r.Context(), competitionID, tmpl)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusCreated, createdCompClass)
}

func (hdlr *compClassHandler) DeleteCompClass(w http.ResponseWriter, r *http.Request) {
	compClassID, err := parseResourceID[domain.CompClassID](r.PathValue("compClassID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = hdlr.compClassUseCase.DeleteCompClass(r.Context(), compClassID)
	if err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (hdlr *compClassHandler) PatchCompClass(w http.ResponseWriter, r *http.Request) {
	compClassID, err := parseResourceID[domain.CompClassID](r.PathValue("compClassID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var patch domain.CompClassPatch
	if !readJSON(w, r, &patch) {
		return
	}

	updatedCompClass, err := hdlr.compClassUseCase.PatchCompClass(r.Context(), compClassID, patch)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, updatedCompClass)
}
