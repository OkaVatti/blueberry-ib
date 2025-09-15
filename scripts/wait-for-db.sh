#!/usr/bin/env bash
set -euo pipefail

# Path to the docker-compose file used for DB
COMPOSE_FILE="apps/backend/docker-compose.yml"
SERVICE_NAME="db"      # the service name inside that compose file
MAX_SECS=60
SLEEP_INTERVAL=1

echo "[wait-for-db] ensuring docker compose service '$SERVICE_NAME' is up using $COMPOSE_FILE"

# get container id for service
container_id=$(docker compose -f "$COMPOSE_FILE" ps -q "$SERVICE_NAME")
if [ -z "$container_id" ]; then
  echo "[wait-for-db] container for service '$SERVICE_NAME' not found. Did you run docker compose up -d?"
  exit 1
fi

echo "[wait-for-db] watching container $container_id for health..."

elapsed=0
while [ $elapsed -lt $MAX_SECS ]; do
  # If health is configured, use it; otherwise fall back to .State.Status
  status=$(docker inspect --format='{{if .State.Health}}{{.State.Health.Status}}{{else}}{{.State.Status}}{{end}}' "$container_id" 2>/dev/null || true)

  if [ "$status" = "healthy" ] || [ "$status" = "running" ]; then
    echo "[wait-for-db] container status: $status — ready."
    exit 0
  fi

  # If status is empty, container might not be ready yet; print and try again
  echo "[wait-for-db] waiting... status='$status' ($elapsed/$MAX_SECS)"
  sleep $SLEEP_INTERVAL
  elapsed=$((elapsed + SLEEP_INTERVAL))
done

echo "[wait-for-db] timeout waiting for DB to become healthy (waited $MAX_SECS seconds)."
echo "[wait-for-db] last known status: $status"
exit 2

