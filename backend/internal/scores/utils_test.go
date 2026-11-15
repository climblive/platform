package scores_test

import (
	"testing"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/scores"
	"github.com/stretchr/testify/assert"
)

func TestCompareResult(t *testing.T) {
	now := time.Now()

	r1 := domain.Result{
		Timestamp:    now,
		CompetitorID: 1,
		Score:        "100p",
		Placement:    1,
		Finalist:     true,
		RankOrder:    0,
	}
	r2 := domain.Result{
		Timestamp:    now.Add(time.Second),
		CompetitorID: 1,
		Score:        "100p",
		Placement:    1,
		Finalist:     true,
		RankOrder:    0,
	}

	assert.True(t, scores.CompareResult(r1, r2))
}
