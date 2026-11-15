package domain_test

import (
	"encoding/json/v2"
	"testing"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCategoryIDJSON(t *testing.T) {
	for _, tc := range []struct {
		name  string
		value any
	}{
		{"Contender", domain.Contender{CategoryID: 42}},
		{"ContenderPatch", domain.ContenderPatch{CategoryID: domain.NewPatch(domain.CategoryID(42))}},
		{"ScoreboardEntry", domain.ScoreboardEntry{CategoryID: 42}},
		{"ContenderEnteredEvent", domain.ContenderEnteredEvent{CategoryID: 42}},
		{"ContenderSwitchedCategoryEvent", domain.ContenderSwitchedCategoryEvent{CategoryID: 42}},
		{"ContenderPublicInfoUpdatedEvent", domain.ContenderPublicInfoUpdatedEvent{CategoryID: 42}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			data, err := json.Marshal(tc.value)
			require.NoError(t, err)

			var fields map[string]any
			require.NoError(t, json.Unmarshal(data, &fields))
			assert.Equal(t, float64(42), fields["categoryId"])
			assert.NotContains(t, fields, "compClassId")
		})
	}
}

func TestContenderPatchCategoryIDJSON(t *testing.T) {
	var patch domain.ContenderPatch
	err := json.Unmarshal([]byte(`{"categoryId":42}`), &patch, json.RejectUnknownMembers(true))
	require.NoError(t, err)
	assert.Equal(t, domain.NewPatch(domain.CategoryID(42)), patch.CategoryID)
}
