package repository

import (
	"context"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
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

func (r contenderRecord) fromDomain(contender domain.Contender) contenderRecord {
	return contenderRecord{
		ID:               e2n(contender.ID),
		OrganizerID:      contender.Ownership.OrganizerID,
		ContestID:        contender.ContestID,
		RegistrationCode: contender.RegistrationCode,
		Name:             e2n(contender.Name),
		Club:             e2n(contender.ClubName),
		ClassID:          e2n(contender.CompClassID),
		Entered:          e2n(contender.Entered),
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
		Entered:             n2e(r.Entered),
		WithdrawnFromFinals: false,
		Disqualified:        r.Disqualified,
		Score:               0,
		Placement:           0,
		ScoreUpdated:        time.Time{},
	}
}

func (d *Database) GetContender(ctx context.Context, tx domain.Transaction, contenderID domain.ResourceID) (domain.Contender, error) {
	transaction, ok := tx.(*transaction)
	if !ok {
		return domain.Contender{}, ErrIncompatibleTransaction
	}

	var record contenderRecord
	err := transaction.db.WithContext(ctx).Raw(`SELECT * FROM contender WHERE id = ?`, contenderID).Scan(&record).Error

	return record.toDomain(), err
}

func (d *Database) GetContenderByCode(ctx context.Context, tx domain.Transaction, registrationCode string) (domain.Contender, error) {
	var record contenderRecord
	err := d.db.WithContext(ctx).Raw(`SELECT * FROM contender WHERE registration_code = ?`, registrationCode).Scan(&record).Error

	return record.toDomain(), err
}

func (d *Database) GetContendersByCompClass(ctx context.Context, tx domain.Transaction, compClassID domain.ResourceID) ([]domain.Contender, error) {
	var records []contenderRecord

	err := d.db.WithContext(ctx).Raw(`SELECT * FROM contender WHERE class_id = ?`, compClassID).Scan(&records).Error
	if err != nil {
		return nil, err
	}

	contenders := make([]domain.Contender, 0)

	for _, record := range records {
		contenders = append(contenders, record.toDomain())
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
		contenders = append(contenders, record.toDomain())
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
			disqualified = @Disqualified`, contenderRecord{}.fromDomain(contender),
	).Scan(&storedRecord).Error
	if err != nil {
		return domain.Contender{}, nil
	}

	return storedRecord.toDomain(), nil
}

func (d *Database) DeleteContender(ctx context.Context, tx domain.Transaction, contenderID domain.ResourceID) error {
	return d.db.WithContext(ctx).Raw(`DELETE FROM contender WHERE id = ?`, contenderID).Error
}
