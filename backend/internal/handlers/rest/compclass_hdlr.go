package rest

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/climblive/platform/backend/internal/domain"
)

type compClassUseCase interface {
	GetCompClassesByContest(ctx context.Context, contestID domain.ContestID) ([]domain.CompClass, error)
	CreateCompClass(ctx context.Context, contestID domain.ContestID, tmpl domain.CompClassTemplate) (domain.CompClass, error)
	DeleteCompClass(ctx context.Context, compClassID domain.CompClassID) error
	PatchCompClass(ctx context.Context, compClassID domain.CompClassID, tmpl domain.CompClassPatch) (domain.CompClass, error)
}

type compClassHandler struct {
	compClassUseCase compClassUseCase
}

func InstallCompClassHandler(mux *Mux, compClassUseCase compClassUseCase) {
	handler := &compClassHandler{
		compClassUseCase: compClassUseCase,
	}

	mux.HandleFunc("GET /contests/{contestID}/comp-classes", handler.GetCompClassesByContest)
	mux.HandleFunc("POST /contests/{contestID}/comp-classes", handler.CreateCompClass)
	mux.HandleFunc("DELETE /comp-classes/{compClassID}", handler.DeleteCompClass)
	mux.HandleFunc("PATCH /comp-classes/{compClassID}", handler.PatchCompClass)
}

func (hdlr *compClassHandler) GetCompClassesByContest(w http.ResponseWriter, r *http.Request) {
	contestID, err := parseResourceID[domain.ContestID](r.PathValue("contestID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	compClasses, err := hdlr.compClassUseCase.GetCompClassesByContest(r.Context(), contestID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, compClasses)
}

func (hdlr *compClassHandler) CreateCompClass(w http.ResponseWriter, r *http.Request) {
	contestID, err := parseResourceID[domain.ContestID](r.PathValue("contestID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var tmpl domain.CompClassTemplate
	err = json.NewDecoder(r.Body).Decode(&tmpl)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	createdCompClass, err := hdlr.compClassUseCase.CreateCompClass(r.Context(), contestID, tmpl)
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
	err = json.NewDecoder(r.Body).Decode(&patch)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	updatedCompClass, err := hdlr.compClassUseCase.PatchCompClass(r.Context(), compClassID, patch)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, updatedCompClass)
}
