package scrubber_test

import (
	"context"
	"testing"
	"testing/synctest"
	"time"

	"github.com/climblive/platform/backend/internal/scrubber"
	"github.com/stretchr/testify/mock"
)

type contenderScrubberMock struct {
	mock.Mock
}

func (m *contenderScrubberMock) ScrubContenders(ctx context.Context, deadline time.Time) (int, error) {
	args := m.Called(ctx, deadline)
	return args.Int(0), args.Error(1)
}

func TestScrubber(t *testing.T) {
	t.Run("RunsOncePerInterval", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			mockedScrubber := new(contenderScrubberMock)

			interval := time.Hour

			start := time.Now()

			mockedScrubber.On("ScrubContenders", mock.Anything, start.Add(2*interval)).
				Return(1, nil).Once()
			mockedScrubber.On("ScrubContenders", mock.Anything, start.Add(3*interval)).
				Return(2, nil).Once()
			mockedScrubber.On("ScrubContenders", mock.Anything, start.Add(4*interval)).
				Return(3, nil).Once()

			scrubber := scrubber.New(mockedScrubber, interval)
			ctx, cancel := context.WithCancel(context.Background())

			wg := scrubber.Run(ctx)

			time.Sleep(interval)

			time.Sleep(interval)

			time.Sleep(interval)

			cancel()
			wg.Wait()

			mockedScrubber.AssertExpectations(t)
		})
	})
}
