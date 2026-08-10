package route

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUpdateSearchProfileIsRejected(t *testing.T) {
	req := httptest.NewRequest(http.MethodPut, "/api/v1/search-profiles/1", strings.NewReader(`{"initial_prompt":"updated","required_conditions":{}}`))
	recorder := httptest.NewRecorder()

	NewRouter().ServeHTTP(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("PUT /search-profiles/{id} status = %d, want %d", recorder.Code, http.StatusBadRequest)
	}
}
