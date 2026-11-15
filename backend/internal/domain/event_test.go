package domain_test

import (
	"slices"
	"testing"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestNewEventFilter(t *testing.T) {
	competitionID := testutils.RandomCompetitionID()
	competitorID := testutils.RandomResourceID[domain.CompetitorID]()

	filter := domain.NewEventFilter(competitionID, competitorID, "A", "B", "C")

	assert.Equal(t, competitionID, filter.CompetitionID)
	assert.Equal(t, competitorID, filter.CompetitorID)

	assert.Equal(t, map[string]struct{}{
		"A": {},
		"B": {},
		"C": {}}, filter.EventTypes)
}

func TestMatchFilter(t *testing.T) {
	competitionID := testutils.RandomCompetitionID()
	t.Run("CompetitionMatchWildcard", func(t *testing.T) {
		filter := domain.NewEventFilter(domain.CompetitionID{}, 0)

		match := filter.Match(testutils.RandomCompetitionID(), 0, "A")

		assert.True(t, match)
	})

	t.Run("CompetitionMatch", func(t *testing.T) {
		filter := domain.NewEventFilter(competitionID, 0)

		match := filter.Match(competitionID, 0, "A")

		assert.True(t, match)
	})

	t.Run("CompetitionNoMatch", func(t *testing.T) {
		filter := domain.NewEventFilter(competitionID, 0)

		match := filter.Match(testutils.RandomCompetitionID(), 0, "A")

		assert.False(t, match)
	})

	t.Run("CompetitorMatchWildcard", func(t *testing.T) {
		filter := domain.NewEventFilter(domain.CompetitionID{}, 0)

		match := filter.Match(testutils.RandomCompetitionID(), testutils.RandomResourceID[domain.CompetitorID](), "A")

		assert.True(t, match)
	})

	t.Run("CompetitorMatch", func(t *testing.T) {
		filter := domain.NewEventFilter(domain.CompetitionID{}, 1337)

		match := filter.Match(testutils.RandomCompetitionID(), domain.CompetitorID(1337), "A")

		assert.True(t, match)
	})

	t.Run("CompetitorNoMatch", func(t *testing.T) {
		filter := domain.NewEventFilter(domain.CompetitionID{}, 1337)

		match := filter.Match(testutils.RandomCompetitionID(), domain.CompetitorID(42), "A")

		assert.False(t, match)
	})

	t.Run("EventTypeMatch", func(t *testing.T) {
		filter := domain.NewEventFilter(domain.CompetitionID{}, 0, "A", "B", "C")

		for eventType := range slices.Values([]string{"A", "B", "C"}) {
			match := filter.Match(testutils.RandomCompetitionID(), testutils.RandomResourceID[domain.CompetitorID](), eventType)

			assert.True(t, match)
		}
	})

	t.Run("EventTypeNoMatch", func(t *testing.T) {
		filter := domain.NewEventFilter(domain.CompetitionID{}, 0, "A", "B", "C")

		match := filter.Match(testutils.RandomCompetitionID(), testutils.RandomResourceID[domain.CompetitorID](), "X")

		assert.False(t, match)
	})
}
