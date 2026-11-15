package testutils

import (
	"testing"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestRandomResourceID(t *testing.T) {
	id1 := RandomResourceID[domain.CompetitorID]()
	id2 := RandomResourceID[domain.CompetitorID]()
	id3 := RandomResourceID[domain.CompetitorID]()

	assert.Greater(t, id1, domain.CompetitorID(0))
	assert.Equal(t, id1+1, id2)
	assert.Equal(t, id2+1, id3)
}
