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
		ID:          e2n(int(tick.ID)),
		OrganizerID: int(tick.Ownership.OrganizerID),
		ContestID:   int(tick.ContestID),
		ContenderID: int(*tick.Ownership.ContenderID),
		ProblemID:   int(tick.ProblemID),
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
		ID: domain.TickID(n2e(r.ID)),
		Ownership: domain.OwnershipData{
			OrganizerID: domain.OrganizerID(r.OrganizerID),
			ContenderID: nillableIntToResourceID[domain.ContenderID](&r.ContenderID),
		},
		Timestamp:    r.Timestamp,
		ContestID:    domain.ContestID(r.ContestID),
		ProblemID:    domain.ProblemID(r.ProblemID),
		Top:          true,
		AttemptsTop:  attempts(r.Flash),
		Zone:         true,
		AttemptsZone: attempts(r.Flash),
	}
}

func (d *Database) GetTicksByContender(ctx context.Context, tx domain.Transaction, contenderID domain.ContenderID) ([]domain.Tick, error) {
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

func (d *Database) DeleteTick(ctx context.Context, tx domain.Transaction, tickID domain.TickID) error {
	err := d.tx(tx).WithContext(ctx).Exec(`DELETE FROM tick WHERE id = ?`, tickID).Error
	if err != nil {
		return errors.Wrap(err, 0)
	}

	return nil
}

func (d *Database) GetTick(ctx context.Context, tx domain.Transaction, tickID domain.TickID) (domain.Tick, error) {
	var record tickRecord
	err := d.tx(tx).WithContext(ctx).Raw(`SELECT * FROM tick WHERE id = ?`, tickID).Scan(&record).Error
	if err != nil {
		return domain.Tick{}, errors.Wrap(err, 0)
	}

	if record.ID == nil {
		return domain.Tick{}, errors.Wrap(domain.ErrNotFound, 0)
	}

	return record.toDomain(), nil
}
