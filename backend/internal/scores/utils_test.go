package scores_test

import (
	"testing"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/scores"
	"github.com/stretchr/testify/assert"
)

func TestCompareScore(t *testing.T) {
	now := time.Now()

	s1 := domain.Score{
		Timestamp:   now,
		ContenderID: 1,
		Score:       100,
		Placement:   1,
		Finalist:    true,
		RankOrder:   0,
	}
	s2 := domain.Score{
		Timestamp:   now.Add(time.Second),
		ContenderID: 1,
		Score:       100,
		Placement:   1,
		Finalist:    true,
		RankOrder:   0,
	}

	assert.True(t, scores.CompareScore(s1, s2))
}
