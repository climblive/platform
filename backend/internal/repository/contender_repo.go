package repository

import (
	"context"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

type contenderRecord struct {
	ID               *int `gorm:"primaryKey;autoIncrement"`
	OrganizerID      int
	ContestID        int
	RegistrationCode string
	Name             *string
	Club             *string
	ClassID          *int
	Entered          *time.Time
	Disqualified     bool
}

func (contenderRecord) TableName() string {
	return "contender"
}

func (r contenderRecord) fromDomain(contender domain.Contender) contenderRecord {
	return contenderRecord{
		ID:               e2n(contender.ID),
		OrganizerID:      contender.Ownership.OrganizerID,
		ContestID:        contender.ContestID,
		RegistrationCode: contender.RegistrationCode,
		Name:             e2n(contender.Name),
		Club:             e2n(contender.ClubName),
		ClassID:          e2n(contender.CompClassID),
		Entered:          contender.Entered,
		Disqualified:     contender.Disqualified,
	}
}

func (r *contenderRecord) toDomain() domain.Contender {
	return domain.Contender{
		ID: n2e(r.ID),
		Ownership: domain.OwnershipData{
			OrganizerID: r.OrganizerID,
			ContenderID: r.ID,
		},
		ContestID:           r.ContestID,
		RegistrationCode:    r.RegistrationCode,
		Name:                n2e(r.Name),
		PublicName:          n2e(r.Name),
		ClubName:            n2e(r.Club),
		CompClassID:         n2e(r.ClassID),
		Entered:             r.Entered,
		WithdrawnFromFinals: false,
		Disqualified:        r.Disqualified,
		Score:               0,
		Placement:           0,
		Finalist:            false,
		ScoreUpdated:        nil,
	}
}

func (d *Database) GetContender(ctx context.Context, tx domain.Transaction, contenderID domain.ResourceID) (domain.Contender, error) {
	var record contenderRecord
	err := d.tx(tx).WithContext(ctx).Raw(`SELECT * FROM contender WHERE id = ?`, contenderID).Scan(&record).Error
	if err != nil {
		return domain.Contender{}, errors.Wrap(err, 0)
	}

	if record.ID == nil {
		return domain.Contender{}, errors.Wrap(domain.ErrNotFound, 0)
	}

	return record.toDomain(), nil
}

func (d *Database) GetContenderByCode(ctx context.Context, tx domain.Transaction, registrationCode string) (domain.Contender, error) {
	var record contenderRecord
	err := d.tx(tx).WithContext(ctx).Raw(`SELECT * FROM contender WHERE registration_code = ?`, registrationCode).Scan(&record).Error
	if err != nil {
		return domain.Contender{}, errors.Wrap(err, 0)
	}

	if record.ID == nil {
		return domain.Contender{}, errors.Wrap(domain.ErrNotFound, 0)
	}

	return record.toDomain(), nil
}

func (d *Database) GetContendersByCompClass(ctx context.Context, tx domain.Transaction, compClassID domain.ResourceID) ([]domain.Contender, error) {
	var records []contenderRecord

	err := d.tx(tx).WithContext(ctx).Raw(`SELECT * FROM contender WHERE class_id = ?`, compClassID).Scan(&records).Error
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	contenders := make([]domain.Contender, 0)

	for _, record := range records {
		contenders = append(contenders, record.toDomain())
	}

	return contenders, nil
}

func (d *Database) GetContendersByContest(ctx context.Context, tx domain.Transaction, contestID domain.ResourceID) ([]domain.Contender, error) {
	var records []contenderRecord

	err := d.tx(tx).WithContext(ctx).Raw(`SELECT * FROM contender WHERE contest_id = ?`, contestID).Scan(&records).Error
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	contenders := make([]domain.Contender, 0)

	for _, record := range records {
		contenders = append(contenders, record.toDomain())
	}

	return contenders, nil
}

func (d *Database) StoreContender(ctx context.Context, tx domain.Transaction, contender domain.Contender) (domain.Contender, error) {
	var err error
	var record contenderRecord = contenderRecord{}.fromDomain(contender)

	err = d.tx(tx).WithContext(ctx).Save(&record).Error
	if err != nil {
		return domain.Contender{}, errors.Wrap(err, 0)
	}

	return record.toDomain(), nil
}

func (d *Database) DeleteContender(ctx context.Context, tx domain.Transaction, contenderID domain.ResourceID) error {
	err := d.tx(tx).WithContext(ctx).Exec(`DELETE FROM contender WHERE id = ?`, contenderID).Error

	return errors.Wrap(err, 0)
}

func (d *Database) GetNumberOfContenders(ctx context.Context, tx domain.Transaction, contestID domain.ResourceID) (int, error) {
	var count int

	err := d.tx(tx).WithContext(ctx).Raw(`SELECT COUNT(*) FROM contender WHERE contest_id = ?`, contestID).Scan(&count).Error
	if err != nil {
		return 0, errors.Wrap(err, 0)
	}

	return count, nil
}
