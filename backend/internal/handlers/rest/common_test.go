package rest

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"uuid"

	"github.com/climblive/platform/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadJSONRejectsOversizedPayload(t *testing.T) {
	body := `{"value":"` + strings.Repeat("a", maxJSONBodySize) + `"}`
	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	var payload struct {
		Value string `json:"value"`
	}

	assert.False(t, readJSON(response, request, &payload))
	assert.Equal(t, http.StatusRequestEntityTooLarge, response.Code)
}

func TestParseCompetitionID(t *testing.T) {
	id := uuid.New()
	parsed, err := parseCompetitionID(id.String())
	require.NoError(t, err)
	assert.Equal(t, domain.CompetitionID(id), parsed)
	for _, value := range []string{"", "1", "not-a-uuid"} {
		_, err := parseCompetitionID(value)
		require.Error(t, err)
	}
	organizerID, err := parseResourceID[domain.OrganizerID]("123")
	require.NoError(t, err)
	assert.Equal(t, domain.OrganizerID(123), organizerID)
}
