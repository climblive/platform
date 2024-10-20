package repository

import (
	"context"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

type contestRecord struct {
	ID                 *int `gorm:"primaryKey;autoIncrement"`
	OrganizerID        int
	Protected          bool
	SeriesID           *int
	Name               string
	Description        *string
	Location           *string
	FinalEnabled       bool
	QualifyingProblems int
	Finalists          int
	Rules              *string
	GracePeriod        int
}

func (contestRecord) TableName() string {
	return "contest"
}

func (r contestRecord) fromDomain(contest domain.Contest) contestRecord {
	return contestRecord{
		ID:                 e2n(int(contest.ID)),
		OrganizerID:        int(contest.Ownership.OrganizerID),
		Protected:          contest.Protected,
		SeriesID:           e2n(int(contest.SeriesID)),
		Name:               contest.Name,
		Description:        e2n(contest.Description),
		Location:           e2n(contest.Location),
		FinalEnabled:       contest.FinalsEnabled,
		QualifyingProblems: contest.QualifyingProblems,
		Finalists:          contest.Finalists,
		Rules:              e2n(contest.Rules),
		GracePeriod:        int(contest.GracePeriod.Seconds()),
	}
}

func (r *contestRecord) toDomain() domain.Contest {
	return domain.Contest{
		ID: domain.ContestID(n2e(r.ID)),
		Ownership: domain.OwnershipData{
			OrganizerID: domain.OrganizerID(r.OrganizerID),
		},
		Protected:          r.Protected,
		SeriesID:           domain.SeriesID(n2e(r.SeriesID)),
		Name:               r.Name,
		Description:        n2e(r.Description),
		Location:           n2e(r.Location),
		FinalsEnabled:      r.FinalEnabled,
		QualifyingProblems: r.QualifyingProblems,
		Finalists:          r.Finalists,
		Rules:              n2e(r.Rules),
		GracePeriod:        time.Duration(r.GracePeriod) * time.Second,
	}
}

func (d *Database) GetContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) (domain.Contest, error) {
	var record contestRecord
	err := d.tx(tx).WithContext(ctx).Raw(`SELECT * FROM contest WHERE id = ?`, contestID).Scan(&record).Error
	if err != nil {
		return domain.Contest{}, errors.Wrap(err, 0)
	}

	if record.ID == nil {
		return domain.Contest{}, errors.Wrap(domain.ErrNotFound, 0)
	}

	return record.toDomain(), nil
}
