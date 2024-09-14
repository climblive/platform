package rest

import (
	"net/http"

	"github.com/climblive/platform/backend/internal/domain"
)

type tickHandler struct {
	tickUseCase domain.TickUseCase
}

func InstallTickHandler(tickUseCase domain.TickUseCase) {
	handler := &tickHandler{
		tickUseCase: tickUseCase,
	}

	http.HandleFunc("GET /contenders/{contenderID}/ticks", handler.GetTicksByContender)
}

func (hdlr *tickHandler) GetTicksByContender(w http.ResponseWriter, r *http.Request) {
	contenderID := parseResourceID(r.PathValue("contenderID"))

	ticks, err := hdlr.tickUseCase.GetTicksByContender(r.Context(), contenderID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, ticks)
}
