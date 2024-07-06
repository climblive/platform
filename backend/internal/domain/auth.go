package domain

import (
	"context"
)

type AuthRole string

const (
	ContenderRole AuthRole = "contender"
	JudgeRole     AuthRole = "judge"
	OrganizerRole AuthRole = "organizer"
	AdminRole     AuthRole = "admin"
)

func (role AuthRole) OneOf(roles ...AuthRole) bool {
	for _, otherRole := range roles {
		if role == otherRole {
			return true
		}
	}

	return false
}

type Authorizer interface {
	HasOwnership(ctx context.Context, resourceOwnership OwnershipData) (*AuthRole, error)
}

type OwnershipData struct {
	OrganizerID ResourceID
	ContenderID *ResourceID
}
