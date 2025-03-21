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
