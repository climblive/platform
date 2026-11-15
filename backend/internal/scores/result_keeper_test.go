package scores_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/events"
	"github.com/climblive/platform/backend/internal/scores"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"uuid"
)

func TestKeeper(t *testing.T) {
	makeMocks := func(bufferCapacity int) (*repositoryMock, *eventBrokerMock, *events.Subscription) {
		mockedRepo := new(repositoryMock)
		mockedEventBroker := new(eventBrokerMock)

		filters := []domain.EventFilter{domain.NewEventFilter(
			0,
			0,
			"CONTENDER_RESULT_UPDATED",
		)}

		subscription := events.NewSubscription(filters, bufferCapacity)
		subscriptionID := uuid.New()

		mockedEventBroker.On("Subscribe", filters, 0).Return(subscriptionID, subscription)

		mockedEventBroker.On("Unsubscribe", subscriptionID).Return()

		return mockedRepo, mockedEventBroker, subscription
	}

	t.Run("StartAndStop", func(t *testing.T) {
		mockedRepo, mockedEventBroker, _ := makeMocks(0)
		keeper := scores.NewResultKeeper(mockedEventBroker, mockedRepo)

		ctx, cancel := context.WithCancel(context.Background())

		wg := keeper.Run(ctx)
		cancel()

		wg.Wait()

		mockedEventBroker.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("SubscriptionUnexpectedlyClosed", func(t *testing.T) {
		mockedRepo, mockedEventBroker, subscription := makeMocks(1)
		keeper := scores.NewResultKeeper(mockedEventBroker, mockedRepo)

		err := subscription.Post(domain.EventEnvelope{
			Data: domain.ContenderResultUpdatedEvent{},
		})
		require.NoError(t, err)

		err = subscription.Post(domain.EventEnvelope{
			Data: domain.ContenderResultUpdatedEvent{},
		})
		require.Error(t, err)

		wg := keeper.Run(context.Background())

		wg.Wait()

		mockedEventBroker.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("GatherResults", func(t *testing.T) {
		mockedRepo, mockedEventBroker, subscription := makeMocks(0)
		keeper := scores.NewResultKeeper(mockedEventBroker, mockedRepo)

		ctx, cancel := context.WithCancel(context.Background())
		now := time.Now()

		wg := keeper.Run(ctx)

		for k := 1; k <= 5; k++ {
			result := domain.Result{
				Timestamp:   now,
				ContenderID: domain.ContenderID(k),
				Score:       fmt.Sprintf("%dp", k*100),
				Placement:   k,
				Finalist:    true,
				RankOrder:   k - 1,
			}

			err := subscription.Post(domain.EventEnvelope{
				Data: domain.ContenderResultUpdatedEvent(result),
			})

			mockedRepo.On("StoreResult", mock.Anything, nil, result).Return(nil)

			require.NoError(t, err)
		}

		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			for k := 1; k <= 5; k++ {
				result, err := keeper.GetResult(domain.ContenderID(k))

				require.NoError(collect, err)
				assert.Equal(collect, now, result.Timestamp)
				assert.Equal(collect, domain.ContenderID(k), result.ContenderID)
				assert.Equal(collect, fmt.Sprintf("%dp", k*100), result.Score)
				assert.Equal(collect, k, result.Placement)
				assert.True(collect, result.Finalist)
				assert.Equal(collect, k-1, result.RankOrder)
			}
		}, time.Second, 10*time.Millisecond)

		cancel()

		wg.Wait()

		mockedEventBroker.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("PersistResults", func(t *testing.T) {
		mockedRepo, mockedEventBroker, subscription := makeMocks(0)
		keeper := scores.NewResultKeeper(mockedEventBroker, mockedRepo)

		ctx, cancel := context.WithCancel(context.Background())
		now := time.Now()

		wg := keeper.Run(ctx)

		for k := 1; k <= 5; k++ {
			result := domain.Result{
				Timestamp:   now,
				ContenderID: domain.ContenderID(k),
				Score:       fmt.Sprintf("%dp", k*100),
				Placement:   k,
				Finalist:    true,
				RankOrder:   k - 1,
			}

			mockedRepo.On("StoreResult", mock.Anything, nil, result).Return(nil)

			err := subscription.Post(domain.EventEnvelope{
				Data: domain.ContenderResultUpdatedEvent(result),
			})

			require.NoError(t, err)
		}

		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			for k := 1; k <= 5; k++ {
				_, err := keeper.GetResult(domain.ContenderID(k))

				require.NoError(collect, err)
			}
		}, time.Second, 10*time.Millisecond)

		keeper.RequestPersist()

		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			for k := 1; k <= 5; k++ {
				_, err := keeper.GetResult(domain.ContenderID(k))

				require.Error(collect, err, domain.ErrNotFound)
			}
		}, time.Second, 10*time.Millisecond)

		cancel()

		wg.Wait()

		mockedEventBroker.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("PersistResultsBeforeShutdown", func(t *testing.T) {
		mockedRepo, mockedEventBroker, subscription := makeMocks(0)
		keeper := scores.NewResultKeeper(mockedEventBroker, mockedRepo)

		ctx, cancel := context.WithCancel(context.Background())
		now := time.Now()

		wg := keeper.Run(ctx)

		for k := 1; k <= 5; k++ {
			result := domain.Result{
				Timestamp:   now,
				ContenderID: domain.ContenderID(k),
				Score:       fmt.Sprintf("%dp", k*100),
				Placement:   k,
				Finalist:    true,
				RankOrder:   k - 1,
			}

			mockedRepo.On("StoreResult", mock.Anything, nil, result).Return(nil)

			err := subscription.Post(domain.EventEnvelope{
				Data: domain.ContenderResultUpdatedEvent(result),
			})

			require.NoError(t, err)
		}

		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			for k := 1; k <= 5; k++ {
				_, err := keeper.GetResult(domain.ContenderID(k))

				require.NoError(collect, err)
			}
		}, time.Second, 10*time.Millisecond)

		cancel()

		wg.Wait()

		mockedEventBroker.AssertExpectations(t)
		mockedRepo.AssertExpectations(t)
	})

	t.Run("PersistResults_KeepInMemoryOnFailure", func(t *testing.T) {
		mockedRepo, mockedEventBroker, subscription := makeMocks(0)
		keeper := scores.NewResultKeeper(mockedEventBroker, mockedRepo)
		errMock := errors.New("mock error")

		mockedRepo.On("StoreResult", mock.Anything, nil, mock.AnythingOfType("domain.Result")).Return(errMock)

		ctx, cancel := context.WithCancel(context.Background())
		now := time.Now()

		wg := keeper.Run(ctx)

		for k := 1; k <= 5; k++ {
			result := domain.Result{
				Timestamp:   now,
				ContenderID: domain.ContenderID(k),
				Score:       fmt.Sprintf("%dp", k*100),
				Placement:   k,
				Finalist:    true,
				RankOrder:   k - 1,
			}

			err := subscription.Post(domain.EventEnvelope{
				Data: domain.ContenderResultUpdatedEvent(result),
			})

			require.NoError(t, err)
		}

		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			for k := 1; k <= 5; k++ {
				_, err := keeper.GetResult(domain.ContenderID(k))

				require.NoError(collect, err)
			}
		}, time.Second, 10*time.Millisecond)

		keeper.RequestPersist()

		assert.EventuallyWithT(t, withLogf(func(collect *CollectTWithLogf) {
			mockedRepo.AssertExpectations(collect)

			for k := 1; k <= 5; k++ {
				_, err := keeper.GetResult(domain.ContenderID(k))

				require.NoError(collect, err)
			}
		}), time.Second, 10*time.Millisecond)

		cancel()

		wg.Wait()

		mockedEventBroker.AssertExpectations(t)
	})

	t.Run("PersistResults_DropIfContenderNotFound", func(t *testing.T) {
		mockedRepo, mockedEventBroker, subscription := makeMocks(0)
		keeper := scores.NewResultKeeper(mockedEventBroker, mockedRepo)

		ctx, cancel := context.WithCancel(context.Background())
		now := time.Now()

		wg := keeper.Run(ctx)

		for k := 1; k <= 5; k++ {
			result := domain.Result{
				Timestamp:   now,
				ContenderID: domain.ContenderID(k),
				Score:       fmt.Sprintf("%dp", k*100),
				Placement:   k,
				Finalist:    true,
				RankOrder:   k - 1,
			}

			mockedRepo.On("StoreResult", mock.Anything, nil, result).Return(domain.ErrNotFound)

			err := subscription.Post(domain.EventEnvelope{
				Data: domain.ContenderResultUpdatedEvent(result),
			})

			require.NoError(t, err)
		}

		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			for k := 1; k <= 5; k++ {
				_, err := keeper.GetResult(domain.ContenderID(k))

				require.NoError(collect, err)
			}
		}, time.Second, 10*time.Millisecond)

		keeper.RequestPersist()

		assert.EventuallyWithT(t, withLogf(func(collect *CollectTWithLogf) {
			mockedRepo.AssertExpectations(collect)

			for k := 1; k <= 5; k++ {
				_, err := keeper.GetResult(domain.ContenderID(k))

				require.Error(collect, err, domain.ErrNotFound)
			}
		}), time.Second, 10*time.Millisecond)

		cancel()

		wg.Wait()

		mockedEventBroker.AssertExpectations(t)
	})
}

func withLogf(condition func(*CollectTWithLogf)) func(*assert.CollectT) {
	return func(c *assert.CollectT) {
		condition(&CollectTWithLogf{c})
	}
}

type CollectTWithLogf struct {
	*assert.CollectT
}

func (c *CollectTWithLogf) Logf(format string, args ...interface{}) {
}
