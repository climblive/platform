package testutils

import (
	"testing"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestRandomResourceID(t *testing.T) {
	id1 := RandomResourceID[domain.ContenderID]()
	id2 := RandomResourceID[domain.ContenderID]()
	id3 := RandomResourceID[domain.ContenderID]()

	assert.Greater(t, id1, domain.ContenderID(0))
	assert.Equal(t, id1+1, id2)
	assert.Equal(t, id2+1, id3)
}
