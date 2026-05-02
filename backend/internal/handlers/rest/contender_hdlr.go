package rest

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/climblive/platform/backend/internal/domain"
)

type contenderUseCase interface {
	GetContender(ctx context.Context, contenderID domain.ContenderID) (domain.Contender, error)
	GetContenderByCode(ctx context.Context, registrationCode string) (domain.Contender, error)
	GetContendersByCompClass(ctx context.Context, compClassID domain.CompClassID) ([]domain.Contender, error)
	GetContendersByContest(ctx context.Context, contestID domain.ContestID) ([]domain.Contender, error)
	PatchContender(ctx context.Context, contenderID domain.ContenderID, patch domain.ContenderPatch) (domain.Contender, error)
	ScrubContender(ctx context.Context, contenderID domain.ContenderID) (domain.Contender, error)
	DeleteContender(ctx context.Context, contenderID domain.ContenderID) error
	CreateContenders(ctx context.Context, contestID domain.ContestID, number int) ([]domain.Contender, error)
}

type contenderHandler struct {
	contenderUseCase contenderUseCase
}

func InstallContenderHandler(mux *Mux, contenderUseCase contenderUseCase) {
	handler := &contenderHandler{
		contenderUseCase: contenderUseCase,
	}

	mux.HandleFunc("GET /contenders/{contenderID}", handler.GetContender)
	mux.HandleFunc("GET /codes/{registrationCode}/contender", handler.GetContenderByCode)
	mux.HandleFunc("GET /compClasses/{compClassID}/contenders", handler.GetContendersByCompClass)
	mux.HandleFunc("GET /contests/{contestID}/contenders", handler.GetContendersByContest)
	mux.HandleFunc("PATCH /contenders/{contenderID}", handler.PatchContender)
	mux.HandleFunc("POST /contenders/{contenderID}/scrub", handler.ScrubContender)
	mux.HandleFunc("DELETE /contenders/{contenderID}", handler.DeleteContender)
	mux.HandleFunc("POST /contests/{contestID}/contenders", handler.CreateContenders)
}

func (hdlr *contenderHandler) GetContender(w http.ResponseWriter, r *http.Request) {
	contenderID, err := parseResourceID[domain.ContenderID](r.PathValue("contenderID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	contender, err := hdlr.contenderUseCase.GetContender(r.Context(), contenderID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, contender)
}

func (hdlr *contenderHandler) GetContenderByCode(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("registrationCode")

	contender, err := hdlr.contenderUseCase.GetContenderByCode(r.Context(), code)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, contender)
}

func (hdlr *contenderHandler) GetContendersByCompClass(w http.ResponseWriter, r *http.Request) {
	compClassID, err := parseResourceID[domain.CompClassID](r.PathValue("compClassID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	contenders, err := hdlr.contenderUseCase.GetContendersByCompClass(r.Context(), compClassID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, contenders)
}

func (hdlr *contenderHandler) GetContendersByContest(w http.ResponseWriter, r *http.Request) {
	contestID, err := parseResourceID[domain.ContestID](r.PathValue("contestID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	contenders, err := hdlr.contenderUseCase.GetContendersByContest(r.Context(), contestID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, contenders)
}

func (hdlr *contenderHandler) PatchContender(w http.ResponseWriter, r *http.Request) {
	contenderID, err := parseResourceID[domain.ContenderID](r.PathValue("contenderID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var patch domain.ContenderPatch
	err = json.NewDecoder(r.Body).Decode(&patch)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	updatedContender, err := hdlr.contenderUseCase.PatchContender(r.Context(), contenderID, patch)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, updatedContender)
}

func (hdlr *contenderHandler) DeleteContender(w http.ResponseWriter, r *http.Request) {
	contenderID, err := parseResourceID[domain.ContenderID](r.PathValue("contenderID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = hdlr.contenderUseCase.DeleteContender(r.Context(), contenderID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusNoContent, nil)
}

func (hdlr *contenderHandler) ScrubContender(w http.ResponseWriter, r *http.Request) {
	contenderID, err := parseResourceID[domain.ContenderID](r.PathValue("contenderID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	contender, err := hdlr.contenderUseCase.ScrubContender(r.Context(), contenderID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, contender)
}

func (hdlr *contenderHandler) CreateContenders(w http.ResponseWriter, r *http.Request) {
	contestID, err := parseResourceID[domain.ContestID](r.PathValue("contestID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var arguments CreateContendersArguments
	err = json.NewDecoder(r.Body).Decode(&arguments)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	contenders, err := hdlr.contenderUseCase.CreateContenders(r.Context(), contestID, arguments.Number)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusCreated, contenders)
}
