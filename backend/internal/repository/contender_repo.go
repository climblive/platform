package repository

import (
	"context"
	"database/sql"

	"github.com/climblive/platform/backend/internal/database"
	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
	"github.com/google/uuid"
)

func (d *Database) GetContender(ctx context.Context, tx domain.Transaction, contenderID domain.ContenderID) (domain.Contender, error) {
	record, err := d.WithTx(tx).GetContender(ctx, uuid.UUID(contenderID))
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return domain.Contender{}, errors.Wrap(domain.ErrNotFound, 0)
	case err != nil:
		return domain.Contender{}, errors.Wrap(err, 0)
	}

	contender := contenderToDomain(record)

	return contender, nil
}

func (d *Database) GetContenderByCode(ctx context.Context, tx domain.Transaction, registrationCode string) (domain.Contender, error) {
	record, err := d.WithTx(tx).GetContenderByCode(ctx, registrationCode)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return domain.Contender{}, errors.Wrap(domain.ErrNotFound, 0)
	case err != nil:
		return domain.Contender{}, errors.Wrap(err, 0)
	}

	contender := contenderToDomain(database.GetContenderRow(record))

	return contender, nil
}

func (d *Database) GetContendersByCompClass(ctx context.Context, tx domain.Transaction, compClassID domain.CompClassID) ([]domain.Contender, error) {
	records, err := d.WithTx(tx).GetContendersByCompClass(ctx, uuid.NullUUID{Valid: true, UUID: uuid.UUID(compClassID)})
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	contenders := make([]domain.Contender, 0)

	for _, record := range records {
		contender := contenderToDomain(database.GetContenderRow(record))

		contenders = append(contenders, contender)
	}

	return contenders, nil
}

func (d *Database) GetContendersByContest(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) ([]domain.Contender, error) {
	records, err := d.WithTx(tx).GetContendersByContest(ctx, uuid.UUID(contestID))
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	contenders := make([]domain.Contender, 0)

	for _, record := range records {
		contender := contenderToDomain(database.GetContenderRow(record))

		contenders = append(contenders, contender)
	}

	return contenders, nil
}

func (d *Database) StoreContender(ctx context.Context, tx domain.Transaction, contender domain.Contender) (domain.Contender, error) {
	if uuid.UUID(contender.ID) == uuid.Nil {
		contender.ID = domain.ContenderID(uuid.New())
	}

	params := database.UpsertContenderParams{
		ID:                  uuid.UUID(contender.ID),
		OrganizerID:         uuid.UUID(contender.Ownership.OrganizerID),
		ContestID:           uuid.UUID(contender.ContestID),
		RegistrationCode:    contender.RegistrationCode,
		Name:                makeNullString(contender.Name),
		ClassID:             makeNullUUID(uuid.UUID(contender.CompClassID)),
		Entered:             makeNullTime(contender.Entered),
		Disqualified:        contender.Disqualified,
		WithdrawnFromFinals: contender.WithdrawnFromFinals,
	}

	_, err := d.WithTx(tx).UpsertContender(ctx, params)
	if err != nil {
		return domain.Contender{}, errors.Wrap(err, 0)
	}

	return contender, err
}

func (d *Database) DeleteContender(ctx context.Context, tx domain.Transaction, contenderID domain.ContenderID) error {
	err := d.WithTx(tx).DeleteContender(ctx, uuid.UUID(contenderID))
	if err != nil {
		return errors.Wrap(err, 0)
	}

	return nil
}

func (d *Database) GetNumberOfContenders(ctx context.Context, tx domain.Transaction, contestID domain.ContestID) (int, error) {
	count, err := d.WithTx(tx).CountContenders(ctx, uuid.UUID(contestID))
	if err != nil {
		return 0, errors.Wrap(err, 0)
	}

	return int(count), nil
}
