package repository

import (
	"context"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

type compClassRecord struct {
	ID          *int `gorm:"primaryKey;autoIncrement"`
	OrganizerID int
	ContestID   int
	Name        string
	Description *string
	Color       *string
	TimeBegin   time.Time
	TimeEnd     time.Time
}

func (compClassRecord) TableName() string {
	return "comp_class"
}

func (r compClassRecord) fromDomain(compClass domain.CompClass) compClassRecord {
	return compClassRecord{
		ID:          e2n(int(compClass.ID)),
		OrganizerID: int(compClass.Ownership.OrganizerID),
		ContestID:   int(compClass.ContestID),
		Name:        compClass.Name,
		Description: e2n(compClass.Description),
		Color:       e2n(string(compClass.Color)),
		TimeBegin:   compClass.TimeBegin,
		TimeEnd:     compClass.TimeEnd,
	}
}

func (r *compClassRecord) toDomain() domain.CompClass {
	return domain.CompClass{
		ID: domain.CompClassID(n2e(r.ID)),
		Ownership: domain.OwnershipData{
			OrganizerID: domain.OrganizerID(r.OrganizerID),
		},
		ContestID:   domain.ContestID(r.ContestID),
		Name:        r.Name,
		Description: n2e(r.Description),
		Color:       domain.ColorRGB(n2e(r.Color)),
		TimeBegin:   r.TimeBegin,
		TimeEnd:     r.TimeEnd,
	}
}

func (d *Database) GetCompClass(ctx context.Context, tx domain.Transaction, compClassID domain.CompClassID) (domain.CompClass, error) {
	var record compClassRecord
	err := d.tx(tx).WithContext(ctx).Raw(`SELECT * FROM comp_class WHERE id = ?`, compClassID).Scan(&record).Error
	if err != nil {
		return domain.CompClass{}, errors.Wrap(err, 0)
	}

	if record.ID == nil {
		return domain.CompClass{}, errors.Wrap(domain.ErrNotFound, 0)
	}

	return record.toDomain(), nil
}

func (d *Database) GetCompClassesByContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) ([]domain.CompClass, error) {
	var records []compClassRecord

	err := d.tx(tx).WithContext(ctx).Raw(`SELECT * FROM comp_class WHERE contest_id = ?`, contestID).Scan(&records).Error
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	compClasses := make([]domain.CompClass, 0)

	for _, record := range records {
		compClasses = append(compClasses, record.toDomain())
	}

	return compClasses, nil
}
