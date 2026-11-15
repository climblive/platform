package domain_test

import (
	"encoding/json/v2"
	"testing"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResultJSON(t *testing.T) {
	result := &domain.Result{
		Timestamp:   time.Date(2026, time.November, 15, 12, 0, 0, 0, time.UTC),
		ContenderID: 42,
		Score:       "123.4p",
		Placement:   2,
		Finalist:    true,
		RankOrder:   1,
	}

	tests := []struct {
		name      string
		value     any
		hasResult bool
	}{
		{name: "ContenderWithResult", value: domain.Contender{Result: result}, hasResult: true},
		{name: "ContenderWithoutResult", value: domain.Contender{}},
		{name: "ScoreboardEntryWithResult", value: domain.ScoreboardEntry{Result: result}, hasResult: true},
		{name: "ScoreboardEntryWithoutResult", value: domain.ScoreboardEntry{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := json.Marshal(tt.value)
			require.NoError(t, err)

			var fields map[string]any
			require.NoError(t, json.Unmarshal(encoded, &fields))
			assert.NotContains(t, fields, "score")

			if !tt.hasResult {
				assert.NotContains(t, fields, "result")
				return
			}

			require.Contains(t, fields, "result")
			encodedResult, err := json.Marshal(fields["result"])
			require.NoError(t, err)
			assert.JSONEq(t, `{"timestamp":"2026-11-15T12:00:00Z","contenderId":42,"score":"123.4p","placement":2,"finalist":true,"rankOrder":1}`, string(encodedResult))
		})
	}
}
