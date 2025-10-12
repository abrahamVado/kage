package trip

import (
	"context"
	"errors"
	"sync"
	"time"

	"kage/backend/internal/contracts"
)

// ErrInvalidTransition occurs when a lifecycle rule is violated.
var ErrInvalidTransition = errors.New("invalid trip state transition")

// Clock abstracts time for deterministic unit tests.
type Clock interface {
	Now() time.Time
}

// RealClock delegates to time.Now.
type RealClock struct{}

// Now returns the current time.
func (RealClock) Now() time.Time { return time.Now() }

type tripState struct {
	state       contracts.TripState
	startedAt   time.Time
	lastResumed time.Time
	lastPaused  time.Time
	totalActive time.Duration
	totalPaused time.Duration
}

// Manager coordinates trip lifecycle transitions.
type Manager struct {
	mu    sync.RWMutex
	repo  EventRepository
	clock Clock
	trips map[string]*tripState
}

// NewManager constructs a Manager with the provided repository.
func NewManager(repo EventRepository, clock Clock) *Manager {
	if clock == nil {
		clock = RealClock{}
	}
	return &Manager{
		repo:  repo,
		clock: clock,
		trips: make(map[string]*tripState),
	}
}

// StartTrip marks the trip as active and records an event.
func (m *Manager) StartTrip(ctx context.Context, tripID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.trips[tripID]; exists {
		return ErrInvalidTransition
	}
	now := m.clock.Now()
	st := &tripState{state: contracts.TripStateActive, startedAt: now, lastResumed: now}
	m.trips[tripID] = st
	return m.persistEvent(ctx, tripID, contracts.TripStateActive, "trip started")
}

// PauseTrip transitions an active trip into the paused state.
func (m *Manager) PauseTrip(ctx context.Context, tripID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	st, ok := m.trips[tripID]
	if !ok || st.state != contracts.TripStateActive {
		return ErrInvalidTransition
	}
	now := m.clock.Now()
	st.totalActive += now.Sub(st.lastResumed)
	st.lastPaused = now
	st.state = contracts.TripStatePaused
	return m.persistEvent(ctx, tripID, contracts.TripStatePaused, "trip paused")
}

// ResumeTrip moves a paused trip back to active.
func (m *Manager) ResumeTrip(ctx context.Context, tripID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	st, ok := m.trips[tripID]
	if !ok || st.state != contracts.TripStatePaused {
		return ErrInvalidTransition
	}
	now := m.clock.Now()
	st.totalPaused += now.Sub(st.lastPaused)
	st.lastResumed = now
	st.state = contracts.TripStateActive
	return m.persistEvent(ctx, tripID, contracts.TripStateActive, "trip resumed")
}

// CancelTrip stops the trip permanently.
func (m *Manager) CancelTrip(ctx context.Context, tripID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	st, ok := m.trips[tripID]
	if !ok || st.state == contracts.TripStateComplete || st.state == contracts.TripStateCanceled {
		return ErrInvalidTransition
	}
	now := m.clock.Now()
	if st.state == contracts.TripStateActive {
		st.totalActive += now.Sub(st.lastResumed)
	}
	if st.state == contracts.TripStatePaused {
		st.totalPaused += now.Sub(st.lastPaused)
	}
	st.state = contracts.TripStateCanceled
	return m.persistEvent(ctx, tripID, contracts.TripStateCanceled, "trip canceled")
}

// CompleteTrip finalizes the trip and stores a completion event.
func (m *Manager) CompleteTrip(ctx context.Context, tripID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	st, ok := m.trips[tripID]
	if !ok || st.state != contracts.TripStateActive {
		return ErrInvalidTransition
	}
	now := m.clock.Now()
	st.totalActive += now.Sub(st.lastResumed)
	st.state = contracts.TripStateComplete
	return m.persistEvent(ctx, tripID, contracts.TripStateComplete, "trip completed")
}

// Metrics describes durations for auditing.
type Metrics struct {
	TotalActive time.Duration
	TotalPaused time.Duration
	StartedAt   time.Time
}

// MetricsFor retrieves aggregated trip timing metrics.
func (m *Manager) MetricsFor(tripID string) (Metrics, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	st, ok := m.trips[tripID]
	if !ok {
		return Metrics{}, false
	}
	return Metrics{
		TotalActive: st.totalActive,
		TotalPaused: st.totalPaused,
		StartedAt:   st.startedAt,
	}, true
}

func (m *Manager) persistEvent(ctx context.Context, tripID string, state contracts.TripState, notes string) error {
	if m.repo == nil {
		return nil
	}
	event := contracts.TripEvent{TripID: tripID, State: state, OccurredAt: m.clock.Now(), Notes: notes}
	return m.repo.RecordEvent(ctx, event)
}
