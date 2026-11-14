package validators

import (
	_ "embed"
	"encoding/json/v2"
	"strings"
	"time"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/go-errors/errors"
)

var errCompetitionConstraintViolation = errors.New("constraint violation")

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
		panic(err)
	}

	validCountryCodes = make(map[string]bool, len(countries))
	for _, country := range countries {
		validCountryCodes[country.Code] = true
	}
}

type CompetitionValidator struct {
}

func (v CompetitionValidator) Validate(competition domain.Competition) error {
	switch {
	case len(strings.TrimSpace(competition.Name)) < 1:
		fallthrough
	case len(competition.Country) != 2 || !validCountryCodes[competition.Country]:
		fallthrough
	case competition.Finalists < 0 || competition.Finalists > 65536:
		fallthrough
	case competition.QualifyingProblems < 0 || competition.QualifyingProblems > 65536:
		fallthrough
	case !competition.UsePoints && competition.PooledPoints:
		fallthrough
	case !competition.UsePoints && competition.QualifyingProblems > 0:
		fallthrough
	case !competition.UsePoints && competition.MaxAttempts > 0:
		fallthrough
	case !competition.UsePoints && competition.PointDeduction > 0:
		fallthrough
	case competition.MaxAttempts < 0 || competition.MaxAttempts > 999:
		fallthrough
	case competition.PointDeduction < 0 || competition.PointDeduction > domain.Tenths(maxAllowedPointValue)*10:
		fallthrough
	case competition.GracePeriod < 0 || competition.GracePeriod > time.Hour:
		fallthrough
	case !isValidNameRetentionTime(competition.NameRetentionTime):
		return errors.Errorf("%w: %w", domain.ErrInvalidData, errCompetitionConstraintViolation)
	}

	return nil
}

func (v CompetitionValidator) IsValidationError(err error) bool {
	return errors.Is(err, errCompetitionConstraintViolation)
}

func isValidNameRetentionTime(d time.Duration) bool {
	if d >= 14*24*time.Hour && d <= 90*24*time.Hour {
		return true
	}

	return false
}
