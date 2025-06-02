package rest

import (
	"context"
	"net/http"

	"github.com/climblive/platform/backend/internal/domain"
)

type raffleUseCase interface {
	GetRaffle(ctx context.Context, raffleID domain.RaffleID) (domain.Raffle, error)
	GetRafflesByContest(ctx context.Context, contestID domain.ContestID) ([]domain.Raffle, error)
	CreateRaffle(ctx context.Context, contestID domain.ContestID) (domain.Raffle, error)
	DrawRaffleWinner(ctx context.Context, raffleID domain.RaffleID) (domain.RaffleWinner, error)
}

type raffleHandler struct {
	raffleUseCase raffleUseCase
}

func InstallRaffleHandler(mux *Mux, raffleUseCase raffleUseCase) {
	handler := &raffleHandler{
		raffleUseCase: raffleUseCase,
	}

	mux.HandleFunc("GET /raffles/{raffleID}", handler.GetRaffle)
	mux.HandleFunc("GET /contests/{contestID}/raffles", handler.GetRaffles)
	mux.HandleFunc("POST /contests/{contestID}/raffles", handler.CreateRaffle)
	mux.HandleFunc("POST /raffles/{raffleID}/winners", handler.DrawRaffleWinner)
}

func (hdlr *raffleHandler) GetRaffle(w http.ResponseWriter, r *http.Request) {
	raffleID, err := parseResourceID[domain.RaffleID](r.PathValue("raffleID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	raffle, err := hdlr.raffleUseCase.GetRaffle(r.Context(), raffleID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, raffle)
}

func (hdlr *raffleHandler) GetRaffles(w http.ResponseWriter, r *http.Request) {
	contestID, err := parseResourceID[domain.ContestID](r.PathValue("contestID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	raffles, err := hdlr.raffleUseCase.GetRafflesByContest(r.Context(), contestID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, raffles)
}

func (hdlr *raffleHandler) CreateRaffle(w http.ResponseWriter, r *http.Request) {
	contestID, err := parseResourceID[domain.ContestID](r.PathValue("contestID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	createdRaffle, err := hdlr.raffleUseCase.CreateRaffle(r.Context(), contestID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusCreated, createdRaffle)
}

func (hdlr *raffleHandler) DrawRaffleWinner(w http.ResponseWriter, r *http.Request) {
	raffleID, err := parseResourceID[domain.RaffleID](r.PathValue("raffleID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	winner, err := hdlr.raffleUseCase.DrawRaffleWinner(r.Context(), raffleID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusCreated, winner)
}
