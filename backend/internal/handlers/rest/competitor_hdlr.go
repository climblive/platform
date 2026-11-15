package rest

import (
	"context"
	"net/http"

	"github.com/climblive/platform/backend/internal/domain"
)

type competitorUseCase interface {
	GetCompetitor(ctx context.Context, competitorID domain.CompetitorID) (domain.Competitor, error)
	GetCompetitorByCode(ctx context.Context, registrationCode string) (domain.Competitor, error)
	GetPointValues(ctx context.Context, competitorID domain.CompetitorID) ([]domain.PointValue, error)
	GetCompetitorsByCategory(ctx context.Context, categoryID domain.CategoryID) ([]domain.Competitor, error)
	GetCompetitorsByCompetition(ctx context.Context, competitionID domain.CompetitionID) ([]domain.Competitor, error)
	PatchCompetitor(ctx context.Context, competitorID domain.CompetitorID, patch domain.CompetitorPatch) (domain.Competitor, error)
	ScrubCompetitor(ctx context.Context, competitorID domain.CompetitorID) (domain.Competitor, error)
	DeleteCompetitor(ctx context.Context, competitorID domain.CompetitorID) error
	CreateCompetitors(ctx context.Context, competitionID domain.CompetitionID, number int) ([]domain.Competitor, error)
}

type competitorHandler struct {
	competitorUseCase competitorUseCase
}

func InstallCompetitorHandler(mux *Mux, competitorUseCase competitorUseCase) {
	handler := &competitorHandler{
		competitorUseCase: competitorUseCase,
	}

	mux.HandleFunc("GET /competitors/{competitorID}", handler.GetCompetitor)
	mux.HandleFunc("GET /competitors/{competitorID}/points", handler.GetPointValues)
	mux.HandleFunc("GET /codes/{registrationCode}/competitor", handler.GetCompetitorByCode)
	mux.HandleFunc("GET /categories/{categoryID}/competitors", handler.GetCompetitorsByCategory)
	mux.HandleFunc("GET /competitions/{competitionID}/competitors", handler.GetCompetitorsByCompetition)
	mux.HandleFunc("PATCH /competitors/{competitorID}", handler.PatchCompetitor)
	mux.HandleFunc("POST /competitors/{competitorID}/scrub", handler.ScrubCompetitor)
	mux.HandleFunc("DELETE /competitors/{competitorID}", handler.DeleteCompetitor)
	mux.HandleFunc("POST /competitions/{competitionID}/competitors", handler.CreateCompetitors)
}

func (hdlr *competitorHandler) GetPointValues(w http.ResponseWriter, r *http.Request) {
	competitorID, err := parseResourceID[domain.CompetitorID](r.PathValue("competitorID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pointValues, err := hdlr.competitorUseCase.GetPointValues(r.Context(), competitorID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, pointValues)
}

func (hdlr *competitorHandler) GetCompetitor(w http.ResponseWriter, r *http.Request) {
	competitorID, err := parseResourceID[domain.CompetitorID](r.PathValue("competitorID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	competitor, err := hdlr.competitorUseCase.GetCompetitor(r.Context(), competitorID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, competitor)
}

func (hdlr *competitorHandler) GetCompetitorByCode(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("registrationCode")

	competitor, err := hdlr.competitorUseCase.GetCompetitorByCode(r.Context(), code)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, competitor)
}

func (hdlr *competitorHandler) GetCompetitorsByCategory(w http.ResponseWriter, r *http.Request) {
	categoryID, err := parseResourceID[domain.CategoryID](r.PathValue("categoryID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	competitors, err := hdlr.competitorUseCase.GetCompetitorsByCategory(r.Context(), categoryID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, competitors)
}

func (hdlr *competitorHandler) GetCompetitorsByCompetition(w http.ResponseWriter, r *http.Request) {
	competitionID, err := parseCompetitionID(r.PathValue("competitionID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	competitors, err := hdlr.competitorUseCase.GetCompetitorsByCompetition(r.Context(), competitionID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, competitors)
}

func (hdlr *competitorHandler) PatchCompetitor(w http.ResponseWriter, r *http.Request) {
	competitorID, err := parseResourceID[domain.CompetitorID](r.PathValue("competitorID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var patch domain.CompetitorPatch
	if !readJSON(w, r, &patch) {
		return
	}

	updatedCompetitor, err := hdlr.competitorUseCase.PatchCompetitor(r.Context(), competitorID, patch)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, updatedCompetitor)
}

func (hdlr *competitorHandler) DeleteCompetitor(w http.ResponseWriter, r *http.Request) {
	competitorID, err := parseResourceID[domain.CompetitorID](r.PathValue("competitorID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = hdlr.competitorUseCase.DeleteCompetitor(r.Context(), competitorID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusNoContent, nil)
}

func (hdlr *competitorHandler) ScrubCompetitor(w http.ResponseWriter, r *http.Request) {
	competitorID, err := parseResourceID[domain.CompetitorID](r.PathValue("competitorID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	competitor, err := hdlr.competitorUseCase.ScrubCompetitor(r.Context(), competitorID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, competitor)
}

func (hdlr *competitorHandler) CreateCompetitors(w http.ResponseWriter, r *http.Request) {
	competitionID, err := parseCompetitionID(r.PathValue("competitionID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var arguments CreateCompetitorsArguments
	if !readJSON(w, r, &arguments) {
		return
	}

	competitors, err := hdlr.competitorUseCase.CreateCompetitors(r.Context(), competitionID, arguments.Number)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusCreated, competitors)
}
