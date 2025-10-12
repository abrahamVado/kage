# Deployment & Operations Overview

## Responsibility
- Orchestrate continuous delivery for backend, frontend, and mobile artifacts with environment parity.
- Maintain observability, incident response procedures, and operational automation.

## High-Level Architecture
- GitHub Actions triggers build pipelines that publish containers to an internal registry and create release candidates.
- ArgoCD promotes workloads across staging and production clusters with progressive delivery gates.
- Terraform manages cloud infrastructure, while Crossplane provisions managed services declaratively.

## Conventions
- State management for configuration uses GitOps repositories with locked main branches and promotion PRs.
- Super comments using `//1.-`, `//2.-` belong in Terraform modules or pipeline scripts explaining sequential automation steps.
- Infrastructure or pipeline files should not exceed 350 lines; break apart modules and reusable actions proactively.

## Local Development Stack
- `deploy/docker-compose.yml` builds dedicated images for backend, frontend, mobile, and database services with live-reload tooling baked in.
- Each service runs on the shared `kage-internal` bridge network so components can communicate over container DNS.
- Bind mounts map the repository into the containers, while named volumes retain dependency caches between restarts.

## Startup Workflow
1. `cd deploy` to keep subsequent commands scoped to the deployment assets.
2. `docker compose up --build` to create the images defined in the service-specific Dockerfiles and start all containers.
3. Watch the compose logs for `database-1  ...  healthy` before expecting the backend to finish bootstrapping against MariaDB.
4. Access the stack on the default ports:
   - Backend API: http://localhost:8080
   - Frontend Next.js dev server: http://localhost:3000
   - Flutter web runner: http://localhost:8082
   - MariaDB: localhost:3306 (credentials `kage` / `kagepass`)

## Teardown Workflow
- `docker compose down` stops the stack without deleting developer data in `database_data`, `frontend_node_modules`, `flutter_pub_cache`, and `go_pkg` volumes.
- `docker compose down --volumes` performs a clean slate teardown by removing cached dependencies and the database state.

## Sanity Verification
- Run `../deploy/scripts/stack-sanity.sh` after `docker compose up` completes to print the container table and confirm that every service reports the `running` state.
- Alternatively, `docker compose ps --format json | jq '.[].State'` highlights non-running containers quickly when troubleshooting.

## Troubleshooting
- If a container exits immediately, inspect logs via `docker compose logs <service>` to surface build or runtime errors.
- Delete named volumes with `docker volume rm` when dependency caches become corrupted or incompatible.
- MariaDB initialization failures are often resolved by ensuring ports `3306` and data directories are not in use by a local installation.
- Flutter hot reload issues typically disappear after running `flutter clean` inside the container shell followed by a restart.
