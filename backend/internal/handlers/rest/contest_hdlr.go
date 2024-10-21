package rest

import (
	"net/http"

	"github.com/climblive/platform/backend/internal/domain"
)

type contestHandler struct {
	contestUseCase domain.ContestUseCase
}

func InstallContestHandler(mux *Mux, contestUseCase domain.ContestUseCase) {
	handler := &contestHandler{
		contestUseCase: contestUseCase,
	}

	mux.HandleFunc("GET /contests/{contestID}", handler.GetContest)
	mux.HandleFunc("GET /contests/{contestID}/scoreboard", handler.GetScoreboard)
}

func (hdlr *contestHandler) GetContest(w http.ResponseWriter, r *http.Request) {
	contestID := parseResourceID[domain.ContestID](r.PathValue("contestID"))

	contest, err := hdlr.contestUseCase.GetContest(r.Context(), contestID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, contest)
}

func (hdlr *contestHandler) GetScoreboard(w http.ResponseWriter, r *http.Request) {
	contestID := parseResourceID[domain.ContestID](r.PathValue("contestID"))

	scoreboard, err := hdlr.contestUseCase.GetScoreboard(r.Context(), contestID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, scoreboard)
}
