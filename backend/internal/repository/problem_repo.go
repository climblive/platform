package repository

import (
	"context"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

type problemRecord struct {
	ID                 *int `gorm:"primaryKey;autoIncrement"`
	OrganizerID        int
	ContestID          int
	Number             int
	HoldColorPrimary   string
	HoldColorSecondary *string
	Name               *string
	Description        *string
	Points             int
	FlashBonus         *int
}

func (problemRecord) TableName() string {
	return "problem"
}

func (r problemRecord) fromDomain(problem domain.Problem) problemRecord {
	return problemRecord{
		ID:                 e2n(problem.ID),
		OrganizerID:        problem.Ownership.OrganizerID,
		ContestID:          problem.ContestID,
		Number:             problem.Number,
		HoldColorPrimary:   problem.HoldColorPrimary,
		HoldColorSecondary: e2n(problem.HoldColorSecondary),
		Name:               e2n(problem.Name),
		Description:        e2n(problem.Description),
		Points:             problem.PointsTop,
		FlashBonus:         e2n(problem.FlashBonus),
	}
}

func (r *problemRecord) toDomain() domain.Problem {
	return domain.Problem{
		ID: n2e(r.ID),
		Ownership: domain.OwnershipData{
			OrganizerID: r.OrganizerID,
		},
		ContestID:          r.ContestID,
		Number:             r.Number,
		HoldColorPrimary:   r.HoldColorPrimary,
		HoldColorSecondary: n2e(r.HoldColorSecondary),
		Name:               n2e(r.Name),
		Description:        n2e(r.Description),
		PointsTop:          r.Points,
		PointsZone:         0,
		FlashBonus:         n2e(r.FlashBonus),
	}
}

func (d *Database) GetProblemsByContest(ctx context.Context, tx domain.Transaction, contestID domain.ResourceID) ([]domain.Problem, error) {
	var records []problemRecord

	err := d.tx(tx).WithContext(ctx).Raw(`SELECT * FROM problem WHERE contest_id = ?`, contestID).Scan(&records).Error
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	problems := make([]domain.Problem, 0)

	for _, record := range records {
		problems = append(problems, record.toDomain())
	}

	return problems, nil
}

func (d *Database) GetProblem(ctx context.Context, tx domain.Transaction, problemID domain.ResourceID) (domain.Problem, error) {
	var record problemRecord
	err := d.tx(tx).WithContext(ctx).Raw(`SELECT * FROM problem WHERE id = ?`, problemID).Scan(&record).Error
	if err != nil {
		return domain.Problem{}, errors.Wrap(err, 0)
	}

	if record.ID == nil {
		return domain.Problem{}, errors.Wrap(domain.ErrNotFound, 0)
	}

	return record.toDomain(), nil
}
