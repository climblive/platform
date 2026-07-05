package scores_test

import (
	"context"
	"testing"
	"testing/synctest"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/events"
	"github.com/climblive/platform/backend/internal/scores"
	"github.com/climblive/platform/backend/internal/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPointValueKeeper(t *testing.T) {
	makeMocks := func(bufferCapacity int) (*eventBrokerMock, *events.Subscription) {
		mockedEventBroker := new(eventBrokerMock)

		subscription := events.NewSubscription(domain.EventFilter{}, bufferCapacity)
		subscriptionID := uuid.New()

		mockedEventBroker.On("Subscribe", domain.NewEventFilter(
			0,
			0,
			"POINT_VALUE_UPDATED",
		), 0).Return(subscriptionID, subscription)

		mockedEventBroker.On("Unsubscribe", subscriptionID).Return()

		return mockedEventBroker, subscription
	}

	t.Run("StartAndStop", func(t *testing.T) {
		mockedEventBroker, _ := makeMocks(0)
		keeper := scores.NewPointValueKeeper(mockedEventBroker, time.Hour)

		ctx, cancel := context.WithCancel(context.Background())

		wg := keeper.Run(ctx)
		cancel()

		wg.Wait()

		mockedEventBroker.AssertExpectations(t)
	})

	t.Run("SubscriptionUnexpectedlyClosed", func(t *testing.T) {
		mockedEventBroker, subscription := makeMocks(1)
		keeper := scores.NewPointValueKeeper(mockedEventBroker, time.Hour)

		err := subscription.Post(domain.EventEnvelope{
			Data: domain.PointValueUpdatedEvent{},
		})
		require.NoError(t, err)

		err = subscription.Post(domain.EventEnvelope{
			Data: domain.PointValueUpdatedEvent{},
		})
		require.Error(t, err)

		wg := keeper.Run(context.Background())

		wg.Wait()

		mockedEventBroker.AssertExpectations(t)
	})

	t.Run("GatherScores", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			mockedEventBroker, subscription := makeMocks(0)
			keeper := scores.NewPointValueKeeper(mockedEventBroker, time.Hour)

			ctx, cancel := context.WithCancel(context.Background())

			wg := keeper.Run(ctx)

			fakedContenderIDs := []domain.ContenderID{
				testutils.RandomResourceID[domain.ContenderID](),
				testutils.RandomResourceID[domain.ContenderID](),
				testutils.RandomResourceID[domain.ContenderID](),
			}

			for _, contenderID := range fakedContenderIDs {
				for k := 1; k <= 3; k++ {
					pointValue := domain.PointValue{
						ContenderID: contenderID,
						ProblemID:   domain.ProblemID(1000 + k),
						Current:     1000 + k,
						Zone1:       10 + k,
						Zone2:       20 + k,
						Top:         1000 + k,
						FlashBonus:  100 + k,
					}

					err := subscription.Post(domain.EventEnvelope{
						Data: domain.PointValueUpdatedEvent(pointValue),
					})

					require.NoError(t, err)
				}
			}

			synctest.Wait()

			for _, contenderID := range fakedContenderIDs {
				pointValues := keeper.GetPointValues(contenderID)

				assert.ElementsMatch(t, []domain.PointValue{
					{
						ContenderID: contenderID,
						ProblemID:   1001,
						Current:     1001,
						Zone1:       11,
						Zone2:       21,
						Top:         1001,
						FlashBonus:  101,
					},
					{
						ContenderID: contenderID,
						ProblemID:   1002,
						Current:     1002,
						Zone1:       12,
						Zone2:       22,
						Top:         1002,
						FlashBonus:  102,
					},
					{
						ContenderID: contenderID,
						ProblemID:   1003,
						Current:     1003,
						Zone1:       13,
						Zone2:       23,
						Top:         1003,
						FlashBonus:  103,
					},
				}, pointValues)
			}

			cancel()

			wg.Wait()

			mockedEventBroker.AssertExpectations(t)
		})
	})

	t.Run("ExpungeExpiredScores", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			mockedEventBroker, subscription := makeMocks(0)
			keeper := scores.NewPointValueKeeper(mockedEventBroker, time.Hour)

			ctx, cancel := context.WithCancel(context.Background())

			wg := keeper.Run(ctx)

			fakedContenderIDs := []domain.ContenderID{
				testutils.RandomResourceID[domain.ContenderID](),
				testutils.RandomResourceID[domain.ContenderID](),
				testutils.RandomResourceID[domain.ContenderID](),
			}

			for _, contenderID := range fakedContenderIDs {
				for k := 1; k <= 3; k++ {
					pointValue := domain.PointValue{
						ContenderID: contenderID,
						ProblemID:   domain.ProblemID(1000 + k),
					}

					err := subscription.Post(domain.EventEnvelope{
						Data: domain.PointValueUpdatedEvent(pointValue),
					})

					require.NoError(t, err)
				}
			}

			synctest.Wait()

			for _, contenderID := range fakedContenderIDs {
				pointValues := keeper.GetPointValues(contenderID)

				assert.ElementsMatch(t, []domain.PointValue{
					{
						ContenderID: contenderID,
						ProblemID:   1001,
					},
					{
						ContenderID: contenderID,
						ProblemID:   1002,
					},
					{
						ContenderID: contenderID,
						ProblemID:   1003,
					},
				}, pointValues)
			}

			<-time.After(time.Hour + time.Minute + time.Nanosecond)

			for _, contenderID := range fakedContenderIDs {
				pointValues := keeper.GetPointValues(contenderID)

				assert.Empty(t, pointValues)
			}

			cancel()

			wg.Wait()

			mockedEventBroker.AssertExpectations(t)
		})
	})
}
