// raw-kafka-baseline/producer
//
// HTTP server with one endpoint that publishes an order to Kafka, with OTel
// trace context manually propagated through Kafka message headers.
//
// Compare against examples/kafka-tracing-demo/api-gateway/main.go (~30 lines).
// This file is ~120 lines and every single block is needed.
package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

// ─── (1) Tracer setup ─────────────────────────────────────────────────────
// ⚠️  GoFr handles this for you. Reading TRACE_EXPORTER, dialing the OTLP
//     collector, picking a sampler — all driven by env vars in GoFr.
func initTracer(ctx context.Context) (*sdktrace.TracerProvider, error) {
	exp, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(getenv("TRACER_URL", "localhost:4317")),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}
	res, err := resource.New(ctx,
		resource.WithAttributes(semconv.ServiceName("raw-kafka-producer")),
	)
	if err != nil {
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	return tp, nil
}

// ─── (2) Manual TextMapCarrier over Kafka headers ─────────────────────────
// ⚠️  GoFr handles this for you. See pkg/gofr/datasource/pubsub/kafka/tracing.go
//     — exactly this carrier, baked in.
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

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tp, err := initTracer(ctx)
	if err != nil {
		log.Fatalf("tracer: %v", err)
	}
	defer tp.Shutdown(context.Background())

	// ─── (3) Kafka writer — connection lifecycle, batching, retries ──────
	// ⚠️  GoFr handles this for you. PUBSUB_BROKER + KAFKA_BATCH_* env vars,
	//     sensible defaults, graceful shutdown.
	w := &kafka.Writer{
		Addr:         kafka.TCP(getenv("KAFKA_BROKER", "localhost:9092")),
		Topic:        "orders",
		Balancer:     &kafka.LeastBytes{},
		BatchSize:    100,
		BatchBytes:   1 << 20,
		BatchTimeout: time.Second,
		RequiredAcks: kafka.RequireAll,
	}
	defer w.Close()

	tracer := otel.Tracer("raw-kafka-producer")

	mux := http.NewServeMux()
	mux.HandleFunc("/order", func(rw http.ResponseWriter, r *http.Request) {
		// ─── (4) Server span — extract any incoming traceparent header ───
		// ⚠️  GoFr handles this for you. The HTTP middleware in
		//     pkg/gofr/http/middleware/tracer.go does it on every request.
		reqCtx := otel.GetTextMapPropagator().Extract(r.Context(),
			propagation.HeaderCarrier(r.Header))
		reqCtx, serverSpan := tracer.Start(reqCtx, "POST /order",
			trace.WithSpanKind(trace.SpanKindServer))
		defer serverSpan.End()

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(rw, err.Error(), 400)
			return
		}

		// ─── (5) Producer span + manual header injection ─────────────────
		// ⚠️  GoFr handles this for you. startPublishSpan() in
		//     pkg/gofr/datasource/pubsub/kafka/tracing.go does exactly this.
		pubCtx, pubSpan := tracer.Start(reqCtx, "kafka-publish",
			trace.WithSpanKind(trace.SpanKindProducer),
			trace.WithAttributes(
				attribute.String("messaging.system", "kafka"),
				attribute.String("messaging.destination.name", "orders"),
				attribute.String("messaging.operation", "publish"),
			),
		)

		var hdrs headerCarrier
		otel.GetTextMapPropagator().Inject(pubCtx, &hdrs)

		payload, _ := json.Marshal(body)
		if err := w.WriteMessages(pubCtx, kafka.Message{
			Key:     []byte(""),
			Value:   payload,
			Headers: []kafka.Header(hdrs),
		}); err != nil {
			pubSpan.RecordError(err)
			pubSpan.End()
			http.Error(rw, err.Error(), 500)
			return
		}
		pubSpan.End()

		rw.WriteHeader(202)
		_, _ = rw.Write([]byte(`{"status":"accepted"}` + "\n"))
	})

	// ─── (6) Graceful shutdown ───────────────────────────────────────────
	// ⚠️  GoFr handles this for you. app.Run() drains everything on SIGTERM.
	srv := &http.Server{Addr: ":18000", Handler: mux}
	go func() {
		log.Println("raw producer listening on :18000")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	_ = srv.Shutdown(context.Background())
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
