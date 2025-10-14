package trip

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"kage/backend/internal/contracts"
)

func TestSQLEventRepositoryRecordEvent(t *testing.T) {
	//1.- Prepare the sqlmock database connection to capture executed queries.
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewSQLEventRepository(db)
	event := contracts.TripEvent{
		TripID:     "trip-456",
		State:      contracts.TripStateComplete,
		OccurredAt: time.Unix(1735689600, 0).UTC(),
		Notes:      "fare settled",
	}

	//2.- Expect the INSERT into trip_events with the event payload.
	mock.ExpectExec(`INSERT INTO trip_events \(trip_id, state, occurred_at, notes\) VALUES \(\?, \?, \?, \?\)`).
		WithArgs(
			event.TripID,
			event.State,
			event.OccurredAt,
			event.Notes,
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := repo.RecordEvent(context.Background(), event); err != nil {
		t.Fatalf("RecordEvent: %v", err)
	}

	//3.- Ensure the mocked expectations were satisfied once the call completes.
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
