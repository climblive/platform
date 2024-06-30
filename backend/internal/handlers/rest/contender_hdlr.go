package rest

import (
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
	http.HandleFunc("GET /contest/{contestID}/contenders", handler.GetContendersByContest)
	http.HandleFunc("PUT /contenders/{contenderID}", handler.UpdateContender)
	http.HandleFunc("DELETE /contenders/{contenderID}", handler.DeleteContender)
	http.HandleFunc("POST /contenders", handler.CreateContenders)
}

func (hdlr *contenderHandler) GetContender(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (hdlr *contenderHandler) GetContenderByCode(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (hdlr *contenderHandler) GetContendersByCompClass(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (hdlr *contenderHandler) GetContendersByContest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (hdlr *contenderHandler) UpdateContender(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (hdlr *contenderHandler) DeleteContender(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (hdlr *contenderHandler) CreateContenders(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
