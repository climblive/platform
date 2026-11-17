package scores_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/scores"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompareCompetitorAttemptPriority(t *testing.T) {
	tests := []struct {
		name          string
		first         scores.Score
		second        scores.Score
		zonesFirst    int
		attemptsFirst int
	}{
		{
			name:          "tops before zones and attempts",
			first:         scores.Score{Tops: 2, AttemptsTops: 10, Zone2s: 2, Zone1s: 2},
			second:        scores.Score{Tops: 1, AttemptsTops: 1, Zone2s: 5, Zone1s: 5},
			zonesFirst:    -1,
			attemptsFirst: -1,
		},
		{
			name:          "zone 2s versus top attempts",
			first:         scores.Score{Tops: 1, AttemptsTops: 3, Zone2s: 2},
			second:        scores.Score{Tops: 1, AttemptsTops: 1, Zone2s: 1},
			zonesFirst:    -1,
			attemptsFirst: 1,
		},
		{
			name:          "zone 1s versus top attempts",
			first:         scores.Score{Tops: 1, AttemptsTops: 3, Zone1s: 2},
			second:        scores.Score{Tops: 1, AttemptsTops: 1, Zone1s: 1},
			zonesFirst:    -1,
			attemptsFirst: 1,
		},
		{
			name:          "zone 1s versus zone 2 attempts",
			first:         scores.Score{Zone2s: 1, AttemptsZone2s: 3, Zone1s: 2},
			second:        scores.Score{Zone2s: 1, AttemptsZone2s: 1, Zone1s: 1},
			zonesFirst:    -1,
			attemptsFirst: 1,
		},
		{
			name:          "zone 2s before zone 1s",
			first:         scores.Score{Zone2s: 2, Zone1s: 2},
			second:        scores.Score{Zone2s: 1, Zone1s: 3},
			zonesFirst:    -1,
			attemptsFirst: -1,
		},
		{
			name:          "top attempts when counts tie",
			first:         scores.Score{Tops: 1, AttemptsTops: 1, Zone2s: 2, Zone1s: 3, AttemptsZone2s: 9},
			second:        scores.Score{Tops: 1, AttemptsTops: 2, Zone2s: 2, Zone1s: 3, AttemptsZone2s: 2},
			zonesFirst:    -1,
			attemptsFirst: -1,
		},
		{
			name:          "zone 2 attempts before zone 1 attempts",
			first:         scores.Score{Zone2s: 1, Zone1s: 1, AttemptsZone2s: 1, AttemptsZone1s: 2},
			second:        scores.Score{Zone2s: 1, Zone1s: 1, AttemptsZone2s: 2, AttemptsZone1s: 1},
			zonesFirst:    -1,
			attemptsFirst: -1,
		},
		{
			name:          "zone 1 attempts",
			first:         scores.Score{Zone1s: 1, AttemptsZone1s: 1},
			second:        scores.Score{Zone1s: 1, AttemptsZone1s: 2},
			zonesFirst:    -1,
			attemptsFirst: -1,
		},
		{
			name:   "equal scores",
			first:  scores.Score{Tops: 1, AttemptsTops: 1},
			second: scores.Score{Tops: 1, AttemptsTops: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			first := scores.Competitor{ID: 1, Score: tt.first}
			second := scores.Competitor{ID: 1, Score: tt.second}
			assert.Equal(t, tt.zonesFirst, first.Compare(second, false))
			assert.Equal(t, -tt.zonesFirst, second.Compare(first, false))
			assert.Equal(t, tt.attemptsFirst, first.Compare(second, true))
			assert.Equal(t, -tt.attemptsFirst, second.Compare(first, true))
		})
	}
}

