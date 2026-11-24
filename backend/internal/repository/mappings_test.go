package repository

import (
	"database/sql"
	"testing"

	"github.com/climblive/platform/backend/internal/database"
	"github.com/climblive/platform/backend/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestNullableCategoryBounds(t *testing.T) {
	for _, tc := range []struct {
		name   string
		value  *int
		stored sql.NullInt32
	}{
		{name: "Unbounded", stored: sql.NullInt32{}},
		{name: "Zero", value: new(0), stored: sql.NullInt32{Int32: 0, Valid: true}},
		{name: "Age", value: new(18), stored: sql.NullInt32{Int32: 18, Valid: true}},
		{name: "BirthYear", value: new(2010), stored: sql.NullInt32{Int32: 2010, Valid: true}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.stored, makeNullIntPointer(tc.value))
			assert.Equal(t, tc.value, nullIntToPointer(tc.stored))
		})
	}
}

func TestCategoryEligibilityMapping(t *testing.T) {
	for _, tc := range []struct {
		name         string
		basis        domain.EligibilityBasis
		minAge       *int
		maxAge       *int
		minBirthYear *int
		maxBirthYear *int
	}{
		{name: "AgeRange", basis: domain.EligibilityBasisAge, minAge: new(18), maxAge: new(99)},
		{name: "AgeZero", basis: domain.EligibilityBasisAge, minAge: new(0)},
		{name: "BirthYearRange", basis: domain.EligibilityBasisBirthYear, minBirthYear: new(2010), maxBirthYear: new(2014)},
		{name: "UnboundedBirthYears", basis: domain.EligibilityBasisBirthYear},
	} {
		t.Run(tc.name, func(t *testing.T) {
			category := categoryToDomain(database.Category{
				Description:      sql.NullString{String: "Retained for rollback", Valid: true},
				EligibilityBasis: database.CompClassEligibilityBasis(tc.basis),
				MinAge:           makeNullIntPointer(tc.minAge),
				MaxAge:           makeNullIntPointer(tc.maxAge),
				MinBirthYear:     makeNullIntPointer(tc.minBirthYear),
				MaxBirthYear:     makeNullIntPointer(tc.maxBirthYear),
			})
			assert.Equal(t, tc.basis, category.EligibilityBasis)
			assert.Equal(t, tc.minAge, category.MinAge)
			assert.Equal(t, tc.maxAge, category.MaxAge)
			assert.Equal(t, tc.minBirthYear, category.MinBirthYear)
			assert.Equal(t, tc.maxBirthYear, category.MaxBirthYear)
			assert.Equal(t, "Retained for rollback", category.Description)
		})
	}
}
