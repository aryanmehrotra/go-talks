// api-gateway is a 1-endpoint HTTP service that accepts an order and publishes
// it to the "orders" Kafka topic. Distributed-tracing context (W3C traceparent)
// is injected into the Kafka message headers automatically by GoFr — there is
// no tracing code in this file, on purpose.
package main

import (
	"encoding/json"

	"gofr.dev/pkg/gofr"
	"kafka-tracing-demo/api-gateway/migrations"
)

type order struct {
	OrderID string `json:"orderId"`
	Item    string `json:"item"`
	Qty     int    `json:"qty"`
}

func main() {
	app := gofr.New()

	app.Migrate(migrations.All())

	app.POST("/order", placeOrder)

	app.Run()
}

func placeOrder(ctx *gofr.Context) (any, error) {
	var o order
	if err := ctx.Bind(&o); err != nil {
		return nil, err
	}

	msg, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}

	// One line. Trace context is injected into Kafka headers under the hood —
	// see pkg/gofr/datasource/pubsub/kafka/tracing.go.
	if err := ctx.GetPublisher().Publish(ctx, "orders", msg); err != nil {
		return nil, err
	}

	return map[string]string{"status": "accepted", "orderId": o.OrderID}, nil
}
