#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost}"
REQUESTS="${REQUESTS:-300}"
CONCURRENCY="${CONCURRENCY:-20}"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "=== Pre-load metrics ==="
"${SCRIPT_DIR}/capacity-metrics.sh"
echo

echo "=== Load scenario: homepage ==="
REQUESTS="$REQUESTS" CONCURRENCY="$CONCURRENCY" BASE_URL="$BASE_URL" "${SCRIPT_DIR}/load-sim.sh" "/"
echo

echo "=== Load scenario: product API ==="
REQUESTS="$REQUESTS" CONCURRENCY="$CONCURRENCY" BASE_URL="$BASE_URL" "${SCRIPT_DIR}/load-sim.sh" "/api/product"
echo

echo "=== Load scenario: auth refresh API ==="
REQUESTS="$REQUESTS" CONCURRENCY="$CONCURRENCY" BASE_URL="$BASE_URL" "${SCRIPT_DIR}/load-sim.sh" "/auth/refresh"
echo

echo "=== Post-load metrics ==="
"${SCRIPT_DIR}/capacity-metrics.sh"
echo

echo "Capacity run completed."
