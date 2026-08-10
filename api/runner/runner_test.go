package runner

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/kikudesuyo/room-finder/api/agent"
)

func TestRunnerRejectsUnmatchedOffersAndSavesMatchedOffers(t *testing.T) {
	store := &memoryStore{}
	saved := 0
	runner := Runner{
		State: store, MaxAttempts: 2, Now: func() time.Time { return time.Date(2026, 8, 11, 1, 0, 0, 0, time.UTC) },
		Crawl: func(context.Context, Profile) ([]Offer, error) {
			return []Offer{{Source: "lifullhomes", SourceOfferID: "1"}, {Source: "lifullhomes", SourceOfferID: "2"}}, nil
		},
		Extract: func(context.Context, Profile, Offer) (agent.ExtractionResult, error) {
			return agent.ExtractionResult{Offer: agent.ExtractedOffer{RentYen: pointer(int64(80000))}}, nil
		},
		Save: func(context.Context, Profile, Offer, agent.ExtractionResult) error { saved++; return nil },
	}
	report, err := runner.RunIfDue(context.Background(), Profile{ID: 1, RequiredConditions: map[string]any{"max_rent_yen": float64(100000)}})
	if err != nil || report.Saved != 2 || saved != 2 {
		t.Fatalf("report=%+v saved=%d err=%v", report, saved, err)
	}
	if _, err := runner.RunIfDue(context.Background(), Profile{ID: 1}); err != nil || !runnerReportSkipped(runner, Profile{ID: 1}) {
		t.Fatalf("daily run was not skipped: err=%v", err)
	}
}

func TestRunnerRetainsRetryableFailure(t *testing.T) {
	store := &memoryStore{}
	runner := Runner{
		State: store, MaxAttempts: 1, Backoff: time.Hour, Now: func() time.Time { return time.Date(2026, 8, 11, 1, 0, 0, 0, time.UTC) },
		Crawl: func(context.Context, Profile) ([]Offer, error) {
			return []Offer{{Source: "lifullhomes", SourceOfferID: "1"}}, nil
		},
		Extract: func(context.Context, Profile, Offer) (agent.ExtractionResult, error) {
			return agent.ExtractionResult{}, RetryError{Err: errors.New("timeout")}
		},
		Save: func(context.Context, Profile, Offer, agent.ExtractionResult) error { return nil },
	}
	report, err := runner.RunIfDue(context.Background(), Profile{ID: 1})
	if err != nil || report.Retained != 1 {
		t.Fatalf("report=%+v err=%v", report, err)
	}
}

func TestFileStateStorePersistsState(t *testing.T) {
	store := &FileStateStore{Path: filepath.Join(t.TempDir(), "runner.json")}
	want := state{LastRunAt: time.Date(2026, 8, 11, 1, 0, 0, 0, time.UTC), Pending: []pendingTask{{Offer: Offer{SourceOfferID: "1"}}}}
	if err := store.Save(context.Background(), want); err != nil {
		t.Fatal(err)
	}
	got, err := store.Load(context.Background())
	if err != nil || len(got.Pending) != 1 || got.LastRunAt.IsZero() {
		t.Fatalf("state=%+v err=%v", got, err)
	}
}

type memoryStore struct{ value state }

func (s *memoryStore) Load(context.Context) (state, error)       { return s.value, nil }
func (s *memoryStore) Save(_ context.Context, value state) error { s.value = value; return nil }

func pointer(value int64) *int64 { return &value }

func runnerReportSkipped(r Runner, profile Profile) bool {
	report, err := r.RunIfDue(context.Background(), profile)
	return err == nil && report.Skipped
}
