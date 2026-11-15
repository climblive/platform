package domain

import (
	"database/sql/driver"
	"uuid"

	"github.com/go-errors/errors"
)

func (id CompetitionID) String() string {
	return uuid.UUID(id).String()
}

func (id CompetitionID) MarshalText() ([]byte, error) {
	return uuid.UUID(id).MarshalText()
}

func (id *CompetitionID) UnmarshalText(text []byte) error {
	parsed, err := uuid.Parse(string(text))
	if err != nil {
		return errors.Wrap(err, 0)
	}
	*id = CompetitionID(parsed)
	return nil
}

func (id CompetitionID) Value() (driver.Value, error) {
	return id[:], nil
}

func (id *CompetitionID) Scan(value any) error {
	data, ok := value.([]byte)
	if !ok || len(data) != len(id) {
		return errors.Wrap(ErrInvalidData, 0)
	}

	copy(id[:], data)
	return nil
}
