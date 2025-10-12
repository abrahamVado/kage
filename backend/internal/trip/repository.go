package trip

import (
	"context"
	"database/sql"

	"kage/backend/internal/contracts"
)

// EventRepository persists trip lifecycle events.
type EventRepository interface {
	RecordEvent(ctx context.Context, event contracts.TripEvent) error
}

// SQLEventRepository writes trip events using database/sql.
type SQLEventRepository struct {
	db *sql.DB
}

// NewSQLEventRepository builds a repository backed by MariaDB.
func NewSQLEventRepository(db *sql.DB) *SQLEventRepository {
	return &SQLEventRepository{db: db}
}

// RecordEvent inserts a trip event row.
func (r *SQLEventRepository) RecordEvent(ctx context.Context, event contracts.TripEvent) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO trip_events (trip_id, state, occurred_at, notes) VALUES (?, ?, ?, ?)`,
		event.TripID, event.State, event.OccurredAt, event.Notes,
	)
	return err
}
