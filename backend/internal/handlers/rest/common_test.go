package rest

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
