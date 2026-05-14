#!/usr/bin/env bash
set -euo pipefail

PROM_URL="${PROM_URL:-http://localhost:9090}"

command -v curl >/dev/null 2>&1 || { echo "curl is required" >&2; exit 1; }

if command -v python3 >/dev/null 2>&1; then
  PYTHON_BIN="python3"
elif command -v python >/dev/null 2>&1; then
  PYTHON_BIN="python"
else
  echo "python3 or python is required" >&2
  exit 1
fi

query_value() {
  local prom_query="$1"
  curl -fsS --get "${PROM_URL}/api/v1/query" --data-urlencode "query=${prom_query}" \
    | "${PYTHON_BIN}" -c 'import json,sys; d=json.load(sys.stdin); r=d.get("data",{}).get("result",[]); print(r[0]["value"][1] if r else "NaN")'
}

cpu_usage="$(query_value '(1 - avg(rate(node_cpu_seconds_total{mode="idle"}[5m]))) * 100')"
memory_usage="$(query_value 'avg((1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)) * 100)')"
request_rate="$(query_value 'sum(rate(http_requests_total[1m]))')"
error_rate="$(query_value '(sum(rate(http_requests_total{status=~"5.."}[1m])) / sum(rate(http_requests_total[1m]))) * 100')"
restart_frequency="$(query_value 'sum(changes(process_start_time_seconds{job=~"frontend|auth-service|product-service|order-service|user-service|chat-service"}[15m]))')"
order_cpu="$(query_value 'sum(rate(process_cpu_seconds_total{job="order-service"}[1m])) * 100')"

cat <<EOF
Capacity metrics snapshot
-------------------------
Prometheus URL:                  ${PROM_URL}
CPU usage (%):                   ${cpu_usage}
Memory utilization (%):          ${memory_usage}
Request rate (req/sec):          ${request_rate}
Error rate (%):                  ${error_rate}
Restart frequency (15m):         ${restart_frequency}
Order-service CPU (% single-vCPU): ${order_cpu}
EOF
