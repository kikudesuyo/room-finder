package runner

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/kikudesuyo/room-finder/api/agent"
)

type APISaver struct {
	BaseURL    string
	HTTPClient *http.Client
}

func (s APISaver) Save(ctx context.Context, profile Profile, offer Offer, result agent.ExtractionResult) error {
	if profile.ID <= 0 {
		return fmt.Errorf("profile id must be positive")
	}
	payload := map[string]any{
		"source": offer.Source, "source_offer_id": offer.SourceOfferID, "source_url": offer.SourceURL,
		"name": result.Offer.Name, "address": result.Offer.Address, "rent_yen": result.Offer.RentYen,
		"management_fee_yen": result.Offer.ManagementFeeYen, "room_layout": result.Offer.RoomLayout,
		"area_square_meters": result.Offer.AreaSquareMeters, "built_year": result.Offer.BuiltYear,
		"floor": result.Offer.Floor, "access": result.Offer.Access, "amenities": result.Offer.Amenities,
		"details": map[string]any{"evidence": result.Evidence}, "captured_at": offer.CapturedAt,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	baseURL := strings.TrimRight(s.BaseURL, "/")
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/api/v1/search-profiles/%d/rental-offers", baseURL, profile.ID), bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := s.HTTPClient
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Do(req)
	if err != nil {
		return RetryError{Err: err}
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		message, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
		httpError := APIError{StatusCode: resp.StatusCode, Message: strings.TrimSpace(string(message))}
		if httpError.Retryable() {
			return RetryError{Err: httpError}
		}
		return httpError
	}
	return nil
}

type APIError struct {
	StatusCode int
	Message    string
}

func (e APIError) Error() string {
	return fmt.Sprintf("save API returned HTTP status %d: %s", e.StatusCode, e.Message)
}

func (e APIError) Retryable() bool {
	return e.StatusCode == http.StatusTooManyRequests || e.StatusCode >= 500
}
