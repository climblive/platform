package repository

import (
	"context"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
)

type contestRecord struct {
	ID                 *int    `gorm:"column:id,primaryKey,autoIncrement:true"`
	OrganizerID        int     `gorm:"column:organizer_id"`
	Protected          bool    `gorm:"column:protected"`
	SeriesID           *int    `gorm:"column:series_id"`
	Name               string  `gorm:"column:name"`
	Description        *string `gorm:"column:description"`
	Location           *string `gorm:"column:location"`
	FinalEnabled       bool    `gorm:"column:final_enabled"`
	QualifyingProblems int     `gorm:"column:qualifying_problems"`
	Finalists          int     `gorm:"column:finalists"`
	Rules              *string `gorm:"column:rules"`
	GracePeriod        int     `gorm:"column:grace_period"`
}

func (r contestRecord) ToRecord(contest domain.Contest) contestRecord {
	return contestRecord{
		ID:                 emptyAsNil(contest.ID),
		OrganizerID:        contest.Ownership.OrganizerID,
		Protected:          contest.Protected,
		SeriesID:           emptyAsNil(contest.SeriesID),
		Name:               contest.Name,
		Description:        emptyAsNil(contest.Description),
		Location:           emptyAsNil(contest.Location),
		FinalEnabled:       contest.FinalEnabled,
		QualifyingProblems: contest.QualifyingProblems,
		Finalists:          contest.Finalists,
		Rules:              emptyAsNil(contest.Rules),
		GracePeriod:        int(contest.GracePeriod.Seconds()),
	}
}

func (r *contestRecord) ToDomain() domain.Contest {
	return domain.Contest{
		ID: nilAsEmpty(r.ID),
		Ownership: domain.OwnershipData{
			OrganizerID: r.OrganizerID,
		},
		Location:           nilAsEmpty(r.Location),
		SeriesID:           nilAsEmpty(r.SeriesID),
		Protected:          r.Protected,
		Name:               r.Name,
		Description:        nilAsEmpty(r.Description),
		FinalEnabled:       r.FinalEnabled,
		QualifyingProblems: r.QualifyingProblems,
		Finalists:          r.Finalists,
		Rules:              nilAsEmpty(r.Rules),
		GracePeriod:        time.Duration(r.GracePeriod) * time.Second,
	}
}

func (d *Database) GetContest(ctx context.Context, tx domain.Transaction, contestID domain.ResourceID) (domain.Contest, error) {
	var record contestRecord
	err := d.db.WithContext(ctx).Raw(`SELECT * FROM contest WHERE id = ?`, contestID).Scan(&record).Error

	return record.ToDomain(), err
}
