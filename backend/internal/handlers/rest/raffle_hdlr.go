package rest

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/climblive/platform/backend/internal/domain"
)

type raffleUseCase interface {
	CreateRaffle(ctx context.Context, contestID domain.ContestID, tmpl domain.RaffleTemplate) (domain.Raffle, error)
}

type raffleHandler struct {
	raffleUseCase raffleUseCase
}

func InstallRaffleHandler(mux *Mux, raffleUseCase raffleUseCase) {
	handler := &raffleHandler{
		raffleUseCase: raffleUseCase,
	}

	mux.HandleFunc("POST /contests/{contestID}/raffles", handler.CreateRaffle)
}

func (hdlr *raffleHandler) CreateRaffle(w http.ResponseWriter, r *http.Request) {
	contestID, err := parseResourceID[domain.ContestID](r.PathValue("contestID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var tmpl domain.RaffleTemplate
	err = json.NewDecoder(r.Body).Decode(&tmpl)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	createdRaffle, err := hdlr.raffleUseCase.CreateRaffle(r.Context(), contestID, tmpl)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusCreated, createdRaffle)
}
