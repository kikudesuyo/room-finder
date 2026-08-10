package agent

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEvaluateRequiredConditionsDoesNotInferMissingValues(t *testing.T) {
	results, err := EvaluateRequiredConditions(ExtractedOffer{}, map[string]any{"max_rent_yen": float64(100000)})
	if err != nil {
		t.Fatal(err)
	}
	if AllConditionsMatched(results) {
		t.Fatal("missing rent was treated as a match")
	}
}

func TestOpenAIClientUsesStrictStructuredOutput(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request map[string]any
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatal(err)
		}
		format := request["text"].(map[string]any)["format"].(map[string]any)
		if format["type"] != "json_schema" || format["strict"] != true {
			t.Fatalf("format = %+v", format)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"output":[{"content":[{"text":"{\"offer\":{\"name\":\"コーポ藤沢\",\"address\":null,\"rent_yen\":80000,\"management_fee_yen\":null,\"room_layout\":\"1LDK\",\"area_square_meters\":null,\"built_year\":null,\"floor\":null,\"access\":[],\"amenities\":[]},\"evidence\":[{\"field\":\"rent_yen\",\"value\":\"8万円\",\"source\":\"家賃欄\"}]}"}]}]}`))
	}))
	client, err := NewOpenAIClient("test-key", "test-model", server.Client())
	if err != nil {
		t.Fatal(err)
	}
	client.endpoint = server.URL
	result, err := client.Extract(context.Background(), ExtractRequest{InitialPrompt: "家賃10万円以下", SourceURL: "https://www.homes.co.jp/chintai/123", SourceText: "家賃 8万円"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Offer.Name == nil || *result.Offer.Name != "コーポ藤沢" || result.Offer.RentYen == nil || *result.Offer.RentYen != 80000 {
		t.Fatalf("result = %+v", result)
	}
}
