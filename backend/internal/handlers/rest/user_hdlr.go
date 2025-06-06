package rest

import (
	"context"
	"net/http"

	"github.com/climblive/platform/backend/internal/domain"
)

type userUseCase interface {
	GetSelf(ctx context.Context) (domain.User, error)
}

type userHandler struct {
	userUseClass userUseCase
}

func InstallUserHandler(mux *Mux, userUseCase userUseCase) {
	handler := &userHandler{
		userUseClass: userUseCase,
	}

	mux.HandleFunc("GET /users/self", handler.GetSelf)
}

func (hdlr *userHandler) GetSelf(w http.ResponseWriter, r *http.Request) {
	user, err := hdlr.userUseClass.GetSelf(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, user)
}
