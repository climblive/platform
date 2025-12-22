package rest

import (
	"context"
	"net/http"

	"github.com/climblive/platform/backend/internal/domain"
)

type userUseCase interface {
	GetSelf(ctx context.Context) (domain.User, error)
	GetUsersByOrganizer(ctx context.Context, organizerID domain.OrganizerID) ([]domain.User, error)
}

type userHandler struct {
	userUseCase userUseCase
}

func InstallUserHandler(mux *Mux, userUseCase userUseCase) {
	handler := &userHandler{
		userUseCase: userUseCase,
	}

	mux.HandleFunc("GET /users/self", handler.GetSelf)
	mux.HandleFunc("GET /organizers/{organizerID}/users", handler.GetUsersByOrganizer)
}

func (hdlr *userHandler) GetSelf(w http.ResponseWriter, r *http.Request) {
	user, err := hdlr.userUseCase.GetSelf(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, user)
}

func (hdlr *userHandler) GetUsersByOrganizer(w http.ResponseWriter, r *http.Request) {
	organizerID, err := parseResourceID[domain.OrganizerID](r.PathValue("organizerID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	users, err := hdlr.userUseCase.GetUsersByOrganizer(r.Context(), organizerID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, users)
}
