package rest

import (
	"net/http"

	"github.com/climblive/platform/backend/internal/domain"
)

type contestHandler struct {
	contestUseCase domain.ContestUseCase
}

func InstallContestHandler(contestUseCase domain.ContestUseCase) {
	handler := &contestHandler{
		contestUseCase: contestUseCase,
	}

	http.HandleFunc("GET /contests/{contestID}", handler.GetContest)
}

func (hdlr *contestHandler) GetContest(w http.ResponseWriter, r *http.Request) {
	contestID := parseResourceID(r.PathValue("contestID"))

	contest, err := hdlr.contestUseCase.GetContest(r.Context(), contestID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, contest)
}