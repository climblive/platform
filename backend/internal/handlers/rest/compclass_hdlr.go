package rest

import (
	"net/http"

	"github.com/climblive/platform/backend/internal/domain"
)

type compClassHandler struct {
	compClassUseCase domain.CompClassUseCase
}

func InstallCompClassHandler(compClassUseCase domain.CompClassUseCase) {
	handler := &compClassHandler{
		compClassUseCase: compClassUseCase,
	}

	http.HandleFunc("GET /contests/{contestID}/compClasses", handler.GetCompClassesByContest)
}

func (hdlr *compClassHandler) GetCompClassesByContest(w http.ResponseWriter, r *http.Request) {
	contestID := parseResourceID(r.PathValue("contestID"))

	compClasses, err := hdlr.compClassUseCase.GetCompClassesByContest(r.Context(), contestID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, compClasses)
}
