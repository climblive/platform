package domain_test

import (
	"encoding/json"
	"testing"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type data[T any] struct {
	Data domain.Patch[T] `json:"data,omitzero"`
}

func TestNewPatch(t *testing.T) {
	patch := domain.NewPatch("Hello, World!")

	assert.True(t, patch.Present)
	assert.Equal(t, "Hello, World!", patch.Value)
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
