#!/usr/bin/env bash
set -euo pipefail

COMPOSE_FILE="${1:-docker-compose.yml}"
ENV_FILE="${2:-.env}"

required_services=(
  gateway
  frontend
  auth-service
  product-service
  order-service
  user-service
  chat-service
  mongo
  postgres
  prometheus
  grafana
  node-exporter
)

required_env_vars=(
  JWT_SECRET
  ADMIN_EMAIL
  POSTGRES_DB
  POSTGRES_USER
  POSTGRES_PASSWORD
  GRAFANA_ADMIN_PASSWORD
)

fail() {
  echo "ERROR: $*" >&2
  exit 1
}

pass() {
  echo "OK: $*"
}

command -v docker >/dev/null 2>&1 || fail "docker is not installed"
docker compose version >/dev/null 2>&1 || fail "docker compose plugin is not available"

if [[ -f "$ENV_FILE" ]]; then
  for var_name in "${required_env_vars[@]}"; do
    grep -Eq "^${var_name}=.+" "$ENV_FILE" || fail "${var_name} missing in $ENV_FILE"
  done
  pass "env file variables are present"
else
  fail "env file '$ENV_FILE' not found"
fi

[[ -f "$COMPOSE_FILE" ]] || fail "compose file '$COMPOSE_FILE' not found"
docker compose -f "$COMPOSE_FILE" config -q || fail "docker compose config validation failed"
pass "compose config is valid"

services_output="$(docker compose -f "$COMPOSE_FILE" config --services)"
for svc in "${required_services[@]}"; do
  echo "$services_output" | grep -qx "$svc" || fail "required service '$svc' is missing"
done
pass "required services are present"

if grep -Eq "wrong-mongo-host" "$COMPOSE_FILE"; then
  fail "invalid Mongo hostname (wrong-mongo-host) found in $COMPOSE_FILE"
fi
pass "MongoDB host sanity check passed"

for svc in gateway auth-service product-service order-service user-service chat-service frontend; do
  if ! awk "/^  ${svc}:/{in_svc=1} in_svc && /^  [a-zA-Z0-9_-]+:/{if(!match(\$0, /^  ${svc}:/)){in_svc=0}} in_svc{print}" "$COMPOSE_FILE" | grep -q "/health"; then
    fail "health check for '$svc' does not use /health endpoint"
  fi
done
pass "health check endpoints are configured"

echo "Pre-deployment validation completed successfully."
