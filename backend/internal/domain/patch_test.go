package domain_test

import (
	"encoding/json/v2"
	"testing"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type data[T comparable] struct {
	Data domain.Patch[T] `json:"data,omitzero"`
}

func TestNewPatch(t *testing.T) {
	patch := domain.NewPatch("Hello, World!")

	assert.True(t, patch.Present)
	assert.Equal(t, "Hello, World!", patch.Value)
}

func TestDistinct(t *testing.T) {
	patch := domain.NewPatch("Hello, World!")

	assert.False(t, patch.PresentAndDistinct("Hello, World!"))
	assert.True(t, patch.PresentAndDistinct("Hello, Universe!"))
}

func TestDistinct_Empty(t *testing.T) {
	patch := domain.Patch[string]{}

	assert.False(t, patch.PresentAndDistinct("Anything"))
}

func TestPatchMarshal(t *testing.T) {
	t.Run("WithPatchValue", func(t *testing.T) {
		data := data[string]{
			Data: domain.NewPatch("Hello, World!"),
		}

		encoded, err := json.Marshal(data)

		require.NoError(t, err)
		assert.Equal(t, `{"data":"Hello, World!"}`, string(encoded))
	})

	t.Run("WithoutPatchValue", func(t *testing.T) {
		data := data[string]{}

		encoded, err := json.Marshal(data)

		require.NoError(t, err)
		assert.Equal(t, `{}`, string(encoded))
	})

	t.Run("WithPatchPointerValue", func(t *testing.T) {
		data := data[*int]{
			Data: domain.NewPatch[*int](nil),
		}

		encoded, err := json.Marshal(data)

		require.NoError(t, err)
		assert.Equal(t, `{"data":null}`, string(encoded))
	})
}

func TestPatchUnmarshal(t *testing.T) {
	t.Run("WithPatchValue", func(t *testing.T) {
		encoded := `{"data":"Hello, World!"}`

		var data data[string]

		err := json.Unmarshal([]byte(encoded), &data)

		require.NoError(t, err)
		assert.True(t, data.Data.Present)
		assert.Equal(t, "Hello, World!", data.Data.Value)
	})

	t.Run("WithoutPatchValue", func(t *testing.T) {
		encoded := `{}`

		var data data[string]

		err := json.Unmarshal([]byte(encoded), &data)

		require.NoError(t, err)
		assert.False(t, data.Data.Present)
		assert.Empty(t, data.Data.Value)
	})

	t.Run("WithPatchPointerNilValue", func(t *testing.T) {
		encoded := `{"data":null}`

		var data data[*int]

		err := json.Unmarshal([]byte(encoded), &data)

		require.NoError(t, err)
		assert.True(t, data.Data.Present)
		assert.Nil(t, data.Data.Value)
	})

	t.Run("WithPatchPointerValue", func(t *testing.T) {
		encoded := `{"data":5}`

		var data data[*int]

		err := json.Unmarshal([]byte(encoded), &data)

		require.NoError(t, err)
		assert.True(t, data.Data.Present)
		require.NotNil(t, data.Data.Value)
		assert.Equal(t, 5, *data.Data.Value)
	})

	t.Run("WithNestedPatch", func(t *testing.T) {
		encoded := `{"data":{"data":{"data":5}}}`

		var data data[data[data[int]]]

		err := json.Unmarshal([]byte(encoded), &data)

		require.NoError(t, err)
		assert.True(t, data.Data.Present)

		nested1 := data.Data.Value

		assert.True(t, nested1.Data.Present)

		nested2 := nested1.Data.Value

		assert.True(t, nested2.Data.Present)
		assert.Equal(t, 5, nested2.Data.Value)
	})
}

func TestCategoryEligibilityPatchJSON(t *testing.T) {
	for _, tc := range []struct {
		name     string
		encoded  string
		expected domain.CategoryPatch
	}{
		{name: "Omitted", encoded: `{}`},
		{name: "ZeroAge", encoded: `{"minAge":0}`, expected: domain.CategoryPatch{MinAge: domain.NewPatch(new(0))}},
		{name: "ClearAgeBounds", encoded: `{"minAge":null,"maxAge":null}`, expected: domain.CategoryPatch{MinAge: domain.NewPatch[*int](nil), MaxAge: domain.NewPatch[*int](nil)}},
		{name: "BirthYearBounds", encoded: `{"eligibilityBasis":"birth_year","minBirthYear":2010,"maxBirthYear":2014}`, expected: domain.CategoryPatch{EligibilityBasis: domain.NewPatch(domain.EligibilityBasisBirthYear), MinBirthYear: domain.NewPatch(new(2010)), MaxBirthYear: domain.NewPatch(new(2014))}},
		{name: "ClearBirthYearBounds", encoded: `{"minBirthYear":null,"maxBirthYear":null}`, expected: domain.CategoryPatch{MinBirthYear: domain.NewPatch[*int](nil), MaxBirthYear: domain.NewPatch[*int](nil)}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var patch domain.CategoryPatch
			require.NoError(t, json.Unmarshal([]byte(tc.encoded), &patch))
			assert.Equal(t, tc.expected, patch)
			encoded, err := json.Marshal(patch)
			require.NoError(t, err)
			assert.JSONEq(t, tc.encoded, string(encoded))
		})
	}
}
