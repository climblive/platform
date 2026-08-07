package domain_test

import (
	"slices"
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

func TestMatchFilter(t *testing.T) {
	t.Run("ContestMatchWildcard", func(t *testing.T) {
		filter := domain.NewEventFilter(0, 0)

		match := filter.Match(testutils.RandomResourceID[domain.ContestID](), 0, "A")

		assert.True(t, match)
	})

	t.Run("ContestMatch", func(t *testing.T) {
		filter := domain.NewEventFilter(1337, 0)

		match := filter.Match(domain.ContestID(1337), 0, "A")

		assert.True(t, match)
	})

	t.Run("ContestNoMatch", func(t *testing.T) {
		filter := domain.NewEventFilter(1337, 0)

		match := filter.Match(domain.ContestID(42), 0, "A")

		assert.False(t, match)
	})

	t.Run("ContenderMatchWildcard", func(t *testing.T) {
		filter := domain.NewEventFilter(0, 0)

		match := filter.Match(testutils.RandomResourceID[domain.ContestID](), testutils.RandomResourceID[domain.ContenderID](), "A")

		assert.True(t, match)
	})

	t.Run("ContenderMatch", func(t *testing.T) {
		filter := domain.NewEventFilter(0, 1337)

		match := filter.Match(testutils.RandomResourceID[domain.ContestID](), domain.ContenderID(1337), "A")

		assert.True(t, match)
	})

	t.Run("ContenderNoMatch", func(t *testing.T) {
		filter := domain.NewEventFilter(0, 1337)

		match := filter.Match(testutils.RandomResourceID[domain.ContestID](), domain.ContenderID(42), "A")

		assert.False(t, match)
	})

	t.Run("EventTypeMatch", func(t *testing.T) {
		filter := domain.NewEventFilter(0, 0, "A", "B", "C")

		for eventType := range slices.Values([]string{"A", "B", "C"}) {
			match := filter.Match(testutils.RandomResourceID[domain.ContestID](), testutils.RandomResourceID[domain.ContenderID](), eventType)

			assert.True(t, match)
		}
	})

	t.Run("EventTypeNoMatch", func(t *testing.T) {
		filter := domain.NewEventFilter(0, 0, "A", "B", "C")

		match := filter.Match(testutils.RandomResourceID[domain.ContestID](), testutils.RandomResourceID[domain.ContenderID](), "X")

		assert.False(t, match)
	})
}
