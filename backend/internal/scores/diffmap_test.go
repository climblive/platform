package scores_test

import (
	"testing"

	"github.com/climblive/platform/backend/internal/scores"
	"github.com/stretchr/testify/assert"
)

func TestDiffMap(t *testing.T) {
	dm := scores.NewDiffMap[int, string](func(v1, v2 string) bool { return v1 == v2 })

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
