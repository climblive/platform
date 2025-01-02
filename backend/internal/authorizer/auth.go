package authorizer

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

type contextKey struct{}

type authenticationResult struct {
	regcode  string
	username string
	err      error
}

type authorizerRepository interface {
	domain.Transactor

	GetContenderByCode(ctx context.Context, tx domain.Transaction, registrationCode string) (domain.Contender, error)
	GetUserByUsername(ctx context.Context, tx domain.Transaction, username string) (domain.User, error)
	StoreUser(ctx context.Context, tx domain.Transaction, user domain.User) (domain.User, error)
	StoreOrganizer(ctx context.Context, tx domain.Transaction, organizer domain.Organizer) (domain.Organizer, error)
	AddUserToOrganizer(ctx context.Context, tx domain.Transaction, userID domain.UserID, organizerID domain.OrganizerID) error
}

type Claims struct {
	Username   string `json:"username"`
	Expiration int64  `json:"exp"`
	Issuer     string `json:"iss"`
}

type JWTDecoder interface {
	Decode(jwt string) (Claims, error)
}

type Authorizer struct {
	repo           authorizerRepository
	regcodePattern *regexp.Regexp
	jwtDecoder     JWTDecoder
}

func NewAuthorizer(repo authorizerRepository, jwtDecoder JWTDecoder) *Authorizer {
	return &Authorizer{
		repo:           repo,
		regcodePattern: regexp.MustCompile(`^Regcode ([A-Za-z0-9]{8})$`),
		jwtDecoder:     jwtDecoder,
	}
}

func (a *Authorizer) HasOwnership(ctx context.Context, resourceOwnership domain.OwnershipData) (domain.AuthRole, error) {
	authenticationResult, ok := ctx.Value(contextKey{}).(authenticationResult)
	if !ok {
		return domain.NilRole, domain.ErrNotAuthenticated
	}

	if authenticationResult.err != nil {
		return domain.NilRole, errors.Errorf("%w: %w", domain.ErrNotAuthenticated, authenticationResult.err)
	}

	if authenticationResult.regcode != "" {
		return a.checkRegcode(ctx, authenticationResult.regcode, resourceOwnership)
	} else if authenticationResult.username != "" {
		return a.checkUsername(ctx, authenticationResult.username, resourceOwnership)
	}

	return domain.NilRole, domain.ErrNoOwnership
}

func (a *Authorizer) checkRegcode(ctx context.Context, regcode string, resourceOwnership domain.OwnershipData) (domain.AuthRole, error) {
	contender, err := a.repo.GetContenderByCode(ctx, nil, strings.ToUpper(regcode))
	if err != nil {
		return domain.NilRole, domain.ErrNotAuthorized
	}

	if resourceOwnership.ContenderID != nil && contender.ID == *resourceOwnership.ContenderID {
		return domain.ContenderRole, nil
	}

	return domain.NilRole, domain.ErrNoOwnership
}

func (a *Authorizer) checkUsername(ctx context.Context, username string, resourceOwnership domain.OwnershipData) (domain.AuthRole, error) {
	user, err := a.repo.GetUserByUsername(ctx, nil, username)

	switch {
	case errors.Is(err, domain.ErrNotFound):
		err := a.createUser(ctx, username)
		if err != nil {
			return domain.NilRole, errors.Wrap(err, 0)
		}

		return domain.NilRole, domain.ErrNoOwnership
	case err != nil:
		return domain.NilRole, domain.ErrNotAuthorized
	}

	if user.Admin {
		return domain.AdminRole, nil
	}

	for _, organizerID := range user.Organizers {
		if organizerID == resourceOwnership.OrganizerID {
			return domain.OrganizerRole, nil
		}
	}

	return domain.NilRole, domain.ErrNoOwnership
}

func (a *Authorizer) createUser(ctx context.Context, username string) error {
	tx, err := a.repo.Begin()
	if err != nil {
		return errors.Wrap(err, 0)
	}

	exec := func() error {
		organizer, err := a.repo.StoreOrganizer(ctx, tx, domain.Organizer{
			Name: fmt.Sprintf("%s's organizer", username),
		})
		if err != nil {
			return errors.Wrap(err, 0)
		}

		user, err := a.repo.StoreUser(ctx, tx, domain.User{
			Name:       username,
			Username:   username,
			Organizers: []domain.OrganizerID{organizer.ID},
		})
		if err != nil {
			return errors.Wrap(err, 0)
		}

		err = a.repo.AddUserToOrganizer(ctx, tx, user.ID, organizer.ID)
		if err != nil {
			return errors.Wrap(err, 0)
		}

		return nil
	}

	err = exec()
	if err != nil {
		tx.Rollback()

		return errors.Wrap(err, 0)
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, 0)
	}

	return nil
}

func (a *Authorizer) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		nextCtx := r.Context()

		var bearer string
		var err error
		var claims Claims

		matches := a.regcodePattern.FindStringSubmatch(header)
		if matches != nil {
			nextCtx = context.WithValue(nextCtx, contextKey{}, authenticationResult{regcode: matches[1]})
			goto Next
		}

		if n, err := fmt.Sscanf(header, "Bearer %s", &bearer); err != nil || n != 1 {
			goto Next
		}

		claims, err = a.jwtDecoder.Decode(bearer)
		nextCtx = context.WithValue(nextCtx, contextKey{}, authenticationResult{username: claims.Username, err: err})

	Next:
		next.ServeHTTP(w, r.WithContext(nextCtx))
	})
}
