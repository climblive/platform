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

type Scrubber struct {
	useCase  contenderScrubberUseCase
	interval time.Duration
	running  atomic.Bool
}

func New(useCase contenderScrubberUseCase, interval time.Duration) *Scrubber {
	return &Scrubber{useCase: useCase, interval: interval, running: atomic.Bool{}}
}

func (s *Scrubber) Run(ctx context.Context, options ...func(*runOptions)) *sync.WaitGroup {
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

		s.running.Store(true)
		defer s.running.Store(false)

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

func (s *Scrubber) GetStatus() domain.ServiceStatus {
	return domain.ServiceStatus{Name: "Scrubber", Healthy: s.running.Load(), CheckedAt: time.Now()}
}
