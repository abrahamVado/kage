# Database Structure

## Connection Lifecycle
- The backend optionally opens a MariaDB-compatible DSN during application wiring. When `DBDSN` is populated, `sql.Open("mysql", cfg.DBDSN)` constructs the shared handle that repositories reuse for persistence, and it is closed during application cleanup. 【F:backend/internal/app/app.go†L29-L55】

## Tables

### accepted_bids
- **Purpose:** Stores the winning bid for a trip after arbitration.
- **Write Path:** `SQLRepository.SaveAcceptedBid` issues `INSERT INTO accepted_bids (bid_id, trip_id, driver_id, price, accepted_at)` with the accepted bid payload. 【F:backend/internal/bidding/repository.go†L24-L31】
- **Columns:**
  - `bid_id` (`VARCHAR`): identifier of the accepted bid, sourced from `contracts.AcceptedBid.BidID`. 【F:backend/internal/contracts/contracts.go†L25-L30】
  - `trip_id` (`VARCHAR`): trip identifier linked to the rider request. 【F:backend/internal/contracts/contracts.go†L25-L30】
  - `driver_id` (`VARCHAR`): driver who won the auction. 【F:backend/internal/contracts/contracts.go†L25-L30】
  - `price` (`DECIMAL`): fare amount of the accepted bid. 【F:backend/internal/contracts/contracts.go†L25-L30】
  - `accepted_at` (`TIMESTAMP`): UTC timestamp when the bid was finalized. 【F:backend/internal/contracts/contracts.go†L25-L30】

### trip_events
- **Purpose:** Captures lifecycle transitions for auditing rider trips.
- **Write Path:** `SQLEventRepository.RecordEvent` executes `INSERT INTO trip_events (trip_id, state, occurred_at, notes)` to persist each event. 【F:backend/internal/trip/repository.go†L24-L31】
- **Columns:**
  - `trip_id` (`VARCHAR`): identifier of the trip undergoing a state change. 【F:backend/internal/contracts/contracts.go†L37-L41】
  - `state` (`ENUM`/`VARCHAR`): new `TripState` value such as `pending`, `active`, or `complete`. 【F:backend/internal/contracts/contracts.go†L32-L36】【F:backend/internal/contracts/contracts.go†L37-L41】
  - `occurred_at` (`TIMESTAMP`): UTC timestamp describing when the event happened. 【F:backend/internal/contracts/contracts.go†L37-L41】
  - `notes` (`TEXT`): optional additional details about the transition. 【F:backend/internal/contracts/contracts.go†L37-L41】

## Testing Hooks
- Repository tests use `sqlmock` to assert the exact SQL shape for both tables, providing living documentation for expected inserts. 【F:backend/internal/bidding/repository_test.go†L9-L46】【F:backend/internal/trip/repository_test.go†L9-L46】
