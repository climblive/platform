package scores_test

import (
	"math/rand"
	"slices"
	"testing"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/climblive/platform/backend/internal/scores"
	"github.com/stretchr/testify/assert"
)

func TestFilterByCompClass(t *testing.T) {
	store := scores.NewMemoryStore()

	store.SaveContender(scores.Contender{
		ID:          1,
		CompClassID: 1,
	})
	store.SaveContender(scores.Contender{
		ID:          2,
		CompClassID: 1,
	})
	store.SaveContender(scores.Contender{
		ID:          3,
		CompClassID: 2,
	})
	store.SaveContender(scores.Contender{
		ID:          4,
		CompClassID: 3,
	})
	store.SaveContender(scores.Contender{
		ID:          5,
		CompClassID: 1,
	})
	store.SaveContender(scores.Contender{
		ID:          6,
		CompClassID: 2,
	})
	store.SaveContender(scores.Contender{
		ID:          7,
		CompClassID: 1,
	})

	t.Run("CompClassOne", func(t *testing.T) {
		var filteredContenders []domain.ContenderID
		for contender := range slices.Values(slices.Collect(store.GetContendersByCompClass(1))) {
			filteredContenders = append(filteredContenders, contender.ID)
		}

		assert.ElementsMatch(t, filteredContenders, []domain.ContenderID{1, 2, 5, 7})
	})

	t.Run("CompClassTwo", func(t *testing.T) {
		var filteredContenders []domain.ContenderID
		for contender := range slices.Values(slices.Collect(store.GetContendersByCompClass(2))) {
			filteredContenders = append(filteredContenders, contender.ID)
		}

		assert.ElementsMatch(t, filteredContenders, []domain.ContenderID{3, 6})
	})

	t.Run("CompClassThree", func(t *testing.T) {
		var filteredContenders []domain.ContenderID
		for contender := range slices.Values(slices.Collect(store.GetContendersByCompClass(3))) {
			filteredContenders = append(filteredContenders, contender.ID)
		}

		assert.ElementsMatch(t, filteredContenders, []domain.ContenderID{4})
	})
}

func TestCompClasses(t *testing.T) {
	store := scores.NewMemoryStore()

	store.SaveContender(scores.Contender{
		ID:          domain.ContenderID(rand.Int()),
		CompClassID: 1,
	})
	store.SaveContender(scores.Contender{
		ID:          domain.ContenderID(rand.Int()),
		CompClassID: 1,
	})
	store.SaveContender(scores.Contender{
		ID:          domain.ContenderID(rand.Int()),
		CompClassID: 2,
	})
	store.SaveContender(scores.Contender{
		ID:          domain.ContenderID(rand.Int()),
		CompClassID: 3,
	})
	store.SaveContender(scores.Contender{
		ID:          domain.ContenderID(rand.Int()),
		CompClassID: 1,
	})
	store.SaveContender(scores.Contender{
		ID:          domain.ContenderID(rand.Int()),
		CompClassID: 2,
	})
	store.SaveContender(scores.Contender{
		ID:          domain.ContenderID(rand.Int()),
		CompClassID: 1,
	})

	var compClasses []domain.CompClassID

	for _, compClassID := range store.GetCompClassIDs() {
		compClasses = append(compClasses, compClassID)
	}

	assert.ElementsMatch(t, compClasses, []domain.CompClassID{1, 2, 3})
}
