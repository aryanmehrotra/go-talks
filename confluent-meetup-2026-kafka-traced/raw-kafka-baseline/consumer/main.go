// raw-kafka-baseline/consumer
//
// Consumer-group reader for the "orders" topic, with manual OTel header
// extraction, span-link creation, and offset commit-on-success.
//
// Compare against examples/kafka-tracing-demo/order-service/main.go (~25 lines).
// This file is ~150 lines.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// ─── (1) Tracer setup ─────────────────────────────────────────────────────
// ⚠️  GoFr handles this for you.
func initTracer(ctx context.Context) (*sdktrace.TracerProvider, error) {
	exp, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(getenv("TRACER_URL", "localhost:4317")),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	return tp, nil
}

// ─── (2) Same TextMapCarrier as the producer ──────────────────────────────
// ⚠️  GoFr handles this for you.
type headerCarrier []kafka.Header

func (c headerCarrier) Get(key string) string {
	for _, h := range c {
		if h.Key == key {
			return string(h.Value)
		}
	}
	return ""
}
func (c *headerCarrier) Set(key, value string) {
	for i, h := range *c {
		if h.Key == key {
			(*c)[i].Value = []byte(value)
			return
		}
	}
	*c = append(*c, kafka.Header{Key: key, Value: []byte(value)})
}
func (c headerCarrier) Keys() []string {
	keys := make([]string, 0, len(c))
	for _, h := range c {
		keys = append(keys, h.Key)
	}
	return keys
}

// ─── (3) Extract a span LINK (not parent) from headers ────────────────────
// Kafka is async + fan-out: a single message can be consumed by N subscribers.
// Modeling each as a child of the producer would be wrong; span LINKS are the
// correct OTel primitive.
//
// ⚠️  GoFr handles this for you. See extractTraceLinks() in
//     pkg/gofr/datasource/pubsub/kafka/tracing.go.
func linksFromHeaders(headers []kafka.Header) []trace.Link {
	carrier := headerCarrier(headers)
	ctx := otel.GetTextMapPropagator().Extract(context.Background(), &carrier)
	sc := trace.SpanContextFromContext(ctx)
	if !sc.IsValid() {
		return nil
	}
	return []trace.Link{{SpanContext: sc}}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tp, err := initTracer(ctx)
	if err != nil {
		log.Fatalf("tracer: %v", err)
	}
	defer tp.Shutdown(context.Background())

	// ─── (4) Consumer-group reader ────────────────────────────────────────
	// You must pick a GroupID for the broker to do partition assignment +
	// rebalancing, decide commit semantics (we want commit-on-success, so
	// we DON'T set CommitInterval), and survive rebalance events.
	//
	// ⚠️  GoFr handles all of this for you. CONSUMER_ID env var is the
	//     GroupID; the framework owns the read loop and commit semantics.
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{getenv("KAFKA_BROKER", "localhost:9092")},
		GroupID:        "raw-baseline-orders",
		Topic:          "orders",
		MinBytes:       1,
		MaxBytes:       10 << 20,
		CommitInterval: 0, // explicit commit only
	})
	defer r.Close()

	tracer := otel.Tracer("raw-kafka-consumer")

	// ─── (5) Read loop with manual offset commit and graceful shutdown ───
	// ⚠️  GoFr handles this for you. Subscriber loop + SIGTERM drain are
	//     in pkg/gofr/datasource/pubsub/kafka/kafka.go (Subscribe + commit).
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	go func() { <-stop; cancel() }()

	log.Println("raw consumer joined group raw-baseline-orders")
	for {
		msg, err := r.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, context.Canceled) {
				return
			}
			log.Printf("fetch: %v", err)
			continue
		}

		// ─── (6) Build a consumer span LINKED to the producer span ───────
		// ⚠️  GoFr handles this for you. startSubscribeSpan().
		links := linksFromHeaders(msg.Headers)
		opts := []trace.SpanStartOption{
			trace.WithSpanKind(trace.SpanKindConsumer),
			trace.WithAttributes(
				attribute.String("messaging.system", "kafka"),
				attribute.String("messaging.destination.name", "orders"),
				attribute.String("messaging.operation", "receive"),
			),
		}
		if len(links) > 0 {
			opts = append(opts, trace.WithLinks(links...))
		}
		_, span := tracer.Start(ctx, "kafka-subscribe", opts...)

		// ─── (7) Business logic ─────────────────────────────────────────
		var order map[string]any
		if err := json.Unmarshal(msg.Value, &order); err != nil {
			log.Printf("bad payload: %v", err)
			span.End()
			// You still must commit, otherwise this poison message stalls
			// the partition forever.
			_ = r.CommitMessages(ctx, msg)
			continue
		}
		log.Printf("processed order: %v", order)

		// ─── (8) Commit-on-success, manually ─────────────────────────────
		// ⚠️  GoFr handles this for you. The framework commits only when the
		//     handler returns nil; non-nil error leaves the offset, so the
		//     message is redelivered after rebalance.
		if err := r.CommitMessages(ctx, msg); err != nil {
			log.Printf("commit: %v", err)
		}
		span.End()
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
