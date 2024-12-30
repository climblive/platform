package repository

import (
	"context"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

type organizerRecord struct {
	ID       *int `gorm:"primaryKey;autoIncrement"`
	Name     string
	Homepage *string
}

func (organizerRecord) TableName() string {
	return "organizer"
}

func (r organizerRecord) fromDomain(organizer domain.Organizer) organizerRecord {
	return organizerRecord{
		ID:       e2n(int(organizer.ID)),
		Name:     organizer.Name,
		Homepage: e2n(organizer.Homepage),
	}
}

func (r *organizerRecord) toDomain() domain.Organizer {
	return domain.Organizer{
		ID: domain.OrganizerID(n2e(r.ID)),
		Ownership: domain.OwnershipData{
			OrganizerID: domain.OrganizerID(n2e(r.ID)),
		},
		Name:     r.Name,
		Homepage: n2e(r.Homepage),
	}
}

func (d *Database) StoreOrganizer(ctx context.Context, tx domain.Transaction, organizer domain.Organizer) (domain.Organizer, error) {
	var err error
	var record organizerRecord = organizerRecord{}.fromDomain(organizer)

	err = d.tx(tx).WithContext(ctx).Save(&record).Error
	if err != nil {
		return domain.Organizer{}, errors.Wrap(err, 0)
	}

	return record.toDomain(), nil
}
