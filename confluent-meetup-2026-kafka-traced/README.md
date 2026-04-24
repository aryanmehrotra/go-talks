# Kafka, Traced. (Confluent Meetup 2026)

> One trace ID. Three services. Zero plumbing — with **GoFr**.

A self-contained, runnable demo for the Confluent Meetup talk **"Kafka, Traced."**
Three GoFr microservices fanning out over Kafka, with a single OpenTelemetry trace
following every message from HTTP entry → publish → broker → subscribe → re-publish
→ subscribe.

```
   curl POST /order
         │
         ▼
 ┌──────────────────┐  publish      ┌─────────┐  subscribe   ┌──────────────────┐  publish    ┌─────────┐  subscribe   ┌──────────────────────┐
 │  api-gateway     │ ─── orders ──▶│ Kafka   │ ──────────▶ │  order-service   │ ── alerts ─▶│ Kafka   │ ──────────▶  │ notification-service │
 │  (HTTP :8000)    │               │  topic  │              │  (validate +     │             │  topic  │              │ (deliver - log)      │
 └──────────────────┘               └─────────┘              │   re-publish)    │             └─────────┘              └──────────────────────┘
                                                              └──────────────────┘
```

All three services emit OTel spans; GoFr injects/extracts the W3C `traceparent`
through Kafka headers automatically. The hosted tracer at
[tracer.gofr.dev](https://tracer.gofr.dev) renders one connected trace.

---

## Quickstart

```bash
git clone https://github.com/aryanmehrotra/go-talks.git
cd go-talks/confluent-meetup-2026-kafka-traced

docker-compose up -d        # 7 containers: zk, kafka, prom, grafana, 3 services

# Trigger the flow
curl -X POST http://localhost:11000/order \
  -H "Content-Type: application/json" \
  -d '{"orderId":"ORD-42","item":"GoFr T-shirt","qty":1}'

# Inspect
open https://tracer.gofr.dev           # GoFr hosted tracer — paste the trace ID from response headers
open http://localhost:11030            # Grafana — "GoFr Pub/Sub" dashboard (anonymous Viewer)
open http://localhost:11090            # Prometheus
curl http://localhost:11121/metrics    # raw counters from api-gateway
```

> **Tracer.** The demo ships with `TRACE_EXPORTER=gofr` so traces land in
> GoFr's hosted tracer at https://tracer.gofr.dev — no local tracing backend
> needed, just internet. For an offline demo, add a `jaeger` service to the
> compose file and set `TRACE_EXPORTER=otlp` / `TRACER_URL=jaeger:4317` in
> each `configs/.env`.

You should see one trace with **5 spans** connected across all 3 services:

```
HTTP POST /order      (api-gateway,           kind=server)
└─ kafka-publish      (api-gateway,           kind=producer, topic=orders)
   └─ kafka-subscribe (order-service,         kind=consumer, topic=orders)   [+ link to producer]
      └─ kafka-publish (order-service,        kind=producer, topic=alerts)
         └─ kafka-subscribe (notification-svc, kind=consumer, topic=alerts)  [+ link to producer]
```

---

## Try the failure mode (the resilience punchline)

```bash
docker-compose stop notification-service
# Hit /order a few times — alerts pile up on the broker
curl -X POST http://localhost:11000/order -H "Content-Type: application/json" \
  -d '{"orderId":"X","item":"y","qty":1}'
# Grafana: subscribe rate for notification-service is 0; publish keeps climbing.
docker-compose start notification-service
# Drain in real time. Trace context survives the rebalance.
```

---

## Folder layout

| Path | What |
|---|---|
| `api-gateway/` | HTTP → Publish (GoFr, ~30 lines) |
| `order-service/` | Subscribe → Publish (GoFr, ~25 lines) |
| `notification-service/` | Subscribe → log (GoFr, ~20 lines) |
| `raw-kafka-baseline/` | **Same flow, written without GoFr.** ~270 lines per service — the "what you're avoiding" comparison shown on stage. Separate `go.mod`; run with `GOWORK=off`. |
| `docker-compose.yaml` | Kafka + Zookeeper + Prometheus + Grafana + 3 services |
| `prometheus/`, `grafana/` | Auto-provisioned dashboards |
| `slides/` | The reveal.js deck used for the talk (`open slides/index.html`). 22 slides. |

---

## How tracing actually works

The mechanism is ~50 lines of code living in
[`pkg/gofr/datasource/pubsub/kafka/tracing.go`](https://github.com/gofr-dev/gofr/blob/main/pkg/gofr/datasource/pubsub/kafka/tracing.go)
in the GoFr framework:

1. A `headerCarrier` type wraps `[]kafka.Header` to satisfy
   `propagation.TextMapCarrier` (the OTel interface).
2. On publish, GoFr calls `otel.GetTextMapPropagator().Inject(ctx, &carrier)` —
   the active span's `traceparent` becomes a Kafka message header.
3. On subscribe, GoFr calls `Extract` on the inbound headers, gets a
   `SpanContext`, and starts the consumer span **as a child of the producer**
   (so `HTTP → publish → subscribe → publish → subscribe` shows up as one
   connected waterfall with a single trace ID, and sampling is inherited via
   `ParentBased`). It **also** attaches a `trace.Link` to the producer, so
   OTel-aware UIs can still model fan-out (one message, many consumer groups).

Read `raw-kafka-baseline/` to see exactly how much code you'd write to do this
yourself. The slide deck does the side-by-side on stage.

---

## Configuration

Each service's `configs/.env` is the entire configuration surface. Notable
knobs:

| Env var | Effect |
|---|---|
| `PUBSUB_BACKEND=KAFKA` | Picks the Kafka driver. Swap to `NATS`, `GOOGLE`, etc. — tracing works the same way. |
| `PUBSUB_BROKER` | Comma-separated brokers |
| `CONSUMER_ID` | Kafka consumer-group ID |
| `TRACE_EXPORTER=gofr` | **Default** — GoFr-hosted tracer; view at tracer.gofr.dev |
| `TRACE_EXPORTER=otlp` | OTLP/gRPC (Jaeger, Tempo, Datadog, Honeycomb…) |
| `TRACER_URL` | OTLP/Zipkin collector endpoint (omit for `gofr` exporter) |
| `TRACER_RATIO` | Head-based sampling (1.0 = all spans) |
| `KAFKA_BATCH_SIZE`, `KAFKA_BATCH_BYTES`, `KAFKA_BATCH_TIMEOUT` | Producer batching |
| `KAFKA_SECURITY_PROTOCOL`, `KAFKA_SASL_*`, `KAFKA_TLS_*` | SASL/TLS for production brokers |

---

## The slides

The reveal.js deck used for the talk lives in `slides/`:

```bash
open slides/index.html
# Navigation: ← / → or Space. Press "f" for fullscreen, "s" for speaker view.
```

22 slides covering: Kafka 101 (roles → topics → partitions → consumer groups →
fan-out → message anatomy), why HTTP tracing doesn't apply to async fan-out,
OTel propagators & carriers, hand-rolling span-link propagation, then the
GoFr version as a direct `diff` — plus the production knobs and live-demo
script.

---

## Speakers

**Aryan Mehrotra** — SSDE @ zop.dev · Maintainer, GoFr

* LinkedIn: <https://www.linkedin.com/in/aryanmehrotra>
* X: <https://x.com/_aryanmehrotra>
* GitHub: <https://github.com/aryanmehrotra>

**Piyush Singh** — SDE-2 @ zop.dev · OpenTofu / Helm modules maintainer

* GitHub: <https://github.com/piyushsingh>

---

## GoFr

* Framework repo: <https://github.com/gofr-dev/gofr>
* Documentation: <https://gofr.dev/docs>
* This example upstream (full commit history):
  [`examples/kafka-tracing-demo`](https://github.com/gofr-dev/gofr/tree/main/examples/kafka-tracing-demo)
* Hosted tracer: <https://tracer.gofr.dev>

---

## References

* W3C Trace Context spec — <https://www.w3.org/TR/trace-context/>
* OpenTelemetry Go — <https://github.com/open-telemetry/opentelemetry-go>
* Messaging semantic conventions — <https://opentelemetry.io/docs/specs/semconv/messaging/>
* `segmentio/kafka-go` (used in the raw baseline) — <https://github.com/segmentio/kafka-go>
