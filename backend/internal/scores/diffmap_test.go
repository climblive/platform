package scores_test

import (
	"testing"

	"github.com/climblive/platform/backend/internal/scores"
	"github.com/stretchr/testify/assert"
)

func TestDiffMap(t *testing.T) {
	dm := scores.NewDiffMap[int](func(v1, v2 string) bool { return v1 == v2 })

	dm.Set(1, "Alice")
	dm.Set(2, "Bob")
	dm.Set(3, "Tom")

	diff := dm.Commit()

	assert.Len(t, diff, 3)
	assert.ElementsMatch(t, []string{"Alice", "Bob", "Tom"}, diff)

	dm.Set(1, "Alicia")
	dm.Set(2, "Bob")
	dm.Set(3, "Tommy")
	dm.Set(4, "Eve")

	diff = dm.Commit()

	assert.Len(t, diff, 3)
	assert.ElementsMatch(t, []string{"Alicia", "Tommy", "Eve"}, diff)

	diff = dm.Commit()
	assert.Len(t, diff, 0)
}

func TestDiffMap_RevertValue(t *testing.T) {
	dm := scores.NewDiffMap[int](func(v1, v2 string) bool { return v1 == v2 })

	dm.Set(1, "Alice")

	diff := dm.Commit()
	assert.ElementsMatch(t, []string{"Alice"}, diff)

	dm.Set(1, "Bob")
	dm.Set(1, "Alice")

	diff = dm.Commit()
	assert.Len(t, diff, 0)
}

func TestDiffMap_GetValue(t *testing.T) {
	dm := scores.NewDiffMap[int](func(v1, v2 string) bool { return v1 == v2 })

	dm.Set(1, "Alice")
	dm.Set(2, "Bob")

	name, found := dm.Get(1)
	assert.True(t, found)
	assert.Equal(t, "Alice", name)

	name, found = dm.Get(2)
	assert.True(t, found)
	assert.Equal(t, "Bob", name)

	name, found = dm.Get(3)
	assert.False(t, found)
	assert.Empty(t, name)

	_ = dm.Commit()

	dm.Set(2, "Bobby")

	name, found = dm.Get(1)
	assert.True(t, found)
	assert.Equal(t, "Alice", name)

	name, found = dm.Get(2)
	assert.True(t, found)
	assert.Equal(t, "Bobby", name)
}
