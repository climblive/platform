package rest

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
	"github.com/xuri/excelize/v2"
)

type competitionUseCase interface {
	GetCompetition(ctx context.Context, competitionID domain.CompetitionID) (domain.Competition, error)
	GetAllCompetitions(ctx context.Context) ([]domain.Competition, error)
	GetCompetitionsByOrganizer(ctx context.Context, organizerID domain.OrganizerID) ([]domain.Competition, error)
	GetScoreboard(ctx context.Context, competitionID domain.CompetitionID) ([]domain.ScoreboardEntry, error)
	PatchCompetition(ctx context.Context, competitionID domain.CompetitionID, patch domain.CompetitionPatch) (domain.Competition, error)
	ArchiveCompetition(ctx context.Context, competitionID domain.CompetitionID) (domain.Competition, error)
	RestoreCompetition(ctx context.Context, competitionID domain.CompetitionID) (domain.Competition, error)
	CreateCompetition(ctx context.Context, organizerID domain.OrganizerID, template domain.CompetitionTemplate) (domain.Competition, error)
	DuplicateCompetition(ctx context.Context, competitionID domain.CompetitionID) (domain.Competition, error)
	TransferCompetition(ctx context.Context, competitionID domain.CompetitionID, newOrganizerID domain.OrganizerID) (domain.Competition, error)
}

type competitionHandler struct {
	competitionUseCase competitionUseCase
	categoryUseCase    categoryUseCase
	tickUseCase        tickUseCase
	problemUseCase     problemUseCase
}

func InstallCompetitionHandler(
	mux *Mux,
	competitionUseCase competitionUseCase,
	categoryUseCase categoryUseCase,
	tickUseCase tickUseCase,
	problemUseCase problemUseCase) {
	handler := &competitionHandler{
		competitionUseCase: competitionUseCase,
		categoryUseCase:    categoryUseCase,
		tickUseCase:        tickUseCase,
		problemUseCase:     problemUseCase,
	}

	mux.HandleFunc("GET /competitions/{competitionID}", handler.GetCompetition)
	mux.HandleFunc("GET /competitions", handler.GetAllCompetitions)
	mux.HandleFunc("GET /competitions/{competitionID}/scoreboard", handler.GetScoreboard)
	mux.HandleFunc("GET /organizers/{organizerID}/competitions", handler.GetCompetitionsByOrganizer)
	mux.HandleFunc("POST /organizers/{organizerID}/competitions", handler.CreateCompetition)
	mux.HandleFunc("POST /competitions/{competitionID}/duplicate", handler.DuplicateCompetition)
	mux.HandleFunc("POST /competitions/{competitionID}/transfer", handler.TransferCompetition)
	mux.HandleFunc("POST /competitions/{competitionID}/archive", handler.ArchiveCompetition)
	mux.HandleFunc("POST /competitions/{competitionID}/restore", handler.RestoreCompetition)
	mux.HandleFunc("GET /competitions/{competitionID}/results", handler.DownloadResults)
	mux.HandleFunc("PATCH /competitions/{competitionID}", handler.PatchCompetition)
}

