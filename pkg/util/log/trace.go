package log

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var propgator propagation.TextMapPropagator
var tracer trace.Tracer

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

// InjectTraceData injects trace data, prepending it to the given data.
func InjectTraceData(ctx context.Context, data []byte) ([]byte, error) {
	td, err := MarshalTrace(ctx)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(td)
	buf.Write([]byte("|"))
	buf.Write(data)

	return buf.Bytes(), nil
}

// ExtractTraceData extracts trace data from the given data
func ExtractTraceData(data []byte) (context.Context, []byte, error) {
	parts := bytes.SplitN(data, []byte("|"), 2)
	ctx, err := UnmarshalTrace(parts[0])
	if err != nil {
		return nil, nil, err
	}
	if len(parts) == 2 {
		return ctx, parts[1], nil
	}
	return ctx, nil, fmt.Errorf("no data found")
}

func discoverServiceName() string {
	svcName := os.Getenv("OTEL_SERVICE_NAME")
	if svcName == "" {
		svcAttrs := os.Getenv("OTEL_RESOURCE_ATTRIBUTES")
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
		// opentelemetry.io/docs/specs/otel/configuration/sdk-environment-variables/
		svcName = "unknown"
		Warn().Str("service", svcName).Msg("Service name not set, OTEL_SERVICE_NAME or OTEL_RESOURCE_ATTRIBUTES to set it")
	} else {
		Debug().Str("service", svcName).Msg("Service name discovered")
	}

	return svcName
}

// setupOTelSDK bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func setupOTelSDK(ctx context.Context) (func(), error) {
	// Configure a new OTLP exporter
	client := otlptracehttp.NewClient()
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
