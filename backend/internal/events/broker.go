package events

import "github.com/climblive/platform/backend/internal/domain"

type broker struct {
}

func NewBroker() domain.EventBroker {
	return &broker{}
}

func (b *broker) Dispatch(contestID domain.ResourceID, event any) {

}
