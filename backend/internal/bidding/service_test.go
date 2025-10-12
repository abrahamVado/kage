package bidding

import (
	"context"
	"errors"
	"testing"
	"time"

	"kage/backend/internal/contracts"
)

type fakeRepo struct {
	saved []contracts.AcceptedBid
	err   error
}

func (f *fakeRepo) SaveAcceptedBid(ctx context.Context, bid contracts.AcceptedBid) error {
	if f.err != nil {
		return f.err
	}
	f.saved = append(f.saved, bid)
	return nil
}

type fakeClock struct {
	now     time.Time
	afterCh chan time.Time
}

func (f *fakeClock) Now() time.Time { return f.now }

func (f *fakeClock) After(d time.Duration) <-chan time.Time {
	if f.afterCh != nil {
		return f.afterCh
	}
	ch := make(chan time.Time)
	return ch
}

func TestRankAndSelect(t *testing.T) {
	fc := &fakeClock{now: time.Date(2023, 6, 1, 12, 0, 0, 0, time.UTC)}
	repo := &fakeRepo{}
	arbiter := NewArbiter(repo, WithClock(fc), WithRadius(10), WithTimeout(time.Minute))

	req := contracts.BidRequest{RiderID: "r1", TripID: "t1", Latitude: 0, Longitude: 0, MaxETA: 30 * time.Minute, MaxPrice: 50}
	bids := []contracts.Bid{
		{ID: "b1", DriverID: "d1", TripID: "t1", Price: 40, Latitude: 0.05, Longitude: 0.05, ETA: 20 * time.Minute, ExpiresAt: fc.now.Add(time.Hour)},
		{ID: "b2", DriverID: "d2", TripID: "t1", Price: 35, Latitude: 5, Longitude: 5, ETA: 15 * time.Minute, ExpiresAt: fc.now.Add(time.Hour)},
	}

	winner, ok, err := arbiter.RankAndSelect(context.Background(), req, bids)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatalf("expected a winner")
	}
	if winner.ID != "b1" {
		t.Fatalf("expected b1 to win got %s", winner.ID)
	}
	if len(repo.saved) != 1 {
		t.Fatalf("expected repository save")
	}
}

func TestRankAndSelectTimeout(t *testing.T) {
	fc := &fakeClock{now: time.Now(), afterCh: make(chan time.Time, 1)}
	fc.afterCh <- time.Now()
	arbiter := NewArbiter(&fakeRepo{}, WithClock(fc), WithTimeout(time.Millisecond))

	_, ok, err := arbiter.RankAndSelect(context.Background(), contracts.BidRequest{}, nil)
	if !errors.Is(err, ErrEvaluationTimeout) {
		t.Fatalf("expected timeout error got %v", err)
	}
	if ok {
		t.Fatalf("expected no winner due to timeout")
	}
}
