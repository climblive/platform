package repository

import (
	"context"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
)

type contenderRecord struct {
	ID               *int       `gorm:"primaryKey,autoIncrement:true"`
	OrganizerID      int        `gorm:"column:organizer_id"`
	ContestID        int        `gorm:"column:contest_id"`
	RegistrationCode string     `gorm:"column:registration_code"`
	Name             *string    `gorm:"column:name"`
	Club             *string    `gorm:"column:club"`
	ClassID          *int       `gorm:"column:class_id"`
	Entered          *time.Time `gorm:"column:entered"`
	Disqualified     bool       `gorm:"column:disqualified"`
}

func (r contenderRecord) ToRecord(contender domain.Contender) contenderRecord {
	return contenderRecord{
		ID:               emptyAsNil(contender.ID),
		OrganizerID:      contender.Ownership.OrganizerID,
		ContestID:        contender.ContestID,
		RegistrationCode: contender.RegistrationCode,
		Name:             emptyAsNil(contender.Name),
		Club:             emptyAsNil(contender.ClubName),
		ClassID:          emptyAsNil(contender.CompClassID),
		Entered:          emptyAsNil(contender.Entered),
		Disqualified:     contender.Disqualified,
	}
}

func (r *contenderRecord) ToDomain() domain.Contender {
	return domain.Contender{
		ID: nilAsEmpty(r.ID),
		Ownership: domain.OwnershipData{
			OrganizerID: r.OrganizerID,
			ContenderID: r.ID,
		},
		ContestID:           r.ContestID,
		CompClassID:         nilAsEmpty(r.ClassID),
		RegistrationCode:    r.RegistrationCode,
		Name:                nilAsEmpty(r.Name),
		PublicName:          nilAsEmpty(r.Name),
		ClubName:            nilAsEmpty(r.Club),
		Entered:             nilAsEmpty(r.Entered),
		WithdrawnFromFinals: false,
		Disqualified:        r.Disqualified,
		Score:               0,
		Placement:           0,
		ScoreUpdated:        time.Time{},
	}
}

func (d *Database) GetContender(ctx context.Context, tx domain.Transaction, contenderID domain.ResourceID) (domain.Contender, error) {
	var record contenderRecord
	err := d.db.WithContext(ctx).Raw(`SELECT * FROM contender WHERE id = ?`, contenderID).Scan(&record).Error

	return record.ToDomain(), err
}

func (d *Database) GetContenderByCode(ctx context.Context, tx domain.Transaction, registrationCode string) (domain.Contender, error) {
	var record contenderRecord
	err := d.db.WithContext(ctx).Raw(`SELECT * FROM contender WHERE registration_code = ?`, registrationCode).Scan(&record).Error

	return record.ToDomain(), err
}

func (d *Database) GetContendersByCompClass(ctx context.Context, tx domain.Transaction, compClassID domain.ResourceID) ([]domain.Contender, error) {
	var records []contenderRecord

	err := d.db.WithContext(ctx).Raw(`SELECT * FROM contender WHERE class_id = ?`, compClassID).Scan(&records).Error
	if err != nil {
		return nil, err
	}

	contenders := make([]domain.Contender, 0)

	for _, record := range records {
		contenders = append(contenders, record.ToDomain())
	}

	return contenders, nil
}

func (d *Database) GetContendersByContest(ctx context.Context, tx domain.Transaction, contestID domain.ResourceID) ([]domain.Contender, error) {
	var records []contenderRecord

	err := d.db.WithContext(ctx).Raw(`SELECT * FROM contender WHERE contest_id = ?`, contestID).Scan(&records).Error
	if err != nil {
		return nil, err
	}

	contenders := make([]domain.Contender, 0)

	for _, record := range records {
		contenders = append(contenders, record.ToDomain())
	}

	return contenders, nil
}

func (d *Database) StoreContender(ctx context.Context, tx domain.Transaction, contender domain.Contender) (domain.Contender, error) {
	var storedRecord contenderRecord

	err := d.db.WithContext(ctx).Raw(`
		INSERT INTO
			contender (id, organizer_id, contest_id, registration_code, name, club, class_id, entered, disqualified)
			VALUES (@ID, @OrganizerID, @ContestID, @RegistrationCode, @Name, @Club, @ClassID, @Entered, @Disqualified)
		ON DUPLICATE KEY UPDATE
			name = @Name,
			club = @Club,
			class_id = @ClassID,
			entered = @Entered,
			disqualified = @Disqualified`, contenderRecord{}.ToRecord(contender),
	).Scan(&storedRecord).Error
	if err != nil {
		return domain.Contender{}, nil
	}

	return storedRecord.ToDomain(), nil
}

func (d *Database) DeleteContender(ctx context.Context, tx domain.Transaction, contenderID domain.ResourceID) error {
	return d.db.WithContext(ctx).Raw(`DELETE FROM contender WHERE id = ?`, contenderID).Error
}
