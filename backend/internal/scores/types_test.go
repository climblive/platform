package scores

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompareContender(t *testing.T) {
	t.Run("ByScore", func(t *testing.T) {
		c1 := Contender{
			Score: 200,
		}

		c2 := Contender{
			Score: 100,
		}

		assert.Less(t, c1.Compare(c2), 0)
		assert.Greater(t, c2.Compare(c1), 0)
	})

	t.Run("TieBreak", func(t *testing.T) {
		c1 := Contender{
			ID:    1,
			Score: 100,
		}

		c2 := Contender{
			ID:    2,
			Score: 100,
		}

		assert.Less(t, c1.Compare(c2), 0)
		assert.Greater(t, c2.Compare(c1), 0)
	})
}
