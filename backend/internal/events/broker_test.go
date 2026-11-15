package events_test

import (
	"testing"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/events"
	"github.com/stretchr/testify/assert"
)

func TestResultEventNames(t *testing.T) {
	tests := []struct {
		name     string
		event    any
		expected string
	}{
		{name: "Single", event: domain.ContenderResultUpdatedEvent{}, expected: "CONTENDER_RESULT_UPDATED"},
		{name: "Batch", event: []domain.ContenderResultUpdatedEvent{}, expected: "[]CONTENDER_RESULT_UPDATED"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, events.EventName(tt.event))
		})
	}
}

func TestBlockingSubscriber(t *testing.T) {
	broker := events.NewBroker()
	filter := domain.EventFilter{
		CompetitionID: 1,
	}

	_, _ = broker.Subscribe([]domain.EventFilter{filter}, 1)

	for range 100 {
		broker.Dispatch(1, domain.ContenderEnteredEvent{
			ContenderID: 1,
		})
	}
}
