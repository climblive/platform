package authorizer

import (
	"context"
	"net/http"
	"regexp"

	"github.com/climblive/platform/backend/internal/domain"
)

type contextKey struct{}

type authorizerRepository interface {
	domain.Transactor

	GetContenderByCode(ctx context.Context, tx domain.Transaction, registrationCode string) (domain.Contender, error)
}

type Authorizer struct {
	repo                       authorizerRepository
	authorizationHeaderPattern *regexp.Regexp
}

func NewAuthorizer(repo authorizerRepository) *Authorizer {
	re := regexp.MustCompile(`^Regcode ([A-Z0-9]{8})$`)

	return &Authorizer{
		repo:                       repo,
		authorizationHeaderPattern: re,
	}
}

func (a *Authorizer) HasOwnership(ctx context.Context, resourceOwnership domain.OwnershipData) (domain.AuthRole, error) {
	registrationCode, ok := ctx.Value(contextKey{}).(string)
	if !ok {
		return domain.NilRole, domain.ErrNotAuthorized
	}

	contender, err := a.repo.GetContenderByCode(ctx, nil, registrationCode)
	if err != nil {
		return domain.NilRole, domain.ErrNotAuthorized
	}

	if resourceOwnership.ContenderID != nil && contender.ID == *resourceOwnership.ContenderID {
		return domain.ContenderRole, nil
	}

	return domain.NilRole, nil
}

func (a *Authorizer) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		matches := a.authorizationHeaderPattern.FindStringSubmatch(r.Header.Get("Authorization"))

		if matches != nil {
			ctx := context.WithValue(r.Context(), contextKey{}, matches[1])

			next.ServeHTTP(w, r.WithContext(ctx))

			return
		}

		next.ServeHTTP(w, r)
	})
}
