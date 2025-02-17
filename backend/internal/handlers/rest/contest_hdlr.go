package rest

import (
	"context"
	"net/http"

	"github.com/climblive/platform/backend/internal/domain"
)

type contestUseCase interface {
	GetContest(ctx context.Context, contestID domain.ContestID) (domain.Contest, error)
	GetContestsByOrganizer(ctx context.Context, organizerID domain.OrganizerID) ([]domain.Contest, error)
	GetScoreboard(ctx context.Context, contestID domain.ContestID) ([]domain.ScoreboardEntry, error)
}

type contestHandler struct {
	contestUseCase contestUseCase
}

func InstallContestHandler(mux *Mux, contestUseCase contestUseCase) {
	handler := &contestHandler{
		contestUseCase: contestUseCase,
	}

	mux.HandleFunc("GET /contests/{contestID}", handler.GetContest)
	mux.HandleFunc("GET /contests/{contestID}/scoreboard", handler.GetScoreboard)
	mux.HandleFunc("GET /organizers/{organizerID}/contests", handler.GetContestsByOrganizer)
}

func (hdlr *contestHandler) GetContest(w http.ResponseWriter, r *http.Request) {
	contestID, err := parseResourceID[domain.ContestID](r.PathValue("contestID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	contest, err := hdlr.contestUseCase.GetContest(r.Context(), contestID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, contest)
}

func (hdlr *contestHandler) GetScoreboard(w http.ResponseWriter, r *http.Request) {
	contestID, err := parseResourceID[domain.ContestID](r.PathValue("contestID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	scoreboard, err := hdlr.contestUseCase.GetScoreboard(r.Context(), contestID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, scoreboard)
}

func (hdlr *contestHandler) GetContestsByOrganizer(w http.ResponseWriter, r *http.Request) {
	organizerID, err := parseResourceID[domain.OrganizerID](r.PathValue("organizerID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	contest, err := hdlr.contestUseCase.GetContestsByOrganizer(r.Context(), organizerID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, contest)
}
