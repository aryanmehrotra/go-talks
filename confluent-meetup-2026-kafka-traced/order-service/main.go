// order-service subscribes to the "orders" topic, performs trivial validation,
// then publishes an enriched event to the "alerts" topic.
//
// The trace context that arrived in the Kafka headers of the inbound message
// is restored on `ctx` automatically by GoFr, and re-injected into the outbound
// message — so a single trace ID survives the bounce through the broker.
// You will not see any tracing code in this file. That is the point.
package main

import (
	"encoding/json"
	"fmt"

	"gofr.dev/pkg/gofr"
)

type order struct {
	OrderID string `json:"orderId"`
	Item    string `json:"item"`
	Qty     int    `json:"qty"`
}

type alert struct {
	OrderID string `json:"orderId"`
	Channel string `json:"channel"`
	Body    string `json:"body"`
}

func main() {
	app := gofr.New()

	app.Subscribe("orders", handleOrder)

	app.Run()
}

func handleOrder(ctx *gofr.Context) error {
	var o order
	if err := ctx.Bind(&o); err != nil {
		ctx.Logger.Errorf("bad order payload: %v", err)
		return nil // don't requeue garbage
	}

	if o.Qty <= 0 {
		ctx.Logger.Warnf("rejecting order %s: invalid qty %d", o.OrderID, o.Qty)
		return nil
	}

	a := alert{
		OrderID: o.OrderID,
		Channel: "slack",
		Body:    fmt.Sprintf("order %s confirmed: %d × %s", o.OrderID, o.Qty, o.Item),
	}

	msg, err := json.Marshal(a)
	if err != nil {
		return err
	}

	ctx.Logger.Infof("processed order %s, fanning out alert", o.OrderID)

	// Same single line. Trace context flows through.
	return ctx.GetPublisher().Publish(ctx, "alerts", msg)
}
