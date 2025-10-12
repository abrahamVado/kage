#!/bin/sh
# //1.- Confirm docker compose is available before attempting any status checks.
if ! command -v docker >/dev/null 2>&1; then
  echo "Docker CLI is required for the sanity check." >&2
  exit 1
fi

# //2.- Print a concise table of the stack so operators can verify each service is healthy.
docker compose -f "$(dirname "$0")/../docker-compose.yml" ps
