package log

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// tracer os the default global tracer
var tracer trace.Tracer
var propgator propagation.TextMapPropagator

func init() {
	setupOTelSDK(context.Background())
}

func MarshalTrace(ctx context.Context) ([]byte, error) {
	carrier := propagation.MapCarrier{}
	propgator.Inject(ctx, carrier)
	return json.Marshal(carrier)
}

func UnmarshalTrace(data []byte) (context.Context, error) {
	carrier := propagation.MapCarrier{}
	err := json.Unmarshal(data, &carrier)
	return propgator.Extract(context.Background(), carrier), err
}

func discoverServiceName() string {
	svcName := os.Getenv("OTEL_SERVICE_NAME")
	if svcName == "" {
		svcAttrs := os.Getenv("OTEL_SERVICE_ATTRIBUTES")
		for _, b := range strings.Split(svcAttrs, ",") {
			kv := strings.Split(b, "=")
			if len(kv) != 2 {
				continue
			}
			if strings.ToLower(kv[0]) == "service.name" {
				svcName = kv[1]
			}
		}
	}
	if svcName == "" {
		svcName = "unknown"
		Warn().Str("service", svcName).Msg("Service name not set, OTEL_SERVICE_NAME or OTEL_SERVICE_ATTRIBUTES to set it")
	}
	return svcName
}

// setupOTelSDK bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func setupOTelSDK(ctx context.Context) (func(), error) {
	// Configure a new OTLP exporter
	client := otlptracegrpc.NewClient()
	exp, err := otlptrace.New(ctx, client)
	if err != nil {
		return nil, err
	}

	// Create a new tracer provider with a batch span processor and the otlp exporter
	tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exp))

	// Register the global Tracer provider
	otel.SetTracerProvider(tp)

	tracer = otel.Tracer(discoverServiceName())
	propgator = propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})

	// Register the W3C trace context and baggage propagators so data is propagated across services/processes
	otel.SetTextMapPropagator(propgator)

	return func() {
		_ = exp.Shutdown(ctx)
		_ = tp.Shutdown(ctx)
	}, err
}
