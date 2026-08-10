package runner

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/kikudesuyo/room-finder/api/agent"
)

type Profile struct {
	ID                 int64
	InitialPrompt      string
	SearchURL          string
	RequiredConditions map[string]any
}

type Offer struct {
	Source        string    `json:"source"`
	SourceOfferID string    `json:"source_offer_id"`
	SourceURL     string    `json:"source_url"`
	SourceText    string    `json:"source_text"`
	CapturedAt    time.Time `json:"captured_at"`
}

type CrawlFunc func(context.Context, Profile) ([]Offer, error)
type ExtractFunc func(context.Context, Profile, Offer) (agent.ExtractionResult, error)
type SaveFunc func(context.Context, Profile, Offer, agent.ExtractionResult) error

type Retryable interface{ Retryable() bool }

type RetryError struct {
	Err error
}

func (e RetryError) Error() string   { return e.Err.Error() }
func (e RetryError) Unwrap() error   { return e.Err }
func (e RetryError) Retryable() bool { return true }

type state struct {
	LastRunAt time.Time     `json:"last_run_at"`
	Pending   []pendingTask `json:"pending"`
}

type pendingTask struct {
	Offer    Offer     `json:"offer"`
	Attempts int       `json:"attempts"`
	NextRun  time.Time `json:"next_run"`
}

type StateStore interface {
	Load(context.Context) (state, error)
	Save(context.Context, state) error
}

type FileStateStore struct {
	Path string
	mu   sync.Mutex
}

func (s *FileStateStore) Load(_ context.Context) (state, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, err := os.ReadFile(s.Path)
	if errors.Is(err, os.ErrNotExist) {
		return state{}, nil
	}
	if err != nil {
		return state{}, err
	}
	var value state
	if err := json.Unmarshal(data, &value); err != nil {
		return state{}, fmt.Errorf("decode runner state: %w", err)
	}
	return value, nil
}

func (s *FileStateStore) Save(_ context.Context, value state) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := os.MkdirAll(filepath.Dir(s.Path), 0o700); err != nil {
		return err
	}
	temporary := s.Path + ".tmp"
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(temporary, data, 0o600); err != nil {
		return err
	}
	return os.Rename(temporary, s.Path)
}

type Runner struct {
	Crawl       CrawlFunc
	Extract     ExtractFunc
	Save        SaveFunc
	State       StateStore
	MaxAttempts int
	Backoff     time.Duration
	Now         func() time.Time
}

type Report struct {
	Skipped       bool
	Processed     int
	Saved         int
	Rejected      int
	Retained      int
	RetryableErrs int
}

func (r Runner) RunIfDue(ctx context.Context, profile Profile) (Report, error) {
	now := r.now()
	current, err := r.State.Load(ctx)
	if err != nil {
		return Report{}, err
	}
	if !current.LastRunAt.IsZero() && sameUTCDate(current.LastRunAt, now) {
		return Report{Skipped: true}, nil
	}
	report, err := r.run(ctx, profile, current, now)
	if err != nil {
		return report, err
	}
	current.LastRunAt = now.UTC()
	if err := r.State.Save(ctx, current); err != nil {
		return report, err
	}
	return report, nil
}

func (r Runner) run(ctx context.Context, profile Profile, current state, now time.Time) (Report, error) {
	if r.Crawl == nil || r.Extract == nil || r.Save == nil || r.State == nil {
		return Report{}, errors.New("runner dependencies are required")
	}
	maxAttempts := r.MaxAttempts
	if maxAttempts <= 0 {
		maxAttempts = 3
	}
	tasks := make([]pendingTask, 0, len(current.Pending))
	for _, task := range current.Pending {
		if !task.NextRun.After(now) {
			tasks = append(tasks, task)
		}
	}
	newOffers, err := r.Crawl(ctx, profile)
	if err != nil {
		if isRetryable(err) {
			for _, offer := range newOffers {
				tasks = append(tasks, pendingTask{Offer: offer, NextRun: now.Add(r.Backoff)})
			}
		}
		return Report{Retained: len(tasks), RetryableErrs: boolInt(isRetryable(err))}, err
	}
	for _, offer := range newOffers {
		if !containsOffer(tasks, offer) {
			tasks = append(tasks, pendingTask{Offer: offer})
		}
	}

	nextPending := make([]pendingTask, 0, len(tasks))
	seen := make(map[string]struct{})
	report := Report{}
	for _, task := range tasks {
		key := task.Offer.Source + "\x00" + task.Offer.SourceOfferID
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		report.Processed++
		extracted, extractErr := r.Extract(ctx, profile, task.Offer)
		if extractErr != nil {
			if isRetryable(extractErr) {
				nextPending = append(nextPending, retryTask(task, now, maxAttempts, r.Backoff))
				report.RetryableErrs++
			}
			continue
		}
		conditions, conditionErr := agent.EvaluateRequiredConditions(extracted.Offer, profile.RequiredConditions)
		if conditionErr != nil || !agent.AllConditionsMatched(conditions) {
			report.Rejected++
			continue
		}
		if saveErr := r.Save(ctx, profile, task.Offer, extracted); saveErr != nil {
			if isRetryable(saveErr) {
				nextPending = append(nextPending, retryTask(task, now, maxAttempts, r.Backoff))
				report.RetryableErrs++
			}
			continue
		}
		report.Saved++
	}
	current.Pending = nextPending
	report.Retained = len(nextPending)
	return report, nil
}

func retryTask(task pendingTask, now time.Time, maxAttempts int, backoff time.Duration) pendingTask {
	task.Attempts++
	if task.Attempts >= maxAttempts {
		task.Attempts = 0
	}
	if backoff <= 0 {
		backoff = time.Minute
	}
	task.NextRun = now.Add(backoff)
	return task
}

func (r Runner) now() time.Time {
	if r.Now != nil {
		return r.Now()
	}
	return time.Now()
}

func isRetryable(err error) bool {
	var retryable Retryable
	return errors.As(err, &retryable) && retryable.Retryable()
}

func sameUTCDate(a, b time.Time) bool {
	a, b = a.UTC(), b.UTC()
	return a.Year() == b.Year() && a.YearDay() == b.YearDay()
}

func containsOffer(tasks []pendingTask, offer Offer) bool {
	for _, task := range tasks {
		if task.Offer.Source == offer.Source && task.Offer.SourceOfferID == offer.SourceOfferID {
			return true
		}
	}
	return false
}

func boolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}
