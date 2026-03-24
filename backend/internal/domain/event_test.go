package domain_test

import (
	"testing"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestNewEventFilter(t *testing.T) {
	contestID := testutils.RandomResourceID[domain.ContestID]()
	contenderID := testutils.RandomResourceID[domain.ContenderID]()

	filter := domain.NewEventFilter(contestID, contenderID, "A", "B", "C")

	assert.Equal(t, contestID, filter.ContestID)
	assert.Equal(t, contenderID, filter.ContenderID)

	assert.Equal(t, map[string]struct{}{
		"A": {},
		"B": {},
		"C": {}}, filter.EventTypes)
}
