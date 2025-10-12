# Backend Services Overview

## Responsibility
- Provide core APIs for authentication, bidding, trip orchestration, and messaging fan-out.
- Enforce domain invariants, auditing, and threat detection before state mutations reach downstream systems.

## High-Level Architecture
- Services follow a hexagonal boundary: transport adapters (REST/gRPC) delegate to application services, which orchestrate domain aggregates.
- Shared libraries expose serialization, cryptography helpers, and message-contract validators consumed by all services.
- Event-driven workflows emit canonical events onto the message bus for other domains to react.

## Conventions
- State management uses immutable domain events and projections; mutable caches must reconcile via event sourcing replays nightly.
- Super comments documenting complex flows must follow `//1.-`, `//2.-` step annotations aligned with service handlers.
- Code files should remain below 500 lines; split bounded contexts into dedicated modules once they exceed 350 lines to avoid monolith creep.
