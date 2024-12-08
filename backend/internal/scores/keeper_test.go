package scores_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/events"
	"github.com/climblive/platform/backend/internal/scores"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestKeeperGatherScores(t *testing.T) {
	makeMocks := func(bufferCapacity int) (*repositoryMock, *eventBrokerMock, *events.Subscription) {
		mockedRepo := new(repositoryMock)
		mockedEventBroker := new(eventBrokerMock)

		subscription := events.NewSubscription(domain.EventFilter{}, bufferCapacity)
		subscriptionID := uuid.New()

		mockedEventBroker.On("Subscribe", domain.NewEventFilter(
			0,
			0,
			"CONTENDER_SCORE_UPDATED",
		), 0).Return(subscriptionID, subscription)

		mockedEventBroker.On("Unsubscribe", subscriptionID).Return()

		return mockedRepo, mockedEventBroker, subscription
	}

	t.Run("StartAndStop", func(t *testing.T) {
		mockedRepo, mockedEventBroker, _ := makeMocks(0)
		keeper := scores.NewScoreKeeper(mockedEventBroker, mockedRepo)

		ctx, cancel := context.WithCancel(context.Background())

		wg := keeper.Run(ctx)
		cancel()

		wg.Wait()

		mockedEventBroker.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("SubscriptionUnexpectedlyClosed", func(t *testing.T) {
		mockedRepo, mockedEventBroker, subscription := makeMocks(1)
		keeper := scores.NewScoreKeeper(mockedEventBroker, mockedRepo)

		err := subscription.Post(domain.EventEnvelope{
			Name: "CONTENDER_SCORE_UPDATED",
			Data: domain.ContenderScoreUpdatedEvent{},
		})
		require.NoError(t, err)

		err = subscription.Post(domain.EventEnvelope{
			Name: "CONTENDER_SCORE_UPDATED",
			Data: domain.ContenderScoreUpdatedEvent{},
		})
		require.Error(t, err)

		wg := keeper.Run(context.Background())

		wg.Wait()

		mockedEventBroker.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("GatherScores", func(t *testing.T) {
		mockedRepo, mockedEventBroker, subscription := makeMocks(0)
		keeper := scores.NewScoreKeeper(mockedEventBroker, mockedRepo)

		ctx, cancel := context.WithCancel(context.Background())
		now := time.Now()

		wg := keeper.Run(ctx)

		for k := 1; k <= 5; k++ {
			err := subscription.Post(domain.EventEnvelope{
				Name: "CONTENDER_SCORE_UPDATED",
				Data: domain.ContenderScoreUpdatedEvent{
					Timestamp:   now,
					ContenderID: domain.ContenderID(k),
					Score:       k * 100,
					Placement:   k,
					Finalist:    true,
					RankOrder:   k - 1,
				},
			})

			require.NoError(t, err)
		}

		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			for k := 1; k <= 5; k++ {
				score, err := keeper.GetScore(domain.ContenderID(k))

				require.NoError(collect, err)
				assert.Equal(collect, now, score.Timestamp)
				assert.Equal(collect, domain.ContenderID(k), score.ContenderID)
				assert.Equal(collect, k*100, score.Score)
				assert.Equal(collect, k, score.Placement)
				assert.True(collect, score.Finalist)
				assert.Equal(collect, k-1, score.RankOrder)
			}
		}, time.Second, 10*time.Millisecond)

		cancel()

		wg.Wait()

		mockedEventBroker.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("PersistScores", func(t *testing.T) {
		mockedRepo, mockedEventBroker, subscription := makeMocks(0)
		keeper := scores.NewScoreKeeper(mockedEventBroker, mockedRepo)

		ctx, cancel := context.WithCancel(context.Background())
		now := time.Now()

		wg := keeper.Run(ctx)

		for k := 1; k <= 5; k++ {
			score := domain.Score{
				Timestamp:   now,
				ContenderID: domain.ContenderID(k),
				Score:       k * 100,
				Placement:   k,
				Finalist:    true,
				RankOrder:   k - 1,
			}

			mockedRepo.On("StoreScore", mock.Anything, mock.Anything, score).Return(score, nil)

			err := subscription.Post(domain.EventEnvelope{
				Name: "CONTENDER_SCORE_UPDATED",
				Data: domain.ContenderScoreUpdatedEvent(score),
			})

			require.NoError(t, err)
		}

		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			for k := 1; k <= 5; k++ {
				_, err := keeper.GetScore(domain.ContenderID(k))

				require.NoError(collect, err)
			}
		}, time.Second, 10*time.Millisecond)

		keeper.RequestPersist()

		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			for k := 1; k <= 5; k++ {
				_, err := keeper.GetScore(domain.ContenderID(k))

				require.Error(collect, err, domain.ErrNotFound)
			}
		}, time.Second, 10*time.Millisecond)

		cancel()

		wg.Wait()

		mockedEventBroker.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("PersistScoresBeforeShutdown", func(t *testing.T) {
		mockedRepo, mockedEventBroker, subscription := makeMocks(0)
		keeper := scores.NewScoreKeeper(mockedEventBroker, mockedRepo)

		ctx, cancel := context.WithCancel(context.Background())
		now := time.Now()

		wg := keeper.Run(ctx)

		for k := 1; k <= 5; k++ {
			score := domain.Score{
				Timestamp:   now,
				ContenderID: domain.ContenderID(k),
				Score:       k * 100,
				Placement:   k,
				Finalist:    true,
				RankOrder:   k - 1,
			}

			mockedRepo.On("StoreScore", mock.Anything, mock.Anything, score).Return(score, nil)

			err := subscription.Post(domain.EventEnvelope{
				Name: "CONTENDER_SCORE_UPDATED",
				Data: domain.ContenderScoreUpdatedEvent(score),
			})

			require.NoError(t, err)
		}

		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			for k := 1; k <= 5; k++ {
				_, err := keeper.GetScore(domain.ContenderID(k))

				require.NoError(collect, err)
			}
		}, time.Second, 10*time.Millisecond)

		cancel()

		wg.Wait()

		mockedEventBroker.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("PersistScoresFailure", func(t *testing.T) {
		mockedRepo, mockedEventBroker, subscription := makeMocks(0)
		keeper := scores.NewScoreKeeper(mockedEventBroker, mockedRepo)
		var errMock error = errors.New("mock error")

		ctx, cancel := context.WithCancel(context.Background())
		now := time.Now()

		wg := keeper.Run(ctx)

		for k := 1; k <= 5; k++ {
			score := domain.Score{
				Timestamp:   now,
				ContenderID: domain.ContenderID(k),
				Score:       k * 100,
				Placement:   k,
				Finalist:    true,
				RankOrder:   k - 1,
			}

			mockedRepo.On("StoreScore", mock.Anything, mock.Anything, score).Return(domain.Score{}, errMock)

			err := subscription.Post(domain.EventEnvelope{
				Name: "CONTENDER_SCORE_UPDATED",
				Data: domain.ContenderScoreUpdatedEvent(score),
			})

			require.NoError(t, err)
		}

		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			for k := 1; k <= 5; k++ {
				_, err := keeper.GetScore(domain.ContenderID(k))

				require.NoError(collect, err)
			}
		}, time.Second, 10*time.Millisecond)

		keeper.RequestPersist()

		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			for k := 1; k <= 5; k++ {
				_, err := keeper.GetScore(domain.ContenderID(k))

				require.NoError(collect, err)
			}
		}, time.Second, 10*time.Millisecond)

		cancel()

		wg.Wait()

		mockedEventBroker.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})
}
