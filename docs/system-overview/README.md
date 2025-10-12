# System Overview

## Authentication Workflow
```text
[Client App] -> [Edge Auth API] -> [Identity Provider]
                          \-> [Token Introspection]
```
- Clients request OAuth 2.1 PKCE authorization, exchanging codes for JWT access tokens.
- Edge gateway validates device fingerprinting and rate limits credential attempts.
- Refresh tokens rotate on every use; revoked tokens propagate via message bus to all services within 500 ms.

## Bidding Lifecycle
1. Rider publishes a trip request with constraints (time, price ceiling, pickup window).
2. Matching service broadcasts bid invitations to eligible drivers through the messaging mesh.
3. Drivers submit sealed bids; scoring service evaluates latency, rating, and surge multipliers.
4. Auction service selects the optimal bid and emits a `trip.bid.awarded` event for downstream consumers.

## Trip Lifecycle Stages
- **Initiation:** Trip service creates trip aggregate, reserving inventory and locking payment method.
- **En Route:** Location service streams telemetry; ETA predictions update the rider app.
- **Pickup Confirmation:** Driver verifies rider, triggering authentication handshake and escrow hold.
- **In Transit:** Fare adjustments apply automatically based on telemetry deltas.
- **Completion:** Payment service captures funds, generates receipts, and updates loyalty balances.
- **Post-Trip:** Feedback service ingests ratings, while analytics processes feed the recommender models.

## Inter-Service Messaging Contracts
- Events follow CloudEvents 1.0 with JSON payloads and schema registry validation.
- Topic taxonomy: `trip.*` for lifecycle, `bid.*` for auctions, `auth.*` for credential updates, `ops.*` for monitoring.
- Command messages require idempotency keys and include `traceparent` headers for distributed tracing.
- Dead-letter queues retain failed messages for 48 hours, triaged via automated replay workflows.
