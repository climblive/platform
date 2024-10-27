package domain_test

import (
	"math/rand"
	"testing"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewEventFilter(t *testing.T) {
	contestID := domain.ContestID(rand.Int())
	contenderID := domain.ContenderID(rand.Int())

	filter := domain.NewEventFilter(contestID, contenderID, "A", "B", "C")

	assert.Equal(t, contestID, filter.ContestID)
	assert.Equal(t, contenderID, filter.ContenderID)

	assert.Equal(t, map[string]struct{}{
		"A": {},
		"B": {},
		"C": {}}, filter.EventTypes)
}
