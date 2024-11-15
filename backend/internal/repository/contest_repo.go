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
	TimeBegin          *time.Time
	TimeEnd            *time.Time
}

func (contestRecord) TableName() string {
	return "contest"
}

//nolint:unused // will be used in future versions
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
		TimeBegin:          r.TimeBegin,
		TimeEnd:            r.TimeEnd,
	}
}

func (d *Database) GetContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) (domain.Contest, error) {
	var record contestRecord
	err := d.tx(tx).WithContext(ctx).Raw(`SELECT contest.*, MIN(cc.time_begin) AS time_begin, MAX(cc.time_end) AS time_end FROM contest
LEFT JOIN comp_class cc ON cc.contest_id = contest.id
WHERE contest.id = ?`, contestID).Scan(&record).Error
	if err != nil {
		return domain.Contest{}, errors.Wrap(err, 0)
	}

	if record.ID == nil {
		return domain.Contest{}, errors.Wrap(domain.ErrNotFound, 0)
	}

	return record.toDomain(), nil
}

func (d *Database) GetContestsCurrentlyRunningOrByStartTime(ctx context.Context, tx domain.Transaction, earliestStartTime, latestStartTime time.Time) ([]domain.Contest, error) {
	var records []contestRecord

	err := d.tx(tx).WithContext(ctx).Raw(`SELECT contest.*, MIN(cc.time_begin) AS time_begin, MAX(cc.time_end) AS time_end FROM contest
JOIN comp_class cc ON cc.contest_id = contest.id
HAVING
	? BETWEEN MIN(cc.time_begin) AND MAX(cc.time_end)
	OR MIN(cc.time_begin) BETWEEN ? AND ?`, time.Now(), earliestStartTime, latestStartTime).Scan(&records).Error
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	contests := make([]domain.Contest, 0)

	for _, record := range records {
		contests = append(contests, record.toDomain())
	}

	return contests, nil
}
