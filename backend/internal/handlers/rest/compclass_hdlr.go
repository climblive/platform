package rest

import (
	"net/http"

	"github.com/climblive/platform/backend/internal/domain"
)

type compClassHandler struct {
	compClassUseCase domain.CompClassUseCase
}

func InstallCompClassHandler(mux *Mux, compClassUseCase domain.CompClassUseCase) {
	handler := &compClassHandler{
		compClassUseCase: compClassUseCase,
	}

	mux.HandleFunc("GET /contests/{contestID}/compClasses", handler.GetCompClassesByContest)
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
