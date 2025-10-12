package contracts

import "time"

// Bid represents a single driver offer for a rider.
type Bid struct {
	ID        string
	DriverID  string
	TripID    string
	Price     float64
	Latitude  float64
	Longitude float64
	ETA       time.Duration
	ExpiresAt time.Time
}

// BidRequest collects context required to evaluate bids for a rider.
type BidRequest struct {
	RiderID   string
	TripID    string
	Latitude  float64
	Longitude float64
	MaxETA    time.Duration
	MaxPrice  float64
}

// AcceptedBid records the winning offer for persistence.
type AcceptedBid struct {
	BidID      string
	TripID     string
	DriverID   string
	Price      float64
	AcceptedAt time.Time
}

// TripState captures the lifecycle stage of a trip.
type TripState string

const (
	TripStatePending  TripState = "pending"
	TripStateActive   TripState = "active"
	TripStatePaused   TripState = "paused"
	TripStateCanceled TripState = "canceled"
	TripStateComplete TripState = "complete"
)

// TripEvent describes a state change event for auditing.
type TripEvent struct {
	TripID     string
	State      TripState
	OccurredAt time.Time
	Notes      string
}
