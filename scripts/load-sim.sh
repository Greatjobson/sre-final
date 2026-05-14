#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost}"
ENDPOINT="${1:-/api/product}"
REQUESTS="${REQUESTS:-500}"
CONCURRENCY="${CONCURRENCY:-25}"

RESULTS_FILE="$(mktemp)"
trap 'rm -f "$RESULTS_FILE"' EXIT

export BASE_URL ENDPOINT RESULTS_FILE

echo "Running load simulation"
echo "  Base URL:     $BASE_URL"
echo "  Endpoint:     $ENDPOINT"
echo "  Requests:     $REQUESTS"
echo "  Concurrency:  $CONCURRENCY"

start_ts="$(date +%s)"

seq "$REQUESTS" | xargs -I{} -P "$CONCURRENCY" sh -c '
  code=$(curl -sS -o /dev/null -w "%{http_code}" "${BASE_URL}${ENDPOINT}" || echo 000)
  echo "$code" >> "$RESULTS_FILE"
'

end_ts="$(date +%s)"
duration=$((end_ts - start_ts))
if (( duration <= 0 )); then
  duration=1
fi

total="$(wc -l < "$RESULTS_FILE" | tr -d ' ')"
ok="$(grep -Ec "^(2|3)" "$RESULTS_FILE" || true)"
err4="$(grep -Ec "^4" "$RESULTS_FILE" || true)"
err5="$(grep -Ec "^5" "$RESULTS_FILE" || true)"
err_net="$(grep -Ec "^000$" "$RESULTS_FILE" || true)"
failed=$((total - ok))
rps="$(awk -v n="$total" -v d="$duration" 'BEGIN { printf "%.2f", n/d }')"
error_pct="$(awk -v e="$failed" -v n="$total" 'BEGIN { if (n==0) { print "0.00" } else { printf "%.2f", (e/n)*100 } }')"

cat <<EOF

Load simulation summary
-----------------------
Total requests:         $total
Successful (2xx/3xx):  $ok
Client errors (4xx):   $err4
Server errors (5xx):   $err5
Network errors (000):  $err_net
Elapsed time (sec):    $duration
Throughput (req/sec):  $rps
Error rate (%):        $error_pct
EOF
