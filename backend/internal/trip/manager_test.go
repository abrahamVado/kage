package trip

import (
	"context"
	"testing"
	"time"

	"kage/backend/internal/contracts"
)

type recordingRepo struct {
	events []contracts.TripEvent
}

func (r *recordingRepo) RecordEvent(ctx context.Context, event contracts.TripEvent) error {
	r.events = append(r.events, event)
	return nil
}

type sequenceClock struct {
	times []time.Time
	idx   int
}

func (s *sequenceClock) Now() time.Time {
	if s.idx >= len(s.times) {
		return s.times[len(s.times)-1]
	}
	t := s.times[s.idx]
	s.idx++
	return t
}

func TestTripLifecycle(t *testing.T) {
	times := []time.Time{
		time.Unix(0, 0),   // start
		time.Unix(1, 0),   // start event
		time.Unix(60, 0),  // pause now
		time.Unix(61, 0),  // pause event
		time.Unix(120, 0), // resume now
		time.Unix(121, 0), // resume event
		time.Unix(240, 0), // complete now
		time.Unix(241, 0), // complete event
	}
	clock := &sequenceClock{times: times}
	repo := &recordingRepo{}
	mgr := NewManager(repo, clock)

	if err := mgr.StartTrip(context.Background(), "trip-1"); err != nil {
		t.Fatalf("start failed: %v", err)
	}
	if err := mgr.PauseTrip(context.Background(), "trip-1"); err != nil {
		t.Fatalf("pause failed: %v", err)
	}
	if err := mgr.ResumeTrip(context.Background(), "trip-1"); err != nil {
		t.Fatalf("resume failed: %v", err)
	}
	if err := mgr.CompleteTrip(context.Background(), "trip-1"); err != nil {
		t.Fatalf("complete failed: %v", err)
	}

	metrics, ok := mgr.MetricsFor("trip-1")
	if !ok {
		t.Fatalf("expected metrics")
	}
	if metrics.TotalActive != 180*time.Second {
		t.Fatalf("unexpected active duration: %v", metrics.TotalActive)
	}
	if metrics.TotalPaused != 60*time.Second {
		t.Fatalf("unexpected paused duration: %v", metrics.TotalPaused)
	}
	if len(repo.events) != 4 {
		t.Fatalf("expected 4 events recorded got %d", len(repo.events))
	}
}

func TestInvalidTransitions(t *testing.T) {
	mgr := NewManager(nil, &sequenceClock{times: []time.Time{time.Now(), time.Now()}})
	if err := mgr.PauseTrip(context.Background(), "missing"); err != ErrInvalidTransition {
		t.Fatalf("expected invalid transition error")
	}
	if err := mgr.CompleteTrip(context.Background(), "missing"); err != ErrInvalidTransition {
		t.Fatalf("expected invalid transition error")
	}
}
