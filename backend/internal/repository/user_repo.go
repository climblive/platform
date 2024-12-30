package repository

import (
	"context"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

type userRecord struct {
	ID         *int `gorm:"primaryKey;autoIncrement"`
	Name       string
	Username   string
	Admin      bool
	Organizers []organizerRecord `gorm:"many2many:user_organizer;foreignKey:id;joinForeignKey:user_id;References:id;joinReferences:organizer_id;"`
}

func (userRecord) TableName() string {
	return "user"
}

func (r userRecord) fromDomain(user domain.User) userRecord {
	record := userRecord{
		ID:       e2n(int(user.ID)),
		Name:     user.Name,
		Username: user.Username,
		Admin:    user.Admin,
	}

	for _, organizerID := range user.Organizers {
		record.Organizers = append(record.Organizers, organizerRecord{
			ID: e2n(int(organizerID)),
		})
	}

	return record
}

func (r *userRecord) toDomain() domain.User {
	entity := domain.User{
		ID:       domain.UserID(n2e(r.ID)),
		Name:     r.Name,
		Username: r.Username,
		Admin:    r.Admin,
	}

	for _, organizer := range r.Organizers {
		entity.Organizers = append(entity.Organizers, domain.OrganizerID(n2e(organizer.ID)))
	}

	return entity
}

func (d *Database) StoreUser(ctx context.Context, tx domain.Transaction, user domain.User) (domain.User, error) {
	var err error
	var record userRecord = userRecord{}.fromDomain(user)

	err = d.tx(tx).WithContext(ctx).Save(&record).Error
	if err != nil {
		return domain.User{}, errors.Wrap(err, 0)
	}

	return record.toDomain(), nil
}

func (d *Database) GetUserByUsername(ctx context.Context, tx domain.Transaction, username string) (domain.User, error) {
	var record userRecord

	err := d.tx(tx).WithContext(ctx).
		Debug().
		Model(&userRecord{}).
		Preload("Organizers").
		Select("*").
		Where("username = ?", username).
		Scan(&record).Error
	if err != nil {
		return domain.User{}, errors.Wrap(err, 0)
	}

	if record.ID == nil {
		return domain.User{}, errors.Wrap(domain.ErrNotFound, 0)
	}

	return record.toDomain(), nil
}
