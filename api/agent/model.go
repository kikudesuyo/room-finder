package agent

import (
	"errors"
	"fmt"
	"strings"
)

type ExtractRequest struct {
	InitialPrompt string
	SourceURL     string
	SourceText    string
}

type ExtractedOffer struct {
	Name             *string  `json:"name"`
	Address          *string  `json:"address"`
	RentYen          *int64   `json:"rent_yen"`
	ManagementFeeYen *int64   `json:"management_fee_yen"`
	RoomLayout       *string  `json:"room_layout"`
	AreaSquareMeters *float64 `json:"area_square_meters"`
	BuiltYear        *int     `json:"built_year"`
	Floor            *string  `json:"floor"`
	Access           []Access `json:"access"`
	Amenities        []string `json:"amenities"`
}

type Access struct {
	Line        string `json:"line"`
	Station     string `json:"station"`
	WalkMinutes *int   `json:"walk_minutes"`
}

type Evidence struct {
	Field  string `json:"field"`
	Value  string `json:"value"`
	Source string `json:"source"`
}

type ExtractionResult struct {
	Offer    ExtractedOffer `json:"offer"`
	Evidence []Evidence     `json:"evidence"`
}

type ConditionResult struct {
	Condition string
	Matched   bool
	Reason    string
}

func EvaluateRequiredConditions(offer ExtractedOffer, conditions map[string]any) ([]ConditionResult, error) {
	results := make([]ConditionResult, 0, len(conditions))
	for name, value := range conditions {
		result, err := evaluateCondition(offer, name, value)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

func AllConditionsMatched(results []ConditionResult) bool {
	for _, result := range results {
		if !result.Matched {
			return false
		}
	}
	return true
}

func evaluateCondition(offer ExtractedOffer, name string, value any) (ConditionResult, error) {
	result := ConditionResult{Condition: name}
	switch name {
	case "max_rent_yen":
		limit, ok := number(value)
		if !ok {
			return result, fmt.Errorf("%s must be a number", name)
		}
		result.Matched = offer.RentYen != nil && *offer.RentYen <= limit
		result.Reason = compareReason(offer.RentYen, limit, "rent_yen")
	case "max_management_fee_yen":
		limit, ok := number(value)
		if !ok {
			return result, fmt.Errorf("%s must be a number", name)
		}
		result.Matched = offer.ManagementFeeYen != nil && *offer.ManagementFeeYen <= limit
		result.Reason = compareReason(offer.ManagementFeeYen, limit, "management_fee_yen")
	case "min_area_square_meters":
		limit, ok := floatNumber(value)
		if !ok {
			return result, fmt.Errorf("%s must be a number", name)
		}
		result.Matched = offer.AreaSquareMeters != nil && *offer.AreaSquareMeters >= limit
		result.Reason = compareFloatReason(offer.AreaSquareMeters, limit, "area_square_meters")
	case "room_layout":
		wanted, ok := value.(string)
		if !ok || strings.TrimSpace(wanted) == "" {
			return result, errors.New("room_layout must be a non-empty string")
		}
		result.Matched = offer.RoomLayout != nil && *offer.RoomLayout == wanted
		result.Reason = fmt.Sprintf("room_layout=%q, required=%q", stringValue(offer.RoomLayout), wanted)
	case "must_have_amenities":
		wanted, ok := stringSlice(value)
		if !ok {
			return result, errors.New("must_have_amenities must be an array of strings")
		}
		missing := make([]string, 0)
		for _, item := range wanted {
			if !contains(offer.Amenities, item) {
				missing = append(missing, item)
			}
		}
		result.Matched = len(missing) == 0
		result.Reason = fmt.Sprintf("missing_amenities=%v", missing)
	case "max_walk_minutes":
		limit, ok := number(value)
		if !ok {
			return result, fmt.Errorf("%s must be a number", name)
		}
		result.Matched = false
		for _, access := range offer.Access {
			if access.WalkMinutes != nil && int64(*access.WalkMinutes) <= limit {
				result.Matched = true
				break
			}
		}
		result.Reason = fmt.Sprintf("max_walk_minutes=%d", limit)
	default:
		return result, fmt.Errorf("unsupported required condition: %s", name)
	}
	return result, nil
}

func number(value any) (int64, bool) {
	switch value := value.(type) {
	case float64:
		return int64(value), true
	case int:
		return int64(value), true
	case int64:
		return value, true
	default:
		return 0, false
	}
}

func floatNumber(value any) (float64, bool) {
	switch value := value.(type) {
	case float64:
		return value, true
	case int:
		return float64(value), true
	case int64:
		return float64(value), true
	default:
		return 0, false
	}
}

func stringSlice(value any) ([]string, bool) {
	values, ok := value.([]any)
	if !ok {
		if strings, ok := value.([]string); ok {
			return strings, true
		}
		return nil, false
	}
	result := make([]string, 0, len(values))
	for _, item := range values {
		stringValue, ok := item.(string)
		if !ok {
			return nil, false
		}
		result = append(result, stringValue)
	}
	return result, true
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func compareReason(value *int64, limit int64, name string) string {
	if value == nil {
		return name + " is unknown"
	}
	return fmt.Sprintf("%s=%d, limit=%d", name, *value, limit)
}

func compareFloatReason(value *float64, limit float64, name string) string {
	if value == nil {
		return name + " is unknown"
	}
	return fmt.Sprintf("%s=%.2f, limit=%.2f", name, *value, limit)
}

func stringValue(value *string) string {
	if value == nil {
		return "unknown"
	}
	return *value
}
