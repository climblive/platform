package rest

import (
	"context"
	"net/http"

	"github.com/climblive/platform/backend/internal/domain"
)

type categoryUseCase interface {
	GetCategoriesByCompetition(ctx context.Context, competitionID domain.CompetitionID) ([]domain.Category, error)
	CreateCategory(ctx context.Context, competitionID domain.CompetitionID, tmpl domain.CategoryTemplate) (domain.Category, error)
	DeleteCategory(ctx context.Context, categoryID domain.CategoryID) error
	GetCategory(ctx context.Context, categoryID domain.CategoryID) (domain.Category, error)
	PatchCategory(ctx context.Context, categoryID domain.CategoryID, patch domain.CategoryPatch) (domain.Category, error)
}

type categoryHandler struct {
	categoryUseCase categoryUseCase
}

func InstallCategoryHandler(mux *Mux, categoryUseCase categoryUseCase) {
	handler := &categoryHandler{
		categoryUseCase: categoryUseCase,
	}

	mux.HandleFunc("GET /categories/{categoryID}", handler.GetCategory)
	mux.HandleFunc("GET /competitions/{competitionID}/categories", handler.GetCategoriesByCompetition)
	mux.HandleFunc("POST /competitions/{competitionID}/categories", handler.CreateCategory)
	mux.HandleFunc("DELETE /categories/{categoryID}", handler.DeleteCategory)
	mux.HandleFunc("PATCH /categories/{categoryID}", handler.PatchCategory)
}

func (hdlr *categoryHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	categoryID, err := parseResourceID[domain.CategoryID](r.PathValue("categoryID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	category, err := hdlr.categoryUseCase.GetCategory(r.Context(), categoryID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, category)
}

func (hdlr *categoryHandler) GetCategoriesByCompetition(w http.ResponseWriter, r *http.Request) {
	competitionID, err := parseResourceID[domain.CompetitionID](r.PathValue("competitionID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	categories, err := hdlr.categoryUseCase.GetCategoriesByCompetition(r.Context(), competitionID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, categories)
}

func (hdlr *categoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	competitionID, err := parseResourceID[domain.CompetitionID](r.PathValue("competitionID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var tmpl domain.CategoryTemplate
	if !readJSON(w, r, &tmpl) {
		return
	}

	createdCategory, err := hdlr.categoryUseCase.CreateCategory(r.Context(), competitionID, tmpl)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusCreated, createdCategory)
}

func (hdlr *categoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	categoryID, err := parseResourceID[domain.CategoryID](r.PathValue("categoryID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = hdlr.categoryUseCase.DeleteCategory(r.Context(), categoryID)
	if err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (hdlr *categoryHandler) PatchCategory(w http.ResponseWriter, r *http.Request) {
	categoryID, err := parseResourceID[domain.CategoryID](r.PathValue("categoryID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var patch domain.CategoryPatch
	if !readJSON(w, r, &patch) {
		return
	}

	updatedCategory, err := hdlr.categoryUseCase.PatchCategory(r.Context(), categoryID, patch)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, updatedCategory)
}
