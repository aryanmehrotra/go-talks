// notification-service is the terminal hop. It subscribes to "alerts" and
// pretends to deliver each one (Slack/email mock — just logs).
//
// Look at this trace in Jaeger: it ends here, joined to the original
// HTTP POST /order via Kafka headers carrying the W3C traceparent.
package main

import (
	"gofr.dev/pkg/gofr"
)

type alert struct {
	OrderID string `json:"orderId"`
	Channel string `json:"channel"`
	Body    string `json:"body"`
}

func main() {
	app := gofr.New()

	app.Subscribe("alerts", deliver)

	app.Run()
}

func deliver(ctx *gofr.Context) error {
	var a alert
	if err := ctx.Bind(&a); err != nil {
		ctx.Logger.Errorf("bad alert payload: %v", err)
		return nil
	}

	ctx.Logger.Infof("📣 [%s] %s", a.Channel, a.Body)

	return nil
}
