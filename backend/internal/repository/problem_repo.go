package repository

import (
	"context"
	"database/sql"

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
