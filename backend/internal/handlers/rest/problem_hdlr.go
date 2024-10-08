package rest

import (
	"net/http"

	"github.com/climblive/platform/backend/internal/domain"
)

type problemHandler struct {
	problemUseCase domain.ProblemUseCase
}

func InstallProblemHandler(mux *Mux, problemUseCase domain.ProblemUseCase) {
	handler := &problemHandler{
		problemUseCase: problemUseCase,
	}

	mux.HandleFunc("GET /contests/{contestID}/problems", handler.GetProblemsByContest)
}

func (hdlr *problemHandler) GetProblemsByContest(w http.ResponseWriter, r *http.Request) {
	contestID := parseResourceID(r.PathValue("contestID"))

	problems, err := hdlr.problemUseCase.GetProblemsByContest(r.Context(), contestID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, problems)
}
