package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/xuri/excelize/v2"
)

type contestUseCase interface {
	GetContest(ctx context.Context, contestID domain.ContestID) (domain.Contest, error)
	GetContestsByOrganizer(ctx context.Context, organizerID domain.OrganizerID) ([]domain.Contest, error)
	GetScoreboard(ctx context.Context, contestID domain.ContestID) ([]domain.ScoreboardEntry, error)
	CreateContest(ctx context.Context, organizerID domain.OrganizerID, template domain.ContestTemplate) (domain.Contest, error)
	DuplicateContest(ctx context.Context, contestID domain.ContestID) (domain.Contest, error)
}

type contestHandler struct {
	contestUseCase   contestUseCase
	compClassUseCase compClassUseCase
}

func InstallContestHandler(mux *Mux, contestUseCase contestUseCase, compcompClassUseCase compClassUseCase) {
	handler := &contestHandler{
		contestUseCase:   contestUseCase,
		compClassUseCase: compcompClassUseCase,
	}

	mux.HandleFunc("GET /contests/{contestID}", handler.GetContest)
	mux.HandleFunc("GET /contests/{contestID}/scoreboard", handler.GetScoreboard)
	mux.HandleFunc("GET /organizers/{organizerID}/contests", handler.GetContestsByOrganizer)
	mux.HandleFunc("POST /organizers/{organizerID}/contests", handler.CreateContest)
	mux.HandleFunc("POST /contests/{contestID}/duplicate", handler.DuplicateContest)
	mux.HandleFunc("GET /contests/{contestID}/results", handler.DownloadResults)
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

	writeBook := func(book *excelize.File) error {
		compClasses, err := hdlr.compClassUseCase.GetCompClassesByContest(r.Context(), contestID)
		if err != nil {
			return err
		}

		scoreboard, err := hdlr.contestUseCase.GetScoreboard(r.Context(), contestID)
		if err != nil {
			return err
		}

		slices.SortFunc(scoreboard, func(a, b domain.ScoreboardEntry) int {
			return a.Score.RankOrder - b.Score.RankOrder
		})

		for _, compClass := range compClasses {
			if _, err := book.NewSheet(compClass.Name); err != nil {
				return err
			}
		}

		nextRowNumbers := make(map[domain.CompClassID]int)

		for _, compClass := range compClasses {
			sheetName := compClass.Name

			err = book.SetColWidth(sheetName, "A", "B", 40)
			if err != nil {
				return err
			}

			err = book.SetCellValue(sheetName, "A1", "Name")
			if err != nil {
				return err
			}

			err = book.SetCellValue(sheetName, "B1", "Club")
			if err != nil {
				return err
			}

			err = book.SetCellValue(sheetName, "C1", "Score")
			if err != nil {
				return err
			}

			err = book.SetCellValue(sheetName, "D1", "Placement")
			if err != nil {
				return err
			}

			nextRowNumbers[compClass.ID] = 2
		}

		for _, entry := range scoreboard {
			var sheetName string
			counter := nextRowNumbers[entry.CompClassID]

			for _, compClass := range compClasses {
				if entry.CompClassID == compClass.ID {
					sheetName = compClass.Name

					break
				}
			}

			err = book.SetCellValue(sheetName, fmt.Sprintf("A%d", counter), entry.PublicName)
			if err != nil {
				return err
			}

			err = book.SetCellValue(sheetName, fmt.Sprintf("B%d", counter), entry.ClubName)
			if err != nil {
				return err
			}

			err = book.SetCellValue(sheetName, fmt.Sprintf("C%d", counter), entry.Score.Score)
			if err != nil {
				return err
			}

			err = book.SetCellValue(sheetName, fmt.Sprintf("D%d", counter), entry.Score.Placement)
			if err != nil {
				return err
			}

			nextRowNumbers[entry.CompClassID]++
		}

		err = book.DeleteSheet("Sheet1")
		if err != nil {
			return err
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
