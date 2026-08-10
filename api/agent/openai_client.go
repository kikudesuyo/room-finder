package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type OpenAIClient struct {
	apiKey     string
	model      string
	httpClient *http.Client
	endpoint   string
}

func NewOpenAIClient(apiKey, model string, httpClient *http.Client) (*OpenAIClient, error) {
	if strings.TrimSpace(apiKey) == "" {
		return nil, errors.New("OPENAI_API_KEY is required")
	}
	if strings.TrimSpace(model) == "" {
		return nil, errors.New("OPENAI_MODEL is required")
	}
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	return &OpenAIClient{apiKey: apiKey, model: model, httpClient: httpClient, endpoint: "https://api.openai.com/v1/responses"}, nil
}

func (c *OpenAIClient) Extract(ctx context.Context, req ExtractRequest) (ExtractionResult, error) {
	body := map[string]any{
		"model": c.model,
		"store": false,
		"input": []map[string]any{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": UserPrompt(req)},
		},
		"text": map[string]any{"format": map[string]any{
			"type": "json_schema", "name": "rental_offer_extraction", "strict": true, "schema": extractionSchema(),
		}},
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return ExtractionResult{}, err
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(payload))
	if err != nil {
		return ExtractionResult{}, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return ExtractionResult{}, err
	}
	defer func() { _ = resp.Body.Close() }()
	responseBody, err := io.ReadAll(io.LimitReader(resp.Body, 10<<20))
	if err != nil {
		return ExtractionResult{}, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return ExtractionResult{}, fmt.Errorf("OpenAI API returned HTTP status %d: %s", resp.StatusCode, strings.TrimSpace(string(responseBody)))
	}
	var response struct {
		Output []struct {
			Content []struct {
				Text string `json:"text"`
			} `json:"content"`
		} `json:"output"`
	}
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return ExtractionResult{}, err
	}
	for _, item := range response.Output {
		for _, content := range item.Content {
			if strings.TrimSpace(content.Text) == "" {
				continue
			}
			var result ExtractionResult
			if err := json.Unmarshal([]byte(content.Text), &result); err != nil {
				return ExtractionResult{}, fmt.Errorf("structured output is invalid: %w", err)
			}
			return result, nil
		}
	}
	return ExtractionResult{}, errors.New("OpenAI response did not contain output text")
}

func extractionSchema() map[string]any {
	nullable := func(kind string) map[string]any { return map[string]any{"type": []string{kind, "null"}} }
	return map[string]any{
		"type": "object", "additionalProperties": false,
		"required": []string{"offer", "evidence"},
		"properties": map[string]any{
			"offer": map[string]any{
				"type": "object", "additionalProperties": false,
				"required": []string{"name", "address", "rent_yen", "management_fee_yen", "room_layout", "area_square_meters", "built_year", "floor", "access", "amenities"},
				"properties": map[string]any{
					"name": nullable("string"), "address": nullable("string"), "rent_yen": nullable("integer"), "management_fee_yen": nullable("integer"),
					"room_layout": nullable("string"), "area_square_meters": nullable("number"), "built_year": nullable("integer"), "floor": nullable("string"),
					"access":    map[string]any{"type": "array", "items": map[string]any{"type": "object", "additionalProperties": false, "required": []string{"line", "station", "walk_minutes"}, "properties": map[string]any{"line": map[string]any{"type": "string"}, "station": map[string]any{"type": "string"}, "walk_minutes": nullable("integer")}}},
					"amenities": map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
				},
			},
			"evidence": map[string]any{"type": "array", "items": map[string]any{"type": "object", "additionalProperties": false, "required": []string{"field", "value", "source"}, "properties": map[string]any{"field": map[string]any{"type": "string"}, "value": map[string]any{"type": "string"}, "source": map[string]any{"type": "string"}}}},
		},
	}
}
