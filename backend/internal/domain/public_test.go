package domain_test

import (
	"encoding/json/v2"
	"testing"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompetitionIDJSON(t *testing.T) {
	for _, tc := range []struct {
		name  string
		value any
	}{
		{"EventFilter", domain.EventFilter{CompetitionID: 42}},
		{"CompClass", domain.CompClass{CompetitionID: 42}},
		{"Contender", domain.Contender{CompetitionID: 42}},
		{"Problem", domain.Problem{CompetitionID: 42}},
		{"Raffle", domain.Raffle{CompetitionID: 42}},
		{"ScoreEngineDescriptor", domain.ScoreEngineDescriptor{CompetitionID: 42}},
		{"RulesUpdatedEvent", domain.RulesUpdatedEvent{CompetitionID: 42}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			data, err := json.Marshal(tc.value)
			require.NoError(t, err)

			var fields map[string]any
			require.NoError(t, json.Unmarshal(data, &fields))
			assert.Equal(t, float64(42), fields["competitionId"])
			assert.NotContains(t, fields, "contestId")
		})
	}

	t.Run("DecodeEventFilter", func(t *testing.T) {
		var filter domain.EventFilter
		err := json.Unmarshal([]byte(`{"competitionId":42}`), &filter, json.RejectUnknownMembers(true))
		require.NoError(t, err)
		assert.Equal(t, domain.CompetitionID(42), filter.CompetitionID)
	})
}
