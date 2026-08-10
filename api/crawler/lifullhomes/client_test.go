package lifullhomes

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestSearchAndFetchOffer(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/robots.txt":
			w.WriteHeader(http.StatusNotFound)
		case "/chintai/123":
			w.Header().Set("Content-Type", "text/html")
			_, _ = w.Write([]byte(`<html><body><script type="application/ld+json">{"name":"コーポ藤沢","address":"神奈川県藤沢市","offers":{"price":80000},"numberOfRooms":"1LDK","floorSize":{"value":35.5}}</script></body></html>`))
		case "/search":
			w.Header().Set("Content-Type", "text/html")
			_, _ = w.Write([]byte(`<html><body><a href="/chintai/123">詳細</a><a href="https://example.com/chintai/999">外部</a></body></html>`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	// The adapter only accepts real LIFULL hosts; rewrite the test server request
	// host while keeping its TLS transport to test parsing without the network.
	client := server.Client()
	adapter := NewClient(client, "test", 0)
	originalHosts := allowedHosts
	allowedHosts = map[string]bool{"127.0.0.1": true}
	defer func() { allowedHosts = originalHosts }()

	refs, err := adapter.Search(context.Background(), SearchRequest{URL: server.URL + "/search"})
	if err != nil {
		t.Fatal(err)
	}
	if len(refs) != 1 || refs[0].SourceOfferID != "123" {
		t.Fatalf("Search() refs = %+v", refs)
	}
	offer, err := adapter.FetchOffer(context.Background(), refs[0])
	if err != nil {
		t.Fatal(err)
	}
	if offer.Name == nil || *offer.Name != "コーポ藤沢" || offer.RentYen == nil || *offer.RentYen != 80000 {
		t.Fatalf("FetchOffer() offer = %+v", offer)
	}
	if offer.CapturedAt.Before(time.Now().Add(-time.Minute)) {
		t.Fatal("FetchOffer() did not capture the current time")
	}
}

func TestRobotsDisallow(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/robots.txt" {
			_, _ = w.Write([]byte("User-agent: *\nDisallow: /chintai/\n"))
			return
		}
		http.Error(w, "should not fetch", http.StatusInternalServerError)
	}))
	defer server.Close()
	client := NewClient(server.Client(), "test", 0)
	originalHosts := allowedHosts
	allowedHosts = map[string]bool{"127.0.0.1": true}
	defer func() { allowedHosts = originalHosts }()

	_, err := client.FetchOffer(context.Background(), OfferReference{SourceOfferID: "123", URL: server.URL + "/chintai/123"})
	if err == nil || !strings.Contains(err.Error(), "disallowed") {
		t.Fatalf("FetchOffer() error = %v", err)
	}
}
