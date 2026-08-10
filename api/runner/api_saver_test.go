package runner

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kikudesuyo/room-finder/api/agent"
)

func TestAPISaverUsesRentalOfferEndpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/search-profiles/7/rental-offers" {
			t.Fatalf("request = %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	err := (APISaver{BaseURL: server.URL}).Save(context.Background(), Profile{ID: 7}, Offer{Source: "lifullhomes", SourceOfferID: "123", SourceURL: "https://www.homes.co.jp/chintai/123", CapturedAt: time.Now()}, agent.ExtractionResult{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestAPIErrorRetryClassification(t *testing.T) {
	if !(APIError{StatusCode: http.StatusTooManyRequests}).Retryable() || !(APIError{StatusCode: http.StatusBadGateway}).Retryable() || (APIError{StatusCode: http.StatusBadRequest}).Retryable() {
		t.Fatal("unexpected API retry classification")
	}
}
