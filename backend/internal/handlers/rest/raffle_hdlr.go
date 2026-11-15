package rest

import (
	"context"
	"net/http"

	"github.com/climblive/platform/backend/internal/domain"
)

type raffleUseCase interface {
	GetRaffle(ctx context.Context, raffleID domain.RaffleID) (domain.Raffle, error)
	GetRafflesByCompetition(ctx context.Context, competitionID domain.CompetitionID) ([]domain.Raffle, error)
	CreateRaffle(ctx context.Context, competitionID domain.CompetitionID) (domain.Raffle, error)
	DeleteRaffle(ctx context.Context, raffleID domain.RaffleID) error
	DrawRaffleWinner(ctx context.Context, raffleID domain.RaffleID) (domain.RaffleWinner, error)
	GetRaffleWinners(ctx context.Context, raffleID domain.RaffleID) ([]domain.RaffleWinner, error)
}

type raffleHandler struct {
	raffleUseCase raffleUseCase
}

func InstallRaffleHandler(mux *Mux, raffleUseCase raffleUseCase) {
	handler := &raffleHandler{
		raffleUseCase: raffleUseCase,
	}

	mux.HandleFunc("GET /raffles/{raffleID}", handler.GetRaffle)
	mux.HandleFunc("GET /competitions/{competitionID}/raffles", handler.GetRaffles)
	mux.HandleFunc("POST /competitions/{competitionID}/raffles", handler.CreateRaffle)
	mux.HandleFunc("DELETE /raffles/{raffleID}", handler.DeleteRaffle)
	mux.HandleFunc("POST /raffles/{raffleID}/winners", handler.DrawRaffleWinner)
	mux.HandleFunc("GET /raffles/{raffleID}/winners", handler.GetRaffleWinners)
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
	competitionID, err := parseCompetitionID(r.PathValue("competitionID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	raffles, err := hdlr.raffleUseCase.GetRafflesByCompetition(r.Context(), competitionID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, raffles)
}

func (hdlr *raffleHandler) CreateRaffle(w http.ResponseWriter, r *http.Request) {
	competitionID, err := parseCompetitionID(r.PathValue("competitionID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	createdRaffle, err := hdlr.raffleUseCase.CreateRaffle(r.Context(), competitionID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusCreated, createdRaffle)
}

func (hdlr *raffleHandler) DeleteRaffle(w http.ResponseWriter, r *http.Request) {
	raffleID, err := parseResourceID[domain.RaffleID](r.PathValue("raffleID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = hdlr.raffleUseCase.DeleteRaffle(r.Context(), raffleID)
	if err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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

func (hdlr *raffleHandler) GetRaffleWinners(w http.ResponseWriter, r *http.Request) {
	raffleID, err := parseResourceID[domain.RaffleID](r.PathValue("raffleID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	winners, err := hdlr.raffleUseCase.GetRaffleWinners(r.Context(), raffleID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, winners)
}
