package repository

import (
	"context"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
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

func (r compClassRecord) fromDomain(compClass domain.CompClass) compClassRecord {
	return compClassRecord{
		ID:          e2n(compClass.ID),
		OrganizerID: compClass.Ownership.OrganizerID,
		ContestID:   compClass.ContestID,
		Name:        compClass.Name,
		Description: e2n(compClass.Description),
		Color:       e2n(string(compClass.Color)),
		TimeBegin:   compClass.TimeBegin,
		TimeEnd:     compClass.TimeEnd,
	}
}

func (r *compClassRecord) toDomain() domain.CompClass {
	return domain.CompClass{
		ID: n2e(r.ID),
		Ownership: domain.OwnershipData{
			OrganizerID: r.OrganizerID,
		},
		ContestID:   r.ContestID,
		Name:        r.Name,
		Description: n2e(r.Description),
		Color:       domain.ColorRGB(n2e(r.Color)),
		TimeBegin:   r.TimeBegin,
		TimeEnd:     r.TimeEnd,
	}
}

func (d *Database) GetCompClass(ctx context.Context, tx domain.Transaction, compClassID domain.ResourceID) (domain.CompClass, error) {
	transaction, ok := tx.(*transaction)
	if !ok {
		return domain.CompClass{}, ErrIncompatibleTransaction
	}

	var record compClassRecord
	err := transaction.db.WithContext(ctx).Raw(`SELECT * FROM comp_class WHERE id = ?`, compClassID).Scan(&record).Error

	return record.toDomain(), err
}
