package bidding

import (
	"context"
	"database/sql"

	"kage/backend/internal/contracts"
)

// Repository persists bidding outcomes.
type Repository interface {
	SaveAcceptedBid(ctx context.Context, bid contracts.AcceptedBid) error
}

// SQLRepository stores accepted bids inside MariaDB using database/sql.
type SQLRepository struct {
	db *sql.DB
}

// NewSQLRepository wires the repository with the given database handle.
func NewSQLRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

// SaveAcceptedBid inserts the accepted bid row.
func (r *SQLRepository) SaveAcceptedBid(ctx context.Context, bid contracts.AcceptedBid) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO accepted_bids (bid_id, trip_id, driver_id, price, accepted_at) VALUES (?, ?, ?, ?, ?)`,
		bid.BidID, bid.TripID, bid.DriverID, bid.Price, bid.AcceptedAt,
	)
	return err
}
