// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package httpstats // import "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

import (
	"time"

	"github.com/gin-gonic/gin"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const ScopeName = "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

const (
	ReadBytesKey  = attribute.Key("http.read_bytes")  // if anything was read from the request body, the total number of bytes read
	ReadErrorKey  = attribute.Key("http.read_error")  // If an error occurred while reading a request, the string of the error (io.EOF is not recorded)
	WroteBytesKey = attribute.Key("http.wrote_bytes") // if anything was written to the response writer, the total number of bytes written
	WriteErrorKey = attribute.Key("http.write_error") // if an error occurred while writing a reply, the string of the error (io.EOF is not recorded)
)

// Server HTTP metrics.
const (
	serverRequestSize  = "http.server.request.size"  // Incoming request bytes total
	serverResponseSize = "http.server.response.size" // Incoming response bytes total
	serverDuration     = "http.server.duration"      // Incoming end to end duration, milliseconds
)

// middleware is an http middleware which wraps the next handler in a span.
type middleware struct {
	meter metric.Meter

	serverLatencyMeasure metric.Float64Histogram
}

// NewMiddleware returns a tracing and metrics instrumentation middleware.
// The handler returned by the middleware wraps a handler
// in a span named after the operation and enriches it with metrics.
func NewMiddleware() gin.HandlerFunc {
	h := middleware{
		meter: otel.GetMeterProvider().Meter(
			ScopeName,
			metric.WithInstrumentationVersion("0.53.0"),
		),
	}

	h.createMeasures()

	return func(c *gin.Context) {
		h.serveHTTP(c)
	}
}

func handleErr(err error) {
	if err != nil {
		otel.Handle(err)
	}
}

func (h *middleware) createMeasures() {
	var err error

	h.serverLatencyMeasure, err = h.meter.Float64Histogram(
		serverDuration,
		metric.WithUnit("ms"),
		metric.WithDescription("Measures the duration of inbound HTTP requests."),
	)
	handleErr(err)
}

// serveHTTP sets up tracing and calls the given next http.Handler with the span
// context injected into the request context.
func (h *middleware) serveHTTP(c *gin.Context) {
	ctx := c.Request.Context()
	requestStartTime := time.Now()

	// Add metrics
	attributes := []attribute.KeyValue{}
	attributes = append(attributes, serverRequestMetrics(c)...)
	route := c.FullPath()
	if route == "" {
		attributes = append(attributes, attribute.Key("http.route").String("not-found"))
	} else {
		attributes = append(attributes, attribute.Key("http.route").String(route))
	}

	c.Next()

	attributes = append(attributes, attribute.Key("http.status_code").Int(c.Writer.Status()))
	o := metric.WithAttributeSet(attribute.NewSet(attributes...))
	elapsedTime := float64(time.Since(requestStartTime)) / float64(time.Millisecond)
	h.serverLatencyMeasure.Record(ctx, elapsedTime, o)
}

func serverRequestMetrics(c *gin.Context) []attribute.KeyValue {
	attrs := make([]attribute.KeyValue, 0, 6)
	attrs = append(attrs, attribute.Key("http.method").String(c.Request.Method))
	attrs = append(attrs, attribute.Key("http.scheme").String(c.Request.URL.Scheme))
	attrs = append(attrs, attribute.Key("net.host.name").String(c.Request.Host))
	attrs = append(attrs, attribute.Key("net.protocol.name").String(c.Request.URL.Scheme))
	return attrs
}
