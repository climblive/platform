package rest

import (
	"encoding/json"
	"net/http"

	"github.com/climblive/platform/backend/internal/domain"
)

type contenderHandler struct {
	contenderUseCase domain.ContenderUseCase
}

func InstallContenderHandler(mux *Mux, contenderUseCase domain.ContenderUseCase) {
	handler := &contenderHandler{
		contenderUseCase: contenderUseCase,
	}

	mux.HandleFunc("GET /contenders/{contenderID}", handler.GetContender)
	mux.HandleFunc("GET /codes/{registrationCode}/contender", handler.GetContenderByCode)
	mux.HandleFunc("GET /compClasses/{compClassID}/contenders", handler.GetContendersByCompClass)
	mux.HandleFunc("GET /contests/{contestID}/contenders", handler.GetContendersByContest)
	mux.HandleFunc("PATCH /contenders/{contenderID}", handler.PatchContender)
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

type createContendersTemplate struct {
	Number int `json:"number"`
}

func (hdlr *contenderHandler) CreateContenders(w http.ResponseWriter, r *http.Request) {
	contestID, err := parseResourceID[domain.ContestID](r.PathValue("contestID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var tmpl createContendersTemplate
	err = json.NewDecoder(r.Body).Decode(&tmpl)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	contenders, err := hdlr.contenderUseCase.CreateContenders(r.Context(), contestID, tmpl.Number)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusCreated, contenders)
}
