package validators

import (
	_ "embed"
	"encoding/json"
	"strings"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

var errContestConstraintViolation = errors.New("constraint violation")

//go:embed countries.json
var countriesJSON []byte

type countryData struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

var validCountryCodes map[string]bool

func init() {
	var countries []countryData
	if err := json.Unmarshal(countriesJSON, &countries); err != nil {
		panic("failed to load countries.json: " + err.Error())
	}

	validCountryCodes = make(map[string]bool, len(countries))
	for _, country := range countries {
		validCountryCodes[country.Code] = true
	}
}

type ContestValidator struct {
}

func (v ContestValidator) Validate(contest domain.Contest) error {
	switch {
	case len(strings.TrimSpace(contest.Name)) < 1:
		fallthrough
	case len(contest.Country) != 2 || !validCountryCodes[contest.Country]:
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
