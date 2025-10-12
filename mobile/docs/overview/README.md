# Mobile Platform Overview

## Responsibility
- Provide native-quality rider and driver apps with offline capabilities and push notifications.
- Manage trip lifecycle operations, bidding participation, and secure credential storage on devices.

## High-Level Architecture
- React Native monorepo structured via feature packages sharing a common UI kit and navigation library.
- Native modules wrap location, camera, and secure enclave APIs, exposed through typed bridges.
- Background sync workers coordinate with the backend through GraphQL mutations and websocket reconnections.

## Conventions
- State management uses Zustand stores synchronized with Redux slices to keep UI responsive while persisting server truth.
- Super comments (`//1.-`, `//2.-`) must annotate asynchronous side-effects and native bridge boundaries.
- Feature modules stay under 400 lines; split navigation stacks and service hooks once files exceed 275 lines.
