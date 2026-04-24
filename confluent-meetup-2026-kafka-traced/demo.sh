#!/usr/bin/env bash
# demo.sh — bring up the stack, hit /order, print the trace ID + all observability URLs.
#
# Usage:
#   ./demo.sh                       # one order with the slide-20 payload
#   ./demo.sh 10                    # fire 10 orders back-to-back
#   ORDER_ID=X ITEM=y QTY=1 ./demo.sh
#
# Requires: docker, curl.

set -euo pipefail

cd "$(dirname "$0")"

GATEWAY="${GATEWAY:-http://localhost:11000}"
GRAFANA="${GRAFANA:-http://localhost:11030}"
PROMETHEUS="${PROMETHEUS:-http://localhost:11090}"
METRICS_GATEWAY="${METRICS_GATEWAY:-http://localhost:11121}"
METRICS_ORDER="${METRICS_ORDER:-http://localhost:11122}"
METRICS_NOTIF="${METRICS_NOTIF:-http://localhost:11123}"
TRACER="${TRACER:-https://tracer.gofr.dev}"
DASHBOARD_UID="${DASHBOARD_UID:-gofr-pubsub}"

ORDER_ID="${ORDER_ID:-DEMO-$(date +%s)}"
ITEM="${ITEM:-sticker}"
QTY="${QTY:-3}"
COUNT="${1:-1}"

echo "▸ Ensuring stack is up…"
docker compose up -d --quiet-pull >/dev/null

printf '▸ Waiting for api-gateway at %s ' "$GATEWAY"
for i in $(seq 1 60); do
  # -f makes curl exit non-zero on >=400, so this also retries through the
  # short window where Kafka topics aren't yet created (publish returns 500).
  if curl -fsS -o /dev/null -X POST "$GATEWAY/order" \
        -H "Content-Type: application/json" \
        -d '{"orderId":"probe","item":"x","qty":1}' 2>/dev/null; then
    printf ' ✓\n'
    break
  fi
  printf '.'
  sleep 1
  if [ "$i" = "60" ]; then
    printf '\n✗ gateway not responding cleanly after 60s — check: docker compose logs api-gateway\n'
    exit 1
  fi
done

echo
for i in $(seq 1 "$COUNT"); do
  id="$ORDER_ID"
  [ "$COUNT" -gt 1 ] && id="${ORDER_ID}-${i}"
  payload=$(printf '{"orderId":"%s","item":"%s","qty":%s}' "$id" "$ITEM" "$QTY")
  echo "▸ POST /order  →  $payload"
  resp=$(curl -fsS -X POST "$GATEWAY/order" \
    -H "Content-Type: application/json" \
    -d "$payload")
  echo "  ← $resp"
done

echo
echo "▸ Collecting trace IDs from api-gateway logs…"
sleep 2
trace_ids=$(
  docker compose logs api-gateway --tail=80 2>&1 \
    | grep -oE '"trace_id":"[a-f0-9]{32}"' \
    | tail -n "$COUNT" \
    | sed -E 's/.*"trace_id":"([a-f0-9]+)".*/\1/'
)

if [ -z "$trace_ids" ]; then
  echo "✗ no trace_id found — tracer may be disabled or logs rolled"
  exit 2
fi

echo
echo "═══════════════════════════════════════════════════════════════════════"
echo "  🔗 TRACES  (hosted tracer — allow ~10s for batch flush)"
echo "═══════════════════════════════════════════════════════════════════════"
echo "$trace_ids" | while IFS= read -r tid; do
  [ -z "$tid" ] && continue
  echo "    $TRACER/?traceID=$tid"
done

echo
echo "═══════════════════════════════════════════════════════════════════════"
echo "  📊 DASHBOARDS"
echo "═══════════════════════════════════════════════════════════════════════"
echo "    Grafana     ▸  $GRAFANA/d/$DASHBOARD_UID   (GoFr Pub/Sub dashboard)"
echo "    Prometheus  ▸  $PROMETHEUS"

echo
echo "═══════════════════════════════════════════════════════════════════════"
echo "  📈 RAW METRICS"
echo "═══════════════════════════════════════════════════════════════════════"
echo "    api-gateway           ▸  $METRICS_GATEWAY/metrics"
echo "    order-service         ▸  $METRICS_ORDER/metrics"
echo "    notification-service  ▸  $METRICS_NOTIF/metrics"

echo
echo "═══════════════════════════════════════════════════════════════════════"
echo "  🛠  DEMO TIPS"
echo "═══════════════════════════════════════════════════════════════════════"
echo "    Burst:   ./demo.sh 100"
echo "    Kill:    docker compose stop notification-service   # alerts pile up"
echo "    Resume:  docker compose start notification-service  # drains in real time"
echo "    Logs:    docker compose logs -f order-service notification-service"
echo
