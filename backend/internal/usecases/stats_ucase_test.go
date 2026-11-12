package usecases_test

import (
	"context"
	"testing"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/usecases"
	"github.com/stretchr/testify/assert"
)

func TestGetGlobalStats(t *testing.T) {
	for _, tc := range []struct {
		name  string
		stats domain.GlobalStats
		err   error
	}{
		{name: "Counts", stats: domain.GlobalStats{Contests: 12, RegisteredContenders: 345, Ticks: 6789}},
		{name: "Empty"},
		{name: "RepositoryError", err: errMock},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			repo := new(repositoryMock)
			repo.On("GetGlobalStats", ctx).Return(tc.stats, tc.err).Once()
			uc := usecases.StatsUseCase{Repo: repo}

			stats, err := uc.GetGlobalStats(ctx)

			assert.Equal(t, tc.stats, stats)
			assert.ErrorIs(t, err, tc.err)
			repo.AssertExpectations(t)
		})
	}
}
