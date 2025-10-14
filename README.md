# Kage

## Description
Kage is a mobility marketplace backend that evaluates driver bids, manages trip lifecycles, and streams rider-driver updates through a unified HTTP and WebSocket surface. The backend wires together a Gin API layer, bidding arbiter, trip manager, authentication guard, and realtime hub so that riders receive the best available offer while trip state changes remain auditable.

## Architecture
- **HTTP API** — Exposes health, bid evaluation, trip state transitions, and trip metrics endpoints while delegating token checks to the validator.
- **Bidding Arbiter** — Filters bids by freshness, budget, ETA, and proximity; scores candidates; and persists the accepted offer with a repository abstraction.
- **Trip Manager** — Coordinates start, pause, resume, cancel, and complete actions, aggregates timing metrics, and optionally records events to storage.
- **Realtime Hub** — Upgrades WebSocket connections for riders and drivers, tracks room membership, and broadcasts payloads while supporting graceful shutdown.
- **Application Composition** — Bootstraps the database (when configured), instantiates the arbiter, trip manager, validator, and hub, registers routes, and provides a cleanup hook.

## Component Diagram
```mermaid
flowchart LR
    subgraph Clients
        RiderApp
        DriverApp
        OpsTools
    end

    RiderApp -- REST --> APIGateway
    DriverApp -- REST --> APIGateway
    OpsTools -- WebSocket --> RealtimeHub

    subgraph Backend
        APIGateway[HTTP API (Gin)]
        Arbiter[Bidding Arbiter]
        TripMgr[Trip Manager]
        AuthSvc[Token Validator]
        RealtimeHub[WebSocket Hub]
    end

    APIGateway -->|Evaluate Bids| Arbiter
    APIGateway -->|Trip Actions & Metrics| TripMgr
    APIGateway -->|Authorization| AuthSvc
    RealtimeHub -->|Broadcasts| RiderApp
    RealtimeHub -->|Broadcasts| DriverApp
```

## Repository Layout
```text
kage/
├── backend/
│   ├── cmd/server/           # Go entrypoint for the HTTP service
│   ├── internal/api/         # REST handlers and route registration
│   ├── internal/app/         # Dependency wiring and lifecycle management
│   ├── internal/auth/        # Shared-secret token validation
│   ├── internal/bidding/     # Bid scoring logic and repositories
│   ├── internal/contracts/   # Domain DTOs shared across services
│   ├── internal/geo/         # Geographic helpers for scoring
│   ├── internal/trip/        # Trip lifecycle orchestration
│   └── internal/ws/          # WebSocket hub for realtime updates
├── docs/                     # Narrative documentation sets
│   ├── handbook/
│   └── system-overview/
├── frontend/                 # Web experience documentation and assets
├── mobile/                   # Mobile docs and Flutter prototype
├── tests/                    # Structural integrity checks and smoke tests
└── package.json              # Workspace tooling metadata
```
