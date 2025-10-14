package bidding

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"kage/backend/internal/contracts"
)

func TestSQLRepositorySaveAcceptedBid(t *testing.T) {
	//1.- Create a sqlmock database handle to observe the executed statement.
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewSQLRepository(db)
	accepted := contracts.AcceptedBid{
		BidID:      "bid-123",
		TripID:     "trip-456",
		DriverID:   "driver-789",
		Price:      4200,
		AcceptedAt: time.Unix(1735689600, 0).UTC(),
	}

	//2.- Expect an INSERT with the accepted bid payload and ensure it succeeds.
	mock.ExpectExec(`INSERT INTO accepted_bids \(bid_id, trip_id, driver_id, price, accepted_at\) VALUES \(\?, \?, \?, \?, \?\)`).
		WithArgs(
			accepted.BidID,
			accepted.TripID,
			accepted.DriverID,
			accepted.Price,
			accepted.AcceptedAt,
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := repo.SaveAcceptedBid(context.Background(), accepted); err != nil {
		t.Fatalf("SaveAcceptedBid: %v", err)
	}

	//3.- Verify the repository issued the expected SQL statement exactly once.
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
