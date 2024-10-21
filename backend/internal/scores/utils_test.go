package scores_test

import (
	"slices"
	"testing"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/scores"
	"github.com/stretchr/testify/assert"
)

func TestFilterByClass(t *testing.T) {
	contenders := make(map[domain.ContenderID]*scores.Contender)

	contenders[1] = &scores.Contender{
		ID:          1,
		CompClassID: 1,
	}
	contenders[2] = &scores.Contender{
		ID:          2,
		CompClassID: 1,
	}
	contenders[3] = &scores.Contender{
		ID:          3,
		CompClassID: 2,
	}
	contenders[4] = &scores.Contender{
		ID:          4,
		CompClassID: 3,
	}
	contenders[5] = &scores.Contender{
		ID:          5,
		CompClassID: 1,
	}
	contenders[6] = &scores.Contender{
		ID:          6,
		CompClassID: 2,
	}
	contenders[7] = &scores.Contender{
		ID:          7,
		CompClassID: 1,
	}

	t.Run("CompClassOne", func(t *testing.T) {
		var filteredContenders []domain.ContenderID
		for contender := range slices.Values(slices.Collect(scores.FilterByClass(contenders, 1))) {
			filteredContenders = append(filteredContenders, contender.ID)
		}

		assert.ElementsMatch(t, filteredContenders, []domain.ContenderID{1, 2, 5, 7})
	})

	t.Run("CompClassTwo", func(t *testing.T) {
		var filteredContenders []domain.ContenderID
		for contender := range slices.Values(slices.Collect(scores.FilterByClass(contenders, 2))) {
			filteredContenders = append(filteredContenders, contender.ID)
		}

		assert.ElementsMatch(t, filteredContenders, []domain.ContenderID{3, 6})
	})

	t.Run("CompClassThree", func(t *testing.T) {
		var filteredContenders []domain.ContenderID
		for contender := range slices.Values(slices.Collect(scores.FilterByClass(contenders, 3))) {
			filteredContenders = append(filteredContenders, contender.ID)
		}

		assert.ElementsMatch(t, filteredContenders, []domain.ContenderID{4})
	})
}
