package rest

import (
	"context"
	"net/http"

	"github.com/climblive/platform/backend/internal/domain"
)

type problemUseCase interface {
	GetProblemsByContest(ctx context.Context, contestID domain.ContestID) ([]domain.Problem, error)
}

type problemHandler struct {
	problemUseCase problemUseCase
}

func InstallProblemHandler(mux *Mux, problemUseCase problemUseCase) {
	handler := &problemHandler{
		problemUseCase: problemUseCase,
	}

	mux.HandleFunc("GET /contests/{contestID}/problems", handler.GetProblemsByContest)
}

func (hdlr *problemHandler) GetProblemsByContest(w http.ResponseWriter, r *http.Request) {
	contestID, err := parseResourceID[domain.ContestID](r.PathValue("contestID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	problems, err := hdlr.problemUseCase.GetProblemsByContest(r.Context(), contestID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, problems)
}
