package scrubber

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/climblive/platform/backend/internal/usecases"
	"github.com/climblive/platform/backend/internal/utils"
)

type scrubber struct {
	useCase  *usecases.ScrubberUseCase
	interval time.Duration
}

func New(useCase *usecases.ScrubberUseCase, interval time.Duration) *scrubber {
	return &scrubber{useCase: useCase, interval: interval}
}

func (s *scrubber) Run(ctx context.Context) *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			next := time.Now().Add(s.interval).Round(time.Hour)

			delay := time.Until(next)
			slog.Info("scrubber scheduled", "next_run", next, "delay", delay, "interval", s.interval)

			select {
			case <-ctx.Done():
				return
			case <-time.After(delay):
				slog.Info("running contender scrubber")
				count, err := s.useCase.ScrubContenders(ctx, time.Now().Add(s.interval).Round(time.Hour))
				if err != nil {
					if stack := utils.GetErrorStack(err); stack != "" {
						slog.Error("scrubber error", "stack", stack)
					}
					slog.Error("failed to scrub contenders", "error", err)
				} else {
					slog.Info("contender scrubber completed", "count", count)
				}
			}
		}
	}()

	return &wg
}
