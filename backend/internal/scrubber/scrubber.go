package scrubber

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
)

type runOptions struct {
	recoverPanics bool
}

func WithPanicRecovery() func(*runOptions) {
	return func(s *runOptions) {
		s.recoverPanics = true
	}
}

type contenderScrubberUseCase interface {
	ScrubContenders(ctx context.Context, deadline time.Time) (int, error)
}

type scrubber struct {
	useCase    contenderScrubberUseCase
	interval   time.Duration
	running    int32
	lastSeenAt int64
}

func New(useCase contenderScrubberUseCase, interval time.Duration) *scrubber {
	return &scrubber{useCase: useCase, interval: interval}
}

func (s *scrubber) Run(ctx context.Context, options ...func(*runOptions)) *sync.WaitGroup {
	config := &runOptions{}
	for _, opt := range options {
		opt(config)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer func() {
			if !config.recoverPanics {
				return
			}

			if r := recover(); r != nil {
				slog.Error("scrubber panicked", "error", r)
			}
		}()

		defer wg.Done()

		atomic.StoreInt32(&s.running, 1)
		defer func() {
			atomic.StoreInt64(&s.lastSeenAt, time.Now().UnixNano())
			atomic.StoreInt32(&s.running, 0)
		}()

		for {
			next := time.Now().Add(s.interval).Round(time.Hour)

			delay := time.Until(next)
			slog.Info("scrubber scheduled", "next_run", next, "delay", delay, "interval", s.interval)

			select {
			case <-ctx.Done():
				return
			case <-time.After(delay):
				slog.Info("running scrubber")

				count, err := s.useCase.ScrubContenders(ctx, time.Now().Add(s.interval).Round(time.Hour))
				if err != nil {
					slog.Error("failed to scrub contenders", "error", err)
				} else if count > 0 {
					slog.Info("scrubber completed", "count", count)
				}
			}
		}
	}()

	return &wg
}

func (s *scrubber) GetStatus() domain.RunnerStatus {
	if atomic.LoadInt32(&s.running) == 1 {
		return domain.RunnerStatus{Healthy: true, CheckedAt: time.Now()}
	}

	ns := atomic.LoadInt64(&s.lastSeenAt)
	if ns == 0 {
		return domain.RunnerStatus{}
	}

	return domain.RunnerStatus{
		Healthy:   false,
		CheckedAt: time.Unix(0, ns),
	}
}
