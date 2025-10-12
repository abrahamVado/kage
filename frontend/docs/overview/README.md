# Frontend Experience Overview

## Responsibility
- Deliver the web application for rider, driver, and support personas with responsive layouts.
- Surface real-time trip updates, bidding feedback, and secure authentication flows.

## High-Level Architecture
- React-based micro-frontends composed through a module federation shell with shared design system tokens.
- Data fetching occurs through GraphQL hooks backed by the backend gateway with optimistic UI updates.
- WebSockets subscribe to trip and bidding channels, reconciling local caches upon event delivery.

## Conventions
- State management relies on Redux Toolkit query slices and persisted caches scoped per persona.
- Components must include `//1.-`, `//2.-` style super comments near complex effects or reducers.
- Component files should not exceed 300 lines; decompose into presentational and hook modules when approaching 250 lines.
