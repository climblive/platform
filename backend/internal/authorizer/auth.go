package authorizer

import (
	"context"

	"github.com/climblive/platform/backend/internal/domain"
)

type authorizer struct {
}

func NewAuthorizer() domain.Authorizer {
	return &authorizer{}
}

func (a *authorizer) HasOwnership(ctx context.Context, resourceOwnership domain.OwnershipData) (*domain.AuthRole, error) {
	role := domain.AdminRole
	return &role, nil
}
