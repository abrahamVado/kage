# Data Platform Overview

## Responsibility
- Maintain authoritative storage for user profiles, trips, bids, and settlement records.
- Provide analytical exports and replication feeds for downstream ML and reporting systems.

## High-Level Architecture
- Hybrid storage using PostgreSQL for OLTP workloads and BigQuery for analytical aggregation.
- Change Data Capture (CDC) streams out of PostgreSQL into Kafka topics, feeding read models and warehouses.
- Schema migrations are orchestrated via Atlas, with automated drift detection in CI/CD.

## Conventions
- State management follows event sourcing snapshots plus relational projections for high-volume reads.
- Super comments employing `//1.-`, `//2.-` annotate complex migration scripts or data repair jobs.
- SQL or migration files must stay under 250 lines; partition large operations into incremental scripts.
