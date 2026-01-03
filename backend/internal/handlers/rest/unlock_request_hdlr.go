package rest

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/climblive/platform/backend/internal/domain"
)

type unlockRequestUseCase interface {
	CreateUnlockRequest(ctx context.Context, template domain.UnlockRequestTemplate) (domain.UnlockRequestID, error)
	GetUnlockRequest(ctx context.Context, id domain.UnlockRequestID) (domain.UnlockRequest, error)
	GetUnlockRequestsByContest(ctx context.Context, contestID domain.ContestID) ([]domain.UnlockRequest, error)
	GetUnlockRequestsByOrganizer(ctx context.Context, organizerID domain.OrganizerID) ([]domain.UnlockRequest, error)
	GetPendingUnlockRequests(ctx context.Context) ([]domain.UnlockRequest, error)
	ReviewUnlockRequest(ctx context.Context, id domain.UnlockRequestID, review domain.UnlockRequestReview) error
}

type unlockRequestHandler struct {
	unlockRequestUseCase unlockRequestUseCase
}

func InstallUnlockRequestHandler(mux *Mux, unlockRequestUseCase unlockRequestUseCase) {
	handler := &unlockRequestHandler{
		unlockRequestUseCase: unlockRequestUseCase,
	}

	mux.HandleFunc("POST /contests/{contestID}/unlock-requests", handler.CreateUnlockRequest)
	mux.HandleFunc("GET /unlock-requests/{id}", handler.GetUnlockRequest)
	mux.HandleFunc("GET /contests/{contestID}/unlock-requests", handler.GetUnlockRequestsByContest)
	mux.HandleFunc("GET /organizers/{organizerID}/unlock-requests", handler.GetUnlockRequestsByOrganizer)
	mux.HandleFunc("GET /unlock-requests", handler.GetPendingUnlockRequests)
	mux.HandleFunc("PATCH /unlock-requests/{id}/review", handler.ReviewUnlockRequest)
}

func (hdlr *unlockRequestHandler) CreateUnlockRequest(w http.ResponseWriter, r *http.Request) {
	contestID, err := parseResourceID[domain.ContestID](r.PathValue("contestID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var template domain.UnlockRequestTemplate
	if err := json.NewDecoder(r.Body).Decode(&template); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	template.ContestID = contestID

	id, err := hdlr.unlockRequestUseCase.CreateUnlockRequest(r.Context(), template)
	if err != nil {
		handleError(w, err)
		return
	}

	request, err := hdlr.unlockRequestUseCase.GetUnlockRequest(r.Context(), id)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusCreated, request)
}

func (hdlr *unlockRequestHandler) GetUnlockRequest(w http.ResponseWriter, r *http.Request) {
	id, err := parseResourceID[domain.UnlockRequestID](r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	request, err := hdlr.unlockRequestUseCase.GetUnlockRequest(r.Context(), id)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, request)
}

func (hdlr *unlockRequestHandler) GetUnlockRequestsByContest(w http.ResponseWriter, r *http.Request) {
	contestID, err := parseResourceID[domain.ContestID](r.PathValue("contestID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	requests, err := hdlr.unlockRequestUseCase.GetUnlockRequestsByContest(r.Context(), contestID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, requests)
}

func (hdlr *unlockRequestHandler) GetUnlockRequestsByOrganizer(w http.ResponseWriter, r *http.Request) {
	organizerID, err := parseResourceID[domain.OrganizerID](r.PathValue("organizerID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	requests, err := hdlr.unlockRequestUseCase.GetUnlockRequestsByOrganizer(r.Context(), organizerID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, requests)
}

func (hdlr *unlockRequestHandler) GetPendingUnlockRequests(w http.ResponseWriter, r *http.Request) {
	requests, err := hdlr.unlockRequestUseCase.GetPendingUnlockRequests(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, requests)
}

func (hdlr *unlockRequestHandler) ReviewUnlockRequest(w http.ResponseWriter, r *http.Request) {
	id, err := parseResourceID[domain.UnlockRequestID](r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var review domain.UnlockRequestReview
	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := hdlr.unlockRequestUseCase.ReviewUnlockRequest(r.Context(), id, review); err != nil {
		handleError(w, err)
		return
	}

	request, err := hdlr.unlockRequestUseCase.GetUnlockRequest(r.Context(), id)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, request)
}
