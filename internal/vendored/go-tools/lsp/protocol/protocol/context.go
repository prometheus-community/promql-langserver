package protocol

import (
	"context"
	"fmt"

	"github.com/prometheus-community/promql-langserver/internal/vendored/go-tools/telemetry"
	"github.com/prometheus-community/promql-langserver/internal/vendored/go-tools/telemetry/export"
	"github.com/prometheus-community/promql-langserver/internal/vendored/go-tools/xcontext"
)

func init() {
	export.AddExporters(logExporter{})
}

type contextKey int

const (
	clientKey = contextKey(iota)
)

func WithClient(ctx context.Context, client Client) context.Context {
	return context.WithValue(ctx, clientKey, client)
}

// logExporter sends the log event back to the client if there is one stored on the
// context.
type logExporter struct{}

func (logExporter) StartSpan(context.Context, *telemetry.Span)   {}
func (logExporter) FinishSpan(context.Context, *telemetry.Span)  {}
func (logExporter) Metric(context.Context, telemetry.MetricData) {}
func (logExporter) Flush()                                       {}

func (logExporter) Log(ctx context.Context, event telemetry.Event) {
	client, ok := ctx.Value(clientKey).(Client)
	if !ok {
		return
	}
	msg := &LogMessageParams{Type: Info, Message: fmt.Sprint(event)}
	if event.Error != nil {
		msg.Type = Error
	}
	go client.LogMessage(xcontext.Detach(ctx), msg)
}
