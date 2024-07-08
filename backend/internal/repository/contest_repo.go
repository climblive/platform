package repository

import (
	"context"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
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

func (r contestRecord) FromDomain(contest domain.Contest) contestRecord {
	return contestRecord{
		ID:                 e2n(contest.ID),
		OrganizerID:        contest.Ownership.OrganizerID,
		Protected:          contest.Protected,
		SeriesID:           e2n(contest.SeriesID),
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

func (r *contestRecord) ToDomain() domain.Contest {
	return domain.Contest{
		ID: n2e(r.ID),
		Ownership: domain.OwnershipData{
			OrganizerID: r.OrganizerID,
		},
		Protected:          r.Protected,
		SeriesID:           n2e(r.SeriesID),
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

func (d *Database) GetContest(ctx context.Context, tx domain.Transaction, contestID domain.ResourceID) (domain.Contest, error) {
	transaction, ok := tx.(*transaction)
	if !ok {
		return domain.Contest{}, ErrIncompatibleTransaction
	}

	var record contestRecord
	err := transaction.db.WithContext(ctx).Raw(`SELECT * FROM contest WHERE id = ?`, contestID).Scan(&record).Error

	return record.ToDomain(), err
}