func TestBasicRankerAttemptPriority(t *testing.T) {
	competitors := []scores.Competitor{
		{ID: 1, Score: scores.Score{Tops: 1, AttemptsTops: 1, Zone1s: 1, AttemptsZone1s: 1}},
		{ID: 2, Score: scores.Score{Tops: 1, AttemptsTops: 2, Zone1s: 2, AttemptsZone1s: 2}},
		{ID: 3, Score: scores.Score{Tops: 1, AttemptsTops: 2, Zone1s: 2, AttemptsZone1s: 2}},
	}
	for _, favorTopAttempts := range []bool{false, true} {
		t.Run(fmt.Sprintf("favorTopAttempts=%t", favorTopAttempts), func(t *testing.T) {
			ranker := scores.NewBasicRanker(1, false, true, false, favorTopAttempts)
			results := ranker.RankCompetitors(slices.Values(competitors))
			expected := []string{"i:2 p:1 r:0 f:🏆", "i:3 p:1 r:1 f:🏆", "i:1 p:3 r:2 f:-"}
			if favorTopAttempts {
				expected = []string{"i:1 p:1 r:0 f:🏆", "i:2 p:2 r:1 f:-", "i:3 p:2 r:2 f:-"}
			}
			assert.Equal(t, expected, prettifyAll(results))
		})
	}
}

func TestAttemptPriorityDoesNotAffectPointsRanking(t *testing.T) {
	competitors := []scores.Competitor{
		{ID: 1, Score: scores.Score{Points: 250, Tops: 1, AttemptsTops: 1, Zone1s: 1}},
		{ID: 2, Score: scores.Score{Points: 250, Tops: 1, AttemptsTops: 2, Zone1s: 2}},
	}
	for _, favorTopAttempts := range []bool{false, true} {
		ranker := scores.NewBasicRanker(1, true, true, false, favorTopAttempts)
		assert.Equal(t, []string{"i:1 p:1 r:0 f:🏆", "i:2 p:1 r:1 f:🏆"},
			prettifyAll(ranker.RankCompetitors(slices.Values(competitors))))
	}
}

func TestAttemptPriorityRulesUpdate(t *testing.T) {
	tests := []struct {
		name         string
		zone1Enabled bool
		zone2Enabled bool
	}{
		{name: "no zones"},
		{name: "zone 1", zone1Enabled: true},
		{name: "zone 2", zone2Enabled: true},
		{name: "two zones", zone1Enabled: true, zone2Enabled: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := scores.NewMemoryStore()
			for id := domain.ProblemID(1); id <= 2; id++ {
				store.SaveProblem(scores.Problem{ID: id, Zone1Enabled: tt.zone1Enabled, Zone2Enabled: tt.zone2Enabled})
			}
			for id := domain.CompetitorID(1); id <= 2; id++ {
				store.SaveCompetitor(scores.Competitor{ID: id, CategoryID: 1})
				store.SaveTick(id, scores.Tick{ID: domain.TickID(id), ProblemID: 1,
					Top: true, AttemptsTop: int(id), Zone1: true, AttemptsZone1: 1, Zone2: true, AttemptsZone2: 1})
			}
			store.SaveTick(2, scores.Tick{ID: 3, ProblemID: 2, Zone1: true, AttemptsZone1: 1, Zone2: true, AttemptsZone2: 1})
			engine := scores.NewDefaultScoreEngine(store)
			runner := scores.NewEffectRunner(engine)
			for _, favorTopAttempts := range []bool{false, true, false} {
				runner.RunEffects(engine.HandleRulesUpdated(domain.RulesUpdatedEvent{
					CompetitionID: 1, Finalists: 1, FavorTopAttempts: favorTopAttempts,
				}))
				expectedWinner := domain.CompetitorID(2)
				if favorTopAttempts || (!tt.zone1Enabled && !tt.zone2Enabled) {
					expectedWinner = 1
				}
				results := engine.GetDirtyResults()
				if !tt.zone1Enabled && !tt.zone2Enabled && len(results) == 0 {
					continue
				}
				require.Len(t, results, 2)
				for _, result := range results {
					assert.Equal(t, result.CompetitorID == expectedWinner, result.Finalist)
					if result.CompetitorID == expectedWinner {
						assert.Equal(t, 1, result.Placement)
					} else {
						assert.Equal(t, 2, result.Placement)
					}
				}
			}
		})
	}
}
