package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
	"github.com/xuri/excelize/v2"
)

type contestUseCase interface {
	GetContest(ctx context.Context, contestID domain.ContestID) (domain.Contest, error)
	GetAllContests(ctx context.Context) ([]domain.Contest, error)
	GetContestsByOrganizer(ctx context.Context, organizerID domain.OrganizerID) ([]domain.Contest, error)
	GetScoreboard(ctx context.Context, contestID domain.ContestID) ([]domain.ScoreboardEntry, error)
	PatchContest(ctx context.Context, contestID domain.ContestID, patch domain.ContestPatch) (domain.Contest, error)
	CreateContest(ctx context.Context, organizerID domain.OrganizerID, template domain.ContestTemplate) (domain.Contest, error)
	DuplicateContest(ctx context.Context, contestID domain.ContestID) (domain.Contest, error)
	TransferContest(ctx context.Context, contestID domain.ContestID, newOrganizerID domain.OrganizerID) (domain.Contest, error)
}

type contestHandler struct {
	contestUseCase   contestUseCase
	compClassUseCase compClassUseCase
	tickUseCase      tickUseCase
	problemUseCase   problemUseCase
}

func InstallContestHandler(
	mux *Mux,
	contestUseCase contestUseCase,
	compClassUseCase compClassUseCase,
	tickUseCase tickUseCase,
	problemUseCase problemUseCase) {
	handler := &contestHandler{
		contestUseCase:   contestUseCase,
		compClassUseCase: compClassUseCase,
		tickUseCase:      tickUseCase,
		problemUseCase:   problemUseCase,
	}

	mux.HandleFunc("GET /contests/{contestID}", handler.GetContest)
	mux.HandleFunc("GET /contests", handler.GetAllContests)
	mux.HandleFunc("GET /contests/{contestID}/scoreboard", handler.GetScoreboard)
	mux.HandleFunc("GET /organizers/{organizerID}/contests", handler.GetContestsByOrganizer)
	mux.HandleFunc("POST /organizers/{organizerID}/contests", handler.CreateContest)
	mux.HandleFunc("POST /contests/{contestID}/duplicate", handler.DuplicateContest)
	mux.HandleFunc("POST /contests/{contestID}/transfer", handler.TransferContest)
	mux.HandleFunc("GET /contests/{contestID}/results", handler.DownloadResults)
	mux.HandleFunc("PATCH /contests/{contestID}", handler.PatchContest)
}

func (hdlr *contestHandler) GetContest(w http.ResponseWriter, r *http.Request) {
	contestID, err := parseResourceID[domain.ContestID](r.PathValue("contestID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	contest, err := hdlr.contestUseCase.GetContest(r.Context(), contestID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, contest)
}

func (hdlr *contestHandler) GetAllContests(w http.ResponseWriter, r *http.Request) {
	contests, err := hdlr.contestUseCase.GetAllContests(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, contests)
}

func (hdlr *contestHandler) GetScoreboard(w http.ResponseWriter, r *http.Request) {
	contestID, err := parseResourceID[domain.ContestID](r.PathValue("contestID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	scoreboard, err := hdlr.contestUseCase.GetScoreboard(r.Context(), contestID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, scoreboard)
}

func (hdlr *contestHandler) GetContestsByOrganizer(w http.ResponseWriter, r *http.Request) {
	organizerID, err := parseResourceID[domain.OrganizerID](r.PathValue("organizerID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	contests, err := hdlr.contestUseCase.GetContestsByOrganizer(r.Context(), organizerID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, contests)
}

func (hdlr *contestHandler) PatchContest(w http.ResponseWriter, r *http.Request) {
	contestID, err := parseResourceID[domain.ContestID](r.PathValue("contestID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var patch domain.ContestPatch
	err = json.NewDecoder(r.Body).Decode(&patch)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	updatedContest, err := hdlr.contestUseCase.PatchContest(r.Context(), contestID, patch)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, updatedContest)
}

func (hdlr *contestHandler) CreateContest(w http.ResponseWriter, r *http.Request) {
	organizerID, err := parseResourceID[domain.OrganizerID](r.PathValue("organizerID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var tmpl domain.ContestTemplate
	err = json.NewDecoder(r.Body).Decode(&tmpl)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	createdContest, err := hdlr.contestUseCase.CreateContest(r.Context(), organizerID, tmpl)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusCreated, createdContest)
}

func (hdlr *contestHandler) DuplicateContest(w http.ResponseWriter, r *http.Request) {
	contestID, err := parseResourceID[domain.ContestID](r.PathValue("contestID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	duplicatedContest, err := hdlr.contestUseCase.DuplicateContest(r.Context(), contestID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusCreated, duplicatedContest)
}

func (hdlr *contestHandler) TransferContest(w http.ResponseWriter, r *http.Request) {
	contestID, err := parseResourceID[domain.ContestID](r.PathValue("contestID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req domain.ContestTransferRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	transferredContest, err := hdlr.contestUseCase.TransferContest(r.Context(), contestID, req.NewOrganizerID)
	if err != nil {
		handleError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, transferredContest)
}

func (hdlr *contestHandler) DownloadResults(w http.ResponseWriter, r *http.Request) {
	contestID, err := parseResourceID[domain.ContestID](r.PathValue("contestID"))
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
		compClasses, err := hdlr.compClassUseCase.GetCompClassesByContest(r.Context(), contestID)
		if err != nil {
			return errors.Wrap(err, 0)
		}

		scoreboard, err := hdlr.contestUseCase.GetScoreboard(r.Context(), contestID)
		if err != nil {
			return errors.Wrap(err, 0)
		}

		problems, err := hdlr.problemUseCase.GetProblemsByContest(r.Context(), contestID)
		if err != nil {
			return errors.Wrap(err, 0)
		}

		ticks, err := hdlr.tickUseCase.GetTicksByContest(r.Context(), contestID)
		if err != nil {
			return errors.Wrap(err, 0)
		}

		slices.SortFunc(scoreboard, func(a, b domain.ScoreboardEntry) int {
			return a.Score.RankOrder - b.Score.RankOrder
		})

		slices.SortFunc(problems, func(a, b domain.Problem) int {
			return a.Number - b.Number
		})

		for _, compClass := range compClasses {
			if _, err := book.NewSheet(sanitizeSheetName(compClass.Name)); err != nil {
				return errors.Wrap(err, 0)
			}
		}

		nextRowNumbers := make(map[domain.CompClassID]int)

		for _, compClass := range compClasses {
			sheetName := sanitizeSheetName(compClass.Name)

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

			nextRowNumbers[compClass.ID] = 2
		}

		for _, entry := range scoreboard {
			var sheetName string
			counter := nextRowNumbers[entry.CompClassID]

			for _, compClass := range compClasses {
				if entry.CompClassID == compClass.ID {
					sheetName = sanitizeSheetName(compClass.Name)

					break
				}
			}

			err = book.SetSheetRow(sheetName, fmt.Sprintf("A%d", counter), &[]any{
				entry.Name,
				entry.Score.Score,
				entry.Score.Placement})
			if err != nil {
				return errors.Wrap(err, 0)
			}

			problemResults := make(map[domain.ProblemID]string, 0)
			for _, tick := range ticks {
				if *tick.Ownership.ContenderID == entry.ContenderID {
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

			nextRowNumbers[entry.CompClassID]++
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
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="contest_%d_results.xlsx"`, contestID))
	w.WriteHeader(http.StatusOK)

	if _, err := book.WriteTo(w); err != nil {
		handleError(w, err)
	}
}
