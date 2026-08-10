package lifullhomes

import "time"

const Source = "lifullhomes"

type SearchRequest struct {
	URL string
}

type OfferReference struct {
	SourceOfferID string
	URL           string
}

type RentalOffer struct {
	Source           string
	SourceOfferID    string
	SourceURL        string
	Name             *string
	Address          *string
	RentYen          *int64
	ManagementFeeYen *int64
	RoomLayout       *string
	AreaSquareMeters *float64
	BuiltYear        *int
	Floor            *string
	Access           []Access
	Amenities        []string
	Details          map[string]any
	CapturedAt       time.Time
}

type Access struct {
	Line        string `json:"line,omitempty"`
	Station     string `json:"station,omitempty"`
	WalkMinutes *int   `json:"walk_minutes,omitempty"`
}

type CrawlError struct {
	URL       string
	Operation string
	Retryable bool
	Err       error
}

func (e *CrawlError) Error() string {
	return e.Operation + " " + e.URL + ": " + e.Err.Error()
}

func (e *CrawlError) Unwrap() error { return e.Err }
