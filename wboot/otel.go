package wboot

import (
	"github.com/go-logr/logr"
	"github.com/guoyk93/rg"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
	"net/http"
)

func setupOTEL() (err error) {
	defer rg.Guard(&err)

	// clear error handler
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {}))
	// using zipkin
	otel.SetTracerProvider(trace.NewTracerProvider(
		trace.WithBatcher(rg.Must(zipkin.New("", zipkin.WithLogr(logr.Discard())))),
	))
	// using b3
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
			b3.New(b3.WithInjectEncoding(b3.B3MultipleHeader|b3.B3SingleHeader)),
		),
	)

	// re-initialize http client
	otelhttp.DefaultClient = &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
	http.DefaultClient = otelhttp.DefaultClient

	return
}
