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
		assert.True(t, data.Data.Valid)
		assert.Equal(t, "Hello, World!", data.Data.Value)
	})

	t.Run("WithoutPatchValue", func(t *testing.T) {
		encoded := `{}`

		var data data[string]

		err := json.Unmarshal([]byte(encoded), &data)

		require.NoError(t, err)
		assert.False(t, data.Data.Valid)
		assert.Empty(t, data.Data.Value)
	})

	t.Run("WithPatchPointerNilValue", func(t *testing.T) {
		encoded := `{"data":null}`

		var data data[*int]

		err := json.Unmarshal([]byte(encoded), &data)

		require.NoError(t, err)
		assert.True(t, data.Data.Valid)
		assert.Nil(t, data.Data.Value)
	})

	t.Run("WithPatchPointerValue", func(t *testing.T) {
		encoded := `{"data":5}`

		var data data[*int]

		err := json.Unmarshal([]byte(encoded), &data)

		require.NoError(t, err)
		assert.True(t, data.Data.Valid)
		require.NotNil(t, data.Data.Value)
		assert.Equal(t, 5, *data.Data.Value)
	})
}
