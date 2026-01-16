package repository

import (
	"context"
	"database/sql"

	"github.com/climblive/platform/backend/internal/database"
	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
	"github.com/google/uuid"
)

func (d *Database) GetProblemsByContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) ([]domain.Problem, error) {
	records, err := d.WithTx(tx).GetProblemsByContest(ctx, uuid.UUID(contestID))
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
	record, err := d.WithTx(tx).GetProblem(ctx, uuid.UUID(problemID))
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
		ContestID: uuid.UUID(contestID),
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
	if uuid.UUID(problem.ID) == uuid.Nil {
		problem.ID = domain.ProblemID(uuid.New())
	}

	params := database.UpsertProblemParams{
		ID:                 uuid.UUID(problem.ID),
		OrganizerID:        uuid.UUID(problem.Ownership.OrganizerID),
		ContestID:          uuid.UUID(problem.ContestID),
		Number:             int32(problem.Number),
		HoldColorPrimary:   problem.HoldColorPrimary,
		HoldColorSecondary: makeNullString(problem.HoldColorSecondary),
		Description:        makeNullString(problem.Description),
		Zone1Enabled:       problem.Zone1Enabled,
		Zone2Enabled:       problem.Zone2Enabled,
		PointsZone1:        makeNullInt32(int32(problem.PointsZone1)),
		PointsZone2:        makeNullInt32(int32(problem.PointsZone2)),
		PointsTop:          int32(problem.PointsTop),
		FlashBonus:         makeNullInt32(int32(problem.FlashBonus)),
	}

	_, err := d.WithTx(tx).UpsertProblem(ctx, params)
	if err != nil {
		return domain.Problem{}, errors.Wrap(err, 0)
	}

	return problem, err
}

func (d *Database) DeleteProblem(ctx context.Context, tx domain.Transaction, problemID domain.ProblemID) error {
	err := d.WithTx(tx).DeleteProblem(ctx, uuid.UUID(problemID))
	if err != nil {
		return errors.Wrap(err, 0)
	}

	return nil
}
