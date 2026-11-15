package rest

import (
	"context"
	"net/http"

	"github.com/climblive/platform/backend/internal/domain"
)

type contenderUseCase interface {
	GetContender(ctx context.Context, contenderID domain.ContenderID) (domain.Contender, error)
	GetContenderByCode(ctx context.Context, registrationCode string) (domain.Contender, error)
	GetPointValues(ctx context.Context, contenderID domain.ContenderID) ([]domain.PointValue, error)
	GetContendersByCompClass(ctx context.Context, compClassID domain.CompClassID) ([]domain.Contender, error)
	GetContendersByCompetition(ctx context.Context, competitionID domain.CompetitionID) ([]domain.Contender, error)
	PatchContender(ctx context.Context, contenderID domain.ContenderID, patch domain.ContenderPatch) (domain.Contender, error)
	ScrubContender(ctx context.Context, contenderID domain.ContenderID) (domain.Contender, error)
	DeleteContender(ctx context.Context, contenderID domain.ContenderID) error
	CreateContenders(ctx context.Context, competitionID domain.CompetitionID, number int) ([]domain.Contender, error)
}

type contenderHandler struct {
	contenderUseCase contenderUseCase
}

func InstallContenderHandler(mux *Mux, contenderUseCase contenderUseCase) {
	handler := &contenderHandler{
		contenderUseCase: contenderUseCase,
	}

	mux.HandleFunc("GET /contenders/{contenderID}", handler.GetContender)
	mux.HandleFunc("GET /contenders/{contenderID}/points", handler.GetPointValues)
	mux.HandleFunc("GET /codes/{registrationCode}/contender", handler.GetContenderByCode)
	mux.HandleFunc("GET /compClasses/{compClassID}/contenders", handler.GetContendersByCompClass)
	mux.HandleFunc("GET /competitions/{competitionID}/contenders", handler.GetContendersByCompetition)
	mux.HandleFunc("PATCH /contenders/{contenderID}", handler.PatchContender)
	mux.HandleFunc("POST /contenders/{contenderID}/scrub", handler.ScrubContender)
	mux.HandleFunc("DELETE /contenders/{contenderID}", handler.DeleteContender)
	mux.HandleFunc("POST /competitions/{competitionID}/contenders", handler.CreateContenders)
}

func (hdlr *contenderHandler) GetPointValues(w http.ResponseWriter, r *http.Request) {
	contenderID, err := parseResourceID[domain.ContenderID](r.PathValue("contenderID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pointValues, err := hdlr.contenderUseCase.GetPointValues(r.Context(), contenderID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, pointValues)
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

func (hdlr *contenderHandler) GetContendersByCompetition(w http.ResponseWriter, r *http.Request) {
	competitionID, err := parseResourceID[domain.CompetitionID](r.PathValue("competitionID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	contenders, err := hdlr.contenderUseCase.GetContendersByCompetition(r.Context(), competitionID)
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
	if !readJSON(w, r, &patch) {
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
	competitionID, err := parseResourceID[domain.CompetitionID](r.PathValue("competitionID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var arguments CreateContendersArguments
	if !readJSON(w, r, &arguments) {
		return
	}

	contenders, err := hdlr.contenderUseCase.CreateContenders(r.Context(), competitionID, arguments.Number)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusCreated, contenders)
}
