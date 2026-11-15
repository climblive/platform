package domain_test

import (
	"encoding/json/v2"
	"testing"
	"uuid"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompetitionID(t *testing.T) {
	id := domain.CompetitionID(uuid.New())
	encoded, err := json.Marshal(id)
	require.NoError(t, err)
	assert.Equal(t, `"`+id.String()+`"`, string(encoded))
	var decoded domain.CompetitionID
	require.NoError(t, json.Unmarshal(encoded, &decoded))
	assert.Equal(t, id, decoded)

	value, err := id.Value()
	require.NoError(t, err)
	var scanned domain.CompetitionID
	require.NoError(t, scanned.Scan(value))
	assert.Equal(t, id, scanned)

	filter := domain.NewEventFilter(domain.CompetitionID{}, 0)
	encoded, err = json.Marshal(filter)
	require.NoError(t, err)
	assert.JSONEq(t, `{}`, string(encoded))
}

func TestCompetitionIDRejectsInvalidValues(t *testing.T) {
	for _, value := range []any{nil, "not binary", int32(1), []byte{}, make([]byte, 15), make([]byte, 17)} {
		t.Run("database", func(t *testing.T) {
			id := domain.CompetitionID(uuid.New())
			before := id
			require.ErrorIs(t, id.Scan(value), domain.ErrInvalidData)
			assert.Equal(t, before, id)
		})
	}
	for _, value := range []string{`1`, `"1"`, `"invalid"`, `""`, `[]`} {
		t.Run(value, func(t *testing.T) {
			var id domain.CompetitionID
			require.Error(t, json.Unmarshal([]byte(value), &id))
		})
	}
}
