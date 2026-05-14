#!/usr/bin/env bash
set -euo pipefail

SINCE="${1:-15m}"
LOG_FILE="$(mktemp)"
trap 'rm -f "$LOG_FILE"' EXIT

docker compose logs --since "$SINCE" --no-color >"$LOG_FILE"

db_pattern='MongoDB:|server selection error|no such host|dial tcp|connection refused|wrong-mongo-host'
restart_pattern='(Restarting|restarting|restart loop|back-off|Back-off|unhealthy|Exited \([0-9]+\))'

db_hits="$(grep -Eic "$db_pattern" "$LOG_FILE" || true)"
restart_hits="$(grep -Eic "$restart_pattern" "$LOG_FILE" || true)"

echo "Log scan window: last $SINCE"
echo "Database failure pattern matches: $db_hits"
echo "Restart-loop pattern matches: $restart_hits"

if (( db_hits > 0 )); then
  echo
  echo "Critical log matches:"
  grep -Ein "$db_pattern" "$LOG_FILE" | head -n 20
  exit 2
fi

if (( restart_hits > 0 )); then
  echo
  echo "Restart-related log matches:"
  grep -Ein "$restart_pattern" "$LOG_FILE" | head -n 20
  exit 1
fi

echo "No critical patterns detected."
