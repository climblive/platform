package repository

import (
	"context"
	"database/sql"

	"github.com/climblive/platform/backend/internal/database"
	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

func (d *Database) GetProblemsByContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) ([]domain.Problem, error) {
	records, err := d.WithTx(tx).GetProblemsByContest(ctx, int32(contestID))
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	problems := make([]domain.Problem, 0)

	for _, record := range records {
		problems = append(problems, problemToDomain(record.Problem))
	}

	return problems, nil
}

func (d *Database) GetProblem(ctx context.Context, tx domain.Transaction, problemID domain.ProblemID) (domain.Problem, error) {
	record, err := d.WithTx(tx).GetProblem(ctx, int32(problemID))
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return domain.Problem{}, errors.Wrap(domain.ErrNotFound, 0)
	case err != nil:
		return domain.Problem{}, errors.Wrap(err, 0)
	}

	return problemToDomain(record.Problem), nil
}

func (d *Database) GetProblemByNumber(ctx context.Context, tx domain.Transaction, contestID domain.ContestID, problemNumber int) (domain.Problem, error) {
	record, err := d.WithTx(tx).GetProblemByNumber(ctx, database.GetProblemByNumberParams{
		ContestID: int32(contestID),
		Number:    int32(problemNumber),
	})
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return domain.Problem{}, errors.Wrap(domain.ErrNotFound, 0)
	case err != nil:
		return domain.Problem{}, errors.Wrap(err, 0)
	}

	return problemToDomain(record.Problem), nil
}

func (d *Database) StoreProblem(ctx context.Context, tx domain.Transaction, problem domain.Problem) (domain.Problem, error) {
	params := database.UpsertProblemParams{
		ID:                 int32(problem.ID),
		OrganizerID:        int32(problem.Ownership.OrganizerID),
		ContestID:          int32(problem.ContestID),
		Number:             int32(problem.Number),
		HoldColorPrimary:   problem.HoldColorPrimary,
		HoldColorSecondary: makeNullString(problem.HoldColorSecondary),
		Name:               makeNullString(problem.Name),
		Description:        makeNullString(problem.Description),
		Points:             int32(problem.PointsTop),
		FlashBonus:         makeNullInt32(int32(problem.FlashBonus)),
	}

	insertID, err := d.WithTx(tx).UpsertProblem(ctx, params)
	if err != nil {
		return domain.Problem{}, errors.Wrap(err, 0)
	}

	if insertID != 0 {
		problem.ID = domain.ProblemID(insertID)
	}

	return problem, err
}
