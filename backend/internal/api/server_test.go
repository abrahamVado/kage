package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"kage/backend/internal/auth"
	"kage/backend/internal/bidding"
	"kage/backend/internal/contracts"
	"kage/backend/internal/trip"
)

type captureRepo struct {
	saved []contracts.AcceptedBid
}

func (c *captureRepo) SaveAcceptedBid(_ context.Context, bid contracts.AcceptedBid) error {
	c.saved = append(c.saved, bid)
	return nil
}

func TestEvaluateBidsAuthorized(t *testing.T) {
	// 1.- Build the HTTP server with the real arbiter configured to persist the winner.
	gin.SetMode(gin.TestMode)
	repo := &captureRepo{}
	arbiter := bidding.NewArbiter(repo, bidding.WithTimeout(5*time.Second))
	manager := trip.NewManager(nil, nil)
	validator := auth.NewValidator("top-secret")
	server := NewServer(arbiter, manager, validator)
	router := gin.New()
	server.RegisterRoutes(router)

	payload := struct {
		Request contracts.BidRequest `json:"request"`
		Bids    []contracts.Bid      `json:"bids"`
	}{
		Request: contracts.BidRequest{TripID: "trip-1", RiderID: "rider-1", Latitude: 1.0, Longitude: 2.0, MaxPrice: 60, MaxETA: 45 * time.Minute},
		Bids: []contracts.Bid{
			{ID: "bid-1", TripID: "trip-1", DriverID: "driver-a", Price: 50, Latitude: 1.01, Longitude: 2.01, ETA: 30 * time.Minute, ExpiresAt: time.Now().Add(time.Hour)},
			{ID: "bid-2", TripID: "trip-1", DriverID: "driver-b", Price: 70, Latitude: 1.02, Longitude: 2.02, ETA: 35 * time.Minute, ExpiresAt: time.Now().Add(time.Hour)},
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	// 2.- Issue an authenticated request and verify that the response exposes the accepted bid.
	req := httptest.NewRequest(http.MethodPost, "/bids/evaluate", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer top-secret")
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", res.Code)
	}
	var envelope struct {
		Winner contracts.Bid `json:"winner"`
	}
	if err := json.Unmarshal(res.Body.Bytes(), &envelope); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if envelope.Winner.ID != "bid-1" {
		t.Fatalf("unexpected winner: %s", envelope.Winner.ID)
	}
	if len(repo.saved) != 1 {
		t.Fatalf("expected winner persisted, got %d saves", len(repo.saved))
	}
}

func TestEvaluateBidsUnauthorized(t *testing.T) {
	// 1.- Create the server with an authentication secret but omit the header.
	gin.SetMode(gin.TestMode)
	repo := &captureRepo{}
	arbiter := bidding.NewArbiter(repo, bidding.WithTimeout(5*time.Second))
	manager := trip.NewManager(nil, nil)
	validator := auth.NewValidator("top-secret")
	server := NewServer(arbiter, manager, validator)
	router := gin.New()
	server.RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodPost, "/bids/evaluate", bytes.NewReader([]byte(`{"request":{},"bids":[]}`)))
	res := httptest.NewRecorder()

	// 2.- Exercise the handler and expect an unauthorized status without persisting bids.
	router.ServeHTTP(res, req)

	if res.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d", res.Code)
	}
	if len(repo.saved) != 0 {
		t.Fatalf("arbiter should not persist bids on unauthorized request")
	}
}
