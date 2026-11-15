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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"uuid"
)

func TestPointValueKeeper(t *testing.T) {
	makeMocks := func(bufferCapacity int) (*eventBrokerMock, *events.Subscription) {
		mockedEventBroker := new(eventBrokerMock)

		filters := []domain.EventFilter{domain.NewEventFilter(
			0,
			0,
			"POINT_VALUE_UPDATED",
		)}

		subscription := events.NewSubscription(filters, bufferCapacity)
		subscriptionID := uuid.New()

		mockedEventBroker.On("Subscribe", filters, 0).Return(subscriptionID, subscription)

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

	t.Run("HealthCheck", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			mockedEventBroker, _ := makeMocks(0)
			keeper := scores.NewPointValueKeeper(mockedEventBroker, time.Hour)

			ctx, cancel := context.WithCancel(context.Background())

			wg := keeper.Run(ctx)

			assert.Equal(t, domain.ServiceStatus{
				Name:      "PointValueKeeper",
				Healthy:   true,
				CheckedAt: time.Now(),
			}, keeper.GetStatus())

			cancel()

			wg.Wait()

			assert.Equal(t, domain.ServiceStatus{
				Name:      "PointValueKeeper",
				Healthy:   false,
				CheckedAt: time.Now(),
			}, keeper.GetStatus())

			mockedEventBroker.AssertExpectations(t)
		})
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

			fakedCompetitorIDs := []domain.CompetitorID{
				testutils.RandomResourceID[domain.CompetitorID](),
				testutils.RandomResourceID[domain.CompetitorID](),
				testutils.RandomResourceID[domain.CompetitorID](),
			}

			for _, competitorID := range fakedCompetitorIDs {
				for k := 1; k <= 3; k++ {
					pointValue := domain.PointValue{
						CompetitorID: competitorID,
						ProblemID:    domain.ProblemID(1000 + k),
						Current:      domain.Tenths(1000 + k),
						Zone1:        domain.Tenths(10 + k),
						Zone2:        domain.Tenths(20 + k),
						Top:          domain.Tenths(1000 + k),
						FlashBonus:   domain.Tenths(100 + k),
					}

					err := subscription.Post(domain.EventEnvelope{
						Data: domain.PointValueUpdatedEvent(pointValue),
					})

					require.NoError(t, err)
				}
			}

			synctest.Wait()

			for _, competitorID := range fakedCompetitorIDs {
				pointValues := keeper.GetPointValues(competitorID)

				assert.ElementsMatch(t, []domain.PointValue{
					{
						CompetitorID: competitorID,
						ProblemID:    1001,
						Current:      1001,
						Zone1:        11,
						Zone2:        21,
						Top:          1001,
						FlashBonus:   101,
					},
					{
						CompetitorID: competitorID,
						ProblemID:    1002,
						Current:      1002,
						Zone1:        12,
						Zone2:        22,
						Top:          1002,
						FlashBonus:   102,
					},
					{
						CompetitorID: competitorID,
						ProblemID:    1003,
						Current:      1003,
						Zone1:        13,
						Zone2:        23,
						Top:          1003,
						FlashBonus:   103,
					},
				}, pointValues)
			}

			cancel()

			wg.Wait()

			mockedEventBroker.AssertExpectations(t)
		})
	})

	t.Run("GatherScores_MultipleUpdates", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			mockedEventBroker, subscription := makeMocks(0)
			keeper := scores.NewPointValueKeeper(mockedEventBroker, time.Hour)

			ctx, cancel := context.WithCancel(context.Background())

			wg := keeper.Run(ctx)

			fakedCompetitorID := testutils.RandomResourceID[domain.CompetitorID]()
			fakedProblemID := testutils.RandomResourceID[domain.ProblemID]()

			for k := 1; k <= 10; k++ {
				pointValue := domain.PointValue{
					CompetitorID: fakedCompetitorID,
					ProblemID:    fakedProblemID,
					Current:      domain.Tenths(1000 + k),
					Zone1:        domain.Tenths(10 + k),
					Zone2:        domain.Tenths(20 + k),
					Top:          domain.Tenths(1000 + k),
					FlashBonus:   domain.Tenths(100 + k),
				}

				err := subscription.Post(domain.EventEnvelope{
					Data: domain.PointValueUpdatedEvent(pointValue),
				})

				require.NoError(t, err)
			}

			synctest.Wait()

			pointValues := keeper.GetPointValues(fakedCompetitorID)

			assert.ElementsMatch(t, []domain.PointValue{
				{
					CompetitorID: fakedCompetitorID,
					ProblemID:    fakedProblemID,
					Current:      1010,
					Zone1:        20,
					Zone2:        30,
					Top:          1010,
					FlashBonus:   110,
				},
			}, pointValues)

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

			fakedCompetitorIDs := []domain.CompetitorID{
				testutils.RandomResourceID[domain.CompetitorID](),
				testutils.RandomResourceID[domain.CompetitorID](),
				testutils.RandomResourceID[domain.CompetitorID](),
			}

			for _, competitorID := range fakedCompetitorIDs {
				for k := 1; k <= 2; k++ {
					pointValue := domain.PointValue{
						CompetitorID: competitorID,
						ProblemID:    domain.ProblemID(1000 + k),
					}

					err := subscription.Post(domain.EventEnvelope{
						Data: domain.PointValueUpdatedEvent(pointValue),
					})

					require.NoError(t, err)
				}
			}

			time.Sleep(5 * time.Minute)

			for _, competitorID := range fakedCompetitorIDs {
				pointValue := domain.PointValue{
					CompetitorID: competitorID,
					ProblemID:    domain.ProblemID(1003),
				}

				err := subscription.Post(domain.EventEnvelope{
					Data: domain.PointValueUpdatedEvent(pointValue),
				})

				require.NoError(t, err)
			}

			synctest.Wait()

			for _, competitorID := range fakedCompetitorIDs {
				pointValues := keeper.GetPointValues(competitorID)

				assert.ElementsMatch(t, []domain.PointValue{
					{
						CompetitorID: competitorID,
						ProblemID:    1001,
					},
					{
						CompetitorID: competitorID,
						ProblemID:    1002,
					},
					{
						CompetitorID: competitorID,
						ProblemID:    1003,
					},
				}, pointValues)
			}

			time.Sleep(time.Hour)

			for _, competitorID := range fakedCompetitorIDs {
				pointValues := keeper.GetPointValues(competitorID)

				assert.Equal(t, []domain.PointValue{
					{
						CompetitorID: competitorID,
						ProblemID:    1003,
					},
				}, pointValues)
			}

			cancel()

			wg.Wait()

			mockedEventBroker.AssertExpectations(t)
		})
	})
}
