# Kage Project Architecture Overview

## High-Level Composition
- **Flutter mobile client** located in `mobile/flutter_app` builds a driver-focused UI that wires authentication, bidding, tracking, and trip flows through Provider view models and shared services.【F:mobile/flutter_app/lib/main.dart†L1-L71】
- **Go backend** in `backend` exposes REST endpoints for bid evaluation and trip lifecycle management, plus a WebSocket hub for realtime rider/driver rooms.【F:backend/internal/api/server.go†L1-L105】【F:backend/internal/ws/hub.go†L1-L131】

## Go Backend Responsibilities
### HTTP API Layer
- The Gin-based server registers `/health`, `/bids/evaluate`, `/trips/:id/state`, and `/trips/:id/metrics` endpoints, enforcing bearer-token checks when a validator is configured.【F:backend/internal/api/server.go†L21-L105】
- Trip state transitions map JSON `action` payloads to manager methods that guard invalid state changes before responding.【F:backend/internal/api/server.go†L61-L100】

### Application Wiring
- `app.Build` optionally opens a MariaDB connection (via the MySQL driver) when a DSN is provided, instantiates the bidding arbiter and trip manager, registers REST and WebSocket routes, and returns a cleanup hook that shuts down background workers and closes shared resources.【F:backend/internal/app/app.go†L1-L88】

### Bidding Engine
- The arbiter filters bids by expiry, rider budgets, ETA, and proximity, ranks remaining candidates with a weighted score, persists the winner, and enforces an evaluation timeout to avoid blocking requests.【F:backend/internal/bidding/service.go†L1-L129】【F:backend/internal/bidding/service.go†L131-L204】

### Trip Lifecycle
- The trip manager tracks per-trip state, supports `start`, `pause`, `resume`, `cancel`, and `complete` actions with validation, and accumulates active/paused durations for auditing metrics.【F:backend/internal/trip/manager.go†L1-L117】【F:backend/internal/trip/manager.go†L119-L164】

### Realtime Hub
- The WebSocket hub upgrades `/ws/:role/:room` connections, maintains room membership, broadcasts JSON payloads to participants, and exposes occupancy stats. Graceful shutdown drains registration/unregistration queues before closing.【F:backend/internal/ws/hub.go†L1-L141】【F:backend/internal/ws/hub.go†L143-L208】

## Flutter Client Responsibilities
### Dependency Graph
- `DriverApp` seeds Provider instances for authentication, secure storage, WebSocket connectivity, geolocation publishing, and feature-specific view models that react to realtime updates.【F:mobile/flutter_app/lib/main.dart†L15-L71】

### Authentication Flow
- `AuthService.login` posts email/password credentials to `$BASE/api/v1/auth/login`, expecting a JSON body with `token` and `user_id`, and raises an exception when the Go backend rejects the request.【F:mobile/flutter_app/lib/core/services/auth_service.dart†L1-L41】
- `AuthViewModel` orchestrates login, saves tokens in secure storage, and starts downstream realtime/location flows upon success (see `features/auth` files).【F:mobile/flutter_app/lib/features/auth/auth_view_model.dart†L1-L81】

### Realtime Bidding and Tracking
- `WebSocketService` manages a single connection to `$BASE/ws?token=...`, buffers outbound JSON payloads, and streams inbound events to multiple listeners.【F:mobile/flutter_app/lib/core/services/websocket_service.dart†L1-L45】
- `BiddingViewModel` listens for `bid_offer` and `bid_selected` messages, keeps driver-facing bid lists sorted by proximity, and emits bid submissions/selections back to the socket.【F:mobile/flutter_app/lib/features/bidding/bidding_view_model.dart†L1-L76】
- `TrackingViewModel` watches for `rider_update` payloads, transforms them into ordered nearby-rider lists, and updates UI state as events arrive.【F:mobile/flutter_app/lib/features/tracking/tracking_view_model.dart†L1-L62】

### Location Streaming and Trip Controls
- `LocationService` requests geolocation permissions, subscribes to `Geolocator` position streams (or test overrides), and forwards periodic `location_update` payloads over the WebSocket connection.【F:mobile/flutter_app/lib/core/services/location_service.dart†L1-L56】
- Trip view models/screens toggle WebSocket commands for lifecycle events (e.g., start/pause/complete) and reflect backend metrics (see `features/trip`).【F:mobile/flutter_app/lib/features/trip/trip_view_model.dart†L1-L81】【F:mobile/flutter_app/lib/features/trip/trip_view.dart†L1-L67】

## Integration Gaps and Considerations
- **Authentication mismatch:** The Flutter client targets `/api/v1/auth/login`, but the Go API only exposes `/bids` and `/trips` routes; implementing a compatible login endpoint (or adjusting the client to match token expectations) is necessary for end-to-end auth.【F:mobile/flutter_app/lib/core/services/auth_service.dart†L16-L37】【F:backend/internal/api/server.go†L23-L105】
- **WebSocket contract drift:** The client connects to `/ws?token=...` and emits generic event types, whereas the Go hub expects `/ws/:role/:room` path parameters; harmonizing connection URLs and payload schemas is required for realtime flows to function.【F:mobile/flutter_app/lib/core/services/websocket_service.dart†L17-L45】【F:backend/internal/ws/hub.go†L59-L100】
- **Trip control parity:** Mobile view models assume the backend broadcasts selection and tracking events; ensure hub broadcasts generated by trip manager or bidding flows align with Flutter message handlers (`bid_offer`, `rider_update`, etc.).【F:mobile/flutter_app/lib/features/bidding/bidding_view_model.dart†L17-L61】【F:mobile/flutter_app/lib/features/tracking/tracking_view_model.dart†L23-L52】【F:backend/internal/ws/hub.go†L100-L140】

## Testing Baseline
- Go unit tests cover bidding heuristics, trip-state transitions, and API wiring. Running `go test ./...` from the `backend` directory currently succeeds, confirming core backend logic is deterministic.【e05581†L1-L9】

## Key Takeaways
1. The repository packages a Go mobility backend with a Flutter driver prototype; aligning HTTP/WebSocket contracts is the main prerequisite for an integrated experience.
2. Backend modules already encapsulate bidding arbitration, trip orchestration, and realtime messaging; clients can focus on UX once transport schemas stabilize.
3. Mobile services/view models are structured for Provider-driven state management, making it straightforward to adjust endpoints when the backend contract finalizes.
