package rest

import (
	"encoding/json"
	"net/http"

	"github.com/climblive/platform/backend/internal/domain"
)

type contenderHandler struct {
	contenderUseCase domain.ContenderUseCase
}

func InstallContenderHandler(contenderUseCase domain.ContenderUseCase) {
	handler := &contenderHandler{
		contenderUseCase: contenderUseCase,
	}

	http.HandleFunc("GET /contenders/{contenderID}", handler.GetContender)
	http.HandleFunc("GET /codes/{registrationCode}/contender", handler.GetContenderByCode)
	http.HandleFunc("GET /compClasses/{compClassID}/contenders", handler.GetContendersByCompClass)
	http.HandleFunc("GET /contests/{contestID}/contenders", handler.GetContendersByContest)
	http.HandleFunc("PUT /contenders/{contenderID}", handler.UpdateContender)
	http.HandleFunc("DELETE /contenders/{contenderID}", handler.DeleteContender)
	http.HandleFunc("POST /contests/{contestID}/contenders", handler.CreateContenders)
}

func (hdlr *contenderHandler) GetContender(w http.ResponseWriter, r *http.Request) {
	contenderID := parseResourceID(r.PathValue("contenderID"))

	contender, err := hdlr.contenderUseCase.GetContender(r.Context(), contenderID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, contender)
}

func (hdlr *contenderHandler) GetContenderByCode(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("registrationCode")

	contender, err := hdlr.contenderUseCase.GetContenderByCode(r.Context(), code)
	if err != nil {
		writeError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, contender)
}

func (hdlr *contenderHandler) GetContendersByCompClass(w http.ResponseWriter, r *http.Request) {
	compClassID := parseResourceID(r.PathValue("compClassID"))

	contenders, err := hdlr.contenderUseCase.GetContendersByCompClass(r.Context(), compClassID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, contenders)
}

func (hdlr *contenderHandler) GetContendersByContest(w http.ResponseWriter, r *http.Request) {
	contestID := parseResourceID(r.PathValue("contestID"))

	contenders, err := hdlr.contenderUseCase.GetContendersByContest(r.Context(), contestID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, contenders)
}

func (hdlr *contenderHandler) UpdateContender(w http.ResponseWriter, r *http.Request) {
	contenderID := parseResourceID(r.PathValue("contenderID"))

	var contender domain.Contender
	err := json.NewDecoder(r.Body).Decode(&contender)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	updatedContender, err := hdlr.contenderUseCase.UpdateContender(r.Context(), contenderID, contender)
	if err != nil {
		writeError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, updatedContender)
}

func (hdlr *contenderHandler) DeleteContender(w http.ResponseWriter, r *http.Request) {
	contenderID := parseResourceID(r.PathValue("contenderID"))

	err := hdlr.contenderUseCase.DeleteContender(r.Context(), contenderID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeResponse(w, http.StatusNoContent, nil)
}

type createContendersTemplate struct {
	Number int `json:"number"`
}

func (hdlr *contenderHandler) CreateContenders(w http.ResponseWriter, r *http.Request) {
	contestID := parseResourceID(r.PathValue("contestID"))

	var tmpl createContendersTemplate
	err := json.NewDecoder(r.Body).Decode(&tmpl)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	contenders, err := hdlr.contenderUseCase.CreateContenders(r.Context(), contestID, tmpl.Number)
	if err != nil {
		writeError(w, err)
		return
	}

	writeResponse(w, http.StatusCreated, contenders)
}
