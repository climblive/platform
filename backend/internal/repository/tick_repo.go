package repository

import (
	"context"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

type tickRecord struct {
	ID          *int `gorm:"primaryKey;autoIncrement"`
	OrganizerID int
	ContestID   int
	ContenderID int
	ProblemID   int
	Flash       bool
	Timestamp   time.Time
}

func (tickRecord) TableName() string {
	return "tick"
}

func (r tickRecord) fromDomain(tick domain.Tick) tickRecord {
	return tickRecord{
		ID:          e2n(tick.ID),
		OrganizerID: tick.Ownership.OrganizerID,
		ContestID:   tick.ContestID,
		ContenderID: *tick.Ownership.ContenderID,
		ProblemID:   tick.ProblemID,
		Flash:       tick.AttemptsTop == 1,
		Timestamp:   tick.Timestamp,
	}
}

func (r *tickRecord) toDomain() domain.Tick {
	attempts := func(isFlash bool) int {
		if isFlash {
			return 1
		}

		return 999
	}

	return domain.Tick{
		ID: n2e(r.ID),
		Ownership: domain.OwnershipData{
			OrganizerID: r.OrganizerID,
			ContenderID: &r.ContenderID,
		},
		Timestamp:    r.Timestamp,
		ContestID:    r.ContestID,
		ProblemID:    r.ProblemID,
		Top:          true,
		AttemptsTop:  attempts(r.Flash),
		Zone:         true,
		AttemptsZone: attempts(r.Flash),
	}
}

func (d *Database) GetTicksByContender(ctx context.Context, tx domain.Transaction, contenderID domain.ResourceID) ([]domain.Tick, error) {
	var records []tickRecord

	err := d.tx(tx).WithContext(ctx).Raw(`SELECT * FROM tick WHERE contender_id = ?`, contenderID).Scan(&records).Error
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	ticks := make([]domain.Tick, 0)

	for _, record := range records {
		ticks = append(ticks, record.toDomain())
	}

	return ticks, nil
}

func (d *Database) StoreTick(ctx context.Context, tx domain.Transaction, tick domain.Tick) (domain.Tick, error) {
	var err error
	var record tickRecord = tickRecord{}.fromDomain(tick)

	err = d.tx(tx).WithContext(ctx).Save(&record).Error
	if err != nil {
		return domain.Tick{}, errors.Wrap(err, 0)
	}

	return record.toDomain(), nil
}

func (d *Database) DeleteTick(ctx context.Context, tx domain.Transaction, tickID domain.ResourceID) error {
	err := d.tx(tx).WithContext(ctx).Exec(`DELETE FROM tick WHERE id = ?`, tickID).Error

	return errors.Wrap(err, 0)
}
