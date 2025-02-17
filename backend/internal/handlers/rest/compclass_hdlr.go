package rest

import (
	"context"
	"net/http"

	"github.com/climblive/platform/backend/internal/domain"
)

type compClassUseCase interface {
	GetCompClassesByContest(ctx context.Context, contestID domain.ContestID) ([]domain.CompClass, error)
}

type compClassHandler struct {
	compClassUseCase compClassUseCase
}

func InstallCompClassHandler(mux *Mux, compClassUseCase compClassUseCase) {
	handler := &compClassHandler{
		compClassUseCase: compClassUseCase,
	}

	mux.HandleFunc("GET /contests/{contestID}/comp-classes", handler.GetCompClassesByContest)
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
