package bidding

import (
	"context"
	"errors"
	"sort"
	"time"

	"kage/backend/internal/contracts"
	"kage/backend/internal/geo"
)

// ErrEvaluationTimeout occurs when ranking exceeds the configured deadline.
var ErrEvaluationTimeout = errors.New("bid evaluation timeout")

// Clock abstracts time for deterministic testing.
type Clock interface {
	Now() time.Time
	After(d time.Duration) <-chan time.Time
}

// RealClock delegates to the time package.
type RealClock struct{}

// Now returns the current time.
func (RealClock) Now() time.Time { return time.Now() }

// After defers to time.After.
func (RealClock) After(d time.Duration) <-chan time.Time { return time.After(d) }

// Arbiter evaluates bids and persists the winner.
type Arbiter struct {
	repo              Repository
	clock             Clock
	evaluationTimeout time.Duration
	radiusKm          float64
}

// Option mutates Arbiter configuration.
type Option func(*Arbiter)

// WithTimeout configures a custom evaluation timeout.
func WithTimeout(d time.Duration) Option {
	return func(a *Arbiter) { a.evaluationTimeout = d }
}

// WithRadius configures the proximity radius filter.
func WithRadius(km float64) Option {
	return func(a *Arbiter) { a.radiusKm = km }
}

// WithClock injects a custom clock for tests.
func WithClock(clock Clock) Option {
	return func(a *Arbiter) { a.clock = clock }
}

// NewArbiter builds the orchestrator with sane defaults.
func NewArbiter(repo Repository, opts ...Option) *Arbiter {
	a := &Arbiter{
		repo:              repo,
		clock:             RealClock{},
		evaluationTimeout: 3 * time.Second,
		radiusKm:          5,
	}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

// RankAndSelect picks the optimal bid and persists it using the repository.
func (a *Arbiter) RankAndSelect(ctx context.Context, req contracts.BidRequest, bids []contracts.Bid) (contracts.Bid, bool, error) {
	resCh := make(chan struct {
		bid contracts.Bid
		ok  bool
		err error
	}, 1)

	go func() {
		// 1.- Filter bids by freshness, budget, travel time, and radius.
		candidates := a.filterBids(req, bids)
		if len(candidates) == 0 {
			resCh <- struct {
				bid contracts.Bid
				ok  bool
				err error
			}{bid: contracts.Bid{}, ok: false, err: nil}
			return
		}

		// 2.- Score the remaining bids, persist the winner, and emit the result.
		winner := a.rankCandidates(req, candidates)
		if a.repo != nil {
			err := a.repo.SaveAcceptedBid(ctx, contracts.AcceptedBid{
				BidID:      winner.ID,
				TripID:     winner.TripID,
				DriverID:   winner.DriverID,
				Price:      winner.Price,
				AcceptedAt: a.clock.Now(),
			})
			if err != nil {
				resCh <- struct {
					bid contracts.Bid
					ok  bool
					err error
				}{bid: contracts.Bid{}, ok: false, err: err}
				return
			}
		}
		resCh <- struct {
			bid contracts.Bid
			ok  bool
			err error
		}{bid: winner, ok: true, err: nil}
	}()

	var timeout <-chan time.Time
	if a.evaluationTimeout > 0 {
		timeout = a.clock.After(a.evaluationTimeout)
	}

	select {
	case <-ctx.Done():
		return contracts.Bid{}, false, ctx.Err()
	case <-timeout:
		return contracts.Bid{}, false, ErrEvaluationTimeout
	case res := <-resCh:
		return res.bid, res.ok, res.err
	}
}

func (a *Arbiter) filterBids(req contracts.BidRequest, bids []contracts.Bid) []contracts.Bid {
	now := a.clock.Now()
	var filtered []contracts.Bid
	for _, bid := range bids {
		if !bid.ExpiresAt.IsZero() && bid.ExpiresAt.Before(now) {
			continue
		}
		if req.MaxPrice > 0 && bid.Price > req.MaxPrice {
			continue
		}
		if req.MaxETA > 0 && bid.ETA > req.MaxETA {
			continue
		}
		if a.radiusKm > 0 && !geo.WithinRadius(req.Latitude, req.Longitude, bid.Latitude, bid.Longitude, a.radiusKm) {
			continue
		}
		filtered = append(filtered, bid)
	}
	return filtered
}

func (a *Arbiter) rankCandidates(req contracts.BidRequest, bids []contracts.Bid) contracts.Bid {
	type candidate struct {
		bid   contracts.Bid
		score float64
	}

	var candidates []candidate
	for _, bid := range bids {
		distance := geo.DistanceBetween(req.Latitude, req.Longitude, bid.Latitude, bid.Longitude)
		score := a.computeScore(req, bid, distance)
		candidates = append(candidates, candidate{bid: bid, score: score})
	}

	sort.SliceStable(candidates, func(i, j int) bool {
		return candidates[i].score > candidates[j].score
	})

	return candidates[0].bid
}

func (a *Arbiter) computeScore(req contracts.BidRequest, bid contracts.Bid, distance float64) float64 {
	priceComponent := 1.0
	if req.MaxPrice > 0 {
		priceComponent = (req.MaxPrice - bid.Price) / req.MaxPrice
	}
	timeComponent := 1.0
	if req.MaxETA > 0 {
		timeComponent = float64(req.MaxETA-bid.ETA) / float64(req.MaxETA)
	}
	if timeComponent < 0 {
		timeComponent = 0
	}
	if priceComponent < 0 {
		priceComponent = 0
	}
	proximityComponent := 1 / (1 + distance)
	return 0.45*priceComponent + 0.35*timeComponent + 0.2*proximityComponent
}
