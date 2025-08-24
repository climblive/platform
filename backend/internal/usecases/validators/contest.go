package validators

import (
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

var errContestConstraintViolation = errors.New("constraint violation")

type ContestValidator struct {
}

func (v ContestValidator) Validate(contest domain.Contest) error {
	switch {
	case len(contest.Name) < 1:
		fallthrough
	case contest.Finalists < 0 || contest.Finalists > 65536:
		fallthrough
	case contest.QualifyingProblems < 0 || contest.QualifyingProblems > 65536:
		fallthrough
	case contest.GracePeriod < 0 || contest.GracePeriod > time.Hour:
		return errors.Errorf("%w: %w", domain.ErrInvalidData, errContestConstraintViolation)
	}

	return nil
}

func (v ContestValidator) IsValidationError(err error) bool {
	return errors.Is(err, errContestConstraintViolation)
}
