package events

import (
	"context"
	"sync"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
	"github.com/google/uuid"
)

var ErrBufferFull = errors.New("buffer full")

type Subscription struct {
	ID             domain.SubscriptionID
	Filter         domain.EventFilter
	mu             sync.Mutex
	cond           *sync.Cond
	buffer         []domain.EventEnvelope
	bufferCapacity int
	closeReason    error
}

func NewSubscription(
	filter domain.EventFilter,
	bufferSize int,
) *Subscription {
	sub := Subscription{
		ID:             uuid.New(),
		Filter:         filter,
		bufferCapacity: bufferSize,
	}

	sub.cond = sync.NewCond(&sub.mu)

	return &sub
}

func (s *Subscription) Await(ctx context.Context) (domain.EventEnvelope, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for {
		if ctx.Err() != nil {
			return domain.EventEnvelope{}, ctx.Err()
		}

		if s.closeReason != nil {
			return domain.EventEnvelope{}, s.closeReason
		}

		event, ok := s.popQueueUnsafe()
		if ok {
			return event, nil
		}

		stop := context.AfterFunc(ctx, func() {
			s.mu.Lock()
			defer s.mu.Unlock()

			s.cond.Broadcast()
		})

		s.cond.Wait()
		stop()
	}
}

func (s *Subscription) popQueueUnsafe() (domain.EventEnvelope, bool) {
	if len(s.buffer) > 0 {
		event := s.buffer[0]
		s.buffer = s.buffer[1:]

		return event, true
	}

	return domain.EventEnvelope{}, false
}

func (s *Subscription) Post(event domain.EventEnvelope) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.bufferCapacity != 0 && len(s.buffer) == s.bufferCapacity {
		s.closeReason = ErrBufferFull

		return ErrBufferFull
	}

	s.buffer = append(s.buffer, event)
	s.cond.Broadcast()

	return nil
}