func (hdlr *competitionHandler) GetCompetition(w http.ResponseWriter, r *http.Request) {
	competitionID, err := parseResourceID[domain.CompetitionID](r.PathValue("competitionID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	competition, err := hdlr.competitionUseCase.GetCompetition(r.Context(), competitionID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, competition)
}

func (hdlr *competitionHandler) GetAllCompetitions(w http.ResponseWriter, r *http.Request) {
	competitions, err := hdlr.competitionUseCase.GetAllCompetitions(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, competitions)
}

func (hdlr *competitionHandler) GetScoreboard(w http.ResponseWriter, r *http.Request) {
	competitionID, err := parseResourceID[domain.CompetitionID](r.PathValue("competitionID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	scoreboard, err := hdlr.competitionUseCase.GetScoreboard(r.Context(), competitionID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, scoreboard)
}

func (hdlr *competitionHandler) GetCompetitionsByOrganizer(w http.ResponseWriter, r *http.Request) {
	organizerID, err := parseResourceID[domain.OrganizerID](r.PathValue("organizerID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	competitions, err := hdlr.competitionUseCase.GetCompetitionsByOrganizer(r.Context(), organizerID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, competitions)
}

func (hdlr *competitionHandler) PatchCompetition(w http.ResponseWriter, r *http.Request) {
	competitionID, err := parseResourceID[domain.CompetitionID](r.PathValue("competitionID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var patch domain.CompetitionPatch
	if !readJSON(w, r, &patch) {
		return
	}

	updatedCompetition, err := hdlr.competitionUseCase.PatchCompetition(r.Context(), competitionID, patch)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, updatedCompetition)
}

func (hdlr *competitionHandler) CreateCompetition(w http.ResponseWriter, r *http.Request) {
	organizerID, err := parseResourceID[domain.OrganizerID](r.PathValue("organizerID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var tmpl domain.CompetitionTemplate
	if !readJSON(w, r, &tmpl) {
		return
	}

	createdCompetition, err := hdlr.competitionUseCase.CreateCompetition(r.Context(), organizerID, tmpl)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusCreated, createdCompetition)
}

func (hdlr *competitionHandler) DuplicateCompetition(w http.ResponseWriter, r *http.Request) {
	competitionID, err := parseResourceID[domain.CompetitionID](r.PathValue("competitionID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	duplicatedCompetition, err := hdlr.competitionUseCase.DuplicateCompetition(r.Context(), competitionID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusCreated, duplicatedCompetition)
}

func (hdlr *competitionHandler) ArchiveCompetition(w http.ResponseWriter, r *http.Request) {
	competitionID, err := parseResourceID[domain.CompetitionID](r.PathValue("competitionID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	competition, err := hdlr.competitionUseCase.ArchiveCompetition(r.Context(), competitionID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, competition)
}

func (hdlr *competitionHandler) RestoreCompetition(w http.ResponseWriter, r *http.Request) {
	competitionID, err := parseResourceID[domain.CompetitionID](r.PathValue("competitionID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	competition, err := hdlr.competitionUseCase.RestoreCompetition(r.Context(), competitionID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, competition)
}

func (hdlr *competitionHandler) TransferCompetition(w http.ResponseWriter, r *http.Request) {
	competitionID, err := parseResourceID[domain.CompetitionID](r.PathValue("competitionID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req domain.CompetitionTransferRequest
	if !readJSON(w, r, &req) {
		return
	}

	transferredCompetition, err := hdlr.competitionUseCase.TransferCompetition(r.Context(), competitionID, req.NewOrganizerID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, transferredCompetition)
}

func (hdlr *competitionHandler) DownloadResults(w http.ResponseWriter, r *http.Request) {
	competitionID, err := parseResourceID[domain.CompetitionID](r.PathValue("competitionID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	book := excelize.NewFile()
	defer func() {
		if err := book.Close(); err != nil {
			handleError(w, err)
			return
		}
	}()

	sanitizeSheetName := func(name string) string {
		invalidCharacters := ":\\/?*[]"

		for _, char := range invalidCharacters {
			name = strings.ReplaceAll(name, string(char), "")
		}

		if len(name) > excelize.MaxSheetNameLength {
			name = name[0:excelize.MaxSheetNameLength]
		}

		return name
	}

	writeBook := func(book *excelize.File) error {
		categories, err := hdlr.categoryUseCase.GetCategoriesByCompetition(r.Context(), competitionID)
		if err != nil {
			return errors.Wrap(err, 0)
		}

		scoreboard, err := hdlr.competitionUseCase.GetScoreboard(r.Context(), competitionID)
		if err != nil {
			return errors.Wrap(err, 0)
		}

		problems, err := hdlr.problemUseCase.GetProblemsByCompetition(r.Context(), competitionID)
		if err != nil {
			return errors.Wrap(err, 0)
		}

		ticks, err := hdlr.tickUseCase.GetTicksByCompetition(r.Context(), competitionID)
		if err != nil {
			return errors.Wrap(err, 0)
		}

		slices.SortFunc(scoreboard, func(a, b domain.ScoreboardEntry) int {
			return a.Result.RankOrder - b.Result.RankOrder
		})

		slices.SortFunc(problems, func(a, b domain.Problem) int {
			return a.Number - b.Number
		})

		for _, category := range categories {
			if _, err := book.NewSheet(sanitizeSheetName(category.Name)); err != nil {
				return errors.Wrap(err, 0)
			}
		}

		nextRowNumbers := make(map[domain.CategoryID]int)

		for _, category := range categories {
			sheetName := sanitizeSheetName(category.Name)

			style, err := book.NewStyle(&excelize.Style{
				Font: &excelize.Font{
					Bold:         true,
					Italic:       false,
					Underline:    "",
					Family:       "",
					Size:         0,
					Strike:       false,
					Color:        "",
					ColorIndexed: 0,
					ColorTheme:   nil,
					ColorTint:    0,
					VertAlign:    "",
					Charset:      nil,
				},
				Border:        nil,
				Fill:          excelize.Fill{},
				Alignment:     nil,
				Protection:    nil,
				NumFmt:        0,
				DecimalPlaces: nil,
				CustomNumFmt:  nil,
				NegRed:        false,
			})
			if err != nil {
				return errors.Wrap(err, 0)
			}

			err = book.SetColWidth(sheetName, "A", "A", 40)
			if err != nil {
				return errors.Wrap(err, 0)
			}

			err = book.SetColWidth(sheetName, "B", "C", 20)
			if err != nil {
				return errors.Wrap(err, 0)
			}

			lastStyledCell, err := excelize.CoordinatesToCellName(3+len(problems), 1)
			if err != nil {
				return errors.Wrap(err, 0)
			}

			err = book.SetCellStyle(sheetName, "A1", lastStyledCell, style)
			if err != nil {
				return errors.Wrap(err, 0)
			}

			err = book.SetSheetRow(sheetName, "A1", &[]string{"Name", "Score", "Placement"})
			if err != nil {
				return errors.Wrap(err, 0)
			}

			problemNumbers := make([]string, 0)
			for _, problem := range problems {
				problemNumbers = append(problemNumbers, fmt.Sprintf("P%d", problem.Number))
			}

			err = book.SetSheetRow(sheetName, "D1", &problemNumbers)
			if err != nil {
				return errors.Wrap(err, 0)
			}

			nextRowNumbers[category.ID] = 2
		}

		for _, entry := range scoreboard {
			var sheetName string
			counter := nextRowNumbers[entry.CategoryID]

			for _, category := range categories {
				if entry.CategoryID == category.ID {
					sheetName = sanitizeSheetName(category.Name)

					break
				}
			}

			var (
				score     string
				placement int
			)
			if entry.Result != nil {
				score = entry.Result.Score
				placement = entry.Result.Placement
			}

			err = book.SetSheetRow(sheetName, fmt.Sprintf("A%d", counter), &[]any{
				entry.Name,
				score,
				placement})
			if err != nil {
				return errors.Wrap(err, 0)
			}

			problemResults := make(map[domain.ProblemID]string, 0)
			for _, tick := range ticks {
				if *tick.Ownership.CompetitorID == entry.CompetitorID {
					result := ""

					switch {
					case tick.Top && tick.AttemptsTop == 1:
						result = "F"
					case tick.Top:
						result = "T"
					}

					problemResults[tick.ProblemID] = result
				}
			}

			resultsRow := make([]string, 0)
			for _, problem := range problems {
				resultsRow = append(resultsRow, problemResults[problem.ID])
			}

			err = book.SetSheetRow(sheetName, fmt.Sprintf("D%d", counter), &resultsRow)
			if err != nil {
				return errors.Wrap(err, 0)
			}

			nextRowNumbers[entry.CategoryID]++
		}

		err = book.DeleteSheet("Sheet1")
		if err != nil {
			return errors.Wrap(err, 0)
		}

		return nil
	}

	err = writeBook(book)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="competition_%d_results.xlsx"`, competitionID))
	w.WriteHeader(http.StatusOK)

	if _, err := book.WriteTo(w); err != nil {
		handleError(w, err)
	}
}
