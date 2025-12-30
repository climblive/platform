package rest

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/google/uuid"
)

type organizerUseCase interface {
	CreateOrganizer(ctx context.Context, template domain.OrganizerTemplate) (domain.Organizer, error)
	GetOrganizer(ctx context.Context, organizerID domain.OrganizerID) (domain.Organizer, error)
	PatchOrganizer(ctx context.Context, organizerID domain.OrganizerID, patch domain.OrganizerPatch) (domain.Organizer, error)
	GetOrganizerInvitesByOrganizer(ctx context.Context, organizerID domain.OrganizerID) ([]domain.OrganizerInvite, error)
	GetOrganizerInvite(ctx context.Context, inviteID domain.OrganizerInviteID) (domain.OrganizerInvite, error)
	CreateOrganizerInvite(ctx context.Context, organizerID domain.OrganizerID) (domain.OrganizerInvite, error)
	DeleteOrganizerInvite(ctx context.Context, inviteID domain.OrganizerInviteID) error
	AcceptOrganizerInvite(ctx context.Context, inviteID domain.OrganizerInviteID) error
}

type organizerHandler struct {
	organizerUseCase organizerUseCase
}

func InstallOrganizerHandler(mux *Mux, organizerUseCase organizerUseCase) {
	handler := &organizerHandler{
		organizerUseCase: organizerUseCase,
	}

	mux.HandleFunc("POST /organizers", handler.CreateOrganizer)
	mux.HandleFunc("GET /organizers/{organizerID}", handler.GetOrganizer)
	mux.HandleFunc("PATCH /organizers/{organizerID}", handler.PatchOrganizer)
	mux.HandleFunc("GET /organizers/{organizerID}/invites", handler.GetOrganizerInvites)
	mux.HandleFunc("POST /organizers/{organizerID}/invites", handler.CreateOrganizerInvite)
	mux.HandleFunc("GET /invites/{inviteID}", handler.GetOrganizerInvite)
	mux.HandleFunc("DELETE /invites/{inviteID}", handler.DeleteOrganizerInvite)
	mux.HandleFunc("POST /invites/{inviteID}/accept", handler.AcceptOrganizerInvite)
}

func (hdlr *organizerHandler) CreateOrganizer(w http.ResponseWriter, r *http.Request) {
	var template domain.OrganizerTemplate
	err := json.NewDecoder(r.Body).Decode(&template)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	organizer, err := hdlr.organizerUseCase.CreateOrganizer(r.Context(), template)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusCreated, organizer)
}

func (hdlr *organizerHandler) GetOrganizer(w http.ResponseWriter, r *http.Request) {
	organizerID, err := parseResourceID[domain.OrganizerID](r.PathValue("organizerID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	organizer, err := hdlr.organizerUseCase.GetOrganizer(r.Context(), organizerID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, organizer)
}

func (hdlr *organizerHandler) PatchOrganizer(w http.ResponseWriter, r *http.Request) {
	organizerID, err := parseResourceID[domain.OrganizerID](r.PathValue("organizerID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var patch domain.OrganizerPatch
	err = json.NewDecoder(r.Body).Decode(&patch)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	updatedOrganizer, err := hdlr.organizerUseCase.PatchOrganizer(r.Context(), organizerID, patch)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, updatedOrganizer)
}

func (hdlr *organizerHandler) GetOrganizerInvites(w http.ResponseWriter, r *http.Request) {
	organizerID, err := parseResourceID[domain.OrganizerID](r.PathValue("organizerID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	invites, err := hdlr.organizerUseCase.GetOrganizerInvitesByOrganizer(r.Context(), organizerID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, invites)
}

func (hdlr *organizerHandler) CreateOrganizerInvite(w http.ResponseWriter, r *http.Request) {
	organizerID, err := parseResourceID[domain.OrganizerID](r.PathValue("organizerID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	invite, err := hdlr.organizerUseCase.CreateOrganizerInvite(r.Context(), organizerID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusCreated, invite)
}

func (hdlr *organizerHandler) GetOrganizerInvite(w http.ResponseWriter, r *http.Request) {
	inviteID, err := uuid.Parse(r.PathValue("inviteID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	invite, err := hdlr.organizerUseCase.GetOrganizerInvite(r.Context(), domain.OrganizerInviteID(inviteID))
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, invite)
}

func (hdlr *organizerHandler) DeleteOrganizerInvite(w http.ResponseWriter, r *http.Request) {
	inviteID, err := uuid.Parse(r.PathValue("inviteID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = hdlr.organizerUseCase.DeleteOrganizerInvite(r.Context(), domain.OrganizerInviteID(inviteID))
	if err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (hdlr *organizerHandler) AcceptOrganizerInvite(w http.ResponseWriter, r *http.Request) {
	inviteID, err := uuid.Parse(r.PathValue("inviteID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = hdlr.organizerUseCase.AcceptOrganizerInvite(r.Context(), domain.OrganizerInviteID(inviteID))
	if err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
